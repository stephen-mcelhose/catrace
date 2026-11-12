package stream

import (
	"fmt"
	"math"
	"sync"

	"github.com/stephen-mcelhose/catrace"
)

// Monitor tracks behavioral drift between a fast-adapting kernel and a slow baseline kernel.
type Monitor struct {
	mu   sync.RWMutex
	Fast *Estimator
	Slow *Estimator
}

func NewMonitor(fastHalfLife, slowHalfLife float64) *Monitor {
	return &Monitor{
		Fast: NewEstimator(fastHalfLife),
		Slow: NewEstimator(slowHalfLife),
	}
}

// ObserveNext atomically advances both estimators on a sequential event stream.
func (m *Monitor) ObserveNext(stateName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Fast.ObserveNext(stateName)
	m.Slow.ObserveNext(stateName)
}

// ResetSequence atomically clears sequence tracking across session boundaries.
func (m *Monitor) ResetSequence() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Fast.ResetSequence()
	m.Slow.ResetSequence()
}

// DriftScore computes the activity-weighted expected Jensen-Shannon divergence
// between the transition rows of the fast and slow kernels.
// If keepStates is non-empty, both kernels are folded via (*Kernel).Trace onto keepStates first.
func (m *Monitor) DriftScore(keepStates []string, tol float64) (float64, error) {
	m.mu.RLock()
	kFast, err := m.Fast.Kernel()
	if err != nil {
		m.mu.RUnlock()
		return 0, err
	}
	kSlow, err := m.Slow.Kernel()
	if err != nil {
		m.mu.RUnlock()
		return 0, err
	}
	fastWeights := m.Fast.rowMassSnapshot()
	slowWeights := m.Slow.rowMassSnapshot()
	m.mu.RUnlock()

	if len(keepStates) > 0 {
		var errF, errS error
		kFast, errF = TraceByName(kFast, keepStates, tol)
		kSlow, errS = TraceByName(kSlow, keepStates, tol)
		if errF != nil {
			return 0, fmt.Errorf("fast trace: %w", errF)
		}
		if errS != nil {
			return 0, fmt.Errorf("slow trace: %w", errS)
		}
	}

	union := make(map[string]int)
	for _, name := range kFast.StateNames {
		if _, ok := union[name]; !ok {
			union[name] = len(union)
		}
	}
	for _, name := range kSlow.StateNames {
		if _, ok := union[name]; !ok {
			union[name] = len(union)
		}
	}

	fastIdx := indexMap(kFast.StateNames)
	slowIdx := indexMap(kSlow.StateNames)
	uDim := len(union)

	var weightedJSD, totalWeight float64

	for state := range union {
		fRow, inFast := fastIdx[state]
		sRow, inSlow := slowIdx[state]
		if !inFast || !inSlow {
			continue
		}

		// Ignore unvisited/terminal rows (such as the current lastNode)
		if fastWeights[state] <= 1e-9 || slowWeights[state] <= 1e-9 {
			continue
		}

		p := make([]float64, uDim)
		q := make([]float64, uDim)
		for j, targetName := range kFast.StateNames {
			p[union[targetName]] = kFast.P.At(fRow, j)
		}
		for j, targetName := range kSlow.StateNames {
			q[union[targetName]] = kSlow.P.At(sRow, j)
		}

		w := fastWeights[state]
		weightedJSD += w * jsd(p, q)
		totalWeight += w
	}

	if totalWeight == 0 {
		return 0, nil
	}
	return weightedJSD / totalWeight, nil
}

// TraceByName folds k onto the named boundary states via (*Kernel).Trace.
// States in keepStates that are absent from k are ignored; at least one must match.
func TraceByName(k *catrace.Kernel, keepStates []string, tol float64) (*catrace.Kernel, error) {
	keepSet := make(map[string]bool, len(keepStates))
	for _, s := range keepStates {
		keepSet[s] = true
	}
	var keepIdx []int
	var keepNames []string
	for i, name := range k.StateNames {
		if keepSet[name] {
			keepIdx = append(keepIdx, i)
			keepNames = append(keepNames, name)
		}
	}
	if len(keepIdx) == 0 {
		return nil, fmt.Errorf("none of keepStates %v found in kernel", keepStates)
	}
	traced, err := k.Trace(keepIdx, tol)
	if err != nil {
		return nil, err
	}
	if len(traced.StateNames) == 0 {
		traced.StateNames = keepNames
	}
	return traced, nil
}

func indexMap(names []string) map[string]int {
	m := make(map[string]int, len(names))
	for i, name := range names {
		m[name] = i
	}
	return m
}

func jsd(p, q []float64) float64 {
	m := make([]float64, len(p))
	for i := range p {
		m[i] = 0.5 * (p[i] + q[i])
	}
	return 0.5*klDiv(p, m) + 0.5*klDiv(q, m)
}

func klDiv(a, b []float64) float64 {
	var sum float64
	for i := range a {
		if a[i] > 1e-12 && b[i] > 1e-12 {
			sum += a[i] * math.Log(a[i]/b[i])
		}
	}
	return sum
}
