# Walkthrough: Blackboard / shared-workspace collaboration

Prior examples either couple peers through each other's world health
(`validator_repair`) or advance a fixed pipeline one active stage at a time
(`prompt_chaining`). This example is different: **many specialists, one shared
board** — agents never message each other; board status is what they perceive
and what their joint posts rewrite.

Run:

```
go run examples/blackboard/main.go
```

---

## The story

Three identical specialists (role labels only) read a shared diagnostic board
and write findings, endorsements, flags, or requests. Structural novelty:
opportunistic contribution via shared state — a `tentative_diagnosis` raises
$P(\texttt{evidence_strong})$ — plus an action rule that counts posts /
endorsements / flags to move the board.

---

## State spaces

**Blackboard (shared world):**

| Layer | States | Meaning |
|-------|--------|---------|
| World | `undiagnosed`, `tentative_diagnosis`, `confirmed_diagnosis`, `contradicted` | Collective agreement; `confirmed_diagnosis` absorbing; `contradicted` may reopen |

**Each specialist** (identical P, D):

| Layer | States | Meaning |
|-------|--------|---------|
| Experience | `evidence_strong`, `evidence_weak` | Whether the board looks actionable (precondition fires) |
| Action | `post_finding`, `endorse_prior`, `flag_contradiction`, `request_more_data` | Contribute, support, challenge, or abstain-with-ask |

**Joint spaces:** $X^3$ (8 experiences) and $G^3$ (64 joint actions). Indexing is
base-$n$ digits: specialist 0 is the high place.

---

## Coupling

- **P_joint:** every specialist observes the **shared board only**. Empty board
  → mostly `evidence_weak`; `tentative_diagnosis` → mostly `evidence_strong`.
  No peer experience is visible.
- **D_joint = D⊗D⊗D:** independent decisions; out-of-band A2A forbidden.
- **A(w, g, ·):** next board depends on current status **and** joint action
  counts (posts / endorsements / flags). Because of the $w$ dependence, $J$ is
  assembled row-wise rather than as one rectangle $A: G\to W$.

---

## Math worked by hand

From `undiagnosed`, a single specialist posts with probability

$$
P(\texttt{post}) = 0.30\cdot 0.30 + 0.70\cdot 0.08 = 0.146.
$$

Three independent specialists: probability **at least one** posts is
$1 - (1-0.146)^3 \approx 0.377$. Under the action rule that case sends mass
$0.80$ to `tentative_diagnosis`, so a rough lower bound on
$J[\texttt{undiagnosed},\texttt{tentative}]$ is $0.377\times 0.80 \approx 0.30$.
The program prints $0.351558$ — higher because the “no post” branch still
leaks $0.08$ to tentative, and multi-post paths still use the same $0.80$ rule.

From `tentative_diagnosis`, two or more endorsements (no flags) send $0.80$ to
`confirmed_diagnosis`. That is why $J[\texttt{tentative},\texttt{confirmed}]$
lands near one third of the row mass once perception biases toward
`evidence_strong` and $D$ favors `endorse_prior`.

---

## Reading the output

```
=== Blackboard: per-specialist kernels (identical; roles are labels) ===

P (board → evidence_strong / evidence_weak):
⎡ 0.3   0.7⎤
  ⎢0.75  0.25⎥
  ⎢0.85  0.15⎥
  ⎣ 0.2   0.8⎦

D (evidence → post / endorse / flag / request):
⎡ 0.3  0.45   0.1  0.15⎤
  ⎣0.08  0.07   0.4  0.45⎦
```

Empty board mostly fails the precondition; a tentative hypothesis flips that.

```
=== Board world kernel J (assembled Σ P_joint · D_joint · A(w,g,·)) ===
⎡0.617127  0.351558  0.031315         0⎤
  ⎢0.046592  0.329206  0.323796  0.300406⎥
  ⎢       0         0         1         0⎥
  ⎣     0.2      0.05         0      0.75⎦
```

Row 0: empty board usually stays empty, but ~35% opens a tentative hypothesis.
Row 1: from tentative, mass splits among stay / confirm (~32%) / contradict (~30%).
Row 2: confirmed is absorbing. Row 3: contradicted reopens 20% of the time.

```
Stationary distribution π:
  undiagnosed            0.000000
  tentative_diagnosis    0.000000
  confirmed_diagnosis    1.000000
  contradicted           0.000000

Entropy rate H(J): 0.000000 bits/step
Recurrent classes: [[2]]
Transient states:  [0 1 3]
MFPT undiagnosed → confirmed_diagnosis: 10.4147 steps
```

Long-run every case eventually confirms once contradicted cases reopen — so
$\pi$ sits on `confirmed_diagnosis` and entropy is zero. The interesting number
is latency: about **10.4 steps** from empty board to confirmed.

```
=== Trace onto {undiagnosed, confirmed_diagnosis, contradicted} ===
⎡ 0.6415455105054607  0.20101400471381675  0.15744048478072253⎤
  ⎢                  0                    1                    0⎥
  ⎣ 0.2034728992805541  0.02413527849086307   0.7723918222285828⎦

IsTraceOf verification: true
```

Collapsing `tentative_diagnosis` shows direct-ish flow from empty → confirmed
(~20%) versus empty → contradicted (~16%) before reopening.

```
=== Sample A(·|board) for key joint actions ===
  undiagnosed + post|ask|ask   → [0.15 0.80 0.05 0.00]
  tentative + end|end|ask      → [0.02 0.15 0.80 0.03]
  tentative + end|end|end      → [0.02 0.15 0.80 0.03]
  tentative + flag|end|ask     → [0.05 0.25 0.10 0.60]
```

One post opens the board; two endorsements confirm; any flag risks contradiction.

---

## What you can change

1. **Raise undiagnosed → evidence_strong** (e.g. `0.30` → `0.50` in `specP`).
   Watch MFPT fall — specialists contribute before a hypothesis exists.
2. **Weaken tentative trigger** (`0.75` → `0.40` on strong). Watch
   `J[tentative, confirmed]` drop and MFPT rise — board no longer enables
   endorsement.
3. **Harden flags** (increase `0.60` contradicted mass when `nFlag >= 1`).
   Watch Trace mass empty→contradicted grow — opportunistic challenge wins more often.
4. **Make confirmed leak** (small reopen to `undiagnosed`). Stationary leaves
   `{1,0,0,0}` shape; entropy becomes nonzero — closed cases are no longer final.
5. **Heterogeneous `specP` / `specD` per specialist** (leave v1 identical).
   Compare MFPT — tests whether modality accuracy, not just the board, drives latency.
