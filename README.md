# catrace

`catrace` is a Go library for finite-state Markov models of autonomous-agent networks, implemented with `gonum`.

This project intentionally extracts only the stochastic / Markov machinery from the source paper and excludes consciousness and philosophical claims.

## Source note

This implementation is derived from the mathematical constructions in:

> Hoffman, Prakash & Chattopadhyay, *Traces of Consciousness*, Preprints 2024.
> https://www.preprints.org/manuscript/202410.1305/v1
> Published under CC BY 4.0.

In particular, the project draws on the paper's treatment of agents as coupled stochastic maps between world, experience, and action spaces, and on the trace-chain construction used to reduce a larger Markov process to an effective process on an observed subset of states.

## Scope

Implemented concepts:

- validated row-stochastic square kernels
- agent triplet model with perception / decision / action maps
- composed kernels:
  - $Q = DAP$
  - $S = APD$
  - $W = PDA$
- trace chains on state subsets
- stationary distributions
- entropy rate
- communicating classes / recurrent classes
- mean first-passage times and commute times
- sampling utilities
- kernel estimation from sequences

## Package layout

- `kernel.go` — core kernel types and composition helpers
- `agent.go` — agent triplet model
- `trace.go` — trace-chain construction
- `stationary.go` — stationary distribution and entropy rate
- `analysis.go` — communicating/recurrent classes
- `passage.go` — first-passage / commute times
- `sample.go` — sampling and estimation
- `util.go` — helpers
- `docs/math_summary.md` — sanitized math summary
- `GLOSSARY.md` — definitions of all math and notation terms
- `examples/simple_agent` — single-agent composition demo
- `examples/trace_analysis` — trace-chain demo

## Core model

An agent is described by three row-stochastic kernels:

- perception $P: W \to X$
- decision $D: X \to G$
- action $A: G \to W$

The main derived kernels are:

- qualia kernel $Q = DAP$ on experiences
- strategy kernel $S = APD$ on actions
- world kernel $W = PDA$ on world states

## Trace chain

For a kernel $L$ and observed subset $S$, the trace kernel is

$\text{Tr}(L) = L_{SS} + L_{SB}(I - L_{BB})^{-1}L_{BS}$

under the block decomposition into observed states $S$ and hidden states $B$.

This is implemented by `(*Kernel).Trace`.

## Example usage

```go
Q, err := agent.QualiaKernel()
pi, err := Q.Stationary(1e-12, 5000)
H, err := Q.EntropyRate(2)
tr, err := parent.Trace([]int{0,1}, 1e-12)
```

## API

The library provides:

- `Kernel` with optional `StateNames`
- `Sample`
- `Trace` and `IsTraceOf`
- `Stationary`
- `EntropyRate`
- `Classes`
- `MeanFirstPassage`
- `CommuteTime`
- `Agent` abstraction with `D`, `A`, `P` and derived kernels `Q`, `S`, `W`

## Scenario write-ups

These examples are presented as short stories rather than abstract state tables, to make the state spaces easier to reason about.

### 1. Single LLM task agent

Story:

An LLM support agent is assigned to handle a task independently. The real task may be routine or genuinely complex, but the agent only sees the prompt and surrounding context, so it can misread the situation. Based on its internal interpretation, it may answer directly, ask a clarifying question, or escalate to a human.

[Full story, state meanings, and interpretation →](docs/patterns/story-single-llm-agent.md)

### 2. LLM agent with hidden support system

Story:

A focal LLM agent is visible to us, but the rest of the support system is hidden in the background: retrieval services, monitoring tools, human reviewers, and other agents. We only observe whether the focal agent appears valid or invalid from the outside. The hidden system may help or hinder it before we see the focal agent again.

[Full story, state meanings, and interpretation →](docs/patterns/story-hidden-support-system.md)

### 3. Two-agent validator / repair pair

Story:

A worker agent performs tasks while a validator agent monitors its health. Either agent may itself be functioning well or badly. When the validator is healthy, it can detect worker problems and attempt repairs — but repair takes effort and can degrade the validator too. When both are degraded, recovery depends on chance.

[Full story, state meanings, and interpretation →](docs/patterns/story-validator-repair.md)

### 4. Self-adjusting / self-healing network nodes

Story:

A network node monitors its own error rate and throttles itself when errors climb. An outer evolutionary loop watches pool throughput and mutates the node's configuration when performance drops. The two loops compete to explain recovery: in one regime the node's own throttle is the primary healer; in the other the evolver's config search is what keeps the system alive.

[Full story, state meanings, and interpretation →](docs/patterns/story-self-healing-nodes.md)

### 5. Three-agent majority-valid coordination network *(not yet implemented)*

Story:

Three agents coordinate on a shared task. As long as at least two are functioning well, the team can usually stabilize itself and recover local failures. Once only one agent remains reliable, recovery becomes much harder and collapse becomes more likely.

[Full story, state meanings, and interpretation →](docs/patterns/story-supervisor.md)

### 6. Four-agent pipeline with escalation *(not yet implemented)*

Story:

A four-stage workflow handles incoming work. Local fixes are fast but unreliable, rerouting helps around isolated failures, and escalation is expensive but can stabilize the whole system. The network may move between fully healthy operation, partial degradation, and systemic breakdown.

[Full story, state meanings, and interpretation →](docs/patterns/story-prompt-chaining.md)

### 7. Routing agent

Story:

A customer support system receives tickets of genuinely different types — billing disputes, technical faults, and general enquiries — but the router agent can only read the ticket text, not the underlying truth. Misclassification sends a ticket to the wrong specialist, who cannot resolve it and re-queues it, often as an apparently different type. The tension is between classification accuracy and the cost of a wrong route: a confident wrong decision causes a longer misrouting loop than an uncertain escalation to human triage.

- world states: `billing_ticket`, `technical_ticket`, `general_ticket` — the true nature of the incoming request
- experience states: `reads_billing`, `reads_technical`, `reads_general` — the router's perceived classification
- actions: `route_billing`, `route_technical`, `route_general`, `escalate_human` — where the router sends the ticket

Interpretation:

- perception captures classification accuracy — diagonal entries are correct reads; off-diagonal entries encode misclassification probability
- decision captures the routing policy given perceived type, including the rate of escalation to human triage regardless of perceived classification
- action captures how the routing choice changes the world state: a correct route resolves the ticket and draws the next from the same pool; a wrong route re-enters the ticket as an apparently different type
- the stationary distribution shows the steady-state mix of ticket types in the queue, including misrouted tickets that have cycled through the wrong specialist
- MFPT between ticket types measures misrouting loop latency

Code: `examples/routing/main.go` — walkthrough at `examples/routing/WALKTHROUGH.md`

#### Played-out version: Story 7

**Version A: statistically typical path**

1. World state `technical_ticket` (37.2% of steady-state — the most common queue type)
2. P[`technical_ticket`, `reads_technical`] = 0.80 → router correctly perceives `reads_technical`
3. D[`reads_technical`, `route_technical`] = 0.80 → router sends to technical specialist
4. A[`route_technical`, `technical_ticket`] = 0.70 → ticket resolved; next is technical (0.80 × 0.80 × 0.70 = 0.448)

This path contributes the bulk of W[`technical_ticket`, `technical_ticket`] = 0.533. In plain English: a technical ticket is correctly identified, resolved, and another technical ticket arrives.

**Version B: misrouting loop**

1. World state `billing_ticket` — genuine billing dispute arrives
2. P[`billing_ticket`, `reads_technical`] = 0.15 → router misreads it as technical
3. D[`reads_technical`, `route_technical`] = 0.80 → router sends to technical specialist
4. A[`route_technical`, `billing_ticket`] = 0.15 → specialist cannot resolve it; re-enters as billing (0.15 × 0.80 × 0.15 = 0.018)

This path contributes to W[`billing_ticket`, `billing_ticket`] via the wrong-route direction. In plain English: a billing ticket is misclassified as technical, the technical specialist cannot help, and the ticket wastes one full cycle before returning to the billing queue.

**Version C: escalation recovery**

1. World state `general_ticket` — general enquiry arrives
2. P[`general_ticket`, `reads_billing`] = 0.10 → router misreads it as billing
3. D[`reads_billing`, `escalate_human`] = 0.10 → router escalates rather than routes
4. A[`escalate_human`, `general_ticket`] = 0.33 → human resolves it correctly (0.10 × 0.10 × 0.33 = 0.0033)

This path contributes to W[`general_ticket`, `general_ticket`] via the escalation route. In plain English: a misread ticket is caught by escalation — the human resolves it regardless of type, at the cost of one escalation cycle rather than a misrouting loop.

Across all three paths: technical tickets dominate the queue (37.2%) because they are classified most accurately. The MFPT billing→technical (3.55 steps) is a direct measure of how long the misrouting loop runs — reduce off-diagonal entries in P and this number grows, confirming better routing.

## Tests and examples

The intended style for this project is:

- short mathematical write-up paired with each example
- small finite-state scenario with named states
- runnable example program

Current files:

- `examples/simple_agent/main.go`
- `examples/trace_analysis/main.go`
- `examples/validator_repair/main.go`
- `examples/self_healing_nodes/main.go`
- `catrace_test.go`

## Build

Requires Go 1.22+.

```
go build ./...
go test ./...
go run examples/simple_agent/main.go
go run examples/trace_analysis/main.go
```
