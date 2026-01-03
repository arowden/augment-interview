/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CapTable } from '../models/CapTable';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class CapTableService {
    /**
     * Get cap table for a fund
     * Returns the cap table entries for the specified fund with pagination support.
     * Each entry shows the owner, units held, percentage ownership, and acquisition date.
     *
     * @param fundId The unique identifier of the fund
     * @param limit Maximum number of entries to return
     * @param offset Number of entries to skip
     * @returns CapTable The cap table for the fund
     * @throws ApiError
     */
    public static getCapTable(
        fundId: string,
        limit: number = 100,
        offset?: number,
    ): CancelablePromise<CapTable> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/funds/{fundId}/cap-table',
            path: {
                'fundId': fundId,
            },
            query: {
                'limit': limit,
                'offset': offset,
            },
            errors: {
                400: `Invalid request parameters`,
                404: `Fund not found`,
                500: `Internal server error`,
            },
        });
    }
}
