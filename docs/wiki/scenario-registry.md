---
title: Scenario Registry
tags: [scenarios, readme, examples, status, registry]
sources: [README.md, plans/network-of-healers.md, examples/prompt_chaining/main.go]
updated: 2026-08-04
---

# Scenario Registry

Authoritative map of `README.md` scenario write-ups to story files, runnable
examples, and wiki pages. **README scenario numbering wins** over plans or
ad-hoc claims. Lint must re-check this page whenever `README.md`,
`examples/*/WALKTHROUGH.md`, or `plans/*.md` change.

## README scenarios

| # | README title | Story file | Example code | Wiki | Status |
|---|--------------|------------|--------------|------|--------|
| 1 | Single LLM task agent | `docs/patterns/story-single-llm-agent.md` | `examples/simple_agent` | [[Example: Simple Agent]] | Implemented |
| 2 | LLM agent with hidden support system | `docs/patterns/story-hidden-support-system.md` | `examples/trace_analysis` | [[Example: Hidden Support System]] | Implemented |
| 3 | Two-agent validator / repair pair | `docs/patterns/story-validator-repair.md` | `examples/validator_repair` | [[Example: Validator Repair]] | Implemented |
| 4 | Self-adjusting / self-healing network nodes | `docs/patterns/story-self-healing-nodes.md` | `examples/self_healing_nodes` | [[Example: Self-Healing Nodes]] | Implemented |
| 5 | Three-agent majority-valid coordination network | `docs/patterns/story-supervisor.md` | — | [[Structural Patterns]] (§8 Supervisor) | Not yet implemented |
| 6 | Four-agent pipeline with escalation | `docs/patterns/story-prompt-chaining.md` | draft only (see below) | [[Structural Patterns]] (§2 Prompt Chaining) | Not yet implemented |

## Planned examples not on the README list

| Plan | Intended path | Notes |
|------|---------------|-------|
| `plans/network-of-healers.md` | `examples/network_of_healers` | Sequel to [[Example: Self-Healing Nodes]] (topology / coupling variants). **Not** README scenario 5. Premise still under critical review; no scenario number until README is deliberately updated. |

## Draft code without walkthrough

| Path | Related README # | Notes |
|------|------------------|-------|
| `examples/prompt_chaining/main.go` | 6 | Partial draft exists; no `WALKTHROUGH.md`, not listed under README “Current files”. Treat as **not implemented** until walkthrough + README status flip. |

## Lint rules for this registry

1. Every `### N.` heading under README “Scenario write-ups” must appear in the table above.
2. Status “Implemented” requires `examples/<name>/main.go` **and** `WALKTHROUGH.md`.
3. A `plans/*.md` file must not claim a README scenario number that maps to a different story.
4. Pattern catalogue status ([[Agentic Patterns Catalogue]]) should not contradict this page for README-linked examples.

## Sources

- `README.md`
- `plans/network-of-healers.md`
- `examples/prompt_chaining/main.go`
