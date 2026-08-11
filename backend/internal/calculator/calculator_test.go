package calculator

import (
	"errors"
	"math"
	"testing"
)

const relativeTolerance = 1e-9

// assertAlmostEqual compares with a relative tolerance because floating point
// arithmetic rarely lands on the exact decimal value: 0.1+0.2 is not 0.3.
func assertAlmostEqual(t *testing.T, got, want float64) {
	t.Helper()

	if got == want {
		return
	}

	diff := math.Abs(got - want)
	if want == 0 {
		if diff >= relativeTolerance {
			t.Errorf("got %v, want %v", got, want)
		}
		return
	}

	if diff/math.Abs(want) >= relativeTolerance {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want float64
	}{
		{name: "positive operands", a: 2, b: 3, want: 5},
		{name: "negative operands", a: -2, b: -3, want: -5},
		{name: "decimal operands", a: 0.1, b: 0.2, want: 0.3},
		{name: "zero is the identity", a: 5, b: 0, want: 5},
		{name: "operands at the maximum", a: MaxOperand, b: MaxOperand, want: 2e15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAlmostEqual(t, Add(tt.a, tt.b), tt.want)
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want float64
	}{
		{name: "positive operands", a: 5, b: 3, want: 2},
		{name: "negative result", a: 3, b: 5, want: -2},
		{name: "subtracting zero", a: 5, b: 0, want: 5},
		{name: "a number minus itself", a: 5, b: 5, want: 0},
		{name: "decimal operands", a: 1.5, b: 0.3, want: 1.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAlmostEqual(t, Subtract(tt.a, tt.b), tt.want)
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want float64
	}{
		{name: "positive operands", a: 4, b: 3, want: 12},
		{name: "by zero", a: 4, b: 0, want: 0},
		{name: "negative times positive", a: -4, b: 3, want: -12},
		{name: "negative times negative", a: -4, b: -3, want: 12},
		{name: "one is the identity", a: 7, b: 1, want: 7},
		{name: "decimal operands", a: 0.5, b: 4, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAlmostEqual(t, Multiply(tt.a, tt.b), tt.want)
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		want    float64
		wantErr error
	}{
		{name: "exact division", a: 10, b: 2, want: 5},
		{name: "non exact division", a: 10, b: 3, want: 3.3333333333333335},
		{name: "zero dividend", a: 0, b: 5, want: 0},
		{name: "negative divisor", a: 10, b: -2, want: -5},
		{name: "decimal operands", a: 7.5, b: 2.5, want: 3},
		{name: "zero divisor", a: 10, b: 0, wantErr: ErrDivisionByZero},
		// A -0.0 literal folds to 0, so the negative zero is built explicitly.
		{name: "negative zero divisor", a: 10, b: math.Copysign(0, -1), wantErr: ErrDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			assertAlmostEqual(t, got, tt.want)
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{name: "perfect square", a: 9, want: 3},
		{name: "zero", a: 0, want: 0},
		{name: "non perfect square", a: 2, want: 1.4142135623730951},
		{name: "decimal operand", a: 0.25, want: 0.5},
		{name: "operand at the maximum", a: MaxOperand, want: 3.1622776601683795e7},
		{name: "negative operand", a: -1, wantErr: ErrNegativeSquareRoot},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sqrt(tt.a)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			assertAlmostEqual(t, got, tt.want)
		})
	}
}

func TestExponent(t *testing.T) {
	tests := []struct {
		name string
		base float64
		exp  float64
		want float64
	}{
		{name: "positive base and exponent", base: 2, exp: 3, want: 8},
		{name: "exponent zero returns one", base: 5, exp: 0, want: 1},
		{name: "exponent one returns the base", base: 5, exp: 1, want: 5},
		{name: "negative exponent", base: 2, exp: -2, want: 0.25},
		{name: "fractional exponent", base: 9, exp: 0.5, want: 3},
		{name: "negative base with integer exponent", base: -2, exp: 3, want: -8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAlmostEqual(t, Exponent(tt.base, tt.exp), tt.want)
		})
	}
}

// TestExponentNonFiniteResults pins down the operand combinations that produce
// values JSON cannot encode. The HTTP layer relies on this contract to reject
// them before serializing a response.
func TestExponentNonFiniteResults(t *testing.T) {
	t.Run("negative base with fractional exponent has no real result", func(t *testing.T) {
		if got := Exponent(-8, 0.5); !math.IsNaN(got) {
			t.Errorf("got %v, want NaN", got)
		}
	})

	t.Run("zero to a negative exponent overflows", func(t *testing.T) {
		if got := Exponent(0, -1); !math.IsInf(got, 1) {
			t.Errorf("got %v, want +Inf", got)
		}
	})

	t.Run("large operands overflow", func(t *testing.T) {
		if got := Exponent(MaxOperand, MaxOperand); !math.IsInf(got, 1) {
			t.Errorf("got %v, want +Inf", got)
		}
	})
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want float64
	}{
		{name: "ten percent of a number", a: 10, b: 200, want: 20},
		{name: "zero percent", a: 0, b: 200, want: 0},
		{name: "one hundred percent returns the number", a: 100, b: 200, want: 200},
		{name: "more than one hundred percent", a: 150, b: 200, want: 300},
		{name: "negative percentage", a: -10, b: 200, want: -20},
		{name: "decimal percentage", a: 12.5, b: 80, want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAlmostEqual(t, Percentage(tt.a, tt.b), tt.want)
		})
	}
}
