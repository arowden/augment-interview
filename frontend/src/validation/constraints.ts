/**
 * Validation constraints derived from the OpenAPI specification.
 *
 * IMPORTANT: The OpenAPI spec (api/openapi.yaml) is the single source of truth
 * for all validation constraints. These constants MUST match the spec exactly.
 * When updating constraints, update the OpenAPI spec first, then update this file.
 *
 * @see api/openapi.yaml components/schemas/* for authoritative definitions
 */

// Name constraints from OpenAPI spec.
// See: api/openapi.yaml Fund.name, CreateFundRequest.name, CapTableEntry.ownerName
export const MIN_NAME_LENGTH = 1;
export const MAX_NAME_LENGTH = 255;

// Pattern: no leading/trailing whitespace (^\S(.*\S)?$)
export const NAME_PATTERN = /^\S(.*\S)?$/;

// Units constraints from OpenAPI spec.
// See: api/openapi.yaml Fund.totalUnits, Transfer.units
export const MIN_UNITS = 1;
export const MAX_UNITS = 2_147_483_647; // PostgreSQL INTEGER max

// Pagination constraints from OpenAPI spec.
// See: api/openapi.yaml components/parameters/Limit, Offset
export const DEFAULT_LIMIT = 100;
export const MAX_LIMIT = 1000;
