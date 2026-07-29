package catrace_test

import (
	"fmt"
	"testing"

	"github.com/stephen-mcelhose/catrace"
	"gonum.org/v1/gonum/mat"
)

func TestTraceMatchesRestrictionTheoremPattern(t *testing.T) {
	parent, err := catrace.NewKernel(mat.NewDense(4, 4, []float64{
		0.60, 0.20, 0.10, 0.10,
		0.15, 0.55, 0.15, 0.15,
		0.20, 0.20, 0.40, 0.20,
		0.10, 0.20, 0.20, 0.50,
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	trace, err := parent.Trace([]int{0, 1}, 1e-12)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := trace.IsTraceOf(parent, []int{0, 1}, 1e-12); err != nil || !ok {
		t.Fatalf("trace verification failed: ok=%v err=%v", ok, err)
	}
	piParent, err := parent.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatal(err)
	}
	piTrace, err := trace.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatal(err)
	}
	restricted, err := catrace.RestrictDistribution(piParent, []int{0, 1}, 1e-12)
	if err != nil {
		t.Fatal(err)
	}
	for i := range piTrace {
		if diff := piTrace[i] - restricted[i]; diff < -1e-8 || diff > 1e-8 {
			t.Fatalf("stationary mismatch at %d: got %g want %g", i, piTrace[i], restricted[i])
		}
	}
}

func ExampleAgent_QualiaKernel() {
	D := mat.NewDense(2, 2, []float64{
		0.9, 0.1,
		0.2, 0.8,
	})
	A := mat.NewDense(2, 2, []float64{
		0.7, 0.3,
		0.4, 0.6,
	})
	P := mat.NewDense(2, 2, []float64{
		0.8, 0.2,
		0.1, 0.9,
	})
	agent := &catrace.Agent{D: D, A: A, P: P}
	Q, _ := agent.QualiaKernel()
	fmt.Printf("%.3f %.3f\n", Q.P.At(0, 0), Q.P.At(0, 1))
	fmt.Printf("%.3f %.3f\n", Q.P.At(1, 0), Q.P.At(1, 1))
	// Output:
	// 0.569 0.431
	// 0.422 0.578
}

func ExampleKernel_Trace() {
	k, _ := catrace.NewKernel(mat.NewDense(3, 3, []float64{
		0.6, 0.3, 0.1,
		0.2, 0.6, 0.2,
		0.4, 0.1, 0.5,
	}), nil)
	tr, _ := k.Trace([]int{0, 1}, 1e-12)
	fmt.Printf("%.3f %.3f\n", tr.P.At(0, 0), tr.P.At(0, 1))
	fmt.Printf("%.3f %.3f\n", tr.P.At(1, 0), tr.P.At(1, 1))
	// Output:
	// 0.680 0.320
	// 0.360 0.640
}
