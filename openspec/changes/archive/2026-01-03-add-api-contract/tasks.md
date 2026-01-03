## 1. OpenAPI Specification (Single Source of Truth)
- [x] 1.1 Create `/api/openapi.yaml` with OpenAPI 3.0.3 header
- [x] 1.2 Define Fund schema with validation constraints (minLength, maxLength, minimum, maximum)
- [x] 1.3 Define CreateFundRequest with required fields and constraints
- [x] 1.4 Define CapTable schema with pagination fields (entries, total, limit, offset)
- [x] 1.5 Define CapTableEntry schema with percentage field
- [x] 1.6 Define Transfer schema with timestamp examples
- [x] 1.7 Define CreateTransferRequest with optional idempotencyKey (uuid)
- [x] 1.8 Define Error schema with code enum and details object
- [x] 1.9 Add Fund CRUD endpoints with error responses
- [x] 1.10 Add CapTable endpoint with limit/offset query parameters
- [x] 1.11 Add Transfer endpoints with idempotency support
- [x] 1.12 Add operationIds for code generation
- [x] 1.13 Add response examples for all error codes

## 2. Validation Constraints (OpenAPI as Single Source of Truth)
All validation rules MUST be defined in the OpenAPI spec. Generated code and implementations derive from this.

### 2.1 Field Constraints
- [x] 2.1.1 Fund.name: minLength=1, maxLength=255, pattern=^\S.*\S$|^\S$ (no leading/trailing whitespace)
- [x] 2.1.2 Fund.totalUnits: minimum=1, maximum=2147483647 (PostgreSQL INTEGER max)
- [x] 2.1.3 CapTableEntry.units: minimum=0, maximum=2147483647
- [x] 2.1.4 CapTableEntry.percentage: minimum=0, maximum=100
- [x] 2.1.5 Transfer.units: minimum=1, maximum=2147483647
- [x] 2.1.6 Pagination.limit: minimum=1, maximum=1000, default=100
- [x] 2.1.7 Pagination.offset: minimum=0, default=0
- [x] 2.1.8 Owner names: minLength=1, maxLength=255

### 2.2 Domain Invariants (x-invariants extension)
- [x] 2.2.1 Document: "Sum of all cap table entry units equals fund.totalUnits"
- [x] 2.2.2 Document: "Transfer.fromOwner != Transfer.toOwner"
- [x] 2.2.3 Document: "Transfer.units <= fromOwner's current holdings"
- [x] 2.2.4 Document: "All cap table percentages sum to 100.0 (within floating point tolerance)"

### 2.3 Format Constraints
- [x] 2.3.1 All IDs: format=uuid
- [x] 2.3.2 All timestamps: format=date-time (RFC 3339)
- [x] 2.3.3 idempotencyKey: format=uuid

## 3. Error Handling
- [x] 3.1 Define error code enum: INVALID_REQUEST, INVALID_FUND, FUND_NOT_FOUND, OWNER_NOT_FOUND, INSUFFICIENT_UNITS, SELF_TRANSFER, DUPLICATE_TRANSFER, INTERNAL_ERROR
- [x] 3.2 Map error codes to HTTP status codes (400, 404, 409, 500)
- [x] 3.3 Add error response examples for each endpoint
- [x] 3.4 Document which validation constraints trigger which error codes

## 4. Code Generation Setup
- [x] 4.1 Add oapi-codegen to Go module tools
- [x] 4.2 Create Makefile target for Go server generation
- [x] 4.3 Add openapi-typescript-codegen to frontend package.json
- [x] 4.4 Create npm script for client generation

## 5. Validation & Verification
- [x] 5.1 Validate spec with openapi-generator-cli validate
- [x] 5.2 Generate Go types and verify compilation
- [x] 5.3 Generate TypeScript client and verify compilation
- [x] 5.4 Verify error types are generated correctly in TypeScript
- [x] 5.5 Verify validation constraints are present in generated types
