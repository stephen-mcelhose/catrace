---
title: Variant Comparison Methodology
tags: [methodology, variant-comparison, hypothesis, experiment, measurement, architecture]
sources: [examples/self_healing_nodes/WALKTHROUGH.md, docs/patterns/story-self-healing-nodes.md]
updated: 2026-08-02
---

# Variant Comparison Methodology

Most tools for evaluating agentic architectures work by running the system and reasoning qualitatively about what happened. catrace offers a complementary approach: **encode a design debate as two kernel parameter regimes and let the Markov chain metrics vote**.

This is catrace's core methodological contribution beyond implementing the math.

## The core idea

An architectural claim like "mechanism X recovers faster than mechanism Y" is hard to evaluate by inspection — X and Y may differ in subtle ways whose interactions are non-obvious. catrace makes the comparison precise:

1. **Fix the topology** — same state spaces, same agents, same coupling structure
2. **Vary only the parameters that encode the claim** — specific entries in P, D, or A kernels
3. **Run both variants through the same analysis pipeline** — `Stationary`, `MeanFirstPassage`, `EntropyRate`, `Trace`
4. **Read the verdict from the numbers** — the metrics either agree or reveal a trade-off

The claim is supported if the metrics vote consistently. A split vote reveals a genuine trade-off: X wins on recovery speed, Y wins on predictability. That is also information.

## Hypothesis template

Every catrace experiment should be stated before running. See `experiments/hypothesis-template.md` for the blank form.

**Claim** — the architectural assertion in plain language. Must be falsifiable: it should be possible for the numbers to go the other way.

**Variables** — which specific kernel entries differ between variants, and why those entries encode the claim. If you can't point to specific matrix cells, the claim is not yet well-formed.

**Metrics and predictions** — for each metric, state which direction supports the claim *before* running. This prevents post-hoc rationalization.

| Metric | Use when you care about... |
|--------|---------------------------|
| π(target state) | Long-run time in a desired state |
| MFPT(bad → good) | Recovery speed from degraded states |
| H(J) | Predictability / legibility of the system |
| Trace π(coarse state) | Observable health from an external vantage point |

**Verdict rule** — specify in advance how many metrics must agree to call the claim supported. Recommend: majority (≥2 of 3) for claims with three metrics; unanimous for two-metric comparisons where the metrics are closely related.

## What makes a well-formed hypothesis

A good hypothesis:
- Names a single mechanism as primary (not "X and Y both matter")
- Identifies exactly which kernel entries encode that mechanism
- Predicts the metric direction before running
- Uses a topology that isolates the mechanism (other coupling held equal)

A poorly-formed hypothesis:
- Is not translatable to kernel entries ("the system is more robust")
- Uses a different topology for each variant (topology changes confound parameter changes)
- Has no metric that could distinguish the variants
- Is stated after seeing the results

## Completed experiments

### nodes-throttle-vs-evolver — Supported (4/4 metrics)

See [[Example: Self-Healing Nodes]] for the full analysis and `experiments/nodes-throttle-vs-evolver/hypothesis.md` for the formalized record.

**Claim:** The node's throttle action is the primary recovery mechanism; the evolver contributes but is not essential.

**Result:** All four metrics (π, MFPT, H, π_trace) supported the claim. Variant A (strong throttle) produced a healthier, faster-recovering, more legible system than Variant B (strong evolver, weak throttle).

**What the split would have looked like:** If π(H·G) favored A but MFPT favored B, the correct interpretation would be: "throttle keeps the system healthier on average, but when things go badly wrong, random mutation search recovers faster from the worst state." That is a genuine architectural trade-off — not a failure of the methodology.

### wiki-knowledge-graph — Not supported (1/4 criteria)

See `experiments/wiki-knowledge-graph/hypothesis.md` for the formalized record.

**Claim:** Trace chain correction of the catrace wiki graph (treating 14 planned pages as hidden states) would substantially rerank structurally important pages relative to a naive 14-node PageRank.

**Result:** Only 1 of 4 structural criteria was met. The mathematical identity (`IsTraceOf = true`) held — the stationary consistency theorem is exact. But the predicted rank shifts for `Structural Patterns` and `Dev-Workflow Patterns` did not materialise at the required magnitudes.

**What this means for the methodology:** A "not supported" verdict is not a failure — it is information. The specific graph structure of this incomplete wiki produced smaller correction effects than predicted. The technique is sound; the claim about this particular graph was wrong. A revised hypothesis might lower the required rank-shift thresholds, or model the planned pages with different assumed link distributions.

## Generalizing to other patterns

Any pattern where two design choices share a topology is a candidate:

| Pattern | Candidate comparison |
|---------|---------------------|
| Evaluator-Optimizer (6) | Strict critic vs. lenient critic — which produces better stationary distribution? |
| Orchestrator-Workers (5) | Early synthesis vs. wait-for-all — which minimizes MFPT to completion? |
| Routing (3) | Confident classifier vs. uncertainty-aware escalator — which minimizes misrouting residence time? |
| D3 Implement-Critic | High-sensitivity critic vs. high-specificity critic — which reduces revision cycle commute time? |

Each comparison follows the same four steps: fix topology → vary parameters → run pipeline → read verdict.

## Relationship to catrace examples

Implemented examples are not variant comparisons by default — they demonstrate a single configuration. A variant comparison requires two runs of the same example under different kernel parameters, with the comparison stated as a hypothesis in advance.

The `experiments/` directory is where these comparisons live, separate from the example walkthroughs.

## Pending experiments

Four further experiments are filed and awaiting runs:

| Slug | Claim | Issue |
|------|-------|-------|
| `kg-grounding-agent-behavior` | Knowledge graph quality systematically shifts agent behavior — a well-structured graph raises π(understood), lowers MFPT to recovery, and reduces entropy rate; only the Perception kernel P changes between variants | — |
| `spectral-gap-mixing-time` | Spectral gap (1 − \|λ₂\|) rank-orders kernels correctly by empirical mixing speed; a kernel with 2× the spectral gap mixes in roughly half the steps | #25 |
| `stationary-sensitivity` | Perturbations to high-π rows cause disproportionately large shifts in stationary distribution; low-π rows have negligible effect | #26 |
| `n-agent-scalability` | The trace chain restores tractability for N-agent networks: a joint kernel built over k^N states can be collapsed to a small observable subset without losing stationary consistency | — |

These experiments extend the methodology into **perception quality** (kg-grounding), **measurement** (spectral gap, sensitivity), and **scale** (N-agent trace collapse). The hypothesis files are in `experiments/<slug>/hypothesis.md`.

## Sources

- `examples/self_healing_nodes/WALKTHROUGH.md`
- `docs/patterns/story-self-healing-nodes.md`
