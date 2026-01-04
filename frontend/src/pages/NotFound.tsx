import { Link } from 'react-router-dom';

export function NotFound() {
  return (
    <div className="flex min-h-[400px] items-center justify-center">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-slate-100">404</h1>
        <p className="mt-2 text-lg text-slate-300">Page not found</p>
        <p className="mt-1 text-sm text-slate-400">
          The page you're looking for doesn't exist.
        </p>
        <Link
          to="/"
          className="mt-6 inline-block rounded-md bg-primary-600 px-4 py-2 text-white hover:bg-primary-700 focus-visible-ring"
        >
          Back to Dashboard
        </Link>
      </div>
    </div>
  );
}
