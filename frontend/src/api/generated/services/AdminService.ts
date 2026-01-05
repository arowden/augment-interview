/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AdminService {
    /**
     * Reset all data
     * Deletes all funds, cap table entries, and transfers from the database.
     * This is a destructive operation intended for development/testing purposes only.
     *
     * @returns any Database reset successfully
     * @throws ApiError
     */
    public static resetDatabase(): CancelablePromise<{
        message?: string;
        deletedFunds?: number;
        deletedTransfers?: number;
        deletedOwnership?: number;
    }> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/reset',
            errors: {
                500: `Internal server error`,
            },
        });
    }
}
