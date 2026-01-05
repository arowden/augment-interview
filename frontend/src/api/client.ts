import { OpenAPI } from './generated';

export function configureApiClient(): void {
  const apiUrl = import.meta.env.VITE_API_URL;
  if (apiUrl) {
    OpenAPI.BASE = apiUrl;
  }
}

configureApiClient();

export * from './generated';
