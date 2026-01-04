import type { Transfer } from '../api/client';

interface TransferHistoryProps {
  transfers: Transfer[];
  isLoading: boolean;
}

export function TransferHistory({ transfers, isLoading }: TransferHistoryProps) {
  if (isLoading) {
    return <TransferHistorySkeleton />;
  }

  if (transfers.length === 0) {
    return (
      <div className="rounded-lg border-2 border-dashed border-dark-600 p-8 text-center bg-dark-800/50">
        <p className="text-slate-400">No transfers yet.</p>
      </div>
    );
  }

  // Sort by date descending (most recent first).
  const sortedTransfers = [...transfers].sort(
    (a, b) => new Date(b.transferredAt).getTime() - new Date(a.transferredAt).getTime()
  );

  return (
    <div className="space-y-3">
      {sortedTransfers.map((transfer) => (
        <div
          key={transfer.id}
          className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-sm">
              <span className="font-medium text-gray-900">{transfer.fromOwner}</span>
              <span className="text-gray-400">&rarr;</span>
              <span className="font-medium text-gray-900">{transfer.toOwner}</span>
            </div>
            <span className="text-sm font-semibold text-primary-600">
              {transfer.units.toLocaleString()} units
            </span>
          </div>
          <p className="mt-1 text-xs text-gray-500">
            {formatDateTime(transfer.transferredAt)}
          </p>
        </div>
      ))}
    </div>
  );
}

function TransferHistorySkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 3 }).map((_, i) => (
        <div
          key={i}
          className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
        >
          <div className="flex items-center justify-between">
            <div className="skeleton h-4 w-1/2 rounded" />
            <div className="skeleton h-4 w-1/4 rounded" />
          </div>
          <div className="skeleton mt-2 h-3 w-1/3 rounded" />
        </div>
      ))}
    </div>
  );
}

function formatDateTime(dateString: string): string {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}
