// Package validation provides shared validation constraints derived from the OpenAPI spec.
//
// IMPORTANT: The OpenAPI specification (api/openapi.yaml) is the single source of truth
// for all validation constraints. These constants MUST match the spec exactly.
// When updating constraints, update the OpenAPI spec first, then update this file.
//
// See: api/openapi.yaml components/schemas/* for authoritative definitions.
package validation

// Name constraints from OpenAPI spec.
// See: api/openapi.yaml Fund.name, CreateFundRequest.name, CapTableEntry.ownerName
const (
	// MinNameLength is the minimum allowed length for names (minLength: 1).
	MinNameLength = 1
	// MaxNameLength is the maximum allowed length for names (maxLength: 255).
	MaxNameLength = 255
)

// Units constraints from OpenAPI spec.
// See: api/openapi.yaml Fund.totalUnits, Transfer.units
const (
	// MinUnits is the minimum allowed value for units (minimum: 1 for creation, 0 for entries).
	MinUnits = 1
	// MaxUnits is the maximum allowed value for units (maximum: 2147483647, PostgreSQL INTEGER max).
	MaxUnits = 2_147_483_647
)

// Pagination constraints from OpenAPI spec.
// See: api/openapi.yaml components/parameters/Limit, Offset
const (
	// DefaultLimit is the default pagination limit (default: 100).
	DefaultLimit = 100
	// MaxLimit is the maximum allowed pagination limit (maximum: 1000).
	MaxLimit = 1000
)

// Calculation constants.
const (
	// PercentageMultiplier converts a ratio to a percentage (0.5 * 100 = 50%).
	PercentageMultiplier = 100.0
)

// ListParams configures pagination for list operations.
// This is the shared pagination type used across all domain packages.
type ListParams struct {
	Limit  int
	Offset int
}

// Normalize applies defaults and constraints to ListParams.
// Returns a new ListParams with normalized values; the original is unchanged.
// Use: params = params.Normalize()
func (p ListParams) Normalize() ListParams {
	if p.Limit <= 0 {
		p.Limit = DefaultLimit
	}
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}
