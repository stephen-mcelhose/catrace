---
title: "Example: Blackboard"
tags: [example, blackboard, shared-workspace, multi-agent, trace, mfpt, opportunistic]
sources: [examples/blackboard/WALKTHROUGH.md, docs/patterns/story-blackboard.md, docs/wiki/raw/agent-patterns-catalog-blackboard.md]
updated: 2026-08-05
---

# Example: Blackboard

**Pattern:** Blackboard (10) — see [[Agentic Patterns Catalogue]], [[Agent Patterns Catalog — Blackboard]]
**Code:** `examples/blackboard/main.go`
**Run:** `go run examples/blackboard/main.go`

Three identical specialists collaborate through a shared diagnostic board with no peer messaging. Board status is what they perceive and what their joint posts rewrite. Unlike [[Example: Validator Repair]] (cross-agent health coupling) or [[Example: Prompt Chaining]] (one active stage), every specialist may contribute each cycle; coordination is opportunistic via shared state.

## Board world

| State | Meaning |
|-------|---------|
| `undiagnosed` | Empty case board |
| `tentative_diagnosis` | Working hypothesis posted |
| `confirmed_diagnosis` | Closed / absorbing |
| `contradicted` | Irreconcilable findings (may reopen) |

## Construction

$$
J[w,w'] = \sum_{x,g} P_\text{joint}[w,x]\, D^{\otimes 3}[x,g]\, A[w,g,w']
$$

$P$ raises $P(\texttt{evidence_strong})$ when the board is tentative. $A$ counts posts / endorsements / flags. Because $A$ depends on current $w$, $J$ is assembled row-wise.

## Analysis outputs (default parameters)

| Quantity | Value | Reading |
|----------|-------|---------|
| MFPT(`undiagnosed`→`confirmed`) | ≈ 10.41 steps | Diagnostic latency |
| π(`confirmed`) | 1.0 | Only recurrent class once confirmed absorbs and contradicted reopens |
| Trace `{undiagnosed, confirmed, contradicted}` | IsTraceOf true | Coarse start / confirmed / contradicted picture |

## Related

- Story: `docs/patterns/story-blackboard.md`
- Walkthrough: `examples/blackboard/WALKTHROUGH.md`
- [[Scenario Registry]] scenario 7
- [[Structural Patterns]] §10
- [[Joint Kernels and Coupling]], [[Trace Chain]], [[catrace API]]

## Sources

- `examples/blackboard/WALKTHROUGH.md`
- `docs/patterns/story-blackboard.md`
- `docs/wiki/raw/agent-patterns-catalog-blackboard.md`
