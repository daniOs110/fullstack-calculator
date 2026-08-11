import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { Calculator } from './Calculator'

function mockJsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response
}

describe('Calculator', () => {
  beforeEach(() => {
    vi.stubEnv('VITE_API_BASE_URL', 'http://localhost:8080')
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllEnvs()
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('shows First number and Second number for Add', () => {
    render(<Calculator />)
    expect(screen.getByLabelText('First number')).toBeInTheDocument()
    expect(screen.getByLabelText('Second number')).toBeInTheDocument()
  })

  it('shows only Number for Square root', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText('Operation'), 'sqrt')

    expect(screen.getByLabelText('Number')).toBeInTheDocument()
    expect(screen.queryByLabelText('Second number')).not.toBeInTheDocument()
  })

  it('shows Base and Exponent for Exponent', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText('Operation'), 'exponent')

    expect(screen.getByLabelText('Base')).toBeInTheDocument()
    expect(screen.getByLabelText('Exponent')).toBeInTheDocument()
  })

  it('shows Percentage labels for Percentage', async () => {
    const user = userEvent.setup()
    render(<Calculator />)

    await user.selectOptions(screen.getByLabelText('Operation'), 'percentage')

    expect(screen.getByLabelText('Percentage (%)')).toBeInTheDocument()
    expect(screen.getByLabelText('Of (number)')).toBeInTheDocument()
  })

  it('shows required errors on empty submit and does not call fetch', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    render(<Calculator />)
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    const alerts = screen.getAllByRole('alert')
    expect(alerts.length).toBeGreaterThanOrEqual(2)
    expect(alerts[0]).toHaveTextContent('required')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('blocks incomplete drafts without calling fetch', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    render(<Calculator />)
    await user.type(screen.getByLabelText('First number'), '-')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    expect(screen.getByRole('alert')).toHaveTextContent(
      'must be a finite number',
    )
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('posts add with { a, b } and renders the result', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch).mockResolvedValue(
      mockJsonResponse(200, { result: 5 }),
    )

    render(<Calculator />)
    await user.type(screen.getByLabelText('First number'), '2')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    await waitFor(() => {
      expect(screen.getByText('5')).toBeInTheDocument()
    })

    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [url, init] = fetchMock.mock.calls[0]!
    expect(url).toBe('http://localhost:8080/api/v1/add')
    expect(init).toMatchObject({ method: 'POST' })
    expect(JSON.parse(String(init?.body))).toEqual({ a: 2, b: 3 })
  })

  it('posts sqrt with only { a }', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch).mockResolvedValue(
      mockJsonResponse(200, { result: 3 }),
    )

    render(<Calculator />)
    await user.type(screen.getByLabelText('First number'), '9')
    await user.type(screen.getByLabelText('Second number'), '99')
    await user.selectOptions(screen.getByLabelText('Operation'), 'sqrt')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    await waitFor(() => {
      expect(screen.getByText('3')).toBeInTheDocument()
    })

    const body = JSON.parse(String(fetchMock.mock.calls[0]![1]?.body))
    expect(body).toEqual({ a: 9 })
    expect(body).not.toHaveProperty('b')
  })

  it('renders a backend error message', async () => {
    const user = userEvent.setup()
    vi.mocked(fetch).mockResolvedValue(
      mockJsonResponse(400, { error: 'division by zero' }),
    )

    render(<Calculator />)
    await user.selectOptions(screen.getByLabelText('Operation'), 'divide')
    await user.type(screen.getByLabelText('First number'), '10')
    await user.type(screen.getByLabelText('Second number'), '0')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    await waitFor(() => {
      expect(screen.getByRole('alert')).toHaveTextContent('division by zero')
    })
  })

  it('clears previous result when the operation changes', async () => {
    const user = userEvent.setup()
    vi.mocked(fetch).mockResolvedValue(mockJsonResponse(200, { result: 5 }))

    render(<Calculator />)
    await user.type(screen.getByLabelText('First number'), '2')
    await user.type(screen.getByLabelText('Second number'), '3')
    await user.click(screen.getByRole('button', { name: 'Calculate' }))

    await waitFor(() => {
      expect(screen.getByText('5')).toBeInTheDocument()
    })

    await user.selectOptions(screen.getByLabelText('Operation'), 'subtract')
    expect(screen.getByText('Result will appear here')).toBeInTheDocument()
    expect(screen.queryByText('5')).not.toBeInTheDocument()
  })
})
