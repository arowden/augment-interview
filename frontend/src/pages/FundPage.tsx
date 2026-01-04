import { useParams, Link } from 'react-router-dom';

import { useFund } from '../hooks/useFunds';
import { useCapTable } from '../hooks/useCapTable';
import { useTransfers, useCreateTransfer } from '../hooks/useTransfers';
import { CapTable } from '../components/CapTable';
import { TransferForm } from '../components/TransferForm';
import { TransferHistory } from '../components/TransferHistory';
import { LoadingSpinner } from '../components/LoadingSpinner';
import type { CreateTransferInput } from '../schemas/transfer';

export function FundPage() {
  const { id } = useParams<{ id: string }>();
  const fundId = id ?? '';

  const { data: fund, isLoading: fundLoading, error: fundError } = useFund(fundId);
  const { data: capTable, isLoading: capTableLoading } = useCapTable(fundId);
  const { data: transferList, isLoading: transfersLoading } = useTransfers(fundId);
  const createTransfer = useCreateTransfer(fundId);

  const handleTransfer = async (data: CreateTransferInput) => {
    await createTransfer.mutateAsync({
      fromOwner: data.fromOwner,
      toOwner: data.toOwner,
      units: data.units,
    });
  };

  if (fundLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (fundError || !fund) {
    return (
      <div
        role="alert"
        className="rounded-lg bg-error-50 p-6 text-center"
      >
        <h2 className="text-lg font-semibold text-error-600">Fund Not Found</h2>
        <p className="mt-2 text-error-600">
          The fund you're looking for doesn't exist or has been removed.
        </p>
        <Link
          to="/"
          className="mt-4 inline-block rounded-md bg-primary-600 px-4 py-2 text-white hover:bg-primary-700 focus-visible-ring"
        >
          Back to Dashboard
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Fund header. */}
      <div>
        <h1 className="text-2xl font-bold text-slate-100">{fund.name}</h1>
        <p className="mt-1 text-sm text-slate-300">
          Total Units: {fund.totalUnits.toLocaleString()}
        </p>
        <p className="text-sm text-slate-400">
          Created: {new Date(fund.createdAt).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
          })}
        </p>
      </div>

      {/* Cap table. */}
      <section>
        <h2 className="mb-4 text-lg font-semibold text-slate-100">Cap Table</h2>
        <CapTable
          entries={capTable?.entries ?? []}
          totalUnits={fund.totalUnits}
          isLoading={capTableLoading}
        />
      </section>

      {/* Transfer form. */}
      <section>
        <TransferForm
          onSubmit={handleTransfer}
          isLoading={createTransfer.isPending}
          error={createTransfer.error}
        />
      </section>

      {/* Transfer history. */}
      <section>
        <h2 className="mb-4 text-lg font-semibold text-slate-100">Transfer History</h2>
        <TransferHistory
          transfers={transferList?.transfers ?? []}
          isLoading={transfersLoading}
        />
      </section>
    </div>
  );
}
