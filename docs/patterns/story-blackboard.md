# Story: Blackboard / shared-workspace collaboration

Story:

A hospital diagnostic system brings together three specialist agents — a radiologist reading scan images, a pathologist reading biopsy results, and a clinical notes reader parsing the patient history — to form a joint diagnosis. No specialist can see the others' raw inputs, but all of them read and write to a shared blackboard of accumulated findings. Each specialist contributes when it has something to add, potentially revising the working diagnosis in the light of prior postings. The system converges when two or more specialists endorse the same conclusion, or it flags a contradiction if their findings are irreconcilable.

State meanings:
- blackboard world states: `undiagnosed`, `tentative_diagnosis`, `confirmed_diagnosis`, `contradicted` — the true state of collective agreement on the case
- specialist experience states: `evidence_strong`, `evidence_weak` — whether the specialist's own input material is clear enough to draw a confident finding
- specialist actions: `post_finding`, `endorse_prior`, `flag_contradiction`, `request_more_data` — how the specialist advances or challenges the shared workspace

Interpretation:
- each specialist's perception is coupled to the blackboard world state: a `tentative_diagnosis` on the board shifts each specialist's experience toward `evidence_strong` if their own data aligns, or toward `evidence_weak` if it conflicts
- decisions are independent within a cycle: $D_\text{joint} = D_\text{rad} \otimes D_\text{path} \otimes D_\text{notes}$; each specialist reads the shared board but does not coordinate its action with the others before posting
- the action kernel captures how the combination of specialist posts advances the blackboard state toward confirmed or contradicted
- the world kernel $W = PDA$ gives the full diagnostic state-transition dynamics; mean first passage time from `undiagnosed` to `confirmed_diagnosis` measures diagnostic latency
- tracing onto `{undiagnosed, confirmed_diagnosis, contradicted}` collapses intermediate agreement states and gives the coarse outcome distribution

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
