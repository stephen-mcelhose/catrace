# ADR-0001: Accept custom starting distributions — `StationaryFrom` and `PersonalizedPageRank`

**Status:** Accepted  
**Date:** 2025-07-17  
**Issue:** [#37 — feat: StationaryFrom — power iteration from a custom initial distribution](https://github.com/stephen-mcelhose/catrace/issues/37)

---

## Context

`Kernel.Stationary` runs power iteration seeded from the uniform distribution `1/n`. This is the correct default for ergodic chains — the uniform start is distribution-agnostic and converges to the unique stationary distribution whenever the chain is irreducible and aperiodic.

However, several important use cases require starting from a *caller-supplied* distribution or a biased iteration rule:

1. **Distribution-independence verification.** Starting from multiple different distributions and confirming convergence to the same fixed point is a practical ergodicity check. `Stationary` today makes this impossible without re-implementing the loop.

2. **Personalized PageRank (PPR).** PPR replaces global ranking with topic- or node-scoped ranking. Its iteration is:

   ```
   x_{t+1} = α·v + (1−α)·x_t·P   until ‖x_{t+1} − x_t‖_∞ < tol
   ```

   where `v` is a restart distribution concentrated on a "seed" set (e.g., goal nodes in a wiki graph). The PPR fixed point is the stationary distribution of the *teleporting chain* `α·v·𝟙ᵀ + (1−α)·P` — it is **not** a stationary distribution of `P` itself and cannot be computed by passing `v` as the start vector to plain power iteration.

3. **Consumer requirement.** [wikigraph issue #22](https://github.com/stephen-mcelhose/wikigraph/issues/22) requires PPR seeded at goal nodes as a structural ranking alternative to MFPT. Implementing PPR cleanly inside wikigraph means duplicating the convergence loop — a maintenance liability.

The core primitive `LeftAction` already exists. The gap is: (a) flexible initialisation and (b) the teleportation blend in the PPR loop body.

---

## Decision

We will add two methods to `stationary.go`.

### `StationaryFrom` — plain power iteration from a custom start

```go
// StationaryFrom runs power iteration initialised from start rather than the
// uniform distribution. start must be a valid probability distribution over
// k.NumStates() states (non-negative, sums to 1 within tol).
//
// Use Stationary for the standard uniform-start case. Use StationaryFrom when
// the initial distribution carries semantic meaning or when verifying that
// convergence is distribution-independent.
func (k *Kernel) StationaryFrom(start []float64, tol float64, maxIter int) ([]float64, error)
```

`Stationary` is refactored to delegate to `StationaryFrom` with a freshly-allocated uniform vector, preserving its existing contract exactly.

### `PersonalizedPageRank` — biased power iteration with teleportation

```go
// PersonalizedPageRank computes the PPR vector for the given restart
// distribution and teleportation weight alpha ∈ [0, 1].
//
// Each iteration step is:
//
//	x ← α·restart + (1−α)·(x·P)
//
// For alpha ∈ (0, 1] convergence is guaranteed because the teleporting chain
// is strongly connected. For alpha = 0 this reduces to plain power iteration
// (equivalent to StationaryFrom(restart, ...)).
func (k *Kernel) PersonalizedPageRank(restart []float64, alpha, tol float64, maxIter int) ([]float64, error)
```

PPR is a **sibling** of `StationaryFrom`, not a wrapper. Its loop body differs: at every step the propagated distribution is blended back toward `restart` by factor `alpha`, so the fixed point satisfies `x* = α·v + (1−α)·x*·P`. Plain power iteration does not produce this fixed point when started from `v`.

### Validation contract (both methods)

**`restart` / `start`:**
- Length must equal `k.NumStates()`.
- All entries must be ≥ 0 (entries in `(-tol, 0)` are clamped to 0 before validation).
- Sum must be within `tol` of 1.0.

**`alpha` (PPR only):**
- Must satisfy `0 ≤ alpha ≤ 1`. Values outside this range are rejected with an error.

Both methods default `maxIter` to 1000 when ≤ 0, consistent with `Stationary`.

---

## Consequences

### Positive

- Removes the need for consumers to duplicate any iteration loop.
- Enables PPR inside wikigraph directly via `k.PersonalizedPageRank(goals, alpha, tol, maxIter)`.
- Enables distribution-independence checks via `StationaryFrom`.
- Zero breaking change: `Stationary` signature and behaviour are unchanged.
- PPR convergence for `alpha ∈ (0, 1]` is unconditional — the teleporting chain is always ergodic, eliminating the non-ergodic-chain caveat that affects `StationaryFrom`.

### Negative / trade-offs

- `StationaryFrom` does **not** guarantee convergence for non-ergodic chains. The doc comment must say so explicitly.
- `tol` serves double duty (validation and convergence criterion) in both methods, matching the existing discipline in `Validate` and `normalizeVector`.
- Adding `PersonalizedPageRank` to `Kernel` couples the core Markov library to the `alpha` parameter and the teleportation idiom. This is accepted: PPR is a well-established Markov operation, not an application-level concern, and the method is clearly scoped by its name.

---

## Alternatives considered

### A. Add a `WithStart` functional option to `Stationary`

Rejected. Options pattern is heavier than warranted for a single variant of one function. A second method is cleaner and keeps the API surface self-documenting.

### B. Export `normalizeVector` and let callers manage iteration themselves

Rejected. Forces every consumer to re-implement the convergence loop — exactly the duplication we want to eliminate.

### C. Keep `PersonalizedPageRank` as a caller-side wrapper

Superseded. PPR requires a distinct loop body (teleportation blend per step); it is not constructible by simply passing `restart` as the `start` argument to `StationaryFrom`. Exposing it on `Kernel` removes boilerplate for every consumer and keeps the blend logic in one tested place.
