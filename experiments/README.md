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
6. Add a row to the table below (include Issue `#N`)
7. Open a GitHub issue (`experiment: …`) and sync the wiki registry — or use the
   project skill `.agents/skills/experiments/SKILL.md`

## Experiments

| Slug | Claim | Status | Verdict | Issue |
|------|-------|--------|---------|-------|
| [nodes-throttle-vs-evolver](nodes-throttle-vs-evolver/hypothesis.md) | Node throttle is primary recovery mechanism; evolver contributes but is not essential | Complete | Supported (4/4 metrics) | — |
| [wiki-knowledge-graph](wiki-knowledge-graph/hypothesis.md) | Trace chain corrects importance distortion caused by missing wiki pages | Complete | Not supported (1/4 criteria) | [#21](https://github.com/stephen-mcelhose/catrace/issues/21) |
| [spectral-gap-mixing-time](spectral-gap-mixing-time/hypothesis.md) | Spectral gap rank-orders kernels correctly by empirical mixing speed | Pending | — | [#25](https://github.com/stephen-mcelhose/catrace/issues/25) |
| [stationary-sensitivity](stationary-sensitivity/hypothesis.md) | Perturbations to high-π rows cause disproportionately large shifts in stationary distribution | Pending | — | [#26](https://github.com/stephen-mcelhose/catrace/issues/26) |
| [n-agent-scalability](n-agent-scalability/hypothesis.md) | Dense joint kernel approach becomes intractable at N=4+ agents; trace collapses are the solution | Pending | — | [#27](https://github.com/stephen-mcelhose/catrace/issues/27) |
| [kg-grounding-agent-behavior](kg-grounding-agent-behavior/hypothesis.md) | Knowledge graph quality systematically shifts knowledge agent behavior — richer graph → higher π(understood), faster recovery, lower entropy | Pending | — | [#32](https://github.com/stephen-mcelhose/catrace/issues/32) |
| [heal-on-critical-path](heal-on-critical-path/hypothesis.md) | On a fixed \(1\to2\) load graph, strong local heal belongs on the downstream sink more than on the upstream feeder | Pending | — (needs `network_of_healers`) | [#33](https://github.com/stephen-mcelhose/catrace/issues/33) |
| [collapse-masks-heterogeneity](collapse-masks-heterogeneity/hypothesis.md) | Strong-sink / weak-feeder makes collapsed pool_ok look healthier than joint upstream-overload mass warrants | Pending | — (needs `network_of_healers`) | [#34](https://github.com/stephen-mcelhose/catrace/issues/34) |

## Methodology

See [docs/wiki/variant-comparison-methodology.md](../docs/wiki/variant-comparison-methodology.md) for the full methodology, hypothesis template guidance, and generalization to other patterns. Maintenance: `.agents/skills/experiments/SKILL.md`.
