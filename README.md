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

### 6. Prompt-chaining document pipeline

Story:

A document intelligence pipeline sends raw text through three specialist stages in sequence: an extractor identifies key claims, a summariser condenses them into a structured brief, and a formatter produces the final delivery-ready report. Between stages a quality gate checks the intermediate output; if it fails, the stage reruns with slightly different framing before passing the result on. The pipeline can stall in a retry loop or abandon a document entirely if the gate repeatedly rejects.

Rather than modeling the pipeline as a single deterministic graph, each stage is a perception-decision-action loop. The stage agent perceives input quality imperfectly — a degraded input may look processable; a clean one may appear noisy — and its policy determines whether to advance, retry, or escalate. The world kernel W = PDA captures the full stage-to-stage transition matrix including retry loops and failure exits.

- world states: `raw`, `extracted`, `summarised`, `formatted`, `failed` — the true pipeline stage the document has reached
- experience states: `input_clear`, `input_noisy` — whether the stage agent perceives its input as processable
- actions: `process`, `retry`, `escalate` — attempt the transformation, reframe and try again, or abandon

Interpretation:

- perception captures how reliably a stage agent reads the quality of its incoming material; `raw` documents are mostly perceived as noisy, `formatted` documents as clear
- decision captures the agent's policy given perceived input quality; `input_clear` leads mostly to `process`, `input_noisy` splits between `retry` and `escalate`
- action captures how the chosen action advances (or fails to advance) the pipeline world state
- tracing onto `{raw, formatted, failed}` collapses the intermediate stages and shows the pipeline's end-to-end success and failure rate directly

Code: `examples/prompt_chaining/main.go` — walkthrough at `examples/prompt_chaining/WALKTHROUGH.md`

#### Played-out version: Story 6

**Version A: statistically typical path**

1. World state `extracted` (29.6% of steady-state time — the most common stage)
2. P[`extracted`, `input_clear`] = 0.60 → agent perceives `input_clear`
3. D[`input_clear`, `process`] = 0.80 → agent chooses `process`
4. A[`process`, `summarised`] = 0.35 → world advances to `summarised` (0.60 × 0.80 × 0.35 = 0.168)

This path contributes the bulk of the `extracted → summarised` entry in W. In plain English: a document whose claims have been extracted is perceived as processable, the agent attempts the summarisation, and it advances to the structured-brief stage.

**Version B: failure path**

1. World state `raw` — unprocessed document entering the pipeline
2. P[`raw`, `input_noisy`] = 0.75 → agent perceives `input_noisy`
3. D[`input_noisy`, `escalate`] = 0.25 → agent chooses `escalate`
4. A[`escalate`, `failed`] = 0.80 → world moves to `failed` (0.75 × 0.25 × 0.80 = 0.15)

This path contributes to W[`raw`, `failed`] = 0.20. In plain English: a raw document looks unprocessable to the agent, the agent gives up rather than retrying, and the document exits the pipeline as a failure.

**Version C: recovery from failure**

1. World state `failed` (14.9% steady-state — the pipeline occasionally revisits failure)
2. P[`failed`, `input_clear`] = 0.15 → agent perceives `input_clear` (rare, but possible)
3. D[`input_clear`, `process`] = 0.80 → agent chooses `process`
4. A[`process`, `extracted`] = 0.30 → world moves to `extracted` (0.15 × 0.80 × 0.30 = 0.036)

This path contributes to W[`failed`, `extracted`]. In plain English: a failed document occasionally re-enters the pipeline when the stage agent happens to perceive its input as workable — even from failure, the chain is not absorbing.

The three paths together show why a 5-state pipeline cannot be understood by inspection: the stationary distribution (12.3% `raw`, 29.6% `extracted`, 26.7% `summarised`, 16.5% `formatted`, 14.9% `failed`) arises from the interaction of all three paths simultaneously. The coarse trace onto {`raw`, `formatted`, `failed`} normalises to 28.2% / 37.8% / 34.0% — the pipeline succeeds roughly as often as it fails at the end-state level.

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
