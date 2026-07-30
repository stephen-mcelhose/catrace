package main

// Scenario: Two-Agent Validator / Repair Pair
//
// One agent (Worker) performs tasks and a second agent (Validator) checks
// and repairs the Worker's output. Each may be in a valid or invalid mode.
//
// Rather than constructing a joint kernel by hand, we define joint P, D, A
// kernels over product state spaces and compose them as:
//
//	J = P_joint · D_joint · A_joint
//
// This is the same PDA composition as a single agent, lifted to the
// multi-agent level. Coupling enters through specific kernels with
// clear physical meaning:
//
//	P_joint (W→X): Validator's perception depends on Worker's world state.
//	               When Worker is degraded, Validator is more likely to see
//	               a problem even if its own internal state is fine.
//	D_joint (X→G): Kronecker product D₁⊗D₂ — independent decisions.
//	               Each agent decides from its own experience only.
//	               No communication between agents at the decision step.
//	A_joint (G→W): Validator's repair action restores Worker's world state.
//	               Non-repair actions affect only the acting agent's state.
//
// Joint world states (W₁×W₂), row-major (worker first):
//
//	"VV" : both Worker and Validator reliable
//	"VI" : Worker reliable, Validator degraded
//	"IV" : Worker degraded, Validator reliable
//	"II" : both degraded
//
// Joint experience states (X₁×X₂):
//
//	"ok·good"   : Worker sees ok, Validator looks good
//	"ok·bad"    : Worker sees ok, Validator looks bad
//	"prob·good" : Worker sees problem, Validator looks good
//	"prob·bad"  : Worker sees problem, Validator looks bad
//
// Joint action states (G₁×G₂), row-major (worker first):
//
//	(produce|validate), (produce|repair), (produce|idle),
//	(self_check|validate), (self_check|repair), (self_check|idle),
//	(idle|validate), (idle|repair), (idle|idle)

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// ── Worker agent (individual P₁, D₁, A₁) ─────────────────────────────

	// P₁: W₁ → X₁ (2×2)
	// Worker usually knows when it is functioning well, but misses some
	// degradation, and occasionally misreads a healthy state as problematic.
	workerP := mat.NewDense(2, 2, []float64{
		0.90, 0.10, // worker_valid:   mostly sees ok
		0.30, 0.70, // worker_invalid: mostly sees problem
	})

	// D₁: X₁ → G₁ (2×3)
	// When things look fine, produce. When a problem is detected, self-check.
	workerD := mat.NewDense(2, 3, []float64{
		0.80, 0.15, 0.05, // sees_ok:      produce / self_check / idle
		0.10, 0.60, 0.30, // sees_problem: self_check / idle
	})

	// A₁: G₁ → W₁ (3×2)
	// Producing tends to keep the worker valid. Idling drifts toward invalid.
	workerA := mat.NewDense(3, 2, []float64{
		0.70, 0.30, // produce:    valid / invalid
		0.50, 0.50, // self_check: neutral
		0.40, 0.60, // idle:       drifts toward invalid
	})

	worker := &catrace.Agent{
		P:      workerP,
		D:      workerD,
		A:      workerA,
		WNames: []string{"worker_valid", "worker_invalid"},
		XNames: []string{"sees_ok", "sees_problem"},
		GNames: []string{"produce", "self_check", "idle"},
	}

	workerW, err := worker.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}

	// ── Validator agent (individual P₂, D₂, A₂) ──────────────────────────

	// P₂: W₂ → X₂ (2×2)
	// A healthy validator is a reliable but slightly less sharp observer than
	// the worker's self-perception — it surveys a wider scope.
	validatorP := mat.NewDense(2, 2, []float64{
		0.85, 0.15, // validator_valid:   mostly looks good
		0.40, 0.60, // validator_invalid: often looks bad
	})

	// D₂: X₂ → G₂ (2×3)
	// When things look fine, validate (monitor). When bad, shift to repair.
	validatorD := mat.NewDense(2, 3, []float64{
		0.60, 0.10, 0.30, // looks_good: validate / repair / idle
		0.20, 0.70, 0.10, // looks_bad:  repair / validate / idle
	})

	// A₂: G₂ → W₂ (3×2)
	// Validating keeps the validator calibrated. Repairing is taxing — some
	// degradation. Idling is neutral.
	validatorA := mat.NewDense(3, 2, []float64{
		0.85, 0.15, // validate: stays calibrated
		0.60, 0.40, // repair:   taxing, some degradation
		0.50, 0.50, // idle:     neutral
	})

	validator := &catrace.Agent{
		P:      validatorP,
		D:      validatorD,
		A:      validatorA,
		WNames: []string{"validator_valid", "validator_invalid"},
		XNames: []string{"looks_good", "looks_bad"},
		GNames: []string{"validate", "repair", "idle"},
	}

	validatorW, err := validator.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}

	// ── Joint P, D, A over product spaces ────────────────────────────────
	//
	// Joint world:      W = W₁×W₂  (4 states: VV, VI, IV, II)
	// Joint experience: X = X₁×X₂  (4 states)
	// Joint actions:    G = G₁×G₂  (9 states)
	//
	// Row-major indexing: pair (i₁, i₂) → index i₁·n₂ + i₂

	jointWNames := []string{"VV", "VI", "IV", "II"}
	jointXNames := []string{"ok·good", "ok·bad", "prob·good", "prob·bad"}
	jointGNames := []string{
		"produce|validate", "produce|repair", "produce|idle",
		"self_check|validate", "self_check|repair", "self_check|idle",
		"idle|validate", "idle|repair", "idle|idle",
	}

	// P_joint (4×4, W→X) ─ COUPLING 1: Validator observes Worker's world state.
	//
	// Worker's perception is independent of Validator state — P₁[w₁,:] unchanged.
	// Validator's perception shifts when Worker is degraded: it can detect Worker
	// errors in addition to its own internal state.
	//
	// P₂_coupled (Validator perception, adjusted for Worker state w₁):
	//   w₁=valid,   w₂=valid:   [0.85, 0.15]  normal perception
	//   w₁=valid,   w₂=invalid: [0.40, 0.60]  degraded perception
	//   w₁=invalid, w₂=valid:   [0.60, 0.40]  Worker signal shifts toward bad
	//   w₁=invalid, w₂=invalid: [0.25, 0.75]  both signals bad
	//
	// P_joint[(w₁,w₂), (x₁,x₂)] = P₁[w₁,x₁] × P₂_coupled[(w₁,w₂), x₂]
	pJoint := mat.NewDense(4, 4, []float64{
		// ok·good  ok·bad  prob·good  prob·bad
		0.765, 0.135, 0.085, 0.015, // VV: 0.90×0.85, 0.90×0.15, 0.10×0.85, 0.10×0.15
		0.360, 0.540, 0.040, 0.060, // VI: 0.90×0.40, 0.90×0.60, 0.10×0.40, 0.10×0.60
		0.180, 0.120, 0.420, 0.280, // IV: 0.30×0.60, 0.30×0.40, 0.70×0.60, 0.70×0.40
		0.075, 0.225, 0.175, 0.525, // II: 0.30×0.25, 0.30×0.75, 0.70×0.25, 0.70×0.75
	})

	// D_joint (4×9, X→G) ─ D₁⊗D₂: INDEPENDENT DECISIONS.
	//
	// Each agent decides based on its own experience only — no communication.
	// D_joint[(x₁,x₂), (g₁,g₂)] = D₁[x₁,g₁] × D₂[x₂,g₂]
	// The Kronecker product of two row-stochastic matrices is row-stochastic.
	dJoint := mat.NewDense(4, 9, nil)
	for x1 := 0; x1 < 2; x1++ {
		for x2 := 0; x2 < 2; x2++ {
			for g1 := 0; g1 < 3; g1++ {
				for g2 := 0; g2 < 3; g2++ {
					row := x1*2 + x2
					col := g1*3 + g2
					dJoint.Set(row, col, workerD.At(x1, g1)*validatorD.At(x2, g2))
				}
			}
		}
	}

	// A_joint (9×4, G→W) ─ COUPLING 2: Validator's repair restores Worker.
	//
	// For non-repair actions (g₂ ≠ repair):
	//   A_joint[(g₁,g₂),(w₁,w₂)] = A₁[g₁,w₁] × A₂[g₂,w₂]
	//   Each agent's action affects only its own world state.
	//
	// For repair actions (g₂ = repair, indices 1, 4, 7):
	//   Validator repair boosts Worker's probability of being valid by ~0.20.
	//   A₁_boosted: produce→[0.90,0.10], self_check→[0.70,0.30], idle→[0.60,0.40]
	//   Validator's own state follows A₂[repair,:] = [0.60, 0.40] unchanged.
	workerABoosted := [3][2]float64{
		{0.90, 0.10}, // produce    + repair boost
		{0.70, 0.30}, // self_check + repair boost
		{0.60, 0.40}, // idle       + repair boost
	}
	aJoint := mat.NewDense(9, 4, nil)
	for g1 := 0; g1 < 3; g1++ {
		for g2 := 0; g2 < 3; g2++ {
			row := g1*3 + g2
			for w1 := 0; w1 < 2; w1++ {
				for w2 := 0; w2 < 2; w2++ {
					col := w1*2 + w2
					var a1w1 float64
					if g2 == 1 { // repair action
						a1w1 = workerABoosted[g1][w1]
					} else {
						a1w1 = workerA.At(g1, w1)
					}
					aJoint.Set(row, col, a1w1*validatorA.At(g2, w2))
				}
			}
		}
	}

	// Compose J = P_joint · D_joint · A_joint.
	//
	// The joint Agent has P:4×4, D:4×9, A:9×4.
	// Dimension checks: D.rows(4)==P.cols(4), D.cols(9)==A.rows(9), A.cols(4)==P.rows(4).
	// WorldKernel() computes P·D·A = (4×4)(4×9)(9×4) = 4×4.
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

	// ── Analysis of joint kernel ──────────────────────────────────────────

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

	// ── Trace onto {VV, II}: coarse healthy / failed view ─────────────────
	//
	// Subset indices: VV=0, II=3. Complement B={VI=1, IV=2}.
	// The trace kernel gives the effective dynamics between VV and II,
	// folding in all transient visits through VI and IV.
	coarseSubset := []int{0, 3}
	traceCoarse, err := J.Trace(coarseSubset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	okCoarse, err := traceCoarse.IsTraceOf(J, coarseSubset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piCoarse, err := traceCoarse.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	restrictedCoarse, err := catrace.RestrictDistribution(piJ, coarseSubset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}

	// ── Monte Carlo: sample J, extract trace sequence, estimate empirically ──

	rng := rand.New(rand.NewSource(42))
	trajectory := make([]int, 200)
	trajectory[0] = 0 // start at VV
	for i := 1; i < len(trajectory); i++ {
		next, err := J.Sample(trajectory[i-1], rng)
		if err != nil {
			log.Fatal(err)
		}
		trajectory[i] = next
	}
	subsetSet := map[int]bool{0: true, 3: true}
	traceSeq := catrace.SampleTraceFromSequence(trajectory, subsetSet)
	// Remap original indices {0,3} to compact indices {0,1} for the estimator.
	compact := map[int]int{0: 0, 3: 1}
	remapped := make([]int, len(traceSeq))
	for i, v := range traceSeq {
		remapped[i] = compact[v]
	}
	est, err := catrace.EstimateKernelFromSequence(remapped, 2, 1e-3)
	if err != nil {
		log.Fatal(err)
	}

	// ── Output ────────────────────────────────────────────────────────────

	fmt.Println("=== Worker individual world kernel (W₁ = P₁·D₁·A₁) ===")
	fmt.Printf("%v\n\n", mat.Formatted(workerW.P, mat.Prefix("  ")))

	fmt.Println("=== Validator individual world kernel (W₂ = P₂·D₂·A₂) ===")
	fmt.Printf("%v\n\n", mat.Formatted(validatorW.P, mat.Prefix("  ")))

	fmt.Println("=== Joint world kernel J = P_joint · D_joint · A_joint ===")
	fmt.Printf("%v\n\n", mat.Formatted(J.P, mat.Prefix("  ")))

	fmt.Printf("stationary(J)        = %.6f  %.6f  %.6f  %.6f\n", piJ[0], piJ[1], piJ[2], piJ[3])
	fmt.Printf("entropy_rate(J)      = %.6f bits/step\n", hJ)
	fmt.Printf("recurrent classes(J) = %v\n\n", classesJ.Recurrent)

	fmt.Println("=== Trace onto {VV, II} — coarse healthy / failed view ===")
	fmt.Printf("%v\n\n", mat.Formatted(traceCoarse.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf                   = %v\n", okCoarse)
	fmt.Printf("stationary(J)|{VV,II} norm  = %.6f  %.6f\n", restrictedCoarse[0], restrictedCoarse[1])
	fmt.Printf("stationary(trace)           = %.6f  %.6f\n\n", piCoarse[0], piCoarse[1])

	if est.Kernel != nil {
		fmt.Println("=== Empirical trace kernel (200-step sample, filtered to {VV,II}) ===")
		fmt.Printf("%v\n", mat.Formatted(est.Kernel.P, mat.Prefix("  ")))
	}
}
