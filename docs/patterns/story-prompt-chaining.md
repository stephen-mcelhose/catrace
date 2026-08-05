# Story: Prompt chaining pipeline

Story:

A diligence desk turns a raw filing into a client-ready report with a fixed prompt chain — the workflow Anthropic calls prompt chaining, not an autonomous agent. Three specialised LLM calls run in order: an extractor pulls key claims from the filing, a summariser turns accepted claims into a structured brief, and a formatter produces the delivery-ready report. Each call sees only the artifact handed to it (not upstream sources or downstream goals). Between calls, ordinary program code acts as a gate: schema / completeness checks that either accept the intermediate and pass it forward, reject it so that stage retries, or escalate the job to a human review queue. The human queue may later re-queue work back to the desk. The path is fixed at design time; what varies is whether each link emits something the gate will accept.

State meanings:
- pipeline world states: `raw`, `extracted`, `summarised`, `formatted`, `failed` — where the artifact sits (`formatted` is shipped/absorbing; `failed` is the human queue, which may re-queue work back to `raw`)
- extractor experience: `filing_clear`, `filing_noisy` — how readable the raw filing looks
- extractor actions: `emit_claims`, `retry_extract` — produce a claims artifact, or reframe and try again
- summariser experience: `claims_clear`, `claims_noisy` — how usable the gated claims look
- summariser actions: `emit_brief`, `retry_summarise`
- formatter experience: `brief_clear`, `brief_noisy` — how usable the gated brief looks
- formatter actions: `emit_report`, `retry_format`
- gate outcomes (programmatic, not an agent): `pass`, `retry_stage`, `escalate` — accept and advance the pipeline world, keep the artifact at the current stage for another LLM attempt, or move to `failed`

Interpretation:
- this is three chained LLM steps plus code gates — matching the usual prompt-chaining workflow — not one shared stage-worker policy and not three peers acting every timestep
- each stage has its own perception and decision kernels: different prompts, different failure modes (extraction misses vs summary drift vs format breakage)
- only the active stage’s call matters on a given step; inactive stages do not “decide” the artifact forward
- the gate is not a fourth LLM: it is a deterministic (or near-deterministic) check on the emitted intermediate. `pass` advances `raw → extracted → summarised → formatted`; `retry_stage` leaves the pipeline world unchanged so the same stage can fire again; `escalate` sends mass to `failed`
- coupling is handoff through the artifact and the gate, not cross-perception among specialists: the summariser never observes the raw filing, only claims that already passed a gate
- the pipeline world kernel is assembled row-by-row from the active stage’s perception × decision × gate (not a single shared `Agent.WorldKernel()`), yielding stage-to-stage probabilities including retry loops and human-queue exits
- tracing onto `{raw, formatted, failed}` collapses mid-chain stages into start / shipped / human-queue
- mean first passage from `raw` to `formatted` measures expected steps to ship (latency cost of the chain)

Non-goals (v1):
- a single shared `(P, D, A)` for all stages (dropped; that is not how prompt chaining is described)
- treating the gate as an LLM critic / evaluator-optimizer loop (that is a different pattern)
- autonomous routing or planner choosing the next stage at runtime (that is routing / orchestrator-workers)
- full simultaneous three-agent product where every specialist acts every tick regardless of pipeline position

#### Played-out version: Story 6

Each step of the chain is one active stage’s perception → decision → gate, which becomes one row-update of the pipeline world kernel $W$.

**Version A: statistically typical first link**

1. Starting world: `raw` — a new filing is on the desk.
2. Perception: the extractor reads `filing_clear` with probability 0.55.
3. Decision: given `filing_clear`, it chooses `emit_claims` with probability 0.85.
4. Gate: on emit, `pass` with probability 0.70 advances the artifact to `extracted`.

Path probability: $0.55 \times 0.85 \times 0.70 = 0.32725$, the largest single contribution to `W[raw, extracted] = 0.406`.

In plain English: a readable filing, a confident extract, and a gate that accepts the claims — the first link of the chain fires cleanly.

**Version B: noisy filing burns a retry**

1. Starting world: `raw`.
2. Perception: the extractor reads `filing_noisy` with probability 0.45.
3. Decision: given `filing_noisy`, it chooses `retry_extract` with probability 0.75.
4. Gate / retry fate: after a retry action, stay at `raw` with probability 0.90.

Path probability: $0.45 \times 0.75 \times 0.90 = 0.30375$, the largest single contribution to `W[raw, raw] = 0.494`.

In plain English: a messy filing makes the extractor reframe instead of emitting; the artifact does not advance, so the same stage will fire again.

**Version C: middle-stage advance after a gate**

1. Starting world: `extracted` — claims already passed the first gate.
2. Perception: the summariser reads `claims_clear` with probability 0.75 (upstream gating helps).
3. Decision: given `claims_clear`, it chooses `emit_brief` with probability 0.90.
4. Gate: on emit, `pass` with probability 0.80 advances to `summarised`.

Path probability: $0.75 \times 0.90 \times 0.80 = 0.540$, the bulk of `W[extracted, summarised] = 0.61`.

In plain English: once claims are gated, the summariser usually sees clear input and the next gate usually lets the brief through — later links are easier than the first.

**Why this matters**

No single path is the pipeline. Version A moves work forward; Version B shows how retries inflate MFPT without ever calling a second specialist; Version C shows why gated intermediates look clearer than raw filings. Together they explain a pipeline $W$ whose stationary mass sits on `formatted` while MFPT(`raw`→`formatted`) ≈ 7.7 steps measures desk latency under retry and escalate pressure.

Concise shorthand:
- `raw → extracted` — extractor emit + gate pass
- `raw → raw` — retry_extract or gate retry_stage
- `raw → failed` — escalate at the first link
- `extracted → summarised` — summariser emit + gate pass
- `failed → raw` — human re-queues work onto the desk

---

Code: `examples/prompt_chaining/main.go` — walkthrough at `examples/prompt_chaining/WALKTHROUGH.md`

Issue: [#5 — example: add prompt-chaining pipeline example](https://github.com/stephen-mcelhose/catrace/issues/5)

[← Back to pattern reference](agentic-patterns-reference.md)
