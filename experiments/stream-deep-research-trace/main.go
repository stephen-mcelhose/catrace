// Command stream-deep-research-trace runs Experiments A/B/C against
// testdata/events_jsonl.jsonl (harveyai/deep-research-starter sample).
package main

import (
	"fmt"
	"math"
	"os"

	"github.com/stephen-mcelhose/catrace/stream"
)

var (
	keepStates = []string{
		"response.created",
		"response.in_progress",
		"response.completed",
		"response.output_item.added",
		"response.output_item.done",
		"response.web_search_call.in_progress",
		"response.web_search_call.searching",
		"response.web_search_call.completed",
		"response.content_part.added",
		"response.content_part.done",
		"response.output_text.done",
		"response.reasoning_summary_part.added",
		"response.reasoning_summary_part.done",
		"response.reasoning_summary_text.done",
	}
	knobs = []string{
		"response.output_item.done",
		"response.web_search_call.completed",
		"response.reasoning_summary_text.done",
	}
	target = "response.completed"
)

func main() {
	path := "testdata/events_jsonl.jsonl"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open %s: %v\n(hint: curl the harvey sample into testdata/)\n", path, err)
		os.Exit(1)
	}
	defer f.Close()

	est := stream.NewEstimator(0)
	mon := stream.NewMonitor(25, 250)
	var driftSamples []float64
	var driftFailed, driftNaN, n int
	_, skipped, err := stream.IngestJSONL(f, func(state string) {
		est.ObserveNext(state)
		mon.ObserveNext(state)
		n++
		if n%50 != 0 {
			return
		}
		// Mid-stream Trace singularity is signal (boundary ill-posed / phase break),
		// not a gap to drop — record as a peak-severity drift event.
		const traceFailSentinel = 1.0
		score, err := mon.DriftScore(keepStates, 1e-12)
		if math.IsNaN(score) {
			driftNaN++
			return
		}
		if err != nil {
			driftFailed++
			driftSamples = append(driftSamples, traceFailSentinel)
			return
		}
		driftSamples = append(driftSamples, score)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ingest: %v\n", err)
		os.Exit(1)
	}

	full, err := est.Kernel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "kernel: %v\n", err)
		os.Exit(1)
	}
	hFull, _ := full.EntropyRate(2)
	classesFull, err := full.Classes(1e-12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "classes(full): %v\n", err)
		os.Exit(1)
	}

	traced, err := stream.TraceByName(full, keepStates, 1e-12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "trace: %v\n", err)
		os.Exit(1)
	}
	validateErr := traced.Validate(1e-9)
	hTraced, _ := traced.EntropyRate(2)
	classesTraced, err := traced.Classes(1e-12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "classes(traced): %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Experiment A: Full vs Traced ===")
	fmt.Printf("events=%d skipped_lines=%d\n", n, skipped)
	fmt.Printf("N_full=%d  |S|=%d\n", full.NumStates(), traced.NumStates())
	fmt.Printf("H_full=%.6f bits  H_traced=%.6f bits\n", hFull, hTraced)
	fmt.Printf("Validate(traced,1e-9): %v\n", validateOK(validateErr))
	fmt.Printf("classes_full: recurrent=%d transient=%d\n", len(classesFull.Recurrent), len(classesFull.Transient))
	fmt.Printf("classes_traced: recurrent=%d transient=%d names=%v\n",
		len(classesTraced.Recurrent), len(classesTraced.Transient), traced.StateNames)

	fmt.Println("\n=== Experiment B: AnalyzeAndSuggest (ε sweep) ===")
	for _, eps := range []float64{0.05, 0.15, 0.30} {
		rep, err := stream.AnalyzeAndSuggest(traced, stream.SuggestConfig{
			TargetState:  target,
			Controllable: knobs,
			Epsilon:      eps,
			TopK:         3,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "suggest eps=%.2f: %v\n", eps, err)
			continue
		}
		fmt.Printf("\neps=%.2f  baseline_mfpt=%.4f\n", eps, rep.BaselineMFPT)
		for i, s := range rep.Suggestions {
			speedup := 0.0
			if s.BaselineMFPT != 0 {
				speedup = 100 * (s.BaselineMFPT - s.ProjectedMFPT) / s.BaselineMFPT
			}
			fmt.Printf("  %d. %s -> %s  base=%.2f proj=%.2f delta=%.2f speedup=%.1f%%\n",
				i+1, s.FromState, s.ToState, s.BaselineMFPT, s.ProjectedMFPT, s.MFPTDelta, speedup)
		}
	}

	fmt.Println("\n=== Experiment C: DriftScore profile ===")
	fmt.Printf("samples=%d  trace_fail_as_signal=%d  nan=%d\n",
		len(driftSamples), driftFailed, driftNaN)
	if len(driftSamples) == 0 {
		fmt.Println("no drift samples")
		return
	}
	min, max, sum := driftSamples[0], driftSamples[0], 0.0
	for _, v := range driftSamples {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	mean := sum / float64(len(driftSamples))
	fmt.Printf("min=%.6f  mean=%.6f  peak=%.6f  nan_observed=%v\n",
		min, mean, max, driftNaN > 0)
}

func validateOK(err error) string {
	if err == nil {
		return "PASS"
	}
	return "FAIL: " + err.Error()
}
