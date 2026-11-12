package stream

import (
	"fmt"
	"math"
	"sync"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

// Estimator maintains exponentially decayed transition counts over an open vocabulary of states
// and constructs un-damped *catrace.Kernel instances backed by *mat.Dense.
type Estimator struct {
	mu       sync.RWMutex
	decay    float64
	nameToID map[string]int
	idToName []string
	counts   [][]float64
	lastNode string
}

// NewEstimator creates a streaming EMA kernel estimator.
// halfLife is measured in observed transitions (<= 0 disables decay).
func NewEstimator(halfLife float64) *Estimator {
	gamma := 1.0
	if halfLife > 0 {
		gamma = math.Exp(-math.Ln2 / halfLife)
	}
	return &Estimator{
		decay:    gamma,
		nameToID: make(map[string]int),
	}
}

func (e *Estimator) intern(name string) int {
	if id, ok := e.nameToID[name]; ok {
		return id
	}
	id := len(e.idToName)
	e.nameToID[name] = id
	e.idToName = append(e.idToName, name)

	n := len(e.idToName)
	for i := range e.counts {
		e.counts[i] = append(e.counts[i], 0)
	}
	e.counts = append(e.counts, make([]float64, n))
	return id
}

// ObserveNext ingests a sequential event stream, recording transition prev -> stateName.
func (e *Estimator) ObserveNext(stateName string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.lastNode != "" {
		e.observeLocked(e.lastNode, stateName, 1.0)
	}
	e.lastNode = stateName
}

// ObserveTransition explicitly records u -> v with weight w.
func (e *Estimator) ObserveTransition(from, to string, weight float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.observeLocked(from, to, weight)
}

func (e *Estimator) observeLocked(from, to string, weight float64) {
	u := e.intern(from)
	v := e.intern(to)

	if e.decay < 1.0 {
		for i := range e.counts {
			for j := range e.counts[i] {
				e.counts[i][j] *= e.decay
			}
		}
	}
	e.counts[u][v] += weight
}

// ResetSequence clears lastNode between independent episodes/sessions without forgetting learned counts.
func (e *Estimator) ResetSequence() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastNode = ""
}

// rowMassSnapshot returns the total decayed outgoing transition mass per state name.
func (e *Estimator) rowMassSnapshot() map[string]float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	out := make(map[string]float64, len(e.idToName))
	for i, name := range e.idToName {
		var sum float64
		for _, v := range e.counts[i] {
			sum += v
		}
		out[name] = sum
	}
	return out
}

// Kernel builds an un-damped row-stochastic *catrace.Kernel with state names attached.
// Unvisited rows become absorbing self-loops (P[i][i] = 1.0) so traps remain visible.
func (e *Estimator) Kernel() (*catrace.Kernel, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	n := len(e.idToName)
	if n == 0 {
		return nil, fmt.Errorf("stream: no states observed")
	}

	flat := make([]float64, n*n)
	for i := 0; i < n; i++ {
		var rowSum float64
		for j := 0; j < n; j++ {
			rowSum += e.counts[i][j]
		}
		if rowSum == 0 {
			flat[i*n+i] = 1.0 // Un-damped absorbing sink
		} else {
			for j := 0; j < n; j++ {
				flat[i*n+j] = e.counts[i][j] / rowSum
			}
		}
	}

	p := mat.NewDense(n, n, flat)
	names := append([]string(nil), e.idToName...)
	return catrace.NewKernel(p, names)
}
