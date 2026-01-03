## 1. OpenAPI Specification (Single Source of Truth)
- [ ] 1.1 Create `/api/openapi.yaml` with OpenAPI 3.0.3 header
- [ ] 1.2 Define Fund schema with validation constraints (minLength, maxLength, minimum, maximum)
- [ ] 1.3 Define CreateFundRequest with required fields and constraints
- [ ] 1.4 Define CapTable schema with pagination fields (entries, total, limit, offset)
- [ ] 1.5 Define CapTableEntry schema with percentage field
- [ ] 1.6 Define Transfer schema with timestamp examples
- [ ] 1.7 Define CreateTransferRequest with optional idempotencyKey (uuid)
- [ ] 1.8 Define Error schema with code enum and details object
- [ ] 1.9 Add Fund CRUD endpoints with error responses
- [ ] 1.10 Add CapTable endpoint with limit/offset query parameters
- [ ] 1.11 Add Transfer endpoints with idempotency support
- [ ] 1.12 Add operationIds for code generation
- [ ] 1.13 Add response examples for all error codes

## 2. Validation Constraints (OpenAPI as Single Source of Truth)
All validation rules MUST be defined in the OpenAPI spec. Generated code and implementations derive from this.

### 2.1 Field Constraints
- [ ] 2.1.1 Fund.name: minLength=1, maxLength=255, pattern=^\S.*\S$|^\S$ (no leading/trailing whitespace)
- [ ] 2.1.2 Fund.totalUnits: minimum=1, maximum=2147483647 (PostgreSQL INTEGER max)
- [ ] 2.1.3 CapTableEntry.units: minimum=0, maximum=2147483647
- [ ] 2.1.4 CapTableEntry.percentage: minimum=0, maximum=100
- [ ] 2.1.5 Transfer.units: minimum=1, maximum=2147483647
- [ ] 2.1.6 Pagination.limit: minimum=1, maximum=1000, default=100
- [ ] 2.1.7 Pagination.offset: minimum=0, default=0
- [ ] 2.1.8 Owner names: minLength=1, maxLength=255

### 2.2 Domain Invariants (x-invariants extension)
- [ ] 2.2.1 Document: "Sum of all cap table entry units equals fund.totalUnits"
- [ ] 2.2.2 Document: "Transfer.fromOwner != Transfer.toOwner"
- [ ] 2.2.3 Document: "Transfer.units <= fromOwner's current holdings"
- [ ] 2.2.4 Document: "All cap table percentages sum to 100.0 (within floating point tolerance)"

### 2.3 Format Constraints
- [ ] 2.3.1 All IDs: format=uuid
- [ ] 2.3.2 All timestamps: format=date-time (RFC 3339)
- [ ] 2.3.3 idempotencyKey: format=uuid

## 3. Error Handling
- [ ] 3.1 Define error code enum: INVALID_REQUEST, INVALID_FUND, FUND_NOT_FOUND, OWNER_NOT_FOUND, INSUFFICIENT_UNITS, SELF_TRANSFER, DUPLICATE_TRANSFER, INTERNAL_ERROR
- [ ] 3.2 Map error codes to HTTP status codes (400, 404, 409, 500)
- [ ] 3.3 Add error response examples for each endpoint
- [ ] 3.4 Document which validation constraints trigger which error codes

## 4. Code Generation Setup
- [ ] 4.1 Add oapi-codegen to Go module tools
- [ ] 4.2 Create Makefile target for Go server generation
- [ ] 4.3 Add openapi-typescript-codegen to frontend package.json
- [ ] 4.4 Create npm script for client generation

## 5. Validation & Verification
- [ ] 5.1 Validate spec with openapi-generator-cli validate
- [ ] 5.2 Generate Go types and verify compilation
- [ ] 5.3 Generate TypeScript client and verify compilation
- [ ] 5.4 Verify error types are generated correctly in TypeScript
- [ ] 5.5 Verify validation constraints are present in generated types
