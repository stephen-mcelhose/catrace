# Story: Supervisor / hierarchical agent network

Story:

A software engineering team operates under a tech lead who monitors project health and coordinates a pool of developer agents. Each developer runs its own perceive→decide→act loop on assigned tasks. The tech lead observes aggregate signals — build health, PR velocity, blocked-ticket count — and intervenes by redistributing work, unblocking stuck developers, or escalating to management when the situation is beyond the team's self-repair capacity. The interesting question is how often the tech lead's intervention is actually necessary: in a healthy team, developers self-unblock often enough that the supervisor's main job is to hold course; in a degraded team, the supervisor's redistribution is the only mechanism that prevents cascading failure.

State meanings:
- team world states: `healthy`, `one_blocked`, `two_blocked`, `critical` — how many developers are blocked on tasks they cannot self-unblock
- supervisor experience states: `velocity_good`, `velocity_degraded` — what the supervisor reads from aggregate PR and build signals
- supervisor actions: `hold_course`, `redistribute`, `escalate` — the supervisor's available interventions
- developer world states: `unblocked`, `blocked` — whether the developer can make progress independently
- developer experience states: `task_clear`, `task_blocked` — what the developer perceives about its current assignment
- developer actions: `progress`, `request_help`, `idle` — developer responses

Interpretation:
- the supervisor's perception is coupled to the aggregate developer world state: multiple blocked developers suppress velocity signals, making a degraded team legible from the outside
- developer decisions are independent: $D_\text{joint} = D_\text{sup} \otimes D_{d_1} \otimes \cdots \otimes D_{d_n}$; no developer coordinates with another within a cycle
- the action kernel couples supervisor and developers: `redistribute` shifts probability mass from `blocked` toward `unblocked` for targeted developers; `escalate` brings an external recovery boost at the cost of management overhead
- tracing onto `{healthy, critical}` collapses intermediate degradation states and gives a coarse functioning-versus-failing picture of the team
- first passage time from `healthy` to `critical` measures how long the team can sustain itself without supervisor intervention

---

Issue: [#9 — example: add supervisor / hierarchical agent network example](https://github.com/stephen-mcelhose/catrace/issues/9)

[← Back to pattern reference](agentic-patterns-reference.md)
