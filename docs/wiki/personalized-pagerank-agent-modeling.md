---
title: Personalized PageRank and Agent Modeling
tags: [ppr, pagerank, agent, intent, goal-directed, stationary, pda, network, modeling, teleporting-kernel]
sources: [docs/adr/0001-stationary-from-custom-start.md]
updated: 2026-08-13
---

# Personalized PageRank and Agent Modeling

This page explains why `PersonalizedPageRank` is a natural model for goal-directed agent behavior in the catrace framework, how it relates to `StationaryFrom`, and what open research questions follow from the connection.

> **Implementation status (2026-08-13):** `TeleportingKernel`, `PersonalizedPageRank`, and `NewTeleportingKernelFromAdj` are all implemented in the catrace API. See [[catrace API]] for signatures. The document-graph use case described in §"Knowledge-graph grounding" below is now concrete — see `tools/firehose-graph/` in `eis-intake-firehose` for a working PPR visualization of a 64-page document corpus using α=0.15 with the 03-recommend gate pages as the restart distribution.

## The core distinction

Both methods iterate on a kernel `P` starting from a distribution `v`:

| | `StationaryFrom(v)` | `PersonalizedPageRank(v, α)` |
|---|---|---|
| Step | `x ← x·P` | `x ← α·v + (1−α)·x·P` |
| Role of `v` | Initial `x₀` only | Biases **every** step |
| Fixed point | Stationary dist of `P` | Stationary dist of teleporting chain |
| Convergence | Requires ergodic chain | Guaranteed for `α ∈ (0, 1]` |
| Result depends on `v`? | No (ergodic chains) | Yes, always |

`StationaryFrom` answers: *where does this chain settle if it starts here?* The start carries meaning but the answer is determined by `P` alone.

`PersonalizedPageRank` answers: *where does this chain settle when an intent keeps pulling it back?* The result is jointly determined by the dynamics `P` and the goal distribution `v`, mediated by `α`.

## Interpretation in the PDA model

In catrace, a single agent is described by the [[PDA Triplet Model]] — three kernels P (perception), D (decision), A (action) — whose composition yields the qualia kernel Q = D·A·P operating on experience space X.

The stationary distribution `π(Q)` describes the **long-run distribution over qualia states under pure dynamics**: where the agent's experience loop settles without any imposed goal.

Running PPR on Q with restart distribution `v` and weight `α` yields a **goal-directed qualia distribution**: where the agent's experience loop settles when it is persistently drawn toward the qualia states encoded in `v`.

| Quantity | Agent interpretation |
|---|---|
| `π(Q)` | Long-run experience under pure dynamics (no intent) |
| `v` | Goal distribution — the qualia states the agent is drawn toward |
| `α` | Intentionality weight — how strongly goals override trained dynamics |
| `1 − α` | Fraction of behavior governed by the trained kernel |
| PPR fixed point | Long-run experience under both dynamics **and** intent |
| `π(Q) − ppr(Q, v, α)` | How much goals actually shift the agent's experience |

A high `α` means the agent frequently "snaps back" to its goal states regardless of where the dynamics would take it. A low `α` means the agent has goals but mostly follows its trained kernel — goals are a gentle bias, not a dominant force. At `α = 1` the agent is entirely goal-driven (the fixed point is just `v`); at `α = 0` PPR reduces to `StationaryFrom`.

## Interpretation in a network of agents

The picture becomes richer in multi-agent settings (see [[Joint Kernels and Coupling]]).

**Orchestrator pull.** In an orchestrator-workers pattern, the orchestrator's preferred state distribution becomes the restart vector `v` for the joint kernel. `α` is the strength of top-down control. High `α` → orchestrator dominates; low `α` → emergent worker dynamics dominate. The PPR vector on the joint kernel gives the long-run state distribution under that control regime.

**Supervision and reward.** If certain joint states correspond to task-complete or correct-output conditions, concentrating `v` on those states and running PPR gives the long-run distribution under a persistent reward signal — without full reinforcement learning machinery. The gap `‖π − ppr‖₁` measures how much the reward signal shifts long-run behavior.

**Coupling strength measurement.** In a two-agent system where one agent's goals influence the other, the PPR parameter `α` can model coupling strength. Sweeping `α` from 0 to 1 and observing how the PPR vector shifts gives a sensitivity profile of the network to the external intent — a structural analogue of influence analysis.

**Knowledge-graph grounding.** For agents grounded in a knowledge graph (e.g. wikigraph issue #22), PPR on the random-walk kernel of the graph with restart concentrated on goal nodes gives structurally-weighted node importance. This is exactly the Personalized PageRank application from Jeh & Widom (2003): the seed set is the query context, `α` is the damping factor, and the PPR vector ranks every node by its structural proximity to the query.

## Relationship to other catrace metrics

**MFPT.** Mean first passage time `m(i, j)` measures expected steps from state `i` to state `j`. This is a per-pair local measure. PPR gives a *global* distribution — how much total long-run mass lands near a seed set, integrated over all paths. The two metrics are complementary: MFPT is a probe of access cost to a specific state; PPR is a probe of structural centrality relative to a set. It is an open question whether `ppr(v, α)` and MFPT to the support of `v` are related by a closed-form expression.

**Entropy rate.** The entropy rate `H(P)` is defined over the stationary distribution of `P`. The teleporting chain `α·v·𝟙ᵀ + (1−α)·P` has its own entropy rate, which can be computed by calling `EntropyRate` on the modified kernel. This "goal-directed entropy rate" measures the uncertainty per step of an agent operating under both its dynamics and its intent — and is lower-bounded by `α·H(point mass) = 0` (fully goal-driven, deterministic snap-back) and upper-bounded by `H(P)` (no goals, pure dynamics).

**Trace chain.** A [[Trace Chain]] projects a kernel onto an observed subset. Running PPR on a trace chain rather than the parent kernel models a goal-directed agent whose goals are expressed entirely in the observed experience space, with the hidden states integrated out.

## Why `StationaryFrom` alone is insufficient for goal-directed modeling

`StationaryFrom(v)` starts from `v` but the agent's goals play no role after step 1. On any ergodic kernel, every starting distribution converges to the same `π`. There is no way to encode that an agent *keeps* preferring certain states — the model has no mechanism for persistent intent.

PPR is the minimal extension that adds this: a single scalar `α` captures the trade-off between trained dynamics and goal-directedness, and the fixed point is analytically meaningful (stationary distribution of the teleporting chain). This makes it the right primitive for modeling intentional agents.

## Open questions

The following are not yet answered in catrace theory or experiments:

1. **PPR on Q vs W.** Running PPR on the qualia kernel Q vs the world kernel W gives different fixed points. What is the relationship? Does PPR commute with the cyclic permutation between Q, S, W?

2. **Choosing α.** In the knowledge-graph application `α = 0.15` (the Google PageRank damping factor) is now confirmed conventional — the `firehose-graph` tool uses it on the eis-intake-firehose document corpus and produces coherent results. For agent modeling the right value remains open: is there a principled way to infer `α` from observed agent behavior?

3. **PPR and MFPT.** Is there a closed-form relationship between the PPR mass assigned to a state and the MFPT to that state from a typical starting point? Such a relationship would unify two of catrace's primary metrics.

4. **Multi-agent PPR.** In a network with `n` agents, each with its own goal distribution, can PPR be defined on the joint kernel with a joint restart distribution that encodes each agent's goals? Does this decompose into per-agent PPR computations under independence assumptions?

5. **Goal-directed entropy rate.** How does `H(teleporting chain)` vary with `α` and `v`? Is it always monotone in `α`? This would characterize how intentionality reduces behavioral uncertainty.

6. **Sensitivity.** The `experiments/stationary-sensitivity` experiment studies how small perturbations to `P` shift `π`. PPR adds a second perturbation axis (`v`, `α`). How do these interact?

## Sources and further reading

- ADR-0001: `docs/adr/0001-stationary-from-custom-start.md`
- Jeh, G. & Widom, J. (2003). Scaling Personalized Web Search. *WWW 2003.*
- Haveliwala, T. (2002). Topic-Sensitive PageRank. *WWW 2002.*
- Brin, S. & Page, L. (1998). The Anatomy of a Large-Scale Hypertextual Web Search Engine. *WWW 1998.*
- Hoffman, D., Prakash, C. & Chattopadhyay, S. (2024). Traces of Consciousness — source of the PDA triplet formalism.
- [[Markov Chain Foundations]]
- [[PDA Triplet Model]]
- [[Joint Kernels and Coupling]]
