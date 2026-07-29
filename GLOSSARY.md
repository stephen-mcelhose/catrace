# Glossary

Math and notation terms used across the walkthroughs and source code.

**Reference:** Levin & Peres, *Markov Chains and Mixing Times* (2nd ed.) — free PDF at https://pages.uoregon.edu/dlevin/MARKOV/markovmixing.pdf. Page numbers below refer to that book.

---

## Stochastic matrix

A rectangular matrix where every **row sums to 1** and all entries are non-negative. Each row is a probability distribution over the next state. Also called a transition matrix or Markov kernel.

```math
\text{row } i: \quad \sum_j M_{ij} = 1, \quad M_{ij} \geq 0
```

*Levin & Peres, §1.1, p. 1.*

---

## Kernel

In this codebase, kernel and stochastic matrix mean the same thing. A kernel maps one state space to another — it encodes "given I am in state i, what is my distribution over output states?" Square kernels (same input and output space) define Markov chains.

---

## Markov chain

A process that jumps between states over discrete time steps, where the next state depends only on the current state — not on how you got there. Fully described by a square stochastic matrix L: entry `L[i,j]` is the probability of moving from state i to state j in one step.

*Levin & Peres, §1.1, p. 1.*

---

## Stationary distribution (π)

A probability distribution over states that does not change when you apply the transition matrix. Written as a row vector, it satisfies:

```math
\pi \cdot Q = \pi, \qquad \sum_i \pi_i = 1
```

Interpretation: if you run the chain for a very long time, π_i is the fraction of steps spent in state i. π is the standard notation — used universally in the Markov chain literature.

*Levin & Peres, §1.5, p. 9.*

---

## Entropy rate

The average uncertainty (in bits) about the next state, weighted by how often each state is visited:

```math
H(Q) = -\sum_i \pi_i \sum_j Q_{ij} \log_2 Q_{ij}
```

Measured in bits/step. A rate of 0 means the chain is fully deterministic. A rate of 1 bit means the next state is a coin flip. Higher values mean more uncertainty per step.

*Levin & Peres, §4.3.*

---

## Communicating class

A set of states where every state can reach every other state (possibly in multiple steps). States that cannot reach each other belong to different communicating classes.

---

## Recurrent class

A communicating class that the chain never leaves. Once you enter it, you stay. Also called an absorbing class. A chain may have multiple recurrent classes.

*Levin & Peres, §1.7.*

---

## Transient state

A state the chain eventually leaves and never returns to. Transient states are not part of any recurrent class.

---

## Ergodic chain

A chain with exactly one recurrent communicating class that is also aperiodic (not stuck in a cycle). Ergodic chains have a unique stationary distribution and converge to it from any starting state.

*Levin & Peres, §1.7.*

---

## Left action

Applying a distribution (row vector) to a transition matrix to get the distribution one step later:

```math
\pi' = \pi \cdot Q
```

"Left" because the distribution multiplies from the left. This is how probability mass flows forward in time.

---

## Matrix multiplication

The core operation for composing kernels. To get entry `(AB)[i,j]`, take row i of A and column j of B, multiply element by element, and sum. The inner dimensions must match: an (m×n) matrix times an (n×p) matrix gives an (m×p) result.

*[3Blue1Brown: Essence of Linear Algebra, Ch. 4](https://www.youtube.com/watch?v=XkY2DOUCWMU)*

---

## Matrix inverse

For a square matrix M, its inverse M⁻¹ satisfies M · M⁻¹ = I (the identity matrix). Not every matrix has an inverse — only those with non-zero determinant. In the trace formula, `(I - L_BB)⁻¹` captures the sum of all hidden-state excursions.

*[3Blue1Brown: Essence of Linear Algebra, Ch. 7](https://www.youtube.com/watch?v=uQhTuRlWMxw)*

---

## Trace (of a kernel)

Given a Markov chain on a large state space and a subset S of states you can observe, the **trace** is the effective chain on S alone. It folds in all hidden excursions through the complement and gives the right long-run frequencies.

Formula: if the kernel is partitioned into blocks over S and its complement B:

```math
\text{Tr}(L) = L_{SS} + L_{SB}(I - L_{BB})^{-1}L_{BS}
```

This is not the same as deleting rows and columns.

---

## Qualia kernel (Q)

The composed kernel `Q = D · A · P` that describes the closed-loop dynamics in **experience space** (X → X). Given how an agent currently reads a situation, Q gives the distribution over how it will read the situation next, after acting and re-observing.

---

## Decision kernel (D)

A rectangular stochastic matrix mapping experience states X to action states G. Row i gives the agent's action distribution when its current experience is state i.

---

## Action effect kernel (A)

A rectangular stochastic matrix mapping action states G to world states W. Row i gives the distribution over world states that result from taking action i.

---

## Perception kernel (P)

A rectangular stochastic matrix mapping world states W to experience states X. Row i gives the agent's experience distribution when the world is in state i.

---

## Strategy kernel (S)

The composed kernel `S = A · P · D` — the same closed loop as Q, but starting from **action space** (G → G).

---

## World kernel (W)

The composed kernel `W = P · D · A` — the same closed loop as Q, but starting from **world space** (W → W).

---

## Fixed point

A distribution π that satisfies π · Q = π — applying one more step leaves it unchanged. The stationary distribution is a fixed point of the transition matrix.

---

## Eigenvalue / eigenvector

For a square matrix M, a vector v and scalar λ satisfying M · v = λ · v. The stationary distribution π is a left eigenvector of Q with eigenvalue 1 — meaning π · Q = 1 · π.

*[3Blue1Brown: Essence of Linear Algebra, Ch. 14](https://www.youtube.com/watch?v=PFDu9oVAE-g)*
