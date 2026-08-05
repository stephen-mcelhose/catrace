package main

// Scenario: Prompt-chaining document pipeline
//
// Anthropic-style prompt chaining (a workflow, not an autonomous agent):
// three specialised LLM calls in a fixed order, with programmatic gates
// between them.
//
//	raw  --extractor-->  extracted  --summariser-->  summarised  --formatter-->  formatted
//	                              \         |                \         |              /
//	                               \---- gate ----/           \---- gate ----/         /
//	                                          \________________ escalate ______________/--> failed
//
// Each stage has its own perception and decision kernels. Only the active
// stage fires on a given step. The gate is ordinary code (pass / retry_stage /
// escalate), not a fourth LLM.
//
// Pipeline world W (5): raw, extracted, summarised, formatted, failed
//
// The pipeline kernel is assembled row-by-row from the active stage's P·D and
// the gate — not as one shared Agent.WorldKernel() over all stages.

import (
	"fmt"
	"log"
	"os"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

const (
	raw = iota
	extracted
	summarised
	formatted
	failed
	nPipeline
)

// stage is one LLM call in the chain plus the programmatic gate that follows it.
type stage struct {
	name string
	// Perception when this stage is active: P(clear), P(noisy).
	pClear float64
	// Decision D: experience → {emit, retry} (2×2 row-stochastic).
	d *mat.Dense
	// Gate given emit: pass (advance), retryStage (stay), escalate (failed).
	gatePass, gateRetry, gateEscalate float64
	// After a retry action: stay at stage vs escalate.
	retryStay, retryEscalate float64
	xNames, gNames           []string
}

func (s stage) validate() error {
	if s.pClear < 0 || s.pClear > 1 {
		return fmt.Errorf("%s: pClear out of range", s.name)
	}
	if _, err := catrace.NewRectKernel(s.d, s.xNames, s.gNames); err != nil {
		return fmt.Errorf("%s: invalid D: %w", s.name, err)
	}
	if abs(s.gatePass+s.gateRetry+s.gateEscalate-1) > 1e-12 {
		return fmt.Errorf("%s: gate probs must sum to 1", s.name)
	}
	if abs(s.retryStay+s.retryEscalate-1) > 1e-12 {
		return fmt.Errorf("%s: retry probs must sum to 1", s.name)
	}
	return nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func roundMatrix(m *mat.Dense, scale float64) {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m.Set(i, j, float64(int(m.At(i, j)*scale+0.5))/scale)
		}
	}
}

// fillActiveRow writes the pipeline transition row when this stage is active
// at fromState and a gate pass advances to nextState.
func (s stage) fillActiveRow(w *mat.Dense, fromState, nextState int) {
	pNoisy := 1 - s.pClear
	pExp := []float64{s.pClear, pNoisy}
	for x := 0; x < 2; x++ {
		for g := 0; g < 2; g++ {
			mass := pExp[x] * s.d.At(x, g)
			if g == 0 { // emit
				w.Set(fromState, nextState, w.At(fromState, nextState)+mass*s.gatePass)
				w.Set(fromState, fromState, w.At(fromState, fromState)+mass*s.gateRetry)
				w.Set(fromState, failed, w.At(fromState, failed)+mass*s.gateEscalate)
			} else { // retry
				w.Set(fromState, fromState, w.At(fromState, fromState)+mass*s.retryStay)
				w.Set(fromState, failed, w.At(fromState, failed)+mass*s.retryEscalate)
			}
		}
	}
}

func main() {
	pipelineNames := []string{"raw", "extracted", "summarised", "formatted", "failed"}

	extractor := stage{
		name:         "extractor",
		pClear:       0.55,
		d:            mat.NewDense(2, 2, []float64{0.85, 0.15, 0.25, 0.75}),
		gatePass:     0.70,
		gateRetry:    0.20,
		gateEscalate: 0.10,
		retryStay:    0.90,
		retryEscalate: 0.10,
		xNames:       []string{"filing_clear", "filing_noisy"},
		gNames:       []string{"emit_claims", "retry_extract"},
	}
	summariser := stage{
		name:         "summariser",
		pClear:       0.75, // claims already gated once
		d:            mat.NewDense(2, 2, []float64{0.90, 0.10, 0.35, 0.65}),
		gatePass:     0.80,
		gateRetry:    0.15,
		gateEscalate: 0.05,
		retryStay:    0.92,
		retryEscalate: 0.08,
		xNames:       []string{"claims_clear", "claims_noisy"},
		gNames:       []string{"emit_brief", "retry_summarise"},
	}
	formatter := stage{
		name:         "formatter",
		pClear:       0.85,
		d:            mat.NewDense(2, 2, []float64{0.92, 0.08, 0.40, 0.60}),
		gatePass:     0.88,
		gateRetry:    0.08,
		gateEscalate: 0.04,
		retryStay:    0.95,
		retryEscalate: 0.05,
		xNames:       []string{"brief_clear", "brief_noisy"},
		gNames:       []string{"emit_report", "retry_format"},
	}

	for _, s := range []stage{extractor, summariser, formatter} {
		if err := s.validate(); err != nil {
			log.Fatal(err)
		}
	}

	// Assemble pipeline world kernel W (5×5).
	wData := mat.NewDense(nPipeline, nPipeline, nil)
	extractor.fillActiveRow(wData, raw, extracted)
	summariser.fillActiveRow(wData, extracted, summarised)
	formatter.fillActiveRow(wData, summarised, formatted)
	// Shipped report is absorbing. Human queue sometimes returns work to the
	// desk (re-queue) so formatted remains reachable and MFPT is finite.
	wData.Set(formatted, formatted, 1)
	wData.Set(failed, failed, 0.75)
	wData.Set(failed, raw, 0.25)

	roundMatrix(wData, 1e6)

	W, err := catrace.NewKernel(wData, pipelineNames)
	if err != nil {
		log.Fatal(err)
	}

	pi, err := W.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}

	h, err := W.EntropyRate(2)
	if err != nil {
		log.Fatal(err)
	}

	classes, err := W.Classes(1e-12)
	if err != nil {
		log.Fatal(err)
	}

	mfpt, err := W.MeanFirstPassage(raw, formatted)
	if err != nil {
		log.Fatal(err)
	}

	subset := []int{raw, formatted, failed}
	tr, err := W.Trace(subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	ok, err := tr.IsTraceOf(W, subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piTr, err := tr.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Prompt chaining: per-stage kernels ===")
	for _, s := range []stage{extractor, summariser, formatter} {
		fmt.Printf("\n%s  P(clear)=%.2f  P(noisy)=%.2f\n", s.name, s.pClear, 1-s.pClear)
		fmt.Printf("  D (%s → %s):\n%v\n", s.xNames[0]+"/"+s.xNames[1], s.gNames[0]+"/"+s.gNames[1], mat.Formatted(s.d, mat.Prefix("    ")))
		fmt.Printf("  gate|emit: pass=%.2f retry_stage=%.2f escalate=%.2f\n", s.gatePass, s.gateRetry, s.gateEscalate)
		fmt.Printf("  after retry action: stay=%.2f escalate=%.2f\n", s.retryStay, s.retryEscalate)
	}

	fmt.Println("\n=== Pipeline world kernel W (assembled from active stage · gate) ===")
	fmt.Printf("%v\n\n", mat.Formatted(W.P, mat.Prefix("  ")))

	fmt.Println("Stationary distribution π:")
	for i, name := range W.StateNames {
		fmt.Printf("  %-12s %.6f\n", name, pi[i])
	}
	fmt.Println()

	fmt.Printf("Entropy rate H(W): %.6f bits/step\n", h)
	fmt.Printf("Recurrent classes: %v\n", classes.Recurrent)
	fmt.Printf("Transient states:  %v\n", classes.Transient)
	fmt.Printf("MFPT raw → formatted: %.4f steps\n\n", mfpt)

	fmt.Println("=== Trace onto {raw, formatted, failed} ===")
	fmt.Printf("%v\n\n", mat.Formatted(tr.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf verification: %v\n", ok)
	fmt.Printf("Stationary of trace: raw=%.6f  formatted=%.6f  failed=%.6f\n", piTr[0], piTr[1], piTr[2])

	for _, kv := range []struct {
		k     *catrace.Kernel
		title string
		file  string
	}{
		{W, "Prompt Chaining — Pipeline World W", "prompt_chaining_W.html"},
		{tr, "Prompt Chaining — Trace {raw, formatted, failed}", "prompt_chaining_trace.html"},
	} {
		html, err := kv.k.ToHTML(&catrace.VisualiseOptions{Title: kv.title})
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(kv.file, html, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Wrote %s\n", kv.file)
	}
}
