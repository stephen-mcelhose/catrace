package catrace

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// Trace computes the induced trace kernel on a subset A of states.
// If P is partitioned as
//
// 	P = [ a  b ]
// 	    [ d  c ]
//
// with rows/columns of a indexed by A, then the trace is
//
// 	P_A = a + b (I - c)^(-1) d,
//
// provided the excursion operator is well-defined.
func (k *Kernel) Trace(subset []int, tol float64) (*Kernel, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if len(subset) == 0 {
		return nil, fmt.Errorf("empty subset")
	}
	for _, idx := range subset {
		if idx < 0 || idx >= n {
			return nil, fmt.Errorf("subset index %d out of range [0,%d)", idx, n)
		}
	}
	A := sortedCopy(subset)
	Ac := complementIndices(n, A)
	a := submatrix(k.P, A, A)
	if len(Ac) == 0 {
		names := make([]string, len(A))
		for i, idx := range A {
			names[i] = k.StateNames[idx]
		}
		return NewKernel(a, names)
	}
	b := submatrix(k.P, A, Ac)
	c := submatrix(k.P, Ac, Ac)
	d := submatrix(k.P, Ac, A)

	m, _ := c.Dims()
	IminusC := mat.NewDense(m, m, nil)
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			v := -c.At(i, j)
			if i == j {
				v += 1
			}
			IminusC.Set(i, j, v)
		}
	}

	var X mat.Dense
	if err := X.Solve(IminusC, d); err != nil {
		return nil, fmt.Errorf("trace solve failed: %w", err)
	}
	var bX mat.Dense
	bX.Mul(b, &X)
	var t mat.Dense
	t.Add(a, &bX)

	names := make([]string, len(A))
	for i, idx := range A {
		names[i] = k.StateNames[idx]
	}
	tr, err := NewKernel(&t, names)
	if err != nil {
		return nil, err
	}
	if err := tr.NormalizeRows(tol); err != nil {
		return nil, err
	}
	return tr, nil
}

// IsTraceOf checks whether k matches the trace of parent on subset within tol.
func (k *Kernel) IsTraceOf(parent *Kernel, subset []int, tol float64) (bool, error) {
	if k == nil || parent == nil {
		return false, fmt.Errorf("nil kernel")
	}
	tr, err := parent.Trace(subset, tol)
	if err != nil {
		return false, err
	}
	r1, c1 := k.P.Dims()
	r2, c2 := tr.P.Dims()
	if r1 != r2 || c1 != c2 {
		return false, nil
	}
	for i := 0; i < r1; i++ {
		for j := 0; j < c1; j++ {
			if !nearlyEqual(k.P.At(i, j), tr.P.At(i, j), tol) {
				return false, nil
			}
		}
	}
	return true, nil
}

// RestrictDistribution normalizes the restriction of a distribution to subset.
func RestrictDistribution(pi []float64, subset []int, tol float64) ([]float64, error) {
	if len(subset) == 0 {
		return nil, fmt.Errorf("empty subset")
	}
	out := make([]float64, len(subset))
	for i, idx := range subset {
		if idx < 0 || idx >= len(pi) {
			return nil, fmt.Errorf("subset index %d out of range", idx)
		}
		out[i] = pi[idx]
	}
	if err := normalizeVector(out, tol); err != nil {
		return nil, err
	}
	return out, nil
}

// CanonicalSubset returns sorted unique subset indices.
func CanonicalSubset(subset []int) []int {
	seen := map[int]bool{}
	out := make([]int, 0, len(subset))
	for _, idx := range subset {
		if !seen[idx] {
			seen[idx] = true
			out = append(out, idx)
		}
	}
	sort.Ints(out)
	return out
}
