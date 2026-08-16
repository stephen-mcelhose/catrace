# Story: Blackboard / shared-workspace collaboration

Story:

Three specialist agents — a radiologist, a pathologist, and a clinical-notes reader — collaborate on a case by reading and writing a shared diagnostic board. They never message each other: coordination is only through the board. Nobody sees anyone else's raw inputs, only the board's current status. When the board is empty, specialists rarely feel triggered to contribute. When a tentative hypothesis is already posted, the board contents raise the chance that a specialist's precondition fires — they feel evidence strong enough to endorse, post a refinement, or flag a contradiction. The case confirms when enough endorsements accumulate; it contradicts when someone flags irreconcilable findings. Confirmed diagnoses stay closed; contradicted cases sometimes reopen for another pass. Roles are labels only in v1 — the three specialists share the same perception and decision kernels so the pattern (opportunistic contribution via shared state) is not confused with heterogeneous accuracy.

Rather than wiring a point-to-point protocol among specialists, each reads the shared board, decides independently, and writes back. Joint dynamics are the ordinary PDA lift over product experience and action spaces with a single shared world: $J = P_\text{joint} \cdot D_\text{joint} \cdot A_\text{joint}$, where $D_\text{joint} = D^{\otimes 3}$.

State meanings:
- blackboard world states: `undiagnosed`, `tentative_diagnosis`, `confirmed_diagnosis`, `contradicted` — collective agreement status on the board (`confirmed_diagnosis` absorbing; `contradicted` may reopen to `undiagnosed`)
- specialist experience states: `evidence_strong`, `evidence_weak` — whether the current board looks actionable to this specialist (precondition fires vs sits out)
- specialist actions: `post_finding`, `endorse_prior`, `flag_contradiction`, `request_more_data` — contribute, support, challenge, or abstain-with-ask
- joint experience / action: product of the three specialists' X and G — they act in the same cycle without seeing each other's private evidence

Interpretation:
- perception couples every specialist to the **shared board only**: `tentative_diagnosis` raises $P(\texttt{evidence_strong})$ (board contents trigger contribution); `undiagnosed` keeps most mass on `evidence_weak`. No specialist observes a peer's experience
- decisions are independent: $D_\text{joint} = D^{\otimes 3}$; out-of-band agent-to-agent calls are forbidden — the catalog Blackboard forbid
- the action kernel encodes opportunistic accretion rules on joint posts: any `post_finding` from `undiagnosed` tends to open `tentative_diagnosis`; from `tentative_diagnosis`, two or more `endorse_prior` tend toward `confirmed_diagnosis`, while any `flag_contradiction` elevates `contradicted`
- the world kernel $J = P_\text{joint} D_\text{joint} A_\text{joint}$ is the board's status dynamics; mean first passage from `undiagnosed` to `confirmed_diagnosis` measures diagnostic latency
- tracing onto `{undiagnosed, confirmed_diagnosis, contradicted}` collapses the intermediate agreement state into a coarse start / confirmed / contradicted view

Non-goals (v1):
- ground-truth disease label in W (board status is the world)
- heterogeneous specialist P/D/A (modality accuracy)
- within-cycle turn order / explicit control shell picking who writes next
- write-race conflict policies and board pruning (catalog costs; park for a later experiment)

#### Played-out version: Story 7

Each cycle is three specialists reading the board, deciding independently, and writing jointly — compressed into one board-status transition of $J$.

**Version A: first finding opens the board**

1. Starting world: `undiagnosed` — empty case board.
2. Perception: each specialist sees `evidence_weak` with probability 0.70 (board does not yet trigger). Joint experience often has zero or one `evidence_strong`.
3. Decision: a specialist with `evidence_strong` chooses `post_finding` with probability 0.30; weak specialists mostly `request_more_data` or `flag_contradiction`.
4. Action effect: if at least one `post_finding` appears in the joint action, the board moves to `tentative_diagnosis` with high probability (0.80 in the action rule).

This path contributes the bulk of mass to `undiagnosed → tentative_diagnosis` in $J$.

In plain English: someone finally had something to write; the shared workspace now carries a working hypothesis.

**Version B: tentative board triggers endorsements**

1. Starting world: `tentative_diagnosis` — a hypothesis is on the board.
2. Perception: each specialist sees `evidence_strong` with probability 0.75 — the board contents fire preconditions.
3. Decision: strong evidence chooses `endorse_prior` with probability 0.45 (or `post_finding` 0.30).
4. Action effect: two or more endorsements in the joint action move toward `confirmed_diagnosis` with probability 0.80.

This path contributes to `tentative_diagnosis → confirmed_diagnosis` in $J$.

In plain English: the board made contribution actionable; enough specialists backed the working hypothesis and the case closed.

**Version C: flag against a tentative hypothesis**

1. Starting world: `tentative_diagnosis`.
2. Perception: most see `evidence_strong` (0.75), but some still see `evidence_weak` (0.25).
3. Decision: weak evidence chooses `flag_contradiction` with probability 0.40.
4. Action effect: any flag in the joint action elevates `contradicted` (probability 0.60 under the action rule).

This path contributes to `tentative_diagnosis → contradicted` in $J`.

In plain English: once a hypothesis is posted, a single challenge can divert the case into contradiction instead of confirmation — opportunistic accretion cuts both ways.

**Why the coupling matters**

Without board-coupled perception, specialists would not become more triggerable when a tentative hypothesis appears, and endorsement mass would not concentrate after the first post. Without joint action rules that count posts, endorsements, and flags, the shared workspace would be costume only. Together they are the blackboard mechanism: coordinate through state, contribute when the board enables you.

Concise shorthand for reading $J$ entries:
- `undiagnosed → tentative_diagnosis` — a specialist posted the first finding
- `tentative_diagnosis → confirmed_diagnosis` — enough endorsements closed the case
- `tentative_diagnosis → contradicted` — a flag challenged the working hypothesis
- `contradicted → undiagnosed` — the contradicted case reopened for another pass
- `confirmed_diagnosis → confirmed_diagnosis` — closed diagnosis stays closed

---

Code: `examples/blackboard/main.go` — walkthrough at `examples/blackboard/WALKTHROUGH.md`

Issue: [#11 — example: add blackboard / shared-workspace collaboration example](https://github.com/stephen-mcelhose/catrace/issues/11)

[← Back to pattern reference](agentic-patterns-reference.md)
