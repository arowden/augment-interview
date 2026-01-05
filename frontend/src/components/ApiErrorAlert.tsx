import { parseApiError } from '../api/errors';

interface ApiErrorAlertProps {
  error: unknown;
  fallbackMessage?: string;
}

/**
 * Displays structured API errors with user-friendly messages.
 * Shows hints and request ID for debugging when available.
 */
export function ApiErrorAlert({
  error,
  fallbackMessage = 'An error occurred',
}: ApiErrorAlertProps) {
  if (!error) return null;

  const parsed = parseApiError(error);
  const hint = parsed.details?.hint as string | undefined;

  return (
    <div
      role="alert"
      className="rounded-lg bg-red-500/10 border border-red-500/20 p-4"
    >
      <div className="flex items-start gap-3">
        <svg className="w-5 h-5 text-red-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div className="flex-1">
          <p className="font-medium text-red-400">
            {parsed.message || fallbackMessage}
          </p>
          {hint && (
            <p className="mt-1.5 text-sm text-slate-400">
              {hint}
            </p>
          )}
          {parsed.requestId && (
            <p className="mt-2 text-xs text-slate-500 font-mono">
              Request ID: {parsed.requestId}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
