# Walkthrough: Self-adjusting / Self-healing Network Nodes

The previous example (`validator_repair`) showed how to build a joint kernel from two
independently modeled agents and how coupling enters at specific kernel positions. This
example uses the same construction to answer a design question: given two competing
explanations for where recovery lives, which one does the model support? The structural
novelties are the **variant comparison** pattern — the same joint kernel architecture run
under two parameter regimes with the difference in stationary distribution, MFPT, and
entropy rate reading out the answer — and the `MeanFirstPassage` primitive, introduced
here for the first time.

---

## The story

A network node monitors its own error rate via an exponential moving average (EMA) and
throttles itself when errors climb. An outer evolutionary loop watches pool-level throughput
and mutates the node's configuration when performance drops. The model asks: **which loop
does the real work?** Variant A encodes the claim that throttling alone is a strong recovery
mechanism; Variant B encodes the claim that the evolver's config search is the primary
recovery path. Running both variants through the same analysis lets the numbers vote on
which claim holds.

---

## State spaces

**Node agent:**

| Layer      | States                              | Meaning                                    |
|------------|-------------------------------------|--------------------------------------------|
| World      | `healthy`, `degraded`, `overloaded` | actual error rate of the node              |
| Experience | `ema_low`, `ema_mid`, `ema_high`    | EMA band the node observes                 |
| Action     | `push`, `throttle`, `idle`          | how hard the node drives itself next cycle |

**Evolver agent:**

| Layer      | States                          | Meaning                                               |
|------------|---------------------------------|-------------------------------------------------------|
| World      | `good_strategy`, `bad_strategy` | whether the current MaxWorkers/Kp config is effective |
| Experience | `high_score`, `low_score`       | pool-level throughput×success² the evolver observes   |
| Action     | `promote`, `mutate`             | keep the best known config or try a variation         |

**Joint world states:**

| State | Node         | Evolver           | Meaning                                               |
|-------|--------------|-------------------|-------------------------------------------------------|
| H·G   | `healthy`    | `good_strategy`   | Node fine, evolver has a good config                  |
| H·B   | `healthy`    | `bad_strategy`    | Node fine despite a bad config — throttle carrying it |
| D·G   | `degraded`   | `good_strategy`   | Good config, but node still struggling                |
| D·B   | `degraded`   | `bad_strategy`    | Double trouble — degraded and no good config          |
| O·G   | `overloaded` | `good_strategy`   | Overloaded despite good config — load spike           |
| O·B   | `overloaded` | `bad_strategy`    | Worst joint state — fully degraded, bad config        |

---

## Coupling

- **$P_\text{joint}$ (perception):** When the node is `degraded` or `overloaded`, the pool
  score the evolver observes drops — regardless of actual strategy quality. A sick node
  makes the evolver think the strategy is bad even when it isn't. Concretely,
  P(`high_score` | `good_strategy`, node=`healthy`) = 0.85 falls to 0.55 when
  `node=degraded` and to 0.20 when `node=overloaded`.

- **$D_\text{joint}$ (decision):** Kronecker product — independent decisions. The node reads
  its own EMA; the evolver reads the pool score. Neither knows the other's reading within a
  cycle.

- **$A_\text{joint}$ (action):** When the evolver `mutate`s, the node's recovery
  probabilities are boosted. A better MaxWorkers or Kp setting found by the search helps
  the node heal faster. `promote` leaves node recovery unchanged — it just preserves the
  existing best config. The boost is applied to the `healthy` entry of the node's action
  row, then the row is renormalized so no entry goes negative regardless of boost size.

---

## Two variants

**Variant A — throttle does it all:** `throttle` is a strong recovery action ($P(\text{healthy} \mid \text{throttle}) = 0.65$). The mutation boost is negligible ($+0.02$). The evolver's search barely changes the picture.

**Variant B — evolver matters:** `throttle` is weak ($P(\text{healthy} \mid \text{throttle}) = 0.45$). The mutation boost is large ($+0.25$ for `throttle`). Without the outer search loop the node gets stuck `degraded`.

---

## Math worked by hand

### $P_\text{joint}$ coupling entry

$P_\text{joint}$ maps joint world states to joint experience states ($6 \times 6$). The
construction is:

$$P_\text{joint}[w_n \cdot w_e,\; x_n \cdot x_e]
  = P_\text{node}[w_n, x_n] \times P_\text{evolver}(x_e \mid w_n, w_e)$$

The coupling is visible in the `D·G → ema_low·high_score` entry. A degraded node
suppresses the score the evolver sees, making a good strategy look bad:

```
P_node[degraded → ema_low]                              = 0.20
P_ev(high_score | node=degraded, good_strategy)         = 0.55

P_joint[D·G, ema_low·high_score] = 0.20 × 0.55         = 0.110
```

Compare with the same column from `H·G`:

```
P_node[healthy → ema_low]                               = 0.80
P_ev(high_score | node=healthy, good_strategy)          = 0.85

P_joint[H·G, ema_low·high_score]  = 0.80 × 0.85        = 0.680
```

The entry drops from $0.680$ to $0.110$ — a factor of $6\times$ — purely because the node
is sick. When the node degrades, the evolver loses the signal it needs to distinguish a
good config from a bad one.

### $A_\text{joint}$ renormalization (Variant B, `throttle·mutate`)

When the evolver `mutate`s while the node `throttle`s, the `healthy` entry of the node's
action row is boosted before renormalization. Let $a_0, a_1, a_2$ be the original row and
$b$ be the boost:

$$a'_0 = a_0 + b, \quad
  \hat{a}_i = \frac{a'_i}{\,a'_0 + a_1 + a_2\,}$$

For Variant B, `throttle` row $= (0.45,\; 0.45,\; 0.10)$, boost $b = 0.25$:

$$a'_0 = 0.45 + 0.25 = 0.70, \quad
  \text{sum} = 0.70 + 0.45 + 0.10 = 1.25$$
$$\hat{a} = \left(\frac{0.70}{1.25},\; \frac{0.45}{1.25},\; \frac{0.10}{1.25}\right)
           = (0.56,\; 0.36,\; 0.08)$$

The joint action entry coupling the evolver's `mutate` to a `healthy` next-node-world:

$$A_\text{joint}[\text{thr·mutate} \to \text{H·G}]
  = \hat{a}_0 \times A_\text{evolver}[\text{mutate} \to \text{good strategy}]
  = 0.56 \times 0.60 = 0.336$$

For Variant A the same calculation gives $\hat{a}_0 = 0.657$ and
$A_\text{joint}[\text{thr·mutate} \to \text{H·G}] = 0.657 \times 0.60 = 0.394$.
The Variant A throttle entry is already $0.65$ before the boost; the mutation adds little.
In Variant B, the $+0.25$ boost does the heavy lifting.

---

## Reading the output

### Stationary distribution

```
Variant A — throttle does it all
──────────────────────────────────

Stationary distribution:
  H·G    0.5240
  H·B    0.1813
  D·G    0.1855
  D·B    0.0634
  O·G    0.0341
  O·B    0.0116
```

```
Variant B — evolver matters
─────────────────────────────

Stationary distribution:
  H·G    0.4080
  H·B    0.1534
  D·G    0.2646
  D·B    0.0892
  O·G    0.0636
  O·B    0.0213
```

In A the system spends 52% of its time in the best joint state. In B this falls to 41% and
degraded states collectively absorb 10 percentage points more. The evolver is visibly
carrying load it doesn't carry in A, but the recovery path runs through a random search,
which keeps the system off the best state more often.

### Mean first passage time: O·B → H·G

```
Variant A:
MFPT O·B → H·G (worst→best) = 2.03 steps
MFPT H·G → H·G (self)       = 0.00 steps

Variant B:
MFPT O·B → H·G (worst→best) = 2.56 steps
MFPT H·G → H·G (self)       = 0.00 steps
```

Recovery from the worst joint state takes 2.03 steps on average in A versus 2.56 in B.
When the evolver must find a good config before the node can heal, the climb back from the
worst state is slower. The self-MFPT prints 0 by library convention — the system is already
at the target state.

### Entropy rate

```
Variant A:
Entropy rate = 1.8726 bits/step

Variant B:
Entropy rate = 2.1247 bits/step
```

A system that self-heals locally is more legible — you have a better guess where it will be
next cycle. A system whose recovery depends on a random search loop wanders more
($2.12$ vs $1.87$ bits/step), even when it eventually gets back to `healthy`.

### Trace onto {H·G, O·B}

```
Variant A:
Trace onto {H·G, O·B}:
⎡  0.9783767833932107  0.021623216606789376⎤
  ⎣  0.9740398504617195   0.02596014953828047⎦
IsTraceOf = true
stationary(trace): H·G=0.9783  O·B=0.0217

Variant B:
Trace onto {H·G, O·B}:
⎡ 0.9507461401758306  0.04925385982416944⎤
  ⎣ 0.9445015403967697  0.05549845960323037⎦
IsTraceOf = true
stationary(trace): H·G=0.9504  O·B=0.0496
```

An outside observer who can only distinguish "everything fine" from "everything failing"
sees the effective transition probabilities between those two states with all intermediate
states summed out. In B the system is more than twice as likely to tip from peak to worst
in any given cycle ($4.9\%$ vs $2.2\%$). `IsTraceOf = true` confirms the trace kernel is
the exact mathematical projection of the joint kernel — not an approximation.

### Empirical trace

```
Variant A:
Empirical trace (500-step sample):
⎡    0.9774400192479754    0.022559980752024422⎤
  ⎣    0.9998333888703766  0.00016661112962345885⎦

Variant B:
Empirical trace (500-step sample):
⎡   0.9504463923748434    0.04955360762515652⎤
  ⎣   0.9999091074350117  9.089256498818396e-05⎦
```

The `O·B → H·G` row is near $1.0$ in both variants because `O·B` is rare — in 500 steps
the chain visits it few enough times that nearly every observed visit is followed immediately
by recovery. Increasing the trajectory length to 10,000 will bring the empirical estimate
close to the analytical $0.974$ (Variant A) and $0.945$ (Variant B).

---

## What you can change

- **Strengthen Variant B's throttle:** In the `nodeA` for Variant B, raise the `throttle`
  row's `healthy` entry from $0.45$ toward $0.65$ (Variant A's value). Watch
  `MFPT O·B → H·G` fall and the `H·G` stationary weight rise. If both converge on Variant
  A's numbers, the throttle has absorbed the evolver's load and the variant distinction
  disappears.

- **Zero the mutation boost:** Set `mutationBoost` to `{0, 0, 0}` in Variant B.
  $A_\text{joint}$ now has no coupling when the evolver `mutate`s — the two agents' action
  effects are fully independent. Observe how the stationary distribution and MFPT change;
  this isolates the evolver's contribution to node recovery from the node's own throttle.

- **Remove the perception coupling:** Change `evolverPCoupled` so that the evolver's score
  is independent of node health (set all three node rows to the `node=healthy` values:
  $\{0.85, 0.15\},\, \{0.20, 0.80\}$). The evolver can no longer be confused by a sick
  node. Watch how the `D·G → D·B` leakage rate changes — the `good_strategy` states should
  become stickier.

- **Trace onto a different subset:** Try `[]int{0, 2, 4}` (all `good_strategy` states) or
  `[]int{1, 3, 5}` (all `bad_strategy` states). The stationary theorem holds for any subset
  whose complement contains no recurrent states — here the full chain is ergodic so any
  nonempty proper subset works.

- **Extend the Monte Carlo trajectory:** Change `500` to `10000`. Observe the empirical
  `O·B → H·G` entry converge from $\approx 1.000$ toward the analytical value as the chain
  accumulates enough visits to `O·B` to estimate the transition reliably.
