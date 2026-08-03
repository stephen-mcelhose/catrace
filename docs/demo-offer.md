# catrace — Demo Offer

> **Want to explore this yourself?**
> Open an AI assistant, point it at this repo, and paste:
> *"Explain this project using the wiki at docs/wiki/index.md"*
>
> The wiki is a set of interconnected pages the AI can navigate — you'll get a walkthrough
> of whichever area you find most interesting. Good starting questions: *"What is a trace chain?"*
> or *"Walk me through the validator-repair scenario."*

---

## What this is

AI agents are often described in ways that are hard to hold accountable:
*"it usually gets it right," "it recovers when something goes wrong," "adding a validator helps."*
These are intuitions, not measurements.

catrace is a Go library for turning those intuitions into numbers. The approach:
model an AI agent as a small state machine, then compute — how often is it in a good state?
how fast does it recover? how predictable is its behavior? The same technique scales up
to knowledge-grounded agents and multi-agent systems.

Three areas worth exploring, each building on the last.

---

## 1. Measuring a single agent

An agent has three tables:

- **Perception** — given the world is in some state, what does the agent *think* it's seeing?
  (encodes misreads, noise, partial context)
- **Decision** — given what it thinks it's seeing, what does it do?
  (the policy)
- **Action effect** — given what it does, how does the world actually change?
  (causal consequence)

Compose the three tables into one loop — perceive → decide → act → perceive again —
and you have a state machine you can compute on directly:

- **Stationary distribution** — where does the agent spend most of its time?
  *"52% in the healthy state, 8% degraded"*
- **Mean first passage time** — how long does recovery take from the worst state?
  *"2.03 steps on average"*
- **Entropy rate** — how predictable is its behavior?
  *"1.87 bits per step"*

These are exact answers, not simulation estimates. The `*.html` graphs in the repo make the
state flows visible — edge weights are the transition probabilities, node size reflects
how often the chain is there.

**The payoff:** design claims become falsifiable before you build anything. *"This agent
should spend most of its time in the valid state"* is now a checkable prediction, not an opinion.

---

## 2. Knowledge graphs as the agent's perception

An agent that reads from a knowledge graph — a wiki, a RAG index, a doc store — is still
a perceive → decide → act system. The difference is that the **perception kernel** is now
shaped by the knowledge graph: what the graph covers, and how well it's connected,
determines what the agent can correctly understand about a task.

A pending experiment in this project (`experiments/kg-grounding-agent-behavior/`) encodes
this as a catrace model. The world states are task types: `routine`, `complex`, `edge-case`.
The experience states are what the agent believes after its knowledge lookup: `understood`,
`partial`, `lost`. The claim under test:

> *A well-structured knowledge graph produces an agent that spends significantly more time
> in the "understood" state, recovers from confusion faster, and behaves more predictably —
> holding policy and action effects constant.*

Only the Perception kernel P changes between variants. Everything else is identical.
The metrics vote on whether the claim holds.

A second experiment (`experiments/wiki-knowledge-graph/`) applies the same machinery
to the catrace wiki itself — modelling the page-link structure as a Markov chain,
computing which pages are structurally load-bearing (PageRank), and using the
**trace chain** to correct for pages that don't exist yet.
A well-structured knowledge graph shows up in the numbers before anyone reads it.

**The thread:** the same tool that measures agent recovery speed can tell you which pages
in your knowledge base your agents are silently depending on.

---

## 3. What happens when agents interact

A validator watching a worker. A self-healing node with an outer evolutionary loop.
Two agents interacting create joint behavior that neither agent fully controls,
and that is very hard to reason about by inspection.

catrace models this by building a *joint* state machine over the combined state space —
every combination of agent states becomes one joint state. Coupling is made explicit:
the validator's perception is wired to the worker's world state (so it can detect degradation);
the validator's repair action is wired to the worker's world state (so repair actually
changes something). Independent assumptions are also explicit: decisions are separate unless
you say otherwise.

The joint system is analyzed with the same tools as a single agent. You can ask:

- Which joint states are reachable at all?
- When the worker degrades, how long before the validator brings it back?
- If you trace down to just "both fine / both failed," what does the effective two-state
  picture look like?

A completed experiment (`experiments/nodes-throttle-vs-evolver/`) found that the
node's local throttle — not the outer evolutionary loop — is the primary recovery mechanism.
The outer loop is a safety net, not a fast path. That changes how much engineering attention
each feedback loop deserves.

**The thread:** adding agents adds interactions. The interactions are where the surprising
behavior lives. Encoding them explicitly is the difference between a system you understand
and one you hope works.

---

## Two disciplines worth dwelling on

**The LLM-wiki** (`docs/wiki/`) — this project's wiki was built incrementally by an AI agent:
ingesting source files, writing interconnected pages, checking for contradictions, running a
lint pass, updating an index. The `log.md` shows every operation in sequence.
It's a live example of AI-assisted knowledge management running inside the project it documents —
and the wiki-knowledge-graph experiment uses catrace to measure the structure the AI agent produced.

**Manual hypothesis refinement** (`experiments/`) — every experiment here requires writing the
claim and the predicted metric directions *before* running anything. State the claim in plain
language. Identify which kernel entries encode it. Predict which way each metric should go.
Then run and check. The completed experiment in `experiments/nodes-throttle-vs-evolver/`
shows the full form: four metrics predicted, four checked, verdict returned.

The discipline matters because it changes what a surprising result means. If you write
the prediction first and the numbers go the other way, you've learned something.
If you write the interpretation after seeing the numbers, you've explained something —
which is not the same thing.
