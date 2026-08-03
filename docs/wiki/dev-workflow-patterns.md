---
title: Dev-Workflow Patterns
tags: [patterns, dev-workflow, coding-agent, implement, verify, critic, plan]
sources: [docs/patterns/agentic-patterns-reference.md, docs/patterns/story-research-plan-implement.md, docs/patterns/story-implement-verify.md, docs/patterns/story-implement-critic.md, docs/patterns/story-plan-implement-critic-verify.md]
updated: 2026-08-02
---

# Dev-Workflow Patterns

These four patterns are specializations of [[Structural Patterns]] for AI coding and development agent loops. Each inherits from one or more structural parents but has a specific semantic meaning in the development context.

All four are planned catrace examples (not yet implemented). They are most naturally analyzed via [[Markov Chain Foundations]] metrics: MFPT from start to completion, stationary distribution over review/revision states, and entropy rate as a measure of predictability. Any two sub-patterns that share a topology (e.g., D2 vs D3) are natural candidates for [[Variant Comparison Methodology]].

---

## D1. Research-Plan-Implement

**Parent patterns:** Prompt Chaining (2), Plan-and-Execute (12)

**Key distinction:** Three compounding stages where each gate amplifies or corrects the previous stage's quality. Exploration quality gates plan quality; plan quality gates implementation success.

**catrace model:**

| Agent | World states | Experience states | Actions |
|-------|-------------|-------------------|---------|
| Researcher | `unexplored`, `explored`, `planned`, `implementing`, `done`, `blocked` | `context_rich`, `context_sparse` | `explore_more`, `commit_to_plan` |
| Planner | (same world) | `plan_coherent`, `plan_conflicts_detected` | `produce_plan`, `revise_plan`, `abandon` |
| Implementer | (same world) | `step_clear`, `step_ambiguous` | `implement_step`, `backtrack`, `request_replan` |

The research-to-plan transition is the first gate: proceeding with `context_sparse` raises the probability of `plan_conflicts_detected`, which forces a costly revision cycle — an effect not visible in the individual agent kernels, only in the composed world kernel.

**Key metric:** MFPT from `unexplored` to `done` as a function of exploration budget. This directly answers: does spending more cycles on research reduce total task latency?

**Issue:** [#14](https://github.com/stephen-mcelhose/catrace/issues/14)

---

## D2. Implement-Verify

**Parent pattern:** Evaluator-Optimizer (6)

**Key distinction from parent:** The evaluator is the environment (objective tool — test suite, type checker), not an agent with its own perception loop. No perception noise on the evaluator side.

**catrace model:**

| Element | States |
|---------|--------|
| World states | `all_failing`, `some_failing`, `all_passing`, `regressed` |
| Implementer experience | `output_bad`, `output_partial`, `output_clean` |
| Implementer actions | `targeted_fix`, `broad_refactor`, `revert_and_retry` |

The `regressed` state is the key catrace insight: it is reachable from `all_passing`. A check suite that passes today may fail tomorrow if the fix was fragile. The stationary distribution reveals how much time the agent spends stuck in regression cycles vs. making progress.

**Contrast with D3:** In D2 the verifier is a deterministic tool (its accuracy is 1); in D3 the critic is an agent with its own imperfect perception, making joint dynamics richer.

**Issue:** [#15](https://github.com/stephen-mcelhose/catrace/issues/15)

---

## D3. Implement-Critic

**Parent pattern:** Evaluator-Optimizer (6)

**Key distinction from parent:** The evaluator is itself an agent with its own perception kernel and policy. The critic's accuracy is a stochastic variable that shapes long-run outcomes.

**catrace model:**

| Agent | States (experience) | Actions |
|-------|--------------------|----|
| Implementer | `requirements_clear`, `requirements_ambiguous` | `implement`, `revise_targeted`, `revise_broad`, `push_back` |
| Critic | `code_sound`, `code_has_issues`, `code_unclear` | `approve`, `request_changes`, `escalate_to_human` |

Decisions are independent within a cycle: D_joint = D_impl⊗D_critic; the implementer submits before knowing the critic's verdict. The critic's perception is coupled to the true world state (draft quality), but with noise — a subtly broken implementation may look fine; a non-idiomatic but correct one may appear broken.

**Key metrics:**
- Stationary distribution over `revision_requested` vs `approved` — a high `revision_requested` mass indicates an over-zealous critic or poor initial implementation
- Commute time between `draft` and `approved` — the average round-trip cost of one full review cycle

**Issue:** [#16](https://github.com/stephen-mcelhose/catrace/issues/16)

---

## D4. Plan-Implement-Critic-Verify

**Parent patterns:** Plan-and-Execute (12), Evaluator-Optimizer (6)

**Key distinction:** Two independent evaluators with different failure modes — objective verifier (cheap, deterministic) + subjective critic (expensive, noisy). Running them in the right order — verify first, then critique — avoids wasting the critic's capacity on code that won't type-check.

**catrace model:**

World states track the full development cycle step: `planning`, `implementing`, `awaiting_critique`, `awaiting_verification`, `step_done`, `plan_revision`, `abandoned`.

Four agents:
- **Planner**: `requirements_clear` / `requirements_partial` → `produce_plan`, `revise_plan`, `abandon`
- **Implementer**: `plan_step_clear` / `plan_step_ambiguous` → `implement`, `request_clarification`, `skip_step`
- **Critic**: `code_sound` / `code_has_issues` → `approve`, `request_changes`, `escalate_to_planner`
- **Verifier**: `checks_passing` / `checks_failing` → `pass`, `fail` (effectively deterministic — models automated tooling)

The action kernel encodes routing logic:
- `fail` from verifier + any critic action → back to `implementing` (tight loop)
- `escalate_to_planner` from critic + `fail` from verifier → `plan_revision` (wide loop)

**Key insight:** The depth of the feedback loop triggered by each failure type is the core modelling insight. Shallow loops (verify → implement) are fast but local; deep loops (critique → plan revision) are slow but can fix structural problems that shallow loops miss. MFPT from `planning` to `step_done` measures per-step latency across all policy variants.

**Issue:** [#17](https://github.com/stephen-mcelhose/catrace/issues/17)

---

## Comparison table

| Pattern | Evaluator type | Evaluator noise | Loop depth | Key catrace metric |
|---------|---------------|-----------------|------------|--------------------|
| D2 | Environment (tool) | None | Single | MFPT to `all_passing`; stationary on `regressed` |
| D3 | Agent (critic) | High | Single | Commute time draft↔approved; stationary on `revision_requested` |
| D4 | Both (tool + agent) | Tool=0, critic=high | Multi-depth | MFPT per step; stationary across all routing states |

## Sources

- `docs/patterns/agentic-patterns-reference.md`
- `docs/patterns/story-research-plan-implement.md`
- `docs/patterns/story-implement-verify.md`
- `docs/patterns/story-implement-critic.md`
- `docs/patterns/story-plan-implement-critic-verify.md`
