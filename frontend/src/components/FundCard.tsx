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
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });

  return (
    <div
      role="button"
      tabIndex={0}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      className="group card-interactive p-5 focus-ring"
      aria-label={`View ${fund.name} details`}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <h3 className="font-medium text-white truncate group-hover:text-indigo-400 transition-colors">
            {fund.name}
          </h3>
          <p className="mt-1 text-sm text-slate-500">
            {formattedDate}
          </p>
        </div>
        <div className="w-10 h-10 rounded-lg bg-slate-800 flex items-center justify-center flex-shrink-0 group-hover:bg-indigo-500/10 transition-colors">
          <svg className="w-5 h-5 text-slate-500 group-hover:text-indigo-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>

      <div className="mt-4 pt-4 border-t border-slate-800/50">
        <div className="flex items-center justify-between">
          <span className="text-sm text-slate-500">Units</span>
          <span className="text-lg font-semibold text-white tabular-nums">
            {fund.totalUnits.toLocaleString()}
          </span>
        </div>
      </div>
    </div>
  );
}
