# Walkthrough: Prompt-Chaining Pipeline

This example introduces a single-agent model with five world states, extending
the two-state world of the simple agent example. The key novelty is using the
world kernel to represent a multi-stage pipeline with retry loops and a failure
exit, then tracing onto the end states (`raw`, `formatted`, `failed`) to read
the end-to-end pass/fail rate directly.

---

## The story

A document intelligence pipeline sends text through three stages — extractor,
summariser, formatter — with a quality gate between each. A stage agent
perceives input quality and decides to process, retry, or escalate. The
structural novelty here is that tracing collapses the intermediate pipeline
stages and gives you the effective end-to-end dynamics as a 3-state system.

---

## State spaces

**Pipeline agent:**

| Layer      | States                                              | Meaning                                    |
|------------|-----------------------------------------------------|--------------------------------------------|
| World      | `raw`, `extracted`, `summarised`, `formatted`, `failed` | True pipeline stage of the document    |
| Experience | `input_clear`, `input_noisy`                        | Stage agent's read of input quality        |
| Action     | `process`, `retry`, `escalate`                      | Transform, reframe and retry, or abandon   |

---

## Math worked by hand

The world kernel entry W[`raw`, `raw`] — probability of staying at the raw
stage after one step — is computed by summing over all experience-action paths:

```
W[raw, raw] = Σ_x Σ_g  P[raw, x] · D[x, g] · A[g, raw]

P[raw, input_clear] = 0.25,  P[raw, input_noisy] = 0.75

From input_clear: D[clear, process]·A[process, raw] + D[clear, retry]·A[retry, raw] + D[clear, escalate]·A[escalate, raw]
                = 0.80 × 0.05 + 0.15 × 0.30 + 0.05 × 0.05
                = 0.040 + 0.045 + 0.0025 = 0.0875

From input_noisy: 0.25×0.05 + 0.50×0.30 + 0.25×0.05
                = 0.0125 + 0.150 + 0.0125 = 0.175

W[raw, raw] = 0.25 × 0.0875 + 0.75 × 0.175
            = 0.021875 + 0.13125 = 0.153125
```

This matches the printed matrix exactly.

---

## Reading the output

```
=== World kernel W = P·D·A ===
⎡           0.153125              0.29125             0.228125               0.1275                  0.2⎤
  ⎢0.12250000000000001  0.29650000000000004               0.2675                0.166  0.14750000000000002⎥
  ⎢0.10500000000000001  0.29950000000000004  0.29000000000000004  0.18800000000000006  0.11750000000000002⎥
  ⎢            0.09625  0.30100000000000005              0.30125  0.19900000000000004  0.10250000000000002⎥
  ⎣0.16187500000000002              0.28975             0.216875               0.1165  0.21500000000000002⎦
```

Rows are starting stage; columns are next stage (raw, extracted, summarised,
formatted, failed). Notice that no row has a high self-loop — the pipeline
always moves. `formatted` row (row 4): 0.096 chance of immediately moving back
to `extracted`, meaning a formatted document re-enters the pipeline as new input.

```
stationary(W)        = 0.123115  0.296395  0.266709  0.165227  0.148554
entropy_rate(W)      = 2.221416 bits/step
recurrent classes(W) = [[0 1 2 3 4]]
transient states(W)  = []
```

The pipeline spends most of its time in the middle stages (`extracted`: 29.6%,
`summarised`: 26.7%), with `formatted` at 16.5% and `failed` at 14.9%. A single
recurrent class confirms ergodicity — every stage is reachable from every other,
including recovery from failure. Entropy rate of 2.22 bits/step reflects high
variability: the pipeline is genuinely hard to predict step-to-step.

```
=== Trace onto {raw, formatted, failed} — end-to-end view ===
⎡0.29317501005976243   0.3432689011758745   0.3635560887643632⎤
  ⎢ 0.2582527310392107   0.4504090970059166  0.29133817195487266⎥
  ⎣ 0.2985476683706165   0.3267857941250987  0.37466653750428475⎦

IsTraceOf                        = true
stationary(W)|{raw,fmt,fail} norm= 0.281795  0.378183  0.340022
stationary(trace)               = 0.281795  0.378183  0.340022
```

Collapsing `extracted` and `summarised` into implicit hidden states reveals the
coarse pipeline health: 37.8% of end-state visits are `formatted` (success),
34.0% are `failed`. `IsTraceOf = true` confirms the 3-state trace is the exact
effective kernel — not an approximation. The stationary of the trace matches the
parent's stationary restricted and normalised to {`raw`, `formatted`, `failed`}.

---

## What you can change

**Lower the escalation rate.** In D, reduce `D[input_noisy, escalate]` from
0.25 to 0.05 and increase `process` or `retry`. Watch `failed` in the
stationary distribution fall and `formatted` rise.

**Make raw input harder to read.** In P, reduce `P[raw, input_clear]` from
0.25 to 0.10. More raw documents will be perceived as noisy, driving more
escalations early and increasing `failed`.

**Add a strong retry effect.** In A, increase `A[retry, extracted]` and reduce
`A[retry, raw]`. Retry becomes a genuine advancement action rather than a
near-stall; watch the mid-stage stationary mass shift toward `summarised`.

**Raise the gate quality.** In A, increase `A[process, formatted]` and
decrease `A[process, extracted]`. Fewer steps are needed to reach `formatted`;
the pipeline becomes faster and more decisive.

**Tighten the trace to just {formatted, failed}.** Change `endStates` to
`[]int{3, 4}`. The 2-state trace gives a single number — the pass rate — as
a row-stochastic 2×2 matrix.
