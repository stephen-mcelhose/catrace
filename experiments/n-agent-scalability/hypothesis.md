---
title: "N-agent joint kernel scalability bounds"
type: analysis
workspace: /Users/stephen.mcelhose.ext/repos/catrace
created_at: 2025-01-24T12:00:00Z
skill: explore
---

# Hypothesis: N-agent joint kernel scalability bounds

> Fill in all sections marked [TODO] **before running catrace**. Fill in Results and Verdict after.

## Claim

The current dense-matrix approach in `catrace` becomes computationally intractable at N=4 agents (3 states each, 81 joint states) and practically impossible at N=5 (243 joint states) due to the $O(k^{2N})$ scaling of the transition matrix size and $O(k^{3N})$ scaling of matrix operations like inversion (used in Trace).

## Context

- **Pattern:** Swarm / Multi-agent Coordination
- **Related example:** `examples/self_healing_nodes`
- **Motivation:** Patterns like Swarm (issue #10) and Three-agent majority-valid (issue #9) require 3+ agents. The library needs clear bounds on where dense representations fail to guide future sparse or factored implementations.

## Variant definitions

| Variant | Label | Description |
|---------|-------|-------------|
| A | N=2 | Joint system of 2 agents ($3^2=9$ states) |
| B | N=3 | Joint system of 3 agents ($3^3=27$ states) |
| C | N=4 | Joint system of 4 agents ($3^4=81$ states) |
| D | N=5 | Joint system of 5 agents ($3^5=243$ states) |

*The repeated unit is the node agent from `self_healing_nodes` ($k=3$ world states).*

## Variable kernel entries

The transition matrix $J$ size scales exponentially.

| Variant | Matrix Dims | Total Entries | Memory (float64) |
|---------|-------------|---------------|------------------|
| N=2 | 9x9 | 81 | ~0.6 KB |
| N=3 | 27x27 | 729 | ~5.8 KB |
| N=4 | 81x81 | 6,561 | ~52 KB |
| N=5 | 243x243 | 59,049 | ~472 KB |
| N=6 | 729x729 | 531,441 | ~4.2 MB |

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| Construction time | wall(Qualia/World) | N=5 >> N=2 | Kronecker product nesting cost |
| Trace time | wall(Trace) | N=5 >> N=2 | $O(S^3)$ where $S=3^N$ |
| Memory usage | size(J.P) | N=5 >> N=2 | $O(3^{2N})$ |

*Prediction: N=4 will show noticeable slowdown in Trace, and N=5 will reach the practical limit for interactive use (seconds per operation).*

## Verdict rule

The claim is supported if wall time for `Trace` on N=5 exceeds 5 seconds or if memory growth follows the predicted $k^{2N}$ curve.

---

## Results

> Fill in after running.

| Metric | N=2 | N=3 | N=4 | N=5 | Predicted direction | Correct? |
|--------|-----|-----|-----|-----|--------------------|----|
| Trace wall time | | | | | A < B < C < D | |
| Memory (entries) | 81 | 729 | 6,561 | 59,049 | A < B < C < D | |

## Verdict

**Claim:** [Pending]
**Metrics in agreement:** [0/3]

## Interpretation

[TODO]
