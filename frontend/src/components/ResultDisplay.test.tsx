import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ResultDisplay } from './ResultDisplay'

describe('ResultDisplay', () => {
  it('shows the empty placeholder', () => {
    render(<ResultDisplay result={null} error={null} />)
    expect(screen.getByText('Result will appear here')).toBeInTheDocument()
  })

  it('shows a numeric result', () => {
    render(<ResultDisplay result={5} error={null} />)
    expect(screen.getByText('Result')).toBeInTheDocument()
    expect(screen.getByText('5')).toBeInTheDocument()
  })

  it('shows an error and prefers it over a result', () => {
    render(<ResultDisplay result={5} error="division by zero" />)
    const alert = screen.getByRole('alert')
    expect(alert).toHaveTextContent('division by zero')
    expect(alert).not.toHaveTextContent(/^Result$/)
    expect(screen.queryByText('Result', { selector: 'strong' })).not.toBeInTheDocument()
  })
})
