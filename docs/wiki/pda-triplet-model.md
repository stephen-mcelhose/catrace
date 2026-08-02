---
title: PDA Triplet Model
tags: [pda, agent, kernel, composition, cyclic-kernels, modeling, hpc]
sources: [docs/math_summary.md, docs/source_writeup.md]
updated: 2026-08-02
---

# PDA Triplet Model

catrace models a single agent as three coupled row-stochastic kernels operating between three finite state spaces:

| Kernel | Map         | Meaning                                                       |
|--------|-------------|---------------------------------------------------------------|
| P      | W → X       | Perception: world state → probability distribution over experiences |
| D      | X → G       | Decision: experience → probability distribution over actions  |
| A      | G → W       | Action effect: action → probability distribution over world states |

These spaces are:
- **W** — world states: what is actually true in the environment
- **X** — experience states: what the agent internally perceives or interprets
- **G** — action states: what the agent does

The triplet is derived from the mathematical framework in Hoffman, Prakash & Chattopadhyay (*Traces of Consciousness*, 2024), extracting only the stochastic/Markov machinery. catrace implements this as the `Agent` struct.

## Derived cyclic kernels

Composing the three kernels gives three square (state-to-same-state) kernels, one for each space:

| Kernel | Formula | Space | Reads as |
|--------|---------|-------|----------|
| Q      | D · A · P | X → X | Experience-to-experience (qualia kernel) |
| S      | A · P · D | G → G | Action-to-action (strategy kernel) |
| W      | P · D · A | W → W | World-to-world (world kernel) |

All three describe the same closed perceive→decide→act→perceive loop, starting from a different point in the cycle. They are cyclic permutations of each other. catrace provides `Agent.QualiaKernel()`, `Agent.StrategyKernel()`, `Agent.WorldKernel()`.

Because Q, S, and W represent the same loop from different perspectives, their stationary distributions are related: if π(Q) is known, one more application of D gives π(S), and one more of P gives π(W). In practice, most analysis is performed on Q (experience space), but W is natural for multi-agent joint kernels (see [[Joint Kernels and Coupling]]).

## Concrete interpretation

The kernel entries have direct semantic meaning:

- **P(w, x)** — given the world is in state w, the probability the agent perceives/interprets it as experience x. High off-diagonal entries encode perception noise.
- **D(x, g)** — given the agent's experience is x, the probability it chooses action g. This is the policy.
- **A(g, w')** — given the agent takes action g, the probability the world moves to state w'. This encodes causal consequence.

An entry Q(x, x') in the composed kernel compresses many possible W→X→G→W→X micro-paths into one effective experience-to-experience probability. See [[Example: Simple Agent]] for a worked walk-through of specific Q entries.

## Validation

catrace's `Agent.Validate()` checks dimensional consistency:
- D.rows == P.cols (experience space consistent)
- D.cols == A.rows (action space consistent)
- A.cols == P.rows (world space consistent)

These checks prevent silent errors in manual kernel specification.

## Multi-agent extension

The PDA triplet extends to multiple agents via product state spaces. For two agents, joint kernels P_joint, D_joint, A_joint are defined over W₁×W₂, X₁×X₂, G₁×G₂. Independent decisions use the Kronecker product D_joint = D₁⊗D₂. Coupling enters through specific kernels — perception coupling in P_joint and action coupling in A_joint — with explicit physical meaning. See [[Joint Kernels and Coupling]] for the full construction.

## Sources

- `docs/math_summary.md`
- `docs/source_writeup.md`
