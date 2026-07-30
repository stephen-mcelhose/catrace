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

State meanings:
- world states = what is actually true about the task
- experience states = how the LLM interprets the task
- action states = what the LLM does next

Interpretation:
- perception captures imperfect prompt interpretation
- decision captures the LLM policy
- action captures how the chosen response changes the real task situation
- the derived kernel $Q = DAP$ tells you how the agent's interpretation evolves from one interaction cycle to the next

Code: `examples/simple_agent/main.go` — walkthrough at `examples/simple_agent/WALKTHROUGH.md`

#### Played-out version: Story 1

The composite kernel $Q = DAP$ is easiest to understand when you walk through concrete paths. Each entry of $Q$ compresses many possible world-experience-action-world-experience micro-stories into one effective next-experience probability.

**Version A: statistically typical path**

1. The real task is `task_routine` — the ticket is genuinely straightforward.
2. Perception: the agent reads the prompt and, with probability 0.85, experiences `looks_routine`.
3. Decision: given `looks_routine`, the agent chooses `answer` with probability 0.8.
4. Action effect: answering directly resolves the issue, so the world stays `task_routine` with probability 0.9.
5. Re-perception: the world is still routine, so the next experience is again `looks_routine` with probability 0.85.

This path contributes the bulk of the probability mass to the transition `looks_routine -> looks_routine` in $Q$.

In plain English: the task really was simple, it looked simple, the agent answered directly, and the situation remained simple when seen again.

**Version B: alternate path — misread complex task**

1. The real task is `task_complex` — the ticket is actually ambiguous or difficult.
2. Perception: but the prompt is incomplete, so with probability 0.25 the agent still experiences `looks_routine`.
3. Decision: given `looks_routine`, the agent chooses `answer` with probability 0.8.
4. Action effect: answering directly rarely resolves a complex task, so the world remains `task_complex` with probability 0.1.
5. Re-perception: the unresolved situation now looks problematic, so the next experience becomes `looks_risky` with probability 0.75.

This path contributes to the transition `looks_routine -> looks_risky` in $Q$.

In plain English: the task was harder than it looked, the agent answered too quickly, the problem persisted, and on the next pass the agent finally saw the difficulty.

**Version C: recovery path — productive caution**

1. The real task is `task_complex`.
2. Perception: this time the agent reads it correctly and experiences `looks_risky` with probability 0.75.
3. Decision: given `looks_risky`, the agent chooses `clarify` with probability 0.3 (or `escalate` with probability 0.6).
4. Action effect: asking a clarifying question moves the world toward `task_routine` with probability 0.4.
5. Re-perception: the now-routine task is perceived as `looks_routine` with probability 0.85.

This path contributes to the transition `looks_risky -> looks_routine` in $Q$.

In plain English: the agent correctly flagged a hard task, asked for more context, the situation improved, and the next reading was routine.

**Why multiple paths help**

Together these paths show that one entry of the composite kernel is not one literal event. It is an aggregation of many possible micro-stories. When you read a probability in $Q$, think: this number compresses many possible world-experience-action-world-experience paths into one effective next-experience probability.

Concise shorthand for reading $Q$ entries:

- `looks_routine -> looks_routine` — the agent correctly handled a manageable task
- `looks_routine -> looks_risky` — the task was harder than it first appeared, or a direct answer did not stabilize it
- `looks_risky -> looks_routine` — clarification or escalation successfully reduced uncertainty
- `looks_risky -> looks_risky` — the problem remained hard even after intervention

### 2. LLM agent with hidden support system

Story:

A focal LLM agent is visible to us, but the rest of the support system is hidden in the background: retrieval services, monitoring tools, human reviewers, and other agents. We only observe whether the focal agent appears valid or invalid from the outside. The hidden system may help or hinder it before we see the focal agent again.

State meanings:
- `A_valid`, `A_invalid` = visible health of the focal agent
- `B_valid`, `B_invalid` = hidden health of the surrounding system

Interpretation:
- the parent kernel models the full visible-plus-hidden system
- the trace onto `{A_valid, A_invalid}` gives the effective observed dynamics of the focal agent alone
- this is not simple deletion of hidden states; it folds hidden excursions into the observed transition probabilities

Code: `examples/trace_analysis/main.go` — walkthrough at `examples/trace_analysis/WALKTHROUGH.md`

### 3. Two-agent validator / repair pair

Story:

One agent performs work and another agent checks or repairs it. Each may itself be in a valid or invalid mode. Rather than specifying the joint kernel by hand, each agent is modeled independently with its own P, D, A triplet. Joint P, D, A kernels are then defined over product state spaces and composed as J = P_joint · D_joint · A_joint — the same PDA composition as a single agent, lifted to the multi-agent level.

State meanings:
- `VV` = both agents reliable
- `VI` = worker reliable, validator degraded
- `IV` = worker degraded, validator reliable
- `II` = both degraded

Interpretation:
- coupling enters through P_joint (Validator observes Worker degradation) and A_joint (Validator repair restores Worker state)
- D_joint = D₁⊗D₂ encodes the "no direct communication" assumption — each agent decides from its own experience
- tracing onto `{VV, II}` gives a coarse healthy-versus-failed operational picture

Code: `examples/validator_repair/main.go` — walkthrough at `examples/validator_repair/WALKTHROUGH.md`

### 4. Three-agent majority-valid coordination network *(not yet implemented)*

Story:

Three agents coordinate on a shared task. As long as at least two are functioning well, the team can usually stabilize itself and recover local failures. Once only one agent remains reliable, recovery becomes much harder and collapse becomes more likely.

State meanings:
- `3V`, `2V1I`, `1V2I`, `3I` = coarse health levels of the team

Interpretation:
- stationary distribution shows how much time the network spends in robust versus degraded regimes
- first-passage time to `3I` measures time-to-systemic-failure
- a trace onto `{3V, 3I}` gives a simplified robust/failure abstraction

### 5. Four-agent pipeline with escalation *(not yet implemented)*

Story:

A four-stage workflow handles incoming work. Local fixes are fast but unreliable, rerouting helps around isolated failures, and escalation is expensive but can stabilize the whole system. The network may move between fully healthy operation, partial degradation, and systemic breakdown.

State meanings:
- all valid
- one invalid
- two invalid
- systemic invalid

Interpretation:
- escalation trades cost for stability
- local fixes trade speed for reliability
- entropy rate measures how noisy or predictable the overall operating regime is

## Tests and examples

The intended style for this project is:

- short mathematical write-up paired with each example
- small finite-state scenario with named states
- runnable example program

Current files:

- `examples/simple_agent/main.go`
- `examples/trace_analysis/main.go`
- `examples/validator_repair/main.go`
- `catrace_test.go`

## Build

Requires Go 1.22+.

```
go build ./...
go test ./...
go run examples/simple_agent/main.go
go run examples/trace_analysis/main.go
```
