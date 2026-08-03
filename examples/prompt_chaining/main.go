package main

// Scenario: Prompt-Chaining Pipeline
//
// A document intelligence pipeline processes raw text through three specialist
// stages in sequence: extractor, summariser, formatter. Between stages, a
// programmatic gate checks output quality; if it fails, the stage reruns
// before passing the result on.
//
// World states (W) = pipeline stage the document has actually reached
//   - "raw"       : unprocessed input
//   - "extracted" : key claims identified
//   - "summarised": structured brief produced
//   - "formatted" : delivery-ready report
//   - "failed"    : pipeline abandoned via escalation
//
// Experience states (X) = stage agent's perception of its input quality
//   - "input_clear" : agent reads input as processable
//   - "input_noisy" : agent reads input as degraded or unclear
//
// Action states (G) = what the stage agent does
//   - "process"  : attempt the transformation
//   - "retry"    : reframe and try again
//   - "escalate" : abandon this document
//
// Composed world kernel W = P·D·A captures the full stage-to-stage
// transition probabilities including retry loops and failure exits.
// Tracing onto {raw, formatted, failed} collapses intermediate stages
// and reveals the pipeline's end-to-end pass/fail rate.

import (
	"fmt"
	"log"
	"os"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// P: W => X  (5×2, perception)
	//
	// How reliably does the stage agent read the quality of its input?
	//   raw       : mostly noisy   — unprocessed documents are hard to assess
	//   extracted : mixed          — extraction gives some signal
	//   summarised: mostly clear   — structured brief is easier to evaluate
	//   formatted : very clear     — near-final output is legible
	//   failed    : almost always noisy — failure state produces no useful signal
	P := mat.NewDense(5, 2, []float64{
		0.25, 0.75, // raw
		0.60, 0.40, // extracted
		0.80, 0.20, // summarised
		0.90, 0.10, // formatted
		0.15, 0.85, // failed
	})

	// D: X => G  (2×3, decision)
	//
	// Given its read of input quality, what does the agent do?
	//   input_clear : mostly process, occasionally retry, rarely escalate
	//   input_noisy : mix of retry and escalate, sometimes process anyway
	D := mat.NewDense(2, 3, []float64{
		0.80, 0.15, 0.05, // input_clear
		0.25, 0.50, 0.25, // input_noisy
	})

	// A: G => W  (3×5, action effect)
	//
	// How does each action affect the pipeline stage?
	//   process  : tends to advance toward later stages
	//   retry    : tends to stay in early/middle stages; sometimes advances
	//   escalate : mostly sends to failed
	A := mat.NewDense(3, 5, []float64{
		0.05, 0.30, 0.35, 0.25, 0.05, // process
		0.30, 0.40, 0.20, 0.05, 0.05, // retry
		0.05, 0.05, 0.05, 0.05, 0.80, // escalate
	})

	agent := &catrace.Agent{
		P:      P,
		D:      D,
		A:      A,
		WNames: []string{"raw", "extracted", "summarised", "formatted", "failed"},
		XNames: []string{"input_clear", "input_noisy"},
		GNames: []string{"process", "retry", "escalate"},
	}

	W, err := agent.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}
	Q, err := agent.QualiaKernel()
	if err != nil {
		log.Fatal(err)
	}

	piW, err := W.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	hW, err := W.EntropyRate(2)
	if err != nil {
		log.Fatal(err)
	}
	classesW, err := W.Classes(1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piQ, err := Q.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}

	// Trace onto {raw, formatted, failed} — end-to-end pipeline view.
	// Indices: raw=0, formatted=3, failed=4.
	endStates := []int{0, 3, 4}
	traceEnd, err := W.Trace(endStates, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	okTrace, err := traceEnd.IsTraceOf(W, endStates, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piTrace, err := traceEnd.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	restricted, err := catrace.RestrictDistribution(piW, endStates, 1e-12)
	if err != nil {
		log.Fatal(err)
	}

	// ── Output ────────────────────────────────────────────────────────────

	fmt.Println("=== World kernel W = P·D·A ===")
	fmt.Printf("%v\n\n", mat.Formatted(W.P, mat.Prefix("  ")))

	fmt.Println("=== Qualia kernel Q = D·A·P ===")
	fmt.Printf("%v\n\n", mat.Formatted(Q.P, mat.Prefix("  ")))

	fmt.Printf("stationary(W)        = ")
	for i, v := range piW {
		fmt.Printf("%.6f", v)
		if i < len(piW)-1 {
			fmt.Printf("  ")
		}
	}
	fmt.Println()
	fmt.Printf("stationary(Q)        = %.6f  %.6f\n", piQ[0], piQ[1])
	fmt.Printf("entropy_rate(W)      = %.6f bits/step\n", hW)
	fmt.Printf("recurrent classes(W) = %v\n", classesW.Recurrent)
	fmt.Printf("transient states(W)  = %v\n\n", classesW.Transient)

	fmt.Println("=== Trace onto {raw, formatted, failed} — end-to-end view ===")
	fmt.Printf("%v\n\n", mat.Formatted(traceEnd.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf                        = %v\n", okTrace)
	fmt.Printf("stationary(W)|{raw,fmt,fail} norm= ")
	for i, v := range restricted {
		fmt.Printf("%.6f", v)
		if i < len(restricted)-1 {
			fmt.Printf("  ")
		}
	}
	fmt.Println()
	fmt.Printf("stationary(trace)               = ")
	for i, v := range piTrace {
		fmt.Printf("%.6f", v)
		if i < len(piTrace)-1 {
			fmt.Printf("  ")
		}
	}
	fmt.Println()

	for _, kv := range []struct {
		k     *catrace.Kernel
		title string
		file  string
	}{
		{W, "Prompt Chaining — world kernel W (5 pipeline stages)", "prompt_chaining_W.html"},
		{traceEnd, "Prompt Chaining — trace onto {raw, formatted, failed}", "prompt_chaining_trace.html"},
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
