# README scenario conventions

Every scenario in the `## Scenario write-ups` section of `README.md` follows
this structure. The README is for browsing — no raw output, no code snippets.
The goal is a reader who has never run the code can still understand what the
example models and why it matters.

See `docs/walkthrough-conventions.md` for the companion standard covering
`WALKTHROUGH.md` files.

---

## Required sections, in order

### Story

Plain English only. No state names in backticks, no kernel notation, no
matrix references. Just the situation, the tension, and why it is interesting.

The story opens with the situation before introducing any mechanism:

> A worker agent performs tasks while a validator agent monitors its health.
> Either agent may itself be functioning well or badly. When the validator is
> healthy, it can detect worker problems and attempt repairs…

If the scenario introduces a structural novelty (joint kernels, trace,
product state spaces), state it in plain English *after* the situation:

> Rather than writing down a 4×4 joint transition matrix by hand, each agent
> is first modeled independently…

### State meanings

One bullet per state space. Cover **all** state spaces for every agent —
world, experience, and action — plus joint world states for multi-agent
scenarios. Each bullet answers the question "what does this state mean in the
story?"

```
- worker world states: `worker_valid`, `worker_invalid` — is the worker actually functioning?
- worker experience states: `sees_ok`, `sees_problem` — does the worker detect its own degradation?
- worker actions: `produce`, `self_check`, `idle`
```

Not:
```
- `VV` = both agents reliable   ← too terse, no story connection
```

### Interpretation

One bullet per coupling point or structural claim. Each bullet explains
**what** the kernel does and **why** in story terms — not just its name.

```
- the validator's perception is coupled to the worker's world state: a degraded
  worker shifts the validator's experience toward `looks_bad` even if the
  validator itself is fine — this is where cross-agent observation enters
```

Not:
```
- coupling enters through P_joint   ← names the fact, doesn't explain it
```

### Code line

Single line linking to the example and its walkthrough:

```
Code: `examples/foo/main.go` — walkthrough at `examples/foo/WALKTHROUGH.md`
```

### Played-out version (required for implemented examples)

A concrete walk through three versions of the story: typical, failure, and
recovery. Each version steps through the full W→X→G→W cycle with explicit
probabilities and ends with a plain-English summary sentence.

**Structure of each version:**

1. Starting world state — name it, describe it in story terms
2. Perception step — name the experience state, give its probability
3. Decision step — name the action chosen, give its probability
4. Action effect — name the next world state, give its probability, show the
   product if it is the interesting number
5. One-sentence plain-English summary of what happened

After the three versions, a short "why this matters" paragraph explaining what
the paths together reveal that no single path shows alone.

Close with a **concise shorthand table** — four to six kernel entries mapped
to one-line plain-English readings.

---

## Style rules

- "Story:" as a standalone label before the story paragraph, not a header.
- State names always in backticks: `looks_routine`, `VV`.
- Kernel expressions in math or code style: $Q = DAP$, `J = P_joint · D_joint · A_joint`.
- No raw terminal output — that belongs in `WALKTHROUGH.md`.
- No code snippets — that belongs in `WALKTHROUGH.md`.
- Versions named **A / B / C**, bolded, with a short descriptive label:
  `**Version A: statistically typical path**`
- Plain-English summary sentence at the end of each version, introduced with
  "In plain English:".
- "Played-out version: Story N" as the subsection header (`####` level).
