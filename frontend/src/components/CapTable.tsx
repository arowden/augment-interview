import { useState } from 'react';
import type { CapTableEntry, Transfer } from '../api/client';
import { OwnerDetailModal } from './OwnerDetailModal';

interface CapTableProps {
  entries: CapTableEntry[];
  totalUnits: number;
  isLoading: boolean;
  transfers?: Transfer[];
}

export function CapTable({ entries, totalUnits, isLoading, transfers = [] }: CapTableProps) {
  const [selectedOwner, setSelectedOwner] = useState<string | null>(null);

  if (isLoading) {
    return <CapTableSkeleton />;
  }

  if (entries.length === 0) {
    return (
      <div className="card border-dashed border-2 border-slate-700 p-8 text-center">
        <p className="text-slate-400">No ownership entries yet.</p>
      </div>
    );
  }

  const selectedEntry = entries.find(e => e.ownerName === selectedOwner);
  const ownerTransfers = transfers.filter(
    t => t.fromOwner === selectedOwner || t.toOwner === selectedOwner
  );

  return (
    <>
      <div className="card overflow-hidden">
        <table
          className="min-w-full"
          aria-label="Cap table showing ownership distribution"
        >
          <thead>
            <tr className="border-b border-slate-800">
              <th
                scope="col"
                className="px-6 py-4 text-left text-sm font-medium text-slate-400"
              >
                Owner
              </th>
              <th
                scope="col"
                className="px-6 py-4 text-right text-sm font-medium text-slate-400"
              >
                Units
              </th>
              <th
                scope="col"
                className="px-6 py-4 text-right text-sm font-medium text-slate-400"
              >
                Percentage
              </th>
              <th
                scope="col"
                className="px-6 py-4 text-right text-sm font-medium text-slate-400"
              >
                Acquired
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/50">
            {entries.map((entry) => (
              <tr key={entry.ownerName} className="hover:bg-slate-800/30 transition-colors">
                <td className="whitespace-nowrap px-6 py-4">
                  <button
                    onClick={() => setSelectedOwner(entry.ownerName)}
                    className="font-medium text-white hover:text-indigo-400 transition-colors text-left"
                    title="Click to view owner details"
                  >
                    {entry.ownerName}
                  </button>
                </td>
                <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-slate-300 tabular-nums">
                  {entry.units.toLocaleString()}
                </td>
                <td className="whitespace-nowrap px-6 py-4 text-right text-sm">
                  <span className="inline-flex items-center gap-2">
                    <span className="text-slate-300">{formatPercentage(entry.units, totalUnits)}</span>
                    <PercentageBar percentage={(entry.units / totalUnits) * 100} />
                  </span>
                </td>
                <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-slate-500">
                  {formatDate(entry.acquiredAt)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {selectedOwner && selectedEntry && (
        <OwnerDetailModal
          ownerName={selectedOwner}
          entry={selectedEntry}
          totalUnits={totalUnits}
          transfers={ownerTransfers}
          onClose={() => setSelectedOwner(null)}
        />
      )}
    </>
  );
}

function PercentageBar({ percentage }: { percentage: number }) {
  return (
    <div className="w-16 h-1.5 bg-slate-700 rounded-full overflow-hidden">
      <div
        className="h-full bg-indigo-500 rounded-full transition-all duration-300"
        style={{ width: `${Math.min(percentage, 100)}%` }}
      />
    </div>
  );
}

function CapTableSkeleton() {
  return (
    <div className="card overflow-hidden">
      <table className="min-w-full">
        <thead>
          <tr className="border-b border-slate-800">
            <th className="px-6 py-4 text-left"><div className="skeleton h-4 w-16" /></th>
            <th className="px-6 py-4 text-right"><div className="skeleton h-4 w-12 ml-auto" /></th>
            <th className="px-6 py-4 text-right"><div className="skeleton h-4 w-20 ml-auto" /></th>
            <th className="px-6 py-4 text-right"><div className="skeleton h-4 w-24 ml-auto" /></th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-800/50">
          {Array.from({ length: 3 }).map((_, i) => (
            <tr key={i}>
              <td className="px-6 py-4"><div className="skeleton h-4 w-32" /></td>
              <td className="px-6 py-4"><div className="skeleton h-4 w-20 ml-auto" /></td>
              <td className="px-6 py-4"><div className="skeleton h-4 w-24 ml-auto" /></td>
              <td className="px-6 py-4"><div className="skeleton h-4 w-24 ml-auto" /></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function formatPercentage(units: number, totalUnits: number): string {
  if (totalUnits === 0) return '0.00%';
  const percentage = (units / totalUnits) * 100;
  return `${percentage.toFixed(2)}%`;
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}
