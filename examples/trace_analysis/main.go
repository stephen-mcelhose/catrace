package main

// Scenario: LLM Agent With Hidden Support System
//
// We observe a single LLM agent's visible task state, while the rest of
// the support system (backend services, knowledge base, human reviewers,
// other agents) operates in the background.
//
// The full parent system has four states:
//
//   Agent A states:
//     - "A_valid"   : the focal agent is operating correctly
//     - "A_invalid" : the focal agent is in an error or degraded state
//
//   System B states (hidden from the observer):
//     - "B_valid"   : the support system is healthy
//     - "B_invalid" : the support system is degraded
//
// We only care about the focal agent A's health states. We trace the
// full 4-state parent kernel onto the subset {A_valid, A_invalid}.
//
// What the trace means:
//   The trace kernel does NOT mean "just delete the B rows and columns."
//   It means: "what are the effective health dynamics of agent A after
//   integrating out all the hidden excursions through system B?"
//
//   In plain English:
//     - If agent A is valid now, what is the probability that the next
//       time we SEE A again in the observed subset, it is valid vs invalid?
//     - That "next observed return" includes hidden detours through B
//       that we cannot see but whose effect is folded into the probabilities.
//
//   This is the interpretive power of the trace:
//     - the states keep their original meaning (valid / invalid)
//     - the hidden system's influence is absorbed into the transition numbers
//     - the long-run behavior is consistent with the parent restricted
//       to the observed subset

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func main() {
	parent, err := catrace.NewKernel(mat.NewDense(4, 4, []float64{
		0.60, 0.20, 0.10, 0.10,
		0.15, 0.55, 0.15, 0.15,
		0.20, 0.20, 0.40, 0.20,
		0.10, 0.20, 0.20, 0.50,
	}), []string{"agent: ok", "agent: error", "support: ok", "support: error"})
	if err != nil {
		log.Fatal(err)
	}

	subset := []int{0, 1}
	trace, err := parent.Trace(subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	ok, err := trace.IsTraceOf(parent, subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piParent, err := parent.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	piTrace, err := trace.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	restricted, err := catrace.RestrictDistribution(piParent, subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	classes, err := trace.Classes(1e-12)
	if err != nil {
		log.Fatal(err)
	}

	rng := rand.New(rand.NewSource(7))
	trajectory := make([]int, 100)
	trajectory[0] = 0
	for i := 1; i < len(trajectory); i++ {
		next, err := parent.Sample(trajectory[i-1], rng)
		if err != nil {
			log.Fatal(err)
		}
		trajectory[i] = next
	}
	subsetSet := map[int]bool{0: true, 1: true}
	traceSeq := catrace.SampleTraceFromSequence(trajectory, subsetSet)
	est, err := catrace.EstimateKernelFromSequence(traceSeq, 2, 1e-3)
	if err != nil {
		log.Fatal(err)
	}
	windows, err := catrace.WindowedTraceEstimates(trajectory, subsetSet, 25, 10, 1e-3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parent kernel")
	fmt.Printf("%v\n\n", mat.Formatted(parent.P, mat.Prefix("  ")))
	fmt.Println("Trace kernel on subset {0,1}")
	fmt.Printf("%v\n\n", mat.Formatted(trace.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf = %v\n", ok)
	fmt.Printf("stationary(parent)|subset normalized = %.6f %.6f\n", restricted[0], restricted[1])
	fmt.Printf("stationary(trace)                    = %.6f %.6f\n", piTrace[0], piTrace[1])
	fmt.Printf("recurrent classes(trace) = %v\n", classes.Recurrent)
	if est.Kernel != nil {
		fmt.Println("Estimated trace kernel from sampled sequence")
		fmt.Printf("%v\n", mat.Formatted(est.Kernel.P, mat.Prefix("  ")))
	}
	fmt.Printf("windowed estimates computed: %d\n", len(windows))

	for _, kv := range []struct {
		k     *catrace.Kernel
		title string
		file  string
	}{
		{parent, "Trace Analysis — parent kernel (4 states)", "trace_analysis_parent.html"},
		{trace, "Trace Analysis — trace kernel on {agent: ok, agent: error}", "trace_analysis_trace.html"},
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
