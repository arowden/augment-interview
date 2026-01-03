/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * Structured error response
 */
export type Error = {
    /**
     * Machine-readable error code
     */
    code: 'INVALID_REQUEST' | 'INVALID_FUND' | 'FUND_NOT_FOUND' | 'OWNER_NOT_FOUND' | 'INSUFFICIENT_UNITS' | 'SELF_TRANSFER' | 'DUPLICATE_TRANSFER' | 'INTERNAL_ERROR';
    /**
     * Human-readable error message
     */
    message: string;
    /**
     * Additional context about the error
     */
    details?: Record<string, any>;
};

