package httpapi

// Operands are pointers so a missing field can be told apart from a zero
// value: {"a": 1} without "b" must fail validation instead of dividing by 0.
type twoOperandsRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

type exponentRequest struct {
	Base *float64 `json:"base"`
	Exp  *float64 `json:"exp"`
}

type singleOperandRequest struct {
	A *float64 `json:"a"`
}

type resultResponse struct {
	Result float64 `json:"result"`
}
