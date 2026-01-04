import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import { createTransferSchema, type CreateTransferInput } from '../schemas/transfer';
import { LoadingSpinner } from './LoadingSpinner';
import { ApiErrorAlert } from './ApiErrorAlert';

interface TransferFormProps {
  onSubmit: (data: CreateTransferInput) => Promise<void>;
  isLoading: boolean;
  error: unknown;
}

export function TransferForm({ onSubmit, isLoading, error }: TransferFormProps) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateTransferInput>({
    resolver: zodResolver(createTransferSchema),
    defaultValues: {
      fromOwner: '',
      toOwner: '',
      units: 1,
    },
  });

  const handleFormSubmit = async (data: CreateTransferInput) => {
    await onSubmit(data);
    reset();
  };

  return (
    <form
      onSubmit={handleSubmit(handleFormSubmit)}
      className="rounded-lg border border-gray-200 bg-white p-6 shadow"
    >
      <h3 className="text-lg font-semibold text-gray-900">Transfer Units</h3>
      <p className="mt-1 text-sm text-gray-600">
        Move ownership units from one party to another.
      </p>

      <div className="mt-4 grid gap-4 sm:grid-cols-3">
        {/* From Owner field. */}
        <div>
          <label
            htmlFor="from-owner"
            className="block text-sm font-medium text-gray-700"
          >
            From Owner
          </label>
          <input
            {...register('fromOwner')}
            type="text"
            id="from-owner"
            aria-invalid={errors.fromOwner ? 'true' : 'false'}
            aria-describedby={errors.fromOwner ? 'from-owner-error' : undefined}
            className={`mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus-visible-ring text-gray-900 placeholder:text-gray-400 ${
              errors.fromOwner ? 'border-error-500' : 'border-gray-300'
            }`}
            placeholder="Founder LLC"
          />
          {errors.fromOwner && (
            <p
              id="from-owner-error"
              role="alert"
              className="mt-1 text-sm text-error-500"
            >
              {errors.fromOwner.message}
            </p>
          )}
        </div>

        {/* To Owner field. */}
        <div>
          <label
            htmlFor="to-owner"
            className="block text-sm font-medium text-gray-700"
          >
            To Owner
          </label>
          <input
            {...register('toOwner')}
            type="text"
            id="to-owner"
            aria-invalid={errors.toOwner ? 'true' : 'false'}
            aria-describedby={errors.toOwner ? 'to-owner-error' : undefined}
            className={`mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus-visible-ring text-gray-900 placeholder:text-gray-400 ${
              errors.toOwner ? 'border-error-500' : 'border-gray-300'
            }`}
            placeholder="Investor A"
          />
          {errors.toOwner && (
            <p
              id="to-owner-error"
              role="alert"
              className="mt-1 text-sm text-error-500"
            >
              {errors.toOwner.message}
            </p>
          )}
        </div>

        {/* Units field. */}
        <div>
          <label
            htmlFor="transfer-units"
            className="block text-sm font-medium text-gray-700"
          >
            Units
          </label>
          <input
            {...register('units', { valueAsNumber: true })}
            type="number"
            id="transfer-units"
            min={1}
            aria-invalid={errors.units ? 'true' : 'false'}
            aria-describedby={errors.units ? 'transfer-units-error' : undefined}
            className={`mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus-visible-ring text-gray-900 placeholder:text-gray-400 ${
              errors.units ? 'border-error-500' : 'border-gray-300'
            }`}
            placeholder="100000"
          />
          {errors.units && (
            <p
              id="transfer-units-error"
              role="alert"
              className="mt-1 text-sm text-error-500"
            >
              {errors.units.message}
            </p>
          )}
        </div>
      </div>

      {/* API error. */}
      {error != null && (
        <div className="mt-4">
          <ApiErrorAlert error={error} fallbackMessage="Failed to execute transfer" />
        </div>
      )}

      {/* Submit button. */}
      <div className="mt-4">
        <button
          type="submit"
          disabled={isLoading}
          className="inline-flex items-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 focus-visible-ring disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading && <LoadingSpinner size="sm" className="mr-2" />}
          Execute Transfer
        </button>
      </div>
    </form>
  );
}
