/**
 * Validation constraints derived from the OpenAPI specification.
 *
 * IMPORTANT: The OpenAPI spec (api/openapi.yaml) is the single source of truth
 * for all validation constraints. These constants MUST match the spec exactly.
 * When updating constraints, update the OpenAPI spec first, then update this file.
 *
 * @see api/openapi.yaml components/schemas/* for authoritative definitions
 */

export const MIN_NAME_LENGTH = 1;
export const MAX_NAME_LENGTH = 255;

export const NAME_PATTERN = /^\S(.*\S)?$/;

export const MIN_UNITS = 1;
export const MAX_UNITS = 2_147_483_647;

export const DEFAULT_LIMIT = 100;
export const MAX_LIMIT = 1000;
