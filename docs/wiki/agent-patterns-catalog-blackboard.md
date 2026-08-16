---
title: Agent Patterns Catalog — Blackboard
tags: [patterns, blackboard, multi-agent, shared-workspace, external-catalog, agent-patterns-catalog]
sources: [docs/wiki/raw/agent-patterns-catalog-blackboard.md]
updated: 2026-08-05
---

# Agent Patterns Catalog — Blackboard

External definition of the **Blackboard** pattern (aliases: *Shared Workspace*, *Collaboration Whiteboard*) from the [Agent Patterns Catalog](https://www.agentpatternscatalog.org/patterns/blackboard/). Status in practice: experimental. Intent: give multiple specialised agents one inspectable, queryable store they read from and write to under structured keys, so coordination is “contribute what you can” without point-to-point knowledge of peers.

In catrace this maps to pattern **10** in [[Structural Patterns]] / [[Agentic Patterns Catalogue]]: a star topology of specialists plus a shared board world state ([[Example: Blackboard]], `docs/patterns/story-blackboard.md`). The catalog stresses **policy-driven conflict resolution**, optional change notification, and a hard forbid on out-of-band A2A — stronger implementation discipline than the Implemented catrace example, which models agreement dynamics (`undiagnosed` → `confirmed_diagnosis` / `contradicted`) and board-triggered perception without spelling out write races, pruning, or an optional controller that chooses who runs next.

Neighbourhood that matters for catrace: Blackboard is an **alternative-to** Supervisor (shared store vs routing agent), **complements** Swarm (indirect shared state vs direct peer mesh), and **generalises** Stigmergic Coordination (structured shared workspace vs environment marks). Multi-agent modelling still uses [[Joint Kernels and Coupling]]: independent D via Kronecker product; coupling only through perception of / action on the board.

## Key Points

- **Problem shape:** Isolation duplicates work; N×N messaging is brittle; unstructured shared memory races — need flexible order with disciplined shared state.
- **Solution:** Shared store (file / DB / in-memory); agents read a slice, write under keys; optional notify; conflict policy (last-write-wins, version-vector, append-only).
- **Forbids:** All cross-agent communication via the blackboard only — no out-of-band agent-to-agent calls.
- **Gives / costs:** Loose coupling + inspectable state vs races and bloat without pruning.
- **Relations:** generalises Stigmergic Coordination; complements Swarm; alternative-to Supervisor (also Topic-Based Routing, Cellular-Automata Agents, DCOP); used-by SOP-Encoded Multi-Agent Workflow.
- **Frameworks cited:** AgentVerse (shared environment), CopilotKit Shared State, Dendron behavior-tree blackboard.
- **catrace status:** Implemented — `examples/blackboard/` ([[Example: Blackboard]], README scenario 7, issue #11).

## Sources

- `docs/wiki/raw/agent-patterns-catalog-blackboard.md`
- https://www.agentpatternscatalog.org/patterns/blackboard/
- https://www.agentpatternscatalog.org/landing/patterns/blackboard/
