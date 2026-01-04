/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Transfer } from './Transfer';
/**
 * Paginated list of transfers for a fund
 */
export type TransferList = {
    /**
     * The fund these transfers belong to
     */
    fundId: string;
    /**
     * Transfers for the current page
     */
    transfers: Array<Transfer>;
    /**
     * Total number of transfers
     */
    total: number;
    /**
     * Maximum transfers per page
     */
    limit: number;
    /**
     * Number of transfers skipped
     */
    offset: number;
};

