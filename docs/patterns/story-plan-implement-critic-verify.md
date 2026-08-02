# Story: Plan-Implement-Critic-Verify

*Sub-class of: [Plan-and-Execute (12)](agentic-patterns-reference.md), [Evaluator-Optimizer (6)](agentic-patterns-reference.md)*

Story:

A full AI development cycle runs four agents in sequence for each unit of work: a planner structures the task into concrete steps; an implementer executes one step; a critic reads the result and flags issues; a verifier runs automated checks. Only when both the critic and the verifier clear a step does the cycle advance to the next. A failed verification routes back to the implementer. A failed critique may route back to the implementer for a targeted fix, or all the way back to the planner if the critique reveals that the plan step itself was misconceived. The system has multiple feedback loops at different depths — a tight loop for verification failures, a wider loop for critique failures, and the widest for plan revision — and the interaction between them determines whether the system converges quickly or oscillates.

This pattern is a direct composition of Plan-and-Execute and Implement-Critic. What makes it distinct is the presence of *two independent evaluators* with different failure modes: the verifier is objective and cheap; the critic is subjective and expensive. Running them in the right order — verify first, then critique — avoids wasting the critic's capacity on code that will not even type-check.

State meanings:
- world states: `planning`, `implementing`, `awaiting_critique`, `awaiting_verification`, `step_done`, `plan_revision`, `abandoned` — the true stage of the development cycle for the current plan step
- planner experience states: `requirements_clear`, `requirements_partial` — quality of the input specification
- planner actions: `produce_plan`, `revise_plan`, `abandon`
- implementer experience states: `plan_step_clear`, `plan_step_ambiguous` — clarity of the current step
- implementer actions: `implement`, `request_clarification`, `skip_step`
- critic experience states: `code_sound`, `code_has_issues` — perceived code quality
- critic actions: `approve`, `request_changes`, `escalate_to_planner`
- verifier experience states: `checks_passing`, `checks_failing` — automated check results
- verifier actions: `pass`, `fail` — verifier output (deterministic given world state)

Interpretation:
- the verifier is effectively the environment: its action kernel maps the true code state to a pass/fail signal with high fidelity, modelling automated tooling
- the critic's perception is noisy: it may miss real bugs or flag non-issues, making its experience state stochastic given the true world state
- all agent decisions within a stage are independent: $D_\text{joint} = D_\text{planner} \otimes D_\text{impl} \otimes D_\text{critic} \otimes D_\text{verifier}$
- the action kernel encodes routing logic: a `fail` from the verifier combined with any critic action sends the world back to `implementing`; an `escalate_to_planner` from the critic combined with `fail` from the verifier sends it to `plan_revision`
- mean first passage time from `planning` to `step_done` measures per-step latency; multiplied across all plan steps it gives expected total delivery time
- the depth of the feedback loop triggered by each failure type is the key modelling insight: shallow loops (verify → implement) are fast but local; deep loops (critique → plan revision) are slow but can fix structural problems that shallow loops miss

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
