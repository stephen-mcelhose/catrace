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

A worker agent performs tasks while a validator agent monitors its health. Either agent may itself be functioning well or badly. When the validator is healthy, it can detect worker problems and attempt repairs — but repair takes effort and can degrade the validator too. When both are degraded, recovery depends on chance.

Rather than writing down a 4×4 joint transition matrix by hand, each agent is first modeled independently with its own P, D, A kernel triplet. Those individual kernels are then lifted to product state spaces and composed into a joint kernel the same way a single agent works: J = P_joint · D_joint · A_joint. Two coupling points enter explicitly: the validator's perception includes the worker's world state, and the validator's repair action restores the worker's state as well as its own.

State meanings:
- worker world states: `worker_valid`, `worker_invalid` — is the worker actually functioning?
- worker experience states: `sees_ok`, `sees_problem` — does the worker detect its own degradation?
- worker actions: `produce`, `self_check`, `idle`
- validator world states: `validator_valid`, `validator_invalid` — is the validator actually functioning?
- validator experience states: `looks_good`, `looks_bad` — does the validator see a problem?
- validator actions: `validate`, `repair`, `idle`
- joint world states: `VV`, `VI`, `IV`, `II` — the (worker, validator) health pair

Interpretation:
- the validator's perception is coupled to the worker's world state: a degraded worker shifts the validator's experience toward `looks_bad` even if the validator itself is fine — this is where cross-agent observation enters the model
- the validator's repair action is coupled to the worker's world state: repair boosts the worker's probability of being valid, not just the validator's own — this is where cross-agent effect enters
- decisions are independent: D_joint = D₁⊗D₂, so each agent decides from its own experience without communicating
- tracing onto `{VV, II}` collapses the mixed states and gives a coarse healthy-versus-failed picture of the pair

Code: `examples/validator_repair/main.go` — walkthrough at `examples/validator_repair/WALKTHROUGH.md`

#### Played-out version: Story 3

The joint kernel J compresses a full W→X→G→W cycle — for both agents simultaneously — into one effective joint-state transition. Walking concrete paths shows how the coupling between agents shapes that compression.

**Version A: stable healthy system**

1. The system is in `VV` — both agents functioning.
2. Perception: the worker experiences `sees_ok` (probability 0.90); the validator, observing a healthy worker, experiences `looks_good` (probability 0.85). Joint experience: `ok·good`.
3. Decision: the worker chooses `produce` (probability 0.80); the validator chooses `validate` (probability 0.60). Joint action: `produce|validate`.
4. Action effect: producing keeps the worker valid (probability 0.70); validating keeps the validator calibrated (probability 0.85). Joint next world: `VV` with probability 0.70 × 0.85 = 0.595.

This path contributes the bulk of probability mass to the `VV → VV` transition in J.

In plain English: both agents were fine, both perceived no problem, they did their normal work, and the system stayed fine.

**Version B: worker degrades undetected**

1. The system is in `VV`.
2. Perception: the worker sees `sees_ok` (0.90); the validator, observing a healthy worker, sees `looks_good` (0.85). Nothing looks wrong yet.
3. Decision: the worker produces (0.80); the validator validates (0.60).
4. Action effect: this time producing fails to maintain worker validity — the worker degrades with probability 0.30. The validator stays calibrated (0.85). Joint next world: `IV` with probability 0.30 × 0.85 = 0.255.

This path contributes to the `VV → IV` transition.

In plain English: everything looked fine, both agents acted normally, but the worker degraded anyway — and the validator, watching a still-healthy worker at the start of the cycle, had no signal to trigger a repair.

**Version C: coupled recovery from IV**

1. The system is in `IV` — worker degraded, validator healthy.
2. Perception: the worker sees `sees_problem` (probability 0.70). The validator, observing a degraded worker, sees `looks_bad` with elevated probability 0.40 — higher than the 0.15 it would have if the worker were fine. Joint experience: `prob·bad`.
3. Decision: the worker self-checks (0.60); the validator repairs (0.70). Joint action: `self_check|repair`.
4. Action effect: validator repair boosts the worker's probability of returning to valid — 0.70 instead of the 0.50 it could manage independently. Repair taxes the validator, holding it valid with probability 0.60. Joint next world: `VV` with probability 0.70 × 0.60 = 0.420.

This path contributes to the `IV → VV` transition.

In plain English: the worker was degraded and the validator could see it. Because the validator observed the worker's condition — not just its own — it knew to repair. The repair helped the worker too, and both ended up healthy.

**Why the coupling matters**

Without coupled perception, the validator in state `IV` would see `looks_bad` with only 0.15 probability (its own baseline) rather than 0.40 — less than a third as likely to trigger repair. Without coupled action, the validator's repair would have no effect on the worker at all. Together, the two coupling points give the system a recovery pathway that neither agent has alone.

Concise shorthand for reading J entries:
- `VV → VV` — both agents working normally; system holds
- `VV → IV` — worker degraded despite a healthy validator; no signal triggered repair in time
- `IV → VV` — validator detected worker degradation and repaired it; full recovery in one cycle
- `II → VV` — even from full degradation, coupled repair can restore both agents in one step

### 4. Self-adjusting / self-healing network nodes

Story:

A network node monitors its own error rate and throttles itself when errors climb. An outer evolutionary loop watches pool throughput and mutates the node's configuration when performance drops. The two loops compete to explain recovery: in one regime the node's own throttle is the primary healer; in the other the evolver's config search is what keeps the system alive. Running the same joint kernel architecture under both parameter regimes lets the stationary distribution, mean first passage time, and entropy rate vote on which mechanism does the work.

State meanings:
- node world states: `healthy`, `degraded`, `overloaded` — actual error rate of the node
- node experience states: `ema_low`, `ema_mid`, `ema_high` — EMA band the node observes
- node actions: `push`, `throttle`, `idle`
- evolver world states: `good_strategy`, `bad_strategy` — whether the current MaxWorkers/Kp config is effective
- evolver experience states: `high_score`, `low_score` — pool-level throughput×success² the evolver observes
- evolver actions: `promote`, `mutate`
- joint world states: `H·G`, `H·B`, `D·G`, `D·B`, `O·G`, `O·B` — (node health · evolver strategy) pairs

Interpretation:
- the evolver's perception is coupled to node health in $P_\text{joint}$: a sick node depresses the pool score the evolver observes, making a good config look bad even when it isn't
- decisions are independent: $D_\text{joint} = D_\text{node} \otimes D_\text{evolver}$; each agent reads its own signal without communicating within a cycle
- when the evolver mutates, it boosts the node's recovery probability in $A_\text{joint}$; `promote` leaves node recovery unchanged
- the variant comparison pattern uses the model as a measuring instrument: run both parameter regimes, read the difference in MFPT and entropy rate, and let the numbers say which mechanism carries the load

Code: `examples/self_healing_nodes/main.go` — walkthrough at `examples/self_healing_nodes/WALKTHROUGH.md`

#### Played-out version: Story 4

The joint kernel $J = P_\text{joint} \cdot D_\text{joint} \cdot A_\text{joint}$ compresses a full $W \to X \to G \to W$ cycle — for both agents simultaneously — into one effective joint-state transition. The paths below use Variant B parameters, where the evolver's contribution is most visible.

**Version A: statistically typical path**

1. The system is in `H·G` — node healthy, evolver running a good config.
2. Perception: the healthy node observes `ema_low` (probability 0.80); the evolver, watching a healthy node's pool score, reads `high_score` (probability 0.85). Joint experience: `ema_low·high_score`, probability $0.80 \times 0.85 = 0.680$.
3. Decision: reading `ema_low`, the node pushes (probability 0.75); reading `high_score`, the evolver elects to promote its current config (probability 0.80). Joint action: `push·promote`, probability $0.75 \times 0.80 = 0.600$.
4. Action effect: `push` keeps the node `healthy` (probability 0.55); `promote` preserves the `good_strategy` (probability 0.85). Joint next world: `H·G`, probability $0.55 \times 0.85 = 0.468$.

This path contributes the bulk of probability mass to the `H·G → H·G` transition in J.

In plain English: the node was healthy, the EMA stayed quiet, the evolver saw good scores and kept the config, and the system remained in its best state.

**Version B: degraded node masks a good config**

1. The system is in `D·G` — the node has drifted into a degraded state, but the evolver's config is actually good.
2. Perception: the degraded node reads `ema_mid` (probability 0.55); the pool score the evolver observes is suppressed by the degradation — even a good config reads as `low_score` with probability 0.45. Joint experience: `ema_mid·low_score`, probability $0.55 \times 0.45 = 0.248$.
3. Decision: reading `ema_mid`, the node throttles (probability 0.55); reading `low_score`, the evolver mutates away from the (actually good) config (probability 0.75). Joint action: `thr·mutate`, probability $0.55 \times 0.75 = 0.413$.
4. Action effect: the mutation boost raises the node's `healthy` probability from 0.45 to 0.56 after renormalization; mutating gives 0.60 probability of landing on a `good_strategy`. Joint next world: `H·G`, probability $0.56 \times 0.60 = 0.336$.

This path contributes to the `D·G → H·G` transition in J.

In plain English: the node's degradation made a good config appear bad; the evolver unnecessarily searched for a new one, but the mutation boost it applied helped the node recover along the way.

**Version C: recovery from the worst state**

1. The system is in `O·B` — the node is overloaded and the evolver has a bad config.
2. Perception: the overloaded node reads `ema_high` (probability 0.75); a bad strategy with an overloaded node always shows `low_score` (probability 1.00). Joint experience: `ema_high·low_score`, probability $0.75 \times 1.00 = 0.750$.
3. Decision: reading `ema_high`, the node idles (probability 0.45); reading `low_score`, the evolver mutates (probability 0.75). Joint action: `idle·mutate`, probability $0.45 \times 0.75 = 0.338$.
4. Action effect: the mutation boost raises the node's `healthy` probability from 0.60 to 0.667 after renormalization; mutating gives 0.60 probability of landing on `good_strategy`. Joint next world: `H·G`, probability $0.667 \times 0.60 = 0.400$.

This path contributes to the `O·B → H·G` transition in J.

In plain English: both mechanisms fired together — the node backed off under EMA pressure and the evolver searched for a better config; the mutation boost made the combined recovery faster than either would have managed alone.

**Why the coupling matters**

Without the perception coupling, a degraded node would not suppress the pool score — the evolver would keep promoting a good config rather than accidentally mutating away from it, and the `D·G → D·B` leakage would disappear. Without the action coupling, mutation would have no effect on node recovery. Together they create a system where the outer loop can accidentally harm state estimation (perception coupling) and simultaneously provide the primary recovery pathway (action coupling). The variant comparison reads this tension directly: in Variant B the evolver is essential; in Variant A the throttle is strong enough that the outer loop is a minor boost on top of self-healing that's already working.

Concise shorthand for reading J entries:
- `H·G → H·G` — node healthy, evolver keeping the best config; system holds
- `D·G → D·B` — degraded node suppresses the score; evolver mutates away from the good config
- `D·G → H·G` — throttle recovers the node; mutation boost (if large enough) accelerates it
- `O·B → H·G` — idle plus a lucky mutation search finds the recovery path in one step
- `O·B → O·B` — throttle fails and mutation finds no better config; system stays stuck

### 5. Three-agent majority-valid coordination network *(not yet implemented)*

Story:

Three agents coordinate on a shared task. As long as at least two are functioning well, the team can usually stabilize itself and recover local failures. Once only one agent remains reliable, recovery becomes much harder and collapse becomes more likely.

State meanings:
- `3V`, `2V1I`, `1V2I`, `3I` = coarse health levels of the team

Interpretation:
- stationary distribution shows how much time the network spends in robust versus degraded regimes
- first-passage time to `3I` measures time-to-systemic-failure
- a trace onto `{3V, 3I}` gives a simplified robust/failure abstraction

### 6. Four-agent pipeline with escalation *(not yet implemented)*

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
