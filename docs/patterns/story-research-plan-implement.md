# Story: Research-Plan-Implement

*Sub-class of: [Prompt Chaining (2)](agentic-patterns-reference.md), [Plan-and-Execute (12)](agentic-patterns-reference.md)*

Story:

A software engineering agent receives a feature request for an unfamiliar codebase. Before writing any code, it explores relevant files, reads documentation, and builds an internal model of the existing architecture. It then drafts a structured implementation plan — which files to change, which interfaces to add, which tests to write. Only once the plan is complete does it execute each step in order. The risk compounds across stages: incomplete exploration produces a plan that conflicts with codebase reality, and a flawed plan sends the implementer down paths it cannot recover from without restarting. An agent that skips exploration to save cycles often produces more rework than one that spends time up front.

State meanings:
- world states: `unexplored`, `explored`, `planned`, `implementing`, `done`, `blocked` — the true stage of the development task
- research experience states: `context_rich`, `context_sparse` — how complete the agent's model of the codebase is after exploration
- research actions: `explore_more`, `commit_to_plan` — whether to keep gathering or move forward
- planner experience states: `plan_coherent`, `plan_conflicts_detected` — whether the drafted plan is internally consistent with what was found
- planner actions: `produce_plan`, `revise_plan`, `abandon` — planner responses
- implementer experience states: `step_clear`, `step_ambiguous` — whether the current implementation step is well-specified
- implementer actions: `implement_step`, `backtrack`, `request_replan` — implementer responses

Interpretation:
- the research-to-plan transition is the first gate: moving forward with `context_sparse` raises the probability of `plan_conflicts_detected`, which forces a costly revision cycle
- the plan-to-implement transition is the second gate: `plan_conflicts_detected` before execution shifts the implementer's experience toward `step_ambiguous` throughout the run
- the planner's perception is coupled to the research world state: more exploration reduces plan conflict probability; the coupling is the core design decision in this pattern
- decisions within a stage are independent: the researcher does not consult the planner within a cycle; the planner does not consult the implementer
- the world kernel $W = PDA$ gives full stage-transition dynamics including retry loops; mean first passage time from `unexplored` to `done` measures expected total task latency under different exploration-budget policies

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
