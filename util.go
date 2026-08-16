package catrace

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

func cloneDense(m *mat.Dense) *mat.Dense {
	if m == nil {
		return nil
	}
	r, c := m.Dims()
	data := make([]float64, r*c)
	copy(data, m.RawMatrix().Data)
	return mat.NewDense(r, c, data)
}

func cloneStrings(xs []string) []string {
	if xs == nil {
		return nil
	}
	out := make([]string, len(xs))
	copy(out, xs)
	return out
}

func defaultNames(prefix string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("%s%d", prefix, i)
	}
	return out
}

func nearlyEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func normalizeVector(x []float64, tol float64) error {
	sum := 0.0
	for _, v := range x {
		if v < -tol {
			return fmt.Errorf("negative entry %g", v)
		}
		if v > 0 {
			sum += v
		}
	}
	if sum <= tol {
		return fmt.Errorf("vector has near-zero sum %g", sum)
	}
	for i, v := range x {
		if v < 0 && math.Abs(v) <= tol {
			v = 0
		}
		x[i] = v / sum
	}
	return nil
}

// cloneProbDist copies v into a new probability distribution of length n.
// When clampTiny is true, entries in (−tol, 0) are clamped to zero before the
// negativity check. The result must sum to 1 within tol. label is used in
// error messages (e.g. "start", "restart").
func cloneProbDist(v []float64, n int, tol float64, label string, clampTiny bool) ([]float64, error) {
	if len(v) != n {
		return nil, fmt.Errorf("%s length %d does not match kernel size %d", label, len(v), n)
	}
	out := make([]float64, n)
	sum := 0.0
	for i, val := range v {
		if clampTiny && val < 0 && math.Abs(val) <= tol {
			val = 0
		}
		if val < 0 {
			return nil, fmt.Errorf("%s[%d] is negative: %g", label, i, val)
		}
		out[i] = val
		sum += val
	}
	if math.Abs(sum-1.0) > tol {
		if clampTiny {
			return nil, fmt.Errorf("%s does not sum to 1 (got %g, tol %g)", label, sum, tol)
		}
		return nil, fmt.Errorf("%s does not sum to 1 (got %g)", label, sum)
	}
	return out, nil
}

func rowSums(m *mat.Dense) []float64 {
	r, c := m.Dims()
	out := make([]float64, r)
	for i := 0; i < r; i++ {
		s := 0.0
		for j := 0; j < c; j++ {
			s += m.At(i, j)
		}
		out[i] = s
	}
	return out
}

func submatrix(m *mat.Dense, rows, cols []int) *mat.Dense {
	out := mat.NewDense(len(rows), len(cols), nil)
	for i, r := range rows {
		for j, c := range cols {
			out.Set(i, j, m.At(r, c))
		}
	}
	return out
}

func complementIndices(n int, subset []int) []int {
	keep := make(map[int]bool, len(subset))
	for _, v := range subset {
		keep[v] = true
	}
	out := make([]int, 0, n-len(subset))
	for i := 0; i < n; i++ {
		if !keep[i] {
			out = append(out, i)
		}
	}
	return out
}

func sortedCopy(xs []int) []int {
	out := make([]int, len(xs))
	copy(out, xs)
	sort.Ints(out)
	return out
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
