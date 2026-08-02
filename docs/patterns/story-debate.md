# Story: Debate / adversarial agents with judge

Story:

Two legal AI agents argue a civil case from opposing positions. One agent argues for the plaintiff, the other for the defendant. A judge agent observes the accumulated argument quality and rules when it sees a clear winner, or calls for more argument when the case remains genuinely contested. Each debater perceives the current strength of its own position relative to its opponent's last move, and decides whether to press its strongest point, concede a minor issue to strengthen credibility, or object to a weak opposing claim. A debater that over-argues a weak position risks losing credibility with the judge; one that concedes too readily collapses its case.

State meanings:
- debate world states: `open`, `plaintiff_leading`, `defendant_leading`, `ruled_plaintiff`, `ruled_defendant` — the true state of accumulated argument weight
- debater experience states: `position_strong`, `position_weak` — each debater's perceived standing given the last exchange
- debater actions: `press_argument`, `concede_minor`, `object` — the debater's rhetorical moves
- judge experience states: `case_clear`, `case_contested` — whether the judge perceives a decisive advantage
- judge actions: `rule`, `call_for_more` — the judge's decision each cycle

Interpretation:
- each debater's perception is coupled to the debate world state: when the world is `plaintiff_leading`, the defendant agent is more likely to experience `position_weak` and vice versa
- debater decisions are independent: $D_\text{joint} = D_\text{plaintiff} \otimes D_\text{defendant}$; neither debater knows the other's planned move within a cycle
- the judge is a third agent whose perception and action kernel enter the joint composition alongside the debaters
- the action kernel captures how the combination of rhetorical moves shifts the world state: two simultaneous strong presses by opposing sides may cancel out; a concession from one side shifts weight toward the other
- mean first passage time from `open` to `ruled_*` measures expected argument length; entropy rate measures how predictable the debate trajectory is given the starting position

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
