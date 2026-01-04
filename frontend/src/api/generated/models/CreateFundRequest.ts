/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * Request body for creating a new fund
 */
export type CreateFundRequest = {
    /**
     * Name of the fund (no leading/trailing whitespace)
     */
    name: string;
    /**
     * Total number of ownership units in the fund
     */
    totalUnits: number;
    /**
     * Name of the initial owner who will receive all units (no leading/trailing whitespace)
     */
    initialOwner: string;
};

