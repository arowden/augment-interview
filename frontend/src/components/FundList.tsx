import type { Fund } from '../api/client';
import { FundCard } from './FundCard';
import { FundCardSkeleton } from './FundCardSkeleton';

interface FundListProps {
  funds: Fund[] | undefined;
  isLoading: boolean;
  error: Error | null;
  onRetry: () => void;
}

export function FundList({ funds, isLoading, error, onRetry }: FundListProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <FundCardSkeleton key={i} />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div
        role="alert"
        className="glass-card p-8 text-center border-error-500/20"
      >
        <div className="w-12 h-12 rounded-full bg-error-500/10 flex items-center justify-center mx-auto">
          <svg className="w-6 h-6 text-error-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <p className="mt-4 text-slate-300 font-medium">
          Unable to load funds
        </p>
        <p className="mt-1 text-sm text-slate-500">
          {error.message || 'An unexpected error occurred'}
        </p>
        <button
          type="button"
          onClick={onRetry}
          className="btn-primary mt-6"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Try Again
        </button>
      </div>
    );
  }

  if (!funds || funds.length === 0) {
    return (
      <div className="glass-card p-12 text-center border-dashed border-2 border-dark-600">
        <div className="w-16 h-16 rounded-2xl bg-dark-800 flex items-center justify-center mx-auto">
          <svg className="w-8 h-8 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
          </svg>
        </div>
        <h3 className="mt-4 text-lg font-heading font-semibold text-slate-200">
          No funds yet
        </h3>
        <p className="mt-2 text-slate-500">
          Create your first fund to get started tracking ownership.
        </p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {funds.map((fund) => (
        <FundCard key={fund.id} fund={fund} />
      ))}
    </div>
  );
}
