/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Fund } from './Fund';
/**
 * Paginated list of funds
 */
export type FundList = {
    /**
     * Funds for the current page
     */
    funds: Array<Fund>;
    /**
     * Total number of funds
     */
    total: number;
    /**
     * Maximum funds per page
     */
    limit: number;
    /**
     * Number of funds skipped
     */
    offset: number;
};

