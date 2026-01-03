# api-contract Specification

## Purpose
TBD - created by archiving change add-api-contract. Update Purpose after archive.
## Requirements
### Requirement: OpenAPI Specification File
The system SHALL provide an OpenAPI 3.0.3 specification file at `/api/openapi.yaml` that defines all HTTP endpoints, request bodies, response schemas, and error formats.

#### Scenario: Spec file exists and is valid
- **WHEN** the spec file is validated with an OpenAPI linter
- **THEN** no errors are reported

#### Scenario: Spec defines all endpoints
- **WHEN** the spec is parsed
- **THEN** it contains endpoints for fund CRUD, cap-table read, and transfer operations

### Requirement: Fund Endpoints
The API SHALL expose endpoints for creating, listing, and retrieving funds.

#### Scenario: Create fund endpoint
- **WHEN** POST /api/funds is called with name, totalUnits, and initialOwner
- **THEN** a 201 response with the created Fund is returned

#### Scenario: List funds endpoint
- **WHEN** GET /api/funds is called
- **THEN** a 200 response with an array of Fund objects is returned

#### Scenario: Get fund endpoint
- **WHEN** GET /api/funds/{fundId} is called with a valid UUID
- **THEN** a 200 response with the Fund object is returned

#### Scenario: Get fund not found
- **WHEN** GET /api/funds/{fundId} is called with a non-existent UUID
- **THEN** a 404 response with Error object containing code FUND_NOT_FOUND is returned

### Requirement: Cap Table Endpoint
The API SHALL expose an endpoint for retrieving a fund's cap table with pagination support.

#### Scenario: Get cap table
- **WHEN** GET /api/funds/{fundId}/cap-table is called
- **THEN** a 200 response with fundId, entries array, total count, limit, and offset is returned

#### Scenario: Get cap table with pagination
- **WHEN** GET /api/funds/{fundId}/cap-table?limit=10&offset=20 is called
- **THEN** a 200 response with at most 10 entries starting from offset 20 is returned

#### Scenario: Get cap table for non-existent fund
- **WHEN** GET /api/funds/{fundId}/cap-table is called with a non-existent fundId
- **THEN** a 404 response with Error object containing code FUND_NOT_FOUND is returned

### Requirement: Transfer Endpoints
The API SHALL expose endpoints for creating and listing transfers with idempotency support.

#### Scenario: Create transfer endpoint
- **WHEN** POST /api/funds/{fundId}/transfers is called with fromOwner, toOwner, units, and optional idempotencyKey
- **THEN** a 201 response with the created Transfer is returned

#### Scenario: Create transfer with idempotency key
- **WHEN** POST /api/funds/{fundId}/transfers is called twice with the same idempotencyKey
- **THEN** the second call returns 200 with the original Transfer (not 201)

#### Scenario: List transfers endpoint
- **WHEN** GET /api/funds/{fundId}/transfers is called
- **THEN** a 200 response with an array of Transfer objects is returned

#### Scenario: Transfer validation error
- **WHEN** POST /api/funds/{fundId}/transfers is called with invalid data
- **THEN** a 400 response with Error object containing appropriate code is returned

### Requirement: Schema Definitions with Validation
The API spec SHALL define schemas with validation constraints for Fund, CreateFundRequest, CapTable, CapTableEntry, Transfer, CreateTransferRequest, and Error.

#### Scenario: Fund schema
- **WHEN** the Fund schema is examined
- **THEN** it contains id (uuid), name (string 1-255 chars), totalUnits (integer min 1), and createdAt (date-time with example)

#### Scenario: CreateFundRequest schema
- **WHEN** the CreateFundRequest schema is examined
- **THEN** it contains required fields name (string 1-255 chars), totalUnits (integer 1-2147483647), and initialOwner (string 1-255 chars)

#### Scenario: CapTableEntry schema
- **WHEN** the CapTableEntry schema is examined
- **THEN** it contains ownerName (string), units (integer min 0), percentage (number), and acquiredAt (date-time)

#### Scenario: CapTable schema with pagination
- **WHEN** the CapTable schema is examined
- **THEN** it contains fundId (uuid), entries array, total (integer), limit (integer), and offset (integer)

#### Scenario: Transfer schema
- **WHEN** the Transfer schema is examined
- **THEN** it contains id (uuid), fundId (uuid), fromOwner, toOwner, units, and transferredAt (date-time with example)

#### Scenario: CreateTransferRequest schema with idempotency
- **WHEN** the CreateTransferRequest schema is examined
- **THEN** it contains required fields fromOwner, toOwner, units (integer min 1), and optional idempotencyKey (uuid)

### Requirement: Error Response Schema
The API SHALL define a structured Error schema with machine-readable error codes.

#### Scenario: Error schema structure
- **WHEN** the Error schema is examined
- **THEN** it contains required code (enum), message (string), and optional details (object)

#### Scenario: Error codes enum
- **WHEN** the Error code enum is examined
- **THEN** it includes INVALID_REQUEST, INVALID_FUND, FUND_NOT_FOUND, OWNER_NOT_FOUND, INSUFFICIENT_UNITS, SELF_TRANSFER, DUPLICATE_TRANSFER, INTERNAL_ERROR

#### Scenario: Error response examples
- **WHEN** 400/404/409/500 responses are examined
- **THEN** each includes example Error objects with appropriate codes

### Requirement: Go Server Code Generation
The system SHALL support generating Go server interfaces and types from the OpenAPI spec using oapi-codegen.

#### Scenario: Generate Go code
- **WHEN** make generate-api is run
- **THEN** internal/http/openapi.gen.go is created with ServerInterface and type definitions

#### Scenario: Generated code compiles
- **WHEN** the generated Go code is built
- **THEN** no compilation errors occur

#### Scenario: Generated interface methods
- **WHEN** ServerInterface is examined
- **THEN** it contains methods for all endpoints with typed parameters and responses

### Requirement: TypeScript Client Code Generation
The system SHALL support generating a TypeScript API client from the OpenAPI spec.

#### Scenario: Generate TypeScript client
- **WHEN** npm run generate-api is run in the frontend directory
- **THEN** src/api/generated/ directory is populated with typed client code

#### Scenario: Generated client compiles
- **WHEN** the frontend TypeScript is compiled
- **THEN** no type errors occur in the generated client

#### Scenario: Generated error types
- **WHEN** the generated TypeScript is examined
- **THEN** Error type includes the code enum for exhaustive error handling

