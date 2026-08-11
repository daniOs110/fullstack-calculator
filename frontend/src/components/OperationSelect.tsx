import { OPERATIONS } from '../lib/operations'
import type { Operation } from '../api/types'

type OperationSelectProps = {
  value: Operation
  disabled?: boolean
  onChange: (operation: Operation) => void
}

export function OperationSelect({
  value,
  disabled = false,
  onChange,
}: OperationSelectProps) {
  return (
    <div className="field">
      <label htmlFor="operation">Operation</label>
      <select
        id="operation"
        name="operation"
        value={value}
        disabled={disabled}
        onChange={(event) => onChange(event.target.value as Operation)}
      >
        {OPERATIONS.map((op) => (
          <option key={op.value} value={op.value}>
            {op.label}
          </option>
        ))}
      </select>
    </div>
  )
}
