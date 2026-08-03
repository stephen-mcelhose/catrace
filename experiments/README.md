# catrace experiments

This directory contains variant-comparison experiments: formal hypothesis records for architectural claims that catrace is used to test.

## What belongs here

An experiment is a *before-the-run* record of:
- A falsifiable architectural claim
- The exact kernel entries that encode each variant
- Predicted metric directions
- The results, once produced

Experiments are distinct from examples. An **example** (`examples/`) demonstrates how to use catrace on a particular pattern. An **experiment** (`experiments/`) asks and answers a specific design question using catrace as a measuring instrument.

## How to file a new experiment

1. Create a subdirectory named after the claim: `experiments/<short-slug>/`
2. Copy `hypothesis-template.md` into it and rename to `hypothesis.md`
3. Fill in Claim, Variables, and Predictions *before running*
4. Run the catrace analysis (either inline Go code in the hypothesis file, or as a standalone `main.go` in the same directory)
5. Fill in Results and Verdict
6. Add a row to the table below

## Experiments

| Slug | Claim | Status | Verdict |
|------|-------|--------|---------|
| [nodes-throttle-vs-evolver](nodes-throttle-vs-evolver/hypothesis.md) | Node throttle is primary recovery mechanism; evolver contributes but is not essential | Complete | Supported (3/3 metrics) |
| [wiki-knowledge-graph](wiki-knowledge-graph/hypothesis.md) | Trace chain corrects importance distortion caused by missing wiki pages | Pending | — |

## Methodology

See [docs/wiki/variant-comparison-methodology.md](../docs/wiki/variant-comparison-methodology.md) for the full methodology, hypothesis template guidance, and generalization to other patterns.
