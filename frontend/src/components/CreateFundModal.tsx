import { useEffect, useRef } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';

import { createFundSchema, type CreateFundInput } from '../schemas/fund';
import { LoadingSpinner } from './LoadingSpinner';
import { ApiErrorAlert } from './ApiErrorAlert';

interface CreateFundModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateFundInput) => Promise<void>;
  isLoading: boolean;
  error: unknown;
}

export function CreateFundModal({
  isOpen,
  onClose,
  onSubmit,
  isLoading,
  error,
}: CreateFundModalProps) {
  const modalRef = useRef<HTMLDivElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateFundInput>({
    resolver: zodResolver(createFundSchema),
    defaultValues: {
      name: '',
      totalUnits: 1000000,
      initialOwner: '',
    },
  });

  // Focus trap and escape key handling.
  useEffect(() => {
    if (!isOpen) return;

    // Focus the first input when modal opens.
    firstInputRef.current?.focus();

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
        return;
      }

      // Focus trap.
      if (event.key === 'Tab' && modalRef.current) {
        const focusableElements = modalRef.current.querySelectorAll<HTMLElement>(
          'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
        );
        const firstElement = focusableElements[0];
        const lastElement = focusableElements[focusableElements.length - 1];

        if (event.shiftKey && document.activeElement === firstElement) {
          event.preventDefault();
          lastElement?.focus();
        } else if (!event.shiftKey && document.activeElement === lastElement) {
          event.preventDefault();
          firstElement?.focus();
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  // Reset form when modal closes.
  useEffect(() => {
    if (!isOpen) {
      reset();
    }
  }, [isOpen, reset]);

  if (!isOpen) return null;

  const handleFormSubmit = async (data: CreateFundInput) => {
    await onSubmit(data);
  };

  return (
    <div
      className="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      {/* Backdrop. */}
      <div
        className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Modal. */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div
          ref={modalRef}
          className="relative w-full max-w-md transform rounded-lg bg-white p-6 shadow-xl transition-all"
        >
          <h2
            id="modal-title"
            className="text-lg font-semibold text-gray-900"
          >
            Create New Fund
          </h2>

          <form
            onSubmit={handleSubmit(handleFormSubmit)}
            className="mt-4 space-y-4"
          >
            {/* Name field. */}
            <div>
              <label
                htmlFor="fund-name"
                className="block text-sm font-medium text-gray-700"
              >
                Fund Name
              </label>
              <input
                {...register('name')}
                ref={(e) => {
                  register('name').ref(e);
                  (firstInputRef as React.MutableRefObject<HTMLInputElement | null>).current = e;
                }}
                type="text"
                id="fund-name"
                aria-invalid={errors.name ? 'true' : 'false'}
                aria-describedby={errors.name ? 'name-error' : undefined}
                className={`mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus-visible-ring text-gray-900 placeholder:text-gray-400 ${
                  errors.name ? 'border-error-500' : 'border-gray-300'
                }`}
                placeholder="Growth Fund I"
              />
              {errors.name && (
                <p
                  id="name-error"
                  role="alert"
                  className="mt-1 text-sm text-error-500"
                >
                  {errors.name.message}
                </p>
              )}
            </div>

            {/* Total Units field. */}
            <div>
              <label
                htmlFor="total-units"
                className="block text-sm font-medium text-gray-700"
              >
                Total Units
              </label>
              <input
                {...register('totalUnits', { valueAsNumber: true })}
                type="number"
                id="total-units"
                aria-invalid={errors.totalUnits ? 'true' : 'false'}
                aria-describedby={errors.totalUnits ? 'units-error' : undefined}
                className={`mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus-visible-ring text-gray-900 placeholder:text-gray-400 ${
                  errors.totalUnits ? 'border-error-500' : 'border-gray-300'
                }`}
                placeholder="1000000"
              />
              {errors.totalUnits && (
                <p
                  id="units-error"
                  role="alert"
                  className="mt-1 text-sm text-error-500"
                >
                  {errors.totalUnits.message}
                </p>
              )}
            </div>

            {/* Initial Owner field. */}
            <div>
              <label
                htmlFor="initial-owner"
                className="block text-sm font-medium text-gray-700"
              >
                Initial Owner
              </label>
              <input
                {...register('initialOwner')}
                type="text"
                id="initial-owner"
                aria-invalid={errors.initialOwner ? 'true' : 'false'}
                aria-describedby={errors.initialOwner ? 'owner-error' : undefined}
                className={`mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus-visible-ring text-gray-900 placeholder:text-gray-400 ${
                  errors.initialOwner ? 'border-error-500' : 'border-gray-300'
                }`}
                placeholder="Founder LLC"
              />
              {errors.initialOwner && (
                <p
                  id="owner-error"
                  role="alert"
                  className="mt-1 text-sm text-error-500"
                >
                  {errors.initialOwner.message}
                </p>
              )}
            </div>

            {/* API error. */}
            {error != null && (
              <ApiErrorAlert error={error} fallbackMessage="Failed to create fund" />
            )}

            {/* Actions. */}
            <div className="mt-6 flex justify-end gap-3">
              <button
                type="button"
                onClick={onClose}
                className="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 focus-visible-ring"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isLoading}
                className="inline-flex items-center rounded-md bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 focus-visible-ring disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading && <LoadingSpinner size="sm" className="mr-2" />}
                Create Fund
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
