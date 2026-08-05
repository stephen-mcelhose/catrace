# Agent instructions

## Examples

When starting, implementing, closing, or linting pattern examples (issues titled
`example: …`), follow `.agents/skills/examples/SKILL.md`. That workflow includes
a real-ish story critique before coding, then walkthrough/README conventions
and Scenario Registry / wiki close-out.

## Walkthroughs

When writing or updating a walkthrough for any example, follow the standards
defined in `docs/walkthrough-conventions.md`.

## README scenarios

When writing or updating a scenario entry in `README.md`, follow the standards
defined in `docs/readme_scenario_conventions.md`.

## Shell: prefer jq for JSON

When parsing JSON from `gh`, APIs, or files in the shell, prefer `jq` / `gh --jq`
over Python or node one-liners so results are easy to read at a glance.

## Reviews: normal default; GAN when asked or merge-critical

Default reviews are correctness + clarity + conventions (fine for docs, plans,
wiki, example polish).

Use **GAN / adversarial** review (try to break the change) when the user asks
for it, or when reviewing merge-critical code: package-root `*.go`,
`experiments/*/main.go`, or a PR about to merge runnable behavior. Follow
`.agents/skills/review/SKILL.md`.
