/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateFundRequest } from '../models/CreateFundRequest';
import type { Fund } from '../models/Fund';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class FundsService {
    /**
     * List all funds
     * Returns a list of all funds in the system.
     * @returns Fund A list of funds
     * @throws ApiError
     */
    public static listFunds(): CancelablePromise<Array<Fund>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/funds',
            errors: {
                500: `Internal server error`,
            },
        });
    }
    /**
     * Create a new fund
     * Creates a new fund with the specified name, total units, and initial owner.
     * The initial owner receives all units in the cap table.
     *
     * @param requestBody
     * @returns Fund Fund created successfully
     * @throws ApiError
     */
    public static createFund(
        requestBody: CreateFundRequest,
    ): CancelablePromise<Fund> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/funds',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Invalid request parameters`,
                500: `Internal server error`,
            },
        });
    }
    /**
     * Get a fund by ID
     * Returns the fund with the specified ID.
     * @param fundId The unique identifier of the fund
     * @returns Fund The requested fund
     * @throws ApiError
     */
    public static getFund(
        fundId: string,
    ): CancelablePromise<Fund> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/funds/{fundId}',
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
}
