import { ApiError } from './generated/core/ApiError';
import type { Error as ApiErrorBody } from './generated/models/Error';

/**
 * Parsed API error with structured information.
 */
export interface ParsedApiError {
  code: ApiErrorBody['code'] | 'UNKNOWN';
  message: string;
  requestId?: string;
  details?: Record<string, unknown>;
  status: number;
}

/**
 * User-friendly messages for each error code.
 */
const ERROR_MESSAGES: Record<ApiErrorBody['code'], string> = {
  INVALID_REQUEST: 'The request was invalid. Please check your input.',
  INVALID_FUND: 'Invalid fund data. Name must be non-empty and units must be positive.',
  FUND_NOT_FOUND: 'The requested fund could not be found.',
  OWNER_NOT_FOUND: 'The specified owner does not exist in this fund.',
  INSUFFICIENT_UNITS: 'The sender does not have enough units for this transfer.',
  SELF_TRANSFER: 'Cannot transfer units to the same owner.',
  DUPLICATE_TRANSFER: 'This transfer has already been processed.',
  INTERNAL_ERROR: 'An unexpected error occurred. Please try again.',
};

/**
 * Parses an API error into a structured format.
 * Works with both ApiError from the generated client and generic errors.
 */
export function parseApiError(error: unknown): ParsedApiError {
  // Handle generated ApiError.
  if (error instanceof ApiError) {
    const body = error.body as ApiErrorBody | undefined;

    if (body && typeof body === 'object' && 'code' in body) {
      return {
        code: body.code,
        message: body.message || ERROR_MESSAGES[body.code] || error.message,
        requestId: body.details?.requestId as string | undefined,
        details: body.details,
        status: error.status,
      };
    }

    // ApiError without structured body.
    return {
      code: 'UNKNOWN',
      message: error.message || 'An unexpected error occurred',
      status: error.status,
    };
  }

  // Handle generic Error.
  if (error instanceof Error) {
    return {
      code: 'UNKNOWN',
      message: error.message || 'An unexpected error occurred',
      status: 0,
    };
  }

  // Handle unknown error types.
  return {
    code: 'UNKNOWN',
    message: 'An unexpected error occurred',
    status: 0,
  };
}

/**
 * Gets a user-friendly error message for an error code.
 */
export function getErrorMessage(code: ApiErrorBody['code']): string {
  return ERROR_MESSAGES[code] || ERROR_MESSAGES.INTERNAL_ERROR;
}

/**
 * Checks if an error is a specific API error code.
 */
export function isApiErrorCode(
  error: unknown,
  code: ApiErrorBody['code']
): boolean {
  return parseApiError(error).code === code;
}
