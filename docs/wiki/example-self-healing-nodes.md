---
title: "Example: Self-Healing Nodes"
tags: [example, self-healing, variant-comparison, mfpt, entropy-rate, joint-kernel, coupling]
sources: [examples/self_healing_nodes/WALKTHROUGH.md, docs/patterns/story-self-healing-nodes.md]
updated: 2026-08-02
---

# Example: Self-Healing Nodes

**Pattern:** Self-Healing / Adaptive (14), Autonomous Agent Loop (7) — see [[Agentic Patterns Catalogue]]
**Code:** `examples/self_healing_nodes/main.go`
**Run:** `go run examples/self_healing_nodes/main.go`

This example introduces catrace's **variant comparison pattern**: run the same [[Joint Kernels and Coupling]] architecture under two parameter regimes, then let the stationary distribution, MFPT, and entropy rate vote on which mechanism does the real work. It is also the first example to use `MeanFirstPassage`.

## The design question

A network node monitors its own error rate (EMA) and throttles when errors climb. An outer evolutionary loop watches pool-level throughput and mutates the node's configuration when performance drops. Which loop is primarily responsible for recovery?

- **Variant A — "throttle does it all":** The node's own throttle action has a strong recovery effect; the evolver's contribution is marginal.
- **Variant B — "evolver matters":** The throttle is weaker; the evolver's mutation search is the primary recovery pathway.

## Two agents

**Node agent:**
- World states: `healthy`, `degraded`, `overloaded`
- Experience states: `ema_low`, `ema_mid`, `ema_high` (EMA bands)
- Actions: `push`, `throttle`, `idle`

**Evolver agent:**
- World states: `good_strategy`, `bad_strategy`
- Experience states: `high_score`, `low_score` (pool-level throughput×success²)
- Actions: `promote`, `mutate`

**Joint world states:** `H·G`, `H·B`, `D·G`, `D·B`, `O·G`, `O·B` (6 states)

## Coupling

**Perception coupling (P_joint):** A sick node depresses the pool score the evolver observes, making a good config look bad even when it isn't. This creates `D·G → D·B` leakage: the evolver mutates away from an actually-good config because the node's degradation masks it.

**Action coupling (A_joint, Variant B only):** When the evolver mutates while the node throttles, the node's `healthy` recovery probability is boosted before renormalization (+0.25 boost). In Variant A the throttle's own `healthy` entry (0.65) is already strong; the mutation adds little.

**Independent decisions:** D_joint = D_node ⊗ D_evolver.

## Results: the numbers vote

### Stationary distribution

| State | Variant A | Variant B |
|-------|-----------|-----------|
| H·G   | **0.524** | 0.408 |
| H·B   | 0.181 | 0.153 |
| D·G   | 0.186 | **0.265** |
| D·B   | 0.063 | 0.089 |
| O·G   | 0.034 | 0.064 |
| O·B   | 0.012 | 0.021 |

In Variant A, 52% of time in the best joint state. In Variant B, 41% — degraded states absorb ~10 more percentage points. The evolver carries load it doesn't carry in A, but recovery runs through random search, keeping the system off the best state more often.

### Mean first passage time (O·B → H·G)

| Variant | MFPT worst→best |
|---------|-----------------|
| A | 2.03 steps |
| B | 2.56 steps |

Recovery from the worst joint state takes 26% longer in Variant B. When the evolver must find a good config before the node can heal, the climb back is slower.

### Entropy rate

| Variant | Entropy rate |
|---------|-------------|
| A | 1.87 bits/step |
| B | 2.12 bits/step |

A system that self-heals locally is more legible — you have a better guess where it will be next cycle. A system whose recovery depends on a random search loop wanders more (0.25 bits/step difference).

### Trace onto {H·G, O·B}

| | H·G → O·B | O·B → H·G | IsTraceOf |
|--|-----------|-----------|-----------|
| Variant A | 0.022 | 0.974 | ✅ |
| Variant B | 0.049 | 0.945 | ✅ |

An outside observer who can only see "everything fine" vs "everything failing" sees more than twice the tip-to-worst rate in Variant B (4.9% vs 2.2% per cycle).

## The variant comparison as catrace methodology

This example shows catrace as a **measuring instrument for architectural claims**:
1. Encode each architectural claim as a parameter regime (A or B)
2. Run the same joint kernel structure under both
3. Read stationary distribution, MFPT, entropy rate
4. The numbers say which claim the model supports

The formalized hypothesis record for this comparison lives in `experiments/nodes-throttle-vs-evolver/hypothesis.md`. The general methodology — including the hypothesis template and guidance on well-formed claims — is documented in [[Variant Comparison Methodology]].

## Connection to math and API

- Joint kernel construction → [[Joint Kernels and Coupling]]
- MFPT via `Kernel.MeanFirstPassage(fromState, toState)` → [[catrace API]], [[Markov Chain Foundations]]
- Entropy rate via `Kernel.EntropyRate(2)` → [[catrace API]]
- Trace and IsTraceOf → [[Trace Chain]]
- A_joint renormalization with boost: standard row-normalization after adding boost amount

## Sources

- `examples/self_healing_nodes/WALKTHROUGH.md`
- `docs/patterns/story-self-healing-nodes.md`
