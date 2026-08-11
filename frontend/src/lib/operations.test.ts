import { describe, expect, it } from 'vitest'
import {
  MAX_FRACTION_DIGITS,
  MAX_INTEGER_DIGITS,
  MAX_OPERAND,
  buildRequestBody,
  parseOperand,
  sanitizeOperandInput,
} from './operations'

describe('sanitizeOperandInput', () => {
  it('keeps a valid number', () => {
    expect(sanitizeOperandInput('12.5')).toBe('12.5')
    expect(sanitizeOperandInput('-3')).toBe('-3')
  })

  it('strips letters', () => {
    expect(sanitizeOperandInput('12abc')).toBe('12')
  })

  it('strips disallowed symbols', () => {
    expect(sanitizeOperandInput('1,2')).toBe('12')
    expect(sanitizeOperandInput('1$2')).toBe('12')
  })

  it('allows only one decimal point', () => {
    expect(sanitizeOperandInput('1.2.3')).toBe('1.23')
  })

  it('allows a leading minus only', () => {
    expect(sanitizeOperandInput('-12')).toBe('-12')
    expect(sanitizeOperandInput('12-3')).toBe('123')
  })

  it(`caps the integer part at ${MAX_INTEGER_DIGITS} digits`, () => {
    const tooLong = '1'.repeat(MAX_INTEGER_DIGITS + 3)
    expect(sanitizeOperandInput(tooLong)).toBe('1'.repeat(MAX_INTEGER_DIGITS))
  })

  it(`caps the fraction part at ${MAX_FRACTION_DIGITS} digits`, () => {
    const fraction = '1'.repeat(MAX_FRACTION_DIGITS + 4)
    expect(sanitizeOperandInput(`1.${fraction}`)).toBe(
      `1.${'1'.repeat(MAX_FRACTION_DIGITS)}`,
    )
  })

  it('returns an empty string for empty input', () => {
    expect(sanitizeOperandInput('')).toBe('')
  })
})

describe('parseOperand', () => {
  it('rejects an empty value', () => {
    expect(parseOperand('')).toEqual({ ok: false, error: 'required' })
    expect(parseOperand('   ')).toEqual({ ok: false, error: 'required' })
  })

  it('rejects incomplete number drafts', () => {
    expect(parseOperand('-')).toEqual({
      ok: false,
      error: 'must be a finite number',
    })
    expect(parseOperand('.')).toEqual({
      ok: false,
      error: 'must be a finite number',
    })
    expect(parseOperand('-.')).toEqual({
      ok: false,
      error: 'must be a finite number',
    })
  })

  it('accepts a valid number', () => {
    expect(parseOperand('12.5')).toEqual({ ok: true, value: 12.5 })
    expect(parseOperand('-3')).toEqual({ ok: true, value: -3 })
  })

  it('rejects values outside MaxOperand', () => {
    const over = String(MAX_OPERAND * 10)
    expect(parseOperand(over)).toEqual({
      ok: false,
      error: `must be between ${-MAX_OPERAND} and ${MAX_OPERAND}`,
    })
  })
})

describe('buildRequestBody', () => {
  it('includes only a and b for add', () => {
    expect(buildRequestBody('add', { a: 2, b: 3 })).toEqual({ a: 2, b: 3 })
  })

  it('includes only a for sqrt', () => {
    expect(buildRequestBody('sqrt', { a: 9, b: 99 })).toEqual({ a: 9 })
  })

  it('includes only base and exp for exponent', () => {
    expect(
      buildRequestBody('exponent', { base: 2, exp: 3, a: 1, b: 1 }),
    ).toEqual({ base: 2, exp: 3 })
  })

  it('throws when a required field is missing', () => {
    expect(() => buildRequestBody('add', { a: 1 })).toThrow(
      'missing value for field "b"',
    )
  })
})
