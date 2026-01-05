import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import { useFunds, useCreateFund } from '../hooks/useFunds';
import { FundList } from '../components/FundList';
import { CreateFundModal } from '../components/CreateFundModal';
import { ResetConfirmModal } from '../components/ResetConfirmModal';
import { AdminService } from '../api/generated';
import type { CreateFundInput } from '../schemas/fund';

export function Dashboard() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isResetModalOpen, setIsResetModalOpen] = useState(false);
  const queryClient = useQueryClient();
  const { data: fundList, isLoading, error, refetch } = useFunds();
  const createFund = useCreateFund();

  const resetMutation = useMutation({
    mutationFn: () => AdminService.resetDatabase(),
    onSuccess: () => {
      queryClient.invalidateQueries();
      setIsResetModalOpen(false);
    },
  });

  const handleCreateFund = async (data: CreateFundInput) => {
    await createFund.mutateAsync({
      name: data.name,
      totalUnits: data.totalUnits,
      initialOwner: data.initialOwner,
    });
    setIsModalOpen(false);
  };

  const handleReset = async () => {
    await resetMutation.mutateAsync();
  };

  const totalFunds = fundList?.funds?.length ?? 0;
  const totalUnits = fundList?.funds?.reduce((sum, f) => sum + f.totalUnits, 0) ?? 0;

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Page header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-white">
            Dashboard
          </h1>
          <p className="mt-1 text-slate-400">
            Manage and track your cap table
          </p>
        </div>
        <div className="flex gap-3">
          <button
            type="button"
            onClick={() => setIsResetModalOpen(true)}
            className="btn-secondary"
            title="Reset all data"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Reset
          </button>
          <button
            type="button"
            onClick={() => setIsModalOpen(true)}
            className="btn-primary"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
            New Fund
          </button>
        </div>
      </div>

      {/* Stats cards */}
      {!isLoading && !error && fundList?.funds && fundList.funds.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-3">
          <div className="card p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-indigo-500/10 flex items-center justify-center">
                <svg className="w-6 h-6 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z" />
                </svg>
              </div>
              <div>
                <p className="stat-label">Total Funds</p>
                <p className="stat-value">{totalFunds}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center">
                <svg className="w-6 h-6 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5m.75-9l3-3 2.148 2.148A12.061 12.061 0 0116.5 7.605" />
                </svg>
              </div>
              <div>
                <p className="stat-label">Total Units</p>
                <p className="stat-value">{totalUnits.toLocaleString()}</p>
              </div>
            </div>
          </div>

          <div className="card p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-xl bg-violet-500/10 flex items-center justify-center">
                <svg className="w-6 h-6 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 14.25v2.25m3-4.5v4.5m3-6.75v6.75m3-9v9M6 20.25h12A2.25 2.25 0 0020.25 18V6A2.25 2.25 0 0018 3.75H6A2.25 2.25 0 003.75 6v12A2.25 2.25 0 006 20.25z" />
                </svg>
              </div>
              <div>
                <p className="stat-label">Avg Units/Fund</p>
                <p className="stat-value">
                  {totalFunds > 0 ? Math.round(totalUnits / totalFunds).toLocaleString() : '0'}
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Fund list */}
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-medium text-white">
            Funds
          </h2>
          {fundList?.funds && fundList.funds.length > 0 && (
            <span className="badge-indigo">
              {fundList.funds.length} {fundList.funds.length === 1 ? 'fund' : 'funds'}
            </span>
          )}
        </div>
        <FundList
          funds={fundList?.funds}
          isLoading={isLoading}
          error={error}
          onRetry={() => refetch()}
        />
      </section>

      <CreateFundModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleCreateFund}
        isLoading={createFund.isPending}
        error={createFund.error}
      />

      <ResetConfirmModal
        isOpen={isResetModalOpen}
        onClose={() => setIsResetModalOpen(false)}
        onConfirm={handleReset}
        isLoading={resetMutation.isPending}
      />
    </div>
  );
}
