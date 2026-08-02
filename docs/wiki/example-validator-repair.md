---
title: "Example: Validator Repair"
tags: [example, multi-agent, joint-kernel, coupling, product-space, trace, stationary]
sources: [examples/validator_repair/WALKTHROUGH.md, docs/patterns/story-validator-repair.md]
updated: 2026-08-02
---

# Example: Validator Repair

**Pattern:** Evaluator-Optimizer (6), Self-Healing / Adaptive (14) — see [[Agentic Patterns Catalogue]]
**Code:** `examples/validator_repair/main.go`
**Run:** `go run examples/validator_repair/main.go`

This example bridges the first two examples: it demonstrates how to **build** a joint kernel from two independently-modeled agents — and why the construction method matters. It is the first catrace example to use [[Joint Kernels and Coupling]].

## The two agents

**Worker:** performs tasks; may be valid or invalid.
**Validator:** monitors the Worker and attempts repair when problems are detected.

Each agent is modeled with its own P, D, A triplet. Joint kernels are defined over the product of their state spaces and composed exactly as for a single agent: J = P_joint · D_joint · A_joint.

## Joint state spaces

| Space | States | Size |
|-------|--------|------|
| W₁×W₂ | VV, VI, IV, II | 4 |
| X₁×X₂ | ok·good, ok·bad, prob·good, prob·bad | 4 |
| G₁×G₂ | 9 worker×validator action pairs | 9 |

Where V=valid, I=invalid: VV=both valid, VI=worker valid+validator invalid, IV=worker invalid+validator valid, II=both invalid.

## Coupling points

### Coupled perception (P_joint)

The Validator's perception is coupled to the Worker's world state:

```
P_joint[(w₁,w₂), (x₁,x₂)] = P₁[w₁, x₁] × P₂_coupled[(w₁,w₂), x₂]
```

When the Worker is degraded (w₁ = invalid), the Validator's probability of perceiving `looks_bad` rises from 0.15 (baseline) to 0.40 — nearly 3× more likely to trigger repair. Without this coupling, the Validator cannot reliably detect Worker degradation.

### Independent decisions (D_joint)

```
D_joint = D_worker ⊗ D_validator
```

Each agent decides based on its own experience without communicating. The Kronecker product factorization encodes this.

### Coupled action effect (A_joint)

When the Validator repairs (g₂ = repair), it raises the Worker's probability of returning to valid by ~0.20:

| Worker action | Worker→valid (independent) | Worker→valid (with Validator repair) |
|---------------|---------------------------|--------------------------------------|
| produce | 0.70 | 0.90 |
| self_check | 0.50 | 0.70 |
| idle | 0.40 | 0.60 |

For non-repair actions, A_joint = A_worker ⊗ A_validator (no coupling).

## Joint kernel J output

```
       VV      VI      IV      II
VV [ 0.475   0.202   0.233   0.090 ]
VI [ 0.484   0.246   0.188   0.082 ]
IV [ 0.412   0.196   0.276   0.116 ]
II [ 0.422   0.227   0.238   0.113 ]
```

Key readings:
- **IV → VV (0.412):** Validator detects degraded Worker and repairs it; 41% recovery per step
- **II → VV (0.422):** Even from full degradation, coupled repair achieves reasonable recovery (repair boost in A_joint)
- **VV → IV (0.233):** Worker degrades even when Validator is fine — one cycle's lag before perception triggers repair

**Stationary distribution:**

```
VV ≈ 0.457,  VI ≈ 0.212,  IV ≈ 0.234,  II ≈ 0.097
```

~46% of time both valid; only ~10% fully degraded. The single recurrent class confirms ergodicity.

## Trace onto {VV, II}

The coarsest operational picture: healthy or failed, with mixed states collapsed.

```
IsTraceOf = true
stationary(J)|{VV,II} normalized = 0.825220  0.174780
stationary(trace)                = 0.825220  0.174780
```

The coarse picture says: when you can only see fully-healthy vs. fully-failed, the system is in the healthy regime 82.5% of the time. `IsTraceOf = true` confirms the mathematical projection is exact. See [[Trace Chain]].

## Why the coupling matters

Without coupled perception: Validator in state IV sees `looks_bad` with 0.15 probability instead of 0.40 — less than 1/3 as likely to trigger repair. Without coupled action: Validator's repair has no effect on the Worker. Together, the two coupling points create a recovery pathway that neither agent has alone.

## Connection to math and API

- Joint kernel construction via `Agent{P: pJoint, D: dJoint, A: aJoint}.WorldKernel()` → [[catrace API]]
- Kronecker product for D_joint → [[Joint Kernels and Coupling]]
- Trace and verification → [[Trace Chain]]
- Stationary, Classes (ergodicity) → [[Markov Chain Foundations]]

## Sources

- `examples/validator_repair/WALKTHROUGH.md`
- `docs/patterns/story-validator-repair.md`
