# Hypothesis: Heal capacity belongs on the downstream critical-path node

> Fill in Results and Verdict after `examples/network_of_healers` supports heterogeneous node kernels on a fixed chain graph.

## Claim

On a fixed linear load graph \(1 \to 2\), putting **strong local self-heal (Variant A `nodeA`)** on the **downstream** node reduces cascade damage more than putting the same strong heal on the **upstream** feeder (with the partner weak — Variant B). Heal capacity is a **vertex** resource; for pipeline load, it belongs on the sink that absorbs spill, not only on the node that sheds it.

This is falsifiable: if strong-upstream / weak-downstream wins on π(HH) and MFPT(OO→HH), the claim is wrong — feeders matter more than sinks.

## Context

- **Pattern:** Self-Healing / Adaptive (14), multi-node data-plane graph
- **Related example / plan:** `plans/network-of-healers.md`, depends on `examples/network_of_healers` (identical-node v1 first)
- **Issue:** [#33](https://github.com/stephen-mcelhose/catrace/issues/33)
- **Prerequisite:** Identical-node chain must already show clear HO vs OH MFPT asymmetry (graph is real)
- **Motivation:** Engineering effort is finite. Strengthening every node's throttle is expensive. This experiment asks whether heal strength should be concentrated on the critical-path **vertex** that receives load — the design question identical nodes cannot ask.

## Variant definitions

Graph fixed: linear chain \(1 \to 2\) with the same spill rates in all variants.
No shared evolver. Per-node P and D identical; only each node's `nodeA` (Variant A vs B rows from `self_healing_nodes`) differs.

| Variant | Label | Node 1 (upstream) | Node 2 (downstream) |
|---------|-------|-------------------|---------------------|
| A | Strong sink | Variant B `nodeA` (weak heal) | Variant A `nodeA` (strong heal) |
| B | Strong feeder | Variant A `nodeA` (strong heal) | Variant B `nodeA` (weak heal) |
| C | Both weak (control) | Variant B | Variant B |
| D | Both strong (control) | Variant A | Variant A |

*Controls C/D bound the effect: if A does not beat B, and does not beat C by a clear margin on sink-stress metrics, the claim fails.*

## Variable kernel entries

Encode Variant A / B `nodeA` exactly as in `examples/self_healing_nodes/main.go` (no mutation boost):

| Kernel | Row (action → world) | Strong (Var A) | Weak (Var B) | Why |
|--------|----------------------|----------------|--------------|-----|
| nodeA | push → (H, D, O) | (0.70, 0.25, 0.05) | (0.55, 0.35, 0.10) | Push risk under load |
| nodeA | throttle → (H, D, O) | (0.65, 0.30, 0.05) | (0.45, 0.45, 0.10) | **Primary heal strength** |
| nodeA | idle → (H, D, O) | (0.85, 0.13, 0.02) | (0.60, 0.35, 0.05) | Idle recovery floor |

Spill / graph edges are **identical** across variants. Only which vertex gets the strong `nodeA` changes.

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| Stationary mass | π(HH) | A > B | Strong sink absorbs spill; feeder strength alone does not protect the pool |
| Stationary mass | π(HO) + π(OO) (downstream sick) | A < B | Downstream is the stressed vertex under \(1\to2\) load |
| Recovery speed | MFPT(OO → HH) | A < B | Strong sink climbs out of overload faster once spill stops |
| Asymmetry | MFPT(HO → HH) / MFPT(OH → HH) | A < B (ratio closer to 1 or HO easier in A) | Strong sink specifically shortens recovery from “downstream sick” |
| Predictability | H(J) bits/step | A ≤ B | Less cascading randomness when the sink can self-stabilize |
| Control check | π(HH): D ≥ A ≥ C | ordering holds | Both-strong best, both-weak worst; A between them and above B |

## Verdict rule

Claim supported if **A beats B on at least 4 of 5** primary metrics (π(HH), downstream-sick mass, MFPT(OO→HH), HO/OH asymmetry, H) **and** the control ordering π(HH): D ≥ A ≥ C holds. A split where B wins π(HH) but A wins MFPT is a **trade-off** (long-run vs recovery), not support.

---

## Results

> Fill in after running (once heterogeneous kernels exist in `network_of_healers`).

| Metric | A (strong sink) | B (strong feeder) | C (both weak) | D (both strong) | Predicted | Correct? |
|--------|-----------------|-------------------|---------------|-----------------|-----------|----------|
| π(HH) | | | | | A > B; D ≥ A ≥ C | |
| π(HO)+π(OO) | | | | | A < B | |
| MFPT(OO → HH) | | | | | A < B | |
| MFPT(HO → HH) | | | | | A < B | |
| MFPT(OH → HH) | | | | | (context) | |
| H(J) | | | | | A ≤ B | |

## Verdict

**Claim:** [pending]
**Metrics in agreement:** [—]

## Interpretation

[After results: did sink-local heal outperform feeder-local heal? Implication for where to spend throttle engineering on a pipeline.]
