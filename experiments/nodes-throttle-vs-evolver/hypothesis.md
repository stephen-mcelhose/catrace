# Hypothesis: Node throttle is the primary recovery mechanism

## Claim

The node's own throttle action is the primary driver of recovery from degraded states. The evolver's mutation search contributes to long-run health but is not the bottleneck: a system where throttle is strong recovers faster, achieves higher steady-state health, and is more predictable than one that depends on the evolver to carry the load.

## Context

- **Pattern:** Self-Healing / Adaptive (14)
- **Related example:** `examples/self_healing_nodes`
- **Motivation:** When designing a self-healing system with two feedback loops (local + outer), it matters which loop does the real work. If the outer loop is the primary recovery pathway, it adds latency (mutation search is slow and probabilistic). If the local loop is primary, the outer loop is a safety net rather than a critical path. This distinction changes how much engineering effort each loop deserves and what failure modes to design for.

## Variant definitions

| Variant | Label | Description |
|---------|-------|-------------|
| A | Strong throttle | Node's throttle action has high self-recovery probability (0.65). Evolver's mutation search provides only a marginal boost (+0.05) when it cooperates with throttle. |
| B | Weak throttle, strong evolver | Node's throttle action has lower self-recovery probability (0.40). Evolver's mutation search provides a substantial boost (+0.25) when it cooperates with throttle. |

*State spaces, agent structure, and all coupling points are identical in both variants. Only the A_joint entries for throttle-related transitions differ.*

## Variable kernel entries

| Kernel | Transition | Variant A | Variant B | What it encodes |
|--------|-----------|-----------|-----------|-----------------|
| A_joint | (throttle, promote) → healthy | 0.65 | 0.40 | Node's self-recovery strength when evolver is not mutating |
| A_joint | (throttle, mutate) → healthy | 0.70 | 0.65 | Node self-recovery + weak evolver boost (A) vs. strong evolver boost (B) |
| A_joint | (push, mutate) → healthy | 0.20 | 0.35 | Evolver's contribution when node is not throttling |
| A_joint | (idle, mutate) → healthy | 0.15 | 0.30 | Evolver's floor contribution |

*All other A_joint entries are identical across variants. A_joint rows are renormalized after applying boost amounts to maintain row-stochasticity.*

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|----|
| Stationary mass | π(H·G) | A > B | Strong local recovery keeps the system in the best joint state more often |
| Recovery speed | MFPT(O·B → H·G) | A < B | Local throttle acts every cycle; mutation search must first find a better config |
| Predictability | H(J) bits/step | A < B | A system that heals locally is more legible than one depending on random search |
| Observable health | π_trace(H·G) after tracing onto {H·G, O·B} | A > B | Coarse observer sees more "everything fine" states in Variant A |

## Verdict rule

All four metrics should agree. The claim is directional — throttle *primary*, not throttle *exclusive* — so even a small consistent advantage for A across all metrics is sufficient to support it. A split (A wins on π but B wins on MFPT) would indicate a genuine trade-off and require a revised claim.

---

## Results

| Metric | Variant A | Variant B | Predicted direction | Correct? |
|--------|-----------|-----------|--------------------|----|
| π(H·G) | 0.524 | 0.408 | A > B | ✅ |
| MFPT(O·B → H·G) | 2.03 steps | 2.56 steps | A < B | ✅ |
| H(J) bits/step | 1.87 | 2.12 | A < B | ✅ |
| π_trace(H·G) on {H·G, O·B} | 0.977 | 0.951 | A > B | ✅ |

Source: `examples/self_healing_nodes/main.go` — run `go run examples/self_healing_nodes/main.go`

## Verdict

**Claim: supported (4/4 metrics)**

## Interpretation

The node's local throttle is the primary recovery mechanism. In Variant B, where the evolver carries the load:

- The system spends ~10 percentage points more time in degraded joint states (D·G, D·B, O·G, O·B combined: 44% vs 29%)
- Recovery from the worst state (O·B) takes 26% longer (2.56 vs 2.03 steps)
- The system is 0.25 bits/step less predictable — the random mutation search adds irreducible entropy

The practical implication: in a self-healing node architecture, the inner throttle loop should be designed to handle the common case reliably. The evolver should be sized as a slow safety net for structural problems, not as a fast-path recovery mechanism.

**A plausible trade-off that did not materialize:** If the evolver's mutation search could discover genuinely better configurations (not just return to baseline), Variant B might win on long-run π even if it loses on MFPT. This experiment's A_joint models restoration to baseline only. An experiment with asymmetric state spaces (where mutation can reach a "better than original" healthy state) would test whether the outer loop can justify its latency cost through configuration improvement.
