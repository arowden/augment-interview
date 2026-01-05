import { Link } from 'react-router-dom';
import { useQueries } from '@tanstack/react-query';

import { useFunds } from '../hooks/useFunds';
import { CapTableService, type Fund, type CapTableEntry } from '../api/client';
import { LoadingSpinner } from '../components/LoadingSpinner';

interface OwnerHolding {
  fundId: string;
  fundName: string;
  units: number;
  percentage: number;
  totalUnits: number;
  acquiredAt: string;
}

interface AggregatedOwner {
  name: string;
  holdings: OwnerHolding[];
  totalFunds: number;
}

export function OwnersPage() {
  const { data: fundList, isLoading: fundsLoading, error: fundsError } = useFunds();
  const funds = fundList?.funds ?? [];

  const capTableQueries = useQueries({
    queries: funds.map((fund) => ({
      queryKey: ['capTables', fund.id],
      queryFn: () => CapTableService.getCapTable(fund.id),
      enabled: funds.length > 0,
    })),
  });

  const isLoading = fundsLoading || capTableQueries.some((q) => q.isLoading);
  const hasError = fundsError || capTableQueries.some((q) => q.error);

  const aggregatedOwners: AggregatedOwner[] = [];
  const ownerMap = new Map<string, OwnerHolding[]>();

  if (!isLoading && !hasError) {
    capTableQueries.forEach((query, index) => {
      const fund = funds[index];
      if (query.data?.entries && fund) {
        query.data.entries.forEach((entry: CapTableEntry) => {
          const holdings = ownerMap.get(entry.ownerName) ?? [];
          holdings.push({
            fundId: fund.id,
            fundName: fund.name,
            units: entry.units,
            percentage: entry.percentage,
            totalUnits: fund.totalUnits,
            acquiredAt: entry.acquiredAt,
          });
          ownerMap.set(entry.ownerName, holdings);
        });
      }
    });

    ownerMap.forEach((holdings, name) => {
      aggregatedOwners.push({
        name,
        holdings,
        totalFunds: holdings.length,
      });
    });

    aggregatedOwners.sort((a, b) => {
      if (b.totalFunds !== a.totalFunds) {
        return b.totalFunds - a.totalFunds;
      }
      return a.name.localeCompare(b.name);
    });
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
    );
  }

  if (hasError) {
    return (
      <div role="alert" className="rounded-lg bg-red-900/20 border border-red-800 p-6 text-center">
        <h2 className="text-lg font-semibold text-red-400">Error Loading Data</h2>
        <p className="mt-2 text-red-300">Failed to load ownership data. Please try again.</p>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-semibold text-white">Owners</h1>
        <p className="mt-1 text-slate-400">
          View all owners and their holdings across funds
        </p>
      </div>

      {/* Stats */}
      {aggregatedOwners.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="card p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-indigo-500/10 flex items-center justify-center">
                <svg className="w-6 h-6 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
                </svg>
              </div>
              <div>
                <p className="stat-label">Total Owners</p>
                <p className="stat-value">{aggregatedOwners.length}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center">
                <svg className="w-6 h-6 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
                </svg>
              </div>
              <div>
                <p className="stat-label">Funds Tracked</p>
                <p className="stat-value">{funds.length}</p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Owners list */}
      {aggregatedOwners.length === 0 ? (
        <div className="card border-dashed border-2 border-slate-700 p-8 text-center">
          <svg className="mx-auto h-12 w-12 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
          </svg>
          <h3 className="mt-4 text-lg font-medium text-slate-300">No owners yet</h3>
          <p className="mt-2 text-slate-400">
            Create a fund to start tracking ownership.
          </p>
          <Link
            to="/"
            className="mt-4 inline-block btn-primary"
          >
            Go to Dashboard
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {aggregatedOwners.map((owner) => (
            <div key={owner.name} className="card overflow-hidden">
              {/* Owner header */}
              <div className="px-6 py-4 bg-slate-800/30 border-b border-slate-800">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-indigo-500/20 flex items-center justify-center">
                      <span className="text-lg font-semibold text-indigo-400">
                        {owner.name.charAt(0).toUpperCase()}
                      </span>
                    </div>
                    <div>
                      <h3 className="text-lg font-medium text-white">{owner.name}</h3>
                      <p className="text-sm text-slate-400">
                        {owner.totalFunds} {owner.totalFunds === 1 ? 'fund' : 'funds'}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Holdings table */}
              <table className="min-w-full" aria-label={`Holdings for ${owner.name}`}>
                <thead>
                  <tr className="border-b border-slate-800/50">
                    <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                      Fund
                    </th>
                    <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase tracking-wider">
                      Units
                    </th>
                    <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase tracking-wider">
                      Ownership
                    </th>
                    <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-slate-400 uppercase tracking-wider">
                      Acquired
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/50">
                  {owner.holdings.map((holding) => (
                    <tr key={holding.fundId} className="hover:bg-slate-800/30 transition-colors">
                      <td className="whitespace-nowrap px-6 py-4">
                        <Link
                          to={`/funds/${holding.fundId}`}
                          className="font-medium text-indigo-400 hover:text-indigo-300 transition-colors"
                        >
                          {holding.fundName}
                        </Link>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-slate-300 tabular-nums">
                        {holding.units.toLocaleString()} / {holding.totalUnits.toLocaleString()}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-right text-sm">
                        <span className="inline-flex items-center gap-2">
                          <span className="text-slate-300">{holding.percentage.toFixed(2)}%</span>
                          <div className="w-16 h-1.5 bg-slate-700 rounded-full overflow-hidden">
                            <div
                              className="h-full bg-indigo-500 rounded-full transition-all duration-300"
                              style={{ width: `${Math.min(holding.percentage, 100)}%` }}
                            />
                          </div>
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-right text-sm text-slate-500">
                        {formatDate(holding.acquiredAt)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}
