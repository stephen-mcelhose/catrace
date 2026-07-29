# Glossary

Math and notation terms used across the walkthroughs and source code.

**References:**
- **[LP]** Levin & Peres, *Markov Chains and Mixing Times* (2nd ed.) — free PDF at https://pages.uoregon.edu/dlevin/MARKOV/markovmixing.pdf
- **[HPC]** Hoffman, Prakash & Chattopadhyay, *Traces of Consciousness*, Preprints 2024 — https://www.preprints.org/manuscript/202410.1305/v1 — the paper this codebase implements; uses kernel, stationary measure, diagram, and community throughout

---

## Stochastic matrix

A rectangular matrix where every **row sums to 1** and all entries are non-negative. Each row is a probability distribution over the next state. Also called a transition matrix or Markov kernel.

```math
\text{row } i: \quad \sum_j M_{ij} = 1, \quad M_{ij} \geq 0
```

*[LP] §1.1, p. 1.*

---

## Kernel

In this codebase, kernel and stochastic matrix mean the same thing. A kernel maps one state space to another — it encodes "given I am in state i, what is my distribution over output states?" Square kernels (same input and output space) define Markov chains.

*[HPC] uses "kernel" as the primary term throughout.*

---

## Markov chain

A process that jumps between states over discrete time steps, where the next state depends only on the current state — not on how you got there. Fully described by a square stochastic matrix L: entry `L[i,j]` is the probability of moving from state i to state j in one step.

*[LP] §1.1, p. 1. [HPC] §2.*

---

## Markov diagram

The weighted directed graph associated with a Markov chain. Each **node** represents a state; each **directed edge** from node i to node j carries weight `L[i,j]` (the transition probability). Edges with probability 0 are omitted.

The diagram makes the chain's structure visible: communicating classes appear as connected components or strongly connected subgraphs, and communities appear as dense clusters with sparse connections between them.

*[HPC] defines this explicitly and uses it to classify free, bound, and confined particles.*

---

## Random walk on a graph

A Markov chain whose state space is the **vertices** of a graph and whose transition probabilities come directly from the edge structure. This is the bridge between graph theory and Markov chain theory — every undirected weighted graph defines a canonical chain, and every chain can be visualized as a walk on its Markov diagram.

### Construction

Given an undirected graph $G = (V, E)$ with non-negative edge weights $w(i, j)$, define the **weighted degree**:

```math
d(i) = \sum_{j \sim i} w(i, j)
```

The random walk transition matrix is:

```math
P(i, j) = \frac{w(i, j)}{d(i)}
```

For an **unweighted** graph ($w = 1$ on every edge): $P(i, j) = 1/\deg(i)$ if $(i,j)\in E$, else $0$.

### Stationary distribution

For an undirected graph the stationary distribution has a closed form — no power iteration needed:

```math
\pi(i) = \frac{d(i)}{\sum_k d(k)} = \frac{d(i)}{2|E|}
```

**High-degree nodes are visited most often.** This is provable from **detailed balance** (reversibility):

```math
\pi(i)\, P(i, j) = \pi(j)\, P(j, i) \quad \text{for all } i, j
```

Detailed balance means the chain looks statistically identical run forwards or backwards — undirected graphs always produce reversible chains.

### Directed graphs

When edges have direction (as in a Markov diagram), detailed balance no longer holds in general and there is no closed-form stationary distribution. You must solve $\pi P = \pi$ numerically — which is what `Stationary()` does via power iteration. All catrace kernels are directed; the undirected case is a special sub-case.

### Commute time and effective resistance

For undirected graphs, commute time has a beautiful closed-form via **effective resistance** — treating the graph as a resistor network where each edge $(i,j)$ has resistance $1/w(i,j)$:

```math
C(i, j) = 2|E| \cdot R_{\text{eff}}(i, j)
```

$R_{\text{eff}}(i,j)$ is the voltage drop from $i$ to $j$ when one unit of current is injected. States that are structurally close in the graph (many short paths between them) have low effective resistance and low commute time. For directed kernels this shortcut does not apply and `CommuteTime()` computes directly from `MeanFirstPassage`.

### Graph Laplacian

The random walk is tightly related to the **graph Laplacian** $L = D - A$, where $D$ is the diagonal degree matrix and $A$ is the adjacency matrix. The random walk matrix is $P = D^{-1}A$, so $I - P = D^{-1}L$. Spectral properties of $L$ (its eigenvalues and eigenvectors) control mixing speed, community structure, and connectivity — this is the foundation of spectral clustering and graph partitioning.

### Why this matters for catrace

Every kernel in catrace **is** a weighted directed graph. The two representations are identical:

| Graph concept         | Kernel concept            |
|-----------------------|---------------------------|
| Vertex                | State                     |
| Directed edge $i→j$   | Positive entry $P(i,j)$   |
| Edge weight           | Transition probability    |
| Strongly connected component | Communicating class |
| High-degree vertex    | Frequently visited state (high $\pi_i$) |
| Community             | Near-decomposable subchain |

Thinking in graphs makes trace intuitive: the trace onto a subset $A$ is what you observe if you can only see certain nodes and all paths through hidden nodes are summed out.

*[LP] §2.1 (random walks on graphs) and §10 (conductance, spectral gap). Lovász, "Random walks on graphs: A survey" (1993) — a readable standalone reference. For the Laplacian connection: Spielman, "Spectral Graph Theory" lecture notes (freely available online).*

---

## Stationary distribution (π)

A probability distribution over states that does not change when you apply the transition matrix. Written as a row vector, it satisfies:

```math
\pi \cdot Q = \pi, \qquad \sum_i \pi_i = 1
```

Interpretation: if you run the chain for a very long time, π_i is the fraction of steps spent in state i. π is the standard notation — used universally in the Markov chain literature.

*[LP] §1.5, p. 9. [HPC] calls this the "stationary measure" and shows that the stationary measure of a trace kernel is the normalized restriction of the parent's stationary measure.*

---

## Entropy rate

The average uncertainty (in bits) about the next state, weighted by how often each state is visited:

```math
H(Q) = -\sum_i \pi_i \sum_j Q_{ij} \log_2 Q_{ij}
```

Measured in bits/step. A rate of 0 means the chain is fully deterministic. A rate of 1 bit means the next state is a coin flip. Higher values mean more uncertainty per step.

*[LP] §4.3. [HPC] uses entropy rate as a signature of binding — bound particles show a lowering of entropy rate within their community.*

---

## Communicating class

A set of states where every state can reach every other state (possibly in multiple steps). States that cannot reach each other belong to different communicating classes. In the Markov diagram, each communicating class corresponds to a strongly connected component.

*[LP] §1.7. [HPC] connects communicating classes to free particles — a disconnected diagram with two components.*

---

## Community

A generalization of communicating class for chains that are not cleanly decomposable. A community is a subset of states that are **highly interconnected** among themselves, with only weak transitions to and from the rest of the chain. Detected via algorithms such as infomap, spectral clustering, or modularity maximization on the Markov diagram.

*[HPC] uses communities to model bound particles — states that cluster together but are not fully isolated.*

Common algorithms for detecting communities in a Markov diagram: [Infomap](https://mapequation.org), spectral clustering, and modularity maximization. None are currently implemented in catrace.

---

## Recurrent class

A communicating class that the chain never leaves. Once you enter it, you stay. Also called an absorbing class. A chain may have multiple recurrent classes.

*[LP] §1.7.*

---

## Transient state

A state the chain eventually leaves and never returns to. Transient states are not part of any recurrent class.

---

## Ergodic chain

A chain with exactly one recurrent communicating class that is also aperiodic (not stuck in a cycle). Ergodic chains have a unique stationary distribution and converge to it from any starting state.

*[LP] §1.7. [HPC] uses ergodicity as a precondition: "if [the chain] is ergodic then [the observer] sees the asymptotic probabilities of states it shares with [the parent]."*

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

This is not the same as deleting rows and columns. The trace order — the partial order where K is below L when K is a trace of L — is the central object of study in [HPC].

*[HPC] §3, defines the trace order and shows it forms a non-Boolean logic.*

---

## Qualia kernel (Q)

The composed kernel `Q = D · A · P` that describes the closed-loop dynamics in **experience space** (X → X). Given how an agent currently reads a situation, Q gives the distribution over how it will read the situation next, after acting and re-observing.

*[HPC] §2, conscious agent (CA) formalism.*

---

## Decision kernel (D)

A rectangular stochastic matrix mapping experience states X to action states G. Row i gives the agent's action distribution when its current experience is state i.

*[HPC] §2.*

---

## Action effect kernel (A)

A rectangular stochastic matrix mapping action states G to world states W. Row i gives the distribution over world states that result from taking action i.

*[HPC] §2.*

---

## Perception kernel (P)

A rectangular stochastic matrix mapping world states W to experience states X. Row i gives the agent's experience distribution when the world is in state i.

*[HPC] §2.*

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

---

## Kosaraju's algorithm

A graph algorithm that finds all **strongly connected components** (SCCs) of a directed graph in two depth-first search passes. In the context of a Markov chain, each SCC corresponds to a communicating class.

**Pass 1:** DFS on the original graph; record finish times.
**Pass 2:** DFS on the reversed graph in reverse finish-time order; each DFS tree is one SCC.

Applied in `Classes()` to decompose the chain:
- SCCs with no outgoing edges to other SCCs → **recurrent classes** (chain never leaves)
- SCCs with outgoing edges → **transient classes** (chain eventually leaves)

After SCCs are found, `Classes()` also computes the **period** of each recurrent class — the GCD of all cycle lengths within it. Period 1 = aperiodic (required for ergodicity and convergence to stationarity). Period > 1 = the chain cycles and does not converge.

Not discussed in [HPC] — standard Markov chain structure theory. *[LP] §1.7.*

---

## Mean first passage time

The expected number of steps to reach state j for the first time, starting from state i:

```math
m(i, j) = \mathbb{E}[\min\{t \geq 1 : X_t = j\} \mid X_0 = i]
```

Computed by solving a linear system: remove state j from the chain, then solve (I − Q) · m = 1 where Q is the sub-matrix of transitions among the remaining states.

In `passage.go` as `MeanFirstPassage(i, j)`.

---

## Commute time

The expected number of steps to travel from state i to state j **and back**:

```math
C(i, j) = m(i, j) + m(j, i)
```

Symmetric by definition — `C(i,j) = C(j,i)`. A natural notion of distance between states: states that are easy to travel between in both directions have low commute time.

On an undirected graph with the standard random walk, commute time has a closed form in terms of the graph's effective resistance. For directed chains it must be computed from mean first passage times directly.

In `passage.go` as `CommuteTime(i, j)`. Connects to the random-walk-on-graphs interpretation: commute time is the graph-theoretic distance the stationary distribution implicitly defines.
