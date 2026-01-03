/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * An investment fund with ownership units
 */
export type Fund = {
    /**
     * Unique identifier for the fund
     */
    id: string;
    /**
     * Name of the fund (no leading/trailing whitespace)
     */
    name: string;
    /**
     * Total number of ownership units in the fund
     */
    totalUnits: number;
    /**
     * Timestamp when the fund was created
     */
    createdAt: string;
};

