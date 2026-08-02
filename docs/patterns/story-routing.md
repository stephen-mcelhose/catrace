# Story: Routing agent

Story:

A customer support system receives tickets of genuinely different types — billing disputes, technical faults, and general enquiries — but the router agent can only read the ticket text, not the underlying truth. Misclassification sends the ticket to the wrong specialist, who cannot resolve it and escalates it back into the queue. The interesting tension is between the router's classification accuracy and the cost of misrouting: a confident wrong decision is worse than an uncertain escalation to a human triage agent.

State meanings:
- world states: `billing_ticket`, `technical_ticket`, `general_ticket` — the true nature of the incoming request
- experience states: `reads_billing`, `reads_technical`, `reads_general` — the router's perceived classification
- actions: `route_billing`, `route_technical`, `route_general`, `escalate_human` — where the router sends the ticket

Interpretation:
- perception captures classification accuracy — the probability that the router reads the ticket type correctly
- decision captures the router's policy given its perceived classification, including its willingness to escalate rather than guess
- action captures how the routing choice changes the world state: a correct route resolves the ticket; a wrong route re-enters it as a different apparent type
- the world kernel $W = PDA$ gives the full ticket-flow transition matrix, including misrouting loops and human escalation exits
- the stationary distribution over world states shows how much of the queue is occupied by tickets that have already been misrouted at least once

---

Issue: [#6 — example: add routing agent example](https://github.com/stephen-mcelhose/catrace/issues/6)

[← Back to pattern reference](agentic-patterns-reference.md)
