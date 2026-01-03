/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * A single entry in the cap table representing an owner's stake
 */
export type CapTableEntry = {
    /**
     * Name of the unit holder (no leading/trailing whitespace)
     */
    ownerName: string;
    /**
     * Number of units owned
     */
    units: number;
    /**
     * Percentage of total fund units owned
     */
    percentage: number;
    /**
     * Timestamp when the ownership was first acquired
     */
    acquiredAt: string;
};

