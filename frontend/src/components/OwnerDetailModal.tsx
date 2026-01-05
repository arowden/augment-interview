import { useEffect, useRef } from 'react';
import type { CapTableEntry, Transfer } from '../api/client';

interface OwnerDetailModalProps {
  ownerName: string;
  entry: CapTableEntry;
  totalUnits: number;
  transfers: Transfer[];
  onClose: () => void;
}

export function OwnerDetailModal({
  ownerName,
  entry,
  totalUnits,
  transfers,
  onClose,
}: OwnerDetailModalProps) {
  const dialogRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [onClose]);

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) onClose();
  };

  const percentage = totalUnits > 0 ? (entry.units / totalUnits) * 100 : 0;

  const sortedTransfers = [...transfers].sort(
    (a, b) => new Date(b.transferredAt).getTime() - new Date(a.transferredAt).getTime()
  );

  const totalReceived = transfers
    .filter(t => t.toOwner === ownerName)
    .reduce((sum, t) => sum + t.units, 0);
  const totalSent = transfers
    .filter(t => t.fromOwner === ownerName)
    .reduce((sum, t) => sum + t.units, 0);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
      onClick={handleBackdropClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby="owner-modal-title"
    >
      <div
        ref={dialogRef}
        className="w-full max-w-lg bg-slate-900 rounded-xl shadow-2xl border border-slate-800 overflow-hidden animate-in fade-in zoom-in-95 duration-200"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-slate-800">
          <div>
            <h2 id="owner-modal-title" className="text-lg font-semibold text-white">
              {ownerName}
            </h2>
            <p className="text-sm text-slate-400">Owner Details</p>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-slate-400 hover:text-white hover:bg-slate-800 rounded-lg transition-colors"
            aria-label="Close modal"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Ownership Stats */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-slate-800/50 rounded-lg p-4">
              <p className="text-sm text-slate-400">Current Holdings</p>
              <p className="text-2xl font-bold text-white tabular-nums">
                {entry.units.toLocaleString()}
              </p>
              <p className="text-sm text-slate-500">units</p>
            </div>
            <div className="bg-slate-800/50 rounded-lg p-4">
              <p className="text-sm text-slate-400">Ownership</p>
              <p className="text-2xl font-bold text-indigo-400 tabular-nums">
                {percentage.toFixed(2)}%
              </p>
              <div className="mt-1 h-1.5 bg-slate-700 rounded-full overflow-hidden">
                <div
                  className="h-full bg-indigo-500 rounded-full"
                  style={{ width: `${Math.min(percentage, 100)}%` }}
                />
              </div>
            </div>
          </div>

          {/* Transfer Stats */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-green-500/10 border border-green-500/20 rounded-lg p-4">
              <p className="text-sm text-green-400">Total Received</p>
              <p className="text-xl font-bold text-green-400 tabular-nums">
                +{totalReceived.toLocaleString()}
              </p>
            </div>
            <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4">
              <p className="text-sm text-red-400">Total Sent</p>
              <p className="text-xl font-bold text-red-400 tabular-nums">
                -{totalSent.toLocaleString()}
              </p>
            </div>
          </div>

          {/* Timeline */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-medium text-slate-300">Transaction History</h3>
              <span className="text-xs text-slate-500">{transfers.length} transactions</span>
            </div>

            {sortedTransfers.length === 0 ? (
              <p className="text-sm text-slate-500 text-center py-4">
                No transactions yet
              </p>
            ) : (
              <div className="space-y-2 max-h-48 overflow-y-auto pr-2">
                {sortedTransfers.map((transfer) => {
                  const isReceived = transfer.toOwner === ownerName;
                  return (
                    <div
                      key={transfer.id}
                      className="flex items-center justify-between p-3 bg-slate-800/30 rounded-lg"
                    >
                      <div className="flex items-center gap-3">
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                          isReceived ? 'bg-green-500/20' : 'bg-red-500/20'
                        }`}>
                          {isReceived ? (
                            <svg className="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                            </svg>
                          ) : (
                            <svg className="w-4 h-4 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M5 10l7-7m0 0l7 7m-7-7v18" />
                            </svg>
                          )}
                        </div>
                        <div>
                          <p className="text-sm text-white">
                            {isReceived ? `From ${transfer.fromOwner}` : `To ${transfer.toOwner}`}
                          </p>
                          <p className="text-xs text-slate-500">
                            {formatDateTime(transfer.transferredAt)}
                          </p>
                        </div>
                      </div>
                      <span className={`text-sm font-medium tabular-nums ${
                        isReceived ? 'text-green-400' : 'text-red-400'
                      }`}>
                        {isReceived ? '+' : '-'}{transfer.units.toLocaleString()}
                      </span>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          {/* Acquisition date */}
          <div className="text-xs text-slate-500 text-center border-t border-slate-800 pt-4">
            First acquired: {formatDate(entry.acquiredAt)}
          </div>
        </div>
      </div>
    </div>
  );
}

function formatDateTime(dateString: string): string {
  return new Date(dateString).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}
