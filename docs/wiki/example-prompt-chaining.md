---
title: "Example: Prompt Chaining"
tags: [example, prompt-chaining, pipeline, gate, trace, mfpt, workflow]
sources: [examples/prompt_chaining/WALKTHROUGH.md, docs/patterns/story-prompt-chaining.md]
updated: 2026-08-05
---

# Example: Prompt Chaining

**Pattern:** Prompt Chaining (2) — see [[Agentic Patterns Catalogue]]
**Code:** `examples/prompt_chaining/main.go`
**Run:** `go run examples/prompt_chaining/main.go`

Anthropic-style prompt chaining as a catrace example: three specialised LLM stages (extract → summarise → format) with programmatic gates, assembled into a pipeline world kernel. Unlike [[Example: Simple Agent]] (one PDA) or [[Example: Validator Repair]] (peers acting each tick), only the active stage fires; the gate is code, not an agent.

## Pipeline world

| State | Meaning |
|-------|---------|
| `raw` | Filing on the desk |
| `extracted` | Claims accepted by first gate |
| `summarised` | Brief accepted by second gate |
| `formatted` | Report shipped (absorbing) |
| `failed` | Human queue (may re-queue to `raw`) |

## Construction

Each stage contributes a row of $W$ when active:

$$
W[s,\cdot] = \sum_{x,g} P_s(x)\, D_s(g\mid x)\, \mathrm{Gate}(\cdot\mid g, s)
$$

`pass` advances to the next pipeline state; `retry_stage` / retry actions keep mass at $s$; `escalate` moves to `failed`.

## Analysis outputs (default parameters)

| Quantity | Value | Reading |
|----------|-------|---------|
| MFPT(`raw`→`formatted`) | ≈ 7.68 steps | Expected desk latency under retry/escalate |
| π(`formatted`) | 1.0 | Only recurrent class once shipping is absorbing |
| Trace `{raw, formatted, failed}` | IsTraceOf true | Coarse start / shipped / human-queue picture |

## Related

- Story: `docs/patterns/story-prompt-chaining.md`
- Walkthrough: `examples/prompt_chaining/WALKTHROUGH.md`
- [[Scenario Registry]] scenario 6
- [[Structural Patterns]] §2
- [[Trace Chain]], [[catrace API]] (`Trace`, `MeanFirstPassage`, `Stationary`)

## Sources

- `examples/prompt_chaining/WALKTHROUGH.md`
- `docs/patterns/story-prompt-chaining.md`
