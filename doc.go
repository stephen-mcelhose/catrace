// Package catrace implements finite-state Markov models for autonomous-agent networks.
//
// # Core concept
//
// The central operation is the trace of a Markov kernel: given a Markov chain on
// a large state space N, the trace onto a subset A ⊆ N is the induced chain you
// would observe if you watched only the states in A and integrated out all hidden
// excursions through the complement. Concretely, if P is partitioned as
//
//	P = [ a  b ]
//	    [ d  c ]
//
// with rows and columns of a indexed by A, then the trace kernel is
//
//	P_A = a + b (I - c)⁻¹ d
//
// The term b (I-c)⁻¹ d captures the contribution of all paths that leave A,
// wander through the hidden states, and eventually return. This is not the same
// as deleting the hidden rows and columns.
//
// A key property: the stationary distribution of the trace kernel equals the
// stationary distribution of the parent kernel restricted and renormalized to A.
//
// # Agent model
//
// An Agent is defined by three rectangular row-stochastic maps that form a
// closed loop:
//
//	D : X → G   (decision:  experience states → action states)
//	A : G → W   (effect:    action states     → world states)
//	P : W → X   (perception: world states     → experience states)
//
// Composing these in different orders gives square kernels on each space:
//
//	QualiaKernel    Q = D·A·P   on X  (experience dynamics)
//	StrategyKernel  S = A·P·D   on G  (action dynamics)
//	WorldKernel     W = P·D·A   on W  (world dynamics)
//
// The three kernels are cyclic permutations of the same product and share
// eigenvalues, providing three perspectives on the same closed-loop system.
//
// # Kernel operations
//
// All analysis methods are defined on [Kernel], a square row-stochastic matrix
// with named states:
//
//   - [Kernel.Trace] — compute the trace onto a subset of states
//   - [Kernel.Stationary] — stationary distribution via power iteration (ergodic chains only)
//   - [Kernel.EntropyRate] — entropy rate in a specified log base
//   - [Kernel.MeanFirstPassage] — expected steps from state i to state j
//   - [Kernel.CommuteTime] — mean first-passage time i→j plus j→i
//   - [Kernel.Classes] — communicating class decomposition (recurrent/transient, period)
//   - [Kernel.Sample] — forward sample one step given a starting state
//   - [Kernel.LeftAction] — evolve a distribution one step: π' = π·P
//
// # Sampling and estimation
//
//   - [EstimateKernelFromSequence] — empirical transition counts with pseudocount smoothing
//   - [SampleTraceFromSequence] — filter a trajectory to observed states
//   - [WindowedTraceEstimates] — sliding-window kernel estimates from a trajectory
package catrace
