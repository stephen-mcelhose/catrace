# Agentic Patterns Reference

Collected from primary sources for use as catrace example candidates.

## Sources

| ID  | Reference                                                                                                   |
|-----|-------------------------------------------------------------------------------------------------------------|
| ANT | Anthropic, *Building Effective Agents*, 2024 — https://www.anthropic.com/engineering/building-effective-agents |
| IBM | IBM Think, *AI Agent Use Cases* — https://www.ibm.com/think/topics/ai-agent-use-cases                       |
| LG  | VThink Technologies, *Common Agentic Patterns in LangGraph* — https://www.linkedin.com/pulse/common-agentic-patterns-langgraph-vthink-technologies-4k4yc |

---

## Structural / Topological Patterns

These describe *how agents are connected and coordinate*, independent of domain.

| #  | Pattern                   | Topology             | Agents / Roles                                       | Description                                                                                                                               | Sources    | Story                                                  | catrace example                       |
|----|---------------------------|----------------------|------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|------------|--------------------------------------------------------|---------------------------------------|
| 1  | Augmented LLM             | Single node          | Agent + tools/memory/retrieval                       | One agent enhanced with external capabilities (search, memory, APIs). The foundation for all higher patterns.                             | ANT        | [story-single-llm-agent](story-single-llm-agent.md)   | `simple_agent`                        |
| 2  | Prompt Chaining           | Linear pipeline      | N sequential agents, optional gate between each      | Output of each step feeds as input to the next. Gates validate intermediate results before continuing.                                    | ANT, LG    | [story-prompt-chaining](story-prompt-chaining.md)      | `prompt_chaining`                     |
| 3  | Routing                   | Hub-and-spoke        | Router (classifier) + N specialised handlers         | A classifier agent inspects the input and directs it to the most appropriate specialist. Each specialist is optimised for one category.   | ANT, LG    | [story-routing](story-routing.md)                      | —                                     |
| 4  | Parallelisation / Fan-out | Fork-join            | N parallel workers + aggregator                      | The same (or partitioned) input is processed by multiple agents simultaneously; outputs are merged. Includes voting as a special case.    | ANT, LG    | [story-parallelisation](story-parallelisation.md)      | —                                     |
| 5  | Orchestrator-Workers      | Centralised hub      | Orchestrator + N workers + synthesiser               | A planner agent dynamically decomposes a task, dispatches subtasks to workers, handles retries/escalation, and merges results.            | ANT, LG    | [story-orchestrator-workers](story-orchestrator-workers.md) | —                                |
| 6  | Evaluator-Optimizer       | Iterative loop       | Generator + evaluator/critic                         | Generator produces a draft; evaluator scores it and provides feedback; loop repeats until a quality threshold is met or a step limit hit. | ANT, LG    | [story-validator-repair](story-validator-repair.md)    | `validator_repair` (partial)          |
| 7  | Autonomous Agent Loop     | Single cyclic agent  | Agent + tool interface + environment                 | One agent runs a perceive→decide→act loop indefinitely, using tool outputs as its next world state signal.                                | ANT        | [story-single-llm-agent](story-single-llm-agent.md)   | `simple_agent`, `self_healing_nodes`  |
| 8  | Supervisor / Hierarchical | Tree                 | Supervisor + sub-agents (possibly multi-level)       | Supervisor assigns goals to sub-agents and monitors progress; sub-agents may themselves be orchestrators.                                 | LG         | [story-supervisor](story-supervisor.md)                | —                                     |
| 9  | Swarm / Peer-to-peer      | Mesh / fully connected | N peer agents with shared state or message bus     | No central coordinator; agents observe shared state and act; emergent global behaviour arises from local rules.                           | —          | [story-swarm](story-swarm.md)                          | —                                     |
| 10 | Blackboard                | Star (shared memory) | N specialist agents + shared blackboard              | Agents read from and write to a central shared workspace; any agent may contribute whenever it can advance the solution.                  | —          | [story-blackboard](story-blackboard.md)                | —                                     |
| 11 | Debate / Adversarial      | Pair or panel        | Proposer + one or more challengers + judge           | Agents argue opposing positions; a judge (or vote) resolves the disagreement to reach a higher-quality answer.                            | —          | [story-debate](story-debate.md)                        | —                                     |
| 12 | Plan-and-Execute          | Two-phase            | Planner agent + executor agent(s)                    | Planner produces a full task plan upfront; executors carry out each step and report back; planner may revise if execution fails.          | —          | [story-plan-and-execute](story-plan-and-execute.md)    | —                                     |
| 13 | Human-in-the-Loop         | Any + human gate     | Any topology above + human reviewer                  | Agent(s) pause at defined checkpoints and request human approval or guidance before continuing.                                           | ANT        | —                                                      | —                                     |
| 14 | Self-Healing / Adaptive   | Coupled peer pair    | Operational agent + monitor/repair agent             | One agent runs a task; a second monitors its health and triggers recovery actions; coupling through perception and action.                | ANT (impl) | [story-self-healing-nodes](story-self-healing-nodes.md) | `validator_repair`, `self_healing_nodes` |

---

## Development Workflow Sub-patterns

These are specialisations of structural patterns above, oriented around AI coding and development agent loops. Each inherits from one or more structural parents but has a specific semantic meaning in the development context.

| #   | Pattern                       | Parent Pattern(s)                  | Key distinction from parent                                                                                 | Story                                                                                    | catrace example |
|-----|-------------------------------|-------------------------------------|-------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|-----------------|
| D1  | Research-Plan-Implement       | Prompt Chaining (2), Plan-and-Execute (12) | Three compounding stages: exploration gates plan quality; plan quality gates implementation success.   | [story-research-plan-implement](story-research-plan-implement.md)                        | —               |
| D2  | Implement-Verify              | Evaluator-Optimizer (6)             | Evaluator is the environment (objective tool), not an agent — no perception noise on the evaluator side.   | [story-implement-verify](story-implement-verify.md)                                      | —               |
| D3  | Implement-Critic              | Evaluator-Optimizer (6)             | Evaluator is an agent with its own imperfect perception — critic accuracy is itself a stochastic variable. | [story-implement-critic](story-implement-critic.md)                                      | —               |
| D4  | Plan-Implement-Critic-Verify  | Plan-and-Execute (12), Evaluator-Optimizer (6) | Two independent evaluators (objective verifier + subjective critic) with different failure modes and loop depths. | [story-plan-implement-critic-verify](story-plan-implement-critic-verify.md) | —               |

---

## Domain / Use-Case Applications

These are real-world scenarios IBM identifies as common deployments. Each maps onto one or more structural patterns above.

| #  | Domain               | Use Case                                          | Structural Pattern(s)          | Source |
|----|----------------------|---------------------------------------------------|--------------------------------|--------|
| A  | Agriculture          | Soil/weather monitoring + autonomous spraying     | Autonomous loop, Routing       | IBM    |
| B  | Banking & Finance    | Continuous risk auditing, fraud detection         | Autonomous loop, Evaluator     | IBM    |
| C  | Banking & Finance    | Loan underwriting automation                      | Prompt chaining, Routing       | IBM    |
| D  | Banking & Finance    | Personalised wealth management advisory           | Augmented LLM                  | IBM    |
| E  | Content Creation     | Autonomous article/report/script generation       | Prompt chaining                | IBM    |
| F  | Content Creation     | Branded asset design + video editing              | Orchestrator-workers           | IBM    |
| G  | Customer Experience  | Proactive support with sentiment + escalation     | Routing, Human-in-the-loop     | IBM    |
| H  | Customer Experience  | Agent backend assistant (retrieval + guidance)    | Augmented LLM                  | IBM    |
| I  | Disaster Response    | Damage assessment from satellite + social data    | Parallelisation, Orchestrator  | IBM    |
| J  | Disaster Response    | Predictive evacuation simulation                  | Plan-and-execute               | IBM    |
| K  | Education            | Adaptive tutoring with real-time path adjustment  | Evaluator-optimizer, Loop      | IBM    |
| L  | Education            | Interview / language simulation                   | Autonomous loop                | IBM    |

---

## catrace Example Coverage

| catrace example        | Structural pattern(s) modelled                                   |
|------------------------|------------------------------------------------------------------|
| `simple_agent`         | Augmented LLM (1), Autonomous Agent Loop (7)                     |
| `trace_analysis`       | (trace machinery demo — not a network pattern)                   |
| `validator_repair`     | Evaluator-Optimizer (6), Self-Healing pair (14)                  |
| `self_healing_nodes`   | Self-Healing / Adaptive (14), Autonomous Agent Loop (7)          |
| `prompt_chaining`      | Prompt Chaining (2)                                              |
| *planned*              | Routing (3), Parallelisation (4), Orchestrator-Workers (5), Supervisor (8), Swarm (9), Blackboard (10), Debate (11), Plan-and-Execute (12) |
