/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * Request body for creating a new transfer
 */
export type CreateTransferRequest = {
    /**
     * Client-generated key for deduplication. If provided, subsequent requests with the same key will return the original transfer.
     */
    idempotencyKey?: string;
    /**
     * Name of the sender (must exist in cap table, no leading/trailing whitespace)
     */
    fromOwner: string;
    /**
     * Name of the recipient (no leading/trailing whitespace)
     */
    toOwner: string;
    /**
     * Number of units to transfer (must not exceed sender's holdings)
     */
    units: number;
};

