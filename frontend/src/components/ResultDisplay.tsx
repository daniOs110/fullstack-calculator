type ResultDisplayProps = {
  result: number | null
  error: string | null
}

export function ResultDisplay({ result, error }: ResultDisplayProps) {
  if (error) {
    return (
      <div className="result result-error" role="alert">
        <strong>Error</strong>
        <p>{error}</p>
      </div>
    )
  }

  if (result === null) {
    return (
      <div className="result result-empty" aria-live="polite">
        <p>Result will appear here</p>
      </div>
    )
  }

  return (
    <div className="result result-ok" aria-live="polite">
      <strong>Result</strong>
      <p>{formatResult(result)}</p>
    </div>
  )
}

function formatResult(value: number): string {
  // Avoid scientific notation for everyday values; keep precision for decimals.
  return Number.isInteger(value) ? String(value) : String(value)
}
