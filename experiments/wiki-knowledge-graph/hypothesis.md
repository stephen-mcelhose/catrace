# Hypothesis: Trace chain corrects importance distortion from missing wiki pages

## Claim

The naive PageRank on the current 14-page wiki systematically undervalues catalogue/index pages
(`structural-patterns`, `dev-workflow-patterns`) and overvalues pages that happen to be the sole
target of a low-fanout page (`example-self-healing-nodes`, boosted by
`variant-comparison-methodology` having only one out-edge). These distortions are artifacts of
incompleteness — the planned example pages do not yet exist, so their inbound and outbound link
flow is missing entirely.

The trace chain, applied to the full graph (14 existing + 14 planned pages) with the planned pages
as the hidden subset, will produce corrected importance scores for the 14 existing pages that
better reflect their structural role in the complete wiki.

## Context

- **Pattern:** Knowledge graph as Markov chain; trace chain for latent-node correction
- **Related example:** `examples/trace_analysis` (same mathematical technique, different domain)
- **Motivation:** The mapping is described in `docs/wiki/variant-comparison-methodology.md` and
  the broader design discussion. Catrace's trace chain is the principled tool for exactly this
  situation: a subset of nodes is observable; the rest are latent but their influence is real.

## The two graphs

### Graph A — existing wiki only (14 nodes, naive PageRank baseline)

Nodes and current naive PageRank scores (binary adjacency, row-stochastic, power iteration):

| Rank | Score  | Page |
|------|--------|------|
| 1    | 0.1477 | Example: Self-Healing Nodes |
| 2    | 0.1162 | Markov Chain Foundations |
| 3    | 0.1099 | PDA Triplet Model |
| 4    | 0.1056 | Joint Kernels and Coupling |
| 5    | 0.0947 | Trace Chain |
| 6    | 0.0747 | Agentic Patterns Catalogue |
| 7    | 0.0698 | Example: Validator Repair |
| 8    | 0.0684 | catrace API |
| 9    | 0.0659 | Example: Simple Agent |
| 10   | 0.0489 | Variant Comparison Methodology |
| 11   | 0.0400 | Dev-Workflow Patterns |
| 12   | 0.0237 | Example: Hidden Support System |
| 13   | 0.0208 | Structural Patterns |
| 14   | 0.0137 | catrace Glossary |

Known distortions in this baseline:
- `Structural Patterns` ranks 13th: no inbound links from the 10 planned example pages it
  catalogues. In the full wiki it would have the highest in-degree of any page.
- `Example: Self-Healing Nodes` ranks 1st partly because `Variant Comparison Methodology` has
  exactly one out-edge (to it). This is a concentration artifact, not genuine structural centrality.
- `catrace Glossary` ranks 14th: only `catrace API` links to it. Planned example pages would
  typically include a glossary reference.

### Graph B — full intended wiki (14 existing + 14 planned = 28 nodes)

**Planned nodes (hidden subset):**

These are the pages referenced or implied by existing pages that do not yet exist:

| Slug | Referenced from | Hypothesized out-edges (to existing pages) |
|------|----------------|------------------------------------------|
| `example-prompt-chaining` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, trace-chain, agentic-patterns-catalogue |
| `example-routing` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, agentic-patterns-catalogue |
| `example-parallelisation` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, joint-kernels-and-coupling, agentic-patterns-catalogue |
| `example-orchestrator-workers` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, agentic-patterns-catalogue |
| `example-supervisor` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, agentic-patterns-catalogue |
| `example-swarm` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, joint-kernels-and-coupling, agentic-patterns-catalogue |
| `example-blackboard` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, joint-kernels-and-coupling, agentic-patterns-catalogue |
| `example-debate` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, joint-kernels-and-coupling, agentic-patterns-catalogue |
| `example-plan-and-execute` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, trace-chain, agentic-patterns-catalogue |
| `example-human-in-the-loop` | structural-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, agentic-patterns-catalogue |
| `example-d1-research-plan-implement` | dev-workflow-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, trace-chain, agentic-patterns-catalogue, variant-comparison-methodology |
| `example-d2-implement-verify` | dev-workflow-patterns | markov-chain-foundations, catrace-api, agentic-patterns-catalogue, variant-comparison-methodology |
| `example-d3-implement-critic` | dev-workflow-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, joint-kernels-and-coupling, agentic-patterns-catalogue, variant-comparison-methodology |
| `example-d4-plan-implement-critic-verify` | dev-workflow-patterns | markov-chain-foundations, pda-triplet-model, catrace-api, joint-kernels-and-coupling, agentic-patterns-catalogue, variant-comparison-methodology |

**Inbound edges from existing pages to planned pages** (edges that need to be added to graph B):
- `structural-patterns` → all 10 structural planned pages
- `dev-workflow-patterns` → all 4 dev-workflow planned pages

All other existing-page edges stay the same as graph A.

## Variable quantities

| Quantity | Graph A (baseline) | Graph B (full) |
|----------|--------------------|----------------|
| Nodes | 14 | 28 |
| Observed subset | all 14 | 14 existing pages |
| Hidden subset | — | 14 planned pages |
| Adjacency source | parsed from wiki wikilinks | parsed + hypothesized planned-page links |

## Predictions

| Metric | Existing page | Predicted direction | Why |
|--------|--------------|--------------------|----|
| π | Structural Patterns | **large increase** | Gains 10 inbound links from planned example pages |
| π | Dev-Workflow Patterns | **large increase** | Gains 4 inbound links from planned example pages |
| π | Markov Chain Foundations | moderate increase | All 14 planned pages link to it |
| π | catrace API | moderate increase | All 14 planned pages link to it |
| π | Variant Comparison Methodology | moderate increase | 6 planned pages link to it (D1-D4 + others) |
| π | Example: Self-Healing Nodes | **decrease** | Loses concentration advantage of being the sole Variant Comparison Methodology target |
| π | Joint Kernels and Coupling | increase | Multi-agent planned examples link to it |
| Rank | Structural Patterns | from 13th to top 5 | Structural change, not parameter tuning |
| Rank | catrace Glossary | no major change | Planned pages unlikely to link to it |

## Verdict rule

Claim is supported if:
1. `Structural Patterns` moves up ≥ 5 rank positions
2. `Dev-Workflow Patterns` moves up ≥ 3 rank positions
3. `Example: Self-Healing Nodes` moves down ≥ 2 rank positions
4. `IsTraceOf` = true (stationary consistency theorem holds for the full 28-node graph)

Criterion 4 is mathematical — it must hold. Criteria 1-3 are structural.

---

## Results

> Fill in after running `go run experiments/wiki-knowledge-graph/main.go`

| Page | Naive π | Trace-corrected π | Rank change |
|------|---------|-------------------|-------------|
| Example: Self-Healing Nodes | 0.1477 | | |
| Markov Chain Foundations | 0.1162 | | |
| PDA Triplet Model | 0.1099 | | |
| Joint Kernels and Coupling | 0.1056 | | |
| Trace Chain | 0.0947 | | |
| Agentic Patterns Catalogue | 0.0747 | | |
| Example: Validator Repair | 0.0698 | | |
| catrace API | 0.0684 | | |
| Example: Simple Agent | 0.0659 | | |
| Variant Comparison Methodology | 0.0489 | | |
| Dev-Workflow Patterns | 0.0400 | | |
| Example: Hidden Support System | 0.0237 | | |
| Structural Patterns | 0.0208 | | |
| catrace Glossary | 0.0137 | | |

`IsTraceOf`: [ ]

## Verdict

> Fill in after running.

**Claim:** [ supported / not supported / trade-off ]
**Criteria met:** [ /4 ]

## Interpretation

> Fill in after running.

## Implementation notes

The implementation (`main.go` in this directory) should use catrace's own `Kernel` type and
`Kernel.Trace` method — this experiment is a self-referential test of the library on its own
documentation. See the [[Trace Chain]] wiki page for the mathematical construction.

Dangling node handling: planned pages that have no out-edges into the hidden subset (all their
out-edges go to existing pages) are not dangling in the trace sense — they contribute entirely
to the b·(I-c)⁻¹·d term. Planned pages have no internal (hidden-to-hidden) transitions, so
c = 0 and (I-c)⁻¹ = I. This simplifies the trace formula to:

```
L_A = a + b · d
```

where b is the existing→planned block and d is the planned→existing block.

Verify this simplification holds given the hypothesized planned-page structure before implementing
the general formula.
