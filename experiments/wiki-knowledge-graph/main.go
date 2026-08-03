package main

// Experiment: wiki knowledge graph — trace chain importance correction.
//
// Applies catrace's Kernel.Trace to the catrace wiki's own wikilink graph to
// correct the importance distortion caused by 14 planned-but-unwritten pages.
//
// Graph A (14 nodes): existing wiki pages only — naive PageRank baseline.
// Graph B (28 nodes): existing + 14 planned pages. Planned pages are the
// hidden subset; their influence is folded back via the trace chain.
//
// Key simplification confirmed: planned pages have no planned→planned edges,
// so c = 0 and the trace formula reduces to L_A = a + b·d.
//
// Run from repo root:
//
//	go run experiments/wiki-knowledge-graph/main.go

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

// existing page indices (0–13) — ordered by naive PageRank rank for readability.
const (
	selfHealingNodes         = 0
	markovChainFoundations   = 1
	pdaTripletModel          = 2
	jointKernelsAndCoupling  = 3
	traceChain               = 4
	agenticPatternsCatalogue = 5
	validatorRepair          = 6
	catraceAPI               = 7
	simpleAgent              = 8
	variantComparison        = 9
	devWorkflowPatterns      = 10
	hiddenSupportSystem      = 11
	structuralPatterns       = 12
	catraceGlossary          = 13
)

// planned page indices in Graph B (14–27).
const (
	promptChaining         = 14
	routing                = 15
	parallelisation        = 16
	orchestratorWorkers    = 17
	supervisor             = 18
	swarm                  = 19
	blackboard             = 20
	debate                 = 21
	planAndExecute         = 22
	humanInTheLoop         = 23
	d1ResearchPlanImpl     = 24
	d2ImplementVerify      = 25
	d3ImplementCritic      = 26
	d4PlanImplCriticVerify = 27
)

var existingNames = []string{
	"Example: Self-Healing Nodes",   // 0
	"Markov Chain Foundations",       // 1
	"PDA Triplet Model",              // 2
	"Joint Kernels and Coupling",     // 3
	"Trace Chain",                    // 4
	"Agentic Patterns Catalogue",     // 5
	"Example: Validator Repair",      // 6
	"catrace API",                    // 7
	"Example: Simple Agent",          // 8
	"Variant Comparison Methodology", // 9
	"Dev-Workflow Patterns",          // 10
	"Example: Hidden Support System", // 11
	"Structural Patterns",            // 12
	"catrace Glossary",               // 13
}

var plannedNames = []string{
	"Example: Prompt Chaining",                 // 14
	"Example: Routing",                          // 15
	"Example: Parallelisation",                  // 16
	"Example: Orchestrator-Workers",             // 17
	"Example: Supervisor",                       // 18
	"Example: Swarm",                            // 19
	"Example: Blackboard",                       // 20
	"Example: Debate",                           // 21
	"Example: Plan-and-Execute",                 // 22
	"Example: Human-in-the-Loop",               // 23
	"Example: D1 Research-Plan-Implement",       // 24
	"Example: D2 Implement-Verify",             // 25
	"Example: D3 Implement-Critic",             // 26
	"Example: D4 Plan-Implement-Critic-Verify", // 27
}

// link sets adj[from][to] = 1 for each destination index.
func link(adj [][]float64, from int, to ...int) {
	for _, t := range to {
		adj[from][t] = 1
	}
}

// buildGraphA returns the 14×14 binary adjacency matrix for the existing wiki
// pages. Edges are derived from [[wikilink]] references in docs/wiki/*.md.
func buildGraphA() *mat.Dense {
	n := 14
	adj := make([][]float64, n)
	for i := range adj {
		adj[i] = make([]float64, n)
	}

	// Wikilink out-edges parsed from docs/wiki/*.md (one entry per source page).
	link(adj, selfHealingNodes,
		markovChainFoundations, jointKernelsAndCoupling, traceChain,
		agenticPatternsCatalogue, catraceAPI, variantComparison)

	link(adj, markovChainFoundations,
		selfHealingNodes, pdaTripletModel, traceChain, devWorkflowPatterns)

	link(adj, pdaTripletModel,
		simpleAgent, jointKernelsAndCoupling)

	link(adj, jointKernelsAndCoupling,
		selfHealingNodes, validatorRepair, pdaTripletModel)

	link(adj, traceChain,
		selfHealingNodes, markovChainFoundations, validatorRepair, hiddenSupportSystem)

	link(adj, agenticPatternsCatalogue,
		selfHealingNodes, pdaTripletModel, jointKernelsAndCoupling, traceChain,
		validatorRepair, catraceAPI, simpleAgent, variantComparison,
		devWorkflowPatterns, structuralPatterns)

	link(adj, validatorRepair,
		markovChainFoundations, jointKernelsAndCoupling, traceChain,
		agenticPatternsCatalogue, catraceAPI)

	link(adj, catraceAPI,
		markovChainFoundations, pdaTripletModel, traceChain,
		agenticPatternsCatalogue, catraceGlossary)

	link(adj, simpleAgent,
		markovChainFoundations, pdaTripletModel, agenticPatternsCatalogue, catraceAPI)

	// variant-comparison-methodology has exactly one out-edge — the concentration
	// artifact that inflates Example: Self-Healing Nodes in naive PageRank.
	link(adj, variantComparison, selfHealingNodes)

	link(adj, devWorkflowPatterns,
		markovChainFoundations, variantComparison, structuralPatterns)

	link(adj, hiddenSupportSystem,
		markovChainFoundations, traceChain, agenticPatternsCatalogue, catraceAPI)

	link(adj, structuralPatterns,
		selfHealingNodes, pdaTripletModel, validatorRepair, simpleAgent,
		variantComparison, devWorkflowPatterns)

	link(adj, catraceGlossary,
		markovChainFoundations, pdaTripletModel, jointKernelsAndCoupling)

	return denseFrom2D(adj)
}

// buildGraphB returns the 28×28 binary adjacency matrix for the full intended
// wiki. Rows 0–13 are the existing pages (same links as Graph A, extended
// where noted). Rows 14–27 are the planned pages; all their out-edges point
// only to existing pages (confirming c = 0).
func buildGraphB() *mat.Dense {
	n := 28
	adj := make([][]float64, n)
	for i := range adj {
		adj[i] = make([]float64, n)
	}

	// Existing pages — identical to Graph A ...
	link(adj, selfHealingNodes,
		markovChainFoundations, jointKernelsAndCoupling, traceChain,
		agenticPatternsCatalogue, catraceAPI, variantComparison)

	link(adj, markovChainFoundations,
		selfHealingNodes, pdaTripletModel, traceChain, devWorkflowPatterns)

	link(adj, pdaTripletModel,
		simpleAgent, jointKernelsAndCoupling)

	link(adj, jointKernelsAndCoupling,
		selfHealingNodes, validatorRepair, pdaTripletModel)

	link(adj, traceChain,
		selfHealingNodes, markovChainFoundations, validatorRepair, hiddenSupportSystem)

	link(adj, agenticPatternsCatalogue,
		selfHealingNodes, pdaTripletModel, jointKernelsAndCoupling, traceChain,
		validatorRepair, catraceAPI, simpleAgent, variantComparison,
		devWorkflowPatterns, structuralPatterns)

	link(adj, validatorRepair,
		markovChainFoundations, jointKernelsAndCoupling, traceChain,
		agenticPatternsCatalogue, catraceAPI)

	link(adj, catraceAPI,
		markovChainFoundations, pdaTripletModel, traceChain,
		agenticPatternsCatalogue, catraceGlossary)

	link(adj, simpleAgent,
		markovChainFoundations, pdaTripletModel, agenticPatternsCatalogue, catraceAPI)

	link(adj, variantComparison, selfHealingNodes)

	// dev-workflow-patterns: same as Graph A, plus the 4 planned D-pages.
	link(adj, devWorkflowPatterns,
		markovChainFoundations, variantComparison, structuralPatterns,
		d1ResearchPlanImpl, d2ImplementVerify, d3ImplementCritic, d4PlanImplCriticVerify)

	link(adj, hiddenSupportSystem,
		markovChainFoundations, traceChain, agenticPatternsCatalogue, catraceAPI)

	// structural-patterns: same as Graph A, plus the 10 planned structural pages.
	link(adj, structuralPatterns,
		selfHealingNodes, pdaTripletModel, validatorRepair, simpleAgent,
		variantComparison, devWorkflowPatterns,
		promptChaining, routing, parallelisation, orchestratorWorkers, supervisor,
		swarm, blackboard, debate, planAndExecute, humanInTheLoop)

	link(adj, catraceGlossary,
		markovChainFoundations, pdaTripletModel, jointKernelsAndCoupling)

	// Planned pages (14–27): all out-edges point to existing pages only.
	// This confirms c = 0 ⟹ (I-c)⁻¹ = I ⟹ trace = a + b·d.
	link(adj, promptChaining,
		markovChainFoundations, pdaTripletModel, catraceAPI, traceChain, agenticPatternsCatalogue)

	link(adj, routing,
		markovChainFoundations, pdaTripletModel, catraceAPI, agenticPatternsCatalogue)

	link(adj, parallelisation,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		jointKernelsAndCoupling, agenticPatternsCatalogue)

	link(adj, orchestratorWorkers,
		markovChainFoundations, pdaTripletModel, catraceAPI, agenticPatternsCatalogue)

	link(adj, supervisor,
		markovChainFoundations, pdaTripletModel, catraceAPI, agenticPatternsCatalogue)

	link(adj, swarm,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		jointKernelsAndCoupling, agenticPatternsCatalogue)

	link(adj, blackboard,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		jointKernelsAndCoupling, agenticPatternsCatalogue)

	link(adj, debate,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		jointKernelsAndCoupling, agenticPatternsCatalogue)

	link(adj, planAndExecute,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		traceChain, agenticPatternsCatalogue)

	link(adj, humanInTheLoop,
		markovChainFoundations, pdaTripletModel, catraceAPI, agenticPatternsCatalogue)

	link(adj, d1ResearchPlanImpl,
		markovChainFoundations, pdaTripletModel, catraceAPI, traceChain,
		agenticPatternsCatalogue, variantComparison)

	link(adj, d2ImplementVerify,
		markovChainFoundations, catraceAPI, agenticPatternsCatalogue, variantComparison)

	link(adj, d3ImplementCritic,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		jointKernelsAndCoupling, agenticPatternsCatalogue, variantComparison)

	link(adj, d4PlanImplCriticVerify,
		markovChainFoundations, pdaTripletModel, catraceAPI,
		jointKernelsAndCoupling, agenticPatternsCatalogue, variantComparison)

	return denseFrom2D(adj)
}

func denseFrom2D(adj [][]float64) *mat.Dense {
	n := len(adj)
	flat := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			flat[i*n+j] = adj[i][j]
		}
	}
	return mat.NewDense(n, n, flat)
}

type pageEntry struct {
	name      string
	piNaive   float64
	piTrace   float64
	rankNaive int
	rankTrace int
}

func main() {
	const tol = 1e-12

	// ── Graph A: naive 14-node kernel ────────────────────────────────────────
	naive, err := catrace.NewRandomWalkKernel(buildGraphA(), existingNames)
	if err != nil {
		log.Fatalf("build Graph A: %v", err)
	}

	piNaive, err := naive.Stationary(tol, 10000)
	if err != nil {
		log.Fatalf("Graph A stationary: %v", err)
	}

	// ── Graph B: full 28-node kernel ─────────────────────────────────────────
	allNames := append(existingNames, plannedNames...)
	full, err := catrace.NewRandomWalkKernel(buildGraphB(), allNames)
	if err != nil {
		log.Fatalf("build Graph B: %v", err)
	}

	// Observed subset = all 14 existing pages (indices 0–13).
	existingSubset := make([]int, 14)
	for i := range existingSubset {
		existingSubset[i] = i
	}

	// ── Trace ─────────────────────────────────────────────────────────────────
	traced, err := full.Trace(existingSubset, tol)
	if err != nil {
		log.Fatalf("Trace: %v", err)
	}

	piTrace, err := traced.Stationary(tol, 10000)
	if err != nil {
		log.Fatalf("trace stationary: %v", err)
	}

	// ── Verify ────────────────────────────────────────────────────────────────
	isTrace, err := traced.IsTraceOf(full, existingSubset, tol)
	if err != nil {
		log.Fatalf("IsTraceOf: %v", err)
	}

	// ── Build rank table ──────────────────────────────────────────────────────
	entries := make([]pageEntry, 14)
	for i := range entries {
		entries[i] = pageEntry{
			name:    existingNames[i],
			piNaive: piNaive[i],
			piTrace: piTrace[i],
		}
	}

	byNaive := rankOrder(entries, func(e pageEntry) float64 { return e.piNaive })
	byTrace := rankOrder(entries, func(e pageEntry) float64 { return e.piTrace })
	for rank, idx := range byNaive {
		entries[idx].rankNaive = rank + 1
	}
	for rank, idx := range byTrace {
		entries[idx].rankTrace = rank + 1
	}

	// ── Print results ─────────────────────────────────────────────────────────
	sep := strings.Repeat("─", 80)
	fmt.Println(sep)
	fmt.Println("Wiki Knowledge Graph — Trace Chain Importance Correction")
	fmt.Println(sep)
	fmt.Println()
	fmt.Printf("%-42s  %8s  %8s  %6s  %6s  %8s\n",
		"Page", "Naive π", "Trace π", "Rank A", "Rank B", "Δ rank")
	fmt.Printf("%-42s  %8s  %8s  %6s  %6s  %8s\n",
		strings.Repeat("-", 42), strings.Repeat("-", 8),
		strings.Repeat("-", 8), strings.Repeat("-", 6),
		strings.Repeat("-", 6), strings.Repeat("-", 8))

	for _, idx := range byNaive {
		e := entries[idx]
		delta := e.rankNaive - e.rankTrace // positive = moved up in rank
		var dStr string
		switch {
		case delta > 0:
			dStr = fmt.Sprintf("+%d", delta)
		case delta < 0:
			dStr = fmt.Sprintf("%d", delta)
		default:
			dStr = "0"
		}
		fmt.Printf("%-42s  %8.4f  %8.4f  %6d  %6d  %8s\n",
			e.name, e.piNaive, e.piTrace, e.rankNaive, e.rankTrace, dStr)
	}

	fmt.Println()
	fmt.Printf("IsTraceOf = %v\n", isTrace)
	fmt.Println()

	// Verdict evaluation
	var sPatterns, devWorkflow, selfHeal pageEntry
	for _, e := range entries {
		switch e.name {
		case "Structural Patterns":
			sPatterns = e
		case "Dev-Workflow Patterns":
			devWorkflow = e
		case "Example: Self-Healing Nodes":
			selfHeal = e
		}
	}

	c1 := sPatterns.rankNaive-sPatterns.rankTrace >= 5
	c2 := devWorkflow.rankNaive-devWorkflow.rankTrace >= 3
	c3 := selfHeal.rankTrace-selfHeal.rankNaive >= 2
	c4 := isTrace

	fmt.Println("── Verdict criteria ──────────────────────────────────────")
	fmt.Printf("C1  Structural Patterns up ≥5 ranks:    %v  (Δ = %+d)\n",
		c1, sPatterns.rankNaive-sPatterns.rankTrace)
	fmt.Printf("C2  Dev-Workflow Patterns up ≥3 ranks:  %v  (Δ = %+d)\n",
		c2, devWorkflow.rankNaive-devWorkflow.rankTrace)
	fmt.Printf("C3  Self-Healing Nodes down ≥2 ranks:   %v  (Δ = %+d)\n",
		c3, selfHeal.rankTrace-selfHeal.rankNaive)
	fmt.Printf("C4  IsTraceOf = true:                   %v\n", c4)
	fmt.Println()

	met := 0
	for _, b := range []bool{c1, c2, c3, c4} {
		if b {
			met++
		}
	}
	verdict := "not supported"
	if met == 4 {
		verdict = "supported"
	} else if met >= 2 {
		verdict = "trade-off"
	}
	fmt.Printf("Claim: %s  (%d/4 criteria met)\n", verdict, met)

	// ── HTML visualisations ───────────────────────────────────────────────────
	dir := "experiments/wiki-knowledge-graph"

	htmlNaive, err := naive.ToHTML(&catrace.VisualiseOptions{
		Title: "Wiki Knowledge Graph — Naive PageRank (14 nodes)",
	})
	if err != nil {
		log.Fatalf("ToHTML naive: %v", err)
	}
	if err := os.WriteFile(dir+"/naive.html", []byte(htmlNaive), 0o644); err != nil {
		log.Fatalf("write naive.html: %v", err)
	}

	htmlTrace, err := traced.ToHTML(&catrace.VisualiseOptions{
		Title: "Wiki Knowledge Graph — Trace-corrected (14 nodes)",
	})
	if err != nil {
		log.Fatalf("ToHTML trace: %v", err)
	}
	if err := os.WriteFile(dir+"/trace.html", []byte(htmlTrace), 0o644); err != nil {
		log.Fatalf("write trace.html: %v", err)
	}

	fmt.Printf("Wrote %s/naive.html\n", dir)
	fmt.Printf("Wrote %s/trace.html\n", dir)
}

// rankOrder returns indices into entries sorted descending by score(e).
func rankOrder(entries []pageEntry, score func(pageEntry) float64) []int {
	idx := make([]int, len(entries))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool {
		return score(entries[idx[a]]) > score(entries[idx[b]])
	})
	return idx
}
