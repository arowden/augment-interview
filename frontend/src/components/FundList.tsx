import { useNavigate } from 'react-router-dom';
import type { Fund } from '../api/client';

interface FundListProps {
  funds: Fund[] | undefined;
  isLoading: boolean;
  error: Error | null;
  onRetry: () => void;
}

export function FundList({ funds, isLoading, error, onRetry }: FundListProps) {
  const navigate = useNavigate();

  if (isLoading) {
    return (
      <div className="card overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-slate-800">
              <th className="text-left text-sm font-medium text-slate-400 px-6 py-4">Name</th>
              <th className="text-right text-sm font-medium text-slate-400 px-6 py-4">Units</th>
              <th className="text-right text-sm font-medium text-slate-400 px-6 py-4">Created</th>
              <th className="w-12"></th>
            </tr>
          </thead>
          <tbody>
            {Array.from({ length: 5 }).map((_, i) => (
              <tr key={i} className="border-b border-slate-800/50 last:border-0">
                <td className="px-6 py-4"><div className="skeleton h-5 w-40" /></td>
                <td className="px-6 py-4 text-right"><div className="skeleton h-5 w-24 ml-auto" /></td>
                <td className="px-6 py-4 text-right"><div className="skeleton h-5 w-28 ml-auto" /></td>
                <td className="px-6 py-4"><div className="skeleton h-5 w-5" /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  if (error) {
    return (
      <div role="alert" className="card p-8 text-center">
        <div className="w-12 h-12 rounded-full bg-red-500/10 flex items-center justify-center mx-auto">
          <svg className="w-6 h-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <p className="mt-4 font-medium text-white">Unable to load funds</p>
        <p className="mt-1 text-sm text-slate-500">
          {error.message || 'Something went wrong'}
        </p>
        <button type="button" onClick={onRetry} className="btn-primary mt-6">
          Try Again
        </button>
      </div>
    );
  }

  if (!funds || funds.length === 0) {
    return (
      <div className="card p-12 text-center border-dashed border-2 border-slate-800">
        <div className="w-14 h-14 rounded-2xl bg-slate-800 flex items-center justify-center mx-auto">
          <svg className="w-7 h-7 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
          </svg>
        </div>
        <h3 className="mt-4 font-medium text-white">No funds yet</h3>
        <p className="mt-1 text-sm text-slate-500">
          Create your first fund to get started
        </p>
      </div>
    );
  }

  return (
    <div className="card overflow-hidden">
      <table className="w-full">
        <thead>
          <tr className="border-b border-slate-800">
            <th className="text-left text-sm font-medium text-slate-400 px-6 py-4">Name</th>
            <th className="text-right text-sm font-medium text-slate-400 px-6 py-4">Units</th>
            <th className="text-right text-sm font-medium text-slate-400 px-6 py-4">Created</th>
            <th className="w-12"></th>
          </tr>
        </thead>
        <tbody>
          {funds.map((fund) => {
            const formattedDate = new Date(fund.createdAt).toLocaleDateString('en-US', {
              month: 'short',
              day: 'numeric',
              year: 'numeric',
            });

            return (
              <tr
                key={fund.id}
                onClick={() => navigate(`/funds/${fund.id}`)}
                className="border-b border-slate-800/50 last:border-0 cursor-pointer hover:bg-slate-800/30 transition-colors group"
              >
                <td className="px-6 py-4">
                  <span className="font-medium text-white group-hover:text-indigo-400 transition-colors">
                    {fund.name}
                  </span>
                </td>
                <td className="px-6 py-4 text-right">
                  <span className="text-slate-300 tabular-nums">
                    {fund.totalUnits.toLocaleString()}
                  </span>
                </td>
                <td className="px-6 py-4 text-right">
                  <span className="text-slate-500">{formattedDate}</span>
                </td>
                <td className="px-6 py-4">
                  <svg className="w-5 h-5 text-slate-600 group-hover:text-indigo-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                  </svg>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
