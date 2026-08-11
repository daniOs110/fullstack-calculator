import type { Operation } from './types'

/**
 * HTTP client for /api/v1/* endpoints.
 * Implementation comes in the next frontend iteration.
 */
export async function calculate(
  _operation: Operation,
  _body: Record<string, number>,
): Promise<number> {
  throw new Error('api client not implemented yet')
}
