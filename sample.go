package catrace

import (
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

func (k *Kernel) Sample(rowIdx int, rng *rand.Rand) (int, error) {
	if k == nil || k.P == nil {
		return 0, fmt.Errorf("nil kernel")
	}
	n, _ := k.P.Dims()
	if rowIdx < 0 || rowIdx >= n {
		return 0, fmt.Errorf("row index %d out of range [0,%d)", rowIdx, n)
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	u := rng.Float64()
	cum := 0.0
	for j := 0; j < n; j++ {
		cum += k.P.At(rowIdx, j)
		if u <= cum {
			return j, nil
		}
	}
	return n - 1, nil
}

type KernelEstimate struct {
	Counts       *mat.Dense
	Kernel       *Kernel
	RowsObserved []bool
}

func EstimateKernelFromSequence(seq []int, nStates int, pseudocount float64) (*KernelEstimate, error) {
	if nStates <= 0 {
		return nil, fmt.Errorf("nStates must be positive")
	}
	counts := mat.NewDense(nStates, nStates, nil)
	rowsObserved := make([]bool, nStates)
	for i := 0; i+1 < len(seq); i++ {
		a, b := seq[i], seq[i+1]
		if a < 0 || a >= nStates || b < 0 || b >= nStates {
			return nil, fmt.Errorf("sequence entry out of bounds at transition %d: %d -> %d", i, a, b)
		}
		counts.Set(a, b, counts.At(a, b)+1)
		rowsObserved[a] = true
	}
	p := mat.NewDense(nStates, nStates, nil)
	complete := true
	for i := 0; i < nStates; i++ {
		sum := 0.0
		for j := 0; j < nStates; j++ {
			sum += counts.At(i, j)
		}
		if sum == 0 {
			complete = false
			continue
		}
		den := sum + pseudocount*float64(nStates)
		for j := 0; j < nStates; j++ {
			p.Set(i, j, (counts.At(i, j)+pseudocount)/den)
		}
	}
	est := &KernelEstimate{Counts: counts, RowsObserved: rowsObserved}
	if complete {
		k, err := NewKernel(p, nil)
		if err != nil {
			return nil, err
		}
		est.Kernel = k
	}
	return est, nil
}

func SampleTraceFromSequence(seq []int, subset map[int]bool) []int {
	out := make([]int, 0, len(seq))
	for _, x := range seq {
		if subset[x] {
			out = append(out, x)
		}
	}
	return out
}

func WindowedTraceEstimates(seq []int, subset map[int]bool, windowSize, step int, pseudocount float64) ([]*KernelEstimate, error) {
	if windowSize <= 1 {
		return nil, fmt.Errorf("windowSize must be > 1")
	}
	if step <= 0 {
		return nil, fmt.Errorf("step must be positive")
	}
	windows := []*KernelEstimate{}
	for start := 0; start+windowSize <= len(seq); start += step {
		window := seq[start : start+windowSize]
		traced := SampleTraceFromSequence(window, subset)
		mapped := make([]int, 0, len(traced))
		for _, v := range traced {
			if v == 0 {
				mapped = append(mapped, 0)
			} else {
				mapped = append(mapped, 1)
			}
		}
		est, err := EstimateKernelFromSequence(mapped, 2, pseudocount)
		if err != nil {
			return nil, err
		}
		windows = append(windows, est)
	}
	return windows, nil
}
