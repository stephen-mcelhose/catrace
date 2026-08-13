---
title: catrace Glossary
tags: [glossary, notation, definitions, terminology]
sources: [GLOSSARY.md]
updated: 2026-08-13
---

# catrace Glossary

This page synthesizes the formal definitions from `GLOSSARY.md`, grouping them by concept area and cross-referencing the wiki pages where each concept is developed in depth.

## Core kernel concepts

**Stochastic matrix / kernel** — A rectangular matrix where every row sums to 1 and all entries are non-negative. Each row is a probability distribution over the next state. In catrace, "kernel" and "stochastic matrix" are synonyms. Square kernels (same input and output space) define Markov chains. See [[Markov Chain Foundations]].

**Fixed point / stationary distribution** — A distribution π satisfying π·Q = π. Applying one more transition step leaves it unchanged. Computed via `Stationary()`. See [[Markov Chain Foundations]].

**Eigenvalue / eigenvector** — The stationary distribution is a left eigenvector of Q with eigenvalue 1: π·Q = 1·π.

## PDA model kernels

**Qualia kernel (Q)** — Q = D·A·P. The closed-loop experience-to-experience kernel. Row i gives the distribution over next experience states given current experience i. The primary object of analysis in single-agent examples. See [[PDA Triplet Model]].

**Strategy kernel (S)** — S = A·P·D. Same loop starting from action space (G→G). Cyclic permutation of Q.

**World kernel (W)** — W = P·D·A. Same loop starting from world space (W→W). The natural object for multi-agent joint analysis. See [[PDA Triplet Model]].

**Decision kernel (D)** — X→G. Row i gives the agent's action distribution when experience is state i. [HPC §2]

**Action effect kernel (A)** — G→W. Row i gives the world-state distribution resulting from action i. [HPC §2]

**Perception kernel (P)** — W→X. Row i gives the experience distribution when the world is in state i. [HPC §2]

## Multi-agent concepts

**Product state space** — Given S₁ and S₂, the product S₁×S₂ has |S₁|·|S₂| states, indexed row-major: pair (i₁, i₂) → index i₁·|S₂| + i₂. See [[Joint Kernels and Coupling]].

**Kronecker product (⊗)** — For matrices M₁ (p×q) and M₂ (r×s), M₁⊗M₂ is (p·r)×(q·s) with entry (i₁r+i₂, j₁s+j₂) = (M₁)_{i₁j₁}·(M₂)_{i₂j₂}. If both inputs are row-stochastic, the product is row-stochastic. Used for D_joint = D₁⊗D₂ (independent decisions). See [[Joint Kernels and Coupling]].

**Joint kernel** — A stochastic matrix over a product state space. Can be built from joint P, D, A or specified directly. See [[Joint Kernels and Coupling]].

**Coupled perception** — When one agent's perception kernel depends on another agent's world state. See [[Joint Kernels and Coupling]].

**Coupled action effect** — When one agent's action changes another agent's world state. See [[Joint Kernels and Coupling]].

**Independent decisions** — D_joint = D₁⊗D₂; each agent chooses based on its own experience only. All catrace multi-agent examples use independent decisions.

## Structural analysis

**Communicating class** — A strongly connected component (SCC) of the directed graph induced by positive-probability transitions. Computed via Kosaraju's algorithm in `Classes()`:

- Pass 1: DFS on the original graph; record finish times.
- Pass 2: DFS on the reversed graph in reverse finish-time order; each DFS tree is one SCC.

SCCs with no outgoing edges to other SCCs are **recurrent classes** (chain never leaves). SCCs with outgoing edges are **transient classes** (chain eventually leaves). After SCC identification, `Classes()` computes the **period** of each recurrent class — the GCD of all cycle lengths. Period 1 = aperiodic (convergence to stationarity). See [[Markov Chain Foundations]].

## Teleportation and personalization

**Teleporting kernel** — The chain `T = α·v·1ᵀ + (1−α)·P` that blends a restart distribution `v` with a base kernel `P`. At each step, with probability `α` the walk resets to `v`; with probability `1−α` it follows `P`. Ergodic for any `α > 0`, regardless of whether `P` is irreducible. Built via `TeleportingKernel(restart, α)` on an existing `*Kernel`, or directly from a raw adjacency matrix via `NewTeleportingKernelFromAdj`. See [[Personalized PageRank and Agent Modeling]] and [[catrace API]].

**Restart distribution** — The probability vector `v` in the teleporting kernel. Encodes where the walk teleports to. In agent modeling `v` represents goal states; in document-graph analysis `v` is concentrated on seed pages (e.g. gate-output documents). Must be a valid probability vector (non-negative, sums to 1).

**Sink node** — A node with no outgoing edges (all-zero row in the adjacency matrix). `NewRandomWalkKernel` errors on sinks because the row of P is undefined. `NewTeleportingKernelFromAdj` handles sinks natively: a sink's row in T collapses to the restart distribution `v` (no artificial uniform edges inserted). In document graphs, leaf pages with no outgoing relative links are common sinks.

**Personalized PageRank (PPR)** — The stationary distribution of the teleporting kernel `T = α·v·1ᵀ + (1−α)·P`. PPR(v, α) biases stationary mass toward the support of `v`, blending structural importance (from P) with goal-relevance (from v). The unique fixed point is guaranteed for `α > 0`. Computed via `PersonalizedPageRank(restart, α, tol, maxIter)` using power iteration. See [[Personalized PageRank and Agent Modeling]].

## Time metrics

**Mean first passage time (MFPT)** — m(i,j): expected steps to reach j from i for the first time. Computed by solving (I−Q)·h = 1 on states excluding j, with h_j = 0. In catrace: `MeanFirstPassage(i, j)`. See [[Markov Chain Foundations]].

**Commute time** — C(i,j) = m(i,j) + m(j,i). Symmetric round-trip expected cost. A natural distance between states. In catrace: `CommuteTime(i, j)`.

## References

- **[LP]** Levin & Peres, *Markov Chains and Mixing Times* (2nd ed.) — https://pages.uoregon.edu/dlevin/MARKOV/markovmixing.pdf
- **[HPC]** Hoffman, Prakash & Chattopadhyay, *Traces of Consciousness*, Preprints 2024 — https://www.preprints.org/manuscript/202410.1305/v1

## Sources

- `GLOSSARY.md`
