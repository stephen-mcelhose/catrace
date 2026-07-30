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

Plain English only in the opening paragraph — no state names in backticks, no
kernel notation. Just the situation and the tension.

If the scenario introduces a structural novelty (joint kernels, trace, product
state spaces), state it in plain English *after* the situation paragraph:

> Rather than writing down a 4×4 joint transition matrix by hand, each agent
> is first modeled independently…

### State meanings

One bullet per state space. Cover **all** state spaces for every agent —
world, experience, and action — plus joint world states for multi-agent
scenarios. Each bullet names the states and answers the question "what does
this mean in the story?"

```
- worker world states: `worker_valid`, `worker_invalid` — is the worker actually functioning?
- worker experience states: `sees_ok`, `sees_problem` — does the worker detect its own degradation?
- worker actions: `produce`, `self_check`, `idle`
```

### Interpretation

One bullet per structural claim. Each bullet explains **what** a kernel does
and **why** in story terms — not just its name.

For single-agent examples, bullets describe what each kernel (P, D, A) captures:

```
- perception captures imperfect prompt interpretation
- decision captures the LLM policy
- action captures how the chosen response changes the real task situation
```

For multi-agent examples, bullets describe coupling points:

```
- the validator's perception is coupled to the worker's world state: a degraded
  worker shifts the validator's experience toward `looks_bad` even if the
  validator itself is fine — this is where cross-agent observation enters
```

### Code line

Single line linking to the example and its walkthrough:

```
Code: `examples/foo/main.go` — walkthrough at `examples/foo/WALKTHROUGH.md`
```

### Played-out version (required for implemented examples)

A concrete walk through three versions of the story: typical, failure, and
recovery. Each version steps through the full W→X→G→W cycle with explicit
probabilities.

**Structure of each version:**

Four numbered steps:

1. Starting world state — name it, describe it in story terms
2. Perception step — name the experience state, give its probability
3. Decision step — name the action chosen, give its probability
4. Action effect — name the next world state, give its probability; show the
   product if it is the interesting number

Then two unnumbered sentences outside the list:

- "This path contributes [the bulk of / to] the [transition name] in [kernel]."
- "In plain English: [one sentence describing what happened in story terms]."

After the three versions, a short "why this matters" paragraph explaining what
the paths together reveal that no single path shows alone.

Close with a **concise shorthand list** — four to six kernel entries as bullets,
each mapped to a one-line plain-English reading:

```
- `looks_routine -> looks_routine` — the agent correctly handled a manageable task
- `looks_routine -> looks_risky` — the task was harder than it first appeared
```

---

## Style rules

- "Story:" as a standalone label before the story paragraph, not a markdown header.
- State names in backticks in state meanings and interpretation: `looks_routine`, `VV`.
- No state names in backticks inside the story paragraph itself.
- Kernel expressions in math or code style: $Q = DAP$, `J = P_joint · D_joint · A_joint`.
- No raw terminal output — that belongs in `WALKTHROUGH.md`.
- No code snippets — that belongs in `WALKTHROUGH.md`.
- Versions named **A / B / C**, bolded, with a short descriptive label:
  `**Version A: statistically typical path**`
- "In plain English:" introduces the summary sentence at the end of each version.
- "Played-out version: Story N" as the subsection header (`####` level).
