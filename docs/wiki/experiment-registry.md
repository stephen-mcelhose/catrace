---
title: Experiment Registry
tags: [experiments, variant-comparison, hypothesis, results, registry]
sources: [experiments/README.md]
updated: 2026-08-05
---

# Experiment Registry

The `experiments/` directory contains formal hypothesis records for architectural claims tested with catrace. An experiment is distinct from an example: examples demonstrate usage; experiments ask and answer a specific design question using catrace as a measuring instrument.

See [[Variant Comparison Methodology]] for the full filing process and hypothesis template.

## Registry

Every **Pending** experiment must have an open GitHub issue (`experiment: …`).
Maintain via `.agents/skills/experiments/SKILL.md`.

| Slug | Claim | Status | Verdict | Issue |
|------|-------|--------|---------|-------|
| `nodes-throttle-vs-evolver` | Node throttle is primary recovery mechanism; evolver contributes but is not essential | Complete | Supported (4/4 metrics) | — |
| `wiki-knowledge-graph` | Trace chain corrects importance distortion caused by missing wiki pages | Complete | Not supported (1/4 criteria) | #21 |
| `kg-grounding-agent-behavior` | Knowledge graph quality systematically shifts knowledge agent behavior — richer graph → higher π(understood), faster recovery, lower entropy | Pending | — | #32 |
| `spectral-gap-mixing-time` | Spectral gap rank-orders kernels correctly by empirical mixing speed | Pending | — | #25 |
| `stationary-sensitivity` | Perturbations to high-π rows cause disproportionately large shifts in stationary distribution | Pending | — | #26 |
| `n-agent-scalability` | Dense joint kernel approach becomes intractable at N=4+ agents; trace collapses are the solution | Pending | — | #27 |
| `heal-on-critical-path` | On a fixed \(1\to2\) load graph, strong local heal belongs on the downstream sink more than on the upstream feeder | Pending | — (needs `network_of_healers`) | #33 |
| `collapse-masks-heterogeneity` | Strong-sink / weak-feeder makes collapsed pool_ok look healthier than joint upstream-overload mass warrants | Pending | — (needs `network_of_healers`) | #34 |

## Completed results

### nodes-throttle-vs-evolver — Supported (4/4 metrics)

The node's local throttle is the primary recovery mechanism. Variant A (strong throttle) outperformed Variant B (strong evolver) on all four metrics: π(H·G), MFPT(O·B→H·G), H(J), and π_trace(H·G). The outer evolutionary loop is a safety net, not a fast path. See [[Example: Self-Healing Nodes]] for the full analysis.

### wiki-knowledge-graph — Not supported (1/4 criteria)

The hypothesis was that trace-chain correction of a 28-node wiki graph (14 existing + 14 planned pages as hidden states) would substantially rerank structurally important pages (`Structural Patterns`, `Dev-Workflow Patterns`) relative to a naive 14-node PageRank. Only 1 of 4 criteria was met. The mathematical identity (`IsTraceOf = true`) held, but the predicted rank shifts did not materialise at the required magnitudes.

This does not invalidate the [[Trace Chain]] technique — the stationary consistency theorem is exact and confirmed. What failed was the specific architectural claim that trace correction on this particular incomplete graph would produce large, predictable reranking. The graph's actual link structure produced smaller shifts than predicted.

## Filing a new experiment

1. Create `experiments/<slug>/`
2. Copy `experiments/hypothesis-template.md` → `hypothesis.md`
3. Fill Claim, Variables, Predictions **before running**
4. Run catrace analysis
5. Fill Results and Verdict
6. Add a row to `experiments/README.md` (with Issue `#N`)
7. Open/sync GitHub issue — prefer `.agents/skills/experiments` **file** / **lint**

## Sources

- `experiments/README.md`
