## Context
The API contract serves as the single source of truth for all HTTP communication between frontend and backend. It must support code generation for both Go server stubs and TypeScript client.

## Goals / Non-Goals
- Goals: Type-safe API contract, automated code generation, consistent error handling, idempotent transfers
- Non-Goals: API versioning, authentication schemes, rate limiting

## Decisions
- Decision: Use OpenAPI 3.0.3 for broad tooling support
- Alternatives considered: OpenAPI 3.1 (newer but less tooling), GraphQL (overkill for CRUD), gRPC (no browser support)

- Decision: Use oapi-codegen for Go (generates Chi-compatible interfaces)
- Alternatives considered: go-swagger (heavier), ogen (less mature)

- Decision: Use openapi-typescript-codegen for frontend
- Alternatives considered: openapi-generator (Java dependency), orval (less stable)

- Decision: Client-generated idempotency keys for transfer deduplication
- Alternatives considered: Server-generated transaction IDs (client can't retry safely), database unique constraints only (no client feedback)

## API Design Choices
- All endpoints under `/api/` prefix for clarity
- UUIDs for all resource identifiers
- ISO 8601 timestamps with timezone (example: "2024-01-15T10:30:00Z")
- Nested resources: `/funds/{fundId}/cap-table`, `/funds/{fundId}/transfers`
- Standard HTTP status codes: 200, 201, 400, 404, 409, 500
- Pagination on cap-table endpoint (default limit: 100, max: 1000)

## Error Handling Design
All error responses use structured Error schema:
```yaml
Error:
  type: object
  required: [code, message]
  properties:
    code:
      type: string
      enum:
        - INVALID_REQUEST
        - INVALID_FUND
        - FUND_NOT_FOUND
        - OWNER_NOT_FOUND
        - INSUFFICIENT_UNITS
        - SELF_TRANSFER
        - DUPLICATE_TRANSFER
        - INTERNAL_ERROR
    message:
      type: string
      description: Human-readable error message
    details:
      type: object
      additionalProperties: true
      description: Additional context (field errors, IDs, etc.)
```

Error code mapping:
- 400 Bad Request: INVALID_REQUEST, INVALID_FUND, INSUFFICIENT_UNITS, SELF_TRANSFER
- 404 Not Found: FUND_NOT_FOUND, OWNER_NOT_FOUND
- 409 Conflict: DUPLICATE_TRANSFER (idempotency key reused with different data)
- 500 Internal: INTERNAL_ERROR

## Schema Validation Constraints
```yaml
CreateFundRequest:
  type: object
  required: [name, totalUnits, initialOwner]
  properties:
    name:
      type: string
      minLength: 1
      maxLength: 255
    totalUnits:
      type: integer
      minimum: 1
      maximum: 2147483647
    initialOwner:
      type: string
      minLength: 1
      maxLength: 255

CreateTransferRequest:
  type: object
  required: [fromOwner, toOwner, units]
  properties:
    idempotencyKey:
      type: string
      format: uuid
      description: Client-generated key for deduplication
    fromOwner:
      type: string
      minLength: 1
      maxLength: 255
    toOwner:
      type: string
      minLength: 1
      maxLength: 255
    units:
      type: integer
      minimum: 1
      maximum: 2147483647
```

## Pagination Design
```yaml
CapTable:
  type: object
  properties:
    fundId:
      type: string
      format: uuid
    entries:
      type: array
      items:
        $ref: '#/components/schemas/CapTableEntry'
    total:
      type: integer
      description: Total number of entries
    limit:
      type: integer
    offset:
      type: integer

# Query parameters
parameters:
  - name: limit
    in: query
    schema:
      type: integer
      default: 100
      minimum: 1
      maximum: 1000
  - name: offset
    in: query
    schema:
      type: integer
      default: 0
      minimum: 0
```

## Generated Server Interface
```go
type ServerInterface interface {
    GetFunds(w http.ResponseWriter, r *http.Request)
    CreateFund(w http.ResponseWriter, r *http.Request)
    GetFund(w http.ResponseWriter, r *http.Request, fundId openapi_types.UUID)
    GetCapTable(w http.ResponseWriter, r *http.Request, fundId openapi_types.UUID, params GetCapTableParams)
    GetTransfers(w http.ResponseWriter, r *http.Request, fundId openapi_types.UUID)
    CreateTransfer(w http.ResponseWriter, r *http.Request, fundId openapi_types.UUID)
}

type GetCapTableParams struct {
    Limit  *int `form:"limit" json:"limit"`
    Offset *int `form:"offset" json:"offset"`
}
```

## Risks / Trade-offs
- Generated code requires regeneration on spec changes → Mitigated by Makefile targets
- Generated code style may not match preferences → Accept generated patterns
- Idempotency key storage required → Store in transfers table with unique constraint

## Open Questions
- None
