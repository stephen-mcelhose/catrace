---
title: Markov Chain Foundations
tags: [markov, mathematics, stochastic, stationary, entropy-rate, mfpt, communicating-classes]
sources: [docs/math_summary.md, docs/source_writeup.md]
updated: 2026-08-04
---

# Markov Chain Foundations

The mathematical substrate of catrace is classical finite-state Markov chain theory. Every kernel in the library — whether it models a single agent's perception, a joint action effect, or a reduced trace — is a row-stochastic matrix obeying these foundations.

## Row-stochastic matrices

A finite-state system is represented by a matrix P where rows index current states, columns index next states, every entry is non-negative, and every row sums to 1. The entry P(i,j) is the probability of moving from state i to state j in one step:

```
P(i,j) = Pr(X_{t+1} = j | X_t = i)
```

All kernels in catrace enforce this invariant. Composed kernels (Q, S, W, joint kernels) are checked to remain stochastic after multiplication. Small numerical drift from floating-point arithmetic is corrected by row normalization.

## Stationary distributions

A stationary distribution π is a probability vector satisfying π·P = π (equivalently, it is a left eigenvector of P with eigenvalue 1, normalized to sum to 1):

```
π·P = π,   Σ_i π_i = 1,   π_i ≥ 0
```

For an ergodic chain (single recurrent communicating class, aperiodic), the stationary distribution is unique and equals the long-run time average: the chain spends fraction π_i of all steps in state i, regardless of where it started. catrace computes π via power iteration (`Stationary`).

For a [[Trace Chain]], the restricted stationary law is obtained by normalizing the parent stationary mass on the observed subset A:

```
π_A(i) = π(i) / Σ_{j∈A} π(j)   for i ∈ A
```

This identity is the key theorem verified by `IsTraceOf`.

## Communicating classes and recurrence

A communicating class is a strongly connected component (SCC) of the directed graph on states where edges carry positive transition probability. Two states i and j communicate if i can reach j and j can reach i.

A class is **recurrent** (closed) if no state in it has positive probability of leaving — the chain never escapes. A class is **transient** if it has positive probability of eventually leaving. catrace uses Kosaraju's algorithm to identify SCCs (`Classes`), then classifies each by whether it has outgoing edges. For each recurrent class it also computes the **period** — the GCD of all return-cycle lengths. Period 1 (aperiodic) is required for convergence to a unique stationary distribution.

An ergodic chain is irreducible (one recurrent class) and aperiodic. The examples `simple_agent`, `validator_repair`, and `self_healing_nodes` all exhibit ergodic joint kernels.

## Entropy rate

For a stationary Markov chain with distribution π, the entropy rate measures the average uncertainty per step:

```
H(P) = -Σ_i π_i Σ_j P_{ij} log_b P_{ij}
```

With b=2 this gives bits per step. A deterministic chain (P has only 0s and 1s) achieves H=0. A fully random chain (all transitions equiprobable) achieves the maximum. catrace computes entropy rate via `EntropyRate`. In the [[Example: Self-Healing Nodes]] variant comparison, entropy rate distinguishes a self-healing system (lower H, more legible) from one depending on random search (higher H, harder to predict).

## Mean first passage time and commute time

The mean first passage time m(i,j) is the expected number of steps to reach j for the first time starting from i:

```
m(i,j) = E_i[min{t ≥ 1 : X_t = j}]
```

It is computed by solving a linear system: remove state j, then solve (I − Q)·h = 1 on the remaining states, where Q is the transition sub-matrix. catrace provides `MeanFirstPassage(i, j)`.

The commute time C(i,j) = m(i,j) + m(j,i) is symmetric and measures the expected round-trip cost between two states. It appears in [[Example: Self-Healing Nodes]] to characterize recovery speed, and is the key metric in [[Dev-Workflow Patterns]] for quantifying the cost of a review cycle.

## Spectral gap and mixing time

For an ergodic chain, the **spectral gap** is 1 − |λ₂|, where λ₂ is the second-largest eigenvalue of P by magnitude (λ₁ = 1 always). It governs how fast the chain converges to its stationary distribution:

- Large spectral gap → fast mixing (chain forgets its starting state quickly)
- Small spectral gap → slow mixing (chain can be "stuck" in regions for many steps)

The theoretical mixing time bound is proportional to 1 / (spectral gap). Whether this bound is tight enough to be practically useful for comparing catrace kernels — i.e., whether spectral gap reliably rank-orders kernels by empirical mixing speed — is the subject of the pending `experiments/spectral-gap-mixing-time` experiment (#25). If confirmed, `Kernel.SpectralGap()` would be a natural addition to the [[catrace API]], complementing `Stationary` with a convergence-speed measure.

**Total variation distance** is the standard measure of mixing progress:

```
TV(μ, π) = (1/2) Σ_i |μ_i − π_i|
```

Starting from any initial distribution μ, TV(μ·Pᵗ, π) → 0 as t → ∞. The mixing time at threshold ε is the first t where TV < ε for all starting states.

## Practical implementation notes

catrace implements classical stochastic kernels throughout. The HPC source paper (see [[PDA Triplet Model]]) contains notation suggestive of an amplitude-based formulation; catrace ignores those artifacts and enforces:

1. All kernels are row-stochastic on construction.
2. Composed kernels are checked for row-stochasticity after multiplication.
3. Trace kernels use the classical block-matrix formula (see [[Trace Chain]]).
4. Small numerical drift is corrected by row normalization.

## Sources

- `docs/math_summary.md`
- `docs/source_writeup.md`
