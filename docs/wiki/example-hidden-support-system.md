---
title: "Example: Hidden Support System"
tags: [example, trace, hidden-states, trace-chain, stationary, sampling, estimation]
sources: [examples/trace_analysis/WALKTHROUGH.md, docs/patterns/story-hidden-support-system.md]
updated: 2026-08-02
---

# Example: Hidden Support System

**Pattern:** Augmented LLM (1) with hidden context — see [[Agentic Patterns Catalogue]]
**Code:** `examples/trace_analysis/main.go`
**Run:** `go run examples/trace_analysis/main.go`

This example is catrace's primary demonstration of the [[Trace Chain]]. A focal LLM agent (A) is visible; the surrounding support system (B: retrieval backends, monitors, reviewers, other agents) is hidden. The question: what are the effective health dynamics of A alone, after accounting for all hidden influence?

## State spaces

| State | Visible? | Meaning |
|-------|----------|---------|
| A_valid (0) | ✅ | Focal agent functioning |
| A_invalid (1) | ✅ | Focal agent degraded |
| B_valid (2) | ❌ | Hidden system functioning |
| B_invalid (3) | ❌ | Hidden system degraded |

The full parent kernel L is 4×4. Observed subset S = {0, 1}; hidden subset B = {2, 3}.

## The parent kernel and block decomposition

```
L = | 0.60  0.20  | 0.10  0.10 |   (A_valid row)
    | 0.15  0.55  | 0.15  0.15 |   (A_invalid row)
    |-------------|------------|
    | 0.20  0.20  | 0.40  0.20 |   (B_valid row)
    | 0.10  0.20  | 0.20  0.50 |   (B_invalid row)
```

Blocks: L_SS = top-left 2×2 (direct A→A), L_SB = top-right 2×2 (A→B excursion), L_BS = bottom-left 2×2 (B→A return), L_BB = bottom-right 2×2 (B internal).

## Trace computation

Using the formula from [[Trace Chain]]:

```
L_A = L_SS + L_SB · (I − L_BB)⁻¹ · L_BS

det(I − L_BB) = 0.26
(I − L_BB)⁻¹ ≈ | 1.923  0.769 |
                | 0.769  2.308 |

L_SB · (I − L_BB)⁻¹ ≈ | 0.269  0.308 |
                        | 0.404  0.462 |

b(I−c)⁻¹d ≈ | 0.085  0.115 |
             | 0.127  0.173 |

Trace ≈ | 0.685  0.315 |
        | 0.277  0.723 |
```

**Key comparison — naive cut vs. trace:**

| Transition | Naive (L_SS) | With hidden excursions |
|-----------|-------------|----------------------|
| A_valid → A_valid | 0.60 | 0.685 ↑ |
| A_valid → A_invalid | 0.20 | 0.315 ↑ |
| A_invalid → A_valid | 0.15 | 0.277 ↑ |
| A_invalid → A_invalid | 0.55 | 0.723 ↑ |

Every entry has shifted. The hidden support system makes A_valid more likely to stay valid (0.685 vs 0.60) and A_invalid more likely to stay invalid (0.723 vs 0.55) — the hidden system **amplifies existing states**. Simple deletion of hidden states would give the wrong dynamics.

## Key theorem verification

```go
ok, err := trace.IsTraceOf(parent, []int{0,1}, 1e-12)
// IsTraceOf = true

stationary(parent)|subset normalized = 0.467532  0.532468
stationary(trace)                    = 0.467532  0.532468
```

Matching to six decimal places. The trace gives the right long-run frequencies for A's states even though B is gone. See [[Markov Chain Foundations]] for the stationary restriction theorem.

## Sampling and estimation

100 steps of the parent chain are sampled; B-state visits are discarded; the remaining A-state subsequence is used to estimate the trace kernel. With 100 steps the estimate is rough (expected — B states consume many of the 100 steps). With 10,000+ steps it converges to the analytical trace. `WindowedTraceEstimates` provides overlapping-window estimates useful for detecting non-stationarity.

## Connection to math and API

- Block-matrix trace via `Kernel.Trace([]int{0,1}, 1e-12)` → [[Trace Chain]]
- Verification via `Kernel.IsTraceOf(parent, subset, tol)` → [[Trace Chain]]
- Sampling via `Kernel.Sample`, `SampleTraceFromSequence`, `EstimateKernelFromSequence` → [[catrace API]]
- Stationary restriction theorem → [[Markov Chain Foundations]]

## Sources

- `examples/trace_analysis/WALKTHROUGH.md`
- `docs/patterns/story-hidden-support-system.md`
