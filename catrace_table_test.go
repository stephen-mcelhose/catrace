package catrace_test

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func nearlyEqualSlice(t *testing.T, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("slice length: got %d want %d", len(got), len(want))
	}
	for i := range got {
		if math.Abs(got[i]-want[i]) > tol {
			t.Fatalf("[%d]: got %g want %g (tol %g)", i, got[i], want[i], tol)
		}
	}
}

func mustKernel(t *testing.T, data []float64, n int) *catrace.Kernel {
	t.Helper()
	k, err := catrace.NewKernel(mat.NewDense(n, n, data), nil)
	if err != nil {
		t.Fatalf("mustKernel: %v", err)
	}
	return k
}

// ergodic2x2 returns [[0.7,0.3],[0.4,0.6]]; stationary = [4/7, 3/7].
func ergodic2x2(t *testing.T) *catrace.Kernel {
	t.Helper()
	return mustKernel(t, []float64{0.7, 0.3, 0.4, 0.6}, 2)
}

// pathGraph returns the 3-state path A-B-C; period 2, Stationary() will not converge.
func pathGraph(t *testing.T) *catrace.Kernel {
	t.Helper()
	return mustKernel(t, []float64{
		0, 1, 0,
		0.5, 0, 0.5,
		0, 1, 0,
	}, 3)
}

// ── kernel.go ─────────────────────────────────────────────────────────────────

func TestNewKernel(t *testing.T) {
	cases := []struct {
		name    string
		p       *mat.Dense
		names   []string
		wantErr bool
	}{
		{
			name: "accepts well-formed row-stochastic matrix",
			p:    mat.NewDense(2, 2, []float64{0.7, 0.3, 0.4, 0.6}),
		},
		{
			name: "single-state kernel is trivially valid",
			p:    mat.NewDense(1, 1, []float64{1.0}),
		},
		{
			name:    "nil matrix",
			wantErr: true,
		},
		{
			name:    "rejects non-square matrix",
			p:       mat.NewDense(2, 3, []float64{0.5, 0.3, 0.2, 0.4, 0.4, 0.2}),
			wantErr: true,
		},
		{
			name:    "rejects mismatched state name count",
			p:       mat.NewDense(2, 2, []float64{0.7, 0.3, 0.4, 0.6}),
			names:   []string{"a", "b", "c"},
			wantErr: true,
		},
		{
			name:    "negative entry",
			p:       mat.NewDense(2, 2, []float64{-0.1, 1.1, 0.4, 0.6}),
			wantErr: true,
		},
		{
			name:    "row does not sum to 1",
			p:       mat.NewDense(2, 2, []float64{0.5, 0.3, 0.4, 0.6}),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k, err := catrace.NewKernel(tc.p, tc.names)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if !tc.wantErr && k == nil {
				t.Fatal("expected non-nil kernel")
			}
		})
	}
}

func TestNewRectKernel(t *testing.T) {
	cases := []struct {
		name              string
		p                 *mat.Dense
		rowNames, colNames []string
		wantErr           bool
	}{
		{
			name: "accepts rectangular row-stochastic matrix",
			p:    mat.NewDense(2, 3, []float64{0.5, 0.3, 0.2, 0.1, 0.6, 0.3}),
		},
		{
			name:    "nil matrix",
			wantErr: true,
		},
		{
			name:     "rejects mismatched row name count",
			p:        mat.NewDense(2, 3, []float64{0.5, 0.3, 0.2, 0.1, 0.6, 0.3}),
			rowNames: []string{"a", "b", "c"},
			wantErr:  true,
		},
		{
			name:     "rejects mismatched column name count",
			p:        mat.NewDense(2, 3, []float64{0.5, 0.3, 0.2, 0.1, 0.6, 0.3}),
			colNames: []string{"x"},
			wantErr:  true,
		},
		{
			name:    "negative entry",
			p:       mat.NewDense(2, 2, []float64{-0.1, 1.1, 0.4, 0.6}),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k, err := catrace.NewRectKernel(tc.p, tc.rowNames, tc.colNames)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if !tc.wantErr && k == nil {
				t.Fatal("expected non-nil RectKernel")
			}
		})
	}
}

func TestKernel_Clone(t *testing.T) {
	k := ergodic2x2(t)

	t.Run("clone has equal entries", func(t *testing.T) {
		c := k.Clone()
		if c == nil {
			t.Fatal("Clone returned nil")
		}
		nearlyEqualSlice(t, []float64{c.P.At(0, 0), c.P.At(0, 1)}, []float64{0.7, 0.3}, 1e-12)
		nearlyEqualSlice(t, []float64{c.P.At(1, 0), c.P.At(1, 1)}, []float64{0.4, 0.6}, 1e-12)
	})

	t.Run("mutating clone does not affect original", func(t *testing.T) {
		c := k.Clone()
		c.P.Set(0, 0, 0.0)
		if math.Abs(k.P.At(0, 0)-0.7) > 1e-12 {
			t.Fatal("mutating clone affected original")
		}
	})

	t.Run("nil kernel clones to nil", func(t *testing.T) {
		var nilK *catrace.Kernel
		if nilK.Clone() != nil {
			t.Fatal("Clone of nil should return nil")
		}
	})
}

func TestKernel_NumStates(t *testing.T) {
	t.Run("nil kernel reports zero states", func(t *testing.T) {
		var nilK *catrace.Kernel
		if nilK.NumStates() != 0 {
			t.Fatal("nil kernel should have 0 states")
		}
	})
	t.Run("reports correct count for 2-state kernel", func(t *testing.T) {
		if ergodic2x2(t).NumStates() != 2 {
			t.Fatal()
		}
	})
	t.Run("reports correct count for 4-state kernel", func(t *testing.T) {
		k := mustKernel(t, []float64{
			0.60, 0.20, 0.10, 0.10,
			0.15, 0.55, 0.15, 0.15,
			0.20, 0.20, 0.40, 0.20,
			0.10, 0.20, 0.20, 0.50,
		}, 4)
		if k.NumStates() != 4 {
			t.Fatal()
		}
	})
}

func TestKernel_NormalizeRows(t *testing.T) {
	cases := []struct {
		name    string
		data    []float64
		n       int
		wantErr bool
		checkFn func(t *testing.T, p *mat.Dense)
	}{
		{
			name: "already-normalized rows are unchanged",
			data: []float64{0.7, 0.3, 0.4, 0.6},
			n:    2,
		},
		{
			name: "rescales raw counts to probabilities",
			data: []float64{3, 1, 1, 1},
			n:    2,
			checkFn: func(t *testing.T, p *mat.Dense) {
				t.Helper()
				nearlyEqualSlice(t, []float64{p.At(0, 0), p.At(0, 1)}, []float64{0.75, 0.25}, 1e-12)
				nearlyEqualSlice(t, []float64{p.At(1, 0), p.At(1, 1)}, []float64{0.5, 0.5}, 1e-12)
			},
		},
		{
			name:    "significant negative",
			data:    []float64{-0.1, 1.1, 0.4, 0.6},
			n:       2,
			wantErr: true,
		},
		{
			name:    "all-zero row cannot be normalized",
			data:    []float64{0, 0, 0.5, 0.5},
			n:       2,
			wantErr: true,
		},
		{
			name: "clamps near-zero negative entries to zero",
			data: []float64{1.0, -1e-15, 0.5, 0.5},
			n:    2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := mat.NewDense(tc.n, tc.n, tc.data)
			// Construct directly: NormalizeRows must work on kernels
			// with unnormalized rows, which NewKernel would reject.
			k := &catrace.Kernel{P: p}
			err := k.NormalizeRows(1e-9)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if err == nil {
				n, m := p.Dims()
				for i := 0; i < n; i++ {
					sum := 0.0
					for j := 0; j < m; j++ {
						sum += p.At(i, j)
					}
					if math.Abs(sum-1.0) > 1e-9 {
						t.Fatalf("row %d sums to %g after NormalizeRows", i, sum)
					}
				}
				if tc.checkFn != nil {
					tc.checkFn(t, p)
				}
			}
		})
	}
}

func TestKernel_Multiply(t *testing.T) {
	k2 := ergodic2x2(t)
	id2 := mustKernel(t, []float64{1, 0, 0, 1}, 2)
	k3 := mustKernel(t, []float64{
		1.0 / 3, 1.0 / 3, 1.0 / 3,
		1.0 / 3, 1.0 / 3, 1.0 / 3,
		1.0 / 3, 1.0 / 3, 1.0 / 3,
	}, 3)

	t.Run("multiplying by identity leaves kernel unchanged", func(t *testing.T) {
		res, err := id2.Multiply(k2)
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, []float64{res.P.At(0, 0), res.P.At(0, 1)}, []float64{0.7, 0.3}, 1e-12)
		nearlyEqualSlice(t, []float64{res.P.At(1, 0), res.P.At(1, 1)}, []float64{0.4, 0.6}, 1e-12)
	})

	t.Run("right-multiplying by identity leaves kernel unchanged", func(t *testing.T) {
		res, err := k2.Multiply(id2)
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, []float64{res.P.At(0, 0), res.P.At(0, 1)}, []float64{0.7, 0.3}, 1e-12)
	})

	t.Run("rejects kernels with incompatible dimensions", func(t *testing.T) {
		_, err := k2.Multiply(k3)
		if err == nil {
			t.Fatal("expected error for dimension mismatch")
		}
	})

	t.Run("rejects nil operand", func(t *testing.T) {
		_, err := k2.Multiply(nil)
		if err == nil {
			t.Fatal("expected error for nil operand")
		}
	})
}

func TestKernel_LeftAction(t *testing.T) {
	k := ergodic2x2(t)

	t.Run("point mass on state i advances to row i of kernel", func(t *testing.T) {
		got, err := k.LeftAction([]float64{1.0, 0.0})
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, got, []float64{0.7, 0.3}, 1e-12)
	})

	t.Run("evolved distribution sums to 1", func(t *testing.T) {
		got, err := k.LeftAction([]float64{0.5, 0.5})
		if err != nil {
			t.Fatal(err)
		}
		sum := got[0] + got[1]
		if math.Abs(sum-1.0) > 1e-12 {
			t.Fatalf("result sums to %g, want 1", sum)
		}
	})

	t.Run("rejects distribution with wrong number of states", func(t *testing.T) {
		_, err := k.LeftAction([]float64{0.5})
		if err == nil {
			t.Fatal("expected error for wrong distribution length")
		}
	})

	t.Run("nil kernel", func(t *testing.T) {
		var nilK *catrace.Kernel
		_, err := nilK.LeftAction([]float64{0.5, 0.5})
		if err == nil {
			t.Fatal("expected error for nil kernel")
		}
	})
}

// ── graph.go ──────────────────────────────────────────────────────────────────

func TestNewRandomWalkKernel(t *testing.T) {
	cases := []struct {
		name    string
		adj     *mat.Dense
		names   []string
		wantErr bool
		checkFn func(t *testing.T, k *catrace.Kernel)
	}{
		{
			name: "each neighbor receives weight proportional to edge weight",
			adj: mat.NewDense(4, 4, []float64{
				0, 1, 1, 1,
				1, 0, 1, 0,
				1, 1, 0, 0,
				1, 0, 0, 0,
			}),
			names: []string{"A", "B", "C", "D"},
			checkFn: func(t *testing.T, k *catrace.Kernel) {
				t.Helper()
				// A has degree 3: each neighbor gets 1/3
				nearlyEqualSlice(t, []float64{k.P.At(0, 1), k.P.At(0, 2), k.P.At(0, 3)},
					[]float64{1.0 / 3, 1.0 / 3, 1.0 / 3}, 1e-12)
				// D has degree 1: always goes to A
				if math.Abs(k.P.At(3, 0)-1.0) > 1e-12 {
					t.Fatalf("P(D→A) = %g, want 1", k.P.At(3, 0))
				}
			},
		},
		{
			name:    "nil matrix",
			wantErr: true,
		},
		{
			name:    "rejects non-square adjacency matrix",
			adj:     mat.NewDense(2, 3, []float64{1, 1, 0, 1, 0, 1}),
			wantErr: true,
		},
		{
			name:    "negative entry",
			adj:     mat.NewDense(2, 2, []float64{0, -1, 1, 0}),
			wantErr: true,
		},
		{
			name:    "isolated node with no outgoing edges",
			adj:     mat.NewDense(2, 2, []float64{0, 1, 0, 0}),
			wantErr: true,
		},
		{
			name: "directed edges are one-way constraints",
			adj: mat.NewDense(3, 3, []float64{
				0, 0.5, 0.5,
				0, 0, 1,
				1, 0, 0,
			}),
			checkFn: func(t *testing.T, k *catrace.Kernel) {
				t.Helper()
				// row 0: each non-zero gets 0.5/1.0 = 0.5
				nearlyEqualSlice(t, []float64{k.P.At(0, 1), k.P.At(0, 2)}, []float64{0.5, 0.5}, 1e-12)
				// row 1: all mass on state 2
				if math.Abs(k.P.At(1, 2)-1.0) > 1e-12 {
					t.Fatalf("P(1→2) = %g, want 1", k.P.At(1, 2))
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k, err := catrace.NewRandomWalkKernel(tc.adj, tc.names)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if err == nil && tc.checkFn != nil {
				tc.checkFn(t, k)
			}
		})
	}
}

// ── agent.go ──────────────────────────────────────────────────────────────────

func TestAgent_Validate(t *testing.T) {
	d := mat.NewDense(2, 2, []float64{0.9, 0.1, 0.2, 0.8})
	a := mat.NewDense(2, 2, []float64{0.7, 0.3, 0.4, 0.6})
	p := mat.NewDense(2, 2, []float64{0.8, 0.2, 0.1, 0.9})

	cases := []struct {
		name    string
		agent   *catrace.Agent
		wantErr bool
	}{
		{
			name:  "accepts agent with consistent dimensions and stochastic maps",
			agent: &catrace.Agent{D: d, A: a, P: p},
		},
		{
			name:    "nil agent",
			agent:   nil,
			wantErr: true,
		},
		{
			name:    "missing decision kernel",
			agent:   &catrace.Agent{D: nil, A: a, P: p},
			wantErr: true,
		},
		{
			name: "action space mismatch between D and A",
			agent: &catrace.Agent{
				D: mat.NewDense(2, 3, []float64{0.5, 0.3, 0.2, 0.1, 0.6, 0.3}),
				A: a, P: p,
			},
			wantErr: true,
		},
		{
			name: "world space mismatch between A and P",
			agent: &catrace.Agent{
				D: d,
				A: mat.NewDense(2, 3, []float64{0.5, 0.3, 0.2, 0.1, 0.6, 0.3}),
				P: p,
			},
			wantErr: true,
		},
		{
			name: "experience space mismatch between D and P",
			agent: &catrace.Agent{
				D: mat.NewDense(3, 2, []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5}),
				A: a,
				P: p,
			},
			wantErr: true,
		},
		{
			name: "decision kernel with rows not summing to 1",
			agent: &catrace.Agent{
				D: mat.NewDense(2, 2, []float64{0.5, 0.3, 0.4, 0.6}),
				A: a, P: p,
			},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.agent.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
		})
	}
}

func TestAgent_DerivedKernels(t *testing.T) {
	agent := &catrace.Agent{
		D: mat.NewDense(2, 2, []float64{0.9, 0.1, 0.2, 0.8}),
		A: mat.NewDense(2, 2, []float64{0.7, 0.3, 0.4, 0.6}),
		P: mat.NewDense(2, 2, []float64{0.8, 0.2, 0.1, 0.9}),
	}
	// Q = D·A·P = [[0.569,0.431],[0.422,0.578]] (verified by ExampleAgent_QualiaKernel)
	// S = A·P·D = [[0.613,0.387],[0.466,0.534]]
	// W = P·D·A = [[0.628,0.372],[0.481,0.519]]

	t.Run("QualiaKernel computes D·A·P correctly", func(t *testing.T) {
		q, err := agent.QualiaKernel()
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, []float64{q.P.At(0, 0), q.P.At(0, 1)}, []float64{0.569, 0.431}, 1e-3)
		nearlyEqualSlice(t, []float64{q.P.At(1, 0), q.P.At(1, 1)}, []float64{0.422, 0.578}, 1e-3)
	})

	t.Run("StrategyKernel computes A·P·D correctly", func(t *testing.T) {
		s, err := agent.StrategyKernel()
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, []float64{s.P.At(0, 0), s.P.At(0, 1)}, []float64{0.613, 0.387}, 1e-3)
		nearlyEqualSlice(t, []float64{s.P.At(1, 0), s.P.At(1, 1)}, []float64{0.466, 0.534}, 1e-3)
	})

	t.Run("WorldKernel computes P·D·A correctly", func(t *testing.T) {
		w, err := agent.WorldKernel()
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, []float64{w.P.At(0, 0), w.P.At(0, 1)}, []float64{0.628, 0.372}, 1e-3)
		nearlyEqualSlice(t, []float64{w.P.At(1, 0), w.P.At(1, 1)}, []float64{0.481, 0.519}, 1e-3)
	})

	t.Run("all three kernels require a valid agent", func(t *testing.T) {
		bad := &catrace.Agent{}
		if _, err := bad.QualiaKernel(); err == nil {
			t.Fatal("expected error from invalid agent")
		}
		if _, err := bad.StrategyKernel(); err == nil {
			t.Fatal("expected error from invalid agent")
		}
		if _, err := bad.WorldKernel(); err == nil {
			t.Fatal("expected error from invalid agent")
		}
	})
}

// ── trace.go ──────────────────────────────────────────────────────────────────

func TestKernel_Trace(t *testing.T) {
	// 3-state kernel from ExampleKernel_Trace; trace on {0,1} = [[0.68,0.32],[0.36,0.64]]
	parent := mustKernel(t, []float64{
		0.6, 0.3, 0.1,
		0.2, 0.6, 0.2,
		0.4, 0.1, 0.5,
	}, 3)

	cases := []struct {
		name    string
		k       *catrace.Kernel
		subset  []int
		wantErr bool
		want    [][]float64
	}{
		{
			name:   "computes excursion-corrected transitions on subset",
			k:      parent,
			subset: []int{0, 1},
			want:   [][]float64{{0.68, 0.32}, {0.36, 0.64}},
		},
		{
			name:   "trace on full state space equals original kernel",
			k:      parent,
			subset: []int{0, 1, 2},
			want:   [][]float64{{0.6, 0.3, 0.1}, {0.2, 0.6, 0.2}, {0.4, 0.1, 0.5}},
		},
		{
			name:    "nil kernel",
			wantErr: true,
		},
		{
			name:    "empty subset",
			k:       parent,
			subset:  []int{},
			wantErr: true,
		},
		{
			name:    "out-of-range index",
			k:       parent,
			subset:  []int{0, 5},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr, err := tc.k.Trace(tc.subset, 1e-12)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if err == nil {
				for i, row := range tc.want {
					got := make([]float64, len(row))
					for j := range row {
						got[j] = tr.P.At(i, j)
					}
					nearlyEqualSlice(t, got, row, 1e-10)
				}
			}
		})
	}
}

func TestKernel_IsTraceOf(t *testing.T) {
	parent := mustKernel(t, []float64{
		0.6, 0.3, 0.1,
		0.2, 0.6, 0.2,
		0.4, 0.1, 0.5,
	}, 3)
	tr, _ := parent.Trace([]int{0, 1}, 1e-12)

	t.Run("returns true when kernel matches computed trace", func(t *testing.T) {
		ok, err := tr.IsTraceOf(parent, []int{0, 1}, 1e-12)
		if err != nil || !ok {
			t.Fatalf("expected true; ok=%v err=%v", ok, err)
		}
	})

	t.Run("returns false when kernel differs from computed trace", func(t *testing.T) {
		other := ergodic2x2(t)
		ok, err := other.IsTraceOf(parent, []int{0, 1}, 1e-12)
		if err != nil || ok {
			t.Fatalf("expected false; ok=%v err=%v", ok, err)
		}
	})

	t.Run("nil parent", func(t *testing.T) {
		_, err := tr.IsTraceOf(nil, []int{0, 1}, 1e-12)
		if err == nil {
			t.Fatal("expected error for nil parent")
		}
	})
}

func TestRestrictDistribution(t *testing.T) {
	cases := []struct {
		name    string
		pi      []float64
		subset  []int
		wantErr bool
		want    []float64
	}{
		{
			name:   "restricts and renormalizes distribution to subset",
			pi:     []float64{0.3, 0.5, 0.2},
			subset: []int{0, 2},
			want:   []float64{0.6, 0.4}, // [0.3,0.2] → normalized → [3/5,2/5]
		},
		{
			name:   "single-element subset always yields probability 1",
			pi:     []float64{0.4, 0.6},
			subset: []int{1},
			want:   []float64{1.0},
		},
		{
			name:    "empty subset",
			pi:      []float64{0.5, 0.5},
			subset:  []int{},
			wantErr: true,
		},
		{
			name:    "index out of range",
			pi:      []float64{0.5, 0.5},
			subset:  []int{0, 5},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := catrace.RestrictDistribution(tc.pi, tc.subset, 1e-12)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if err == nil {
				nearlyEqualSlice(t, got, tc.want, 1e-12)
			}
		})
	}
}

func TestCanonicalSubset(t *testing.T) {
	cases := []struct {
		name   string
		input  []int
		want   []int
	}{
		{"already sorted",  []int{0, 1, 2}, []int{0, 1, 2}},
		{"unsorted",        []int{2, 0, 1}, []int{0, 1, 2}},
		{"duplicates",      []int{1, 1, 0}, []int{0, 1}},
		{"single element",  []int{3},       []int{3}},
		{"empty",           []int{},        []int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := catrace.CanonicalSubset(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

// ── stationary.go ─────────────────────────────────────────────────────────────

func TestKernel_Stationary(t *testing.T) {
	t.Run("converges to unique stationary distribution for ergodic chain", func(t *testing.T) {
		pi, err := ergodic2x2(t).Stationary(1e-12, 5000)
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, pi, []float64{4.0 / 7, 3.0 / 7}, 1e-8)
	})

	t.Run("uniform chain has uniform stationary distribution", func(t *testing.T) {
		k, _ := catrace.NewKernel(mat.NewDense(2, 2, []float64{0.5, 0.5, 0.5, 0.5}), nil)
		pi, err := k.Stationary(1e-12, 5000)
		if err != nil {
			t.Fatal(err)
		}
		nearlyEqualSlice(t, pi, []float64{0.5, 0.5}, 1e-12)
	})

	t.Run("nil kernel", func(t *testing.T) {
		var nilK *catrace.Kernel
		_, err := nilK.Stationary(1e-12, 5000)
		if err == nil {
			t.Fatal("expected error for nil kernel")
		}
	})

	t.Run("period-2 path graph does not converge", func(t *testing.T) {
		_, err := pathGraph(t).Stationary(1e-12, 5000)
		if err == nil {
			t.Fatal("expected non-convergence error for period-2 chain")
		}
	})

	t.Run("reports non-convergence when iteration budget is exhausted", func(t *testing.T) {
		_, err := ergodic2x2(t).Stationary(1e-12, 1)
		if err == nil {
			t.Fatal("expected error: 1 iteration is insufficient to converge")
		}
	})

}

func TestKernel_EntropyRate(t *testing.T) {
	t.Run("deterministic chain has zero entropy", func(t *testing.T) {
		// Identity matrix: always stay; H = 0.
		k := mustKernel(t, []float64{1, 0, 0, 1}, 2)
		h, err := k.EntropyRate(2)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(h) > 1e-12 {
			t.Fatalf("got %g, want 0", h)
		}
	})

	t.Run("uniform 2x2 has 1 bit", func(t *testing.T) {
		k, _ := catrace.NewKernel(mat.NewDense(2, 2, []float64{0.5, 0.5, 0.5, 0.5}), nil)
		h, err := k.EntropyRate(2)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(h-1.0) > 1e-10 {
			t.Fatalf("got %g, want 1.0 bit", h)
		}
	})

	t.Run("ergodic 2x2 entropy is positive and < 1 bit", func(t *testing.T) {
		h, err := ergodic2x2(t).EntropyRate(2)
		if err != nil {
			t.Fatal(err)
		}
		if h <= 0 || h >= 1 {
			t.Fatalf("got %g, want in (0,1)", h)
		}
	})

	t.Run("base 0 is invalid", func(t *testing.T) {
		_, err := ergodic2x2(t).EntropyRate(0)
		if err == nil {
			t.Fatal("expected error for base=0")
		}
	})

	t.Run("base 1 is invalid", func(t *testing.T) {
		_, err := ergodic2x2(t).EntropyRate(1)
		if err == nil {
			t.Fatal("expected error for base=1")
		}
	})
}

// ── analysis.go ───────────────────────────────────────────────────────────────

func TestKernel_Classes(t *testing.T) {
	t.Run("nil kernel", func(t *testing.T) {
		var nilK *catrace.Kernel
		_, err := nilK.Classes(1e-9)
		if err == nil {
			t.Fatal("expected error for nil kernel")
		}
	})

	t.Run("ergodic chain forms one aperiodic recurrent class", func(t *testing.T) {
		d, err := ergodic2x2(t).Classes(1e-9)
		if err != nil {
			t.Fatal(err)
		}
		if len(d.SCCs) != 1 {
			t.Fatalf("SCCs: got %d want 1", len(d.SCCs))
		}
		if len(d.Recurrent) != 1 {
			t.Fatalf("Recurrent: got %d want 1", len(d.Recurrent))
		}
		if len(d.Transient) != 0 {
			t.Fatalf("Transient: got %v want []", d.Transient)
		}
		// period must be 1 (self-loops in both rows)
		if d.Periods[0] != 1 {
			t.Fatalf("period: got %d want 1", d.Periods[0])
		}
	})

	t.Run("bipartite path graph forms one period-2 recurrent class", func(t *testing.T) {
		d, err := pathGraph(t).Classes(1e-9)
		if err != nil {
			t.Fatal(err)
		}
		if len(d.Recurrent) != 1 {
			t.Fatalf("Recurrent: got %d want 1", len(d.Recurrent))
		}
		if len(d.Transient) != 0 {
			t.Fatalf("Transient: got %v want []", d.Transient)
		}
		// find the SCC index for the recurrent class
		var idx int
		for i, scc := range d.SCCs {
			for _, s := range scc {
				if s == 0 {
					idx = i
				}
			}
		}
		if d.Periods[idx] != 2 {
			t.Fatalf("period: got %d want 2", d.Periods[idx])
		}
	})

	t.Run("two absorbing states form two separate recurrent classes", func(t *testing.T) {
		k := mustKernel(t, []float64{1, 0, 0, 1}, 2)
		d, err := k.Classes(1e-9)
		if err != nil {
			t.Fatal(err)
		}
		if len(d.Recurrent) != 2 {
			t.Fatalf("Recurrent: got %d want 2", len(d.Recurrent))
		}
		if len(d.Transient) != 0 {
			t.Fatalf("Transient: got %v want []", d.Transient)
		}
	})

	t.Run("state with no path back to itself is classified transient", func(t *testing.T) {
		// States 0,1 form a closed ergodic class; state 2 is transient.
		k := mustKernel(t, []float64{
			0.7, 0.3, 0,
			0.4, 0.6, 0,
			0.5, 0.5, 0,
		}, 3)
		d, err := k.Classes(1e-9)
		if err != nil {
			t.Fatal(err)
		}
		if len(d.Recurrent) != 1 {
			t.Fatalf("Recurrent: got %d want 1", len(d.Recurrent))
		}
		if !reflect.DeepEqual(d.Transient, []int{2}) {
			t.Fatalf("Transient: got %v want [2]", d.Transient)
		}
	})
}

// ── passage.go ────────────────────────────────────────────────────────────────

func TestKernel_MeanFirstPassage(t *testing.T) {
	k := ergodic2x2(t) // [[0.7,0.3],[0.4,0.6]]

	cases := []struct {
		name    string
		i, j    int
		wantErr bool
		want    float64
	}{
		// m(0,1): (I-Q)h = 1 with Q=[[0.7]] → 0.3h=1 → h=10/3
		{"expected steps from 0 to 1 is 10/3", 0, 1, false, 10.0 / 3},
		// m(1,0): (I-Q)h = 1 with Q=[[0.6]] → 0.4h=1 → h=2.5
		{"expected steps from 1 to 0 is 2.5", 1, 0, false, 2.5},
		// i==j: 0 by definition
		{"passage time from a state to itself is zero", 0, 0, false, 0},
		{"rejects target index outside state space", 0, 5, true, 0},
		{"rejects negative source index", -1, 0, true, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := k.MeanFirstPassage(tc.i, tc.j)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if err == nil && math.Abs(got-tc.want) > 1e-10 {
				t.Fatalf("got %g want %g", got, tc.want)
			}
		})
	}

	t.Run("nil kernel", func(t *testing.T) {
		var nilK *catrace.Kernel
		_, err := nilK.MeanFirstPassage(0, 1)
		if err == nil {
			t.Fatal("expected error for nil kernel")
		}
	})
}

func TestKernel_CommuteTime(t *testing.T) {
	k := ergodic2x2(t) // m(0,1)=10/3, m(1,0)=2.5; CT=35/6

	t.Run("equals sum of forward and backward mean first-passage times", func(t *testing.T) {
		ct, err := k.CommuteTime(0, 1)
		if err != nil {
			t.Fatal(err)
		}
		want := 10.0/3 + 2.5
		if math.Abs(ct-want) > 1e-10 {
			t.Fatalf("got %g want %g", ct, want)
		}
	})

	t.Run("commute time is symmetric", func(t *testing.T) {
		ct01, err := k.CommuteTime(0, 1)
		if err != nil {
			t.Fatal(err)
		}
		ct10, err := k.CommuteTime(1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(ct01-ct10) > 1e-10 {
			t.Fatalf("CT not symmetric: CT(0,1)=%g CT(1,0)=%g", ct01, ct10)
		}
	})

	t.Run("nil kernel", func(t *testing.T) {
		var nilK *catrace.Kernel
		_, err := nilK.CommuteTime(0, 1)
		if err == nil {
			t.Fatal("expected error for nil kernel")
		}
	})
}

// ── sample.go ─────────────────────────────────────────────────────────────────

func TestKernel_Sample(t *testing.T) {
	t.Run("nil kernel errors", func(t *testing.T) {
		var nilK *catrace.Kernel
		_, err := nilK.Sample(0, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("rejects row index outside state space", func(t *testing.T) {
		k := ergodic2x2(t)
		_, err := k.Sample(5, nil)
		if err == nil {
			t.Fatal("expected error for out-of-range row")
		}
	})

	t.Run("row with all mass on one state always samples that state", func(t *testing.T) {
		// Row 0 has all mass on state 1.
		k, _ := catrace.NewKernel(mat.NewDense(2, 2, []float64{0, 1, 1, 0}), nil)
		rng := rand.New(rand.NewSource(42))
		for i := 0; i < 20; i++ {
			got, err := k.Sample(0, rng)
			if err != nil || got != 1 {
				t.Fatalf("iter %d: got %d err %v", i, got, err)
			}
		}
	})
}

func TestEstimateKernelFromSequence(t *testing.T) {
	cases := []struct {
		name       string
		seq        []int
		nStates    int
		pseudo     float64
		wantErr    bool
		wantKernel bool
	}{
		{
			name:       "all states observed as sources produces a complete kernel",
			seq:        []int{0, 1, 0, 1, 0},
			nStates:    2,
			wantKernel: true,
		},
		{
			name:    "unobserved source state leaves kernel nil",
			seq:     []int{0, 0, 0},
			nStates: 2,
			wantKernel: false,
		},
		{
			// pseudocount applies within observed rows only; unobserved rows still leave Kernel nil
			name:       "pseudocount does not fill unobserved rows",
			seq:        []int{0, 0, 0},
			nStates:    2,
			pseudo:     1,
			wantKernel: false,
		},
		{
			name:    "nStates=0",
			seq:     []int{0, 1},
			nStates: 0,
			wantErr: true,
		},
		{
			name:    "out-of-bounds entry",
			seq:     []int{0, 5},
			nStates: 2,
			wantErr: true,
		},
		{
			name:       "empty sequence produces no transitions",
			seq:        []int{},
			nStates:    2,
			wantKernel: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			est, err := catrace.EstimateKernelFromSequence(tc.seq, tc.nStates, tc.pseudo)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v got err=%v", tc.wantErr, err)
			}
			if err != nil {
				return
			}
			if (est.Kernel != nil) != tc.wantKernel {
				t.Fatalf("wantKernel=%v got Kernel=%v", tc.wantKernel, est.Kernel)
			}
		})
	}

	t.Run("perfectly alternating sequence produces deterministic transition estimates", func(t *testing.T) {
		est, err := catrace.EstimateKernelFromSequence([]int{0, 1, 0, 1, 0}, 2, 0)
		if err != nil || est.Kernel == nil {
			t.Fatalf("err=%v kernel=%v", err, est.Kernel)
		}
		if math.Abs(est.Kernel.P.At(0, 1)-1.0) > 1e-12 {
			t.Fatalf("P[0,1] = %g, want 1", est.Kernel.P.At(0, 1))
		}
		if math.Abs(est.Kernel.P.At(1, 0)-1.0) > 1e-12 {
			t.Fatalf("P[1,0] = %g, want 1", est.Kernel.P.At(1, 0))
		}
	})
}

func TestSampleTraceFromSequence(t *testing.T) {
	cases := []struct {
		name   string
		seq    []int
		subset map[int]bool
		want   []int
	}{
		{
			name:   "retains only states in the observed subset",
			seq:    []int{0, 1, 2, 0, 1, 2},
			subset: map[int]bool{0: true, 2: true},
			want:   []int{0, 2, 0, 2},
		},
		{
			name:   "full subset returns sequence unchanged",
			seq:    []int{0, 1, 0},
			subset: map[int]bool{0: true, 1: true},
			want:   []int{0, 1, 0},
		},
		{
			name:   "disjoint subset produces empty trace",
			seq:    []int{0, 1, 2},
			subset: map[int]bool{3: true},
			want:   []int{},
		},
		{
			name:   "empty sequence yields empty trace",
			seq:    []int{},
			subset: map[int]bool{0: true},
			want:   []int{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := catrace.SampleTraceFromSequence(tc.seq, tc.subset)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestWindowedTraceEstimates(t *testing.T) {
	seq := []int{0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
	subset := map[int]bool{0: true, 1: true}

	t.Run("partitions sequence into overlapping windows and estimates each", func(t *testing.T) {
		windows, err := catrace.WindowedTraceEstimates(seq, subset, 4, 2, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(windows) == 0 {
			t.Fatal("expected at least one window")
		}
	})

	t.Run("rejects window size less than 2 (no transitions possible)", func(t *testing.T) {
		_, err := catrace.WindowedTraceEstimates(seq, subset, 1, 1, 0)
		if err == nil {
			t.Fatal("expected error for windowSize=1")
		}
	})

	t.Run("rejects non-positive step size", func(t *testing.T) {
		_, err := catrace.WindowedTraceEstimates(seq, subset, 4, 0, 0)
		if err == nil {
			t.Fatal("expected error for step=0")
		}
	})

	t.Run("rejects empty subset", func(t *testing.T) {
		_, err := catrace.WindowedTraceEstimates(seq, map[int]bool{}, 4, 1, 0)
		if err == nil {
			t.Fatal("expected error for empty subset")
		}
	})

	t.Run("window longer than seq produces no windows", func(t *testing.T) {
		windows, err := catrace.WindowedTraceEstimates(seq, subset, 100, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(windows) != 0 {
			t.Fatalf("got %d windows, want 0", len(windows))
		}
	})
}

// ── invariants ────────────────────────────────────────────────────────────────

func TestInvariant_StationaryFixedPoint(t *testing.T) {
	cases := []struct {
		name string
		k    *catrace.Kernel
	}{
		{"ergodic 2x2", ergodic2x2(t)},
		{"uniform 2x2", mustKernel(t, []float64{0.5, 0.5, 0.5, 0.5}, 2)},
		{"generic 4x4", mustKernel(t, []float64{
			0.60, 0.20, 0.10, 0.10,
			0.15, 0.55, 0.15, 0.15,
			0.20, 0.20, 0.40, 0.20,
			0.10, 0.20, 0.20, 0.50,
		}, 4)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pi, err := tc.k.Stationary(1e-12, 5000)
			if err != nil {
				t.Fatalf("Stationary failed: %v", err)
			}
			next, err := tc.k.LeftAction(pi)
			if err != nil {
				t.Fatal(err)
			}
			nearlyEqualSlice(t, next, pi, 1e-8)
		})
	}
}

func TestInvariant_CommuteTimeSymmetric(t *testing.T) {
	k := mustKernel(t, []float64{
		0.60, 0.20, 0.10, 0.10,
		0.15, 0.55, 0.15, 0.15,
		0.20, 0.20, 0.40, 0.20,
		0.10, 0.20, 0.20, 0.50,
	}, 4)
	n := k.NumStates()
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			ctij, err := k.CommuteTime(i, j)
			if err != nil {
				t.Fatalf("CT(%d,%d): %v", i, j, err)
			}
			ctji, err := k.CommuteTime(j, i)
			if err != nil {
				t.Fatalf("CT(%d,%d): %v", j, i, err)
			}
			if math.Abs(ctij-ctji) > 1e-10 {
				t.Fatalf("CT(%d,%d)=%g ≠ CT(%d,%d)=%g", i, j, ctij, j, i, ctji)
			}
		}
	}
}

func TestInvariant_AtLeastOneRecurrentClass(t *testing.T) {
	cases := []struct {
		name string
		k    *catrace.Kernel
	}{
		{"ergodic 2x2", ergodic2x2(t)},
		{"period-2 path graph", pathGraph(t)},
		{"two absorbing states", mustKernel(t, []float64{1, 0, 0, 1}, 2)},
		{"one transient state", mustKernel(t, []float64{
			0.7, 0.3, 0,
			0.4, 0.6, 0,
			0.5, 0.5, 0,
		}, 3)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := tc.k.Classes(1e-9)
			if err != nil {
				t.Fatalf("Classes failed: %v", err)
			}
			if len(d.Recurrent) == 0 {
				t.Fatal("every kernel must have at least one recurrent class")
			}
		})
	}
}
