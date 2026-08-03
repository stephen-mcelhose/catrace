# Hypothesis: Trace collapse restores tractability for N-agent networks

## Claim

The trace chain is the principled answer to state-space explosion in multi-agent catrace models.
A network of N agents with k world states each produces a joint kernel over k^N states — but the
vast majority of those states are not directly observable or actionable. The claim is that
collapsing the joint kernel onto a small observable subset via `Trace` is exact (not an
approximation), stationary-consistent, and produces an interpretable kernel that answers the
question you actually care about — regardless of how large the full joint space is.

In other words: the state-space explosion is not a problem to be solved by sparse matrices or
approximations. It is dissolved by the trace.

## Context

- **Pattern:** Swarm / Peer-to-peer (10), any N > 2 agent topology
- **Related example:** `examples/self_healing_nodes` (2-agent trace to {H·G, O·B})
- **Motivation:** Patterns like Swarm (#10) and Supervisor/Hierarchical (#9) require 3+ agents.
  The joint state space grows as k^N, but the questions we want to answer — "is the majority
  healthy?", "how fast does the swarm recover from full degradation?" — live in a 2- or 3-state
  coarse view. The trace is how you get there without approximation.

## Variant definitions

| Variant | Label       | Description                                                       |
|---------|-------------|-------------------------------------------------------------------|
| A       | Full view   | 4-agent joint kernel J over 3^4 = 81 world states; full analysis |
| B       | Coarse view | Trace of J onto 2 states: majority-healthy vs majority-degraded   |

*The repeated unit is the node agent from `self_healing_nodes` (3 world states: Healthy, Overload,
Failed). Four independent nodes, coupled only through a shared evolver action effect (global
config quality).*

## Variable kernel entries

The only thing that changes between A and B is the level of description — the underlying system
is identical.

| Property           | Variant A (full)             | Variant B (trace)             |
|--------------------|------------------------------|-------------------------------|
| State space        | 81 joint states              | 2 coarse states               |
| Kernel size        | 81 × 81 = 6,561 entries      | 2 × 2 = 4 entries             |
| Stationary dist.   | π over 81 states             | π_A (normalized restriction)  |
| Interpretability   | Requires summing many states | Directly readable             |

**Observable subset definition:**

- *Majority-healthy*: ≥ 3 of 4 nodes in Healthy state
- *Majority-degraded*: ≥ 3 of 4 nodes in Overload or Failed state

All other joint states (mixed configurations) are hidden — they are folded into the trace via
the `b·(I-c)⁻¹·d` term.

## Predictions

State the expected direction for each metric *before running*:

| Metric                        | Expression                            | Predicted direction          | Why                                                                               |
|-------------------------------|---------------------------------------|------------------------------|-----------------------------------------------------------------------------------|
| Stationary consistency        | `IsTraceOf(J, subset, tol)`           | true                         | Stationary-restriction theorem is exact by construction                           |
| π agreement                   | π(J)\|subset normalized vs π(trace)   | Equal to machine precision   | Same theorem                                                                      |
| Majority-healthy mass         | π_trace(majority-healthy)             | > 0.5                        | Independent nodes each spend majority of time healthy; joint majority follows     |
| Recovery speed                | MFPT(majority-degraded → majority-healthy) | < MFPT on single node   | 4 nodes recovering independently; first to reach healthy shifts the majority fast |
| Coarse kernel interpretability | Majority-healthy self-loop           | > 0.7                        | Strong individual self-healing should aggregate to strong collective self-healing  |

## Verdict rule

Claim is supported if:
1. `IsTraceOf = true` (necessary — if this fails the trace formula is broken)
2. π agreement is exact (to `tol = 1e-9`)
3. The coarse kernel is interpretable: majority-healthy self-loop > 0.5 and MFPT(degraded→healthy) < MFPT(healthy→degraded)

All three must hold. This is not a "majority vote" experiment — the trace is either exact or it
isn't.

---

## Results

> Fill in after running.

| Metric                                      | Value | Predicted         | Correct? |
|---------------------------------------------|-------|-------------------|----------|
| IsTraceOf                                   |       | true              |          |
| π agreement (max abs diff)                  |       | < 1e-9            |          |
| π_trace(majority-healthy)                   |       | > 0.5             |          |
| MFPT(majority-degraded → majority-healthy)  |       | finite, < single-node MFPT | |
| Majority-healthy self-loop                  |       | > 0.7             |          |

## Verdict

**Claim:** [supported / not supported]
**Metrics in agreement:** [N/3]

## Interpretation

[What the numbers mean architecturally. If supported: the trace is the right tool for N-agent
analysis and the state-space size is irrelevant to the quality of the answer. If not supported:
identify which step fails — construction, stationary computation, or the trace formula itself.]

## Scaling note (secondary)

For reference, the joint state space sizes for this agent unit (k=3):

| N agents | Joint states | Matrix entries | Memory (float64) |
|----------|-------------|----------------|------------------|
| 2        | 9           | 81             | ~0.6 KB          |
| 3        | 27          | 729            | ~5.8 KB          |
| 4        | 81          | 6,561          | ~52 KB           |
| 5        | 243         | 59,049         | ~472 KB          |
| 6        | 729         | 531,441        | ~4.2 MB          |

The trace output is always small (2–4 states). Construction of the full joint kernel is a one-time
cost; all subsequent analysis is on the collapsed chain.
