# Plan: Network of Self-healing Nodes (`examples/network_of_healers`)

## What this example introduces

`self_healing_nodes` asked which of two loops — the node's own throttle or the
evolver's config search — does the real recovery work. This example asks the next
question: **how does the topology of connections between nodes change the system's
ability to recover?**

Four configurations of two self-healing nodes plus one shared evolver are run as
variants and compared head-to-head. The structural novelty is that coupling between
nodes is topology-specific, encoded entirely in $A_\text{joint}$: the same node and
evolver kernels are reused across all four topologies, and only the inter-node action
effects change.

This is a **planned** example (not yet a README scenario). As of 2026-08-04,
`README.md` scenario 5 is the three-agent majority-valid coordination network
(`docs/patterns/story-supervisor.md`, not implemented). Do not claim a README
scenario number for this plan until the modelling premise is settled and the
README scenario list is deliberately updated. See wiki `[[Scenario Registry]]`.

---

## State space

### Per-node model (reused from `self_healing_nodes`)

Each node is the same 3-state agent:

| Layer      | States                              |
|------------|-------------------------------------|
| World      | `healthy (H)`, `degraded (D)`, `overloaded (O)` |
| Experience | `ema_low`, `ema_mid`, `ema_high`    |
| Action     | `push`, `throttle`, `idle`          |

Use Variant B parameters from `self_healing_nodes` as the baseline — the harder
regime where self-healing alone is unreliable. This makes topology differences more
visible.

### Shared evolver

One evolver monitors the aggregate pool score across both nodes:

| Layer      | States                          |
|------------|---------------------------------|
| World      | `good_strategy`, `bad_strategy` |
| Experience | `high_score`, `low_score`       |
| Action     | `promote`, `mutate`             |

Pool score now degrades as a function of **how many** nodes are sick, not just one.

### Joint world space (2 nodes + evolver)

Joint world state = `(node1_health · node2_health · evolver_strategy)`.
Ordered row-major: node1 outermost, evolver innermost.

| Index | State  | Node 1     | Node 2     | Evolver           |
|-------|--------|------------|------------|-------------------|
| 0     | HH·G   | healthy    | healthy    | good_strategy     |
| 1     | HH·B   | healthy    | healthy    | bad_strategy      |
| 2     | HD·G   | healthy    | degraded   | good_strategy     |
| 3     | HD·B   | healthy    | degraded   | bad_strategy      |
| 4     | HO·G   | healthy    | overloaded | good_strategy     |
| 5     | HO·B   | healthy    | overloaded | bad_strategy      |
| 6     | DH·G   | degraded   | healthy    | good_strategy     |
| 7     | DH·B   | degraded   | healthy    | bad_strategy      |
| 8     | DD·G   | degraded   | degraded   | good_strategy     |
| 9     | DD·B   | degraded   | degraded   | bad_strategy      |
| 10    | DO·G   | degraded   | overloaded | good_strategy     |
| 11    | DO·B   | degraded   | overloaded | bad_strategy      |
| 12    | OH·G   | overloaded | healthy    | good_strategy     |
| 13    | OH·B   | overloaded | healthy    | bad_strategy      |
| 14    | OD·G   | overloaded | degraded   | good_strategy     |
| 15    | OD·B   | overloaded | degraded   | bad_strategy      |
| 16    | OO·G   | overloaded | overloaded | good_strategy     |
| 17    | OO·B   | overloaded | overloaded | bad_strategy      |

18 joint world states. 9 joint experience states (3×3 node experiences × 2 evolver
experiences = 18, but evolver experience is derived from world + node healths so also
18). Joint action space = 3×3×2 = 18.

### Extending to 3 nodes

With 3 nodes + evolver: $3^3 \times 2 = 54$ world states, $3^3 \times 2 = 54$
experience states, $3^3 \times 2 = 54$ action states. The joint kernel J is 54×54.
This is tractable for catrace's matrix operations but the coupling loops become
4-deep. **Treat 3-node extension as a stretch goal** — implement 2-node first and
verify the topology comparison is legible, then decide whether 3-node adds enough.

---

## Kernel construction

### $P_\text{joint}$ (18×18)

The evolver's perceived pool score now degrades with the **count** of sick nodes.
Define a severity score $s = \text{num\_degraded} + 2 \times \text{num\_overloaded}$
(overloaded counts double because it is the worse state). Scale the evolver's
`high_score` probability linearly down from its baseline:

```
evolverHighScore(s) = max(baserate - s × penaltyPerSick, 0)
```

where `baserate = 0.85` (both healthy, good strategy) and `penaltyPerSick = 0.20`.

So:
- s=0 (HH): P(high_score | good) = 0.85
- s=1 (one degraded): 0.65
- s=2 (two degraded or one overloaded): 0.45
- s=3 (one degraded + one overloaded): 0.25
- s=4 (two overloaded): 0.05

Construction loop (pseudo-code):
```
for n1w, n2w, ew in joint_world_states:
    s = severity(n1w, n2w)
    for n1x, n2x, ex in joint_experience_states:
        P_joint[n1w·n2w·ew, n1x·n2x·ex] =
            nodeP[n1w, n1x] * nodeP[n2w, n2x] *
            evolverPCoupled(s, ew, ex)
```

### $D_\text{joint}$ (18×18)

Kronecker product — all three agents decide independently from their own experience:

```
D_joint[n1x·n2x·ex, n1g·n2g·eg] =
    nodeD[n1x, n1g] * nodeD[n2x, n2g] * evolverD[ex, eg]
```

### $A_\text{joint}$ (18×18) — topology-specific

This is where the four topologies diverge.

**Baseline (independent):** pure Kronecker product. No inter-node coupling.

```
A_joint[n1g·n2g·eg, n1w'·n2w'·ew'] =
    nodeA_eff(n1g, eg)[n1w'] * nodeA_eff(n2g, eg)[n2w'] * evolverA[eg, ew']
```

where `nodeA_eff(ng, eg)` applies the mutation boost if `eg=mutate`, same as
`self_healing_nodes`.

**Linear chain:** node 1 feeds into node 2. When node 1 is overloaded, its excess
load spills to node 2, increasing node 2's probability of degrading. Node 2 has no
effect on node 1.

Modify node 2's effective action row based on node 1's **current action** (the action
just chosen, not the next world state):

```
spillBoost = 0 if n1g ∈ {push, throttle} and node1 not overloaded
           = 0.15 if n1g = push and node1 = overloaded (spill load)
           = 0.05 if n1g = throttle and node1 = overloaded

node2A_eff[ng2][degraded] += spillBoost
node2A_eff[ng2][overloaded] += spillBoost × 0.5
renormalize node2A_eff row
```

More precisely, spill is conditioned on node 1's **world state** (since that is what
determines whether it is actually shedding load):

```
for each n1g, n2g, eg:
    if node1_is_overloaded(n1g, current_n1w):  # n1w is the from-world-state
        node2_row = spillAdjust(nodeA_eff(n2g, eg), spillRate)
    else:
        node2_row = nodeA_eff(n2g, eg)
    A_joint[n1g·n2g·eg, n1w'·n2w'·ew'] =
        nodeA_eff(n1g, eg)[n1w'] * node2_row[n2w'] * evolverA[eg, ew']
```

Note: the from-world-state is encoded in the A_joint row index. A_joint rows are
joint actions `(n1g, n2g, eg)`, but we need to condition spill on from-world-state
too. This means the loop must iterate over from-world-states as well — A_joint
becomes a function of `(from_world, joint_action)` in the coupling case. Encode
this by building a separate effective A matrix for each from-world partition.

**Alternative simpler encoding:** Treat the spill as a fixed fractional boost to
node 2's `degraded` entry whenever node 1's **action** is `push` (aggressive action
implies node 1 is pushing hard, which transfers load). This avoids conditioning on
world state and keeps the construction loop clean at the cost of some realism:

```
spillRate = 0.10   # push action from node1 spills 10pp onto node2's degraded entry
if n1g = push:
    node2_row = boost(nodeA_eff(n2g, eg), degraded, spillRate)
else:
    node2_row = nodeA_eff(n2g, eg)
```

**Recommendation:** use the action-conditioned version (simpler) for the initial
implementation and note in the walkthrough why it's a simplification.

**Parallel redundant:** nodes share a load pool. When one node degrades, the other
absorbs its work, increasing its own overload probability. Symmetric: each node's
action affects the other's effective A row.

```
mutualSpill = 0.08   # symmetric — each node spills onto the other

for n1g, n2g, eg:
    node1_row = boost(nodeA_eff(n1g, eg), overloaded, mutualSpill if n2g=push else 0)
    node2_row = boost(nodeA_eff(n2g, eg), overloaded, mutualSpill if n1g=push else 0)
    renormalize both rows
    A_joint[n1g·n2g·eg, n1w'·n2w'·ew'] =
        node1_row[n1w'] * node2_row[n2w'] * evolverA[eg, ew']
```

**Ring:** same as parallel redundant for 2 nodes (a 2-node ring is bidirectional by
definition). With 3 nodes a ring differs from parallel — include this topology in the
3-node extension, not the 2-node implementation. For 2 nodes, **merge ring with
parallel redundant** or drop ring and use 3 topologies only.

**Recommended final topology set (2-node):**

| Variant | Name           | Inter-node coupling in $A_\text{joint}$              |
|---------|----------------|------------------------------------------------------|
| 1       | Independent    | None — pure Kronecker                                |
| 2       | Linear chain   | node1 push → boosts node2 degraded/overloaded prob   |
| 3       | Parallel load  | each node push → boosts other's overloaded prob      |
| 4       | Protective pair| each node throttle → reduces other's overloaded prob |

Variant 4 (protective pair) is the optimistic dual of variant 3: nodes that throttle
actively shed load to neighbors, reducing their stress. This gives a nice 2×2 grid:
independent × load-transfer direction × whether transfer is harmful or helpful.

---

## Analysis to run per topology

```go
// 1. Stationary distribution
pi, _ := J.Stationary(1e-12, 5000)

// 2. MFPT: worst joint state → best joint state
//    worst = OO·B (index 17), best = HH·G (index 0)
mfptWorstBest, _ := J.MeanFirstPassage(17, 0)

// 3. MFPT: any-degraded → all-healthy
//    Compute for each degraded starting state, report min/max/mean
degradedStates := []int{2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17}
for _, s := range degradedStates:
    mfpt, _ := J.MeanFirstPassage(s, 0)

// 4. Entropy rate
h, _ := J.EntropyRate(2)

// 5. Trace onto {HH·G=0, OO·B=17}
trace, _ := J.Trace([]int{0, 17}, 1e-12)
ok, _   := trace.IsTraceOf(J, []int{0, 17}, 1e-12)
piTrace, _ := trace.Stationary(1e-12, 5000)

// 6. Coarse trace onto {all-healthy, any-degraded, majority-failed}
//    all-healthy  = {0,1}           (HH·G, HH·B)
//    majority-failed = {8..17}      (DD, DO, OD, OO — both nodes sick)
coarseSubset := []int{0, 1, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
// Or three-state: all-healthy={0,1}, one-sick={2..7}, both-sick={8..17}
```

**Comparative output table** (printed for all 4 topologies):

```
Topology           | π(HH·G)  | π(OO·B)  | MFPT OO·B→HH·G | H (bits/step)
Independent        |  0.XXXX  |  0.XXXX  |   X.XX         |  X.XXXX
Linear chain       |  0.XXXX  |  0.XXXX  |   X.XX         |  X.XXXX
Parallel load      |  0.XXXX  |  0.XXXX  |   X.XX         |  X.XXXX
Protective pair    |  0.XXXX  |  0.XXXX  |   X.XX         |  X.XXXX
```

---

## Code structure

```go
type topology struct {
    name string
    // interNodeSpill: how much n1's push action boosts n2's degraded prob (linear)
    // interNodeSpillSymmetric: whether n2 also affects n1 (parallel/protective)
    // spillSign: +1 for harmful (load transfer), -1 for protective (load shedding)
    n1ToN2Spill float64
    n2ToN1Spill float64
    spillTriggerAction int // 0=push, 1=throttle, 2=idle
    spillTargetState   int // 0=healthy, 1=degraded, 2=overloaded
}

func runTopology(t topology) results { ... }

func main() {
    topologies := []topology{
        {name: "Independent",     n1ToN2Spill: 0,    n2ToN1Spill: 0   },
        {name: "Linear chain",    n1ToN2Spill: 0.10, n2ToN1Spill: 0,
         spillTriggerAction: push, spillTargetState: degraded},
        {name: "Parallel load",   n1ToN2Spill: 0.08, n2ToN1Spill: 0.08,
         spillTriggerAction: push, spillTargetState: overloaded},
        {name: "Protective pair", n1ToN2Spill: -0.08, n2ToN1Spill: -0.08,
         spillTriggerAction: throttle, spillTargetState: overloaded},
    }
    for _, t := range topologies {
        runTopology(t)
    }
    printComparisonTable(results)
}
```

Kernel construction helpers:

```go
// effectiveNodeRow returns the node's A row (3 entries) after applying
// mutation boost (if evolver mutated) and inter-node spill adjustment,
// then renormalizing.
func effectiveNodeRow(baseA [3][3]float64, ng int, eg int,
    mutationBoost [3]float64, spillDelta float64, spillTarget int) [3]float64
```

---

## WALKTHROUGH.md structure (per conventions)

**§1 Opening paragraph:** Positions this after `self_healing_nodes`. That example
asked which internal loop heals; this asks whether and how the wiring between nodes
changes system resilience. Structural novelty: topology as a model parameter, not a
code architecture choice.

**§2 The story:** Two sentences. Two nodes with the same self-healing mechanism,
four ways to wire them. Which wiring keeps the system healthiest?

**§3 State spaces:** Per-agent tables for Node 1, Node 2 (same), Evolver. Joint
world states table (18 rows). Note the systematic naming convention.

**§4 Coupling:** Three bullets — P_joint (aggregate severity score), D_joint
(Kronecker), A_joint (topology-specific spill, one sub-bullet per topology variant).

**§5 Math by hand:** Show the linear-chain $A_\text{joint}$ entry for the
`push·push·mutate → HD·G` transition:

```
node1A_eff[push] = {0.55, 0.35, 0.10}   (Variant B baseline)
spill onto node2 degraded: +0.10
node2A_eff[push] = {0.55, 0.35+0.10, 0.10} = {0.55, 0.45, 0.10} → renorm by 1.10
               = {0.50, 0.409, 0.091}

evolverA[mutate → good] = 0.60

A_joint[push·push·mutate → HD·G] = 0.55 × 0.409 × 0.60 = 0.135
```

Compare with independent topology (no spill):

```
A_joint[push·push·mutate → HD·G] = 0.55 × 0.35 × 0.60 = 0.116
```

The spill raises the probability of node 2 degrading by ~16% relative.

**§6 Reading the output:** Show raw terminal output for all 4 topologies — stationary
distributions, MFPT table, entropy rates, comparative summary table. Two sentences of
interpretation per result.

**§7 What you can change:** Five experiments:
1. Increase spill rate — watch MFPT grow for load topologies
2. Make protective pair's boost stronger — does it beat independent?
3. Try asymmetric spill (n1→n2 strong, n2→n1 weak) — directional dependency
4. Extend to 3 nodes — add a third node agent and observe state space growth
5. Change evolver severity penalty — make evolver blind to node count

---

## README entry (draft — assign number only when adding to README)

**Story paragraph:** Plain English, no backticks. Four ways to wire two self-healing
nodes together: unconnected, in series, sharing load symmetrically, or helping each
other throttle. The wiring is invisible to the nodes themselves — each still reads
only its own EMA and decides independently. The coupling lives entirely in how one
node's action changes the other node's world.

**State meanings:** Node 1/2 world/experience/action (same as self_healing_nodes),
Evolver, Joint world states (18 states with naming convention).

**Interpretation bullets:**
- $P_\text{joint}$ couples evolver perception to aggregate node health (severity score)
- $D_\text{joint} = D_1 \otimes D_2 \otimes D_\text{evolver}$ — nodes decide independently
- $A_\text{joint}$ encodes topology: independent (Kronecker), load-transfer (spill boost),
  protective (spill reduction)
- Running all four topologies as variants isolates the topological effect from the
  per-node self-healing effect

**Code line:** `examples/network_of_healers/main.go`

**Played-out version:** Three paths, all using linear-chain topology:
- Version A: HH·G — both healthy, node 1 pushes, load spills but node 2 handles it
- Version B: HO·G — node 2 already overloaded, node 1's push tips it toward OO·B
- Version C: OO·B → HH·G — both throttle and idle, evolver mutates, recovery in ~3 steps

---

## Open questions / risks

| Risk | Mitigation |
|------|-----------|
| State space grows to 54 for 3 nodes — matrix ops may be slow | Benchmark `J.Stationary` on 54×54; catrace uses gonum dense, should be fine |
| Spill conditioned on joint action only (not from-world) is a simplification | Note clearly in walkthrough; offer world-conditioned version in "what you can change" |
| Protective pair may need careful parameter tuning to show meaningful improvement over independent | Pre-run mentally: if spill reduces overloaded by 0.08, node's idle baseline is 0.05 — boost can't push below 0; renorm handles it |
| 18-state chain may make trace onto 3 coarse buckets computationally noisier | Use `IsTraceOf` check; if numerical issues arise, use 2-state trace {HH·G, OO·B} only |
| Naming 18 states in the joint tables risks reader overload in WALKTHROUGH | Use abbreviated table, refer reader to code comments for full enumeration |
