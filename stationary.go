package catrace

import (
	"fmt"
	"math"
)

// Stationary computes a stationary distribution pi such that pi P = pi.
// It uses power iteration on an initial uniform distribution.
//
// This method requires the chain to be ergodic (irreducible and aperiodic).
// If the chain has multiple recurrent classes, power iteration from a uniform
// start will not converge to a unique stationary distribution and the method
// will return an error after maxIter iterations. Use Classes to inspect the
// chain structure before calling Stationary on reducible chains.
func (k *Kernel) Stationary(tol float64, maxIter int) ([]float64, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if n == 0 {
		return nil, fmt.Errorf("empty kernel")
	}
	uniform := make([]float64, n)
	for i := range uniform {
		uniform[i] = 1.0 / float64(n)
	}
	return k.StationaryFrom(uniform, tol, maxIter)
}

// StationaryFrom runs power iteration initialised from start rather than the
// uniform distribution. start must be a valid probability distribution over
// k.NumStates() states: all entries ≥ 0 (entries in (−tol, 0) are clamped to
// zero) and the sum within tol of 1.
//
// Use Stationary for the standard uniform-start case. Use StationaryFrom when
// the initial distribution carries semantic meaning, or to verify that
// convergence is distribution-independent (ergodicity check).
//
// Note: for non-ergodic chains convergence is not guaranteed. Use Classes to
// inspect chain structure before calling StationaryFrom on reducible chains.
func (k *Kernel) StationaryFrom(start []float64, tol float64, maxIter int) ([]float64, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if n == 0 {
		return nil, fmt.Errorf("empty kernel")
	}
	if len(start) != n {
		return nil, fmt.Errorf("start length %d does not match kernel size %d", len(start), n)
	}
	if maxIter <= 0 {
		maxIter = 1000
	}
	pi := make([]float64, n)
	sum := 0.0
	for i, v := range start {
		if v < 0 && math.Abs(v) <= tol {
			v = 0
		}
		if v < 0 {
			return nil, fmt.Errorf("start[%d] is negative: %g", i, v)
		}
		pi[i] = v
		sum += v
	}
	if math.Abs(sum-1.0) > tol {
		return nil, fmt.Errorf("start does not sum to 1 (got %g, tol %g)", sum, tol)
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
			if d := math.Abs(next[i] - pi[i]); d > delta {
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

// PersonalizedPageRank computes the Personalized PageRank vector for the given
// restart distribution and teleportation weight alpha ∈ [0, 1].
//
// Each iteration step blends the propagated distribution back toward restart:
//
//	x ← α·restart + (1−α)·(x·P)
//
// The fixed point satisfies x* = α·restart + (1−α)·x*·P, which is the
// stationary distribution of the teleporting chain α·restart·𝟙ᵀ + (1−α)·P.
//
// For alpha ∈ (0, 1] convergence is guaranteed regardless of chain structure,
// because the teleporting chain is strongly connected. For alpha = 0 this
// reduces to plain power iteration from restart (equivalent to
// StationaryFrom(restart, tol, maxIter)).
func (k *Kernel) PersonalizedPageRank(restart []float64, alpha, tol float64, maxIter int) ([]float64, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if n == 0 {
		return nil, fmt.Errorf("empty kernel")
	}
	if len(restart) != n {
		return nil, fmt.Errorf("restart length %d does not match kernel size %d", len(restart), n)
	}
	if alpha < 0 || alpha > 1 {
		return nil, fmt.Errorf("alpha %g is out of range [0, 1]", alpha)
	}
	if maxIter <= 0 {
		maxIter = 1000
	}
	// Validate and copy restart — do not mutate the caller's slice.
	v := make([]float64, n)
	sum := 0.0
	for i, val := range restart {
		if val < 0 && math.Abs(val) <= tol {
			val = 0
		}
		if val < 0 {
			return nil, fmt.Errorf("restart[%d] is negative: %g", i, val)
		}
		v[i] = val
		sum += val
	}
	if math.Abs(sum-1.0) > tol {
		return nil, fmt.Errorf("restart does not sum to 1 (got %g, tol %g)", sum, tol)
	}
	x := make([]float64, n)
	copy(x, v)
	for iter := 0; iter < maxIter; iter++ {
		propagated, err := k.LeftAction(x)
		if err != nil {
			return nil, err
		}
		next := make([]float64, n)
		for i := range next {
			next[i] = alpha*v[i] + (1-alpha)*propagated[i]
		}
		delta := 0.0
		for i := range x {
			if d := math.Abs(next[i] - x[i]); d > delta {
				delta = d
			}
		}
		x = next
		if delta <= tol {
			return x, nil
		}
	}
	return x, fmt.Errorf("PersonalizedPageRank did not converge within %d iterations", maxIter)
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
