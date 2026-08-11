package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/daniOs110/fullstack-calculator/internal/calculator"
)

// maxBodyBytes caps the request body; these endpoints never need more than a
// couple of numbers, so anything larger is rejected before it is parsed.
const maxBodyBytes = 4 << 10

type operand struct {
	name  string
	value *float64
}

// decodeTwoOperands reads and validates an {"a", "b"} body. It writes the
// error response itself, so callers only need to check the boolean.
func decodeTwoOperands(w http.ResponseWriter, r *http.Request) (float64, float64, bool) {
	var req twoOperandsRequest
	if !decodeJSON(w, r, &req) {
		return 0, 0, false
	}

	if !validateOperands(w, operand{"a", req.A}, operand{"b", req.B}) {
		return 0, 0, false
	}

	return *req.A, *req.B, true
}

// decodeExponentOperands reads and validates a {"base", "exp"} body.
func decodeExponentOperands(w http.ResponseWriter, r *http.Request) (float64, float64, bool) {
	var req exponentRequest
	if !decodeJSON(w, r, &req) {
		return 0, 0, false
	}

	if !validateOperands(w, operand{"base", req.Base}, operand{"exp", req.Exp}) {
		return 0, 0, false
	}

	return *req.Base, *req.Exp, true
}

// decodeSingleOperand reads and validates an {"a"} body.
func decodeSingleOperand(w http.ResponseWriter, r *http.Request) (float64, bool) {
	var req singleOperandRequest
	if !decodeJSON(w, r, &req) {
		return 0, false
	}

	if !validateOperands(w, operand{"a", req.A}) {
		return 0, false
	}

	return *req.A, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, decodeErrorMessage(err))
		return false
	}

	if decoder.More() {
		writeError(w, http.StatusBadRequest, "body must contain a single JSON object")
		return false
	}

	return true
}

func validateOperands(w http.ResponseWriter, operands ...operand) bool {
	for _, op := range operands {
		if err := validateOperand(op); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return false
		}
	}
	return true
}

func validateOperand(op operand) error {
	if op.value == nil {
		return fmt.Errorf("field %q is required", op.name)
	}

	value := *op.value

	// NaN fails every comparison, so it has to be caught before the range check.
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("field %q must be a finite number", op.name)
	}

	if math.Abs(value) > calculator.MaxOperand {
		return fmt.Errorf("field %q must be between %g and %g", op.name, -calculator.MaxOperand, calculator.MaxOperand)
	}

	return nil
}

func decodeErrorMessage(err error) string {
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	var maxBytesErr *http.MaxBytesError

	switch {
	case errors.Is(err, io.EOF):
		return "body must contain a JSON object"
	case errors.As(err, &typeErr):
		if strings.HasPrefix(typeErr.Value, "number") {
			return fmt.Sprintf("field %q is out of range", typeErr.Field)
		}
		return fmt.Sprintf("field %q must be a number", typeErr.Field)
	case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
		return "body must be valid JSON"
	case errors.As(err, &maxBytesErr):
		return "body is too large"
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		return "body contains an unexpected field"
	default:
		return "body could not be parsed"
	}
}
