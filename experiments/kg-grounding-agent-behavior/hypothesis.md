# Hypothesis: Knowledge graph quality systematically shifts knowledge agent behavior

## Claim

The structure of a knowledge graph is the primary determinant of a knowledge agent's behavioral
quality. An agent grounded by a well-structured graph — high coverage, foundational pages
well-connected — spends significantly more time in correct-understanding states, recovers faster
from confusion, and produces more predictable behavior than the same agent grounded by a sparse
or gap-ridden graph. The policy (D) and action effects (A) are held constant: only the
quality of the grounding changes.

## Context

- **Pattern:** Single LLM task agent (1), extended to include a knowledge graph as the mechanism
  shaping perception
- **Related examples:** `examples/simple_agent`, `experiments/wiki-knowledge-graph`
- **Motivation:** Knowledge graph quality is often treated as a qualitative concern — "more
  complete is better." catrace lets us make this precise: KG quality maps directly to entries
  in the Perception kernel P, and from there we get exact predictions about stationary
  distribution, recovery time, and entropy rate. This experiment establishes whether the effect
  is large enough to be architecturally decisive, and whether all three behavioral metrics agree.

  A secondary motivation: the `wiki-knowledge-graph` experiment measures KG quality using the
  trace chain (corrected PageRank scores). If KG quality provably shifts agent behavior metrics
  in this model, the trace-corrected scores from that experiment become a principled input to P
  — creating a pipeline from *measure your graph* to *predict your agent's behavior*.

## State spaces

| Space | States | Meaning |
|-------|--------|---------|
| W (world) | `routine`, `complex`, `edge-case` | The actual nature of the task the agent receives |
| X (experience) | `understood`, `partial`, `lost` | What the agent believes about the task after KG lookup |
| G (action) | `answer`, `clarify`, `escalate` | What the agent does |

The world cycles back through P after each action: answering/clarifying/escalating changes the
task context the agent faces next (routine tasks addressed, complex tasks partially resolved,
escalated tasks removed from the queue).

## Variant definitions

| Variant | Label | Description |
|---------|-------|-------------|
| A | Rich KG | Knowledge graph has high coverage of all task types. Foundational pages are well-linked and have high trace-corrected importance. Agent perceives world state accurately most of the time. |
| B | Sparse KG | Knowledge graph has significant gaps — missing pages, few cross-references, low trace-corrected importance on key topics. Agent frequently misreads task complexity. |

*All other kernel entries (D, A) are identical. Only P entries differ.*

## Variable kernel entries

The Perception kernel P maps world states → experience states. KG quality is encoded here:
a well-linked page for a topic makes the agent more likely to correctly assess a task of that type.

| Kernel | Entry (world → experience) | Variant A (rich KG) | Variant B (sparse KG) | What it encodes |
|--------|---------------------------|---------------------|----------------------|-----------------|
| P | `routine` → `understood` | 0.80 | 0.50 | Routine tasks are well-covered in A; patchy in B |
| P | `routine` → `partial` | 0.15 | 0.35 | A rarely misreads routine; B often does |
| P | `routine` → `lost` | 0.05 | 0.15 | A almost never loses a routine task; B sometimes does |
| P | `complex` → `understood` | 0.50 | 0.20 | A can stitch together multi-page context; B often can't |
| P | `complex` → `partial` | 0.40 | 0.50 | B more often gets only part of the picture |
| P | `complex` → `lost` | 0.10 | 0.30 | B frequently has no grounding for complex tasks |
| P | `edge-case` → `understood` | 0.25 | 0.10 | A has some coverage of edge cases; B almost none |
| P | `edge-case` → `partial` | 0.45 | 0.30 | |
| P | `edge-case` → `lost` | 0.30 | 0.60 | B is often lost on edge cases |

*All P rows sum to 1.0 in both variants.*

## Predictions

| Metric | Expression | Predicted direction | Why |
|--------|-----------|--------------------|-----|
| Stationary mass | π(`understood`) | A > B | Rich KG means the agent correctly reads most tasks most of the time |
| Stationary mass | π(`lost`) | A < B | Sparse KG leaves many tasks with no grounding |
| Recovery speed | MFPT(`lost` → `understood`) | A < B | From confusion, A recovers quickly on the next task; B is more likely to stay lost |
| Predictability | H(Q) bits/step | A < B | A well-grounded agent has a more legible action distribution; B's random misperceptions add entropy |
| Observable outcome | π_trace(`answered`) after tracing onto {`answered`, `escalated`} | A > B | External observer sees more direct answers from A, more escalations from B |

## Verdict rule

Claim is supported if ≥ 4 of 5 metrics agree. The direction on both stationary mass metrics
must be correct (π(understood) and π(lost)) for the claim to be structurally coherent — a split
there would indicate a modeling error, not a trade-off.

---

## Results

> Fill in after running `go run experiments/kg-grounding-agent-behavior/main.go`

| Metric | Variant A | Variant B | Predicted direction | Correct? |
|--------|-----------|-----------|--------------------|---------|
| π(`understood`) | | | A > B | |
| π(`lost`) | | | A < B | |
| MFPT(`lost` → `understood`) | | | A < B | |
| H(Q) | | | A < B | |
| π_trace(`answered`) | | | A > B | |

## Verdict

**Claim:** [ supported / not supported / trade-off ]
**Metrics in agreement:** [ /5 ]

## Interpretation

> Fill in after running.

## Connection to wiki-knowledge-graph

The `wiki-knowledge-graph` experiment measures trace-corrected PageRank scores for the catrace
wiki. If this experiment is supported, those scores become a principled source for P entries:
pages with high corrected importance → low perception noise for tasks in their topic area;
pages missing entirely → high `lost` probability for tasks that would have relied on them.

This closes a loop: use catrace to measure your knowledge graph, use those measurements to
parameterize a catrace agent model, use the agent model to predict behavioral impact.
The pipeline is entirely within catrace's existing API.
