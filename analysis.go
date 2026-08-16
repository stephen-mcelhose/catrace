package catrace

import "fmt"

// ClassDecomposition summarizes the communicating class structure of a Markov kernel.
//
// SCCs holds all strongly connected components identified by Kosaraju's algorithm.
// Recurrent holds the closed SCCs; once entered, these classes cannot be escaped.
// Transient holds all states that eventually leave their SCC with probability 1.
// Periods maps each recurrent SCC index (into SCCs) to its period.
type ClassDecomposition struct {
	SCCs      [][]int     // all strongly connected components, each sorted
	Recurrent [][]int     // closed (recurrent) components
	Transient []int       // transient states, sorted
	Periods   map[int]int // period of each recurrent class keyed by index into SCCs
}

// Classes decomposes k into communicating classes using Kosaraju's algorithm.
// An edge i→j is included when k.P[i,j] exceeds tol.
// Returns the full SCC decomposition, recurrent and transient states, and the
// period of each recurrent class.
func (k *Kernel) Classes(tol float64) (*ClassDecomposition, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	adj, rev := transitionDigraph(k, n, tol)
	sccs := kosarajuSCC(adj, rev, n)
	recurrent, transient, periods := classifyRecurrent(adj, sccs)
	return &ClassDecomposition{
		SCCs:      sccs,
		Recurrent: recurrent,
		Transient: sortedCopy(transient),
		Periods:   periods,
	}, nil
}

// transitionDigraph builds the forward and reverse digraphs of positive
// transitions in k (edge i→j when P[i,j] > tol).
func transitionDigraph(k *Kernel, n int, tol float64) (adj, rev [][]int) {
	adj = make([][]int, n)
	rev = make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if k.P.At(i, j) > tol {
				adj[i] = append(adj[i], j)
				rev[j] = append(rev[j], i)
			}
		}
	}
	return adj, rev
}

func kosarajuSCC(adj, rev [][]int, n int) [][]int {
	visited := make([]bool, n)
	order := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if !visited[i] {
			dfsFinishOrder(adj, visited, &order, i)
		}
	}
	for i := range visited {
		visited[i] = false
	}
	var sccs [][]int
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if !visited[v] {
			comp := []int{}
			dfsCollectComp(rev, visited, &comp, v)
			sccs = append(sccs, sortedCopy(comp))
		}
	}
	return sccs
}

func dfsFinishOrder(adj [][]int, visited []bool, order *[]int, v int) {
	visited[v] = true
	for _, w := range adj[v] {
		if !visited[w] {
			dfsFinishOrder(adj, visited, order, w)
		}
	}
	*order = append(*order, v)
}

func dfsCollectComp(rev [][]int, visited []bool, comp *[]int, v int) {
	visited[v] = true
	*comp = append(*comp, v)
	for _, w := range rev[v] {
		if !visited[w] {
			dfsCollectComp(rev, visited, comp, w)
		}
	}
}

func classifyRecurrent(adj [][]int, sccs [][]int) (recurrent [][]int, transient []int, periods map[int]int) {
	classOf := make(map[int]int)
	for idx, comp := range sccs {
		for _, v := range comp {
			classOf[v] = idx
		}
	}
	recurrent = [][]int{}
	transient = []int{}
	periods = map[int]int{}
	for idx, comp := range sccs {
		if isClosedClass(adj, comp, idx, classOf) {
			recurrent = append(recurrent, comp)
			periods[idx] = periodOfClass(adj, comp)
		} else {
			transient = append(transient, comp...)
		}
	}
	return recurrent, transient, periods
}

func isClosedClass(adj [][]int, comp []int, idx int, classOf map[int]int) bool {
	for _, v := range comp {
		for _, w := range adj[v] {
			if classOf[w] != idx {
				return false
			}
		}
	}
	return true
}

func periodOfClass(adj [][]int, comp []int) int {
	if len(comp) == 0 {
		return 1
	}
	dist := bfsDistancesInComp(adj, comp, comp[0])
	return periodFromDistances(adj, comp, dist)
}

func bfsDistancesInComp(adj [][]int, comp []int, root int) map[int]int {
	inComp := map[int]bool{}
	for _, v := range comp {
		inComp[v] = true
	}
	dist := map[int]int{root: 0}
	queue := []int{root}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, w := range adj[v] {
			if !inComp[w] {
				continue
			}
			if _, ok := dist[w]; !ok {
				dist[w] = dist[v] + 1
				queue = append(queue, w)
			}
		}
	}
	return dist
}

func periodFromDistances(adj [][]int, comp []int, dist map[int]int) int {
	inComp := map[int]bool{}
	for _, v := range comp {
		inComp[v] = true
	}
	g := 0
	for _, v := range comp {
		for _, w := range adj[v] {
			if !inComp[w] {
				continue
			}
			cycleLen := dist[v] + 1 - dist[w]
			g = gcd(g, cycleLen)
		}
	}
	if g == 0 {
		return 1
	}
	if g < 0 {
		return -g
	}
	return g
}
