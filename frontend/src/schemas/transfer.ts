import { z } from 'zod';
import {
  MIN_NAME_LENGTH,
  MAX_NAME_LENGTH,
  NAME_PATTERN,
  MIN_UNITS,
  MAX_UNITS,
} from '../validation/constraints';

export const createTransferSchema = z
  .object({
    fromOwner: z
      .string()
      .min(MIN_NAME_LENGTH, 'From owner is required')
      .max(MAX_NAME_LENGTH, `From owner must be at most ${MAX_NAME_LENGTH} characters`)
      .regex(NAME_PATTERN, 'From owner cannot have leading or trailing whitespace'),
    toOwner: z
      .string()
      .min(MIN_NAME_LENGTH, 'To owner is required')
      .max(MAX_NAME_LENGTH, `To owner must be at most ${MAX_NAME_LENGTH} characters`)
      .regex(NAME_PATTERN, 'To owner cannot have leading or trailing whitespace'),
    units: z
      .number({ invalid_type_error: 'Units must be a number' })
      .int('Units must be a whole number')
      .min(MIN_UNITS, `Units must be at least ${MIN_UNITS}`)
      .max(MAX_UNITS, `Units must be at most ${MAX_UNITS.toLocaleString()}`),
  })
  .refine((data) => data.fromOwner !== data.toOwner, {
    message: 'Cannot transfer to same owner',
    path: ['toOwner'],
  });

export type CreateTransferInput = z.infer<typeof createTransferSchema>;
