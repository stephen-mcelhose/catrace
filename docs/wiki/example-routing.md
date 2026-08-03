---
title: "Example: Routing Agent"
tags: [example, routing, classification, mfpt, trace, single-agent]
sources: [examples/routing/WALKTHROUGH.md, docs/patterns/story-routing.md]
updated: 2026-08-05
---

# Example: Routing Agent

**Pattern:** Routing (3) — see [[Agentic Patterns Catalogue]]
**Code:** `examples/routing/main.go`
**Run:** `go run examples/routing/main.go`

This example demonstrates the routing topology: a single agent classifies
incoming work and dispatches it. Misclassification creates feedback loops where
tickets cycle through wrong specialists before reaching the right one. The
example is the first to use `MeanFirstPassage` as a primary analysis tool,
measuring misrouting loop latency directly.

## The scenario

A customer support router classifies tickets (billing, technical, general) and
dispatches them to specialists. Wrong routing sends a ticket to a specialist
who cannot resolve it; the ticket re-enters the queue, often as an apparently
different type. The world kernel W captures the full ticket-flow dynamics
including these misrouting loops.

## State spaces

| Layer      | States                                                                           | Meaning                           |
|------------|----------------------------------------------------------------------------------|-----------------------------------|
| World      | `billing_ticket`, `technical_ticket`, `general_ticket`                           | True ticket type                  |
| Experience | `reads_billing`, `reads_technical`, `reads_general`                              | Router's perceived classification |
| Action     | `route_billing`, `route_technical`, `route_general`, `escalate_human`            | Routing decision                  |

## Key results

**Stationary distribution (W):**

```
billing_ticket   0.304149
technical_ticket 0.371701
general_ticket   0.324150
```

Technical tickets dominate the queue (37.2%) because the router classifies them
most accurately (P diagonal = 0.80 vs 0.75 for billing and general), retaining
them in the correct type longer.

**Entropy rate:** 1.472 bits/step — lower than the 5-state pipeline (2.22 bits),
reflecting a more predictable 3-state symmetric system.

**MFPT:**
- billing→technical: 3.55 steps — how long before a billing ticket reaches the technical queue through misclassification
- technical→general: 4.37 steps — slightly longer because technical tickets resist crossing to general

These numbers are direct measurements of misrouting loop latency. Improving
classification accuracy (higher P diagonal) increases both MFPTs — confirming
that better routing means tickets take longer to cross type boundaries.

## Trace onto {billing_ticket, technical_ticket}

Treating general enquiries as a hidden background process:

```
IsTraceOf = true
stationary(trace) = 0.450024  0.549976
                    billing    technical
```

With general collapsed, technical holds 55% of the effective queue vs 45% for
billing. The 2-state kernel captures the billing/technical dynamics exactly
without approximation. See [[Trace Chain]].

## Connection to math and API

- `Agent.WorldKernel()` → [[catrace API]], [[PDA Triplet Model]]
- `Kernel.MeanFirstPassage()` → [[catrace API]], [[Markov Chain Foundations]]
- `Kernel.Trace()` and `IsTraceOf()` → [[Trace Chain]]
- `Kernel.Stationary()`, `EntropyRate()`, `Classes()` → [[Markov Chain Foundations]]

## Sources

- `examples/routing/WALKTHROUGH.md`
- `docs/patterns/story-routing.md`
