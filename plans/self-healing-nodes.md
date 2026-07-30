# Plan: Self-adjusting / self-healing network nodes example

Inspired by https://github.com/caseylmanus/workerpool

---

## Story

Three network nodes route traffic independently. Each node monitors its own
error rate via an EMA (exponential moving average) and throttles itself when
errors climb — the same proportional backpressure the workerpool applies via
its PID controller. When the whole pool is performing poorly, an outer
evolutionary loop promotes a healthier configuration — the same mutation /
elitism cycle the workerpool uses to search strategy space.

The example models this as two coupled agents:

- **Node agent** — perceives its own health (EMA error band), decides how
  hard to push (normal / throttled / idle), and its action affects its own
  world state and its neighbors' load
- **Evolver agent** — perceives pool-level score (throughput × success²),
  decides to promote the current config or mutate toward a healthier one,
  and its action shifts the node agent's decision probabilities

The joint kernel J over (node_health × evolver_mode) is the object of
analysis. The questions catrace can answer that simulation cannot:

- What is the stationary distribution over joint health × evolver states?
  Where does the system actually spend its time?
- Does the evolutionary search concentrate probability mass near healthy
  configurations, or does it wander?
- What is the mean first passage time from a degraded state back to healthy?
- Tracing onto node health alone: what does an outside observer see, with
  the evolver's internal mode summed out?

---

## State spaces

### Node agent

| Layer   | States                              | Meaning                                      |
|---------|-------------------------------------|----------------------------------------------|
| World   | `healthy`, `degraded`, `overloaded` | Actual node condition                        |
| Experience | `ema_low`, `ema_mid`, `ema_high` | EMA error band: <0.10 / 0.10–0.40 / >0.40   |
| Action  | `push`, `throttle`, `idle`          | How hard the node works this cycle           |

EMA bands are derived from the workerpool's alpha=0.2 update rule. Given
true error rate p, the stationary EMA distribution is geometric — the
probability of landing in each band can be computed analytically.

### Evolver agent

| Layer   | States                     | Meaning                                           |
|---------|----------------------------|---------------------------------------------------|
| World   | `good_strategy`, `bad_strategy` | Whether the current config is actually performing |
| Experience | `high_score`, `low_score`  | Pool score = mbps × (succ/total)²                |
| Action  | `promote`, `mutate`        | Elitist clone vs. random mutation (75/25 split)   |

### Joint world states (product)

`healthy × good`, `healthy × bad`, `degraded × good`, `degraded × bad`,
`overloaded × good`, `overloaded × bad` — 6 joint states.

---

## Kernel derivation

### P_node (3×3): world → experience

Derived analytically from the EMA dynamics. With alpha=0.2 and true error
rates estimated per world state:

- `healthy` → error rate ≈ 0.05 → EMA mostly low
- `degraded` → error rate ≈ 0.35 → EMA spread across low/mid
- `overloaded` → error rate ≈ 0.70 → EMA mostly high

Compute the stationary EMA distribution for each true error rate to fill
the rows of P_node. These are the numbers that would otherwise require a
simulation run to estimate.

### D_node (3×3): experience → action

Derived from the PID logic. Higher EMA → more throttle:

- `ema_low`  → push with high probability (≈0.75), rarely throttle
- `ema_mid`  → split between push and throttle
- `ema_high` → throttle dominates; idle possible if saturated

### A_node (3×3): action → world

- `push` from `healthy`: stays healthy with high probability; small chance
  of degradation under sustained load
- `push` from `degraded`: likely stays degraded or tips to overloaded
- `throttle` from `degraded`: recovery pathway — boosts probability of
  returning to healthy (self-healing)
- `idle`: conservative; slowly drifts toward healthy

### P_evolver (2×2): world → experience

- `good_strategy` → `high_score` with high probability
- `bad_strategy`  → `low_score` with high probability (with some noise)

### D_evolver (2×2): experience → action

- `high_score` → `promote` (elitism: keep best)
- `low_score`  → `mutate` with 75% probability (workerpool's 0.75 mutation
  rate hardcoded in `selectNextStrategy`)

### A_evolver (2×2): action → world — **coupled to node**

This is where the two agents interact. `mutate` slightly shifts the node
agent's world distribution toward a better strategy — modeled as an
increased probability of transitioning from `degraded` → `healthy` in
A_node for the next cycle.

- `promote`: world state probabilities unchanged
- `mutate`: stochastically improves node world distribution (±Δ on the
  healthy transition probability, reflecting that mutation sometimes finds
  a better config)

### Joint kernels

- **P_joint** (6×4): node perception × evolver perception. Evolver
  perception depends on node world state — a `degraded` or `overloaded`
  node directly lowers the pool score. Coupling enters here.
- **D_joint** (4×9): Kronecker product D_node ⊗ D_evolver. Each agent
  decides independently from its own experience (no communication between
  node and evolver within a cycle).
- **A_joint** (9×6): mostly Kronecker product, except when evolver action
  is `mutate` — in that case, A_node's recovery probabilities are boosted.

J = P_joint · D_joint · A_joint → 6×6 world kernel.

---

## Analysis targets

1. **Stationary distribution** of J — how much time the system spends in
   each of the 6 joint states. The key question: does the evolver
   concentrate mass in `healthy × good` or does mutation keep the system
   wandering?

2. **Entropy rate** of J — how predictable the joint system is. Low entropy
   = the system settles; high entropy = the evolutionary search keeps
   churning.

3. **Mean first passage time** from `overloaded × bad` back to
   `healthy × good` — the expected recovery time from the worst joint state.

4. **Trace onto node health** `{healthy, degraded, overloaded}` — the
   effective dynamics visible to an outside observer who can't see whether
   the evolver is in a good or bad strategy. This is what a monitoring
   system would actually see.

5. **IsTraceOf** check: verify that the stationary distribution of the trace
   kernel is the normalized restriction of the joint kernel's stationary
   distribution.

---

## Code structure

```
examples/
  self_healing_nodes/
    main.go        — kernel construction and catrace analysis
    WALKTHROUGH.md — narrative explanation of the model and output
```

`main.go` structure:

1. Define individual state names and build P_node, D_node, A_node
2. Define P_evolver, D_evolver, A_evolver
3. Build P_joint, D_joint (Kronecker), A_joint (with mutate coupling)
4. Compose J = WorldKernel via catrace.Agent
5. Run: Stationary, EntropyRate, Classes, MeanFirstPassage, CommuteTime
6. Build trace kernel onto node health states
7. Verify IsTraceOf
8. Print all results with plain-English interpretation

---

## Open questions before implementation

- Should the node agent represent a single node or a pool? A single node
  is cleaner for the matrix sizes; the pool interpretation requires
  justifying the aggregation.
- The mutation Δ on A_node recovery probability needs a defensible value.
  Options: derive from the workerpool's mutation range on MaxWorkers/Kp,
  or treat it as a free parameter and show sensitivity.
- Three world states for the node (healthy / degraded / overloaded) vs. two
  (healthy / unhealthy). Three maps better to the workerpool's EMA bands
  and gives a richer trace, but doubles the joint state space.
