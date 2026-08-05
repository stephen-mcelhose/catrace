---
title: Joint Kernels and Coupling
tags: [joint-kernel, product-space, kronecker, coupling, multi-agent, perception, action]
sources: [GLOSSARY.md, docs/wiki/raw/agent-patterns-catalog-blackboard.md]
updated: 2026-08-05
---

# Joint Kernels and Coupling

When modeling multiple agents, catrace lifts the [[PDA Triplet Model]] to product state spaces. The joint kernel J = P_joint · D_joint · A_joint is a single square Markov kernel over the product world space W₁×W₂ (and analogously for X and G spaces). It is analyzed with the same tools as any kernel: `Stationary`, `EntropyRate`, `Classes`, `Trace`.

## Product state spaces

Given two finite state spaces S₁ and S₂, the product S₁×S₂ has |S₁|·|S₂| states indexed in row-major order: pair (i₁, i₂) maps to index i₁·|S₂| + i₂.

For the two-agent `validator_repair` example:

| Space   | Elements                                 | Size |
|---------|------------------------------------------|------|
| W₁×W₂  | VV, VI, IV, II                           | 4    |
| X₁×X₂  | ok·good, ok·bad, prob·good, prob·bad     | 4    |
| G₁×G₂  | all 9 worker×validator action pairs      | 9    |

The `self_healing_nodes` example uses W = {H,D,O}×{G,B} with 6 joint world states.

## Independent decisions: Kronecker product

When agents decide independently — each based solely on its own experience, without within-cycle communication — the joint decision kernel factors as:

```
D_joint = D₁ ⊗ D₂
```

The Kronecker product of two row-stochastic matrices is row-stochastic; no normalization is needed. The probability of joint action (g₁, g₂) given joint experience (x₁, x₂) is simply:

```
D_joint[(x₁,x₂), (g₁,g₂)] = D₁[x₁, g₁] × D₂[x₂, g₂]
```

In all catrace multi-agent examples, decisions are independent (no within-cycle communication). This is the "no communication" assumption. Coupled decisions — where the joint policy is not a Kronecker product — would require D_joint entries to be specified by hand.

## Coupled perception

Coupling enters P_joint when one agent's perception depends on another agent's world state. In `validator_repair`, the Validator's perception is coupled to the Worker's world state:

```
P_joint[(w₁,w₂), (x₁,x₂)] = P₁[w₁, x₁] × P₂_coupled[(w₁,w₂), x₂]
```

When the Worker is degraded (w₁ = invalid), P₂_coupled raises the Validator's probability of perceiving `looks_bad`, even if the Validator's own state is fine. Without this coupling, the Validator cannot reliably detect Worker degradation.

In `self_healing_nodes`, the evolver's perception is coupled to node health: a sick node depresses the pool score the evolver observes, making a good config look bad even when it isn't (the `D·G → D·B` leakage).

The Planned Blackboard pattern ([[Structural Patterns]] §10; [[Agent Patterns Catalog — Blackboard]]) is the same shape: specialists couple perception to a shared board world state while keeping D_joint a Kronecker product — board-mediated coordination without within-cycle peer messaging.

## Coupled action effect

Coupling enters A_joint when one agent's action changes another agent's world state. In `validator_repair`, the Validator's repair action restores the Worker's world state — raising Worker's probability of becoming valid by ~0.20 above what it would achieve independently:

```
For repair actions (g₂ = repair):
  A_joint[(g₁,repair), (w₁,w₂)] uses a boosted P(worker→valid)
  
For non-repair actions:
  A_joint[(g₁,g₂), (w₁,w₂)] = A₁[g₁,w₁] × A₂[g₂,w₂]  (Kronecker)
```

In `self_healing_nodes` (Variant B), the evolver's `mutate` action boosts the node's `healthy` probability in A_joint; `promote` leaves node recovery unchanged.

## Two construction routes

A joint kernel can be built two ways:

1. **From joint P, D, A** (preferred): Define P_joint, D_joint, A_joint over product spaces and compose J = P·D·A. Coupling enters through specific kernels with physical meaning. Used in `validator_repair` and `self_healing_nodes`.

2. **Directly**: Specify the 4×4 (or larger) matrix by hand. Simpler but loses the P, D, A decomposition. Used in `trace_analysis`.

The composed approach is preferred because changing one agent's kernel propagates through the composition automatically, making structural experiments tractable.

## Analysis

Once J is built, the same tools apply as for single-agent kernels:

- `Stationary` — long-run joint state distribution
- `EntropyRate` — predictability of the joint system
- `Classes` — whether the joint system is ergodic
- `Trace` — reduce to observable subset (e.g., {VV, II} or {H·G, O·B})
- `MeanFirstPassage` — recovery speed from degraded states

See [[Example: Validator Repair]] and [[Example: Self-Healing Nodes]] for worked analyses.

## Sources

- `GLOSSARY.md`
- `docs/wiki/raw/agent-patterns-catalog-blackboard.md` (Blackboard as planned coupled-perception shape; see [[Agent Patterns Catalog — Blackboard]])
