---
title: Structural Patterns
tags: [patterns, structural, topology, agentic, markov, modeling]
sources: [docs/patterns/agentic-patterns-reference.md, docs/patterns/story-single-llm-agent.md, docs/patterns/story-hidden-support-system.md, docs/patterns/story-validator-repair.md, docs/patterns/story-self-healing-nodes.md, docs/patterns/story-prompt-chaining.md, docs/patterns/story-routing.md, docs/patterns/story-parallelisation.md, docs/patterns/story-orchestrator-workers.md, docs/patterns/story-supervisor.md, docs/patterns/story-swarm.md, docs/patterns/story-blackboard.md, docs/patterns/story-debate.md, docs/patterns/story-plan-and-execute.md]
updated: 2026-08-02
---

# Structural Patterns

Structural patterns describe how agents are connected and how they coordinate, independent of application domain. Each is expressed in catrace as a [[PDA Triplet Model]] over named state spaces. This page covers all 14 structural patterns.

## 1. Augmented LLM

**Topology:** Single node with external tools/memory/retrieval. The foundation for all higher patterns.

**catrace model:** Single-agent P, D, A over two experience states (`looks_routine`, `looks_risky`) and three actions (`answer`, `clarify`, `escalate`). The Q=DAP kernel captures how the agent's interpretation evolves from cycle to cycle.

**Key insight:** The Q entry `looks_routine → looks_risky` compresses many micro-paths: task was actually complex, agent answered directly, problem persisted, next cycle sees difficulty. See [[Example: Simple Agent]].

**Implemented:** `examples/simple_agent`

---

## 2. Prompt Chaining

**Topology:** Linear pipeline of N sequential agents, with optional gate between each.

**catrace model:** World states trace pipeline progress (`raw`, `extracted`, `summarised`, `formatted`, `failed`). Perception captures how reliably a stage agent reads incoming material quality (`input_clear`, `input_noisy`). Actions include `process`, `retry`, `escalate`. The world kernel W=PDA gives the full stage-to-stage transition matrix including retry loops and failure exits. Tracing onto {raw, formatted, failed} reveals end-to-end success/failure rates.

**Key insight:** The gate between stages is a soft probabilistic filter. A stage that perceives `input_noisy` but processes anyway has a high probability of producing output that the next stage will also perceive as noisy — cascading degradation is the primary failure mode.

**Not yet implemented** (issue #5).

---

## 3. Routing

**Topology:** Hub-and-spoke — router (classifier) + N specialised handlers.

**catrace model:** World states are true ticket types (`billing_ticket`, `technical_ticket`, `general_ticket`). Experience states are the router's perceived classification. Actions route to specialists or escalate to human triage. The world kernel W=PDA gives the ticket-flow transition matrix including misrouting loops. The stationary distribution over world states shows how much of the queue is occupied by tickets that have already been misrouted at least once.

**Key insight:** A confident wrong decision is worse than an uncertain escalation — misclassification creates a cycle where the wrong specialist cannot resolve the ticket and re-enters it into the queue as an apparently different type.

**Not yet implemented** (issue #6).

---

## 4. Parallelisation / Fan-out

**Topology:** Fork-join — N parallel workers + aggregator.

**catrace model:** Three independent content-moderation checkers (spam, hate, misinformation) each with their own P, D, A triplet operating on the joint product state space. Decision kernel is D_joint = D₁⊗D₂⊗D₃ (independent votes). The joint action kernel maps the 8 possible flag/pass combinations to a majority-vote moderation outcome (`removed`, `approved`). Entropy rate measures coordination efficiency: a well-tuned swarm has low entropy.

**Key insight:** Two checkers that fail on the same content type are much less useful than two that fail independently. The joint kernel captures error correlation that the individual detection rates cannot express.

**Not yet implemented** (issue #7).

---

## 5. Orchestrator-Workers

**Topology:** Centralised hub — orchestrator + N workers + synthesiser.

**catrace model:** Orchestrator world states track true task completion (`no_tasks_done`, `some_tasks_done`, `all_tasks_done`, `synthesis_done`). Worker world states are `idle`, `working`, `complete`, `failed`. The orchestrator's perception is coupled to aggregate worker world state: many failed workers drive the experience toward `progress_stalled` even if one thread quietly succeeds. The orchestrator's `retry_failed` action resets failed workers toward `idle`. Mean first passage time from `no_tasks_done` to `synthesis_done` measures expected delivery latency.

**Key insight:** The tension between waiting for complete results and synthesising early with incomplete information is directly measurable: compare MFPTs under different `synthesise` threshold policies.

**Not yet implemented** (issue #8).

---

## 6. Evaluator-Optimizer

**Topology:** Iterative loop — generator + evaluator/critic.

**catrace model:** See the specializations [[Dev-Workflow Patterns]] (D2, D3, D4) and the implemented `validator_repair` example ([[Example: Validator Repair]]). The core structure is a feedback loop where the evaluator's score drives the next generation cycle.

**Implemented (partial):** `examples/validator_repair`.

---

## 7. Autonomous Agent Loop

**Topology:** Single cyclic agent — perceive→decide→act→perceive.

**catrace model:** This is the native expression of the [[PDA Triplet Model]]. The Q=DAP kernel is the autonomous loop kernel. Every single-agent catrace example (simple_agent, self_healing_nodes) is an autonomous agent loop.

**Implemented:** `examples/simple_agent`, `examples/self_healing_nodes`.

---

## 8. Supervisor / Hierarchical

**Topology:** Tree — supervisor + sub-agents (possibly multi-level).

**catrace model:** Team world states track blocking count (`healthy`, `one_blocked`, `two_blocked`, `critical`). Supervisor experience is aggregate velocity (`velocity_good`, `velocity_degraded`). Supervisor actions are `hold_course`, `redistribute`, `escalate`. Developer agents each have `unblocked`/`blocked` states. Tracing onto {healthy, critical} collapses intermediate degradation states. First passage time from `healthy` to `critical` measures how long the team self-sustains without supervisor intervention.

**Key insight:** The primary question — how often is the supervisor's intervention actually necessary? — is answered by comparing the team's self-unblocking rate (in the developer D kernels) against the first-passage time from healthy to critical under `hold_course`.

**Not yet implemented** (issue #9).

---

## 9. Swarm / Peer-to-peer

**Topology:** Mesh / fully connected — N peer agents with shared world state.

**catrace model:** Zone world states track true aggregate coverage (`uncovered`, `partially_covered`, `well_covered`, `complete`). Each drone agent has local experience (`local_empty`, `local_partial`, `local_found`) and actions (`explore_new_area`, `reinforce_local`, `signal_found`). Decisions are independent within a cycle — agents coordinate only through the shared world state, not directly. The world kernel W=PDA gives coverage dynamics; the stationary distribution shows long-run coverage levels.

**Key insight:** Coverage collapse (all drones simultaneously targeting the same zone) is visible in the stationary distribution as high mass on partially_covered — the emergent failure mode that independent local decisions cannot self-correct.

**Not yet implemented** (issue #10).

---

## 10. Blackboard

**Topology:** Star — N specialist agents + shared blackboard workspace.

**catrace model:** Blackboard world states track collective agreement (`undiagnosed`, `tentative_diagnosis`, `confirmed_diagnosis`, `contradicted`). Each specialist's perception is coupled to the blackboard world state: a `tentative_diagnosis` shifts experience toward `evidence_strong` if the specialist's own data aligns, or toward `evidence_weak` if it conflicts. Decisions are independent (D_joint = D_rad⊗D_path⊗D_notes). The action kernel captures how combinations of specialist posts advance the board state. Mean first passage time from `undiagnosed` to `confirmed_diagnosis` measures diagnostic latency.

**Not yet implemented** (issue #11).

---

## 11. Debate / Adversarial

**Topology:** Pair/panel + judge.

**catrace model:** Debate world states track accumulated argument weight (`open`, `plaintiff_leading`, `defendant_leading`, `ruled_plaintiff`, `ruled_defendant`). Each debater's perception is coupled to the world state (leading side sees `position_strong`). The judge is a third agent whose perception (`case_clear`, `case_contested`) and action (`rule`, `call_for_more`) enter the joint composition alongside the debaters. Two simultaneous strong presses from opposing sides may cancel in the action kernel; a concession from one side shifts world weight toward the other.

**Not yet implemented** (issue #12).

---

## 12. Plan-and-Execute

**Topology:** Two-phase — planner agent + executor agent(s).

**catrace model:** World states track booking progress (`planning`, `executing_step_1`, `executing_step_2`, `executing_step_3`, `partial_failure`, `complete`, `abandoned`). The planner's perception is coupled to the executor's world state: `partial_failure` drives the planner toward `partial_context`, forcing a revision cycle. Commute time between `partial_failure` and `complete` measures the round-trip cost of a failed step.

**Key insight:** A rigid plan that cannot absorb partial failure costs more to revise than a looser one — the MFPT from `planning` to `complete` under a strict plan policy vs. a flexible one measures this directly.

**Not yet implemented** (issue #13).

---

## 14. Self-Healing / Adaptive

**Topology:** Coupled peer pair — operational agent + monitor/repair agent.

**catrace model:** The richest implemented pattern. See [[Example: Validator Repair]] (two-agent explicit coupling) and [[Example: Self-Healing Nodes]] (variant comparison: which loop carries the load?). Coupling enters through perception (monitor observes operational agent's world state) and action (monitor's repair restores operational agent's state).

**Key catrace insight — variant comparison:** Run the same joint kernel architecture under two parameter regimes; read the difference in stationary distribution, MFPT, and entropy rate to determine which mechanism carries the recovery load. See [[Variant Comparison Methodology]] for the general technique and hypothesis template.

**Implemented:** `examples/validator_repair`, `examples/self_healing_nodes`.

---

## Sources

- `docs/patterns/agentic-patterns-reference.md`
- `docs/patterns/story-single-llm-agent.md`
- `docs/patterns/story-hidden-support-system.md`
- `docs/patterns/story-validator-repair.md`
- `docs/patterns/story-self-healing-nodes.md`
- `docs/patterns/story-prompt-chaining.md`
- `docs/patterns/story-routing.md`
- `docs/patterns/story-parallelisation.md`
- `docs/patterns/story-orchestrator-workers.md`
- `docs/patterns/story-supervisor.md`
- `docs/patterns/story-swarm.md`
- `docs/patterns/story-blackboard.md`
- `docs/patterns/story-debate.md`
- `docs/patterns/story-plan-and-execute.md`
