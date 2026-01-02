## 1. OpenAPI Specification
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

## 2. Error Handling
- [ ] 2.1 Define error code enum: INVALID_REQUEST, INVALID_FUND, FUND_NOT_FOUND, OWNER_NOT_FOUND, INSUFFICIENT_UNITS, SELF_TRANSFER, DUPLICATE_TRANSFER, INTERNAL_ERROR
- [ ] 2.2 Map error codes to HTTP status codes (400, 404, 409, 500)
- [ ] 2.3 Add error response examples for each endpoint

## 3. Code Generation Setup
- [ ] 3.1 Add oapi-codegen to Go module tools
- [ ] 3.2 Create Makefile target for Go server generation
- [ ] 3.3 Add openapi-typescript-codegen to frontend package.json
- [ ] 3.4 Create npm script for client generation

## 4. Validation
- [ ] 4.1 Validate spec with openapi-generator-cli validate
- [ ] 4.2 Generate Go types and verify compilation
- [ ] 4.3 Generate TypeScript client and verify compilation
- [ ] 4.4 Verify error types are generated correctly in TypeScript
