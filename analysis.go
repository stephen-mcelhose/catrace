package catrace

import "fmt"

// ClassDecomposition summarizes communicating classes and recurrence structure.
type ClassDecomposition struct {
	SCCs      [][]int
	Recurrent [][]int
	Transient []int
	Periods   map[int]int
}

func (k *Kernel) Classes(tol float64) (*ClassDecomposition, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	n := k.NumStates()
	adj := make([][]int, n)
	rev := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if k.P.At(i, j) > tol {
				adj[i] = append(adj[i], j)
				rev[j] = append(rev[j], i)
			}
		}
	}

	visited := make([]bool, n)
	order := make([]int, 0, n)
	var dfs1 func(int)
	dfs1 = func(v int) {
		visited[v] = true
		for _, w := range adj[v] {
			if !visited[w] {
				dfs1(w)
			}
		}
		order = append(order, v)
	}
	for i := 0; i < n; i++ {
		if !visited[i] {
			dfs1(i)
		}
	}

	for i := range visited {
		visited[i] = false
	}
	var sccs [][]int
	var dfs2 func(int, *[]int)
	dfs2 = func(v int, comp *[]int) {
		visited[v] = true
		*comp = append(*comp, v)
		for _, w := range rev[v] {
			if !visited[w] {
				dfs2(w, comp)
			}
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if !visited[v] {
			comp := []int{}
			dfs2(v, &comp)
			sccs = append(sccs, sortedCopy(comp))
		}
	}

	classOf := make(map[int]int, n)
	for idx, comp := range sccs {
		for _, v := range comp {
			classOf[v] = idx
		}
	}
	recurrent := [][]int{}
	transient := []int{}
	periods := map[int]int{}
	for idx, comp := range sccs {
		closed := true
		for _, v := range comp {
			for _, w := range adj[v] {
				if classOf[w] != idx {
					closed = false
					break
				}
			}
			if !closed {
				break
			}
		}
		if closed {
			recurrent = append(recurrent, comp)
			periods[idx] = periodOfClass(adj, comp)
		} else {
			transient = append(transient, comp...)
		}
	}
	return &ClassDecomposition{SCCs: sccs, Recurrent: recurrent, Transient: sortedCopy(transient), Periods: periods}, nil
}

func periodOfClass(adj [][]int, comp []int) int {
	if len(comp) == 0 {
		return 1
	}
	inComp := map[int]bool{}
	for _, v := range comp {
		inComp[v] = true
	}
	root := comp[0]
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
