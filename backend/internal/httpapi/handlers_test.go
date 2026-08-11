package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/daniOs110/fullstack-calculator/internal/calculator"
)

func newTestHandler() http.Handler {
	return NewRouter([]string{"http://localhost:5173"})
}

func postJSON(t *testing.T, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rr, req)
	return rr
}

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rr, req)
	return rr
}

func assertJSONContentType(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()

	contentType := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}
}

func assertResult(t *testing.T, rr *httptest.ResponseRecorder, want float64) {
	t.Helper()

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}
	assertJSONContentType(t, rr)

	var resp resultResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v; body = %s", err, rr.Body.String())
	}

	if resp.Result != want && math.Abs(resp.Result-want) >= 1e-9 {
		t.Errorf("result = %v, want %v", resp.Result, want)
	}
}

func assertError(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantMessage string) {
	t.Helper()

	if rr.Code != wantStatus {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, wantStatus, rr.Body.String())
	}
	assertJSONContentType(t, rr)

	var resp errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v; body = %s", err, rr.Body.String())
	}
	if resp.Error != wantMessage {
		t.Errorf("error = %q, want %q", resp.Error, wantMessage)
	}
}

func TestHappyPathEndpoints(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		want float64
	}{
		{name: "add", path: "/api/v1/add", body: `{"a":2,"b":3}`, want: 5},
		{name: "subtract", path: "/api/v1/subtract", body: `{"a":5,"b":3}`, want: 2},
		{name: "multiply", path: "/api/v1/multiply", body: `{"a":4,"b":3}`, want: 12},
		{name: "divide", path: "/api/v1/divide", body: `{"a":10,"b":4}`, want: 2.5},
		{name: "exponent", path: "/api/v1/exponent", body: `{"base":2,"exp":10}`, want: 1024},
		{name: "sqrt", path: "/api/v1/sqrt", body: `{"a":9}`, want: 3},
		{name: "percentage", path: "/api/v1/percentage", body: `{"a":10,"b":200}`, want: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertResult(t, postJSON(t, tt.path, tt.body), tt.want)
		})
	}
}

func TestDomainErrors(t *testing.T) {
	t.Run("divide by zero", func(t *testing.T) {
		assertError(t, postJSON(t, "/api/v1/divide", `{"a":10,"b":0}`),
			http.StatusBadRequest, calculator.ErrDivisionByZero.Error())
	})

	t.Run("sqrt of negative", func(t *testing.T) {
		assertError(t, postJSON(t, "/api/v1/sqrt", `{"a":-9}`),
			http.StatusBadRequest, calculator.ErrNegativeSquareRoot.Error())
	})

	t.Run("exponent with no real result", func(t *testing.T) {
		assertError(t, postJSON(t, "/api/v1/exponent", `{"base":-8,"exp":0.5}`),
			http.StatusBadRequest, "the result is not a finite number")
	})
}

func TestAddInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantMsg string
	}{
		{name: "missing field", body: `{"a":1}`, wantMsg: `field "b" is required`},
		{name: "null field", body: `{"a":1,"b":null}`, wantMsg: `field "b" is required`},
		{name: "wrong type", body: `{"a":"x","b":2}`, wantMsg: `field "a" must be a number`},
		{name: "unexpected field", body: `{"a":1,"b":2,"c":3}`, wantMsg: "body contains an unexpected field"},
		{name: "empty body", body: ``, wantMsg: "body must contain a JSON object"},
		{name: "malformed json", body: `{not json`, wantMsg: "body must be valid JSON"},
		{name: "concatenated objects", body: `{"a":1,"b":2}{"a":3,"b":4}`, wantMsg: "body must contain a single JSON object"},
		{
			name:    "operand above MaxOperand",
			body:    `{"a":1e20,"b":2}`,
			wantMsg: fmt.Sprintf("field %q must be between %g and %g", "a", -calculator.MaxOperand, calculator.MaxOperand),
		},
		{name: "float64 overflow", body: `{"a":1e400,"b":2}`, wantMsg: `field "a" is out of range`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertError(t, postJSON(t, "/api/v1/add", tt.body), http.StatusBadRequest, tt.wantMsg)
		})
	}

	t.Run("body larger than limit", func(t *testing.T) {
		// Build a valid JSON object whose size exceeds maxBodyBytes.
		padding := strings.Repeat("0", maxBodyBytes)
		body := `{"a":1,"b":2,"pad":"` + padding + `"}`
		assertError(t, postJSON(t, "/api/v1/add", body), http.StatusBadRequest, "body is too large")
	})
}

func TestDistinctBodyShapes(t *testing.T) {
	t.Run("exponent missing base", func(t *testing.T) {
		assertError(t, postJSON(t, "/api/v1/exponent", `{"exp":2}`),
			http.StatusBadRequest, `field "base" is required`)
	})

	t.Run("exponent missing exp", func(t *testing.T) {
		assertError(t, postJSON(t, "/api/v1/exponent", `{"base":2}`),
			http.StatusBadRequest, `field "exp" is required`)
	})

	t.Run("sqrt missing a", func(t *testing.T) {
		assertError(t, postJSON(t, "/api/v1/sqrt", `{}`),
			http.StatusBadRequest, `field "a" is required`)
	})
}

func TestValidateOperandNonFinite(t *testing.T) {
	nan := math.NaN()
	posInf := math.Inf(1)
	negInf := math.Inf(-1)

	tests := []struct {
		name  string
		value float64
	}{
		{name: "NaN", value: nan},
		{name: "positive infinity", value: posInf},
		{name: "negative infinity", value: negInf},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.value
			err := validateOperand(operand{name: "a", value: &value})
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			want := `field "a" must be a finite number`
			if err.Error() != want {
				t.Errorf("error = %q, want %q", err.Error(), want)
			}
		})
	}
}

func TestNegativeZeroNormalizedInResponse(t *testing.T) {
	rr := postJSON(t, "/api/v1/multiply", `{"a":-4,"b":0}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	raw := bytes.TrimSpace(rr.Body.Bytes())
	if strings.Contains(string(raw), "-0") {
		t.Fatalf("response body contains negative zero: %s", raw)
	}

	var resp resultResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Result != 0 || math.Signbit(resp.Result) {
		t.Errorf("result = %v (signbit=%v), want +0", resp.Result, math.Signbit(resp.Result))
	}
}

func TestRouter(t *testing.T) {
	t.Run("wrong method returns 405", func(t *testing.T) {
		rr := get(t, "/api/v1/add")
		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("unknown route returns 404", func(t *testing.T) {
		rr := get(t, "/api/v1/unknown")
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("health returns ok", func(t *testing.T) {
		rr := get(t, "/health")
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body = %s", rr.Code, http.StatusOK, rr.Body.String())
		}
		assertJSONContentType(t, rr)

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		var resp map[string]string
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp["status"] != "ok" {
			t.Errorf("status = %q, want %q", resp["status"], "ok")
		}
	})
}
