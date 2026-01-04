import { useState } from 'react';

import { useFunds, useCreateFund } from '../hooks/useFunds';
import { FundList } from '../components/FundList';
import { CreateFundModal } from '../components/CreateFundModal';
import type { CreateFundInput } from '../schemas/fund';

export function Dashboard() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { data: fundList, isLoading, error, refetch } = useFunds();
  const createFund = useCreateFund();

  const handleCreateFund = async (data: CreateFundInput) => {
    await createFund.mutateAsync({
      name: data.name,
      totalUnits: data.totalUnits,
      initialOwner: data.initialOwner,
    });
    setIsModalOpen(false);
  };

  const totalFunds = fundList?.funds?.length ?? 0;
  const totalUnits = fundList?.funds?.reduce((sum, f) => sum + f.totalUnits, 0) ?? 0;

  return (
    <div className="space-y-8">
      {/* Page header */}
      <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-heading font-bold text-slate-100">
            Fund Dashboard
          </h1>
          <p className="mt-2 text-slate-400">
            Manage your investment funds and track ownership distribution.
          </p>
        </div>
        <button
          type="button"
          onClick={() => setIsModalOpen(true)}
          className="btn-primary"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          Create Fund
        </button>
      </div>

      {/* Stats overview */}
      {!isLoading && !error && fundList?.funds && fundList.funds.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div className="glass-card p-6">
            <p className="stat-label">Total Funds</p>
            <p className="stat-value mt-2">{totalFunds}</p>
          </div>
          <div className="glass-card p-6">
            <p className="stat-label">Total Units</p>
            <p className="stat-value mt-2">{totalUnits.toLocaleString()}</p>
          </div>
          <div className="glass-card p-6 sm:col-span-2 lg:col-span-1">
            <p className="stat-label">Avg Units per Fund</p>
            <p className="stat-value mt-2">
              {totalFunds > 0 ? Math.round(totalUnits / totalFunds).toLocaleString() : '0'}
            </p>
          </div>
        </div>
      )}

      {/* Fund list */}
      <section>
        <h2 className="text-xl font-heading font-semibold text-slate-200 mb-4">
          Your Funds
        </h2>
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
    </div>
  );
}
