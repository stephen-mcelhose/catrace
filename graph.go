package catrace

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// NewRandomWalkKernel constructs a Kernel from a weighted adjacency matrix.
//
// adj[i][j] is the weight of the edge from node i to node j. Zero means no
// edge. For undirected graphs the matrix should be symmetric.
//
// The transition probability is:
//
//	P(i, j) = adj[i][j] / sum_k adj[i][k]
//
// In plain terms: from node i, move to each neighbor with probability
// proportional to the edge weight. This is the standard random walk on a
// graph — the graph and the Markov kernel are two views of the same object.
//
// For undirected graphs the stationary distribution has a closed form:
//
//	π(i) = degree(i) / sum_k degree(k)
//
// meaning high-degree nodes are visited most often. This can be verified by
// calling Stationary on the returned kernel.
//
// Returns an error if any row of adj is all-zero (an isolated node with no
// outgoing edges, which would make the row of P undefined).
func NewRandomWalkKernel(adj *mat.Dense, names []string) (*Kernel, error) {
	if adj == nil {
		return nil, fmt.Errorf("nil adjacency matrix")
	}
	r, c := adj.Dims()
	if r != c {
		return nil, fmt.Errorf("adjacency matrix must be square, got %dx%d", r, c)
	}

	p := mat.NewDense(r, r, nil)
	for i := 0; i < r; i++ {
		sum := 0.0
		for j := 0; j < r; j++ {
			v := adj.At(i, j)
			if v < 0 {
				return nil, fmt.Errorf("adjacency matrix has negative entry at (%d,%d): %g", i, j, v)
			}
			sum += v
		}
		if sum == 0 {
			name := fmt.Sprintf("%d", i)
			if len(names) > i {
				name = names[i]
			}
			return nil, fmt.Errorf("node %q has no outgoing edges (row %d is all zero)", name, i)
		}
		for j := 0; j < r; j++ {
			p.Set(i, j, adj.At(i, j)/sum)
		}
	}

	return NewKernel(p, names)
}
