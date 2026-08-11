import { sanitizeOperandInput } from '../lib/operations'

type OperandInputProps = {
  id: string
  label: string
  value: string
  error?: string
  disabled?: boolean
  onChange: (value: string) => void
}

export function OperandInput({
  id,
  label,
  value,
  error,
  disabled = false,
  onChange,
}: OperandInputProps) {
  const showClear = value !== '' && !disabled

  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <div className="field-control">
        <input
          id={id}
          name={id}
          type="text"
          inputMode="decimal"
          autoComplete="off"
          spellCheck={false}
          value={value}
          disabled={disabled}
          aria-invalid={error ? true : undefined}
          aria-describedby={error ? `${id}-error` : undefined}
          onChange={(event) => onChange(sanitizeOperandInput(event.target.value))}
        />
        {showClear ? (
          <button
            type="button"
            className="field-clear"
            aria-label={`Clear ${label}`}
            onClick={() => onChange('')}
          >
            ×
          </button>
        ) : null}
      </div>
      {error ? (
        <p id={`${id}-error`} className="field-error" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  )
}
