# Wiki Log

Append-only chronological log of all wiki operations.

---

## [2026-08-02] init | catrace wiki bootstrapped

Raw source layer: existing docs in place (README.md, GLOSSARY.md, docs/, examples/).
Wiki root: docs/wiki/. No pages written yet — ready for first ingest.

---

## [2026-08-02] ingest | docs/math_summary.md + docs/source_writeup.md

Pages written: markov-chain-foundations.md, trace-chain.md, pda-triplet-model.md

---

## [2026-08-02] ingest | GLOSSARY.md

Pages written: catrace-glossary.md, joint-kernels-and-coupling.md
Propagated to: markov-chain-foundations.md, pda-triplet-model.md, trace-chain.md (wikilinks added in authored pages)

---

## [2026-08-02] ingest | README.md

Pages written: catrace-api.md
Propagated to: pda-triplet-model.md, trace-chain.md (API method references)

---

## [2026-08-02] ingest | docs/patterns/agentic-patterns-reference.md + all story-*.md files (17 files)

Pages written: agentic-patterns-catalogue.md, structural-patterns.md, dev-workflow-patterns.md
Propagated to: pda-triplet-model.md, catrace-api.md (pattern-to-API cross-references)

---

## [2026-08-02] ingest | examples/*/WALKTHROUGH.md (4 files)

Pages written: example-simple-agent.md, example-hidden-support-system.md, example-validator-repair.md, example-self-healing-nodes.md
Propagated to: structural-patterns.md, catrace-api.md, markov-chain-foundations.md, trace-chain.md, joint-kernels-and-coupling.md (example wikilinks added throughout)
index.md updated: 13 pages catalogued

---

## [2026-08-02] lint | 13 pages checked, 0 issues found, 0 fixed

All pages have frontmatter (title, tags, sources, updated). All wikilinks resolve to existing pages. Index complete. No orphans, contradictions, or stale claims detected. No missing cross-references found.

---

## [2026-08-02] new | variant-comparison-methodology.md

Promoted the variant comparison technique from a section inside example-self-healing-nodes.md to a first-class wiki page.

Pages written: variant-comparison-methodology.md
Propagated to: example-self-healing-nodes.md, structural-patterns.md, agentic-patterns-catalogue.md, dev-workflow-patterns.md
index.md updated: 14 pages catalogued

Experiment scaffold added (outside wiki root, in experiments/):
- experiments/README.md
- experiments/hypothesis-template.md
- experiments/nodes-throttle-vs-evolver/hypothesis.md (worked example, Verdict: supported 4/4 metrics)

---

## [2026-08-03] lint | 17 pages checked, 4 issues found, 4 fixed

---

## [2026-08-04] ingest | experiments/README.md

Pages written: experiment-registry.md
Propagated to: variant-comparison-methodology.md (worked example section expanded to "Completed experiments" with wiki-knowledge-graph not-supported result added), trace-chain.md (wiki-knowledge-graph experiment note added referencing [[Experiment Registry]])
index.md updated: 15 pages catalogued

---

## [2026-08-04] lint | 14 pages checked, 3 issues found, 3 fixed

Issues found and fixed:
- **variant-comparison-methodology.md**: "All three metrics (π, MFPT, H)" was stale — the `nodes-throttle-vs-evolver` hypothesis.md records four metrics (π, MFPT, H, π_trace); updated to "All four metrics"
- **variant-comparison-methodology.md**: Pending experiments table listed 3 entries; `kg-grounding-agent-behavior` (just filed) was missing; added as first row; updated count from "Three" to "Four"; updated summary sentence
- **pda-triplet-model.md**: P(w,x) entry had no reference to the knowledge-graph-grounding use case; added note that for knowledge-grounded agents, P is shaped by KG coverage/connectivity, with cross-reference to kg-grounding experiment and [[Variant Comparison Methodology]]

Unfixable (outside wiki scope):
- `experiments/README.md` records `nodes-throttle-vs-evolver` as "Supported (3/3 metrics)" — contradicts hypothesis.md (4/4) and this wiki (now 4); requires fix in experiments/README.md directly

Issues found and fixed:
- **catrace-api.md**: missing `graph.go` (NewRandomWalkKernel) and `visualise.go` (ToHTML) from package layout table; both added with method documentation sections
- **markov-chain-foundations.md**: AGENTS.md domain lists "spectral theory, mixing" but page had no content on it; added "Spectral gap and mixing time" section (spectral gap definition, mixing time bound, total variation distance, link to pending experiment #25)
- **variant-comparison-methodology.md**: only referenced nodes-throttle-vs-evolver; added "Pending experiments" table for spectral-gap-mixing-time (#25), stationary-sensitivity (#26), n-agent-scalability

No orphans. No contradictions. No index gaps.

---

## [2026-08-04] lint | exhaustive | 16 pages checked, 8 issues found, 8 fixed

Cross-checks added: README scenarios ↔ plans ↔ examples (beyond orphans/wikilinks).

Issues found and fixed:
- **plans/network-of-healers.md**: falsely claimed "This is scenario 5 in README.md"; README scenario 5 is majority-valid / `story-supervisor.md`. Corrected plan to "planned, no README number"; renamed README-entry section to draft.
- **Missing [[Scenario Registry]]**: no wiki anchor for README numbering; created `scenario-registry.md` (README 1–6 map, planned non-README examples, draft-code rules).
- **structural-patterns.md**: skipped pattern 13 (Human-in-the-Loop) while catalogue listed it; added §13 stub from agentic-patterns-reference.
- **structural-patterns.md**: Prompt Chaining / Supervisor lacked README scenario cross-links; noted scenario 6 draft `prompt_chaining` and scenario 5 ≠ network_of_healers.
- **agentic-patterns-catalogue.md**: Prompt Chaining showed bare "—"; noted draft without walkthrough; Supervisor annotated as README scenario 5.
- **index.md**: stale `updated: 2026-08-02` and missing Scenario Registry row; refreshed (16 pages).
- **AGENTS.md**: lint discipline did not cross-check README scenarios or `plans/*`; added rules + plans as lint-only raw source.
- **Stale frontmatter dates** on pages edited 2026-08-04 but still marked 2026-08-02 (markov-chain-foundations, pda-triplet-model, variant-comparison-methodology).

Already resolved (no change needed):
- `experiments/README.md` already reports Supported (4/4 metrics) — prior unfixed 3/3 contradiction cleared.

Unfixable / deferred (human judgment):
- Whether `network_of_healers` should ever become a README scenario (modelling premise under review).
- Whether to promote `examples/prompt_chaining/main.go` to a full example (needs WALKTHROUGH + README status flip).

Propagated: [[Scenario Registry]] linked from index, catalogue, catrace-api, structural-patterns, example-self-healing-nodes.
No orphans. All wikilinks resolve.

---

## [2026-08-04] ingest | heterogeneous network-of-healers experiments

Filed (pending; blocked on `examples/network_of_healers` identical-node v1 + heterogeneous kernels):

- `experiments/heal-on-critical-path/hypothesis.md` — strong heal on downstream sink vs upstream feeder
- `experiments/collapse-masks-heterogeneity/hypothesis.md` — pool_ok masking under strong-sink / weak-feeder

Propagated to: experiments/README.md, [[Experiment Registry]], [[Variant Comparison Methodology]] pending table, [[Scenario Registry]] planned-examples note, `plans/network-of-healers.md` Y-statement parked section.

---

## [2026-08-04] new | experiments maintenance skill + issue sync

- Added `.agents/skills/experiments/SKILL.md` (file / lint / complete / sync; Pending ⇒ open issue policy)
- Opened issues: #32 `kg-grounding-agent-behavior`, #33 `heal-on-critical-path`, #34 `collapse-masks-heterogeneity`
- Wired Issue column across `experiments/README.md`, [[Experiment Registry]], [[Variant Comparison Methodology]] pending table

---

## [2026-08-05] lint | exhaustive | 18 pages checked, 1 issue found, 0 fixed

Checks: frontmatter, wikilinks, orphans, index coverage, README scenarios ↔ [[Scenario Registry]], experiments disk ↔ README ↔ wiki registry ↔ Pending issue coverage.

Clean:
- All `[[wikilinks]]` resolve; no content orphans; index complete (16 catalogued pages + meta)
- All 8 experiment slugs aligned across disk / `experiments/README.md` / [[Experiment Registry]]
- All Pending experiments have open issues (#25–27, #32–34); Complete rows OK

Unfixable (needs human judgment):
- GitHub #24 `experiment: side-by-side kernel comparison visualisation` — open `experiment:` issue with **no** `experiments/<slug>/` hypothesis (orphan issue per experiments skill). File a hypothesis or close/retitle.
- #21 open with `wiki-knowledge-graph` — resolved next entry (status → Active, not premature Complete)

---

## [2026-08-05] lint | experiments status vocabulary

Added **Active** status (Pending / Active / Complete) to `.agents/skills/experiments/SKILL.md`, `experiments/README.md`, [[Experiment Registry]], [[Variant Comparison Methodology]].

Fixed premature Complete: `wiki-knowledge-graph` → **Active** (provisional not supported on PR #22; `main` hypothesis still empty). Moved narrative from Completed → Active sections; updated [[Trace Chain]] note.

Skill v1.2.0: mandatory § Status check & maintain on every lint/sync (derive status from main Results + open PRs; fix registries; report per-slug table).

---

## [2026-08-05] implement | prompt_chaining example (#5)

Pages written: example-prompt-chaining.md
Propagated to: [[Structural Patterns]] §2, [[Agentic Patterns Catalogue]], [[Scenario Registry]] (scenario 6 → Implemented), [[Wiki Index]]
Raw sources: docs/patterns/story-prompt-chaining.md, examples/prompt_chaining/{main.go,WALKTHROUGH.md}, README.md scenario 6

---

## [2026-08-05] lint | 17 content pages checked, 3 issues found, 2 fixed

Checked: orphans/wikilinks, index coverage, Scenario Registry ↔ README, Implemented = main.go+WALKTHROUGH, catalogue ↔ registry, experiments status maintain.

Clean:
- No content orphans; all index rows present; no real broken wikilinks (false positives: `[[Page Slug]]` in AGENTS, `[[0, 1]]` code output)
- README scenarios 1–6 match [[Scenario Registry]]; scenario 6 Implemented with `examples/prompt_chaining/{main.go,WALKTHROUGH.md}`
- Catalogue Prompt Chaining ✅ matches registry
- Experiments: nodes Complete; wiki-knowledge-graph Active (PR #22 open, Results empty on main); others Pending with issues

Fixed:
- `AGENTS.md` walkthrough count 4 → 5
- [[catrace API]] Implemented examples table: added `prompt_chaining`

Unfixable / needs human:
- Active `wiki-knowledge-graph` still models “14 existing + 14 planned”; [[Example: Prompt Chaining]] landing on main invalidates Graph A baseline — noted on [[Experiment Registry]]; re-baseline before final Verdict
