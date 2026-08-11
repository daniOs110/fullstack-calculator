import type {
  ErrorResponse,
  Operation,
  ResultResponse,
} from './types'

export class ApiError extends Error {
  readonly status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

/**
 * Returns the API origin. An empty string means same-origin requests, which
 * is what Docker Compose uses: nginx proxies /api/* to the backend service.
 */
function apiBaseUrl(): string {
  const base = import.meta.env.VITE_API_BASE_URL
  if (base === undefined || base === null) {
    return ''
  }
  if (typeof base !== 'string') {
    throw new Error('VITE_API_BASE_URL must be a string')
  }
  return base.trim().replace(/\/$/, '')
}

function apiUrl(path: string): string {
  const normalized = path.startsWith('/') ? path : `/${path}`
  const base = apiBaseUrl()
  return base === '' ? normalized : `${base}${normalized}`
}

/**
 * Calls POST /api/v1/{operation} with a body that must contain only the fields
 * that operation expects — never extra keys (backend rejects unknown fields).
 */
export async function calculate(
  operation: Operation,
  body: Record<string, number>,
): Promise<number> {
  const response = await fetch(apiUrl(`/api/v1/${operation}`), {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
    body: JSON.stringify(body),
  })

  const payload: unknown = await response.json().catch(() => null)

  if (!response.ok) {
    const message =
      isErrorResponse(payload) && payload.error
        ? payload.error
        : `request failed with status ${response.status}`
    throw new ApiError(message, response.status)
  }

  if (!isResultResponse(payload)) {
    throw new ApiError('unexpected response from server', response.status)
  }

  return payload.result
}

function isErrorResponse(value: unknown): value is ErrorResponse {
  return (
    typeof value === 'object' &&
    value !== null &&
    'error' in value &&
    typeof (value as ErrorResponse).error === 'string'
  )
}

function isResultResponse(value: unknown): value is ResultResponse {
  return (
    typeof value === 'object' &&
    value !== null &&
    'result' in value &&
    typeof (value as ResultResponse).result === 'number'
  )
}
