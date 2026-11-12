package stream

import (
	"fmt"
	"math"
	"sort"

	"github.com/stephen-mcelhose/catrace"
)

type KnobSuggestion struct {
	FromState     string  `json:"from"`
	ToState       string  `json:"to"`
	Epsilon       float64 `json:"epsilon"`
	BaselineMFPT  float64 `json:"baseline_mfpt"`
	ProjectedMFPT float64 `json:"projected_mfpt"`
	MFPTDelta     float64 `json:"mfpt_delta"`
}

type SuggestConfig struct {
	TargetState  string
	Controllable []string
	Epsilon      float64
	TopK         int
}

type DiagnosticReport struct {
	BoundaryStates []string         `json:"boundary_states"`
	EntropyRate    float64          `json:"entropy_rate_bits"`
	Classes        any              `json:"classes"`
	BaselineMFPT   float64          `json:"baseline_mfpt"`
	Suggestions    []KnobSuggestion `json:"suggestions"`
}

// AnalyzeAndSuggest inspects a (typically traced) *catrace.Kernel, reports its entropy rate,
// communicating classes, and ranks single-row perturbations minimizing MeanFirstPassage to TargetState.
func AnalyzeAndSuggest(k *catrace.Kernel, cfg SuggestConfig) (*DiagnosticReport, error) {
	names := k.StateNames
	n := k.NumStates()
	if n == 0 || len(names) != n {
		return nil, fmt.Errorf("stream: invalid kernel dimensions (%d) or state names (%d)", n, len(names))
	}

	idxMap := indexMap(names)
	targetIdx, ok := idxMap[cfg.TargetState]
	if !ok {
		return nil, fmt.Errorf("target state %q not found in kernel StateNames %v", cfg.TargetState, names)
	}

	entropy, _ := k.EntropyRate(2)
	classes, err := k.Classes(1e-12)
	if err != nil {
		return nil, fmt.Errorf("classes analysis failed: %w", err)
	}

	baseMFPT, err := avgMFPTToTarget(k, targetIdx)
	if err != nil {
		return nil, err
	}

	allowedFrom := make(map[int]bool)
	if len(cfg.Controllable) == 0 {
		for i := 0; i < n; i++ {
			if i != targetIdx {
				allowedFrom[i] = true
			}
		}
	} else {
		for _, name := range cfg.Controllable {
			if idx, exists := idxMap[name]; exists && idx != targetIdx {
				allowedFrom[idx] = true
			}
		}
	}

	eps := cfg.Epsilon
	if eps <= 0 || eps >= 1 {
		eps = 0.10
	}

	// Inject mass eps/(1-eps) into (i, j) so after NormalizeRows, P'_i = (1-eps)*P_i + eps*e_j
	injectMass := eps / (1.0 - eps)
	var suggestions []KnobSuggestion

	for i := 0; i < n; i++ {
		if !allowedFrom[i] {
			continue
		}
		for j := 0; j < n; j++ {
			if i == j || k.P.At(i, j) > 0.95 {
				continue
			}

			pk := k.Clone()
			pk.P.Set(i, j, pk.P.At(i, j)+injectMass)
			if err := pk.NormalizeRows(1e-9); err != nil {
				continue
			}

			projMFPT, err := avgMFPTToTarget(pk, targetIdx)
			if err != nil || math.IsInf(projMFPT, 1) || math.IsNaN(projMFPT) {
				continue
			}

			delta := projMFPT - baseMFPT
			if delta < -1e-6 {
				suggestions = append(suggestions, KnobSuggestion{
					FromState:     names[i],
					ToState:       names[j],
					Epsilon:       eps,
					BaselineMFPT:  baseMFPT,
					ProjectedMFPT: projMFPT,
					MFPTDelta:     delta,
				})
			}
		}
	}

	sort.Slice(suggestions, func(a, b int) bool {
		return suggestions[a].MFPTDelta < suggestions[b].MFPTDelta
	})

	if cfg.TopK > 0 && len(suggestions) > cfg.TopK {
		suggestions = suggestions[:cfg.TopK]
	}

	return &DiagnosticReport{
		BoundaryStates: names,
		EntropyRate:    entropy,
		Classes:        classes,
		BaselineMFPT:   baseMFPT,
		Suggestions:    suggestions,
	}, nil
}

func avgMFPTToTarget(k *catrace.Kernel, targetIdx int) (float64, error) {
	n := k.NumStates()
	var sum float64
	var count int
	for i := 0; i < n; i++ {
		if i == targetIdx {
			continue
		}
		val, err := k.MeanFirstPassage(i, targetIdx)
		if err != nil {
			// Skip states that cannot reach targetIdx (singular linear system for trapped classes)
			continue
		}
		if !math.IsInf(val, 1) && !math.IsNaN(val) {
			sum += val
			count++
		}
	}
	if count == 0 {
		return math.Inf(1), nil
	}
	return sum / float64(count), nil
}
