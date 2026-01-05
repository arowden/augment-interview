export function FundCardSkeleton() {
  return (
    <div className="card p-5">
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <div className="skeleton h-5 w-3/4" />
          <div className="skeleton h-4 w-1/3 mt-2" />
        </div>
        <div className="skeleton w-10 h-10 rounded-lg" />
      </div>
      <div className="mt-4 pt-4 border-t border-slate-800/50">
        <div className="flex items-center justify-between">
          <div className="skeleton h-4 w-12" />
          <div className="skeleton h-6 w-20" />
        </div>
      </div>
    </div>
  );
}
