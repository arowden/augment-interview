import { OpenAPI } from './generated';

// Configure API client with base URL from environment.
export function configureApiClient(): void {
  const apiUrl = import.meta.env.VITE_API_URL;
  if (apiUrl) {
    OpenAPI.BASE = apiUrl;
  }
  // If VITE_API_URL is not set, keep the default '/api' for same-origin requests.
}

// Initialize on module load.
configureApiClient();

// Re-export everything from generated client for convenience.
export * from './generated';
