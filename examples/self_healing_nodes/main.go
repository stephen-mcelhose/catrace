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
// Coupling:
//
//	P_joint (W→X): evolver's score perception depends on node world state.
//	               A degraded or overloaded node lowers the pool score the
//	               evolver observes, independent of actual strategy quality.
//
//	D_joint (X→G): Kronecker product D_node ⊗ D_evolver — independent
//	               decisions. Each agent decides from its own reading only:
//	               the node from its EMA, the evolver from the pool score.
//	               No communication within a cycle.
//
//	A_joint (G→W): mostly Kronecker product, but when the evolver mutates,
//	               node recovery probabilities are boosted — mutation sometimes
//	               discovers a better MaxWorkers or Kp that helps the node heal.

import (
	"fmt"
	"log"
	"math/rand"

	catrace "github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// ── Node agent (individual P₁, D₁, A₁) ───────────────────────────────
	//
	// World states: healthy=0, degraded=1, overloaded=2
	// Experience:   ema_low=0 (<10% errors), ema_mid=1 (10-40%), ema_high=2 (>40%)
	// Actions:      push=0, throttle=1, idle=2

	// P_node (3×3, W→X): EMA band given actual node health.
	// EMA alpha=0.20 (workerpool): ~5 steps to forget the past.
	// healthy:    low true error rate → EMA mostly low
	// degraded:   moderate errors     → EMA spread across low/mid
	// overloaded: high error rate     → EMA mostly high
	nodeP := mat.NewDense(3, 3, []float64{
		// ema_low  ema_mid  ema_high
		0.80, 0.15, 0.05, // healthy
		0.20, 0.55, 0.25, // degraded
		0.05, 0.20, 0.75, // overloaded
	})

	// D_node (3×3, X→G): action given EMA band.
	// PID logic: higher EMA → more throttle / idle.
	// Mirrors workerpool: computePIDThrottle returns sleep proportional to EMA.
	nodeD := mat.NewDense(3, 3, []float64{
		// push  throttle  idle
		0.75, 0.20, 0.05, // ema_low  — mostly keep pushing
		0.30, 0.55, 0.15, // ema_mid  — shift to throttle
		0.05, 0.50, 0.45, // ema_high — throttle or idle to recover
	})

	// A_node (3×3, G→W): world transition given action.
	// push:     works hard, sustains throughput, small risk of overload
	// throttle: self-healing — reduces load, best recovery from degraded
	// idle:     maximal recovery, no throughput
	nodeA := mat.NewDense(3, 3, []float64{
		// healthy  degraded  overloaded
		0.70, 0.25, 0.05, // push
		0.65, 0.30, 0.05, // throttle
		0.85, 0.13, 0.02, // idle
	})

	nodeAgent := &catrace.Agent{
		P:      nodeP,
		D:      nodeD,
		A:      nodeA,
		WNames: []string{"healthy", "degraded", "overloaded"},
		XNames: []string{"ema_low", "ema_mid", "ema_high"},
		GNames: []string{"push", "throttle", "idle"},
	}
	nodeW, err := nodeAgent.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}

	// ── Evolver agent (individual P₂, D₂, A₂) ────────────────────────────
	//
	// World states: good_strategy=0, bad_strategy=1
	// Experience:   high_score=0, low_score=1
	// Actions:      promote=0, mutate=1

	// P_evolver (2×2, W→X): score perception given actual strategy quality.
	// Score = mbps × (successes/total)² — workerpool's Evaluate().
	evolverP := mat.NewDense(2, 2, []float64{
		// high_score  low_score
		0.85, 0.15, // good_strategy — usually sees a high score
		0.20, 0.80, // bad_strategy  — usually sees a low score
	})

	// D_evolver (2×2, X→G): promote vs mutate given score.
	// Mirrors workerpool selectNextStrategy: 75% mutate when score is poor.
	evolverD := mat.NewDense(2, 2, []float64{
		// promote  mutate
		0.80, 0.20, // high_score — elitism: keep what's working
		0.25, 0.75, // low_score  — search: try a random variation
	})

	// A_evolver (2×2, G→W): strategy quality given action.
	// promote: preserves best known config → likely stays good
	// mutate:  random walk → sometimes better, sometimes worse
	evolverA := mat.NewDense(2, 2, []float64{
		// good  bad
		0.85, 0.15, // promote
		0.60, 0.40, // mutate
	})

	evolverAgent := &catrace.Agent{
		P:      evolverP,
		D:      evolverD,
		A:      evolverA,
		WNames: []string{"good_strategy", "bad_strategy"},
		XNames: []string{"high_score", "low_score"},
		GNames: []string{"promote", "mutate"},
	}
	evolverW, err := evolverAgent.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}

	// ── Joint state names ─────────────────────────────────────────────────

	jointWNames := []string{"H·G", "H·B", "D·G", "D·B", "O·G", "O·B"}
	jointXNames := []string{
		"low·hi", "low·lo",
		"mid·hi", "mid·lo",
		"hi·hi", "hi·lo",
	}
	jointGNames := []string{
		"push·promote", "push·mutate",
		"thr·promote", "thr·mutate",
		"idle·promote", "idle·mutate",
	}

	// ── P_joint (6×6, W→X) ── COUPLING 1: evolver score depends on node health
	//
	// P_joint[(nw,ew), (nx,ex)] = P_node[nw,nx] × P_evolver_coupled[(nw,ew),ex]
	//
	// When the node is degraded or overloaded, the pool score the evolver
	// observes drops regardless of actual strategy quality:
	//   healthy:    evolver P unchanged (base rates)
	//   degraded:   high_score probability reduced by 0.30
	//   overloaded: high_score probability reduced by 0.65

	// evolverPCoupled[node_world][evolver_world] = [high_score, low_score]
	evolverPCoupled := [3][2][2]float64{
		// node=healthy: base rates
		{{0.85, 0.15}, {0.20, 0.80}},
		// node=degraded: good_strategy 0.85→0.55, bad_strategy 0.20→0.00
		{{0.55, 0.45}, {0.00, 1.00}},
		// node=overloaded: good_strategy 0.85→0.20, bad_strategy always low
		{{0.20, 0.80}, {0.00, 1.00}},
	}

	pJoint := mat.NewDense(6, 6, nil)
	for nw := 0; nw < 3; nw++ {
		for ew := 0; ew < 2; ew++ {
			row := nw*2 + ew
			for nx := 0; nx < 3; nx++ {
				for ex := 0; ex < 2; ex++ {
					col := nx*2 + ex
					pJoint.Set(row, col, nodeP.At(nw, nx)*evolverPCoupled[nw][ew][ex])
				}
			}
		}
	}

	// ── D_joint (6×6, X→G) ── Kronecker product: INDEPENDENT DECISIONS
	//
	// Node reads its own EMA. Evolver reads the pool score.
	// Neither knows the other's reading at decision time.
	dJoint := mat.NewDense(6, 6, nil)
	for nx := 0; nx < 3; nx++ {
		for ex := 0; ex < 2; ex++ {
			row := nx*2 + ex
			for ng := 0; ng < 3; ng++ {
				for eg := 0; eg < 2; eg++ {
					col := ng*2 + eg
					dJoint.Set(row, col, nodeD.At(nx, ng)*evolverD.At(ex, eg))
				}
			}
		}
	}

	// ── A_joint (6×6, G→W) ── COUPLING 2: mutate boosts node recovery
	//
	// promote (eg=0): node's own recovery probabilities unchanged.
	//                 The existing best config is preserved as-is.
	//
	// mutate  (eg=1): node recovery probabilities are boosted.
	//                 Mutation sometimes discovers a better MaxWorkers or Kp
	//                 setting that helps the node heal faster.
	nodeABoosted := [3][3]float64{
		{0.75, 0.20, 0.05}, // push     + mutation find
		{0.75, 0.22, 0.03}, // throttle + mutation find
		{0.90, 0.09, 0.01}, // idle     + mutation find
	}

	aJoint := mat.NewDense(6, 6, nil)
	for ng := 0; ng < 3; ng++ {
		for eg := 0; eg < 2; eg++ {
			row := ng*2 + eg
			for nw := 0; nw < 3; nw++ {
				for ew := 0; ew < 2; ew++ {
					col := nw*2 + ew
					var na float64
					if eg == 1 { // mutate
						na = nodeABoosted[ng][nw]
					} else {
						na = nodeA.At(ng, nw)
					}
					aJoint.Set(row, col, na*evolverA.At(eg, ew))
				}
			}
		}
	}

	// ── Compose J = P_joint · D_joint · A_joint ───────────────────────────
	//
	// Dimensions: P(6×6) · D(6×6) · A(6×6) = 6×6 joint world kernel.
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

	// ── Analysis of joint kernel ───────────────────────────────────────────

	piJ, err := J.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	hJ, err := J.EntropyRate(2)
	if err != nil {
		log.Fatal(err)
	}
	classesJ, err := J.Classes(1e-12)
	if err != nil {
		log.Fatal(err)
	}

	// ── Trace onto {H·G=0, O·B=5}: best and worst joint states ────────────
	//
	// An outside observer who can only tell "everything fine" from "everything
	// failing" — what are the effective transition probabilities between those
	// two states, with all intermediate states summed out?
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
	restrictedExtreme, err := catrace.RestrictDistribution(piJ, extremeSubset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}

	// ── Monte Carlo: sample J, extract empirical trace ─────────────────────

	rng := rand.New(rand.NewSource(42))
	trajectory := make([]int, 500)
	trajectory[0] = 0 // start at H·G
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
	for i, v := range traceSeq {
		remapped[i] = compact[v]
	}
	est, err := catrace.EstimateKernelFromSequence(remapped, 2, 1e-3)
	if err != nil {
		log.Fatal(err)
	}

	// ── Output ─────────────────────────────────────────────────────────────

	fmt.Println("=== Node world kernel (Q_node = P·D·A) ===")
	fmt.Printf("%v\n\n", mat.Formatted(nodeW.P, mat.Prefix("  ")))

	fmt.Println("=== Evolver world kernel (Q_evolver = P·D·A) ===")
	fmt.Printf("%v\n\n", mat.Formatted(evolverW.P, mat.Prefix("  ")))

	fmt.Println("=== Joint world kernel J = P_joint · D_joint · A_joint ===")
	fmt.Printf("%v\n\n", mat.Formatted(J.P, mat.Prefix("  ")))

	fmt.Printf("stationary(J):\n")
	for i, p := range piJ {
		fmt.Printf("  %-8s %.4f\n", jointWNames[i], p)
	}
	fmt.Printf("entropy_rate(J)      = %.6f bits/step\n", hJ)
	fmt.Printf("recurrent classes(J) = %v\n\n", classesJ.Recurrent)

	fmt.Println("=== Trace onto {H·G, O·B} — peak and worst states ===")
	fmt.Printf("%v\n\n", mat.Formatted(traceExtreme.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf                      = %v\n", okExtreme)
	fmt.Printf("stationary(J)|{H·G,O·B} normed = H·G:%.4f  O·B:%.4f\n", restrictedExtreme[0], restrictedExtreme[1])
	fmt.Printf("stationary(trace)              = H·G:%.4f  O·B:%.4f\n\n", piExtreme[0], piExtreme[1])

	if est.Kernel != nil {
		fmt.Println("=== Empirical trace kernel (500-step sample) ===")
		fmt.Printf("%v\n", mat.Formatted(est.Kernel.P, mat.Prefix("  ")))
	}
}
