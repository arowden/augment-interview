import { parseApiError } from '../api/errors';

interface ApiErrorAlertProps {
  error: unknown;
  fallbackMessage?: string;
}

/**
 * Displays structured API errors with user-friendly messages.
 * Shows request ID for debugging when available.
 */
export function ApiErrorAlert({
  error,
  fallbackMessage = 'An error occurred',
}: ApiErrorAlertProps) {
  if (!error) return null;

  const parsed = parseApiError(error);

  return (
    <div
      role="alert"
      className="rounded-md bg-error-50 p-3 text-sm"
    >
      <p className="font-medium text-error-600">
        {parsed.message || fallbackMessage}
      </p>
      {parsed.requestId && (
        <p className="mt-1 text-xs text-error-500">
          Request ID: {parsed.requestId}
        </p>
      )}
    </div>
  );
}
