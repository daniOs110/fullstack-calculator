/** Shared request/response shapes for the calculator API. */

export type TwoOperandsRequest = {
  a: number
  b: number
}

export type ExponentRequest = {
  base: number
  exp: number
}

export type SingleOperandRequest = {
  a: number
}

export type ResultResponse = {
  result: number
}

export type ErrorResponse = {
  error: string
}

export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'exponent'
  | 'sqrt'
  | 'percentage'
