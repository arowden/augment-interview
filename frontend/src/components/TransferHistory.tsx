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
      <div className="card border-dashed border-2 border-slate-700 p-8 text-center">
        <div className="w-12 h-12 rounded-full bg-slate-800 flex items-center justify-center mx-auto">
          <svg className="w-6 h-6 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" />
          </svg>
        </div>
        <p className="mt-3 text-slate-400">No transfers yet</p>
        <p className="mt-1 text-sm text-slate-500">Transfers will appear here once you move units between owners.</p>
      </div>
    );
  }

  const sortedTransfers = [...transfers].sort(
    (a, b) => new Date(b.transferredAt).getTime() - new Date(a.transferredAt).getTime()
  );

  return (
    <div className="card overflow-hidden">
      <div className="divide-y divide-slate-800/50">
        {sortedTransfers.map((transfer) => (
          <div
            key={transfer.id}
            className="p-4 hover:bg-slate-800/30 transition-colors"
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-indigo-500/10 flex items-center justify-center">
                  <svg className="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                </div>
                <div>
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-medium text-white">{transfer.fromOwner}</span>
                    <svg className="w-4 h-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                    </svg>
                    <span className="font-medium text-white">{transfer.toOwner}</span>
                  </div>
                  <p className="text-xs text-slate-500 mt-0.5">
                    {formatDateTime(transfer.transferredAt)}
                  </p>
                </div>
              </div>
              <div className="text-right">
                <span className="text-sm font-semibold text-indigo-400 tabular-nums">
                  {transfer.units.toLocaleString()}
                </span>
                <p className="text-xs text-slate-500">units</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function TransferHistorySkeleton() {
  return (
    <div className="card overflow-hidden">
      <div className="divide-y divide-slate-800/50">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="p-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="skeleton w-10 h-10 rounded-full" />
                <div>
                  <div className="skeleton h-4 w-48" />
                  <div className="skeleton h-3 w-32 mt-1.5" />
                </div>
              </div>
              <div className="text-right">
                <div className="skeleton h-4 w-20 ml-auto" />
                <div className="skeleton h-3 w-12 ml-auto mt-1" />
              </div>
            </div>
          </div>
        ))}
      </div>
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
