# Walkthrough: Routing Agent

This example introduces the routing topology — a single agent that classifies
incoming work and dispatches it to the correct handler. Unlike the simple agent
(two world states) and the prompt-chaining pipeline (five stages), the routing
agent has a symmetric three-state world where misclassification creates feedback
loops: a wrongly-routed ticket re-enters the queue and may be misrouted again.
The structural novelty is using the world kernel to quantify misrouting latency
and the MFPT to measure how long cross-type cycles take.

---

## The story

A customer support router classifies tickets as billing, technical, or general
and dispatches them to the matching specialist. If it misclassifies, the
specialist cannot resolve the ticket and re-queues it — often as an apparently
different type. The key tension is classification accuracy versus the cost of a
wrong route: a confident wrong decision causes a longer misrouting loop than an
uncertain escalation to human triage.

---

## State spaces

**Routing agent:**

| Layer      | States                                                                           | Meaning                                    |
|------------|----------------------------------------------------------------------------------|--------------------------------------------|
| World      | `billing_ticket`, `technical_ticket`, `general_ticket`                           | True type of the incoming ticket           |
| Experience | `reads_billing`, `reads_technical`, `reads_general`                              | Router's perceived classification          |
| Action     | `route_billing`, `route_technical`, `route_general`, `escalate_human`            | Where the router sends the ticket          |

---

## Math worked by hand

The world kernel entry W[`billing_ticket`, `billing_ticket`] — probability of a
billing ticket remaining in the billing queue after one routing cycle:

```
W[billing, billing] = Σ_x Σ_g  P[billing, x] · D[x, g] · A[g, billing]

P[billing, reads_billing]=0.75, P[billing, reads_technical]=0.15, P[billing, reads_general]=0.10

From reads_billing:  0.80×0.70 + 0.05×0.15 + 0.05×0.10 + 0.10×0.33
                   = 0.560 + 0.0075 + 0.005 + 0.033 = 0.6055

From reads_technical: 0.05×0.70 + 0.80×0.15 + 0.05×0.10 + 0.10×0.33
                    = 0.035 + 0.120 + 0.005 + 0.033 = 0.193

From reads_general:  0.05×0.70 + 0.05×0.15 + 0.80×0.10 + 0.10×0.33
                   = 0.035 + 0.0075 + 0.080 + 0.033 = 0.1555

W[billing, billing] = 0.75×0.6055 + 0.15×0.193 + 0.10×0.1555
                    = 0.454125 + 0.028950 + 0.015550 = 0.498625
```

This matches the printed matrix entry (0.49863). The self-loop is ~0.50: a
billing ticket has roughly even odds of staying in the billing queue versus
being rerouted somewhere else within one cycle.

---

## Reading the output

```
=== World kernel W = P·D·A ===
⎡0.49862500000000004  0.28900000000000003  0.21237500000000004⎤
  ⎢             0.2305              0.53275  0.23675000000000002⎥
  ⎣           0.206125  0.26462500000000005              0.52925⎦
```

Rows are starting ticket type; columns are next type. The diagonal entries
(0.499, 0.533, 0.529) are all around 0.5 — tickets more likely than not to
stay in their correct type, but substantial cross-type flow. Technical tickets
have the highest self-retention (0.533) because P gives them the best
classification accuracy (0.80).

```
stationary(W)        = 0.304149  0.371701  0.324150
entropy_rate(W)      = 1.472504 bits/step
recurrent classes(W) = [[0 1 2]]
transient states(W)  = []
```

Technical tickets dominate the queue at 37.2% despite equal incoming rates —
because they are classified most accurately, they persist in their correct
state longer. Billing and general are roughly equal at ~30% and ~32%. A single
recurrent class confirms every ticket type is reachable from every other.
Entropy rate 1.47 bits/step is lower than the prompt-chaining pipeline (2.22)
because the router's behavior is more predictable: diagonal entries near 0.5
means most transitions are concentrated in two outcomes.

```
MFPT billing→technical  = 3.553381 steps
MFPT technical→general  = 4.365416 steps
```

A billing ticket takes on average 3.55 cycles to arrive in the technical queue —
through misclassification or escalation feedback. The technical→general path
takes 4.37 cycles, slightly longer because technical tickets are well-retained
and resist crossing to general. These numbers quantify misrouting latency
directly: reducing cross-type leakage in P will increase both MFPTs.

```
=== Trace onto {billing_ticket, technical_ticket} — collapse general ===
⎡ 0.591616602496017  0.4083833975039831⎤
  ⎣0.3341645645246947  0.6658354354753053⎦

IsTraceOf                              = true
stationary(W)|{billing,technical} norm = 0.450024  0.549976
stationary(trace)                      = 0.450024  0.549976
```

Treating general enquiries as a hidden process gives the effective
billing/technical two-state picture: billing retains ~59% of cycles (within
the coarse view), technical retains ~67%. The normalized stationary puts
technical at 55% vs billing at 45% — reflecting technical's higher accuracy
and self-retention. `IsTraceOf = true` confirms the 2-state kernel is exact.

---

## What you can change

**Improve classification accuracy.** Increase the diagonal entries of P (e.g.,
`P[billing, reads_billing]` from 0.75 to 0.90). Watch the self-loop entries
in W rise and the stationary distribution become more uniform (better routing
reduces the technical dominance).

**Raise the escalation rate.** In D, increase `escalate_human` from 0.10 to
0.30 across all perceived types. The human resolves all ticket types equally,
so increasing escalation drives the stationary distribution toward uniform.
MFPT cross-type paths will shorten.

**Model an asymmetric queue.** Set `P[billing, reads_technical]` high (0.40)
while keeping others accurate. Billing tickets frequently look like technical
ones; watch billing disappear from the stationary distribution as they migrate
into the technical queue.

**Change the A kernel to make wrong routes unrecoverable.** Set
`A[route_billing, technical_ticket]` to 0 and redistribute. Now a billing
ticket routed to technical never re-enters as technical — watch the world
kernel lose its cross-type paths.

**Trace onto a single ticket type.** Use `Trace([]int{1}, tol)` to collapse
billing and general entirely. The 1×1 trace is just the probability of
staying in `technical_ticket` per cycle — the diagonal of the trace gives this
directly.
