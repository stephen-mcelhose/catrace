# Walkthrough: Prompt-chaining document pipeline

Prior examples either used one agent (`simple_agent`) or coupled peers that act
in the same timestep (`validator_repair`, `self_healing_nodes`). This example
is different: it models Anthropic-style **prompt chaining** — a fixed sequence
of specialised LLM calls with **programmatic gates** between them — by
assembling a pipeline world kernel from per-stage perception/decision maps
plus gate code, rather than one shared `Agent.WorldKernel()`.

Run:

```
go run examples/prompt_chaining/main.go
```

---

## The story

A diligence desk turns a raw filing into a client-ready report: extract claims,
summarise a brief, format the report. Each step is its own LLM call on the
artifact the previous gate accepted. Gates are ordinary checks (pass, retry
stage, or escalate to a human queue that may re-queue). Structural novelty:
only the **active** stage fires; the pipeline world records where the artifact
sits.

---

## State spaces

**Pipeline world:**

| Layer | States | Meaning |
|-------|--------|---------|
| World | `raw`, `extracted`, `summarised`, `formatted`, `failed` | Artifact position; `formatted` absorbing (shipped); `failed` human queue (may return to `raw`) |

**Extractor** (active at `raw`):

| Layer | States | Meaning |
|-------|--------|---------|
| Experience | `filing_clear`, `filing_noisy` | How readable the raw filing looks |
| Action | `emit_claims`, `retry_extract` | Produce claims vs reframe |

**Summariser** (active at `extracted`):

| Layer | States | Meaning |
|-------|--------|---------|
| Experience | `claims_clear`, `claims_noisy` | How usable gated claims look |
| Action | `emit_brief`, `retry_summarise` | Produce brief vs reframe |

**Formatter** (active at `summarised`):

| Layer | States | Meaning |
|-------|--------|---------|
| Experience | `brief_clear`, `brief_noisy` | How usable gated brief looks |
| Action | `emit_report`, `retry_format` | Produce report vs reframe |

**Gate** (code, not an agent): `pass` → next pipeline state; `retry_stage` → stay;
`escalate` → `failed`.

---

## Coupling

- **Handoff:** a gate `pass` advances the pipeline world; the next stage only
  ever sees the accepted intermediate, never the raw upstream filing.
- **Active stage only:** rows of `W` for `raw` / `extracted` / `summarised` are
  built from that stage's perception × decision × gate; inactive specialists
  do not act.
- **No cross-perception:** specialists do not observe each other's internal
  experience — coupling is through the artifact position and the gate, not
  $P_\text{joint}$ among peers.

---

## Math worked by hand

The extractor row of `W` when the filing is `raw`. Perception:
$P(\texttt{filing_clear})=0.55$, $P(\texttt{filing_noisy})=0.45$. Decision:

```
clear → emit 0.85, retry 0.15
noisy → emit 0.25, retry 0.75
```

Gate given emit: pass $0.70$, retry_stage $0.20$, escalate $0.10$.
After retry action: stay $0.90$, escalate $0.10$.

Advance probability `raw → extracted`:

```
0.55 × 0.85 × 0.70 + 0.45 × 0.25 × 0.70
= 0.32725 + 0.07875
= 0.406
```

Stay at `raw`:

```
0.55×0.85×0.20 + 0.45×0.25×0.20 + 0.55×0.15×0.90 + 0.45×0.75×0.90
= 0.0935 + 0.0225 + 0.07425 + 0.30375
= 0.494
```

Escalate to `failed`:

```
0.55×0.85×0.10 + 0.45×0.25×0.10 + 0.55×0.15×0.10 + 0.45×0.75×0.10
= 0.1
```

Check: $0.494 + 0.406 + 0.1 = 1$.

---

## Reading the output

```
=== Pipeline world kernel W (assembled from active stage · gate) ===
⎡  0.494    0.406        0        0     0.1  ⎤
⎢      0  0.332875     0.61        0  0.057125⎥
⎢      0        0   0.21746  0.74096  0.04158 ⎥
⎢      0        0        0        1        0 ⎥
⎣   0.25        0        0        0     0.75 ⎦

Stationary distribution π:
  raw          0.000000
  extracted    0.000000
  summarised   0.000000
  formatted    1.000000
  failed       0.000000

Entropy rate H(W): 0.000000 bits/step
Recurrent classes: [[3]]
Transient states:  [0 1 2 4]
MFPT raw → formatted: 7.6838 steps

=== Trace onto {raw, formatted, failed} ===
⎡ 0.494000  0.351509  0.154491 ⎤
⎢        0         1         0 ⎥
⎣     0.25         0      0.75 ⎦

IsTraceOf verification: true
Stationary of trace: raw=0.000000  formatted=1.000000  failed=0.000000
```

In the long run almost all mass sits in `formatted` — once a report ships, it
stays shipped, and the human queue eventually re-queues enough work that the
chain finishes. MFPT ≈ 7.7 steps is the expected desk latency from inbox to
shipped report under these gate and retry rates. The trace collapses
mid-pipeline stages so you see start / shipped / human-queue dynamics; from
`raw`, about 35% of the *effective* one-step mass in the trace picture heads
toward shipping paths versus ~15% toward the queue (with the rest retrying in
place), before intermediate stages are forgotten.

---

## What you can change

1. **Tighten the extractor gate** — raise `gatePass` from 0.70 toward 0.90.
   Watch MFPT fall and `W[raw, failed]` shrink. Meaning: fewer bad claims reach
   later stages, but you also reject more at the first link.
2. **Noisier filings** — lower extractor `pClear` from 0.55. Watch
   `W[raw, raw]` and MFPT rise. Meaning: upstream document quality dominates
   chain latency.
3. **Stricter human re-queue** — lower `failed → raw` from 0.25 toward 0.
   Watch MFPT blow up (and eventually become undefined if `failed` is fully
   absorbing while still competing with `formatted`). Meaning: a dead-letter
   queue that never returns work is a different reliability story than a
   retrying desk.
4. **Weaker summariser emit-when-clear** — drop summariser `D[clear, emit]`
   from 0.90. Watch mass linger in `extracted`. Meaning: a timid middle stage
   becomes the bottleneck even when gates are loose.
