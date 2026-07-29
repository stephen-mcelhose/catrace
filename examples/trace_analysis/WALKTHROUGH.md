# Walkthrough: LLM Agent with Hidden Support System

This file walks through the math in `main.go` step by step — what the parent kernel represents,
what the trace operation actually computes, and what the output numbers tell you.

Run the example to follow along:

```
go run examples/trace_analysis/main.go
```

---

## The scenario

We can observe one focal LLM agent (A), but behind it is a support system (B) we cannot see —
retrieval backends, monitoring tools, human reviewers, other agents. The full system has four
states: the focal agent can be valid or invalid, and so can the hidden system.

We only see whether A is valid or invalid. We never see B directly. The question the math
answers: *what are the effective health dynamics of A alone, after accounting for all the hidden
influence of B?*

---

## Step 1: The parent kernel (4 states)

```go
parent, err := catrace.NewKernel(mat.NewDense(4, 4, []float64{
    0.60, 0.20, 0.10, 0.10,   // A_valid   -> ...
    0.15, 0.55, 0.15, 0.15,   // A_invalid -> ...
    0.20, 0.20, 0.40, 0.20,   // B_valid   -> ...
    0.10, 0.20, 0.20, 0.50,   // B_invalid -> ...
}), []string{"A_valid", "A_invalid", "B_valid", "B_invalid"})
```

States indexed 0–3. Columns have the same order as rows.

```math
L = \begin{bmatrix} 0.60 & 0.20 & 0.10 & 0.10 \\ 0.15 & 0.55 & 0.15 & 0.15 \\ 0.20 & 0.20 & 0.40 & 0.20 \\ 0.10 & 0.20 & 0.20 & 0.50 \end{bmatrix}
\quad
\begin{array}{l} \leftarrow \text{A\_valid} \\ \leftarrow \text{A\_invalid} \\ \leftarrow \text{B\_valid} \\ \leftarrow \text{B\_invalid} \end{array}
```

Every row sums to 1. Reading row `A_valid`: when the focal agent is currently healthy, next
step it stays healthy 60% of the time, degrades 20%, or the system transitions into B states
the other 20%.

Output:
```
Parent kernel
⎡ 0.6   0.2   0.1   0.1 ⎤
⎢ 0.15  0.55  0.15  0.15⎥
⎢ 0.2   0.2   0.4   0.2 ⎥
⎣ 0.1   0.2   0.2   0.5 ⎦
```

---

## Step 2: What the trace is NOT

The naive approach would be to delete rows and columns for B_valid and B_invalid:

```
WRONG — just cutting rows/cols:
              A_valid  A_invalid
A_valid   [   0.60      0.20   ]   ← rows no longer sum to 1
A_invalid [   0.15      0.55   ]
```

These rows don't sum to 1, so this isn't a valid kernel. More importantly, it ignores what
happens when the chain visits B states: the hidden system may modify A's state before we see
A again.

The trace accounts for all paths through B states between consecutive observations of A.

---

## Step 3: The trace formula

```go
subset := []int{0, 1}
trace, err := parent.Trace(subset, 1e-12)
```

Let S = `{A_valid, A_invalid}` (indices 0, 1) and B = `{B_valid, B_invalid}` (indices 2, 3).

Partition the parent kernel into four blocks:

```
L = [ L_SS  |  L_SB ]
    [ L_BS  |  L_BB ]
```

From the parent matrix:

```math
L_{SS} = \begin{bmatrix} 0.60 & 0.20 \\ 0.15 & 0.55 \end{bmatrix} \quad (A \to A \text{ direct})
\qquad
L_{SB} = \begin{bmatrix} 0.10 & 0.10 \\ 0.15 & 0.15 \end{bmatrix} \quad (A \to B)
```

```math
L_{BS} = \begin{bmatrix} 0.20 & 0.20 \\ 0.10 & 0.20 \end{bmatrix} \quad (B \to A)
\qquad
L_{BB} = \begin{bmatrix} 0.40 & 0.20 \\ 0.20 & 0.50 \end{bmatrix} \quad (B \to B)
```

The trace formula sums all paths that start in S, bounce through B any number of times, then
return to S:

```math
\text{Tr}(L) = L_{SS} + L_{SB}(I - L_{BB})^{-1}L_{BS}
```

The `(I - L_BB)⁻¹` term is a geometric series — it captures one B-hop, two B-hops, three
B-hops, and so on, all at once:

```math
(I - L_{BB})^{-1} = I + L_{BB} + L_{BB}^2 + L_{BB}^3 + \cdots
```

> **References:**
> - Matrix inverses: [3Blue1Brown: Essence of Linear Algebra, Ch. 7](https://www.youtube.com/watch?v=uQhTuRlWMxw)
> - Markov chain notation and stationary distributions: Levin & Peres, *Markov Chains and Mixing Times* (2nd ed.), §1.5, p. 9. Free PDF: https://pages.uoregon.edu/dlevin/MARKOV/markovmixing.pdf

### Computing (I - L_BB)⁻¹

```math
I - L_{BB} = \begin{bmatrix} 0.60 & -0.20 \\ -0.20 & 0.50 \end{bmatrix}
```

det(I − L_BB) = 0.60 × 0.50 − 0.20 × 0.20 = 0.30 − 0.04 = 0.26

```math
(I - L_{BB})^{-1} = \frac{1}{0.26} \begin{bmatrix} 0.50 & 0.20 \\ 0.20 & 0.60 \end{bmatrix} = \begin{bmatrix} 1.923 & 0.769 \\ 0.769 & 2.308 \end{bmatrix}
```

### Computing L_SB · (I - L_BB)⁻¹

```math
\begin{bmatrix} 0.10 & 0.10 \\ 0.15 & 0.15 \end{bmatrix} \cdot \begin{bmatrix} 1.923 & 0.769 \\ 0.769 & 2.308 \end{bmatrix} = \begin{bmatrix} 0.269 & 0.308 \\ 0.404 & 0.462 \end{bmatrix}
```

Row A_valid:   0.10×1.923 + 0.10×0.769 = 0.1923 + 0.0769 = 0.269
               0.10×0.769 + 0.10×2.308 = 0.0769 + 0.2308 = 0.308

Row A_invalid: 0.15×1.923 + 0.15×0.769 = 0.2885 + 0.1154 = 0.404
               0.15×0.769 + 0.15×2.308 = 0.1154 + 0.3462 = 0.462

### Computing · L_BS

```math
\begin{bmatrix} 0.269 & 0.308 \\ 0.404 & 0.462 \end{bmatrix} \cdot \begin{bmatrix} 0.20 & 0.20 \\ 0.10 & 0.20 \end{bmatrix} = \begin{bmatrix} 0.085 & 0.115 \\ 0.127 & 0.173 \end{bmatrix}
```

Row A_valid:   0.269×0.20 + 0.308×0.10 = 0.0538 + 0.0308 = 0.085
               0.269×0.20 + 0.308×0.20 = 0.0538 + 0.0615 = 0.115

Row A_invalid: 0.404×0.20 + 0.462×0.10 = 0.0808 + 0.0462 = 0.127
               0.404×0.20 + 0.462×0.20 = 0.0808 + 0.0923 = 0.173

### Final trace

```math
\text{Tr}(L) = L_{SS} + L_{SB}(I - L_{BB})^{-1}L_{BS} = \begin{bmatrix} 0.60 & 0.20 \\ 0.15 & 0.55 \end{bmatrix} + \begin{bmatrix} 0.085 & 0.115 \\ 0.127 & 0.173 \end{bmatrix} = \begin{bmatrix} 0.685 & 0.315 \\ 0.277 & 0.723 \end{bmatrix}
```

This matches the output exactly:

```
Trace kernel on subset {0,1}
⎡ 0.6846   0.3154 ⎤
⎣ 0.2769   0.7231 ⎦
```

**Reading the trace:**

| From \ To | A_valid | A_invalid |
|-----------|---------|-----------|
| A_valid   | 0.685   | 0.315     |
| A_invalid | 0.277   | 0.723     |

Compare to the naive cut `[0.60, 0.20 / 0.15, 0.55]`: every entry has shifted. The hidden
support system makes A_valid more likely to stay valid (0.685 vs 0.60) and A_invalid more
likely to stay invalid (0.723 vs 0.55) — the hidden system amplifies existing states.

---

## Step 4: Verifying the trace is correct

```go
ok, err := trace.IsTraceOf(parent, subset, 1e-12)
// IsTraceOf = true
```

`IsTraceOf` checks the defining property: the stationary distribution of the trace kernel,
normalized over the subset, must equal the stationary distribution of the parent kernel
restricted to the same subset.

```go
piParent, err := parent.Stationary(1e-12, 5000)
piTrace,  err := trace.Stationary(1e-12, 5000)
restricted, _ := catrace.RestrictDistribution(piParent, subset, 1e-12)

// stationary(parent)|subset normalized = 0.467532 0.532468
// stationary(trace)                    = 0.467532 0.532468
```

They match to six decimal places. This is the key invariant: the trace gives the right
long-run frequencies for the observable states, even though the hidden states are gone.

---

## Step 5: Sampling and estimation

```go
rng := rand.New(rand.NewSource(7))
trajectory := make([]int, 100)
trajectory[0] = 0
for i := 1; i < len(trajectory); i++ {
    next, _ := parent.Sample(trajectory[i-1], rng)
    trajectory[i] = next
}
```

This simulates 100 steps of the full 4-state parent chain starting from A_valid.

```go
subsetSet := map[int]bool{0: true, 1: true}
traceSeq := catrace.SampleTraceFromSequence(trajectory, subsetSet)
```

`SampleTraceFromSequence` extracts only the steps where the chain is in `{A_valid, A_invalid}`,
in order. Steps through B states are discarded — all we keep are the A-state returns.

```go
est, _ := catrace.EstimateKernelFromSequence(traceSeq, 2, 1e-3)
```

This counts consecutive pairs `(A_i → A_j)` in the observed trace sequence and forms a
frequency estimate.

Output:
```
Estimated trace kernel from sampled sequence
⎡ 0.429   0.571 ⎤
⎣ 0.355   0.645 ⎦
```

With only 100 steps of the parent chain — and many of those spent in B states — the trace
sequence is short and noisy. The estimated kernel is in the right ballpark but not close to
the true trace. With 10,000+ steps it converges.

```go
windows, _ := catrace.WindowedTraceEstimates(trajectory, subsetSet, 25, 10, 1e-3)
// windowed estimates computed: 8
```

`WindowedTraceEstimates` cuts the trajectory into overlapping windows of 25 steps, each offset
by 10, and estimates the trace kernel in each window. The 8 windows each give a local estimate —
useful for checking whether the transition probabilities are stable over time or drifting.

---

## Further reading

See also the project [Glossary](../../GLOSSARY.md) for definitions of all math terms used here.

---

## What you can change

To experiment:

- **Strengthen B's influence on A**: raise entries in `L_BS` (e.g. `B_valid → A_valid = 0.5`).
  Watch the trace differ more from the naive cut.
- **Isolate B from A**: set `L_SB` entries to 0 (no path from A into B). The trace collapses
  to the naive cut `L_SS` because there are no hidden excursions.
- **Longer trajectory**: increase the trajectory length from 100 to 10,000. Watch the estimated
  kernel converge to the true trace.
- **Add a third observed state**: extend S to `{A_valid, A_invalid, A_warning}`. The trace
  computation scales — the formula is the same, the matrices are just larger.
