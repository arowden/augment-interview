import type { CapTableEntry } from '../api/client';

interface CapTableProps {
  entries: CapTableEntry[];
  totalUnits: number;
  isLoading: boolean;
}

export function CapTable({ entries, totalUnits, isLoading }: CapTableProps) {
  if (isLoading) {
    return <CapTableSkeleton />;
  }

  if (entries.length === 0) {
    return (
      <div className="rounded-lg border-2 border-dashed border-dark-600 p-8 text-center bg-dark-800/50">
        <p className="text-slate-400">No ownership entries yet.</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow">
      <table
        className="min-w-full divide-y divide-gray-200"
        aria-label="Cap table showing ownership distribution"
      >
        <thead className="bg-gray-50">
          <tr>
            <th
              scope="col"
              className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
            >
              Owner
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500"
            >
              Units
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500"
            >
              Percentage
            </th>
            <th
              scope="col"
              className="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500"
            >
              Acquired
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200 bg-white">
          {entries.map((entry) => (
            <tr key={entry.ownerName}>
              <td className="whitespace-nowrap px-6 py-4 text-sm font-medium text-gray-900">
                {entry.ownerName}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-gray-600">
                {entry.units.toLocaleString()}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-gray-600">
                {formatPercentage(entry.units, totalUnits)}
              </td>
              <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-gray-500">
                {formatDate(entry.acquiredAt)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function CapTableSkeleton() {
  return (
    <div className="overflow-hidden rounded-lg border border-gray-200 bg-white shadow">
      <div className="divide-y divide-gray-200">
        <div className="bg-gray-50 px-6 py-3">
          <div className="skeleton h-4 w-full rounded" />
        </div>
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="flex gap-4 px-6 py-4">
            <div className="skeleton h-4 w-1/4 rounded" />
            <div className="skeleton h-4 w-1/4 rounded" />
            <div className="skeleton h-4 w-1/4 rounded" />
            <div className="skeleton h-4 w-1/4 rounded" />
          </div>
        ))}
      </div>
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
