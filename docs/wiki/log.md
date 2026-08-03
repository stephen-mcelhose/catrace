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
