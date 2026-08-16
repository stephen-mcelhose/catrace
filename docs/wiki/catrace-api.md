---
title: catrace API
tags: [api, go, library, kernel, agent, trace, stationary, ppr, pagerank]
sources: [README.md]
updated: 2026-08-05
---

# catrace API

catrace is a Go library (requires Go 1.22+, uses `gonum`) for finite-state Markov models of autonomous-agent networks. It implements the [[PDA Triplet Model]] mathematical framework and provides analysis tools grounded in [[Markov Chain Foundations]].

## Package layout

| File           | Contents                                              |
|----------------|-------------------------------------------------------|
| `kernel.go`    | Core `Kernel` type, composition helpers, validation   |
| `agent.go`     | `Agent` struct with P, D, A fields and derived kernels |
| `trace.go`     | Trace chain construction and verification             |
| `stationary.go`| Stationary distribution, StationaryFrom, PersonalizedPageRank, TeleportingKernel, entropy rate |
| `analysis.go`  | Communicating/recurrent class decomposition           |
| `graph.go`     | `NewRandomWalkKernel`, `NewTeleportingKernelFromAdj` — build kernels from adjacency matrices |
| `passage.go`   | Mean first-passage times and commute times            |
| `sample.go`    | Sampling, kernel estimation, windowed estimates       |
| `visualise.go` | `ToHTML` — self-contained D3 force-directed graph of any kernel |
| `util.go`      | Helpers                                               |

## Core type: Kernel

`Kernel` wraps a `gonum` dense matrix and enforces row-stochasticity on construction. Optional `StateNames []string` enable human-readable output.

```go
K, err := catrace.NewKernel(mat.NewDense(n, n, data))
```

## Agent

`Agent` holds three kernels (P, D, A) and validates dimensional consistency on construction. Derived kernels are computed on demand:

```go
agent := catrace.Agent{P: p, D: d, A: a}
Q, err := agent.QualiaKernel()    // Q = D·A·P  (experience²)
S, err := agent.StrategyKernel()  // S = A·P·D  (action²)
W, err := agent.WorldKernel()     // W = P·D·A  (world²)
```

See [[PDA Triplet Model]] for the mathematical relationship between Q, S, and W.

## Analysis methods

### Stationary distribution

```go
pi, err := K.Stationary(tol float64, maxIter int)
// Power iteration from a uniform start until convergence within tol

pi, err := K.StationaryFrom(start []float64, tol float64, maxIter int)
// Power iteration from an explicit probability distribution start
```

`Stationary` uses a uniform initial distribution. `StationaryFrom` accepts any valid probability vector as the starting point — useful when the initial distribution carries semantic meaning (e.g. the agent's prior state), or to verify ergodicity by confirming that different starting distributions converge to the same fixed point. On an ergodic kernel, both methods return the same result. See [[Markov Chain Foundations]] for theory.

### Personalized PageRank

```go
ppr, err := K.PersonalizedPageRank(restart []float64, alpha, tol float64, maxIter int)
// Fixed point of x = alpha·restart + (1-alpha)·x·P
```

Each iteration blends the propagated distribution back toward `restart` with weight `alpha`. The fixed point is the stationary distribution of the teleporting chain `alpha·restart·𝟙ᵀ + (1-alpha)·P`. Convergence is guaranteed for any `alpha ∈ (0, 1]` regardless of chain structure.

| Parameter | Meaning |
|---|---|
| `restart` | Goal distribution — the states the agent is persistently drawn toward |
| `alpha`   | Teleportation weight — strength of the pull back to `restart` (standard PageRank: 0.15) |
| `alpha=0` | Degenerates to `StationaryFrom(restart, ...)` |
| `alpha=1` | Returns `restart` unchanged in one step |

See [[Personalized PageRank and Agent Modeling]] for the full agent-modeling interpretation.

### Teleporting kernel

```go
tk, err := K.TeleportingKernel(restart []float64, alpha float64)
// Returns a new *Kernel T = alpha·restart·𝟙ᵀ + (1-alpha)·P
```

Constructs the teleporting Markov chain whose stationary distribution equals the `PersonalizedPageRank` vector for the same `restart` and `alpha`. Useful for visualisation: pass the returned kernel to `ToHTML` and nodes will be sized by PPR score. Use `MinEdge` in `VisualiseOptions` to suppress the low-weight teleportation arcs that connect every node back to the seed set.

| `alpha` | Effect on T                          |
| ------- | ------------------------------------ |
| 0       | T equals P unchanged                 |
| (0, 1)  | Blend of P and the restart broadcast |
| 1       | Every row of T equals `restart`      |

State names are copied from the original kernel. The returned kernel is valid for all other `Kernel` methods (`Classes`, `EntropyRate`, `CommuteTime`, etc.).

### Entropy rate

```go
H, err := K.EntropyRate(base float64)
// H(K) = -Σ_i π_i Σ_j K_{ij} log_b K_{ij}
// base=2 → bits/step
```

### Communicating classes

```go
classes, err := K.Classes(tol float64)
// Returns []CommunicatingClass with {States, IsRecurrent, Period}
```

Uses Kosaraju's algorithm (two DFS passes). Each class carries IsRecurrent and Period (GCD of cycle lengths; 1 = aperiodic). See [[catrace Glossary]] for algorithm details.

### Trace chain

```go
trace, err := parent.Trace(subset []int, tol float64)
// L_A = a + b·(I-c)⁻¹·d
ok, err := trace.IsTraceOf(parent, subset, tol)
// Entry-by-entry verification within tol
```

See [[Trace Chain]] for the full mathematical construction.

### Mean first passage time and commute time

```go
m, err := K.MeanFirstPassage(from, to int)
c, err := K.CommuteTime(i, j int)
// c = m(i→j) + m(j→i)
```

`MeanFirstPassage` solves the linear system (I−Q)·h = 1 on non-target states. Self-MFPT (from == to) returns 0 by library convention. See [[Markov Chain Foundations]].

### Random walk kernel

```go
K, err := catrace.NewRandomWalkKernel(adj [][]float64)
// adj[i][j] = edge weight; zero = no edge. Rows normalized to row-stochastic.
// Stationary distribution has closed form: π_i ∝ degree(i)
```

Builds a Kernel from a weighted adjacency matrix. Rows are normalised to row-stochastic; errors on sink nodes (zero rows). Used in the wiki-knowledge-graph experiment to construct a PageRank-style kernel from wikilink structure.

### Teleporting kernel from adjacency

```go
K, err := catrace.NewTeleportingKernelFromAdj(adj *mat.Dense, restart []float64, alpha float64, names []string)
// T[i][j] = α·restart[j] + (1−α)·(adj[i][j]/rowsum[i])   if rowsum[i] > 0
// T[i][j] = restart[j]                                      if rowsum[i] = 0
```

Combines row-normalisation and teleportation in a single pass. Sink nodes (zero rows) collapse entirely to the restart distribution — no artificial uniform edges inserted. This is the correct primitive for graphs with dangling nodes such as document corpora, where leaf pages have no outgoing links.

| Comparison                  | `NewRandomWalkKernel` + `TeleportingKernel` | `NewTeleportingKernelFromAdj`    |
| --------------------------- | ------------------------------------------- | -------------------------------- |
| Handles sinks               | No — errors                                 | Yes — sinks become `restart`     |
| Sink treatment              | Manual 1/n pre-fill required                | Semantically driven by `restart` |
| Steps                       | 2                                           | 1                                |

### Visualisation

```go
html, err := K.ToHTML(opts *catrace.VisualiseOptions)
// Returns a self-contained HTML file: D3 force-directed graph
// Nodes scaled by stationary mass; edge width by transition probability
// StateNames used as node labels if set
```

Renders any kernel as an interactive browser graph. Each example writes its kernel(s) to `*.html` output files. No external dependencies — the HTML file is fully self-contained.

### Sampling and estimation

```go
next, err := K.Sample(state int, rng *rand.Rand)
seq := catrace.SampleTraceFromSequence(trajectory, subsetSet)
est, err := catrace.EstimateKernelFromSequence(seq, n, smoothing)
windows, err := catrace.WindowedTraceEstimates(traj, subset, winLen, step, smooth)
```

`LeftAction` applies a distribution vector to the kernel for one-step forecasting:

```go
next, err := K.LeftAction(dist []float64)
```

## Example usage

```go
// Single agent
Q, err := agent.QualiaKernel()
pi, err := Q.Stationary(1e-12, 5000)
H, err := Q.EntropyRate(2)
classes, err := Q.Classes(1e-12)

// Trace chain
tr, err := parent.Trace([]int{0, 1}, 1e-12)
ok, err := tr.IsTraceOf(parent, []int{0, 1}, 1e-12)

// First passage
m, err := J.MeanFirstPassage(worstState, bestState)
```

## Implemented examples

| Example              | Key API methods demonstrated               | Pattern(s)            |
|----------------------|--------------------------------------------|-----------------------|
| `simple_agent`       | QualiaKernel, Stationary, EntropyRate, Classes, LeftAction | Augmented LLM, Autonomous Loop |
| `trace_analysis`     | Trace, IsTraceOf, Stationary, Sample, EstimateKernelFromSequence | Hidden support system |
| `validator_repair`   | WorldKernel (joint), Trace, Stationary     | Evaluator-Optimizer, Self-Healing |
| `self_healing_nodes` | WorldKernel (joint), MeanFirstPassage, EntropyRate, Trace | Self-Healing, Autonomous Loop |
| `prompt_chaining`    | NewKernel (assembled W), Trace, IsTraceOf, Stationary, MeanFirstPassage | Prompt Chaining |
| `blackboard`         | NewKernel (assembled J), Trace, IsTraceOf, Stationary, MeanFirstPassage | Blackboard |

See [[Agentic Patterns Catalogue]] for the full pattern coverage map.
See [[Scenario Registry]] for README scenario numbering and implementation status.

## Sources

- `README.md`
