# Story: Single LLM task agent

Story:

An LLM support agent is assigned to handle a task independently. The real task may be routine or genuinely complex, but the agent only sees the prompt and surrounding context, so it can misread the situation. Based on its internal interpretation, it may answer directly, ask a clarifying question, or escalate to a human.

State meanings:
- world states = what is actually true about the task
- experience states = how the LLM interprets the task
- action states = what the LLM does next

Interpretation:
- perception captures imperfect prompt interpretation
- decision captures the LLM policy
- action captures how the chosen response changes the real task situation
- the derived kernel $Q = DAP$ tells you how the agent's interpretation evolves from one interaction cycle to the next

#### Played-out version: Story 1

The composite kernel $Q = DAP$ is easiest to understand when you walk through concrete paths. Each entry of $Q$ compresses many possible world-experience-action-world-experience micro-stories into one effective next-experience probability.

**Version A: statistically typical path**

1. The real task is `task_routine` — the ticket is genuinely straightforward.
2. Perception: the agent reads the prompt and, with probability 0.85, experiences `looks_routine`.
3. Decision: given `looks_routine`, the agent chooses `answer` with probability 0.8.
4. Action effect: answering directly resolves the issue, so the world stays `task_routine` with probability 0.9.
5. Re-perception: the world is still routine, so the next experience is again `looks_routine` with probability 0.85.

This path contributes the bulk of the probability mass to the transition `looks_routine -> looks_routine` in $Q$.

In plain English: the task really was simple, it looked simple, the agent answered directly, and the situation remained simple when seen again.

**Version B: alternate path — misread complex task**

1. The real task is `task_complex` — the ticket is actually ambiguous or difficult.
2. Perception: but the prompt is incomplete, so with probability 0.25 the agent still experiences `looks_routine`.
3. Decision: given `looks_routine`, the agent chooses `answer` with probability 0.8.
4. Action effect: answering directly rarely resolves a complex task, so the world remains `task_complex` with probability 0.1.
5. Re-perception: the unresolved situation now looks problematic, so the next experience becomes `looks_risky` with probability 0.75.

This path contributes to the transition `looks_routine -> looks_risky` in $Q$.

In plain English: the task was harder than it looked, the agent answered too quickly, the problem persisted, and on the next pass the agent finally saw the difficulty.

**Version C: recovery path — productive caution**

1. The real task is `task_complex`.
2. Perception: this time the agent reads it correctly and experiences `looks_risky` with probability 0.75.
3. Decision: given `looks_risky`, the agent chooses `clarify` with probability 0.3 (or `escalate` with probability 0.6).
4. Action effect: asking a clarifying question moves the world toward `task_routine` with probability 0.4.
5. Re-perception: the now-routine task is perceived as `looks_routine` with probability 0.85.

This path contributes to the transition `looks_risky -> looks_routine` in $Q$.

In plain English: the agent correctly flagged a hard task, asked for more context, the situation improved, and the next reading was routine.

**Why multiple paths help**

Together these paths show that one entry of the composite kernel is not one literal event. It is an aggregation of many possible micro-stories. When you read a probability in $Q$, think: this number compresses many possible world-experience-action-world-experience paths into one effective next-experience probability.

Concise shorthand for reading $Q$ entries:

- `looks_routine -> looks_routine` — the agent correctly handled a manageable task
- `looks_routine -> looks_risky` — the task was harder than it first appeared, or a direct answer did not stabilize it
- `looks_risky -> looks_routine` — clarification or escalation successfully reduced uncertainty
- `looks_risky -> looks_risky` — the problem remained hard even after intervention

---

Code: `examples/simple_agent/main.go` — walkthrough at `examples/simple_agent/WALKTHROUGH.md`

[← Back to pattern reference](agentic-patterns-reference.md)
