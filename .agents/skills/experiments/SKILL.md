---
name: experiments
description: >
  Maintain catrace variant-comparison experiments: file new hypotheses; check and
  maintain Pending/Active/Complete status against main hypothesis files, open PRs,
  and GitHub issues; sync experiments/README.md and wiki Experiment Registry;
  lint drift. Use when the user says "file an experiment", "lint experiments",
  "sync experiment issues", "maintain experiments", "check experiment status",
  "experiment registry", or after adding/editing experiments/*/hypothesis.md.
  Also run the status-check section whenever llm-wiki lint touches experiments.
license: MIT
metadata:
  version: "1.2.0"
---

# Experiments maintenance

catrace experiments are formal hypothesis records under `experiments/<slug>/`.
They are distinct from examples: examples teach the API; experiments ask a
falsifiable design question.

**Core duty of this skill:** not only file experiments, but **continuously check
and maintain** their Status, Issue links, and wiki/README alignment. Every
**lint** or **sync** invocation MUST run § Status check & maintain and apply
fixes it can make without human judgment.

## Status vocabulary

| Status | Means |
|--------|--------|
| **Pending** | Hypothesis filed; no run and no open implementation PR yet |
| **Active** | Work in flight — runner and/or Results exist on a branch/PR, but filled Results+Verdict are **not** fully on `main` |
| **Complete** | Results + Verdict present in `hypothesis.md` **on `main`**; registry matches the file |

Blockers (e.g. needs `network_of_healers`) stay in notes / Verdict column — do not invent a fourth status unless it proves necessary.

**Hard rule:** Never mark **Complete** if `main`'s `hypothesis.md` still has empty Results/Verdict placeholders, even if a PR already concluded the run.

## When to use

| Trigger | Operation |
|---------|-----------|
| New architectural claim to test | **file** |
| "lint / sync / maintain / check experiment status" | **lint** (includes status maintain) |
| Opening a runner PR / filling Results on a branch | **activate** |
| After Results + Verdict land on `main` | **complete** |
| llm-wiki lint when experiments changed | **lint** status section |

## Canonical locations

| Path | Role |
|------|------|
| `experiments/<slug>/hypothesis.md` | Source of truth for claim, variants, predictions, results |
| `experiments/hypothesis-template.md` | Template for new files |
| `experiments/README.md` | Human registry table |
| `docs/wiki/experiment-registry.md` | Wiki registry (keep in sync with README) |
| `docs/wiki/variant-comparison-methodology.md` | Methodology + Active / Completed / Pending tables |
| GitHub issues titled `experiment: …` | Tracking / discussion; body links hypothesis path |
| Open PRs touching `experiments/<slug>/` | Evidence for **Active** until merge |

---

## Status check & maintain (REQUIRED on every lint/sync)

For **each** slug under `experiments/*/hypothesis.md` (skip template):

### 1. Gather evidence

| Evidence | How |
|----------|-----|
| E1 Results on `main` | Read `experiments/<slug>/hypothesis.md` on current `main` (or working tree if on main). Results/Verdict are "filled" if not placeholders like `[ ]`, `[TODO]`, `[ supported / not supported`, empty table cells for all metrics |
| E2 Runner on `main` | `experiments/<slug>/main.go` (or equivalent) exists on main |
| E3 Open PR | `gh pr list --state open --search "experiments/<slug>"` (also match PR body/files) |
| E4 Open issue | Registry Issue column + `gh issue view N`; title should start with `experiment:` |
| E5 Registry status | Row in `experiments/README.md` and `docs/wiki/experiment-registry.md` |

### 2. Derive correct status

```
if E1 (Results+Verdict filled on main):
    status = Complete
else if E3 (open PR for this slug) OR (Results filled only on a non-main branch):
    status = Active
else:
    status = Pending
```

Notes:
- Provisional verdicts from a PR go in the Verdict column as  
  `provisional: … — PR #N` while status stays **Active**.
- E2 alone without E1 does not force Complete (runner without filled verdict → Active if PR, else Pending).

### 3. Maintain (apply fixes)

| Drift | Fix |
|-------|-----|
| Registry Status ≠ derived status | Update README + wiki registry Status (and Verdict notes) |
| Complete but E1 false | Flip to Active (if E3) or Pending; move wiki narrative from Completed → Active |
| Active/Pending but E1 true | Flip to Complete; move narrative to Completed; clear `provisional:` |
| Pending/Active missing open issue | Create issue (ask if >3 missing) or reopen; write `#N` into registries |
| Complete with open issue, no follow-up | Comment verdict + close issue (or ask user) |
| Open PR exists, status still Pending | **activate** (Status Active + provisional note + link PR on issue) |
| Wiki "Completed results" includes Active slug | Move to "Active experiments" |
| Wiki "Active" includes Complete slug | Move to Completed; remove provisional wording |
| Methodology pending table lists Complete slug | Remove from pending or mark done |
| `[[Trace Chain]]` / other pages claim a final verdict for Active slug | Soften to provisional + PR link |

### 4. Report

Always print a per-slug status table in the lint report:

```
| Slug | Derived | Registry | Issue | Open PR | Action taken |
```

---

## Issue policy

1. **Every Pending or Active experiment MUST have an open GitHub issue** titled  
   `experiment: <short claim phrase>` with a `## See` link to  
   `experiments/<slug>/hypothesis.md`.
2. **Complete** experiments: close the issue when the merging PR lands; leave
   open only for explicit follow-up noted in Interpretation.
3. Registry tables MUST include an Issue column (`#N` or `—` only for legacy
   Complete with no issue).
4. Do not invent README scenario numbers in experiment issues (see
   `docs/wiki/scenario-registry.md`).

## Operation: file

1. Choose a kebab-case `slug` from the claim.
2. Create `experiments/<slug>/hypothesis.md` from the template.
3. Fill Claim, Context, Variants, Predictions, Verdict rule **before** any run.
4. Add rows (Status: **Pending**) to README + wiki registry (+ methodology pending if relevant).
5. Open GitHub issue; put `#N` in registries.
6. Append wiki `log.md`.

## Operation: activate

When a runner PR opens or Results appear only on a branch:

1. Set Status **Active** in README + wiki registry.
2. Verdict column: `provisional: <verdict> — PR #N` if known.
3. Keep write-ups under wiki **Active experiments**, not Completed.
4. Comment on the issue with the PR URL.

## Operation: lint

Run in order:

1. **§ Status check & maintain** (mandatory — check and fix).
2. Disk vs README vs wiki registry slug sets.
3. Orphan `experiment:` issues (no hypothesis path) → flag.
4. Cross-links: Scenario Registry / plans that mention experiment slugs still valid.

Report:

```
## Experiments lint
### Status
| Slug | Derived | Registry | Issue | Open PR | Action |
### Other
- …
```

Append `docs/wiki/log.md` when wiki pages change.

## Operation: complete

After Results + Verdict are on **`main`**:

1. Status **Complete**; final Verdict (no `provisional:`).
2. Wiki Completed results; remove from Active.
3. Close issue (comment with verdict) unless follow-up remains.
4. Log wiki update.

## Operation: sync

Alias for **lint** with emphasis on PR/issue URLs and blocker notes still accurate.
Always includes § Status check & maintain.

## Rules

- **Check and maintain status every time** this skill runs lint/sync — do not
  only report drift.
- Never fill Results when only filing predictions unless the user asks to run.
- Identical topology across variants; only listed kernel entries differ.
- Coordinate with `llm-wiki` for wiki prose; this skill owns
  experiment ↔ issue ↔ PR ↔ **status** contract.
