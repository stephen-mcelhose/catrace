# Story: Implement-Verify

*Sub-class of: [Evaluator-Optimizer (6)](agentic-patterns-reference.md)*

Story:

A coding agent attempts to fix a failing test suite. After each edit, it runs the test suite and type checker — objective tools that return pass/fail signals, not qualitative feedback. The agent cannot see whether its current approach is fundamentally correct; it only sees whether the checks pass. It loops — implement, observe results, implement again — until all checks pass or a step budget is exhausted. The tension is subtle: a fully passing check suite is not the same as a correct fix. An agent that optimises too hard for green tests may produce a solution that satisfies the verifier but regresses in unobserved ways, or patches symptoms rather than causes.

Unlike the Evaluator-Optimizer pattern, the evaluator here is not an agent with its own perception loop — it is the environment itself. The verifier's output becomes the implementer's world state signal, closing the loop purely through tool feedback rather than agent-to-agent communication.

State meanings:
- world states: `all_failing`, `some_failing`, `all_passing`, `regressed` — the true state of the test suite and type checker
- implementer experience states: `output_bad`, `output_partial`, `output_clean` — what the agent reads from the tool output
- implementer actions: `targeted_fix`, `broad_refactor`, `revert_and_retry` — implementation strategies available each cycle

Interpretation:
- perception captures how reliably the agent reads tool output — noisy compiler messages, misleading test names, and partial failures all reduce the signal quality
- decision captures the agent's repair strategy given its perceived check state
- action captures how the chosen strategy shifts the true world state: targeted fixes are cheap but may miss root causes; broad refactors are expensive but can resolve structural issues; reverts reset to a known state at the cost of discarding progress
- the world kernel $W = PDA$ gives the full check-state transition dynamics including the probability of accidental regression; its stationary distribution shows how much time the agent spends in each check state over a long run
- the `regressed` state is the key catrace insight: it is reachable from `all_passing` — a check-suite that passes today may fail tomorrow if the fix was fragile

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
