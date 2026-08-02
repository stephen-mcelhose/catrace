# Story: Parallelisation / fan-out with voting

Story:

A content moderation system evaluates every submitted video simultaneously with three independent specialist checkers: one for spam, one for hate speech, and one for misinformation. Each checker returns a signal without knowing what the others found. An aggregator then applies a majority-vote rule: two or more flags removes the content; two or more passes approves it. The system's overall accuracy depends not just on each checker's individual detection rate but on how their errors correlate — two checkers that fail on the same content type are much less useful than two that fail independently.

Rather than modelling the checkers as a single merged agent, each is given its own P, D, A triplet operating on the joint product state space. The decision kernel is the Kronecker product D₁ ⊗ D₂ ⊗ D₃ — checkers decide independently with no within-cycle communication. Coupling enters only in the joint action kernel, where the aggregator maps the full vector of checker votes to a moderation outcome.

State meanings:
- world states: `safe`, `spam`, `hate`, `misinfo` — the true content type
- checker experience states: `flagged`, `clear` — what each checker perceives about the content
- checker actions: `vote_flag`, `vote_pass` — each checker's individual vote
- joint actions: the 8 possible flag/pass combinations across the three checkers
- aggregation outcome: `removed` (majority flag), `approved` (majority pass)

Interpretation:
- perception captures each checker's detection accuracy for its specialist signal — naturally different across content types
- decisions are independent: $D_\text{joint} = D_1 \otimes D_2 \otimes D_3$; no checker consults the others within a cycle
- the action kernel maps joint votes to a moderation decision, encoding the majority rule
- the world kernel $W = PDA$ gives the full content-state transition matrix including false-positive removals and false-negative approvals
- entropy rate measures how unpredictable the moderation outcome is; low entropy means the system reliably agrees on clear cases

---

Issue: [#7 — example: add parallelisation / fan-out with voting example](https://github.com/stephen-mcelhose/catrace/issues/7)

[← Back to pattern reference](agentic-patterns-reference.md)
