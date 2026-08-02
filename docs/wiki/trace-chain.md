---
title: Trace Chain
tags: [trace, markov, reduction, hidden-states, observation, kernel]
sources: [docs/math_summary.md]
updated: 2026-08-02
---

# Trace Chain

The trace chain is the central mathematical tool in catrace. Given a Markov kernel L on a large state space N and an observed subset A ⊆ N, the trace chain on A is an effective Markov kernel that captures the long-run dynamics of the chain restricted to A — including all the hidden influence of states outside A.

## The block-matrix formula

Partition the states into the observed subset A and its complement B = N \ A. Reorder the matrix accordingly:

```
L = | a   b |
    | d   c |
```

where block a is A×A (direct A-to-A transitions), b is A×B (A exits into B), c is B×B (B internal), and d is B×A (B returns to A).

The trace chain on A is:

```
L_A = a + b · (I − c)⁻¹ · d
```

when (I−c)⁻¹ exists (guaranteed when B contains no recurrent states).

**Interpretation:**
- `a` captures direct one-step transitions within A.
- `b · (I−c)⁻¹ · d` folds in all excursions: every path that leaves A into B and eventually returns, summed over all possible excursion lengths and routes.

The result L_A is itself row-stochastic. catrace computes it via `(*Kernel).Trace(subset []int, tol float64)`.

## Key theorem: stationary consistency

The stationary distribution of the trace chain equals the parent's stationary distribution restricted and normalized to A:

```
π_A(i) = π(i) / Σ_{j∈A} π(j)   for i ∈ A
```

This means: if you observe a long run of the parent chain and record only the visits to A (in order), the empirical long-run frequencies of A-states converge to exactly π_A. The trace is not an approximation — it is the exact effective kernel for the observed subchain.

catrace verifies this with `(*Kernel).IsTraceOf(parent, subset, tol)`, which recomputes the trace formula and checks the result entry-by-entry.

## Usage in catrace examples

Every example demonstrates trace in a different context:

| Example | Parent chain | Observed subset | What trace reveals |
|---------|-------------|-----------------|-------------------|
| [[Example: Hidden Support System]] | Full 4-state system (focal agent A + hidden system B) | {A_valid, A_invalid} | Effective health dynamics of A alone |
| [[Example: Validator Repair]] | Joint 4-state kernel over (Worker, Validator) | {VV, II} | Coarse healthy-vs-failed picture |
| [[Example: Self-Healing Nodes]] | Joint 6-state kernel over (Node, Evolver) | {H·G, O·B} | Best-vs-worst-state effective transitions |

In all three cases, `IsTraceOf = true` confirms the analytical trace matches what would be observed empirically in a long simulation.

## Sampling and estimation

catrace also provides empirical trace tools:
- `SampleTraceFromSequence` — extract only the A-state visits from a simulated trajectory.
- `EstimateKernelFromSequence` — count consecutive A-state pairs to form a frequency estimate.
- `WindowedTraceEstimates` — estimate the trace kernel in sliding windows, useful for detecting drift.

Short trajectories produce noisy estimates; 10,000+ steps of the parent chain are needed for convergence in typical examples. See [[Markov Chain Foundations]] for the theoretical guarantee.

## Sources

- `docs/math_summary.md`
