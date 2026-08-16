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

// NewTeleportingKernelFromAdj constructs a teleporting Markov kernel directly
// from a raw weighted adjacency matrix, combining row-normalisation and
// teleportation in a single step:
//
//	T[i][j] = α·v[j] + (1−α)·(adj[i][j] / rowsum[i])   if rowsum[i] > 0
//	T[i][j] = v[j]                                        if rowsum[i] = 0
//
// Sink nodes (rows that sum to zero — pages with no outgoing links) collapse
// entirely to the restart distribution v. No artificial uniform edges are
// inserted; the teleportation term carries them instead. This is semantically
// cleaner than pre-filling sink rows with 1/n before calling
// NewRandomWalkKernel.
//
// The returned kernel is the same object as TeleportingKernel would produce on
// a pre-normalised input, but avoids the two-step NewRandomWalkKernel →
// TeleportingKernel pipeline that errors on sink nodes.
//
// alpha must be in [0, 1]. restart must be a valid probability vector (non-
// negative, sums to 1 within 1e-9). adj must be square with non-negative
// entries.
func NewTeleportingKernelFromAdj(adj *mat.Dense, restart []float64, alpha float64, names []string) (*Kernel, error) {
	if adj == nil {
		return nil, fmt.Errorf("nil adjacency matrix")
	}
	r, c := adj.Dims()
	if r != c {
		return nil, fmt.Errorf("adjacency matrix must be square, got %dx%d", r, c)
	}
	n := r
	if alpha < 0 || alpha > 1 {
		return nil, fmt.Errorf("alpha %g is out of range [0, 1]", alpha)
	}
	const tol = 1e-9
	v, err := cloneProbDist(restart, n, tol, "restart", false)
	if err != nil {
		return nil, err
	}
	data, err := teleportMatrix(adj, v, alpha)
	if err != nil {
		return nil, err
	}
	return NewKernel(mat.NewDense(n, n, data), names)
}

// teleportMatrix assembles the flat n×n entries of the teleporting kernel T.
func teleportMatrix(adj *mat.Dense, v []float64, alpha float64) ([]float64, error) {
	n := len(v)
	data := make([]float64, n*n)
	for i := 0; i < n; i++ {
		row, err := teleportRow(adj, i, v, alpha)
		if err != nil {
			return nil, err
		}
		copy(data[i*n:(i+1)*n], row)
	}
	return data, nil
}

func teleportRow(adj *mat.Dense, i int, v []float64, alpha float64) ([]float64, error) {
	n := len(v)
	rowSum := 0.0
	for j := 0; j < n; j++ {
		val := adj.At(i, j)
		if val < 0 {
			return nil, fmt.Errorf("adjacency matrix has negative entry at (%d,%d): %g", i, j, val)
		}
		rowSum += val
	}
	row := make([]float64, n)
	for j := 0; j < n; j++ {
		if rowSum > 0 {
			row[j] = alpha*v[j] + (1-alpha)*adj.At(i, j)/rowSum
		} else {
			// Sink: no outgoing links — collapse entirely to restart.
			row[j] = v[j]
		}
	}
	return row, nil
}
