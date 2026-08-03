# Hypothesis: Spectral gap as predictor of empirical mixing time

## Claim

A Markov kernel's spectral gap (1 − |λ₂|) is a tight enough bound to be practically useful for comparing `catrace` kernels. Specifically, the spectral gap rank-orders kernels correctly by their empirical mixing speed (time to reach a total variation distance threshold from stationary). A kernel with a significantly larger spectral gap will mix in fewer steps.

## Context

- **Pattern:** Core Theory / Measurement
- **Related example:** `examples/simple_agent`, `examples/self_healing_nodes`, `examples/validator_repair`
- **Motivation:** The `Stationary()` method tells us where a chain ends up, but not how fast it gets there. If spectral gap is a reliable predictor, it would be a valuable addition to the `catrace` API (`Kernel.SpectralGap()`), allowing architects to compare the convergence performance of different agent designs without full simulation.

## Variant definitions

Instead of A/B variants of a single system, the "variants" are existing kernels from the `catrace` examples, which provide a range of mixing behaviors.

| Variant | Kernel Label | Description |
|---------|--------------|-------------|
| A | `simple_agent.Q` | Small, well-connected 3-state chain. |
| B | `self_healing_nodes.J` (Var A) | 6-state joint chain with strong local recovery. |
| C | `self_healing_nodes.J` (Var B) | 6-state joint chain with weak local recovery. |
| D | `validator_repair.J` | Larger joint chain representing a consensus validator. |

## Variable kernel entries

The kernels differ in their entire structure and transition probabilities as defined in their respective examples. This experiment tests if the spectral gap generalizes across these diverse structures.

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| Mixing Time | steps to TV(π_t, π) < 0.01 | A < B < C < D (rank) | We expect the simple agent to mix fastest and the complex validator repair to mix slowest. |
| Spectral Gap | 1 - |λ₂| | A > B > C > D (rank) | Larger gap should correlate with faster mixing (fewer steps). |

## Verdict rule

The claim is supported if the rank-ordering of kernels by Spectral Gap is the inverse of their rank-ordering by Empirical Mixing Time for at least 3 out of the 4 variants.

---

## Results

> Fill in after running.

| Variant | Spectral Gap | Empirical Mixing Time (steps) | Rank (Gap) | Rank (Speed) |
|---------|--------------|-------------------------------|------------|--------------|
| A | | | | |
| B | | | | |
| C | | | | |
| D | | | | |

## Verdict

**Claim:** [pending]

## Interpretation

[To be filled after results are obtained.]
