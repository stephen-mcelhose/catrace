package catrace

import (
	"fmt"
	"math"
)

// Stationary computes a stationary distribution pi such that pi P = pi.
// It uses power iteration on an initial uniform distribution.
func (k *Kernel) Stationary(tol float64, maxIter int) ([]float64, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if n == 0 {
		return nil, fmt.Errorf("empty kernel")
	}
	if maxIter <= 0 {
		maxIter = 1000
	}
	pi := make([]float64, n)
	for i := range pi {
		pi[i] = 1.0 / float64(n)
	}
	for iter := 0; iter < maxIter; iter++ {
		next, err := k.LeftAction(pi)
		if err != nil {
			return nil, err
		}
		if err := normalizeVector(next, tol); err != nil {
			return nil, err
		}
		delta := 0.0
		for i := range pi {
			d := math.Abs(next[i] - pi[i])
			if d > delta {
				delta = d
			}
		}
		pi = next
		if delta <= tol {
			return pi, nil
		}
	}
	return pi, fmt.Errorf("stationary iteration did not converge within %d iterations", maxIter)
}

// EntropyRate returns the entropy rate of the chain in the specified log base.
// For base 2 the unit is bits per step.
func (k *Kernel) EntropyRate(base float64) (float64, error) {
	if base <= 0 || base == 1 {
		return 0, fmt.Errorf("invalid logarithm base %g", base)
	}
	pi, err := k.Stationary(1e-12, 5000)
	if err != nil {
		return 0, err
	}
	logBase := math.Log(base)
	h := 0.0
	n := k.NumStates()
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			p := k.P.At(i, j)
			if p > 0 {
				h -= pi[i] * p * (math.Log(p) / logBase)
			}
		}
	}
	return h, nil
}
