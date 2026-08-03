---
title: "Example: Prompt-Chaining Pipeline"
tags: [example, prompt-chaining, pipeline, trace, stationary, single-agent]
sources: [examples/prompt_chaining/WALKTHROUGH.md, docs/patterns/story-prompt-chaining.md]
updated: 2026-08-05
---

# Example: Prompt-Chaining Pipeline

**Pattern:** Prompt Chaining (2) — see [[Agentic Patterns Catalogue]]
**Code:** `examples/prompt_chaining/main.go`
**Run:** `go run examples/prompt_chaining/main.go`

This example extends the single-agent model to five world states, demonstrating
how a multi-stage pipeline with retry loops and a failure exit is captured by
the world kernel W = PDA. It is the first catrace example where the world state
space is large enough that the stationary distribution is non-obvious by
inspection.

## The scenario

A document intelligence pipeline sends text through three specialist stages —
extractor, summariser, formatter — with a quality gate between each. A stage
agent perceives input quality and decides to process, retry, or escalate. The
full pipeline dynamics, including retry loops and failure exits, are encoded
in a single 5×5 world kernel.

## State spaces

| Layer      | States                                                      | Meaning                               |
|------------|-------------------------------------------------------------|---------------------------------------|
| World      | `raw`, `extracted`, `summarised`, `formatted`, `failed`     | True pipeline stage of the document   |
| Experience | `input_clear`, `input_noisy`                                | Stage agent's read of input quality   |
| Action     | `process`, `retry`, `escalate`                              | Transform, retry, or abandon          |

## Key results

**Stationary distribution (W):**

```
raw        0.123115
extracted  0.296395
summarised 0.266709
formatted  0.165227
failed     0.148554
```

The pipeline spends most time in intermediate stages (`extracted` + `summarised`
= 56.3%), with `formatted` at 16.5% and `failed` at 14.9%. No stage dominates,
reflecting a genuinely mixed pipeline with retry pressure.

**Entropy rate:** 2.221 bits/step — high variability, the pipeline is hard to
predict step-to-step.

**Communicating classes:** single recurrent class containing all 5 states.
Every pipeline stage is reachable from every other, including recovery from
`failed`. This is confirmed by `Classes()` returning `[[0 1 2 3 4]]` with no
transient states.

## Trace onto {raw, formatted, failed}

Collapsing `extracted` and `summarised` into hidden states gives the coarse
end-to-end pipeline picture:

```
IsTraceOf = true
stationary(trace) = 0.281795  0.378183  0.340022
                    raw        formatted  failed
```

The pipeline reaches `formatted` 37.8% of the time and `failed` 34.0% of the
time at the coarse checkpoint level. `IsTraceOf = true` confirms the 3-state
kernel is exact — it captures all excursions through `extracted` and
`summarised` without approximation. See [[Trace Chain]].

## Modeling notes

The A kernel (G→W) is memoryless: it does not know the current pipeline stage.
`A[process, :]` encodes the average advancement effect of processing across
all stages. This is a necessary simplification of the linear pipeline structure
into the Markov framework: the stage-to-stage progression emerges from the
interaction of P, D, and A rather than being directly specified. The resulting
W matrix correctly shows advancing transitions (extracted→summarised is the
modal next step from extracted) despite A being stage-agnostic.

## Connection to math and API

- `Agent.WorldKernel()` → [[catrace API]], [[PDA Triplet Model]]
- `Kernel.Trace()` and `IsTraceOf()` → [[Trace Chain]]
- `Kernel.Stationary()`, `EntropyRate()`, `Classes()` → [[Markov Chain Foundations]]

## Sources

- `examples/prompt_chaining/WALKTHROUGH.md`
- `docs/patterns/story-prompt-chaining.md`
