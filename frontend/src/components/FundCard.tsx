import { useNavigate } from 'react-router-dom';
import type { Fund } from '../api/client';

interface FundCardProps {
  fund: Fund;
}

export function FundCard({ fund }: FundCardProps) {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate(`/funds/${fund.id}`);
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      navigate(`/funds/${fund.id}`);
    }
  };

  const formattedDate = new Date(fund.createdAt).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  return (
    <div
      role="button"
      tabIndex={0}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      className="group cursor-pointer glass-card p-6 hover:border-primary-500/30 focus-ring"
      aria-label={`View ${fund.name} fund details`}
    >
      {/* Card header */}
      <div className="flex items-start justify-between">
        <div className="flex-1 min-w-0">
          <h3 className="text-lg font-heading font-semibold text-slate-100 truncate group-hover:text-primary-400 transition-colors">
            {fund.name}
          </h3>
          <p className="mt-1 text-sm text-slate-500">
            Created {formattedDate}
          </p>
        </div>
        <div className="ml-4 flex-shrink-0">
          <div className="w-10 h-10 rounded-xl bg-primary-500/10 flex items-center justify-center group-hover:bg-primary-500/20 transition-colors">
            <svg className="w-5 h-5 text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="mt-6 pt-4 border-t border-white/[0.06]">
        <div className="flex items-baseline justify-between">
          <span className="text-sm text-slate-400">Total Units</span>
          <span className="text-xl font-heading font-semibold text-slate-200">
            {fund.totalUnits.toLocaleString()}
          </span>
        </div>
      </div>

      {/* Decorative gradient line */}
      <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-transparent via-primary-500/50 to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded-b-2xl" />
    </div>
  );
}
