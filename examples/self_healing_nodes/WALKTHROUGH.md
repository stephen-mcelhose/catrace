# Walkthrough: Self-adjusting / Self-healing Network Nodes

Inspired by [caseylmanus/workerpool](https://github.com/caseylmanus/workerpool).

---

## The story

A network node monitors its own error rate via an exponential moving average
(EMA) and throttles itself when errors climb — the same proportional
backpressure the workerpool applies via its PID controller. An outer
evolutionary loop watches pool-level throughput and mutates the node's
configuration when performance drops — the same elitist mutation cycle the
workerpool uses to search strategy space.

The model asks: **which loop does the real work?**

---

## Two agents

**Node agent** — the fast inner loop.
Perceives its own health through its EMA error band, decides how hard to push,
and heals itself by throttling. Runs every fetch cycle.

**Evolver agent** — the slow outer loop.
Perceives pool-level score (throughput × success²), decides to promote the
best known config or mutate toward a new one. Runs every generation window
(5 seconds in the workerpool). When mutation finds a better config, the node's
recovery probabilities improve.

---

## Joint world states

Each state is a pair: **(node health · evolver strategy)**.

| State | Node       | Evolver        | Meaning                                        |
|-------|------------|----------------|------------------------------------------------|
| H·G   | healthy    | good strategy  | Node fine, evolver has a good config           |
| H·B   | healthy    | bad strategy   | Node fine despite a bad config — throttle carrying it |
| D·G   | degraded   | good strategy  | Good config, but node still struggling         |
| D·B   | degraded   | bad strategy   | Double trouble — degraded and no good config   |
| O·G   | overloaded | good strategy  | Overloaded despite good config — load spike    |
| O·B   | overloaded | bad strategy   | Worst joint state — fully degraded, bad config |

---

## Coupling

The two agents are coupled at two points:

**P_joint (perception):** When the node is degraded or overloaded, the pool
score the evolver observes drops — regardless of actual strategy quality. A
sick node makes the evolver think the strategy is bad even when it isn't.

**A_joint (action):** When the evolver mutates, the node's recovery
probabilities are boosted. A better MaxWorkers or Kp setting found by the
search helps the node heal faster. Promote leaves node recovery unchanged —
it just preserves the existing best config.

**D_joint (decision):** Kronecker product — independent decisions. The node
reads its own EMA; the evolver reads the pool score. Neither knows the
other's reading within a cycle.

---

## Two variants

The same model is run with two different claims about where healing lives.

### Variant A — throttle does it all

The node's throttle action is a strong recovery mechanism on its own. The
mutation boost is negligible. The evolver's search barely changes the picture.

### Variant B — evolver matters

The node's throttle is weak — slowing down alone doesn't reliably fix the
problem. The mutation boost is large — finding a better config is the
primary recovery path.

---

## Reading the output

### Stationary distribution

Where the system spends its time in the long run.

| | Variant A | Variant B |
|--|-----------|-----------|
| H·G (best)      | 52.4% | 40.8% |
| H·B             | 18.1% | 15.3% |
| D·G             | 18.6% | 26.5% |
| D·B             |  6.3% |  8.9% |
| O·G             |  3.4% |  6.4% |
| O·B (worst)     |  1.2% |  2.1% |
| **Degraded total** | **24.9%** | **35.4%** |

In B, the system spends 10 percentage points more time degraded. The evolver
is visibly carrying load it doesn't carry in A.

### Mean first passage time: O·B → H·G

How many steps on average to recover from the worst state to the best.

| Variant A | Variant B |
|-----------|-----------|
| 2.03 steps | 2.56 steps |

Recovery is slower in B — when the evolver must find a good config before the
node can heal, it takes longer to climb back from the worst state.

### Entropy rate

How unpredictable the system's next state is, given where it is now.
A fair coin is 1.0 bit/step. Higher means more wandering.

| Variant A | Variant B |
|-----------|-----------|
| 1.87 bits/step | 2.12 bits/step |

A system that self-heals locally is more *legible* — you have a better guess
where it will be next cycle. A system that depends on a random search loop to
recover wanders more, even when it eventually gets back to health.

### Trace onto {H·G, O·B}

An outside observer who can only tell "everything fine" from "everything
failing" — what are the effective transition probabilities between those two
states, with all intermediate states summed out?

| | H·G → O·B | O·B → H·G |
|--|-----------|-----------|
| Variant A | 2.2% | 97.4% |
| Variant B | 4.9% | 94.5% |

In B, the system is more than twice as likely to tip from peak to worst in any
given cycle. The `IsTraceOf` check confirms the trace kernel is the exact
mathematical projection of the joint kernel — not an approximation.
