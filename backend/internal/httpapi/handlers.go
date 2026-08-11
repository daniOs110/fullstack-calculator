package httpapi

import (
	"net/http"

	"github.com/daniOs110/fullstack-calculator/internal/calculator"
)

func handleAdd(w http.ResponseWriter, r *http.Request) {
	a, b, ok := decodeTwoOperands(w, r)
	if !ok {
		return
	}

	writeResult(w, calculator.Add(a, b))
}

func handleSubtract(w http.ResponseWriter, r *http.Request) {
	a, b, ok := decodeTwoOperands(w, r)
	if !ok {
		return
	}

	writeResult(w, calculator.Subtract(a, b))
}

func handleMultiply(w http.ResponseWriter, r *http.Request) {
	a, b, ok := decodeTwoOperands(w, r)
	if !ok {
		return
	}

	writeResult(w, calculator.Multiply(a, b))
}

func handleDivide(w http.ResponseWriter, r *http.Request) {
	a, b, ok := decodeTwoOperands(w, r)
	if !ok {
		return
	}

	// Every error the calculator reports is caused by the operands themselves.
	result, err := calculator.Divide(a, b)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResult(w, result)
}

func handleExponent(w http.ResponseWriter, r *http.Request) {
	base, exp, ok := decodeExponentOperands(w, r)
	if !ok {
		return
	}

	writeResult(w, calculator.Exponent(base, exp))
}

func handleSqrt(w http.ResponseWriter, r *http.Request) {
	a, ok := decodeSingleOperand(w, r)
	if !ok {
		return
	}

	result, err := calculator.Sqrt(a)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResult(w, result)
}

func handlePercentage(w http.ResponseWriter, r *http.Request) {
	a, b, ok := decodeTwoOperands(w, r)
	if !ok {
		return
	}

	writeResult(w, calculator.Percentage(a, b))
}
