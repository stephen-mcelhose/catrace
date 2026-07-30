package main

// Scenario: Self-adjusting / Self-healing Network Nodes
//
// Inspired by github.com/caseylmanus/workerpool.
//
// A network node monitors its own error rate via an exponential moving average
// (EMA) and throttles itself when errors climb — the same proportional
// backpressure the workerpool applies via its PID controller. An outer
// evolutionary loop watches pool-level throughput and mutates the node's
// configuration when performance drops — the same elitist mutation cycle the
// workerpool uses to search strategy space.
//
// Two agents:
//
//	Node agent    — perceives its own health (EMA error band), decides how
//	                hard to push (push / throttle / idle). Self-throttling is
//	                self-healing: the node recovers by reducing its own load.
//
//	Evolver agent — perceives pool-level score (throughput × success²),
//	                decides to promote (keep best config) or mutate (try a
//	                random variation). When mutation finds a better config,
//	                the node's recovery probabilities improve.
//
// Joint world states (node × evolver), row-major (node first):
//
//	"H·G" : healthy × good_strategy  — node fine, evolver has a good config
//	"H·B" : healthy × bad_strategy   — node fine, evolver has a bad config
//	"D·G" : degraded × good_strategy
//	"D·B" : degraded × bad_strategy
//	"O·G" : overloaded × good_strategy
//	"O·B" : overloaded × bad_strategy  — worst joint state
//
// Two variants are run side by side:
//
//	Variant A — "throttle does it all":
//	  Node A kernel is strong (throttle reliably heals). Mutation boost is
//	  negligible. The evolver's strategy search barely moves the needle.
//
//	Variant B — "evolver matters":
//	  Node A kernel is weak (throttle alone is unreliable). Mutation boost
//	  is large. The node gets stuck degraded without the outer search loop.

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	catrace "github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

// scenario holds the two parameters that differ between variants.
type scenario struct {
	name string
	// nodeA: action → world for the node (3×3, push/throttle/idle → healthy/degraded/overloaded)
	nodeA [3][3]float64
	// mutationBoost: extra probability of healthy added to nodeA rows when evolver mutates
	mutationBoost [3]float64
}

func main() {
	variants := []scenario{
		{
			name: "Variant A — throttle does it all",
			// Throttle is a strong recovery action: self-healing built into the node.
			// Mutation boost is tiny — the evolver barely changes the picture.
			nodeA: [3][3]float64{
				{0.70, 0.25, 0.05}, // push
				{0.65, 0.30, 0.05}, // throttle  ← strong self-heal
				{0.85, 0.13, 0.02}, // idle
			},
			mutationBoost: [3]float64{0.02, 0.02, 0.02}, // negligible
		},
		{
			name: "Variant B — evolver matters",
			// Throttle is weak: the node can't reliably self-heal by slowing down alone.
			// Mutation boost is large — finding a better config is the real recovery path.
			nodeA: [3][3]float64{
				{0.55, 0.35, 0.10}, // push      ← higher failure risk
				{0.45, 0.45, 0.10}, // throttle  ← weak self-heal
				{0.60, 0.35, 0.05}, // idle
			},
			mutationBoost: [3]float64{0.20, 0.25, 0.20}, // significant
		},
	}

	for _, v := range variants {
		fmt.Printf("\n%s\n%s\n", v.name, strings.Repeat("─", len(v.name)))
		runScenario(v)
	}
}

func runScenario(v scenario) {
	// ── Node agent ────────────────────────────────────────────────────────
	//
	// World states: healthy=0, degraded=1, overloaded=2
	// Experience:   ema_low=0, ema_mid=1, ema_high=2
	// Actions:      push=0, throttle=1, idle=2

	nodeP := mat.NewDense(3, 3, []float64{
		0.80, 0.15, 0.05, // healthy    → EMA mostly low
		0.20, 0.55, 0.25, // degraded   → EMA spread
		0.05, 0.20, 0.75, // overloaded → EMA mostly high
	})

	nodeD := mat.NewDense(3, 3, []float64{
		0.75, 0.20, 0.05, // ema_low  → mostly push
		0.30, 0.55, 0.15, // ema_mid  → shift to throttle
		0.05, 0.50, 0.45, // ema_high → throttle or idle
	})

	nodeA := mat.NewDense(3, 3, flat3x3(v.nodeA))

	// ── Evolver agent ─────────────────────────────────────────────────────
	//
	// World states: good_strategy=0, bad_strategy=1
	// Experience:   high_score=0, low_score=1
	// Actions:      promote=0, mutate=1

	// P_evolver (embedded in P_joint via evolverPCoupled below):
	//   good_strategy → high_score 0.85, low_score 0.15
	//   bad_strategy  → high_score 0.20, low_score 0.80

	// Mirrors workerpool selectNextStrategy: 75% mutate on poor score.
	evolverD := mat.NewDense(2, 2, []float64{
		0.80, 0.20, // high_score → elitism: keep what's working
		0.25, 0.75, // low_score  → search: try a random variation
	})

	evolverA := mat.NewDense(2, 2, []float64{
		0.85, 0.15, // promote → likely stays good
		0.60, 0.40, // mutate  → improvement less certain
	})

	// ── Joint state names ─────────────────────────────────────────────────

	jointWNames := []string{"H·G", "H·B", "D·G", "D·B", "O·G", "O·B"}
	jointXNames := []string{"low·hi", "low·lo", "mid·hi", "mid·lo", "hi·hi", "hi·lo"}
	jointGNames := []string{
		"push·promote", "push·mutate",
		"thr·promote", "thr·mutate",
		"idle·promote", "idle·mutate",
	}

	// ── P_joint (6×6) ── COUPLING 1: evolver score depends on node health ─
	//
	// A degraded or overloaded node lowers the pool score the evolver observes.
	evolverPCoupled := [3][2][2]float64{
		{{0.85, 0.15}, {0.20, 0.80}}, // node=healthy:    base rates
		{{0.55, 0.45}, {0.00, 1.00}}, // node=degraded:   -0.30 on high_score
		{{0.20, 0.80}, {0.00, 1.00}}, // node=overloaded: -0.65 on high_score
	}

	pJoint := mat.NewDense(6, 6, nil)
	for nw := 0; nw < 3; nw++ {
		for ew := 0; ew < 2; ew++ {
			for nx := 0; nx < 3; nx++ {
				for ex := 0; ex < 2; ex++ {
					pJoint.Set(nw*2+ew, nx*2+ex,
						nodeP.At(nw, nx)*evolverPCoupled[nw][ew][ex])
				}
			}
		}
	}

	// ── D_joint (6×6) ── Kronecker product: INDEPENDENT DECISIONS ─────────

	dJoint := mat.NewDense(6, 6, nil)
	for nx := 0; nx < 3; nx++ {
		for ex := 0; ex < 2; ex++ {
			for ng := 0; ng < 3; ng++ {
				for eg := 0; eg < 2; eg++ {
					dJoint.Set(nx*2+ex, ng*2+eg,
						nodeD.At(nx, ng)*evolverD.At(ex, eg))
				}
			}
		}
	}

	// ── A_joint (6×6) ── COUPLING 2: mutate boosts node recovery ──────────
	//
	// promote (eg=0): node A unchanged.
	// mutate  (eg=1): healthy probability in node A boosted by mutationBoost[ng].

	aJoint := mat.NewDense(6, 6, nil)
	for ng := 0; ng < 3; ng++ {
		for eg := 0; eg < 2; eg++ {
			// Build the node-world row for this (ng, eg) action pair.
			// For mutate (eg=1): boost the healthy entry, then renormalize
			// so the row sums to 1 — avoids negative entries regardless of boost size.
			nodeRow := [3]float64{nodeA.At(ng, 0), nodeA.At(ng, 1), nodeA.At(ng, 2)}
			if eg == 1 {
				nodeRow[0] += v.mutationBoost[ng]
				total := nodeRow[0] + nodeRow[1] + nodeRow[2]
				nodeRow[0] /= total
				nodeRow[1] /= total
				nodeRow[2] /= total
			}
			for nw := 0; nw < 3; nw++ {
				for ew := 0; ew < 2; ew++ {
					aJoint.Set(ng*2+eg, nw*2+ew, nodeRow[nw]*evolverA.At(eg, ew))
				}
			}
		}
	}

	// ── Compose J = P_joint · D_joint · A_joint ───────────────────────────

	jointAgent := &catrace.Agent{
		P:      pJoint,
		D:      dJoint,
		A:      aJoint,
		WNames: jointWNames,
		XNames: jointXNames,
		GNames: jointGNames,
	}
	J, err := jointAgent.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}

	// ── Analysis ──────────────────────────────────────────────────────────

	piJ, err := J.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	hJ, err := J.EntropyRate(2)
	if err != nil {
		log.Fatal(err)
	}

	// Trace onto {H·G=0, O·B=5}: peak and worst states.
	extremeSubset := []int{0, 5}
	traceExtreme, err := J.Trace(extremeSubset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	okExtreme, err := traceExtreme.IsTraceOf(J, extremeSubset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piExtreme, err := traceExtreme.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}

	// MFPT from O·B (worst=5) back to H·G (best=0)
	mfptWorstBest, err := J.MeanFirstPassage(5, 0)
	if err != nil {
		log.Fatal(err)
	}
	mfptSelf, err := J.MeanFirstPassage(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Monte Carlo empirical trace
	rng := rand.New(rand.NewSource(42))
	trajectory := make([]int, 500)
	trajectory[0] = 0
	for i := 1; i < len(trajectory); i++ {
		next, err := J.Sample(trajectory[i-1], rng)
		if err != nil {
			log.Fatal(err)
		}
		trajectory[i] = next
	}
	subsetSet := map[int]bool{0: true, 5: true}
	traceSeq := catrace.SampleTraceFromSequence(trajectory, subsetSet)
	compact := map[int]int{0: 0, 5: 1}
	remapped := make([]int, len(traceSeq))
	for i, val := range traceSeq {
		remapped[i] = compact[val]
	}
	est, err := catrace.EstimateKernelFromSequence(remapped, 2, 1e-3)
	if err != nil {
		log.Fatal(err)
	}

	// ── Output ─────────────────────────────────────────────────────────────

	fmt.Println("\nStationary distribution:")
	for i, p := range piJ {
		fmt.Printf("  %-6s %.4f\n", jointWNames[i], p)
	}

	fmt.Printf("\nMFPT O·B → H·G (worst→best) = %.2f steps\n", mfptWorstBest)
	fmt.Printf("MFPT H·G → H·G (self)       = %.2f steps\n", mfptSelf)
	fmt.Printf("Entropy rate                 = %.4f bits/step\n", hJ)

	fmt.Printf("\nTrace onto {H·G, O·B}:\n")
	fmt.Printf("%v\n", mat.Formatted(traceExtreme.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf = %v\n", okExtreme)
	fmt.Printf("stationary(trace): H·G=%.4f  O·B=%.4f\n", piExtreme[0], piExtreme[1])

	if est.Kernel != nil {
		fmt.Printf("\nEmpirical trace (500-step sample):\n")
		fmt.Printf("%v\n", mat.Formatted(est.Kernel.P, mat.Prefix("  ")))
	}
}

func flat3x3(a [3][3]float64) []float64 {
	return []float64{
		a[0][0], a[0][1], a[0][2],
		a[1][0], a[1][1], a[1][2],
		a[2][0], a[2][1], a[2][2],
	}
}


