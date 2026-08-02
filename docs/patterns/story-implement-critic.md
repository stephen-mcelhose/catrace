# Story: Implement-Critic

*Sub-class of: [Evaluator-Optimizer (6)](agentic-patterns-reference.md)*

Story:

A coding agent drafts an implementation and submits it to a separate critic agent for review. The critic cannot run the code — it can only read it. It evaluates for correctness, security vulnerabilities, adherence to project conventions, and alignment with the original requirements. The implementer revises based on the critique and resubmits. The dynamic is adversarial but cooperative: the critic is actively looking for problems the implementer may have missed; the implementer is trying to fix real issues, not just satisfy the critic's surface preferences. The failure mode on the critic's side is false precision — flagging style issues as blocking problems and delaying delivery without improving correctness. The failure mode on the implementer's side is mechanical compliance — making exactly the changes the critic asked for without understanding why, which often introduces new problems.

Unlike Implement-Verify, the evaluator here is itself an agent with its own perception kernel and policy. The critic's perception may be imperfect: it can miss real bugs and flag non-issues. This makes the joint dynamics richer — the critic's accuracy is itself a stochastic variable that shapes the long-run outcome.

State meanings:
- world states: `draft`, `under_review`, `revision_requested`, `approved`, `escalated` — the true state of the implementation review cycle
- implementer experience states: `requirements_clear`, `requirements_ambiguous` — whether the implementer had a clear specification to work from
- implementer actions: `implement`, `revise_targeted`, `revise_broad`, `push_back` — implementer responses to critique
- critic experience states: `code_sound`, `code_has_issues`, `code_unclear` — what the critic perceives after reading the implementation
- critic actions: `approve`, `request_changes`, `escalate_to_human` — critic decisions

Interpretation:
- the critic's perception is coupled to the true world state: `draft` quality shapes what the critic perceives, but with noise — a subtly broken implementation may look fine; a non-idiomatic but correct one may appear broken
- implementer and critic decisions are independent within a cycle: the implementer does not know what the critic will flag before submitting; $D_\text{joint} = D_\text{impl} \otimes D_\text{critic}$
- the action kernel captures how each (implementer action, critic action) pair advances the review state: a `revise_targeted` paired with `request_changes` sends the cycle around again; a `revise_broad` has a higher probability of reaching `approved` but risks introducing new issues visible to the critic
- the stationary distribution shows how much time the cycle spends in `revision_requested` versus `approved` — a high `revision_requested` mass indicates either poor initial implementation quality or an over-zealous critic
- commute time between `draft` and `approved` measures the average round-trip cost of one full review cycle

---

Issue: [#16 — example: add implement-critic example](https://github.com/stephen-mcelhose/catrace/issues/16)

[← Back to pattern reference](agentic-patterns-reference.md)
