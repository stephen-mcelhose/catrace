package main

// Scenario: Routing Agent
//
// A customer support system receives tickets of genuinely different types —
// billing disputes, technical faults, and general enquiries — but the router
// agent can only read the ticket text, not the underlying truth.
// Misclassification sends the ticket to the wrong specialist, who cannot
// resolve it; the ticket re-enters the queue as an apparent different type.
//
// World states (W) = true nature of the incoming ticket
//   - "billing_ticket"   : genuine billing dispute
//   - "technical_ticket" : genuine technical fault
//   - "general_ticket"   : genuine general enquiry
//
// Experience states (X) = the router agent's perceived classification
//   - "reads_billing"   : router thinks this is a billing issue
//   - "reads_technical" : router thinks this is a technical issue
//   - "reads_general"   : router thinks this is a general enquiry
//
// Action states (G) = where the router sends the ticket
//   - "route_billing"   : send to billing specialist
//   - "route_technical" : send to technical specialist
//   - "route_general"   : send to general support
//   - "escalate_human"  : escalate to human triage agent
//
// The world kernel W = P·D·A captures the full ticket-flow transition matrix,
// including misrouting loops where a ticket sent to the wrong specialist
// re-enters the queue as a different apparent type.
// The stationary distribution shows how much of the queue is occupied by
// tickets that have already been misrouted at least once.

import (
	"fmt"
	"log"
	"os"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// P: W => X  (3×3, perception — classification accuracy)
	//
	// Diagonal entries are correct classifications; off-diagonal are misreads.
	// A billing ticket is mostly perceived correctly (0.75) but sometimes
	// misread as technical (0.15) or general (0.10).
	P := mat.NewDense(3, 3, []float64{
		0.75, 0.15, 0.10, // billing_ticket  → reads_billing / reads_technical / reads_general
		0.10, 0.80, 0.10, // technical_ticket
		0.10, 0.15, 0.75, // general_ticket
	})

	// D: X => G  (3×4, decision — routing policy)
	//
	// The router mostly routes to its perceived classification.
	// It escalates to human triage ~10% of the time regardless of perceived type.
	D := mat.NewDense(3, 4, []float64{
		0.80, 0.05, 0.05, 0.10, // reads_billing   → route_billing / route_technical / route_general / escalate
		0.05, 0.80, 0.05, 0.10, // reads_technical
		0.05, 0.05, 0.80, 0.10, // reads_general
	})

	// A: G => W  (4×3, action effect)
	//
	// Correct routing resolves the ticket; the next ticket is drawn from the
	// true-type pool (roughly proportional to incoming mix).
	// Wrong routing sends the ticket to the wrong specialist; the specialist
	// cannot resolve it and re-queues it, often as an apparent different type.
	// Escalation to human triage resolves any ticket type with roughly equal
	// probability — the human handles it regardless of type.
	A := mat.NewDense(4, 3, []float64{
		0.70, 0.20, 0.10, // route_billing   — mostly billing stays billing; leakage to technical/general
		0.15, 0.70, 0.15, // route_technical
		0.10, 0.15, 0.75, // route_general
		0.33, 0.34, 0.33, // escalate_human  — human resolves; next ticket roughly uniform
	})

	agent := &catrace.Agent{
		P:      P,
		D:      D,
		A:      A,
		WNames: []string{"billing_ticket", "technical_ticket", "general_ticket"},
		XNames: []string{"reads_billing", "reads_technical", "reads_general"},
		GNames: []string{"route_billing", "route_technical", "route_general", "escalate_human"},
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

	// Mean first passage: how long does a billing ticket take to cycle to
	// technical (billing→technical), and how long does a technical ticket
	// take to reach general (technical→general)?
	// These capture misrouting loop latency.
	mfptBillingToTechnical, err := W.MeanFirstPassage(0, 1)
	if err != nil {
		log.Fatal(err)
	}
	mfptTechnicalToGeneral, err := W.MeanFirstPassage(1, 2)
	if err != nil {
		log.Fatal(err)
	}

	// Trace onto {billing_ticket, technical_ticket} — collapse general.
	// Shows the effective billing/technical dynamics when general enquiries
	// are treated as a hidden background process.
	subsetBT := []int{0, 1}
	traceBT, err := W.Trace(subsetBT, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	okTrace, err := traceBT.IsTraceOf(W, subsetBT, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piTraceBT, err := traceBT.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	restrictedBT, err := catrace.RestrictDistribution(piW, subsetBT, 1e-12)
	if err != nil {
		log.Fatal(err)
	}

	// ── Output ────────────────────────────────────────────────────────────

	fmt.Println("=== World kernel W = P·D·A ===")
	fmt.Printf("%v\n\n", mat.Formatted(W.P, mat.Prefix("  ")))

	fmt.Println("=== Qualia kernel Q = D·A·P ===")
	fmt.Printf("%v\n\n", mat.Formatted(Q.P, mat.Prefix("  ")))

	fmt.Printf("stationary(W)        = %.6f  %.6f  %.6f\n", piW[0], piW[1], piW[2])
	fmt.Printf("stationary(Q)        = %.6f  %.6f  %.6f\n", piQ[0], piQ[1], piQ[2])
	fmt.Printf("entropy_rate(W)      = %.6f bits/step\n", hW)
	fmt.Printf("recurrent classes(W) = %v\n", classesW.Recurrent)
	fmt.Printf("transient states(W)  = %v\n\n", classesW.Transient)

	fmt.Printf("MFPT billing→technical  = %.6f steps\n", mfptBillingToTechnical)
	fmt.Printf("MFPT technical→general  = %.6f steps\n\n", mfptTechnicalToGeneral)

	fmt.Println("=== Trace onto {billing_ticket, technical_ticket} — collapse general ===")
	fmt.Printf("%v\n\n", mat.Formatted(traceBT.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf                              = %v\n", okTrace)
	fmt.Printf("stationary(W)|{billing,technical} norm = %.6f  %.6f\n", restrictedBT[0], restrictedBT[1])
	fmt.Printf("stationary(trace)                      = %.6f  %.6f\n", piTraceBT[0], piTraceBT[1])

	for _, kv := range []struct {
		k     *catrace.Kernel
		title string
		file  string
	}{
		{W, "Routing Agent — world kernel W (3 ticket types)", "routing_W.html"},
		{traceBT, "Routing Agent — trace onto {billing, technical}", "routing_trace.html"},
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
