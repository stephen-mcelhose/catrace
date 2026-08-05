---
name: experiments
description: >
  Maintain catrace variant-comparison experiments: file new hypotheses, sync
  experiments/README.md and wiki Experiment Registry, ensure every pending
  experiment has a GitHub issue, close/update issues when verdicts land, and
  lint drift between hypothesis files / registry / issues. Use when the user
  says "file an experiment", "lint experiments", "sync experiment issues",
  "maintain experiments", "experiment registry", or after adding/editing
  experiments/*/hypothesis.md.
license: MIT
metadata:
  version: "1.0.0"
---

# Experiments maintenance

catrace experiments are formal hypothesis records under `experiments/<slug>/`.
They are distinct from examples: examples teach the API; experiments ask a
falsifiable design question.

## When to use

| Trigger | Operation |
|---------|-----------|
| New architectural claim to test | **file** |
| After editing hypotheses / registry | **lint** / **sync** |
| After filling Results + Verdict | **complete** |
| "Do we have issues for experiments?" | **lint** (issue coverage) |

## Canonical locations

| Path | Role |
|------|------|
| `experiments/<slug>/hypothesis.md` | Source of truth for claim, variants, predictions, results |
| `experiments/hypothesis-template.md` | Template for new files |
| `experiments/README.md` | Human registry table |
| `docs/wiki/experiment-registry.md` | Wiki registry (keep in sync with README) |
| `docs/wiki/variant-comparison-methodology.md` | Pending experiments table + methodology |
| GitHub issues titled `experiment: …` | Tracking / discussion; body links hypothesis path |

## Issue policy

1. **Every Pending experiment MUST have an open GitHub issue** titled  
   `experiment: <short claim phrase>` with a `## See` link to  
   `experiments/<slug>/hypothesis.md`.
2. **Complete** experiments: issue may be closed with a comment pointing at
   Verdict; or left open only if follow-up work remains — note that in the
   hypothesis Interpretation.
3. Registry tables SHOULD include an Issue column (`#N` or `—` for complete
   with no issue).
4. Do not invent README scenario numbers in experiment issues (see
   `docs/wiki/scenario-registry.md`).

## Operation: file

1. Choose a kebab-case `slug` from the claim.
2. Create `experiments/<slug>/hypothesis.md` from
   `experiments/hypothesis-template.md`.
3. Fill Claim, Context, Variants, Variable entries, Predictions, Verdict rule
   **before** any run. Leave Results empty.
4. Add a row to `experiments/README.md` (Status: Pending, Verdict: —).
5. Add the same row to `docs/wiki/experiment-registry.md`.
6. Add a row to the Pending table in
   `docs/wiki/variant-comparison-methodology.md` if the claim fits that page.
7. **Open a GitHub issue** via `gh`:

   ```bash
   gh issue create --title "experiment: <short claim>" --body "$(cat <<'EOF'
   ## Claim

   <1–3 sentences from hypothesis>

   ## Motivation

   <from hypothesis Context>

   ## Experiment design

   <variants + key predictions>

   ## Blockers

   <e.g. needs examples/network_of_healers — or "none">

   ## See

   `experiments/<slug>/hypothesis.md`
   EOF
   )"
   ```

8. Put `#N` into README + wiki registry Issue columns.
9. Append `docs/wiki/log.md`:  
   `## [YYYY-MM-DD] ingest | experiments/<slug>`

## Operation: lint

Inventory and report (fix what you can):

1. **Disk vs README** — every `experiments/*/hypothesis.md` (except template)
   has a README row; every README slug has a directory.
2. **README vs wiki registry** — same slugs, status, verdict.
3. **Issue coverage** — for each Pending row, an open issue exists and links
   the hypothesis path. Create missing issues (ask before bulk-opening if >3).
4. **Stale Complete** — if Verdict is filled but Status still Pending, flip to
   Complete and update wiki.
5. **Orphan issues** — `gh issue list` titles matching `experiment:` whose
   hypothesis path is missing → flag for close or restore.
6. **Methodology pending table** — not every experiment must appear, but any
   listed slug must exist on disk.

Report format:

```
## Experiments lint
- N hypotheses on disk
- M issues missing / fixed
- K registry drifts fixed
- Unfixable: …
```

Append a short lint entry to `docs/wiki/log.md` when you change wiki pages.

## Operation: complete

After Results + Verdict are written in the hypothesis:

1. Set Status/Verdict in `experiments/README.md` and
   `docs/wiki/experiment-registry.md`.
2. Move detail into wiki "Completed results" if architecturally notable.
3. Comment on the GitHub issue with verdict summary; close if no follow-up.
4. Log the wiki update.

## Operation: sync

Shorthand: run **lint**, then ensure issue bodies still point at the current
hypothesis path and blockers (e.g. blocked on `network_of_healers`).

## Rules

- Never run the experiment and fill Results in the same step as filing
  predictions without an explicit user request to execute.
- Identical topology across variants; only listed kernel entries differ
  (see template).
- Prefer linking plans (`plans/*.md`) in Context when the example is not
  implemented yet.
- Coordinate with `llm-wiki` when wiki pages change; this skill owns the
  experiment↔issue contract.
