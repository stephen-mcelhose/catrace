# Hypothesis: Pool collapse masks upstream failure under a strong downstream healer

> Fill in Results and Verdict after `examples/network_of_healers` supports heterogeneous node kernels and coarse pool labels.

## Claim

On a linear load graph \(1 \to 2\), a **strong-heal downstream + weak-heal upstream** configuration makes the **collapsed** (dashboard) story look healthier than the joint story warrants: coarse mass on `pool_ok` / Trace π(HH) stays high while joint mass on **upstream-overloaded** states (OH, OD, OO contribution from node1) remains material. The same collapse under **identical weak** nodes does not hide upstream failure as severely.

In short: heterogeneity lets the dashboard lie in a specific way — “pool looks fine / mostly fine” while the feeder burns — which identical-node collapse does not demonstrate as sharply.

## Context

- **Pattern:** Self-Healing / Adaptive (14); observation / Trace collapse
- **Related example / plan:** `plans/network-of-healers.md` (“collapsing the story”); pairs with `experiments/heal-on-critical-path`
- **Issue:** [#34](https://github.com/stephen-mcelhose/catrace/issues/34)
- **Prerequisite:** Coarse labels `{pool_ok, pool_stressed, pool_down}` (or Trace `{HH, OO}`) defined and working on identical-node J
- **Motivation:** Ops acts on collapsed labels. If heal is uneven, yellow/green tiles can mean “sink is coping” rather than “system is healthy.” This experiment falsifies whether that masking is large enough to matter for architectural trust in the dashboard.

## Variant definitions

Graph fixed: linear chain \(1 \to 2\). No evolver.

| Variant | Label | Node 1 (upstream) | Node 2 (downstream) |
|---------|-------|-------------------|---------------------|
| A | Masking config | Variant B (weak) | Variant A (strong) |
| B | Identical weak | Variant B | Variant B |

*Optional third run (report-only, not required for verdict): identical strong (A,A) — expect little to hide because little fails.*

## Variable kernel entries

Same `nodeA` strong/weak tables as `experiments/heal-on-critical-path/hypothesis.md`. Graph spill identical across variants.

## Coarse observation map (freeze before run)

| Coarse label | Joint worlds |
|--------------|--------------|
| `pool_ok` | `{HH}` |
| `pool_stressed` | `{HD, DH, HO, OH, DD}` |
| `pool_down` | `{DO, OD, OO}` |

Also report Trace onto `{HH, OO}` for API consistency with other examples.

**Hidden failure mass** (joint, not coarse):

\[
\pi_{\text{upstream\_O}} := \pi(\mathrm{OH}) + \pi(\mathrm{OD}) + \pi(\mathrm{OO})
\]

(node1 overloaded, any node2).

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| Dashboard health | π(`pool_ok`) | A > B | Strong sink keeps HH more often |
| Masking gap | π(`pool_ok`) / (1 − π_upstream_O) or simpler: π(`pool_ok`) − (1 − π_upstream_O) | gap larger in A (more optimistic) | Coarse “OK” outruns true feeder health under A |
| Hidden feeder stress | π_upstream_O | A has π_upstream_O **not much smaller** than B, while π(pool_ok) is **larger** | Sink heal improves tile without fixing feeder |
| Trace optimism | π_trace(HH) on {HH, OO} | A > B | Extreme Trace also looks better under A |
| Stress composition | Among `pool_stressed`, share of OH vs HO | OH share higher in A than in B | When stressed, A’s yellow tile is more often “feeder bad, sink OK” |

Primary falsifier: if under A, π_upstream_O drops in lockstep with the rise in π(pool_ok) (no gap), masking is not happening — sink heal is truly system heal.

## Verdict rule

Claim supported if **all three** hold:

1. π(`pool_ok`): A > B  
2. π_upstream_O does **not** fall as much as π(pool_ok) rises (define:  
   \(\Delta\pi_{\mathrm{ok}} - |\Delta\pi_{\mathrm{upstream\_O}}| > \delta\) with \(\delta = 0.05\) as the minimum interesting gap — tune \(\delta\) only before first run)  
3. Among `pool_stressed` mass, OH/(OH+HO) is higher in A than in B (when both denominators > 0)

If (1) holds but (2)–(3) fail, revise claim to “strong sink heals the whole pool” (no special masking).

---

## Results

> Fill in after running.

| Metric | Variant A (masking) | Variant B (identical weak) | Predicted | Correct? |
|--------|---------------------|----------------------------|-----------|----------|
| π(pool_ok) | | | A > B | |
| π_upstream_O | | | A not ≪ B relative to Δπ(ok) | |
| Masking gap (δ rule) | | | > 0.05 | |
| OH/(OH+HO) in stressed | | | A > B | |
| π_trace(HH) | | | A > B | |

## Verdict

**Claim:** [pending]
**Criteria met:** [—/3]

## Interpretation

[After results: should ops trust pool_ok when heal strength is known to be uneven? Implication for per-node dashboards vs pool tiles.]
