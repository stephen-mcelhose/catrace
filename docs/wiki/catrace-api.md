---
title: catrace API
tags: [api, go, library, kernel, agent, trace, stationary]
sources: [README.md]
updated: 2026-08-02
---

# catrace API

catrace is a Go library (requires Go 1.22+, uses `gonum`) for finite-state Markov models of autonomous-agent networks. It implements the [[PDA Triplet Model]] mathematical framework and provides analysis tools grounded in [[Markov Chain Foundations]].

## Package layout

| File           | Contents                                              |
|----------------|-------------------------------------------------------|
| `kernel.go`    | Core `Kernel` type, composition helpers, validation   |
| `agent.go`     | `Agent` struct with P, D, A fields and derived kernels |
| `trace.go`     | Trace chain construction and verification             |
| `stationary.go`| Stationary distribution and entropy rate              |
| `analysis.go`  | Communicating/recurrent class decomposition           |
| `passage.go`   | Mean first-passage times and commute times            |
| `sample.go`    | Sampling, kernel estimation, windowed estimates       |
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
// Power iteration until convergence within tol
```

Returns the unique stationary distribution for an ergodic kernel. See [[Markov Chain Foundations]] for theory.

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

See [[Agentic Patterns Catalogue]] for the full pattern coverage map.

## Sources

- `README.md`
