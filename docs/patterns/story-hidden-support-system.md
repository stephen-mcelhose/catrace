# Story: LLM agent with hidden support system

Story:

A focal LLM agent is visible to us, but the rest of the support system is hidden in the background: retrieval services, monitoring tools, human reviewers, and other agents. We only observe whether the focal agent appears valid or invalid from the outside. The hidden system may help or hinder it before we see the focal agent again.

State meanings:
- `A_valid`, `A_invalid` = visible health of the focal agent
- `B_valid`, `B_invalid` = hidden health of the surrounding system

Interpretation:
- the parent kernel models the full visible-plus-hidden system
- the trace onto `{A_valid, A_invalid}` gives the effective observed dynamics of the focal agent alone
- this is not simple deletion of hidden states; it folds hidden excursions into the observed transition probabilities

---

Code: `examples/trace_analysis/main.go` — walkthrough at `examples/trace_analysis/WALKTHROUGH.md`

[← Back to pattern reference](agentic-patterns-reference.md)
