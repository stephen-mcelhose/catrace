# Source Write-up

This project's mathematical framing is based on the preprint section at:

- https://www.preprints.org/manuscript/202410.1305#sec-preprints-h2-13

Briefly, the implementation takes from that source the finite-state stochastic machinery for representing an agent as three coupled kernels:

- perception $P : W \Rightarrow X$
- decision $D : X \Rightarrow G$
- action $A : G \Rightarrow W$

These compose into effective Markov kernels such as $Q = DAP$ on experience states and support trace-chain reduction of larger systems onto observed subsets. The code and documentation intentionally keep only this stochastic / Markov structure and do not attempt to implement the broader philosophical claims of the paper.
