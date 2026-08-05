# Plan: Network of Self-healing Nodes (`examples/network_of_healers`)

## What this example introduces

`self_healing_nodes` studied a **single** node plus an outer config loop. This
example drops the shared evolver and asks a different question:

> Several identical self-healing nodes sit on a real-ish dependency graph.
> Each node's world is pushed around by its **in-neighbors and out-neighbors**.
> The true joint story is large; the story we actually tell (ops dashboard,
> SLO, "is the pool OK?") is a **collapsed** one. What does that collapse
> preserve, and what does it hide?

Structural novelties relative to existing examples:

1. **Graph-shaped coupling** — edges mean load / backpressure into a node's
   world dynamics, not a free ±spill parameter dressed up as "topology."
2. **Collapsing the story** — full joint \(J\) is built, then `Trace` reduces
   it to coarse operator-visible buckets; `IsTraceOf` confirms the collapse.
3. **Optional second act** — estimate the collapsed kernel from a long joint
   trajectory observed only through those coarse labels
   (`EstimateKernelFromSequence`), and compare to the exact Trace.

This does **not** need to re-litigate throttle-vs-evolver. Reuse one per-node
parameter set (default: Variant B from `self_healing_nodes` — weak local heal,
so neighbor stress is visible). If Variant A vs B later makes a nice contrast
under the same graph, that is optional spice, not the thesis.

This is a **planned** example (not yet a README scenario). As of 2026-08-04,
`README.md` scenario 5 is the three-agent majority-valid coordination network
(`docs/patterns/story-supervisor.md`). Do not claim a README scenario number
until the premise is settled and the README list is deliberately updated. See
wiki `[[Scenario Registry]]`.

---

## Y-statement

**We are building** a multi-node example where identical self-healing workers
sit on an explicit dependency graph and the joint health story is collapsed to
coarse pool labels, **for** readers who already understand one-node
`self_healing_nodes` and want the next real-ish step (wiring + partial
observation), **so that** we can measure what graph edges do to recovery and
what a dashboard-sized story loses — **without** folding in a shared evolver
or human operator yet.

**Why this route (not the alternatives):**

| Alternative | Why not (for v1) |
|-------------|------------------|
| Keep shared evolver from `self_healing_nodes` | Mixes **control plane** (mutate/promote) with **data-plane graph** (load edges); re-opens throttle-vs-evolver on a larger state space before graph + collapse are clear |
| “Topology” as ±spill floats only | Dresses a parameter as wiring; HO vs OH asymmetry under a real \(1\to2\) edge is the signature we want |
| Operator/on-call as a third agent in v1 | Valuable sequel; first we need a clean observation map (collapse) those agents would read |
| Heterogeneous per-node kernels (different weights) | Parked — see below; identical nodes isolate *graph* effects from *node identity* effects |

**Parked for implementation sequencing (hypotheses already filed):** whether each
self-healing node should carry different weights / Variant A vs B kernels /
asymmetric heal strength. That question is about **node identity**, not wiring.
Run after identical-node graph + collapse are working:

| Experiment | Claim | Issue |
|------------|-------|-------|
| `experiments/heal-on-critical-path` | Strong heal belongs on the downstream sink more than on the upstream feeder | #33 |
| `experiments/collapse-masks-heterogeneity` | Strong-sink / weak-feeder makes `pool_ok` look healthier than joint upstream-overload warrants | #34 |

---

## Real-ish scenario (no shared evolver)

A small worker pool / pipeline: each machine is the same self-healing node
(EMA error band → push / throttle / idle). Work has directed edges:

- **Upstream → downstream (load):** when an upstream node pushes or is
  overloaded, it dumps work onto its out-neighbors and raises their chance of
  degrading / overloading.
- **Downstream → upstream (backpressure, optional on some graphs):** when a
  downstream node is overloaded, upstream push becomes less effective / more
  risky (queues fill).

Each node still **perceives only its own EMA** and **decides independently**
(Kronecker \(D\)). Coupling lives in how neighbors change the next world —
i.e. in \(A_\text{joint}\), conditioned on the current joint world (who is
actually overloaded) and the joint action (who is pushing).

There is **no** pool-wide evolver agent in v1. The "outside view" is not an
agent in the product space; it is the **observation map** used for Trace
(collapsing the story). See "Operator / evolver (deferred discussion)" below
for how a real operator or config loop would re-enter later.

---

## State space

### Per-node model (reused from `self_healing_nodes`)

Each node is the same 3-state agent:

| Layer      | States                              |
|------------|-------------------------------------|
| World      | `healthy (H)`, `degraded (D)`, `overloaded (O)` |
| Experience | `ema_low`, `ema_mid`, `ema_high`    |
| Action     | `push`, `throttle`, `idle`          |

Default kernels: Variant B `nodeP`, `nodeD`, `nodeA` from
`examples/self_healing_nodes/main.go` (no `mutationBoost` — no evolver).

### Joint world space (2 nodes, no evolver)

Joint world = `(node1_health · node2_health)`, row-major, node1 outermost.

| Index | State | Node 1     | Node 2     |
|-------|-------|------------|------------|
| 0     | HH    | healthy    | healthy    |
| 1     | HD    | healthy    | degraded   |
| 2     | HO    | healthy    | overloaded |
| 3     | DH    | degraded   | healthy    |
| 4     | DD    | degraded   | degraded   |
| 5     | DO    | degraded   | overloaded |
| 6     | OH    | overloaded | healthy    |
| 7     | OD    | overloaded | degraded   |
| 8     | OO    | overloaded | overloaded |

9 joint world states. Joint experience = 3×3 = 9. Joint action = 3×3 = 9.
All of \(P_\text{joint}\), \(D_\text{joint}\), \(A_\text{joint}\), \(J\) are 9×9.

### Extending to 3 nodes

\(3^3 = 27\) world / experience / action states; \(J\) is 27×27. Tractable.
Use 3 nodes when a graph distinction needs it (e.g. chain vs fork vs ring).
**Implement 2-node first**; add a third node once collapse + one graph contrast
are legible.

---

## Graphs (variants)

A graph is a set of directed influence rules applied when building
\(A_\text{joint}\). Recommended v1 set:

| Variant | Name | Edges | Meaning |
|---------|------|-------|---------|
| 1 | Independent | none | Control: two self-healers, no load coupling |
| 2 | Linear chain | \(1 \to 2\) load | Pipeline: node1 feeds node2 |
| 3 | Bidirectional | \(1 \leftrightarrow 2\) load | Shared pool / mutual spill |

Optional later (3-node): chain \(1\to2\to3\), fan-in \(1,2\to3\), ring.
Do **not** invent a "protective pair" as a fourth fake topology unless it
corresponds to a real edge semantics (e.g. explicit load-shed edge).

### Coupling rule (world-conditioned load)

Prefer realism over the old action-only spill shortcut:

```
When building the next-world row for node v, start from nodeA[v's action].
For each in-neighbor u:
  if u's current world is overloaded (or degraded, with smaller weight)
     AND u's action is push:
       boost v's degraded/overloaded mass by spillRate, then renormalize
```

So spill requires **both** "neighbor is actually sick" and "neighbor is still
pushing" — closer to shedding load under stress. Document the rates in code
constants (`spillOverloaded`, `spillDegraded`).

Backpressure (optional, bidirectional variant only): if out-neighbor is
overloaded, damp the local node's `healthy` mass after `push` (queues won't
drain). Keep off for linear-chain v1 to isolate one mechanism.

---

## Kernel construction

### \(P_\text{joint}\) (9×9) — local perception

Each node perceives from **its own** world only (Kronecker of `nodeP`):

```
P_joint[(w1,w2), (x1,x2)] = nodeP[w1,x1] * nodeP[w2,x2]
```

No cross-perception in v1 (nodes do not see neighbors' EMAs). That keeps
"collapsing the story" as an **external** observation problem, not something
the nodes themselves solve.

### \(D_\text{joint}\) (9×9) — independent decisions

```
D_joint[(x1,x2), (g1,g2)] = nodeD[x1,g1] * nodeD[x2,g2]
```

### \(A_\text{joint}\) (9×9) — graph coupling

PDA composition needs \(A: G \to W\), but neighbor influence depends on the
**current** world. Two implementation options (pick one; document in code):

**Preferred for correctness:** build \(J\) by enumerating
`(w, x, g) → w'` micro-paths (or build a world-dependent family of A slices
and form \(W = \sum\) appropriately). Practically: nest loops over
`w1,w2, g1,g2, w1',w2'` with

```
row1 = adjust(nodeA[g1], influences from neighbors given (w1,w2,g1,g2), graph)
row2 = adjust(nodeA[g2], ...)
A_eff[(g1,g2), (w1',w2') | (w1,w2)] = row1[w1'] * row2[w2']
```

then fold with P and D into \(J[(w1,w2),(w1',w2')]\).

**Simpler approximation (acceptable for first cut):** ignore world conditioning
in A; spill only from `gi == push`. Note explicitly as a simplification in the
walkthrough (same caveat as the old plan).

Compose \(J = P_\text{joint}\, D_\text{joint}\, A_\text{joint}\) only when A is
not world-conditioned; otherwise assemble \(J\) from the nested loops above
(still a valid row-stochastic square kernel on joint W).

---

## Collapsing the story (primary analysis)

### Operator-visible coarse states

Define an observation partition of the 9 joint worlds — the "dashboard":

| Coarse label | Joint worlds | Ops meaning |
|--------------|--------------|-------------|
| `pool_ok` | `{HH}` | All green |
| `pool_stressed` | `{HD, DH, HO, OH, DD}` | Partial / uneven sickness |
| `pool_down` | `{DO, OD, OO}` | At least one overloaded + partner sick, or both overloaded |

(Exact membership can be tuned; freeze before reading metrics.)

**Important:** Trace in catrace is defined on a **subset** of states (induced
chain), not on an arbitrary partition into superstates. Two approaches:

1. **Subset Trace (exact API):** pick representative extremes
   `Trace([]int{HH, OO})` — peak vs worst — same pattern as
   `self_healing_nodes`. Show `IsTraceOf`, π on the subset.
2. **Partition collapse (story):** for the 3-bucket dashboard, either
   - lump by analyzing π mass and MFPT **on the full J** with sets as
     targets (MFPT into the set `pool_ok`), and report set-occupancy from π; or
   - build an explicit lumped chain only if we add a helper / do it by hand
     for the walkthrough.

v1 recommendation: use (1) for the Trace API demo, and report (2) as
**set statistics on J** (π(set), mean MFPT from each `pool_down` state into
`HH`) so the "collapsed story" is narrative + numbers without faking a
partition-Trace the library does not provide.

Optional second act: sample a long trajectory on J, map each step to
`{pool_ok, pool_stressed, pool_down}`, estimate a 3×3 kernel, compare
qualitatively to the set-to-set flow implied by J.

### Metrics per graph

```
π on all 9 states; π(pool_ok), π(pool_down)
MFPT(OO → HH), MFPT(HO → HH), MFPT(OH → HH)   // asymmetry under chain
H(J)
Trace onto {HH, OO}: IsTraceOf, π_trace
```

Comparative table:

```
Graph          | π(HH) | π(OO) | MFPT OO→HH | MFPT HO→HH | MFPT OH→HH | H
Independent    |  ...  |  ...  |    ...     |    ...     |    ...     | ...
Linear chain   |  ...  |  ...  |    ...     |    ...     |    ...     | ...
Bidirectional  |  ...  |  ...  |    ...     |    ...     |    ...     | ...
```

Under linear chain, expect **HO→HH vs OH→HH asymmetry** (downstream sick vs
upstream sick) — that is the signature that the graph is doing real work.
Independent should be symmetric.

---

## Code structure

```go
type graph struct {
    name string
    // loadEdge[u][v] = spill rate from u onto v when u overloaded & pushing
    loadEdge [2][2]float64
}

func buildJoint(g graph) *catrace.Kernel  // returns J on 9 joint worlds
func analyze(name string, J *catrace.Kernel)
```

Reuse `nodeP` / `nodeD` / `nodeA` constants from Variant B; no evolver blocks.

HTML: `ToHTML` for J and for the {HH, OO} trace per graph (optional).

---

## WALKTHROUGH.md structure (per conventions)

**§1 Opening:** After `self_healing_nodes` (one node + evolver). Here: many
nodes, graph-shaped load, and collapsing the joint story to what ops would see.

**§2 Story:** Two (or three) identical self-healing workers on a dependency
graph. Each only sees its own EMA. Load still crosses edges. The dashboard
collapses nine joint worlds into a few pool labels — what is lost?

**§3 State spaces:** One node table (shared by all); joint 9-state table.
No evolver table.

**§4 Coupling:** \(P\) local Kronecker; \(D\) Kronecker; \(A\) / \(J\) carry
graph load edges (world-conditioned spill).

**§5 Math by hand:** One linear-chain entry, e.g. contribution to
`HH → HD` when node1 is healthy-pushing vs when node1 is overloaded-pushing
(show the spill term turning on).

**§6 Output:** Per-graph π, MFPT table (including HO vs OH), entropy, Trace
{HH, OO}. Interpret asymmetry under chain.

**§7 What you can change:**
1. Spill rates — strengthen cascade
2. Add backpressure edge
3. Switch Variant A nodeA (strong local heal) — does graph matter less?
4. Three-node chain
5. Empirical 3-bucket estimated kernel vs set occupancy on J

---

## README entry (draft — assign number only when adding to README)

**Story:** Several identical self-healing workers share a dependency graph.
Each watches only its own error EMA and throttles or pushes on its own. Load
still crosses edges: an overloaded upstream node that keeps pushing stresses
its downstream neighbor. The full joint health story is large; the story we
usually tell is smaller — whether the pool looks OK — and that collapse hides
which side of the graph is burning.

**State meanings:** Per-node world / experience / action (same as
self_healing_nodes). Joint worlds HH…OO. Coarse pool labels for the collapsed
story.

**Interpretation:**
- Nodes are independent decision-makers (\(D\) Kronecker)
- Graph edges enter as load (and optional backpressure) in joint world dynamics
- Trace / set statistics collapse the joint story to operator-visible buckets
- Graph contrast (esp. HO vs OH MFPT) shows wiring is not cosmetic

**Code:** `examples/network_of_healers/main.go`

---

## Layers: dashboard vs operator vs evolver

Three real-world roles often get smashed together. v1 only implements the
first as a **collapse map**; the others are sequels.

| Layer | Real world | In the model | v1? |
|-------|------------|--------------|-----|
| **Data plane** | Workers + dependency edges | Nodes + load / backpressure in \(A\) / \(J\) | Yes |
| **Observation** | Grafana / SLO tile (green·yellow·red) | Collapse map: Trace subset `{HH, OO}` + set stats on `{pool_ok, pool_stressed, pool_down}` — **not** an agent; no actions | Yes |
| **Control plane** | Autoscaler / workerpool evolver | Shared config agent (`self_healing_nodes`): \(P\) from aggregate severity; `mutate` boosts recovery in \(A\) | Deferred |
| **Exception path** | On-call human | Slow agent: experience = coarse pool label; actions = restart / shed traffic / page (hits many node worlds) | Deferred |
| **Mesh / LB** | How work is wired | Graph edges and spill rates | Yes (as topology) |

**Dashboard** only *sees* a partition of the true joint state. **On-call**
acts on that coarse view, rarely and heavily. **Evolver / autoscaler** is a
continuous control plane over shared config — real, but it products the state
space by ×2 and reintroduces throttle-vs-mutate on top of the graph.

v1 separates **graph physics** (nodes + edges) from **control plane**
(operator / evolver) so the collapse story stays clean.

---

## Open questions / risks

| Risk | Mitigation |
|------|-----------|
| World-conditioned A does not factor as a single \(G\to W\) matrix | Assemble \(J\) from nested loops; still expose a square Kernel |
| Partition "Trace" is not what `Trace([]int)` does | Use subset Trace for API; set π / set MFPT for dashboard narrative |
| Independent vs chain metrics may be too close under Variant B | Raise spillRate until HO/OH asymmetry is obvious |
| 3-node stretch grows walkthrough surface area | Gate behind working 2-node collapse story |
| Naming 9 states still dense | Abbreviated table + code comments |
| Premise drift vs README scenario 5 | Keep unnumbered until deliberate README edit; Scenario Registry |
| Heterogeneous node weights / per-node Variant A vs B | **Parked for code** — hypotheses filed: `heal-on-critical-path`, `collapse-masks-heterogeneity` |
