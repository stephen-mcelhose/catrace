package catrace

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// MeanFirstPassage returns the expected number of steps to hit target j starting from i.
func (k *Kernel) MeanFirstPassage(i, j int) (float64, error) {
	if k == nil || k.P == nil {
		return 0, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if i < 0 || i >= n || j < 0 || j >= n {
		return 0, fmt.Errorf("indices out of range")
	}
	if i == j {
		return 0, nil
	}
	idx := make([]int, 0, n-1)
	pos := map[int]int{}
	for s := 0; s < n; s++ {
		if s == j {
			continue
		}
		pos[s] = len(idx)
		idx = append(idx, s)
	}
	Q := submatrix(k.P, idx, idx)
	m := len(idx)
	A := mat.NewDense(m, m, nil)
	b := mat.NewDense(m, 1, nil)
	for r := 0; r < m; r++ {
		for c := 0; c < m; c++ {
			v := -Q.At(r, c)
			if r == c {
				v += 1
			}
			A.Set(r, c, v)
		}
		b.Set(r, 0, 1)
	}
	var h mat.Dense
	if err := h.Solve(A, b); err != nil {
		return 0, fmt.Errorf("failed to solve mean first-passage system: %w", err)
	}
	return h.At(pos[i], 0), nil
}

// CommuteTime returns m(i,j)+m(j,i).
func (k *Kernel) CommuteTime(i, j int) (float64, error) {
	mij, err := k.MeanFirstPassage(i, j)
	if err != nil {
		return 0, err
	}
	mji, err := k.MeanFirstPassage(j, i)
	if err != nil {
		return 0, err
	}
	return mij + mji, nil
}
