# Walkthrough conventions

Every example under `examples/` must have a `WALKTHROUGH.md`. This document
defines what goes in it and what goes elsewhere.

---

## Division of labour

| Layer | Content |
|-------|---------|
| `README.md` | Story in plain English, state meanings, played-out paths, conceptual interpretation |
| `WALKTHROUGH.md` | State space table, math worked by hand, raw terminal output, what you can change |

Do not duplicate prose between the two. The README is for browsing; the
walkthrough is for someone with the code running in front of them.

---

## Required sections, in order

### 1. Opening paragraph

One short paragraph situating this example relative to the others — what
problem it introduces that the previous examples did not cover.

### 2. The story

Plain English. No state names, no kernel notation. Just the situation.

### 3. State space table

A table covering **all** state spaces for every agent in the example.

| State | Agent | Layer | Meaning |
|-------|-------|-------|---------|
| `healthy` | Node | World | Node is actually functioning |
| `ema_low` | Node | Experience | Error rate EMA below 10% |
| `push` | Node | Action | Node works at full rate |
| … | … | … | … |

For multi-agent examples, also include the joint world states.

### 4. Coupling

Explain where and how agents are coupled — which kernel (P, D, or A) carries
the coupling and what it means in story terms. One bullet per coupling point.

### 5. Math worked by hand

For at least one key kernel entry or analysis result, show the arithmetic
explicitly. Do not just restate the formula — compute a number.

Example from `simple_agent`:

```
π₀ × 0.73 + π₁ × 0.448 = π₀
0.448 × π₁ = 0.27 × π₀
π₁/π₀ = 0.27 / 0.448 ≈ 0.603
```

### 6. Reading the output

Show the **raw terminal output** for each analysis result, then interpret it
in one or two plain-English sentences. Do not paraphrase the numbers — show
them exactly as the program prints them.

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
the expected direction.

---

## Style rules

- Use `---` to separate major sections.
- State names always in backticks: `healthy`, `ema_low`, `push`.
- Kernel names in code or math style: `P_joint`, $Q = DAP$.
- Tables for state spaces and output comparisons; prose for interpretation.
- Never explain what a stochastic matrix is — link to `GLOSSARY.md` instead.
- "The story" header, not "The scenario".
