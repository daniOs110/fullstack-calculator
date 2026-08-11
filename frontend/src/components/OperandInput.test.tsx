import { useState } from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { OperandInput } from './OperandInput'

function ControlledOperandInput({
  label = 'First number',
  initial = '',
  error,
  onChange,
}: {
  label?: string
  initial?: string
  error?: string
  onChange?: (value: string) => void
}) {
  const [value, setValue] = useState(initial)
  return (
    <OperandInput
      id="a"
      label={label}
      value={value}
      error={error}
      onChange={(next) => {
        setValue(next)
        onChange?.(next)
      }}
    />
  )
}

describe('OperandInput', () => {
  it('renders the label', () => {
    render(<ControlledOperandInput />)
    expect(screen.getByLabelText('First number')).toBeInTheDocument()
  })

  it('keeps only numeric characters while typing', async () => {
    const user = userEvent.setup()
    const onChange = vi.fn()

    render(<ControlledOperandInput onChange={onChange} />)

    await user.type(screen.getByLabelText('First number'), '12ab3')

    expect(screen.getByLabelText('First number')).toHaveValue('123')
    expect(onChange).toHaveBeenCalled()
    const calls = onChange.mock.calls.map((call) => call[0] as string)
    expect(calls.every((value) => /^[-.\d]*$/.test(value))).toBe(true)
  })

  it('shows a clear button that empties the value', async () => {
    const user = userEvent.setup()

    render(<ControlledOperandInput initial="42" />)

    expect(screen.getByLabelText('First number')).toHaveValue('42')
    await user.click(screen.getByRole('button', { name: 'Clear First number' }))
    expect(screen.getByLabelText('First number')).toHaveValue('')
  })

  it('renders a validation error', () => {
    render(<ControlledOperandInput error="required" />)
    expect(screen.getByRole('alert')).toHaveTextContent('required')
  })
})
