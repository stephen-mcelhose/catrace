package main

// Scenario: Blackboard / shared-workspace collaboration
//
// Three identical specialists (radiologist, pathologist, notes — labels only)
// read and write a shared diagnostic board. No agent-to-agent messaging:
// coordination is only through board status.
//
//	J[w, w'] = Σ_{x,g} P_joint[w,x] · D_joint[x,g] · A[w,g,w']
//
//	P_joint (W→X³): each specialist's perception depends on the shared board.
//	               A tentative hypothesis raises P(evidence_strong) — the
//	               board contents trigger opportunistic contribution.
//	D_joint (X³→G³): D⊗D⊗D — independent decisions; no within-cycle talk.
//	A[w,g,·]:        joint posts advance the board by counting
//	               post_finding / endorse_prior / flag_contradiction.
//	               (Next board depends on current board + joint action, so
//	               A is not a pure G→W rectangle; J is assembled row-wise.)
//
// Board world W (4):
//	undiagnosed, tentative_diagnosis, confirmed_diagnosis, contradicted
//
// confirmed_diagnosis is absorbing. contradicted may reopen to undiagnosed
// so MFPT undiagnosed → confirmed stays finite.

import (
	"fmt"
	"log"
	"os"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

const (
	undiagnosed = iota
	tentative
	confirmed
	contradicted
	nBoard

	nExp    = 2 // evidence_strong, evidence_weak
	nAct    = 4 // post, endorse, flag, request
	nJointX = nExp * nExp * nExp // 8
	nJointG = nAct * nAct * nAct // 64
)

const (
	post = iota
	endorse
	flag
	request
)

func main() {
	boardNames := []string{
		"undiagnosed",
		"tentative_diagnosis",
		"confirmed_diagnosis",
		"contradicted",
	}

	// Per-specialist P: board → experience (identical specialists).
	// Rows: undiagnosed, tentative, confirmed, contradicted.
	// Cols: evidence_strong, evidence_weak.
	specP := mat.NewDense(nBoard, nExp, []float64{
		0.30, 0.70, // undiagnosed: board empty — precondition rarely fires
		0.75, 0.25, // tentative: board contents trigger contribution
		0.85, 0.15, // confirmed: still looks actionable (absorbing)
		0.20, 0.80, // contradicted: little left to add
	})

	// Per-specialist D: experience → action.
	specD := mat.NewDense(nExp, nAct, []float64{
		// post  endorse  flag  request
		0.30, 0.45, 0.10, 0.15, // evidence_strong
		0.08, 0.07, 0.40, 0.45, // evidence_weak
	})

	// P_joint (4×8): product of three identical board→experience maps.
	pJoint := mat.NewDense(nBoard, nJointX, nil)
	for w := 0; w < nBoard; w++ {
		for x := 0; x < nJointX; x++ {
			e0, e1, e2 := decode3(x, nExp)
			pJoint.Set(w, x, specP.At(w, e0)*specP.At(w, e1)*specP.At(w, e2))
		}
	}

	// D_joint (8×64) = D⊗D⊗D.
	dJoint := mat.NewDense(nJointX, nJointG, nil)
	for x := 0; x < nJointX; x++ {
		e0, e1, e2 := decode3(x, nExp)
		for g := 0; g < nJointG; g++ {
			a0, a1, a2 := decode3(g, nAct)
			dJoint.Set(x, g, specD.At(e0, a0)*specD.At(e1, a1)*specD.At(e2, a2))
		}
	}

	// Assemble J row-wise: next board depends on (current w, joint g).
	jData := mat.NewDense(nBoard, nBoard, nil)
	for w := 0; w < nBoard; w++ {
		for x := 0; x < nJointX; x++ {
			px := pJoint.At(w, x)
			if px == 0 {
				continue
			}
			for g := 0; g < nJointG; g++ {
				pdg := px * dJoint.At(x, g)
				if pdg == 0 {
					continue
				}
				nPost, nEndorse, nFlag := countActions(g)
				dist := boardNext(w, nPost, nEndorse, nFlag)
				for wp := 0; wp < nBoard; wp++ {
					jData.Set(w, wp, jData.At(w, wp)+pdg*dist[wp])
				}
			}
		}
	}
	roundMatrix(jData, 1e6)

	J, err := catrace.NewKernel(jData, boardNames)
	if err != nil {
		log.Fatal(err)
	}

	pi, err := J.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	h, err := J.EntropyRate(2)
	if err != nil {
		log.Fatal(err)
	}
	classes, err := J.Classes(1e-12)
	if err != nil {
		log.Fatal(err)
	}
	mfpt, err := J.MeanFirstPassage(undiagnosed, confirmed)
	if err != nil {
		log.Fatal(err)
	}

	subset := []int{undiagnosed, confirmed, contradicted}
	tr, err := J.Trace(subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	ok, err := tr.IsTraceOf(J, subset, 1e-12)
	if err != nil {
		log.Fatal(err)
	}
	piTr, err := tr.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Blackboard: per-specialist kernels (identical; roles are labels) ===")
	fmt.Println("\nP (board → evidence_strong / evidence_weak):")
	fmt.Printf("%v\n", mat.Formatted(specP, mat.Prefix("  ")))
	fmt.Println("\nD (evidence → post / endorse / flag / request):")
	fmt.Printf("%v\n", mat.Formatted(specD, mat.Prefix("  ")))

	fmt.Println("\n=== Coupling ===")
	fmt.Println("P_joint: each specialist reads shared board only; tentative raises P(evidence_strong)")
	fmt.Println("D_joint = D⊗D⊗D: independent decisions; no out-of-band A2A")
	fmt.Println("A(w,g,·): count posts/endorsements/flags → next board status")

	fmt.Println("\n=== Board world kernel J (assembled Σ P_joint · D_joint · A(w,g,·)) ===")
	fmt.Printf("%v\n\n", mat.Formatted(J.P, mat.Prefix("  ")))

	fmt.Println("Stationary distribution π:")
	for i, name := range J.StateNames {
		fmt.Printf("  %-22s %.6f\n", name, pi[i])
	}
	fmt.Println()
	fmt.Printf("Entropy rate H(J): %.6f bits/step\n", h)
	fmt.Printf("Recurrent classes: %v\n", classes.Recurrent)
	fmt.Printf("Transient states:  %v\n", classes.Transient)
	fmt.Printf("MFPT undiagnosed → confirmed_diagnosis: %.4f steps\n\n", mfpt)

	fmt.Println("=== Trace onto {undiagnosed, confirmed_diagnosis, contradicted} ===")
	fmt.Printf("%v\n\n", mat.Formatted(tr.P, mat.Prefix("  ")))
	fmt.Printf("IsTraceOf verification: %v\n", ok)
	fmt.Printf("Stationary of trace: undiagnosed=%.6f  confirmed=%.6f  contradicted=%.6f\n",
		piTr[0], piTr[1], piTr[2])

	fmt.Println("\n=== Sample A(·|board) for key joint actions ===")
	samples := []struct {
		label string
		w     int
		acts  [3]int
	}{
		{"undiagnosed + post|ask|ask", undiagnosed, [3]int{post, request, request}},
		{"tentative + end|end|ask", tentative, [3]int{endorse, endorse, request}},
		{"tentative + end|end|end", tentative, [3]int{endorse, endorse, endorse}},
		{"tentative + flag|end|ask", tentative, [3]int{flag, endorse, request}},
	}
	for _, s := range samples {
		nPost, nEndorse, nFlag := 0, 0, 0
		for _, a := range s.acts {
			switch a {
			case post:
				nPost++
			case endorse:
				nEndorse++
			case flag:
				nFlag++
			}
		}
		d := boardNext(s.w, nPost, nEndorse, nFlag)
		fmt.Printf("  %-28s → [%.2f %.2f %.2f %.2f]\n", s.label, d[0], d[1], d[2], d[3])
	}

	for _, kv := range []struct {
		k     *catrace.Kernel
		title string
		file  string
	}{
		{J, "Blackboard — board world kernel J", "blackboard_J.html"},
		{tr, "Blackboard — trace {undiagnosed, confirmed, contradicted}", "blackboard_trace.html"},
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

func countActions(g int) (nPost, nEndorse, nFlag int) {
	a0, a1, a2 := decode3(g, nAct)
	for _, a := range []int{a0, a1, a2} {
		switch a {
		case post:
			nPost++
		case endorse:
			nEndorse++
		case flag:
			nFlag++
		}
	}
	return
}

// boardNext returns the distribution over next board status given current
// board and counts of post / endorse / flag among the three specialists.
func boardNext(w, nPost, nEndorse, nFlag int) [nBoard]float64 {
	var d [nBoard]float64
	switch w {
	case undiagnosed:
		if nPost >= 1 {
			d = [nBoard]float64{0.15, 0.80, 0.05, 0.00}
		} else {
			d = [nBoard]float64{0.90, 0.08, 0.02, 0.00}
		}
	case tentative:
		if nFlag >= 1 {
			d = [nBoard]float64{0.05, 0.25, 0.10, 0.60}
		} else if nEndorse >= 2 {
			d = [nBoard]float64{0.02, 0.15, 0.80, 0.03}
		} else if nEndorse >= 1 && nPost >= 1 {
			d = [nBoard]float64{0.05, 0.40, 0.45, 0.10}
		} else if nPost >= 1 {
			d = [nBoard]float64{0.05, 0.70, 0.15, 0.10}
		} else {
			d = [nBoard]float64{0.10, 0.75, 0.10, 0.05}
		}
	case confirmed:
		d = [nBoard]float64{0, 0, 1, 0}
	case contradicted:
		// Reopen path keeps MFPT undiagnosed→confirmed finite.
		d = [nBoard]float64{0.20, 0.05, 0.00, 0.75}
	}
	return d
}

func decode3(idx, base int) (int, int, int) {
	i2 := idx % base
	idx /= base
	i1 := idx % base
	i0 := idx / base
	return i0, i1, i2
}

func roundMatrix(m *mat.Dense, scale float64) {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m.Set(i, j, float64(int(m.At(i, j)*scale+0.5))/scale)
		}
	}
}
