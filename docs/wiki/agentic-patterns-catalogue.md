---
title: Agentic Patterns Catalogue
tags: [patterns, agentic, structural, topology, dev-workflow, catalogue]
sources: [docs/patterns/agentic-patterns-reference.md, docs/wiki/raw/agent-patterns-catalog-blackboard.md]
updated: 2026-08-05
---

# Agentic Patterns Catalogue

catrace provides Markov chain models for 14 structural agentic patterns and 4 development-workflow sub-patterns. Each pattern has a story file (state spaces, semantic interpretation, played-out micro-paths) and, for implemented patterns, a runnable example with WALKTHROUGH.

See [[Structural Patterns]] and [[Dev-Workflow Patterns]] for depth. This page is the overview index. For README scenario numbers (1–7) and plan-vs-README integrity, see [[Scenario Registry]].

## Structural patterns

These describe *how agents are connected and coordinate*, independent of domain.

| # | Pattern | Topology | catrace example | Story | Status |
|---|---------|----------|-----------------|-------|--------|
| 1 | Augmented LLM | Single node | `simple_agent` | [[Example: Simple Agent]] | ✅ Implemented |
| 2 | Prompt Chaining | Linear pipeline | `prompt_chaining` | [[Example: Prompt Chaining]] | ✅ Implemented |
| 3 | Routing | Hub-and-spoke | — | [[Structural Patterns]] | 🔲 Planned |
| 4 | Parallelisation / Fan-out | Fork-join | — | [[Structural Patterns]] | 🔲 Planned |
| 5 | Orchestrator-Workers | Centralised hub | — | [[Structural Patterns]] | 🔲 Planned |
| 6 | Evaluator-Optimizer | Iterative loop | `validator_repair` (partial) | [[Example: Validator Repair]] | ✅ Implemented |
| 7 | Autonomous Agent Loop | Single cyclic agent | `simple_agent`, `self_healing_nodes` | [[Example: Simple Agent]], [[Example: Self-Healing Nodes]] | ✅ Implemented |
| 8 | Supervisor / Hierarchical | Tree | — (README scenario 5) | [[Structural Patterns]] | 🔲 Planned |
| 9 | Swarm / Peer-to-peer | Mesh | — | [[Structural Patterns]] | 🔲 Planned |
| 10 | Blackboard (Shared Workspace / Collaboration Whiteboard) | Star (shared memory) | `blackboard` | [[Example: Blackboard]]; [[Agent Patterns Catalog — Blackboard]] | ✅ Implemented |
| 11 | Debate / Adversarial | Pair/panel + judge | — | [[Structural Patterns]] | 🔲 Planned |
| 12 | Plan-and-Execute | Two-phase | — | [[Structural Patterns]] | 🔲 Planned |
| 13 | Human-in-the-Loop | Any + human gate | — | — | 🔲 Planned |
| 14 | Self-Healing / Adaptive | Coupled peer pair | `validator_repair`, `self_healing_nodes` | [[Example: Validator Repair]], [[Example: Self-Healing Nodes]] | ✅ Implemented |

## Dev-workflow sub-patterns

Specializations of structural patterns for AI coding/development agent loops.

| # | Pattern | Parent(s) | catrace example | Status |
|---|---------|-----------|-----------------|--------|
| D1 | Research-Plan-Implement | Prompt Chaining (2), Plan-and-Execute (12) | — | 🔲 Planned |
| D2 | Implement-Verify | Evaluator-Optimizer (6) | — | 🔲 Planned |
| D3 | Implement-Critic | Evaluator-Optimizer (6) | — | 🔲 Planned |
| D4 | Plan-Implement-Critic-Verify | Plan-and-Execute (12), Evaluator-Optimizer (6) | — | 🔲 Planned |

See [[Dev-Workflow Patterns]] for the key distinction of each sub-pattern.

## Sources (external)

The pattern catalogue draws from three primary sources:
- **ANT** — Anthropic, *Building Effective Agents*, 2024
- **IBM** — IBM Think, *AI Agent Use Cases*
- **LG** — VThink Technologies, *Common Agentic Patterns in LangGraph*

Additional catalog ingest (Blackboard only so far):
- **APC** — [Agent Patterns Catalog — Blackboard](https://www.agentpatternscatalog.org/patterns/blackboard/) (aliases Shared Workspace / Collaboration Whiteboard; forbids out-of-band A2A; alternative-to Supervisor; complements Swarm; generalises Stigmergic Coordination). Wiki page: [[Agent Patterns Catalog — Blackboard]]. catrace status: ✅ Implemented — `examples/blackboard/` ([[Example: Blackboard]], README scenario 7).

## catrace modeling approach

Every pattern is expressed as a [[PDA Triplet Model]] over appropriate state spaces. The catrace methodology:
1. Name the world, experience, and action states with semantic meaning
2. Specify P, D, A kernels (single agent) or P_joint, D_joint, A_joint (multi-agent)
3. Derive the world kernel W = P·D·A (or joint J) and analyze with [[catrace API]] tools
4. Apply [[Trace Chain]] to collapse unobservable or irrelevant states

For multi-agent patterns, decisions are always independent (Kronecker product D_joint = D₁⊗D₂); coupling enters only in perception (P_joint) and action effects (A_joint). See [[Joint Kernels and Coupling]].

When comparing two design choices that share the same topology, use the [[Variant Comparison Methodology]]: encode each choice as a parameter regime and let the Markov chain metrics vote.

## Sources

- `docs/patterns/agentic-patterns-reference.md`
- `docs/wiki/raw/agent-patterns-catalog-blackboard.md`
- [[Agent Patterns Catalog — Blackboard]]
