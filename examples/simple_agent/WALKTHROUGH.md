# Walkthrough: Single LLM Task Agent

This file walks through the math in `main.go` step by step — what each matrix means,
how the composition works, and what the output numbers tell you.

Run the example to follow along:

```
go run examples/simple_agent/main.go
```

---

## The scenario

An LLM support agent handles tasks one at a time. At any moment:

- The **world** has a real state it cannot fully see
- The **agent** forms an internal interpretation (experience)
- Based on that, it **acts**
- The action changes the world, the agent re-observes, and the cycle repeats

The question the math answers: *in the long run, how does the agent's interpretation of tasks
evolve, step over step?*

---

## Step 1: Define the three maps

Each map is a **stochastic matrix** — a rectangular grid of probabilities where every row sums to 1.
Row `i` gives the probability distribution over output states when the input is state `i`.

### D — Decision kernel (X → G)

```go
D := mat.NewDense(2, 3, []float64{
    0.8, 0.2, 0.0,   // looks_routine -> answer 80%, clarify 20%, escalate 0%
    0.1, 0.3, 0.6,   // looks_risky   -> answer 10%, clarify 30%, escalate 60%
})
```

Rows are experience states, columns are actions:

```math
D = \begin{bmatrix} 0.8 & 0.2 & 0.0 \\ 0.1 & 0.3 & 0.6 \end{bmatrix}
\quad
\begin{array}{l} \leftarrow \text{looks\_routine} \\ \leftarrow \text{looks\_risky} \end{array}
```

A routine-looking task gets answered directly most of the time. A risky-looking task usually
gets escalated.

### A — Action effect kernel (G → W)

```go
A := mat.NewDense(3, 2, []float64{
    0.9, 0.1,   // answer   -> task_routine 90%, task_complex 10%
    0.4, 0.6,   // clarify  -> task_routine 40%, task_complex 60%
    0.2, 0.8,   // escalate -> task_routine 20%, task_complex 80%
})
```

Rows are actions, columns are world states:

```math
A = \begin{bmatrix} 0.9 & 0.1 \\ 0.4 & 0.6 \\ 0.2 & 0.8 \end{bmatrix}
\quad
\begin{array}{l} \leftarrow \text{answer} \\ \leftarrow \text{clarify} \\ \leftarrow \text{escalate} \end{array}
```

Answering directly usually stabilizes a routine task. Clarifying or escalating often reveals
the task is genuinely complex.

### P — Perception kernel (W → X)

```go
P := mat.NewDense(2, 2, []float64{
    0.85, 0.15,   // task_routine -> looks_routine 85%, looks_risky 15%
    0.25, 0.75,   // task_complex -> looks_routine 25%, looks_risky 75%
})
```

Rows are world states, columns are experience states:

```math
P = \begin{bmatrix} 0.85 & 0.15 \\ 0.25 & 0.75 \end{bmatrix}
\quad
\begin{array}{l} \leftarrow \text{task\_routine} \\ \leftarrow \text{task\_complex} \end{array}
```

The agent mostly reads the task correctly, but not always.

---

## Step 2: Compose Q = D · A · P

```go
Q, err := agent.QualiaKernel()
```

Q is the **qualia kernel** — the closed-loop experience-to-experience dynamics. It lives entirely
in the experience space `{looks_routine, looks_risky}` and answers: given how the agent
currently reads the task, what is the distribution over how it will read the task *next*, after
acting and re-observing?

Compute it in two stages.

### Stage 1: D · A  (2×3 · 3×2 = 2×2)

Each entry `(DA)_{ij}` sums over all possible actions — take row `i` from D, column `j` from A,
multiply element by element, add:

Row `looks_routine`:
```
→ task_routine:  0.8×0.9 + 0.2×0.4 + 0.0×0.2 = 0.72 + 0.08 + 0.00 = 0.80
→ task_complex:  0.8×0.1 + 0.2×0.6 + 0.0×0.8 = 0.08 + 0.12 + 0.00 = 0.20
```

Row `looks_risky`:
```
→ task_routine:  0.1×0.9 + 0.3×0.4 + 0.6×0.2 = 0.09 + 0.12 + 0.12 = 0.33
→ task_complex:  0.1×0.1 + 0.3×0.6 + 0.6×0.8 = 0.01 + 0.18 + 0.48 = 0.67
```

```math
D \cdot A = \begin{bmatrix} 0.80 & 0.20 \\ 0.33 & 0.67 \end{bmatrix}
```

### Stage 2: (D·A) · P  (2×2 · 2×2 = 2×2)

Row `looks_routine`:
```
→ looks_routine:  0.80×0.85 + 0.20×0.25 = 0.68 + 0.05 = 0.73
→ looks_risky:    0.80×0.15 + 0.20×0.75 = 0.12 + 0.15 = 0.27
```

Row `looks_risky`:
```
→ looks_routine:  0.33×0.85 + 0.67×0.25 = 0.2805 + 0.1675 = 0.448
→ looks_risky:    0.33×0.15 + 0.67×0.75 = 0.0495 + 0.5025 = 0.552
```

```math
Q = D \cdot A \cdot P = \begin{bmatrix} 0.73 & 0.27 \\ 0.448 & 0.552 \end{bmatrix}
```

This matches the program output:

```
Q = D*A*P
⎡ 0.730   0.270 ⎤
⎣ 0.448   0.552 ⎦
```

**Reading Q:**

| From \ To     | looks_routine | looks_risky |
|---------------|---------------|-------------|
| looks_routine | 0.73          | 0.27        |
| looks_risky   | 0.45          | 0.55        |

- An agent that sees a routine task will still see it as routine 73% of the time next cycle.
- An agent that sees a risky task will flip to routine 45% of the time — mostly by escalating or
  clarifying and having the world settle.

---

## Step 3: Stationary distribution

```go
pi, err := Q.Stationary(1e-12, 5000)
// stationary(Q) = 0.623955 0.376045
```

The stationary distribution π is the long-run fraction of steps the agent spends in each
experience state. It satisfies:

```math
\pi \cdot Q = \pi, \qquad \pi_0 + \pi_1 = 1
```

Running another step doesn't change the distribution — π is a fixed point of Q.

Solve by substituting Q:

```
π₀ × 0.73 + π₁ × 0.448 = π₀
              ↓
0.448 × π₁ = 0.27 × π₀
       π₁/π₀ = 0.27 / 0.448 ≈ 0.603
```

With π₀ + π₁ = 1:

```
π₀ = 1 / (1 + 0.603) ≈ 0.624
π₁ = 1 - 0.624       ≈ 0.376
```

```math
\pi = \begin{bmatrix} 0.624 & 0.376 \end{bmatrix}
```

**Interpretation:** In steady state, the agent experiences the task as routine about 62% of
the time and as risky about 38% of the time. This reflects the combined effect of the perception
noise, the decision policy, and the action consequences.

---

## Step 4: Entropy rate

```go
h, err := Q.EntropyRate(2)
// entropy_rate(Q) = 0.898142 bits/step
```

Entropy rate measures the average uncertainty per step:

```math
H(Q) = -\sum_i \pi_i \sum_j Q_{ij} \log_2 Q_{ij}
```

Per-row entropies (binary entropy of each row):

```
H(row 0) = -[0.73×log₂(0.73) + 0.27×log₂(0.27)] ≈ 0.84 bits
H(row 1) = -[0.45×log₂(0.45) + 0.55×log₂(0.55)] ≈ 0.99 bits
```

Weighted by stationary distribution:

```
H(Q) = 0.624 × 0.84 + 0.376 × 0.99 ≈ 0.90 bits/step
```

**Interpretation:** Close to 1 bit/step — there is substantial uncertainty in how the agent's
experience evolves each cycle. The two rows have similar entropy (neither transition is nearly
deterministic), so the weighting doesn't change the picture much.

An entropy rate of 0 would mean the agent's interpretation is fully predictable step-to-step
(a deterministic policy with perfect perception). A rate of 1 bit would mean the next experience
is a coin flip regardless of the current one.

---

## Step 5: Communicating classes

```go
classes, err := Q.Classes(1e-12)
// recurrent classes = [[0 1]]
// transient states  = []
```

Both experience states belong to one recurrent communicating class — the chain is **ergodic**.
From any starting experience, you can reach any other experience, and you never get trapped.

This means:
- the stationary distribution π is unique
- the long-run averages are the same regardless of starting state
- the entropy rate is a well-defined single number for the whole system

If there were transient states, the agent could start in an experience that it eventually leaves
and never returns to.

---

## Step 6: Left action (one-step forecast)

```go
next, err := Q.LeftAction([]float64{1, 0})
// left_action([1,0], Q) = 0.730000 0.270000
```

`LeftAction` applies a distribution over current states to Q and returns the distribution
one step later:

```math
\text{next} = \begin{bmatrix} 1 & 0 \end{bmatrix} \cdot Q = \begin{bmatrix} 0.73 & 0.27 \end{bmatrix}
```

Starting from 100% `looks_routine`, after one cycle the agent will be in `looks_routine` with
probability 0.73 and `looks_risky` with probability 0.27. This is just reading off the first
row of Q directly.

---

## The other kernels: S and W

```go
S, err := agent.StrategyKernel()   // S = A·P·D  (action space view)
W, err := agent.WorldKernel()      // W = P·D·A  (world space view)
```

Q, S, and W are **cyclic permutations of the same composition**. They describe the same closed
loop from three different starting points:

| Kernel | Starts at | Sequence        | Ends at |
|--------|-----------|-----------------|---------|
| Q      | X         | X → G → W → X  | X       |
| S      | G         | G → W → X → G  | G       |
| W      | W         | W → X → G → W  | W       |

Their stationary distributions are related by the same cycle — if you know π(Q), you can
recover the stationary distributions of S and W by one more application of D or P respectively.

---

## What you can change

To experiment with this example, try:

- **Make perception perfect** (`P = identity 2×2`): watch the entropy rate drop as the
  agent's interpretation tracks the world exactly.
- **Make D deterministic** (`looks_routine → always answer, looks_risky → always escalate`):
  watch how the stationary distribution shifts.
- **Add a third experience state** (`looks_neutral`): you'll need to grow D from 2 rows to 3
  and re-run the composition.
- **Make A more forgiving** (raise the `clarify → task_routine` entry): see how the stationary
  distribution shifts toward `looks_routine` and the entropy rate changes.

---

## Further reading

See also the project [Glossary](../../GLOSSARY.md) for definitions of all math terms used here.

- **Matrix multiplication** — [3Blue1Brown: Essence of Linear Algebra, Ch. 4](https://www.youtube.com/watch?v=XkY2DOUCWMU) — visual, no prior knowledge assumed. Directly relevant to Step 2.
- **Eigenvalues and eigenvectors** — [3Blue1Brown: Essence of Linear Algebra, Ch. 14](https://www.youtube.com/watch?v=PFDu9oVAE-g) — the stationary distribution in Step 3 is an eigenvector of Q with eigenvalue 1.
- **Markov chain notation and stationary distributions** — Levin & Peres, *Markov Chains and Mixing Times* (2nd ed.), §1.5, p. 9. Free PDF: https://pages.uoregon.edu/dlevin/MARKOV/markovmixing.pdf — the standard graduate reference; uses π throughout and defines `π = πP` exactly as we do here.
