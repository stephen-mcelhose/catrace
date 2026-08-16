package catrace

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
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
	if maxIter <= 0 {
		maxIter = 1000
	}
	pi, err := cloneProbDist(start, n, tol, "start", true)
	if err != nil {
		return nil, err
	}
	if err := requireErgodicAperiodic(k, tol); err != nil {
		return nil, err
	}
	return powerIterate(k, pi, tol, maxIter)
}

func requireErgodicAperiodic(k *Kernel, tol float64) error {
	cd, err := k.Classes(tol)
	if err != nil {
		return fmt.Errorf("ergodicity check: %w", err)
	}
	if len(cd.Recurrent) != 1 {
		return fmt.Errorf("StationaryFrom requires an ergodic kernel: found %d recurrent classes", len(cd.Recurrent))
	}
	var period int
	for _, p := range cd.Periods {
		period = p
	}
	if period != 1 {
		return fmt.Errorf("StationaryFrom requires an aperiodic kernel: recurrent class has period %d", period)
	}
	return nil
}

func powerIterate(k *Kernel, start []float64, tol float64, maxIter int) ([]float64, error) {
	pi := start
	for iter := 0; iter < maxIter; iter++ {
		next, err := k.LeftAction(pi)
		if err != nil {
			return nil, err
		}
		if err := normalizeVector(next, tol); err != nil {
			return nil, err
		}
		delta := maxAbsDiff(pi, next)
		pi = next
		if delta <= tol {
			return pi, nil
		}
	}
	return pi, fmt.Errorf("stationary iteration did not converge within %d iterations", maxIter)
}

func maxAbsDiff(a, b []float64) float64 {
	delta := 0.0
	for i := range a {
		if d := math.Abs(b[i] - a[i]); d > delta {
			delta = d
		}
	}
	return delta
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
	if alpha < 0 || alpha > 1 {
		return nil, fmt.Errorf("alpha %g is out of range [0, 1]", alpha)
	}
	if maxIter <= 0 {
		maxIter = 1000
	}
	v, err := cloneProbDist(restart, n, tol, "restart", true)
	if err != nil {
		return nil, err
	}
	return personalizedPowerIterate(k, v, alpha, tol, maxIter)
}

func personalizedPowerIterate(k *Kernel, restart []float64, alpha, tol float64, maxIter int) ([]float64, error) {
	n := len(restart)
	x := make([]float64, n)
	copy(x, restart)
	for iter := 0; iter < maxIter; iter++ {
		propagated, err := k.LeftAction(x)
		if err != nil {
			return nil, err
		}
		next := make([]float64, n)
		for i := range next {
			next[i] = alpha*restart[i] + (1-alpha)*propagated[i]
		}
		delta := maxAbsDiff(x, next)
		x = next
		if delta <= tol {
			return x, nil
		}
	}
	return x, fmt.Errorf("PersonalizedPageRank did not converge within %d iterations", maxIter)
}

// TeleportingKernel constructs the teleporting Markov chain
//
//	T = α·restart·𝟙ᵀ + (1−α)·P
//
// whose stationary distribution equals the PersonalizedPageRank vector for the
// same restart and alpha. The teleporting chain is strongly connected for any
// alpha ∈ (0, 1] and any stochastic restart, so its stationary distribution
// exists and is unique.
//
// For alpha = 1 every row of T equals restart (the chain ignores P entirely).
// For alpha = 0 T equals P unchanged.
//
// TeleportingKernel is useful for visualisation: call ToHTML on the returned
// kernel and nodes will be sized by PPR score. Use MinEdge in VisualiseOptions
// to suppress low-weight teleportation arcs.
func (k *Kernel) TeleportingKernel(restart []float64, alpha float64) (*Kernel, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	if n == 0 {
		return nil, fmt.Errorf("empty kernel")
	}
	if alpha < 0 || alpha > 1 {
		return nil, fmt.Errorf("alpha %g is out of range [0, 1]", alpha)
	}
	const tol = 1e-9
	v, err := cloneProbDist(restart, n, tol, "restart", false)
	if err != nil {
		return nil, err
	}
	data := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			data[i*n+j] = alpha*v[j] + (1-alpha)*k.P.At(i, j)
		}
	}
	return NewKernel(mat.NewDense(n, n, data), k.StateNames)
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
