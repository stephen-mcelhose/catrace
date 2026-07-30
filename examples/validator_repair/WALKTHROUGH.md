# Walkthrough: Two-Agent Validator / Repair Pair

This example fills the gap between the first two examples. The first example modeled a
single agent with P, D, A kernels. The second applied trace to a joint kernel given
directly. This example shows how to **build** a joint kernel from two independently
modeled agents — and why it matters how you build it.

## The story

One agent (Worker) performs tasks. A second agent (Validator) monitors the Worker and,
when things go wrong, attempts to repair it. Each may be in a valid or invalid mode
independently.

We model each agent separately using P, D, A kernel triplets. We then define joint kernels
over the product of their state spaces and compose them exactly as we would for a single
agent:

```
J = P_joint · D_joint · A_joint
```

Coupling between the agents enters through specific kernels with clear physical meaning —
not through a post-hoc blending parameter.

---

## State spaces

### Individual state spaces

**Worker:**

| Space      | States                             |
|------------|------------------------------------|
| World W₁   | `worker_valid`, `worker_invalid`   |
| Experience X₁ | `sees_ok`, `sees_problem`       |
| Actions G₁ | `produce`, `self_check`, `idle`    |

**Validator:**

| Space      | States                               |
|------------|--------------------------------------|
| World W₂   | `validator_valid`, `validator_invalid` |
| Experience X₂ | `looks_good`, `looks_bad`         |
| Actions G₂ | `validate`, `repair`, `idle`         |

### Product state spaces

| Space        | Size | Elements                                          |
|--------------|------|---------------------------------------------------|
| W = W₁ × W₂ | 4    | VV, VI, IV, II                                    |
| X = X₁ × X₂ | 4    | ok·good, ok·bad, prob·good, prob·bad              |
| G = G₁ × G₂ | 9    | all 9 worker-action × validator-action pairs      |

State indexing is row-major: pair (i₁, i₂) maps to index i₁·n₂ + i₂.
So VV=0, VI=1, IV=2, II=3.

---

## Individual agents

Each agent is modeled and its world kernel computed independently before any coupling is
introduced. This lets you reason about each agent's intrinsic dynamics.

**Worker world kernel W₁ = P₁·D₁·A₁:**

Running the example prints this 2×2 matrix. Read it as: given the worker is currently
valid (or invalid), where does it tend to be on the next cycle, operating entirely alone?
The Worker is self-stabilizing when valid (W₁[valid,valid] ≈ 0.64) but tends to
drift toward invalid when it idles.

**Validator world kernel W₂ = P₂·D₂·A₂:**

The Validator is more self-stable — validating keeps it calibrated (W₂[valid,valid] ≈
0.71). Repairing degrades it slightly because repair is taxing.

These two 2×2 matrices represent the agents in isolation. The next step is coupling.

---

## Building the joint kernel

### Why not just blend the individual world kernels?

A naive approach would be to compute W₁ and W₂ separately, form their Kronecker product
W₁⊗W₂ (the uncoupled joint dynamics), and then blend with a hand-crafted coupling matrix
C to add cross-agent effects:

```
J = α·(W₁⊗W₂) + (1-α)·C          ← do NOT do this
```

This is wrong for a subtle reason: the coupling matrix C has no P, D, A decomposition.
It is an ad-hoc stochastic matrix. If you change D₂ (how the Validator decides), it has
zero effect on J — the Validator's decision policy and J are completely disconnected.

The correct approach keeps the coupling inside the P, D, A framework.

### P_joint: coupled perception (4×4, W→X)

Worker's perception is unchanged by Validator's state. Worker sees what it sees.

Validator's perception, however, depends on both its own internal state AND the Worker's
world state. When the Worker is degraded, the Validator can observe Worker errors in
addition to its own internal signals — it is more likely to report `looks_bad`.

Formally, the Validator has a modified perception P₂_coupled[(w₁,w₂), x₂]:

| w₁     | w₂     | P₂_coupled looks_good | P₂_coupled looks_bad |
|--------|--------|------------------------|----------------------|
| valid  | valid  | 0.85                   | 0.15                 |
| valid  | invalid| 0.40                   | 0.60                 |
| invalid| valid  | 0.60                   | 0.40                 |
| invalid| invalid| 0.25                   | 0.75                 |

The bottom two rows are shifted from the top two: Worker degradation pushes Validator's
perception toward `looks_bad` even when Validator's own state is fine.

The full P_joint entry is:

```
P_joint[(w₁,w₂), (x₁,x₂)] = P₁[w₁, x₁] × P₂_coupled[(w₁,w₂), x₂]
```

For example, state VV (w₁=valid, w₂=valid):

```
ok·good   = P₁[valid, sees_ok]    × P₂[valid, looks_good] = 0.90 × 0.85 = 0.765
ok·bad    = P₁[valid, sees_ok]    × P₂[valid, looks_bad]  = 0.90 × 0.15 = 0.135
prob·good = P₁[valid, sees_prob]  × P₂[valid, looks_good] = 0.10 × 0.85 = 0.085
prob·bad  = P₁[valid, sees_prob]  × P₂[valid, looks_bad]  = 0.10 × 0.15 = 0.015
                                                                     sum = 1.000 ✓
```

State IV (w₁=invalid, w₂=valid) — note the shifted Validator perception:

```
ok·good   = 0.30 × 0.60 = 0.180
ok·bad    = 0.30 × 0.40 = 0.120
prob·good = 0.70 × 0.60 = 0.420
prob·bad  = 0.70 × 0.40 = 0.280    sum = 1.000 ✓
```

### D_joint: independent decisions (4×9, X→G)

The joint decision kernel is the Kronecker product D₁⊗D₂. This encodes the assumption
that **agents do not communicate before deciding** — each agent chooses its action based
solely on its own experience.

```
D_joint[(x₁,x₂), (g₁,g₂)] = D₁[x₁, g₁] × D₂[x₂, g₂]
```

For joint experience ok·good (Worker sees ok, Validator looks good):

```
D₁[sees_ok, :]    = [0.80, 0.15, 0.05]   (produce, self_check, idle)
D₂[looks_good, :] = [0.60, 0.10, 0.30]   (validate, repair, idle)
```

Joint action probabilities (row-major, worker first):

```
produce|validate   = 0.80 × 0.60 = 0.480
produce|repair     = 0.80 × 0.10 = 0.080
produce|idle       = 0.80 × 0.30 = 0.240
self_check|validate= 0.15 × 0.60 = 0.090
...                                         sum = 1.000 ✓
```

The Kronecker product of two row-stochastic matrices is row-stochastic — no normalization
needed.

In the code this is a simple nested loop, not a hand-coded matrix:

```go
dJoint.Set(row, col, workerD.At(x1, g1)*validatorD.At(x2, g2))
```

### A_joint: coupled action effects (9×4, G→W)

This is the second coupling point and the more consequential one. When the Validator
repairs, it does not merely restore its own state — **it restores the Worker's world state
too**. Validator repair raises Worker's probability of being valid by approximately 0.20
above what it would achieve independently.

For non-repair actions (g₂ ≠ repair), the effects are independent:

```
A_joint[(g₁,g₂), (w₁,w₂)] = A₁[g₁, w₁] × A₂[g₂, w₂]
```

For repair actions (g₂ = repair), Worker gets a boost:

| Worker action | Independent A₁[g₁, valid] | Boosted A₁[g₁, valid] |
|---------------|---------------------------|------------------------|
| produce       | 0.70                      | 0.90                   |
| self_check    | 0.50                      | 0.70                   |
| idle          | 0.40                      | 0.60                   |

Example: joint action (produce|repair):

```
Independent:  VV = 0.70 × 0.60 = 0.420  (Worker produces, Validator repairs itself)
Coupled:      VV = 0.90 × 0.60 = 0.540  (Validator repair also fixes Worker)
```

The repair boost is visible in the VV entry: 0.540 vs 0.420 independent. This is where
the functional relationship between agents lives.

### Composition: J = P_joint · D_joint · A_joint

The joint Agent struct holds the three joint kernels. Its `WorldKernel()` method computes
J = P·D·A over the product spaces:

```
P_joint (4×4)  ·  D_joint (4×9)  ·  A_joint (9×4)  =  J (4×4)
```

The dimension checks that catrace's `Agent.Validate()` performs:
- D.rows (4) == P.cols (4) ✓   — joint experience space consistent
- D.cols (9) == A.rows (9) ✓   — joint action space consistent
- A.cols (4) == P.rows (4) ✓   — joint world space consistent

---

## Reading the joint kernel J

The output (example values — yours will match):

```
      VV        VI        IV        II
VV  [ 0.475    0.202     0.233     0.090 ]
VI  [ 0.484    0.246     0.188     0.082 ]
IV  [ 0.412    0.196     0.276     0.116 ]
II  [ 0.422    0.227     0.238     0.113 ]
```

Reading the rows in plain English:

- **VV → VV (0.475):** When both are reliable, there is a nearly 50/50 chance they stay
  jointly reliable next cycle. The system is not especially sticky at VV because the
  Worker can degrade on any cycle.
- **VV → IV (0.233):** Worker degrades with non-trivial probability even when Validator is
  fine. Validator's perception catches this, but perception leads to decision which takes
  a full cycle before repair can help.
- **IV → VV (0.412):** The Validator is still reliable and can repair the Worker. Recovery
  to VV from IV is relatively fast — 41% per step.
- **II → VV (0.422):** Even from full degradation, recovery to VV happens with reasonable
  probability per step. This is higher than you might expect from a degraded validator
  because the repair boost in A_joint helps.

**Stationary distribution:**

```
VV ≈ 0.457,  VI ≈ 0.212,  IV ≈ 0.234,  II ≈ 0.097
```

The system spends about 46% of the time with both agents reliable, and only 10% fully
degraded. The single recurrent class confirms the chain is ergodic — the system always
recovers eventually.

---

## Trace onto {VV, II}

We trace J onto the two extreme states — fully reliable and fully degraded — discarding
the mixed states VI and IV. This is the coarsest operational picture: healthy or failed.

Subset indices: VV=0, II=3. Complement B = {VI=1, IV=2}.

The trace formula:

```math
\text{Tr}(J) = J_{SS} + J_{SB}(I - J_{BB})^{-1}J_{BS}
```

where S = {VV,II} and B = {VI,IV}. The (I - J_BB)^{-1} term sums over all possible
detours through the mixed states — how many times the system visits VI or IV before
returning to VV or II.

**Output:**

```
IsTraceOf = true
stationary(J)|{VV,II} norm = 0.825220  0.174780
stationary(trace)          = 0.825220  0.174780
```

`IsTraceOf = true` confirms the library's verification: the trace kernel computed from the
formula matches (within tolerance) what you would observe if you ran the parent chain and
recorded only the {VV,II} states.

The stationary distributions match exactly — the long-run frequencies you observe in the
coarse view are consistent with the parent's stationary distribution restricted and
normalized to {VV,II}.

---

## Empirical trace

The Monte Carlo section samples 200 steps from J, filters to {VV,II}, and estimates the
trace kernel from those observations. With 200 steps the estimate is rough but in the
right ballpark. Increasing the trajectory length to 10,000 will bring the empirical
estimate close to the analytical trace kernel.

---

## What you can change

- **Increase repair boost:** Change `workerABoosted` values toward 1.0. Observe that
  stationary weight on VV increases and time-in-II decreases.
- **Remove coupling in P_joint:** Set rows IV and II of P_joint to use the uncoupled P₂
  values (0.85/0.15 and 0.40/0.60). Validator can no longer observe Worker degradation.
  Watch how the stationary distribution changes.
- **Add communication in D_joint:** Replace D_joint with a non-Kronecker matrix where
  joint experience (prob·good) — Worker sees a problem but Validator looks fine — shifts
  Validator toward repair even without explicit bad signal. This models information
  sharing between agents at the decision step.
- **Trace onto a different subset:** Try `[]int{0, 1, 2}` (all states except II) or
  `[]int{1, 2}` (the mixed states only). The stationary theorem holds for any subset whose
  complement has no recurrent states.
