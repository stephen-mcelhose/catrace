# Story: Plan-and-execute

Story:

A travel booking agent receives a complex multi-leg trip request. A planner agent constructs a full itinerary upfront — flights, hotels, car hire, and activities in sequence — before any booking is attempted. An executor agent then works through the plan step by step, attempting each booking in order. When a step fails due to unavailability, the executor reports back to the planner, which must revise the downstream itinerary without re-doing the already-confirmed steps. The tension is between the planner's investment in a globally coherent plan and the executor's reality: the world changes between planning and execution, and a rigid plan that cannot absorb partial failure costs more to revise than a looser one would have.

State meanings:
- world states: `planning`, `executing_step_1`, `executing_step_2`, `executing_step_3`, `partial_failure`, `complete`, `abandoned` — the true stage of the booking process
- planner experience states: `full_context`, `partial_context` — how complete the planner's view of availability and constraints is at planning time
- planner actions: `produce_plan`, `revise_plan`, `abandon` — planner responses to the current situation
- executor experience states: `step_available`, `step_unavailable` — what the executor finds when it attempts the next booking
- executor actions: `confirm_booking`, `skip_step`, `report_failure` — executor responses to each booking attempt

Interpretation:
- the planner's perception is coupled to the executor's world state: a partial failure drives the planner's experience toward `partial_context`, forcing a revision cycle
- decisions are sequential rather than simultaneous: the planner acts first, then the executor; within a cycle each agent reads its own signal
- the action kernel captures how planner and executor actions jointly advance the booking world state: a `report_failure` from the executor triggers a `revise_plan` from the planner, resetting to an earlier execution stage
- the world kernel $W = PDA$ gives the full booking-state transition matrix; mean first passage time from `planning` to `complete` measures expected booking latency under different revision-tolerance policies
- commute time between `partial_failure` and `complete` measures the round-trip cost of a failed step: how long it takes to reach failure from a mid-execution state and recover back to completion

---

Issue: [#13 — example: add plan-and-execute example](https://github.com/stephen-mcelhose/catrace/issues/13)

[← Back to pattern reference](agentic-patterns-reference.md)
