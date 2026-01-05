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
    <form onSubmit={handleSubmit(handleFormSubmit)} className="card p-6">
      <h3 className="text-lg font-semibold text-white">Transfer Units</h3>
      <p className="mt-1 text-sm text-slate-400">
        Move ownership units from one party to another.
      </p>

      <div className="mt-4 grid gap-4 sm:grid-cols-3">
        {/* From Owner field */}
        <div>
          <label htmlFor="from-owner" className="label">
            From Owner
          </label>
          <input
            {...register('fromOwner')}
            type="text"
            id="from-owner"
            aria-invalid={errors.fromOwner ? 'true' : 'false'}
            aria-describedby={errors.fromOwner ? 'from-owner-error' : undefined}
            className={`input ${errors.fromOwner ? 'input-error' : ''}`}
            placeholder="Founder LLC"
          />
          {errors.fromOwner && (
            <p id="from-owner-error" role="alert" className="mt-1.5 text-sm text-red-400">
              {errors.fromOwner.message}
            </p>
          )}
        </div>

        {/* To Owner field */}
        <div>
          <label htmlFor="to-owner" className="label">
            To Owner
          </label>
          <input
            {...register('toOwner')}
            type="text"
            id="to-owner"
            aria-invalid={errors.toOwner ? 'true' : 'false'}
            aria-describedby={errors.toOwner ? 'to-owner-error' : undefined}
            className={`input ${errors.toOwner ? 'input-error' : ''}`}
            placeholder="Investor A"
          />
          {errors.toOwner && (
            <p id="to-owner-error" role="alert" className="mt-1.5 text-sm text-red-400">
              {errors.toOwner.message}
            </p>
          )}
        </div>

        {/* Units field */}
        <div>
          <label htmlFor="transfer-units" className="label">
            Units
          </label>
          <input
            {...register('units', { valueAsNumber: true })}
            type="number"
            id="transfer-units"
            min={1}
            aria-invalid={errors.units ? 'true' : 'false'}
            aria-describedby={errors.units ? 'transfer-units-error' : undefined}
            className={`input ${errors.units ? 'input-error' : ''}`}
            placeholder="100000"
          />
          {errors.units && (
            <p id="transfer-units-error" role="alert" className="mt-1.5 text-sm text-red-400">
              {errors.units.message}
            </p>
          )}
        </div>
      </div>

      {/* API error */}
      {error != null && (
        <div className="mt-4">
          <ApiErrorAlert error={error} fallbackMessage="Failed to execute transfer" />
        </div>
      )}

      {/* Submit button */}
      <div className="mt-6">
        <button type="submit" disabled={isLoading} className="btn-primary">
          {isLoading && <LoadingSpinner size="sm" className="mr-2" />}
          Save
        </button>
      </div>
    </form>
  );
}
