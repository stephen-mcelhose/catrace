# Mathematical Summary (Sanitized)

This note extracts only the finite-state stochastic machinery relevant to modeling networks of autonomous agents from the uploaded preprint, while intentionally excluding consciousness and philosophical claims.

> [!INFO]
> Source for this mathematical extraction and terminology:
> - https://www.preprints.org/manuscript/202410.1305#sec-preprints-h2-13
> 
> The write-up below keeps only the finite-state stochastic content used for implementation: coupled kernels between world, experience, and action spaces, together with the trace-chain reduction on observed subsets.

## 1. Finite-state stochastic maps

A finite-state system is represented by a row-stochastic matrix #$#P#$# where

- rows index current states,
- columns index next states,
- each entry is nonnegative,
- each row sums to 1.

Thus #$#P(i,j) = \Pr(X_{t+1}=j \mid X_t=i)#$#.

## 2. Agent as three coupled kernels

A single agent is modeled using three finite stochastic maps among finite sets:

- world states #$#W#$#,
- experiences #$#X#$#,
- actions #$#G#$#.

The three row-stochastic kernels are:

- #$#P: W \to X#$#  (perception)
- #$#D: X \to G#$#  (decision)
- #$#A: G \to W#$#  (action-to-world effect)

> [!INFO]
> Plain English for the paper-style notation:
> - #$#P : W \Rightarrow X#$# means perception maps a world state to a probability distribution over experiences.
> - #$#D : X \Rightarrow G#$# means decision maps a current experience to a probability distribution over actions.
> - #$#A : G \Rightarrow W#$# means action maps a chosen action to a probability distribution over resulting world states.
> - #$#Q : X \Rightarrow X#$# is not raw perception; it is the full closed-loop experience-to-experience kernel after perception, decision, and action are composed.
> - notation such as #$#X_\sigma#$# usually means the experience space associated with a particular agent, subsystem, or observer index #$#\sigma#$#.

With row-stochastic convention, the composed square kernels are:

- qualia kernel on #$#X#$#: #$#Q = D A P#$#
- strategy kernel on #$#G#$#: #$#S = A P D#$#
- world kernel on #$#W#$#: #$#K_W = P D A#$#

These are ordinary matrix products with compatible dimensions.

## 3. Trace chain on a subset

For a parent Markov kernel #$#P#$# on state set #$#N#$#, and a subset #$#A \subseteq N#$#, reorder states so that

#$#P = \begin{bmatrix} a & b \\ d & c \end{bmatrix}#$#

where block #$#a#$# is indexed by #$#A#$# and block #$#c#$# is indexed by #$#A^c#$#.

The trace chain on #$#A#$# is

#$#P_A = a + b (I-c)^{-1} d#$#

when #$#(I-c)^{-1}#$# exists.

Interpretation:
- #$#a#$# = direct transitions within the observed subset,
- #$#b (I-c)^{-1} d#$# = excursions outside the subset that eventually return.

## 4. Stationary distributions

A stationary distribution #$#\pi#$# satisfies

#$#\pi P = \pi, \qquad \sum_i \pi_i = 1, \qquad \pi_i \ge 0.#$#

For a trace chain #$#P_A#$#, the restricted stationary law is obtained by normalizing the parent stationary mass on #$#A#$#:

#$#\pi_A(i) = \frac{\pi(i)}{\sum_{j\in A} \pi(j)} \quad \text{for } i\in A.#$#

This is the practical identity checked in the examples.

## 5. Communicating classes and recurrence

A communicating class is a strongly connected component of the directed graph induced by positive-probability transitions.

A class is recurrent/closed if no state in the class has a positive-probability transition to a state outside the class.

These structures are useful for:
- identifying invariant subsystems,
- checking reducibility,
- understanding whether a stationary distribution is unique.

## 6. Entropy rate

For stationary distribution #$#\pi#$#, the entropy rate is

#$#H(P) = -\sum_i \pi_i \sum_j P_{ij} \log_b P_{ij}#$#

where #$#b=2#$# gives bits per step.

## 7. First-passage and commute times

Mean first-passage times are computed by solving the standard linear system on the non-target states.

If #$#j#$# is the target and #$#h_i = \mathbb{E}_i[T_j]#$# for #$#i \ne j#$#, then

#$#h_i = 1 + \sum_{k \ne j} P_{ik} h_k#$#

with #$#h_j = 0#$#.

Commute time between states #$#i#$# and #$#j#$# is

#$#C_{ij} = m_{ij} + m_{ji}#$#

where #$#m_{ij}#$# is the mean first-passage time from #$#i#$# to #$#j#$#.

## 8. Practical math corrections used in code

The paper contains notation/history suggestive of an earlier amplitude-based formulation. For implementation, the code uses only classical stochastic kernels.

Practical corrections:

1. All kernels are enforced to be row-stochastic.
2. No amplitude-like coefficients such as #$#1/\sqrt{2}#$# are used unless they already produce valid classical probabilities after squaring/normalization; in practice they are ignored as obsolete artifacts.
3. All composed kernels are checked to remain stochastic.
4. Trace kernels are computed using the classical block-matrix formula above.
5. Small numerical drift is corrected by row normalization when appropriate.

## 9. Modeling viewpoint for small agent networks

For 2-, 3-, and 4-agent examples, the simplest useful state spaces are mostly binary:

- world valid / invalid,
- local observation valid / invalid,
- action sets such as repair / hold / escalate.

Larger multi-agent systems can then be represented either by:

- direct joint-state kernels over the product state space, or
- modular agent kernels composed with a world update kernel.
