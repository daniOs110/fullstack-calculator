# PROMPTS

This file summarizes the key prompts used during development and the accepted adjustments made in the final implementation.

## 1) Backend bootstrap and architecture

**Prompt used:**

> I've started the project. What steps should we follow next? I propose beginning with the backend structure — but before writing any code, explain the approach and wait for my approval.
>
> Follow-up decisions from the architecture proposal:
>
> - Use `chi` as the router (not plain `net/http` or Echo/Gin).
> - Prefer explicit handlers that call shared decode/validation helpers, rather than fully generic dynamic dispatch or duplicated validation.
> - Scope this first pass to folder structure, wiring, and a `/health` endpoint only.
>
> Later clarification on configuration: why does `config.go` hardcode `PORT` and `ALLOWED_ORIGINS`? Shouldn't those values come from my `.env`? Approved adding `godotenv` so a local `.env` is loaded when running with `go run`.

**Accepted decisions:**

- Use `cmd/server` + `internal/*` package layout
- Use `chi` router with middleware stack
- Keep handlers thin and route arithmetic through pure logic functions
- Add `/health` endpoint in initial skeleton

**Adjusted details:**

- Added `godotenv` so local `.env` is loaded when running with `go run`
- Kept environment variables as source of truth; `.env` is fallback for local only

## 2) Calculator logic and numeric boundaries

**Prompt used:**

> Proceed with option 1: implement `internal/calculator` and, in the same step, share the planned test-case list for my review before writing tests.
>
> On non-finite results (`NaN` / `Inf`): account for the language and `float64` limitations. Bound the magnitude of incoming operands so those cases are less likely, and reject invalid input early with an error the frontend can handle (with double validation later on the client). Confirm `/percentage` as `a% of b`. Approved `±1e15` as the operand limit.

**Accepted decisions:**

- Pure functions in `internal/calculator`
- Explicit errors for division by zero and negative square root
- Percentage semantics: `a% of b` (`a / 100 * b`)

**Adjusted details:**

- Added domain bound `MaxOperand = 1e15`
- Preserved behavior for IEEE-754 non-finite outputs and rejected them in HTTP layer

## 3) HTTP API handlers and validation strategy

**Prompt used:**

> Approved — proceed with the handlers. Also: negative zero should be normalized in API responses (return `0`, not `-0`).
>
> Clarification while reviewing the HTTP layer: in Java we often keep request/response DTOs in a shared `common` package. Is that possible and idiomatic in Go? And are these handlers the equivalent of Spring controllers — where is the HTTP method (GET/POST) defined?

**Accepted decisions:**

- Operation-specific request DTOs
- Strict decoding with unknown field rejection
- Shared decode/validation helpers to avoid boilerplate duplication

**Adjusted details:**

- Enforced exact request shapes per endpoint
- Normalized negative zero in responses (return `0` instead of `-0`)

## 4) Backend testing workflow

**Prompt used:**

> Remind me of the testing process from the project rules. After the calculator logic is in place: the proposed test cases look good — please add them.
>
> After the handlers landed: please restate the `httpapi` test-case list (happy paths, domain errors, input validation, distinct body shapes, non-finite operand checks, negative-zero normalization, and router cases). Approved — implement those tests.

**Accepted decisions:**

- Case-first workflow (approval before writing tests)
- Table-driven tests for calculator logic
- Handler tests using `httptest`

**Adjusted details:**

- Added non-finite edge tests for exponent behavior
- Added route/method and strict payload validation tests

## 5) Frontend architecture and UX

**Prompt used:**

> Let's move on to the frontend. Explain the approach and wait for approval before writing code. For the first pass, start with a small skeleton only (folders, stubs, `.env.example`). Show concrete examples of the three UI styles you proposed (form-based, calculator-keypad style, and hybrid) — the written proposal alone was hard to visualize.
>
> Approved the form-based UI (option A). One hard requirement: when switching operations (e.g. to Square root), field B must not only be hidden visually — it must not be included in the request body, so we do not break the backend's strict unknown-field validation.
>
> UX follow-ups after trying the UI: replace generic A/B labels with operation-specific English labels; filter input so only numbers can be typed (option A); cap input at 15 integer digits and 10 decimals; add a clear (×) button on each input.

We discussed three UI approaches before implementing: (1) a form with an operation selector and dynamic fields, (2) a classic calculator keypad, and (3) a hybrid of both. We chose the form-based approach because it maps cleanly to all seven API operations (including sqrt, exponent, and percentage), stays easy to test, and fits the short delivery budget better than keypad state machines.

**Accepted decisions:**

- Form-based UI (operation selector + dynamic fields)
- Strong client-side validation plus backend validation
- Build API payload with only active operation fields

**Adjusted details:**

- Operation-specific labels (clearer than generic A/B)
- Numeric-only typing sanitization
- Input constraints: max 15 integer digits and 10 decimals
- Per-input clear button (x)

## 6) Frontend testing strategy

**Prompt used:**

> Please prepare the frontend test-case list for my review before implementing any tests.
>
> Those cases look good — please develop them.

**Accepted decisions:**

- Use Vitest + React Testing Library + jsdom
- Cover sanitize/parse helpers, body shaping, component rendering, and submit flows

**Adjusted details:**

- Added test setup cleanup hook to avoid DOM contamination between test cases
- Added coverage script and validation of endpoint-specific payloads (e.g., sqrt sends only `{a}`)

## 7) Docker Compose deployment

**Prompt used:**

> Approved option 1 (nginx reverse proxy). Frontend on port 80, backend on port 8080.
>
> Follow-up after first Docker use: explain how to start Docker Desktop and verify the stack, and why `client.ts` / `vite-env.d.ts` changed for empty `VITE_API_BASE_URL` in Compose.

**Accepted decisions:**

- Multi-stage Dockerfiles for backend and frontend
- nginx serves frontend on port 80
- nginx proxies `/api/*` to backend service

**Adjusted details:**

- Frontend API client supports empty `VITE_API_BASE_URL` for same-origin proxy mode
- Compose ports fixed to frontend `80` and backend `8080`

## 8) Documentation policy followed

**Prompt used:**

> That works — please draft the full `README.md` and `PROMPTS.md`. Keep both documents in English.
>
> Follow-up on `PROMPTS.md`: replace each "Prompt intent" summary with a "Prompt used" block containing the actual prompts from our conversation (cleaned up for professional tone, without changing the technical substance). Keep "Accepted decisions" and "Adjusted details". In the frontend section, briefly note that we compared three UI approaches and chose the form-based one.

**Accepted decisions:**

- README includes local + Docker run flows
- README includes backend/frontend testing and coverage commands
- README includes cURL examples for all required endpoints

**Security rules respected throughout:**

- Real `.env` files were never read or included in prompts
- `.env.example` files used for documented configuration

