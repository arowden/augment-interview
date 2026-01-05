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

  useEffect(() => {
    if (!isOpen) return;

    firstInputRef.current?.focus();

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
        return;
      }

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
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-slate-950/80 backdrop-blur-sm"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Modal */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div
          ref={modalRef}
          className="relative w-full max-w-md card p-6 animate-slide-up"
        >
          <div className="flex items-center justify-between mb-6">
            <h2 id="modal-title" className="text-lg font-semibold text-white">
              Create Fund
            </h2>
            <button
              type="button"
              onClick={onClose}
              className="w-8 h-8 rounded-lg flex items-center justify-center text-slate-400 hover:text-white hover:bg-slate-800 transition-colors"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-5">
            {/* Name field */}
            <div>
              <label htmlFor="fund-name" className="label">
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
                className={`input ${errors.name ? 'input-error' : ''}`}
                placeholder="Growth Fund I"
              />
              {errors.name && (
                <p id="name-error" role="alert" className="mt-1.5 text-sm text-red-400">
                  {errors.name.message}
                </p>
              )}
            </div>

            {/* Total Units field */}
            <div>
              <label htmlFor="total-units" className="label">
                Total Units
              </label>
              <input
                {...register('totalUnits', { valueAsNumber: true })}
                type="number"
                id="total-units"
                aria-invalid={errors.totalUnits ? 'true' : 'false'}
                aria-describedby={errors.totalUnits ? 'units-error' : undefined}
                className={`input ${errors.totalUnits ? 'input-error' : ''}`}
                placeholder="1000000"
              />
              {errors.totalUnits && (
                <p id="units-error" role="alert" className="mt-1.5 text-sm text-red-400">
                  {errors.totalUnits.message}
                </p>
              )}
            </div>

            {/* Initial Owner field */}
            <div>
              <label htmlFor="initial-owner" className="label">
                Initial Owner
              </label>
              <input
                {...register('initialOwner')}
                type="text"
                id="initial-owner"
                aria-invalid={errors.initialOwner ? 'true' : 'false'}
                aria-describedby={errors.initialOwner ? 'owner-error' : undefined}
                className={`input ${errors.initialOwner ? 'input-error' : ''}`}
                placeholder="Founder LLC"
              />
              {errors.initialOwner && (
                <p id="owner-error" role="alert" className="mt-1.5 text-sm text-red-400">
                  {errors.initialOwner.message}
                </p>
              )}
            </div>

            {/* API error */}
            {error != null && (
              <ApiErrorAlert error={error} fallbackMessage="Failed to create fund" />
            )}

            {/* Actions */}
            <div className="flex justify-end gap-3 pt-2">
              <button type="button" onClick={onClose} className="btn-secondary">
                Cancel
              </button>
              <button type="submit" disabled={isLoading} className="btn-primary">
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
