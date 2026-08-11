// Package calculator implements the arithmetic operations exposed by the API,
// independent from any transport or serialization concern.
package calculator

import (
	"errors"
	"math"
)

// MaxOperand bounds the magnitude accepted for any operand. It sits just below
// 2^53, the largest integer a float64 represents exactly; past that point
// x+1 == x can hold and results stop matching what a user would expect.
const MaxOperand = 1e15

// Errors reported when an operation is undefined for the given operands.
var (
	ErrDivisionByZero     = errors.New("division by zero")
	ErrNegativeSquareRoot = errors.New("square root of a negative number")
)

// Add returns a + b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns a - b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns a * b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns a / b, or ErrDivisionByZero when b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// Exponent returns base raised to exp. Operand combinations without a real
// result, such as a negative base with a fractional exponent, yield NaN, and
// large exponents overflow to infinity; callers are expected to reject
// non-finite results.
func Exponent(base, exp float64) float64 {
	return math.Pow(base, exp)
}

// Sqrt returns the square root of a, or ErrNegativeSquareRoot when a is
// negative.
func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return math.Sqrt(a), nil
}

// Percentage returns a percent of b.
func Percentage(a, b float64) float64 {
	return a / 100 * b
}
