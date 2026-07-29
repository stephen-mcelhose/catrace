package catrace

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// RectKernel represents a row-stochastic rectangular kernel.
// RowNames and ColNames label the row and column spaces respectively.
type RectKernel struct {
	P        *mat.Dense // transition matrix
	RowNames []string
	ColNames []string
}

// Kernel represents a row-stochastic square Markov kernel.
// StateNames labels both the row and column indices.
type Kernel struct {
	P          *mat.Dense // transition matrix
	StateNames []string
}

// NewRectKernel constructs a validated rectangular row-stochastic kernel.
// rowNames and colNames may be nil, in which case default labels are assigned.
func NewRectKernel(p *mat.Dense, rowNames, colNames []string) (*RectKernel, error) {
	if p == nil {
		return nil, fmt.Errorf("nil kernel matrix")
	}
	r, c := p.Dims()
	if len(rowNames) > 0 && len(rowNames) != r {
		return nil, fmt.Errorf("row name count %d does not match row count %d", len(rowNames), r)
	}
	if len(colNames) > 0 && len(colNames) != c {
		return nil, fmt.Errorf("column name count %d does not match column count %d", len(colNames), c)
	}
	rk := &RectKernel{
		P:        cloneDense(p),
		RowNames: cloneStrings(rowNames),
		ColNames: cloneStrings(colNames),
	}
	if len(rk.RowNames) == 0 {
		rk.RowNames = defaultNames("r", r)
	}
	if len(rk.ColNames) == 0 {
		rk.ColNames = defaultNames("c", c)
	}
	if err := rk.Validate(1e-9); err != nil {
		return nil, err
	}
	return rk, nil
}

// NewKernel constructs a validated square row-stochastic kernel.
// names may be nil, in which case default state labels are assigned.
func NewKernel(p *mat.Dense, names []string) (*Kernel, error) {
	if p == nil {
		return nil, fmt.Errorf("nil kernel matrix")
	}
	r, c := p.Dims()
	if r != c {
		return nil, fmt.Errorf("kernel must be square, got %dx%d", r, c)
	}
	if len(names) > 0 && len(names) != r {
		return nil, fmt.Errorf("state name count %d does not match state count %d", len(names), r)
	}
	k := &Kernel{P: cloneDense(p), StateNames: cloneStrings(names)}
	if len(k.StateNames) == 0 {
		k.StateNames = defaultNames("s", r)
	}
	if err := k.Validate(1e-9); err != nil {
		return nil, err
	}
	return k, nil
}

// Validate checks that every row is non-negative and sums to 1 within tol.
func (k *RectKernel) Validate(tol float64) error {
	if k == nil || k.P == nil {
		return fmt.Errorf("nil kernel")
	}
	r, c := k.P.Dims()
	for i := 0; i < r; i++ {
		sum := 0.0
		for j := 0; j < c; j++ {
			v := k.P.At(i, j)
			if v < -tol {
				return fmt.Errorf("negative entry at (%d,%d): %g", i, j, v)
			}
			sum += v
		}
		if math.Abs(sum-1.0) > tol {
			return fmt.Errorf("row %d sums to %g, want 1 within tol %g", i, sum, tol)
		}
	}
	return nil
}

// Validate checks that k is square and every row is non-negative and sums to 1 within tol.
func (k *Kernel) Validate(tol float64) error {
	if k == nil || k.P == nil {
		return fmt.Errorf("nil kernel")
	}
	r, c := k.P.Dims()
	if r != c {
		return fmt.Errorf("kernel must be square, got %dx%d", r, c)
	}
	return (&RectKernel{P: k.P}).Validate(tol)
}

// Clone returns a deep copy of k.
func (k *Kernel) Clone() *Kernel {
	if k == nil {
		return nil
	}
	return &Kernel{P: cloneDense(k.P), StateNames: cloneStrings(k.StateNames)}
}

// NumStates returns the number of states in the kernel.
func (k *Kernel) NumStates() int {
	if k == nil || k.P == nil {
		return 0
	}
	r, _ := k.P.Dims()
	return r
}

// NormalizeRows rescales each row to sum to 1, clamping near-zero negatives to 0.
// Returns an error if any row contains a significantly negative entry or a near-zero sum.
func (k *Kernel) NormalizeRows(tol float64) error {
	if k == nil || k.P == nil {
		return fmt.Errorf("nil kernel")
	}
	r, c := k.P.Dims()
	for i := 0; i < r; i++ {
		sum := 0.0
		for j := 0; j < c; j++ {
			v := k.P.At(i, j)
			if v < 0 && math.Abs(v) <= tol {
				v = 0
				k.P.Set(i, j, 0)
			}
			if v < 0 {
				return fmt.Errorf("negative entry at (%d,%d): %g", i, j, v)
			}
			sum += v
		}
		if sum <= tol {
			return fmt.Errorf("row %d has near-zero sum %g", i, sum)
		}
		for j := 0; j < c; j++ {
			k.P.Set(i, j, k.P.At(i, j)/sum)
		}
	}
	return k.Validate(math.Max(tol, 1e-9))
}

// Multiply returns the matrix product k·other as a new Kernel.
// Both operands must be square and of the same dimension.
func (k *Kernel) Multiply(other *Kernel) (*Kernel, error) {
	if k == nil || other == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	r1, c1 := k.P.Dims()
	r2, c2 := other.P.Dims()
	if c1 != r2 {
		return nil, fmt.Errorf("dimension mismatch: %dx%d times %dx%d", r1, c1, r2, c2)
	}
	if r1 != c2 {
		return nil, fmt.Errorf("result is not square: %dx%d", r1, c2)
	}
	var out mat.Dense
	out.Mul(k.P, other.P)
	res, err := NewKernel(&out, k.StateNames)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// LeftAction evolves dist one step by computing π·P, returning the next distribution.
// dist must have length equal to NumStates.
func (k *Kernel) LeftAction(dist []float64) ([]float64, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n, _ := k.P.Dims()
	if len(dist) != n {
		return nil, fmt.Errorf("distribution length %d does not match kernel size %d", len(dist), n)
	}
	out := make([]float64, n)
	for j := 0; j < n; j++ {
		s := 0.0
		for i := 0; i < n; i++ {
			s += dist[i] * k.P.At(i, j)
		}
		out[j] = s
	}
	return out, nil
}
