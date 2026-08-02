# Story: Prompt chaining pipeline

Story:

A document intelligence pipeline processes raw text through three specialist stages in sequence: an extractor identifies key claims, a summariser condenses them into a structured brief, and a formatter produces the final delivery-ready report. Each stage only sees its own input and cannot reach backward or forward. Between stages a programmatic gate checks whether the intermediate output meets quality criteria; if it fails, the stage reruns with slightly different framing before passing the result on. The pipeline can stall in a retry loop or exit early if the gate repeatedly rejects.

State meanings:
- pipeline world states: `raw`, `extracted`, `summarised`, `formatted`, `failed` — the true stage the document has reached
- agent experience states: `input_clear`, `input_noisy` — whether the stage agent perceives its input as processable
- agent actions: `process`, `retry`, `escalate` — attempt the transformation, reframe and try again, or abandon

Interpretation:
- perception captures how reliably a stage agent reads the quality of its incoming material
- decision captures the agent's policy given its perceived input quality
- action captures how the chosen action advances (or fails to advance) the pipeline world state
- the world kernel $W = PDA$ gives the full stage-to-stage transition probabilities including retry loops and failure exits
- tracing onto `{raw, formatted, failed}` collapses intermediate stages and shows the pipeline's end-to-end success and failure rate

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
