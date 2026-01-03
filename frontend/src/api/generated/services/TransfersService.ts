/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateTransferRequest } from '../models/CreateTransferRequest';
import type { Transfer } from '../models/Transfer';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class TransfersService {
    /**
     * List transfers for a fund
     * Returns all transfers for the specified fund, ordered by transfer date descending.
     * @param fundId The unique identifier of the fund
     * @returns Transfer A list of transfers
     * @throws ApiError
     */
    public static listTransfers(
        fundId: string,
    ): CancelablePromise<Array<Transfer>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/funds/{fundId}/transfers',
            path: {
                'fundId': fundId,
            },
            errors: {
                400: `Invalid request parameters`,
                404: `Fund not found`,
                500: `Internal server error`,
            },
        });
    }
    /**
     * Create a transfer
     * Creates a new transfer of units between owners within a fund.
     *
     * ## Idempotency
     * Use the optional `idempotencyKey` field to safely retry transfer requests.
     * - First request with a key: Creates the transfer and returns 201
     * - Subsequent requests with the same key and data: Returns the original transfer with 200
     * - Subsequent requests with the same key but different data: Returns 409 Conflict
     *
     * ## Validation
     * - `fromOwner` must exist in the cap table with sufficient units
     * - `toOwner` must be different from `fromOwner`
     * - `units` must be positive and not exceed the sender's holdings
     *
     * @param fundId The unique identifier of the fund
     * @param requestBody
     * @returns Transfer Idempotent request - returning existing transfer
     * @throws ApiError
     */
    public static createTransfer(
        fundId: string,
        requestBody: CreateTransferRequest,
    ): CancelablePromise<Transfer> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/funds/{fundId}/transfers',
            path: {
                'fundId': fundId,
            },
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Invalid transfer request`,
                404: `Owner not found in cap table`,
                409: `Idempotency key already used with different data`,
                500: `Internal server error`,
            },
        });
    }
}
