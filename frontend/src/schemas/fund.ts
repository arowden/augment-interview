import { z } from 'zod';
import {
  MIN_NAME_LENGTH,
  MAX_NAME_LENGTH,
  NAME_PATTERN,
  MIN_UNITS,
  MAX_UNITS,
} from '../validation/constraints';

export const createFundSchema = z.object({
  name: z
    .string()
    .min(MIN_NAME_LENGTH, 'Name is required')
    .max(MAX_NAME_LENGTH, `Name must be at most ${MAX_NAME_LENGTH} characters`)
    .regex(NAME_PATTERN, 'Name cannot have leading or trailing whitespace'),
  totalUnits: z
    .number({ invalid_type_error: 'Total units must be a number' })
    .int('Total units must be a whole number')
    .min(MIN_UNITS, `Total units must be at least ${MIN_UNITS}`)
    .max(MAX_UNITS, `Total units must be at most ${MAX_UNITS.toLocaleString()}`),
  initialOwner: z
    .string()
    .min(MIN_NAME_LENGTH, 'Initial owner is required')
    .max(MAX_NAME_LENGTH, `Initial owner must be at most ${MAX_NAME_LENGTH} characters`)
    .regex(NAME_PATTERN, 'Initial owner cannot have leading or trailing whitespace'),
});

export type CreateFundInput = z.infer<typeof createFundSchema>;
