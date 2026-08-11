import type { Operation } from '../api/types'

/** Same bound as backend calculator.MaxOperand. */
export const MAX_OPERAND = 1e15

/** Keeps typed integers within a magnitude compatible with MaxOperand. */
export const MAX_INTEGER_DIGITS = 15

/** Fractional digits allowed while typing. */
export const MAX_FRACTION_DIGITS = 10

export type FieldKey = 'a' | 'b' | 'base' | 'exp'

export type OperationMeta = {
  value: Operation
  label: string
  fields: readonly FieldKey[]
  fieldLabels: Readonly<Partial<Record<FieldKey, string>>>
}

export const OPERATIONS: readonly OperationMeta[] = [
  {
    value: 'add',
    label: 'Add',
    fields: ['a', 'b'],
    fieldLabels: { a: 'First number', b: 'Second number' },
  },
  {
    value: 'subtract',
    label: 'Subtract',
    fields: ['a', 'b'],
    fieldLabels: { a: 'First number', b: 'Second number' },
  },
  {
    value: 'multiply',
    label: 'Multiply',
    fields: ['a', 'b'],
    fieldLabels: { a: 'First number', b: 'Second number' },
  },
  {
    value: 'divide',
    label: 'Divide',
    fields: ['a', 'b'],
    fieldLabels: { a: 'First number', b: 'Second number' },
  },
  {
    value: 'exponent',
    label: 'Exponent',
    fields: ['base', 'exp'],
    fieldLabels: { base: 'Base', exp: 'Exponent' },
  },
  {
    value: 'sqrt',
    label: 'Square root',
    fields: ['a'],
    fieldLabels: { a: 'Number' },
  },
  {
    value: 'percentage',
    label: 'Percentage (a% of b)',
    fields: ['a', 'b'],
    fieldLabels: { a: 'Percentage (%)', b: 'Of (number)' },
  },
] as const

export function metaFor(operation: Operation): OperationMeta {
  const found = OPERATIONS.find((op) => op.value === operation)
  if (!found) {
    throw new Error(`unknown operation: ${operation}`)
  }
  return found
}

/**
 * Builds the JSON body with ONLY the fields the selected operation needs.
 * Hidden fields must not appear here — the backend rejects unknown keys.
 */
export function buildRequestBody(
  operation: Operation,
  values: Partial<Record<FieldKey, number>>,
): Record<string, number> {
  const body: Record<string, number> = {}
  for (const field of metaFor(operation).fields) {
    const value = values[field]
    if (value === undefined) {
      throw new Error(`missing value for field "${field}"`)
    }
    body[field] = value
  }
  return body
}

/**
 * Keeps only a valid partial number while the user types or pastes:
 * optional leading "-", digits, one ".", max integer/fraction lengths.
 * Letters and other symbols are dropped.
 */
export function sanitizeOperandInput(raw: string): string {
  if (raw === '') {
    return ''
  }

  let sign = ''
  let rest = raw
  if (rest.startsWith('-')) {
    sign = '-'
    rest = rest.slice(1)
  }

  rest = rest.replace(/[^\d.]/g, '')

  const dotIndex = rest.indexOf('.')
  let integerPart: string
  let fractionPart: string | undefined

  if (dotIndex === -1) {
    integerPart = rest
  } else {
    integerPart = rest.slice(0, dotIndex)
    fractionPart = rest.slice(dotIndex + 1).replace(/\./g, '')
  }

  integerPart = integerPart.slice(0, MAX_INTEGER_DIGITS)

  if (fractionPart !== undefined) {
    fractionPart = fractionPart.slice(0, MAX_FRACTION_DIGITS)
    return `${sign}${integerPart}.${fractionPart}`
  }

  return `${sign}${integerPart}`
}

export function parseOperand(
  raw: string,
): { ok: true; value: number } | { ok: false; error: string } {
  const trimmed = raw.trim()
  if (trimmed === '') {
    return { ok: false, error: 'required' }
  }

  if (trimmed === '-' || trimmed === '.' || trimmed === '-.') {
    return { ok: false, error: 'must be a finite number' }
  }

  const value = Number(trimmed)
  if (!Number.isFinite(value)) {
    return { ok: false, error: 'must be a finite number' }
  }

  if (Math.abs(value) > MAX_OPERAND) {
    return {
      ok: false,
      error: `must be between ${-MAX_OPERAND} and ${MAX_OPERAND}`,
    }
  }

  return { ok: true, value }
}
