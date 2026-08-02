# Story: Self-adjusting / self-healing network nodes

Story:

A network node monitors its own error rate and throttles itself when errors climb. An outer evolutionary loop watches pool throughput and mutates the node's configuration when performance drops. The two loops compete to explain recovery: in one regime the node's own throttle is the primary healer; in the other the evolver's config search is what keeps the system alive. Running the same joint kernel architecture under both parameter regimes lets the stationary distribution, mean first passage time, and entropy rate vote on which mechanism does the work.

State meanings:
- node world states: `healthy`, `degraded`, `overloaded` — actual error rate of the node
- node experience states: `ema_low`, `ema_mid`, `ema_high` — EMA band the node observes
- node actions: `push`, `throttle`, `idle`
- evolver world states: `good_strategy`, `bad_strategy` — whether the current MaxWorkers/Kp config is effective
- evolver experience states: `high_score`, `low_score` — pool-level throughput×success² the evolver observes
- evolver actions: `promote`, `mutate`
- joint world states: `H·G`, `H·B`, `D·G`, `D·B`, `O·G`, `O·B` — (node health · evolver strategy) pairs

Interpretation:
- the evolver's perception is coupled to node health in $P_\text{joint}$: a sick node depresses the pool score the evolver observes, making a good config look bad even when it isn't
- decisions are independent: $D_\text{joint} = D_\text{node} \otimes D_\text{evolver}$; each agent reads its own signal without communicating within a cycle
- when the evolver mutates, it boosts the node's recovery probability in $A_\text{joint}$; `promote` leaves node recovery unchanged
- the variant comparison pattern uses the model as a measuring instrument: run both parameter regimes, read the difference in MFPT and entropy rate, and let the numbers say which mechanism carries the load

#### Played-out version: Story 4

The joint kernel $J = P_\text{joint} \cdot D_\text{joint} \cdot A_\text{joint}$ compresses a full $W \to X \to G \to W$ cycle — for both agents simultaneously — into one effective joint-state transition. The paths below use Variant B parameters, where the evolver's contribution is most visible.

**Version A: statistically typical path**

1. The system is in `H·G` — node healthy, evolver running a good config.
2. Perception: the healthy node observes `ema_low` (probability 0.80); the evolver, watching a healthy node's pool score, reads `high_score` (probability 0.85). Joint experience: `ema_low·high_score`, probability $0.80 \times 0.85 = 0.680$.
3. Decision: reading `ema_low`, the node pushes (probability 0.75); reading `high_score`, the evolver elects to promote its current config (probability 0.80). Joint action: `push·promote`, probability $0.75 \times 0.80 = 0.600$.
4. Action effect: `push` keeps the node `healthy` (probability 0.55); `promote` preserves the `good_strategy` (probability 0.85). Joint next world: `H·G`, probability $0.55 \times 0.85 = 0.468$.

This path contributes the bulk of probability mass to the `H·G → H·G` transition in J.

In plain English: the node was healthy, the EMA stayed quiet, the evolver saw good scores and kept the config, and the system remained in its best state.

**Version B: degraded node masks a good config**

1. The system is in `D·G` — the node has drifted into a degraded state, but the evolver's config is actually good.
2. Perception: the degraded node reads `ema_mid` (probability 0.55); the pool score the evolver observes is suppressed by the degradation — even a good config reads as `low_score` with probability 0.45. Joint experience: `ema_mid·low_score`, probability $0.55 \times 0.45 = 0.248$.
3. Decision: reading `ema_mid`, the node throttles (probability 0.55); reading `low_score`, the evolver mutates away from the (actually good) config (probability 0.75). Joint action: `thr·mutate`, probability $0.55 \times 0.75 = 0.413$.
4. Action effect: the mutation boost raises the node's `healthy` probability from 0.45 to 0.56 after renormalization; mutating gives 0.60 probability of landing on a `good_strategy`. Joint next world: `H·G`, probability $0.56 \times 0.60 = 0.336$.

This path contributes to the `D·G → H·G` transition in J.

In plain English: the node's degradation made a good config appear bad; the evolver unnecessarily searched for a new one, but the mutation boost it applied helped the node recover along the way.

**Version C: recovery from the worst state**

1. The system is in `O·B` — the node is overloaded and the evolver has a bad config.
2. Perception: the overloaded node reads `ema_high` (probability 0.75); a bad strategy with an overloaded node always shows `low_score` (probability 1.00). Joint experience: `ema_high·low_score`, probability $0.75 \times 1.00 = 0.750$.
3. Decision: reading `ema_high`, the node idles (probability 0.45); reading `low_score`, the evolver mutates (probability 0.75). Joint action: `idle·mutate`, probability $0.45 \times 0.75 = 0.338$.
4. Action effect: the mutation boost raises the node's `healthy` probability from 0.60 to 0.667 after renormalization; mutating gives 0.60 probability of landing on `good_strategy`. Joint next world: `H·G`, probability $0.667 \times 0.60 = 0.400$.

This path contributes to the `O·B → H·G` transition in J.

In plain English: both mechanisms fired together — the node backed off under EMA pressure and the evolver searched for a better config; the mutation boost made the combined recovery faster than either would have managed alone.

**Why the coupling matters**

Without the perception coupling, a degraded node would not suppress the pool score — the evolver would keep promoting a good config rather than accidentally mutating away from it, and the `D·G → D·B` leakage would disappear. Without the action coupling, mutation would have no effect on node recovery. Together they create a system where the outer loop can accidentally harm state estimation (perception coupling) and simultaneously provide the primary recovery pathway (action coupling). The variant comparison reads this tension directly: in Variant B the evolver is essential; in Variant A the throttle is strong enough that the outer loop is a minor boost on top of self-healing that's already working.

Concise shorthand for reading J entries:
- `H·G → H·G` — node healthy, evolver keeping the best config; system holds
- `D·G → D·B` — degraded node suppresses the score; evolver mutates away from the good config
- `D·G → H·G` — throttle recovers the node; mutation boost (if large enough) accelerates it
- `O·B → H·G` — idle plus a lucky mutation search finds the recovery path in one step
- `O·B → O·B` — throttle fails and mutation finds no better config; system stays stuck

---

Code: `examples/self_healing_nodes/main.go` — walkthrough at `examples/self_healing_nodes/WALKTHROUGH.md`

[← Back to pattern reference](agentic-patterns-reference.md)
