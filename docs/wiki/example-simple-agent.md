---
title: "Example: Simple Agent"
tags: [example, simple-agent, pda, qualia-kernel, stationary, entropy-rate, ergodic]
sources: [examples/simple_agent/WALKTHROUGH.md, docs/patterns/story-single-llm-agent.md]
updated: 2026-08-05
---

# Example: Simple Agent

**Pattern:** Augmented LLM (1), Autonomous Agent Loop (7) — see [[Agentic Patterns Catalogue]]
**Code:** `examples/simple_agent/main.go`
**Run:** `go run examples/simple_agent/main.go`

This example is the minimal working demonstration of the [[PDA Triplet Model]]. A single LLM support agent handles tasks one at a time. The Q=DAP kernel captures how the agent's interpretation of tasks evolves step by step.

## State spaces

| Space | States | Meaning |
|-------|--------|---------|
| W (world) | `task_routine`, `task_complex` | Actual task difficulty |
| X (experience) | `looks_routine`, `looks_risky` | Agent's internal interpretation |
| G (action) | `answer`, `clarify`, `escalate` | What the agent does next |

## The three kernels

**P (2×2, W→X):** Imperfect perception. A routine task reads as routine 85% of the time; a complex task reads as risky 75% of the time. The 15% and 25% off-diagonal entries encode reading errors.

**D (2×3, X→G):** Decision policy. Given `looks_routine`: answer (0.8), clarify (0.2), escalate (0.0). Given `looks_risky`: answer (0.1), clarify (0.3), escalate (0.6).

**A (3×2, G→W):** Action effects. `answer` keeps the world routine 90% of the time but barely helps complex tasks (10%). `clarify` → routine 40% / complex 60%. `escalate` → routine 20% / complex 80% (escalation often reveals genuine complexity).

## Composed kernel Q = D·A·P

```
         looks_routine  looks_risky
Q = | 0.730         0.270       |  (from looks_routine)
    | 0.448 (≈)     0.552 (≈)   |  (from looks_risky)
```

An agent that sees `looks_routine` stays in `looks_routine` 73% of the next cycle. An agent that sees `looks_risky` returns to `looks_routine` 45% of the time — mostly by escalating or clarifying and having the world settle.

**Reading Q entries as compressed micro-paths:**
- `looks_routine → looks_routine`: task was simple, agent answered, world stayed simple, agent re-read as routine
- `looks_routine → looks_risky`: task was actually complex, agent answered directly, problem persisted, next cycle sees difficulty
- `looks_risky → looks_routine`: agent escalated/clarified, situation improved, re-read as routine
- `looks_risky → looks_risky`: problem remained hard even after intervention

## Analysis outputs

### Stationary distribution

```
π = [0.624  0.376]
```

In steady state the agent spends ~62% of time seeing `looks_routine` and ~38% seeing `looks_risky`. This reflects the combined effect of perception noise, decision policy, and action consequences — not the true task distribution (which is a separate input).

### Entropy rate

```
H(Q) ≈ 0.898 bits/step
```

Close to 1 bit/step — substantial uncertainty per cycle. Row entropies: H(row_0) ≈ 0.84, H(row_1) ≈ 0.99. Neither row is nearly deterministic, so the weighted average stays high. An entropy rate of 0 would mean a fully predictable, deterministic policy; 1 bit would mean the next experience is always a coin flip.

### Communicating classes

```
recurrent classes = [[0, 1]]
transient states  = []
```

One recurrent class, no transient states — the chain is **ergodic**. The stationary distribution is unique; long-run averages are independent of starting state.

## The cyclic kernels S and W

Q, S, and W are cyclic permutations of the same loop:

| Kernel | Starts at | Sequence        |
|--------|-----------|-----------------|
| Q      | X         | X → G → W → X  |
| S      | G         | G → W → X → G  |
| W      | W         | W → X → G → W  |

Their stationary distributions are related by one additional kernel application. See [[PDA Triplet Model]].

## Connection to math and API

- Q is computed via `Agent.QualiaKernel()` → [[catrace API]]
- Stationary via `Kernel.Stationary(1e-12, 5000)` → [[Markov Chain Foundations]]
- Entropy rate via `Kernel.EntropyRate(2)` → [[Markov Chain Foundations]]
- Ergodicity confirmed via `Kernel.Classes(1e-12)` → [[catrace API]]

## Sources

- `examples/simple_agent/WALKTHROUGH.md`
- `docs/patterns/story-single-llm-agent.md`
