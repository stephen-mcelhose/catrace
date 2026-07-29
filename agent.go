package catrace

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// Agent represents the three stochastic maps of the finite-state agent model.
//
// The maps form a closed loop over three spaces:
//
//	D : X → G  (decision:   experience → action)
//	A : G → W  (effect:     action     → world)
//	P : W → X  (perception: world      → experience)
type Agent struct {
	D *mat.Dense // decision kernel X → G
	A *mat.Dense // effect kernel G → W
	P *mat.Dense // perception kernel W → X

	XNames []string // experience state labels
	GNames []string // action state labels
	WNames []string // world state labels
}

// Validate checks that D, A, P are non-nil, dimensionally consistent, and row-stochastic.
func (a *Agent) Validate() error {
	if a == nil {
		return fmt.Errorf("nil agent")
	}
	if a.D == nil || a.A == nil || a.P == nil {
		return fmt.Errorf("agent requires non-nil D, A, P matrices")
	}
	x, g1 := a.D.Dims()
	g2, w1 := a.A.Dims()
	w2, x2 := a.P.Dims()
	if g1 != g2 {
		return fmt.Errorf("D columns (%d) must equal A rows (%d)", g1, g2)
	}
	if w1 != w2 {
		return fmt.Errorf("A columns (%d) must equal P rows (%d)", w1, w2)
	}
	if x != x2 {
		return fmt.Errorf("D rows (%d) must equal P columns (%d)", x, x2)
	}
	if _, err := NewRectKernel(a.D, a.XNames, a.GNames); err != nil {
		return fmt.Errorf("invalid D: %w", err)
	}
	if _, err := NewRectKernel(a.A, a.GNames, a.WNames); err != nil {
		return fmt.Errorf("invalid A: %w", err)
	}
	if _, err := NewRectKernel(a.P, a.WNames, a.XNames); err != nil {
		return fmt.Errorf("invalid P: %w", err)
	}
	return nil
}

// QualiaKernel computes Q = D·A·P, a square kernel on the experience space X.
func (a *Agent) QualiaKernel() (*Kernel, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	var da, dap mat.Dense
	da.Mul(a.D, a.A)
	dap.Mul(&da, a.P)
	return NewKernel(&dap, a.XNames)
}

// StrategyKernel computes S = A·P·D, a square kernel on the action space G.
func (a *Agent) StrategyKernel() (*Kernel, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	var ap, apd mat.Dense
	ap.Mul(a.A, a.P)
	apd.Mul(&ap, a.D)
	return NewKernel(&apd, a.GNames)
}

// WorldKernel computes W = P·D·A, a square kernel on the world space W.
func (a *Agent) WorldKernel() (*Kernel, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	var pd, pda mat.Dense
	pd.Mul(a.P, a.D)
	pda.Mul(&pd, a.A)
	return NewKernel(&pda, a.WNames)
}
