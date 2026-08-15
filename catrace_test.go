package catrace_test

import (
	"fmt"
	"testing"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func TestTraceMatchesRestrictionTheoremPattern(t *testing.T) {
	parent, err := catrace.NewKernel(mat.NewDense(4, 4, []float64{
		0.60, 0.20, 0.10, 0.10,
		0.15, 0.55, 0.15, 0.15,
		0.20, 0.20, 0.40, 0.20,
		0.10, 0.20, 0.20, 0.50,
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	trace, err := parent.Trace([]int{0, 1}, 1e-12)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := trace.IsTraceOf(parent, []int{0, 1}, 1e-12); err != nil || !ok {
		t.Fatalf("trace verification failed: ok=%v err=%v", ok, err)
	}
	piParent, err := parent.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatal(err)
	}
	piTrace, err := trace.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatal(err)
	}
	restricted, err := catrace.RestrictDistribution(piParent, []int{0, 1}, 1e-12)
	if err != nil {
		t.Fatal(err)
	}
	for i := range piTrace {
		if diff := piTrace[i] - restricted[i]; diff < -1e-8 || diff > 1e-8 {
			t.Fatalf("stationary mismatch at %d: got %g want %g", i, piTrace[i], restricted[i])
		}
	}
}

// ExampleNewRandomWalkKernel demonstrates how a graph becomes a Markov kernel.
//
// Consider this small network:
//
//	D – A – B
//	    |  /
//	    C
//
// A is the hub — it connects to B, C, and D.
// D is a dead end — it connects only to A.
//
// The random walk assigns each node a transition probability proportional to
// its edge weights. From A (degree 3), you go to each neighbor with prob 1/3.
// From D (degree 1), you always go to A.
//
// The stationary distribution equals the normalized degree of each node:
//
//	degree(A)=3, degree(B)=2, degree(C)=2, degree(D)=1  →  total=8
//	π = [3/8, 2/8, 2/8, 1/8] = [0.375, 0.250, 0.250, 0.125]
//
// A is visited most often. D is visited least — it is a structural dead end.
// ExampleKernel_StationaryFrom shows that a skewed starting distribution
// converges to the same stationary distribution as the uniform start.
//
// The same 4-node graph as ExampleNewRandomWalkKernel (A–B–C–D) is used.
// Starting from [0.7, 0.1, 0.1, 0.1] (mass concentrated on A) converges
// to the degree-proportional stationary distribution π = [3/8, 2/8, 2/8, 1/8].
func ExampleKernel_StationaryFrom() {
	adj := mat.NewDense(4, 4, []float64{
		0, 1, 1, 1, // A connects to B, C, D
		1, 0, 1, 0, // B connects to A, C
		1, 1, 0, 0, // C connects to A, B
		1, 0, 0, 0, // D connects to A only
	})
	k, _ := catrace.NewRandomWalkKernel(adj, []string{"A", "B", "C", "D"})

	// Skewed start: most mass on A.
	pi, _ := k.StationaryFrom([]float64{0.7, 0.1, 0.1, 0.1}, 1e-12, 5000)
	fmt.Printf("π(A) = %.3f\n", pi[0])
	fmt.Printf("π(B) = %.3f\n", pi[1])
	fmt.Printf("π(C) = %.3f\n", pi[2])
	fmt.Printf("π(D) = %.3f\n", pi[3])
	// Output:
	// π(A) = 0.375
	// π(B) = 0.250
	// π(C) = 0.250
	// π(D) = 0.125
}

// ExampleKernel_PersonalizedPageRank shows how seeding on a leaf node
// concentrates PPR mass near that node compared to the global stationary.
//
// The same 4-node graph (A–B–C–D) is used. D is a structural dead-end
// (degree 1). Seeding entirely on D with alpha=0.15 biases the PPR vector
// toward D and its immediate neighbour A, compared with the global π which
// assigns D the least mass.
func ExampleKernel_PersonalizedPageRank() {
	adj := mat.NewDense(4, 4, []float64{
		0, 1, 1, 1, // A connects to B, C, D
		1, 0, 1, 0, // B connects to A, C
		1, 1, 0, 0, // C connects to A, B
		1, 0, 0, 0, // D connects to A only
	})
	k, _ := catrace.NewRandomWalkKernel(adj, []string{"A", "B", "C", "D"})

	// Seed entirely on D (the leaf). alpha=0.15 is the standard PageRank damping.
	ppr, _ := k.PersonalizedPageRank([]float64{0, 0, 0, 1}, 0.15, 1e-12, 5000)
	fmt.Printf("ppr(A) = %.3f\n", ppr[0])
	fmt.Printf("ppr(B) = %.3f\n", ppr[1])
	fmt.Printf("ppr(C) = %.3f\n", ppr[2])
	fmt.Printf("ppr(D) = %.3f\n", ppr[3])
}

func ExampleNewRandomWalkKernel() {
	// Symmetric adjacency matrix (1 = edge present, 0 = no edge).
	adj := mat.NewDense(4, 4, []float64{
		//  A  B  C  D
		0, 1, 1, 1, // A connects to B, C, D
		1, 0, 1, 0, // B connects to A, C
		1, 1, 0, 0, // C connects to A, B
		1, 0, 0, 0, // D connects to A only
	})

	k, err := catrace.NewRandomWalkKernel(adj, []string{"A", "B", "C", "D"})
	if err != nil {
		panic(err)
	}

	// Each row of the kernel is the adjacency row divided by the node's degree.
	fmt.Printf("P(A→B) = %.3f\n", k.P.At(0, 1)) // degree(A)=3, so 1/3
	fmt.Printf("P(D→A) = %.3f\n", k.P.At(3, 0)) // degree(D)=1, so 1/1

	// Stationary distribution: π(i) = degree(i) / total degree.
	pi, _ := k.Stationary(1e-12, 5000)
	fmt.Printf("π(A)   = %.3f\n", pi[0])
	fmt.Printf("π(B)   = %.3f\n", pi[1])
	fmt.Printf("π(C)   = %.3f\n", pi[2])
	fmt.Printf("π(D)   = %.3f\n", pi[3])
	// Output:
	// P(A→B) = 0.333
	// P(D→A) = 1.000
	// π(A)   = 0.375
	// π(B)   = 0.250
	// π(C)   = 0.250
	// π(D)   = 0.125
}

func ExampleAgent_QualiaKernel() {
	D := mat.NewDense(2, 2, []float64{
		0.9, 0.1,
		0.2, 0.8,
	})
	A := mat.NewDense(2, 2, []float64{
		0.7, 0.3,
		0.4, 0.6,
	})
	P := mat.NewDense(2, 2, []float64{
		0.8, 0.2,
		0.1, 0.9,
	})
	agent := &catrace.Agent{D: D, A: A, P: P}
	Q, _ := agent.QualiaKernel()
	fmt.Printf("%.3f %.3f\n", Q.P.At(0, 0), Q.P.At(0, 1))
	fmt.Printf("%.3f %.3f\n", Q.P.At(1, 0), Q.P.At(1, 1))
	// Output:
	// 0.569 0.431
	// 0.422 0.578
}

func ExampleKernel_Trace() {
	k, _ := catrace.NewKernel(mat.NewDense(3, 3, []float64{
		0.6, 0.3, 0.1,
		0.2, 0.6, 0.2,
		0.4, 0.1, 0.5,
	}), nil)
	tr, _ := k.Trace([]int{0, 1}, 1e-12)
	fmt.Printf("%.3f %.3f\n", tr.P.At(0, 0), tr.P.At(0, 1))
	fmt.Printf("%.3f %.3f\n", tr.P.At(1, 0), tr.P.At(1, 1))
	// Output:
	// 0.680 0.320
	// 0.360 0.640
}
