/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
/**
 * A record of units transferred between owners
 */
export type Transfer = {
    /**
     * Unique identifier for the transfer
     */
    id: string;
    /**
     * The fund the transfer belongs to
     */
    fundId: string;
    /**
     * Name of the sender (no leading/trailing whitespace)
     */
    fromOwner: string;
    /**
     * Name of the recipient (no leading/trailing whitespace)
     */
    toOwner: string;
    /**
     * Number of units transferred
     */
    units: number;
    /**
     * Timestamp when the transfer was executed
     */
    transferredAt: string;
};

