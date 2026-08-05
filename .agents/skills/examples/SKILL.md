---
name: examples
description: >
  Deliver catrace pattern examples end-to-end: critically sharpen the story to
  "real-ish", implement examples/*/main.go, write WALKTHROUGH.md and README
  scenario entries, sync Scenario Registry / pattern reference / wiki, lint
  Planned/Draft/Implemented drift, and close example: GitHub issues. Use when
  the user says "work an example", "start example", "implement example",
  "finish prompt_chaining", "example issue", "lint examples", "close example",
  "real-ish story", or when picking up any open issue titled "example: …".
  Examples teach the API; experiments ask falsifiable design questions — use
  the experiments skill for hypothesis registry work, not this one.
license: MIT
metadata:
  version: "1.0.0"
---

# Examples delivery

catrace examples are runnable teaching demos under `examples/<name>/`.
They are distinct from experiments: examples teach the API; experiments ask a
falsifiable design question (see `.agents/skills/experiments/SKILL.md`).

**Core duty of this skill:** take an `example:` GitHub issue from story sketch
to Implemented — with a mandatory **real-ish** story critique before coding —
then keep Scenario Registry / pattern tables / wiki honest.

Do **not** duplicate walkthrough or README section recipes here. Point at the
convention docs and enforce the Done surface.

## Status vocabulary

| Status | Means |
|--------|--------|
| **Planned** | Open `example:` issue + story file; no `examples/<name>/main.go` (or only stubs that never ran) |
| **Draft** | `main.go` exists but missing `WALKTHROUGH.md` and/or README scenario flip |
| **Implemented** | `main.go` **and** `WALKTHROUGH.md` **and** README scenario entry (story + played-out paths) on the branch that will merge |

Hard rule (same as wiki Scenario Registry): never call an example Implemented
without both `main.go` and `WALKTHROUGH.md`.

## When to use

| Trigger | Operation |
|---------|-----------|
| "work / start / pick up example #N" (or named pattern) | **start** — then stop for human accept unless user said continue |
| Story accepted; build the demo | **implement** |
| PR / merge ready; Done surface | **close** |
| "lint examples", registry drift, after wiki lint touches examples | **lint** |
| New pattern candidate needs issue + story stub | **file** |

## Canonical locations

| Path | Role |
|------|------|
| `docs/patterns/story-*.md` | Story + state meanings + interpretation; links the tracking issue |
| `docs/patterns/agentic-patterns-reference.md` | Pattern table; `catrace example` column |
| GitHub issues titled `example: …` | Tracking; body has Pattern, state spaces, features, AC checkboxes |
| `examples/<name>/main.go` | Runnable demo |
| `examples/<name>/WALKTHROUGH.md` | Hands-on narrative — follow `docs/walkthrough-conventions.md` |
| `README.md` scenario write-ups | Browsing layer — follow `docs/readme_scenario_conventions.md` |
| `docs/wiki/scenario-registry.md` | Status map README ↔ story ↔ code ↔ wiki |
| `docs/wiki/AGENTS.md` | Wiki ingest rules after raw sources change |
| `plans/*.md` | Optional; use when start needs a longer real-ish plan (e.g. multi-agent) |

---

## Operation: start (REQUIRED before first implement)

Human gate: after **start**, present the critique + proposed story edits and
**wait for accept** unless the user explicitly said to continue into implement.

### 1. Load context

1. Resolve issue (`gh issue view N`) and linked story path.
2. Read story, pattern-reference row, any existing `examples/<name>/`, and
   Scenario Registry / README scenario number claims.
3. Skim one nearby Implemented example for tone (`simple_agent`,
   `validator_repair`, or `self_healing_nodes`).

### 2. Real-ish critique (try to break the sketch)

Pressure-test the story. Fail the gate if answers are weak:

| # | Question | Fail if… |
|---|----------|----------|
| 1 | **Mechanism, not costume** — Can every coupling / transition be said as “when X happens, Y’s next world shifts *because* …”? | Coupling is only ±noise with a job-title label |
| 2 | **Observable tension** — What question does a reader ask after running it (MFPT, Trace collapse, who carries recovery)? | Only answer is “demo of pattern N” |
| 3 | **Perception honesty** — What can each agent *not* see? | Hidden state / partial observation omitted for convenience |
| 4 | **Non-goals** — What is parked (extra agent, evolver, heterogeneity)? | v1 quietly includes three theses |
| 5 | **State-space fit** — After rewrite, do W / X / G sets still match? | Preserving sketch states out of loyalty |
| 6 | **AC sync** — Do issue checkboxes and “features exercised” match the revised story? | Coding against stale AC |

Write a short critique block for the user:

```
## Example start — <name> (#N)
### Critique
- …
### Real-ish scenario (proposed)
…
### Non-goals
- …
### State spaces (revised if needed)
- …
### AC deltas (issue body)
- …
### Ready for implement? (awaiting accept)
```

### 3. Apply story adjustments (after accept, or immediately if user said continue)

1. Update `docs/patterns/story-*.md` (story paragraph, state meanings,
   interpretation). Keep the Issue footer link.
2. Update the GitHub issue body if AC / state spaces / features changed
   (`gh issue edit`).
3. Optional: write or trim `plans/<name>.md` when the real-ish rewrite is large
   (Y-statement, non-goals, state tables) — same spirit as
   `plans/network-of-healers.md`, not required for every single-agent pipeline.
4. Do **not** mark Implemented; status stays Planned or Draft.

---

## Operation: implement

Prerequisites: **start** completed (or user explicitly skipped with a
real-ish story already in place).

1. Implement `examples/<name>/main.go` against the **revised** AC.
2. Prefer patterns from existing examples: named states printed, `Trace` /
   `Stationary` / `MeanFirstPassage` / entropy as the issue requires, clear
   stdout a walkthrough can quote verbatim.
3. Write `WALKTHROUGH.md` per `docs/walkthrough-conventions.md` (do not invent
   a different section order).
4. Add or update the README scenario entry per
   `docs/readme_scenario_conventions.md` (story, state meanings, interpretation,
   code line, played-out paths). Assign a scenario number only by deliberate
   README edit — README numbering wins over plans (Scenario Registry rule).
5. Update `docs/patterns/agentic-patterns-reference.md` `catrace example`
   column when the example is on the merge path.
6. Run the example; paste real output into the walkthrough (no paraphrased
   numbers).

Normal review for docs/polish; GAN only if user asks or the change is
merge-critical library / experiment runner code (see review skill).

---

## Operation: close

Run when the example is ready to merge or just merged.

### Definition of Done checklist

- [ ] `examples/<name>/main.go` runs; stdout matches AC
- [ ] `examples/<name>/WALKTHROUGH.md` present and convention-compliant
- [ ] README scenario entry complete (including played-out version)
- [ ] Story file matches what the code models
- [ ] Pattern reference `catrace example` column updated
- [ ] `docs/wiki/scenario-registry.md` status → Implemented (or Draft if still missing pieces)
- [ ] Wiki: ingest new/changed story + walkthrough; example page; pattern catalogue status; `log.md` — via `llm-wiki` skill
- [ ] Close the `example:` issue (comment with example path + PR if any)

Status becomes **Implemented** only when code + walkthrough + README scenario
are all present.

---

## Operation: lint

Check and **fix** drift you can resolve without human judgment.

### 1. Gather evidence (per known example / open `example:` issue)

| Evidence | How |
|----------|-----|
| E1 story | `docs/patterns/story-*.md` linked from issue |
| E2 code | `examples/<name>/main.go` |
| E3 walkthrough | `examples/<name>/WALKTHROUGH.md` |
| E4 README | Scenario write-up under README with code line to this example |
| E5 registry | Row in `docs/wiki/scenario-registry.md` |
| E6 pattern table | `catrace example` column in agentic-patterns-reference |
| E7 issue | Open/closed `example:` issue |

### 2. Derive status

```
if E2 and E3 and E4:
    status = Implemented
else if E2:
    status = Draft
else:
    status = Planned
```

### 3. Maintain

| Drift | Fix |
|-------|-----|
| Registry status ≠ derived | Update Scenario Registry (and pattern catalogue claims) |
| Implemented but E3 or E4 missing | Downgrade to Draft; do not claim Implemented in wiki |
| Draft with no open issue | Flag; offer to open/reopen `example: …` |
| Pattern table lists example but E2 missing | Clear or mark planned |
| Orphan `examples/*/main.go` with no story/issue | Flag |
| README scenario number claimed by a plan for a different story | Fix plan or registry; README wins |

### 4. Report

```
## Examples lint
| Name / issue | Derived | Registry | Code | Walkthrough | README | Action |
```

Coordinate wiki prose fixes with `llm-wiki`; this skill owns the
example ↔ issue ↔ registry **status** contract.

---

## Operation: file

When promoting a new pattern candidate (not yet an issue):

1. Ensure `docs/patterns/story-<slug>.md` exists (or draft it with state
   meanings + interpretation + empty Issue footer).
2. Open GitHub issue titled `example: add <pattern> example` with body:

```markdown
## Pattern

**<Name>** — … in [agentic-patterns-reference](docs/patterns/agentic-patterns-reference.md)

Story: [docs/patterns/story-….md](docs/patterns/story-….md)

## Description

…

## State spaces

- **World:** …
- **Experience:** …
- **Actions:** …

## catrace features exercised

- …

## Acceptance criteria

- [ ] Runnable `examples/<name>/main.go`
- [ ] …
- [ ] Walkthrough at `examples/<name>/WALKTHROUGH.md` following conventions
- [ ] README scenario entry updated with story link
```

3. Link issue in the story footer.
4. Leave pattern-reference example column as `—` until implement/close.
5. Run **start** before coding — filing alone does not skip real-ish critique.

---

## Rules

- Prefer `jq` / `gh --jq` for GitHub JSON (project rule).
- **Never skip start** on a first-time example unless the user explicitly
  confirms the story is already real-ish.
- Do not invent README scenario numbers; edit README deliberately and sync
  Scenario Registry.
- Do not mark Complete/Implemented for experiments here — wrong skill.
- Examples teach; if the work becomes “variant A vs B falsifies claim X”,
  file an experiment instead (or in addition) with the experiments skill.
