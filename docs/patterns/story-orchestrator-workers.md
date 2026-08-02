# Story: Orchestrator-workers

Story:

A research orchestrator receives a complex investigation brief and must coordinate several worker agents to complete it. The orchestrator cannot do the research itself — it can only assign subtasks, monitor reported status, and eventually synthesise results. Workers may succeed, return partial findings, or fail silently. The orchestrator perceives aggregate task progress, decides whether to assign new work, retry a failed subtask, or synthesise with what it has. The tension is between waiting for complete results and synthesising early with incomplete information: waiting too long stalls delivery; synthesising too early produces a shallow report.

State meanings:
- orchestrator world states: `no_tasks_done`, `some_tasks_done`, `all_tasks_done`, `synthesis_done` — true task completion level
- orchestrator experience states: `progress_on_track`, `progress_stalled`, `progress_near_complete` — what the orchestrator perceives from worker status reports
- orchestrator actions: `assign_tasks`, `retry_failed`, `synthesise` — the orchestrator's available moves
- worker world states: `idle`, `working`, `complete`, `failed` — true worker status
- worker experience states: `task_clear`, `task_ambiguous` — how well the worker understands its assigned subtask
- worker actions: `execute`, `request_clarification`, `report_done`, `report_fail` — worker responses

Interpretation:
- the orchestrator's perception is coupled to the aggregate worker world state: many failed workers drive the orchestrator's experience toward `progress_stalled` even if one thread is quietly succeeding
- worker decisions are independent within a cycle: $D_\text{joint} = D_\text{orch} \otimes D_{w_1} \otimes \cdots \otimes D_{w_n}$
- the action kernel couples orchestrator and workers: a `retry_failed` action from the orchestrator resets failed workers toward `idle`, giving them another chance
- the world kernel $W = PDA$ gives the full joint-state transition, including the probability of reaching `synthesis_done` from each starting configuration
- mean first passage time from `no_tasks_done` to `synthesis_done` measures expected delivery latency under different retry and synthesis policies

---

Issue: [#8 — example: add orchestrator-workers example](https://github.com/stephen-mcelhose/catrace/issues/8)

[← Back to pattern reference](agentic-patterns-reference.md)
