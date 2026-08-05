# Wiki Schema

This wiki is maintained by an LLM using the llm-wiki skill, following the Karpathy pattern.

## Domain

catrace — a Go library for modelling agentic systems as Markov chains using a
perception-decision-action (P, D, A) kernel triplet. The wiki covers the mathematical
foundations, the catrace API, the modeling methodology for single and multi-agent
systems, the catalogue of agentic patterns, and the implemented examples.

## Raw sources

Raw sources are the existing project documents, read in place. The LLM never edits them.

| Path                                         | Content                                      |
|----------------------------------------------|----------------------------------------------|
| `README.md`                                  | Project overview and scenario list           |
| `GLOSSARY.md`                                | Term definitions                             |
| `docs/math_summary.md`                       | Mathematical summary of the library          |
| `docs/source_writeup.md`                     | Detailed source writeup                      |
| `docs/markovmixing.txt`                      | Academic paper on Markov mixing times        |
| `docs/patterns/agentic-patterns-reference.md`| Full pattern reference table                 |
| `docs/patterns/story-*.md`                   | Individual pattern stories (17 files)        |
| `examples/*/WALKTHROUGH.md`                  | Example walkthroughs (5 files)               |
| `plans/*.md`                                 | Implementation plans (lint cross-check only) |
| `docs/wiki/raw/`                             | Future external sources (papers, URLs)       |

## Wiki page domains

| Domain                   | Topic                                                                 |
|--------------------------|-----------------------------------------------------------------------|
| Mathematical foundations | Markov chains, kernel algebra, spectral theory, mixing, MFPT          |
| catrace API              | Kernel type, Trace, Stationary, MeanFirstPassage, EntropyRate, Classes |
| Modeling methodology     | P/D/A triplet, Q/S/W cyclic kernels, joint kernels, agent coupling    |
| Agentic patterns         | Structural patterns, dev workflow sub-patterns                        |
| Examples                 | Per-example synthesis pages linking math ↔ API ↔ pattern             |

## Conventions

- **Page slugs**: `kebab-case.md` derived from the concept title
- **Frontmatter**: every page must have `title`, `tags`, `sources`, `updated`
- **Cross-references**: `[[Page Slug]]` wikilinks (never relative paths)
- **Sources section**: every page ends with `## Sources` listing its raw inputs
- **Dates**: `date -u +%Y-%m-%d` at execution time

## Ingest discipline

Raw sources are living documents — story files are added as examples are implemented,
walkthroughs are updated, and the pattern reference evolves. Keeping the wiki current
requires a simple rule:

**When a raw source changes, re-ingest it before closing the PR.**

Specifically:
- New story file created → ingest it immediately, propagate to pattern and example pages
- Existing story file updated → re-ingest it, propagate to affected wiki pages
- New walkthrough added → ingest it, link to its pattern page and API pages it demonstrates
- `GLOSSARY.md` or `docs/math_summary.md` updated → re-ingest, propagate to math and API pages
- After any batch of raw source changes → run lint to catch stale claims
- `README.md` scenario list changes → update [[Scenario Registry]] and re-lint
- New or edited `plans/*.md` → lint against [[Scenario Registry]] (plans must not invent README scenario numbers)

The lint operation is the safety net. Run it whenever in doubt.

**Exhaustive lint must also cross-check** (beyond orphans/wikilinks):

1. README `### N.` scenarios ↔ [[Scenario Registry]] table
2. “Implemented” status requires `main.go` **and** `WALKTHROUGH.md`
3. `plans/*.md` scenario-number claims vs README (README wins)
4. [[Agentic Patterns Catalogue]] / [[Structural Patterns]] status vs registry

## Operations

Run these via the `llm-wiki` skill:

- `ingest <source>` — read a raw source, write a summary page, propagate to related pages
- `query <question>` — synthesize an answer from wiki pages, optionally write back
- `lint` — audit for orphans, contradictions, stale claims, missing links
