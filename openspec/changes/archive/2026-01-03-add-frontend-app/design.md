## Context
The frontend is a single-page application providing CRUD operations for funds and transfers. It consumes the Go API and must handle loading, error, and success states for all operations with proper accessibility.

## Goals / Non-Goals
- Goals: Usable UI, type-safe API integration, responsive design, clear feedback, accessibility compliance (WCAG 2.1 AA), error boundaries, form validation
- Non-Goals: Offline support, PWA features, animations, dark mode

## Decisions
- Decision: Vite for build tooling (fast HMR, modern defaults)
- Alternatives considered: Create React App (deprecated), Next.js (SSR overkill)

- Decision: TanStack Query for server state management with configured defaults
- Alternatives considered: SWR (less features), Redux (boilerplate), plain fetch (no caching)

- Decision: Tailwind CSS for styling (utility-first, fast development)
- Alternatives considered: CSS modules (more files), styled-components (runtime cost)

- Decision: React Hook Form + Zod for form validation
- Alternatives considered: Formik (larger bundle), plain React (no validation), Yup (less type-safe)

- Decision: openapi-typescript-codegen for client generation
- Alternatives considered: openapi-generator (Java), axios templates (no typing)

- Decision: React Error Boundaries for graceful error handling
- Alternatives considered: No error boundaries (crashes propagate), global try/catch (no UI recovery)

## Directory Structure
```
frontend/
  src/
    api/
      generated/          # Generated from OpenAPI (DO NOT EDIT)
      client.ts           # Configured API instance
    components/
      ErrorBoundary.tsx   # Error boundary wrapper
      Layout.tsx          # App shell with skip link
      FundList.tsx        # Fund list display
      FundCard.tsx        # Individual fund card
      CapTable.tsx        # Ownership table with a11y
      TransferForm.tsx    # Transfer form with validation
      TransferHistory.tsx # Transfer list
      CreateFundModal.tsx # Modal with focus trap
    hooks/
      useFunds.ts         # Fund query/mutation hooks
      useCapTable.ts      # Cap table query hook
      useTransfers.ts     # Transfer query/mutation hooks
    schemas/
      fund.ts             # Zod schema for fund form
      transfer.ts         # Zod schema for transfer form
    pages/
      Dashboard.tsx       # Main fund list page
      FundPage.tsx        # Fund detail page
    App.tsx               # Router and providers
    main.tsx              # Entry point
    queryClient.ts        # Configured QueryClient
    index.css             # Tailwind imports
```

## QueryClient Configuration
```typescript
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 5 * 60 * 1000,
      gcTime: 10 * 60 * 1000,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});
```

## Error Boundary Component
```typescript
import { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false };

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('ErrorBoundary caught:', error, info.componentStack);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback ?? (
        <div role="alert" aria-live="assertive">
          <h2>Something went wrong</h2>
          <button onClick={() => this.setState({ hasError: false })}>
            Try again
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}
```

## Form Validation with Zod
```typescript
import { z } from 'zod';

export const transferSchema = z.object({
  fromOwner: z.string().min(1, 'From owner is required').trim(),
  toOwner: z.string().min(1, 'To owner is required').trim(),
  units: z.number().min(1, 'Units must be at least 1'),
}).refine(
  (data) => data.fromOwner !== data.toOwner,
  { message: 'Cannot transfer to same owner', path: ['toOwner'] }
);

export type TransferInput = z.infer<typeof transferSchema>;
```

## Query Hook Pattern with Optimistic Updates
```typescript
export function useFunds() {
  return useQuery({
    queryKey: ['funds'],
    queryFn: () => FundsService.listFunds(),
  });
}

export function useCreateTransfer(fundId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: TransfersService.createTransfer,
    onMutate: async (newTransfer) => {
      await queryClient.cancelQueries({ queryKey: ['capTable', fundId] });
      const previousCapTable = queryClient.getQueryData(['capTable', fundId]);
      return { previousCapTable };
    },
    onError: (err, newTransfer, context) => {
      if (context?.previousCapTable) {
        queryClient.setQueryData(['capTable', fundId], context.previousCapTable);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['capTable', fundId] });
      queryClient.invalidateQueries({ queryKey: ['transfers', fundId] });
    },
  });
}
```

## Accessibility Requirements
```typescript
<table role="grid" aria-label="Cap table showing ownership">
  <thead>
    <tr>
      <th scope="col">Owner</th>
      <th scope="col">Units</th>
      <th scope="col">Percentage</th>
    </tr>
  </thead>
  <tbody>
    {entries.map((entry) => (
      <tr key={entry.id}>
        <td>{entry.ownerName}</td>
        <td>{entry.units.toLocaleString()}</td>
        <td aria-label={`${percentage}%`}>{percentage}%</td>
      </tr>
    ))}
  </tbody>
</table>

<form
  onSubmit={handleSubmit}
  aria-labelledby="transfer-form-heading"
>
  <h2 id="transfer-form-heading">Execute Transfer</h2>

  <label htmlFor="fromOwner">From Owner</label>
  <input
    id="fromOwner"
    aria-describedby={errors.fromOwner ? 'fromOwner-error' : undefined}
    aria-invalid={!!errors.fromOwner}
  />
  {errors.fromOwner && (
    <span id="fromOwner-error" role="alert">{errors.fromOwner.message}</span>
  )}
</form>
```

## Skip Link for Keyboard Navigation
```typescript
export function Layout({ children }) {
  return (
    <>
      <a
        href="#main-content"
        className="sr-only focus:not-sr-only focus:absolute focus:p-4"
      >
        Skip to main content
      </a>
      <header>...</header>
      <main id="main-content" tabIndex={-1}>
        {children}
      </main>
    </>
  );
}
```

## Testing Strategy
```
frontend/
  src/
    __tests__/
      FundList.test.tsx     # Component tests with RTL
      TransferForm.test.tsx # Form validation tests
    __mocks__/
      handlers.ts           # MSW handlers
      server.ts             # MSW setup
  e2e/
    fund-flow.spec.ts       # Playwright E2E tests
```

Dependencies:
- @testing-library/react
- @testing-library/user-event
- msw (Mock Service Worker)
- @playwright/test

## Component Responsibilities
- **ErrorBoundary**: Catches render errors, shows recovery UI
- **Layout**: Header, navigation, skip link, main content area
- **FundList**: Maps funds to FundCards, handles empty state
- **FundCard**: Displays fund summary, keyboard accessible, links to detail
- **CapTable**: Renders ownership entries with proper table semantics
- **TransferForm**: Zod validation, React Hook Form, accessible errors
- **TransferHistory**: Lists transfers chronologically
- **CreateFundModal**: Focus trap, escape key handling, aria-modal

## Risks / Trade-offs
- Generated client tied to spec version → Regenerate on spec changes
- Optimistic updates may show incorrect state briefly → Rollback on error
- Error boundaries don't catch async errors → Use query error states for those

## Open Questions
- None
