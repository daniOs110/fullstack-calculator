import { useMemo, useState, type FormEvent } from 'react'
import { ApiError, calculate } from '../api/client'
import type { Operation } from '../api/types'
import {
  buildRequestBody,
  metaFor,
  parseOperand,
  type FieldKey,
} from '../lib/operations'
import { OperandInput } from './OperandInput'
import { OperationSelect } from './OperationSelect'
import { ResultDisplay } from './ResultDisplay'

const EMPTY_VALUES: Record<FieldKey, string> = {
  a: '',
  b: '',
  base: '',
  exp: '',
}

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [values, setValues] = useState<Record<FieldKey, string>>(EMPTY_VALUES)
  const [result, setResult] = useState<number | null>(null)
  const [apiError, setApiError] = useState<string | null>(null)
  const [submitted, setSubmitted] = useState(false)
  const [loading, setLoading] = useState(false)

  const operationMeta = metaFor(operation)
  const fields = operationMeta.fields

  const fieldErrors = useMemo(() => {
    const errors: Partial<Record<FieldKey, string>> = {}
    for (const field of fields) {
      const parsed = parseOperand(values[field])
      if (!parsed.ok) {
        errors[field] = parsed.error
      }
    }
    return errors
  }, [fields, values])

  const hasLocalErrors = fields.some((field) => fieldErrors[field])

  function handleOperationChange(next: Operation) {
    setOperation(next)
    setResult(null)
    setApiError(null)
    setSubmitted(false)
    // Values for inactive fields stay in React state only; buildRequestBody
    // never includes them in the JSON payload.
  }

  function handleFieldChange(field: FieldKey, value: string) {
    setValues((current) => ({ ...current, [field]: value }))
    setResult(null)
    setApiError(null)
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitted(true)
    setApiError(null)
    setResult(null)

    if (hasLocalErrors) {
      return
    }

    const parsedValues: Partial<Record<FieldKey, number>> = {}
    for (const field of fields) {
      const parsed = parseOperand(values[field])
      if (!parsed.ok) {
        return
      }
      parsedValues[field] = parsed.value
    }

    // Critical: only active fields go into the body (e.g. sqrt => { a } only).
    const body = buildRequestBody(operation, parsedValues)

    setLoading(true)
    try {
      const nextResult = await calculate(operation, body)
      setResult(nextResult)
    } catch (error) {
      if (error instanceof ApiError) {
        setApiError(error.message)
      } else if (error instanceof Error) {
        setApiError(error.message)
      } else {
        setApiError('unexpected error')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="calculator" onSubmit={handleSubmit} noValidate>
      <OperationSelect
        value={operation}
        disabled={loading}
        onChange={handleOperationChange}
      />

      {fields.map((field) => (
        <OperandInput
          key={`${operation}-${field}`}
          id={`operand-${field}`}
          label={operationMeta.fieldLabels[field] ?? field}
          value={values[field]}
          disabled={loading}
          error={submitted ? fieldErrors[field] : undefined}
          onChange={(value) => handleFieldChange(field, value)}
        />
      ))}

      <button type="submit" className="submit" disabled={loading}>
        {loading ? 'Calculating…' : 'Calculate'}
      </button>

      <ResultDisplay result={result} error={apiError} />
    </form>
  )
}
