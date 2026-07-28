package main

// Scenario: Single LLM Task Agent
//
// An LLM support agent is assigned to handle a task independently.
//
// The agent operates in a partially observable environment:
//
//   World states (W) = what is actually true in the task environment
//     - "task_routine"   : the ticket is straightforward
//     - "task_complex"    : the ticket is ambiguous or difficult
//
//   Experience states (X) = what the LLM agent believes or perceives from the prompt/context
//     - "looks_routine"  : the agent reads the task as simple
//     - "looks_risky"    : the agent reads the task as unclear or high-stakes
//
//   Action states (G) = what the agent does next
//     - "answer"         : respond directly to the user
//     - "clarify"        : ask a clarifying question
//     - "escalate"       : escalate to human review
//
// Perception kernel P: W => X
//   The real task condition does not perfectly determine what the LLM thinks.
//   An actually complex task may still look routine if the prompt is incomplete.
//   A routine task may look risky if the wording is strange.
//
// Decision kernel D: X => G
//   Given its internal reading of the task, the LLM chooses an action
//   probabilistically. If it thinks the task looks routine, it usually
//   answers directly. If it thinks the task looks risky, it may escalate.
//
// Action effect kernel A: G => W
//   The chosen action changes the real task situation. Answering directly
//   may resolve a routine task. Asking a clarifying question may turn an
//   ambiguous task into a routine one. Escalating may move a risky task
//   into a safer reviewed state.
//
// The composite kernel Q = DAP : X => X describes the full closed-loop
// experience dynamics: "given how the agent currently interprets the task,
// what is the distribution over how it will interpret the task next, after
// acting and re-observing?"

import (
	"fmt"
	"log"

	"catrace"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// D: X => G  (decision: experience -> action distribution)
	//   looks_routine -> mostly answer, sometimes clarify
	//   looks_risky   -> sometimes answer, sometimes clarify, often escalate
	D := mat.NewDense(2, 3, []float64{
		0.8, 0.2, 0.0,
		0.1, 0.3, 0.6,
	})

	// A: G => W  (action effect: action -> world state distribution)
	//   answer  -> mostly routine, sometimes complex
	//   clarify -> tends to move toward routine
	//   escalate -> tends to move toward routine (reviewed)
	A := mat.NewDense(3, 2, []float64{
		0.9, 0.1,
		0.4, 0.6,
		0.2, 0.8,
	})

	// P: W => X  (perception: world state -> experience distribution)
	//   task_routine  -> mostly looks_routine, sometimes looks_risky
	//   task_complex  -> sometimes looks_routine, mostly looks_risky
	P := mat.NewDense(2, 2, []float64{
		0.85, 0.15,
		0.25, 0.75,
	})

	agent := &catrace.Agent{
		D:      D,
		A:      A,
		P:      P,
		XNames: []string{"looks_routine", "looks_risky"},
		GNames: []string{"answer", "clarify", "escalate"},
		WNames: []string{"task_routine", "task_complex"},
	}

	Q, err := agent.QualiaKernel()
	if err != nil {
		log.Fatal(err)
	}
	S, err := agent.StrategyKernel()
	if err != nil {
		log.Fatal(err)
	}
	W, err := agent.WorldKernel()
	if err != nil {
		log.Fatal(err)
	}
	pi, err := Q.Stationary(1e-12, 5000)
	if err != nil {
		log.Fatal(err)
	}
	h, err := Q.EntropyRate(2)
	if err != nil {
		log.Fatal(err)
	}
	classes, err := Q.Classes(1e-12)
	if err != nil {
		log.Fatal(err)
	}
	next, err := Q.LeftAction([]float64{1, 0})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Q = D*A*P")
	fmt.Printf("%v\n\n", mat.Formatted(Q.P, mat.Prefix("  ")))

	fmt.Println("S = A*P*D")
	fmt.Printf("%v\n\n", mat.Formatted(S.P, mat.Prefix("  ")))

	fmt.Println("W = P*D*A")
	fmt.Printf("%v\n\n", mat.Formatted(W.P, mat.Prefix("  ")))

	fmt.Printf("stationary(Q) = %.6f %.6f\n", pi[0], pi[1])
	fmt.Printf("entropy_rate(Q) = %.6f bits/step\n", h)
	fmt.Printf("left_action([1,0], Q) = %.6f %.6f\n", next[0], next[1])
	fmt.Printf("recurrent classes = %v\n", classes.Recurrent)
	fmt.Printf("transient states = %v\n", classes.Transient)
}
