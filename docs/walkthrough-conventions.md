# Walkthrough conventions

Every example under `examples/` must have a `WALKTHROUGH.md`. This document
defines what goes in it and what goes elsewhere.

---

## Division of labour

| Layer | Content |
|-------|---------|
| `README.md` | Story in plain English, state meanings, played-out paths, conceptual interpretation |
| `WALKTHROUGH.md` | Context paragraph, state space tables, coupling explanation, math by hand, raw terminal output, what you can change |

Do not duplicate prose between the two. The README is for browsing; the
walkthrough is for someone with the code running in front of them.

---

## Required sections, in order

### 1. Opening paragraph

One short paragraph situating this example relative to the others — what
problem it introduces that the previous examples did not cover. This is
positioning, not story.

Example from `validator_repair`:

> This example fills the gap between the first two examples. The first example
> modeled a single agent with P, D, A kernels. The second applied trace to a
> joint kernel given directly. This example shows how to **build** a joint
> kernel from two independently modeled agents — and why it matters how you
> build it.

### 2. The story

Two to four sentences only. Enough to orient a reader who has not yet read the
README, plus one sentence naming the structural novelty this example
introduces. Do not retell the README story in full.

### 3. State space tables

One table per agent, plus a joint world states table for multi-agent examples.
Use this format per agent:

**Agent name:**

| Layer      | States                           | Meaning                          |
|------------|----------------------------------|----------------------------------|
| World      | `healthy`, `degraded`, ...       | What is actually true            |
| Experience | `ema_low`, `ema_mid`, ...        | What the agent perceives         |
| Action     | `push`, `throttle`, `idle`       | What the agent does              |

For multi-agent examples, follow with a joint world states table:

| State | Node       | Evolver        | Meaning                     |
|-------|------------|----------------|-----------------------------|
| H·G   | healthy    | good strategy  | Node fine, good config      |
| …     | …          | …              | …                           |

### 4. Coupling

*Multi-agent examples only.* One bullet per coupling point explaining which
kernel (P, D, or A) carries the coupling and what it means in story terms.
Single-agent examples omit this section entirely.

### 5. Math worked by hand

For at least one key result — a kernel entry, the stationary distribution, or
an entropy rate — show the arithmetic explicitly. Do not just restate the
formula. Compute a number.

Example from `simple_agent`:

```
π₀ × 0.73 + π₁ × 0.448 = π₀
0.448 × π₁ = 0.27 × π₀
π₁/π₀ = 0.27 / 0.448 ≈ 0.603
```

### 6. Reading the output

Show the **raw terminal output** for each analysis result, then interpret it
in one or two plain-English sentences. Show numbers exactly as the program
prints them — do not paraphrase.

```
stationary(J):
  H·G    0.5240
  H·B    0.1813
  ...
```

The system spends 52% of its time in the best joint state…

### 7. What you can change

Three to five concrete experiments the reader can run by editing the matrices.
For each: what to change, what to watch, and what it means if the number moves
in the expected direction.

---

## Style rules

- Use `---` to separate major sections.
- State names always in backticks: `healthy`, `ema_low`, `push`.
- Kernel names in code or math style: `P_joint`, $Q = DAP$.
- Tables for state spaces and output comparisons; prose for interpretation.
- Never explain what a stochastic matrix is — link to `GLOSSARY.md` instead.
- "The story" header, not "The scenario".
