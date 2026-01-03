/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CapTableEntry } from './CapTableEntry';
/**
 * Paginated cap table for a fund
 */
export type CapTable = {
    /**
     * The fund this cap table belongs to
     */
    fundId: string;
    /**
     * Cap table entries for the current page
     */
    entries: Array<CapTableEntry>;
    /**
     * Total number of entries in the cap table
     */
    total: number;
    /**
     * Maximum entries per page
     */
    limit: number;
    /**
     * Number of entries skipped
     */
    offset: number;
};

