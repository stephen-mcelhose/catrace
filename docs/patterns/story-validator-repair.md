# Story: Two-agent validator / repair pair

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

---

Code: `examples/validator_repair/main.go` — walkthrough at `examples/validator_repair/WALKTHROUGH.md`

[← Back to pattern reference](agentic-patterns-reference.md)
