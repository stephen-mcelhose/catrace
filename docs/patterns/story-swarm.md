# Story: Swarm / peer-to-peer agents

Story:

A fleet of autonomous search drones is deployed over a disaster zone to locate survivors. There is no central coordinator — each drone decides for itself where to go next based on two signals: what it directly observes in its local area, and a shared coverage heat-map that all drones read and write. A drone that finds a dense uncovered region moves there; one that finds its local area well-covered disperses to a sparser zone. Emergent global coverage arises entirely from local decisions and the shared map, with no agent ever knowing the full picture. The system can get stuck if all drones simultaneously interpret the shared map as pointing them to the same region, producing a coverage collapse elsewhere.

State meanings:
- zone world states: `uncovered`, `partially_covered`, `well_covered`, `complete` — the true aggregate coverage level of the disaster zone
- drone experience states: `local_empty`, `local_partial`, `local_found` — what each individual drone observes in its immediate search area
- drone actions: `explore_new_area`, `reinforce_local`, `signal_found` — where the drone focuses its next search cycle

Interpretation:
- each drone's perception is coupled to the true zone world state through local sampling: a well-covered zone makes `local_found` more likely for any given drone, but individual drones still see noise
- decisions are independent: $D_\text{joint} = \bigotimes_i D_i$; no drone communicates directly with another within a cycle — they coordinate only through the shared world state
- the action kernel captures how each drone's movement choice changes the aggregate zone coverage level
- the world kernel $W = PDA$ gives the coverage dynamics; its stationary distribution shows the long-run fraction of time the zone spends in each coverage level
- entropy rate measures coordination efficiency: a well-tuned swarm has low entropy (drones reliably disperse to uncovered areas); a poorly tuned one has high entropy (drones cluster unpredictably)

---

Issue: *not yet filed*

[← Back to pattern reference](agentic-patterns-reference.md)
