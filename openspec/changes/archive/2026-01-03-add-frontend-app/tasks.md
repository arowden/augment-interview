## 1. Project Setup
- [x] 1.1 Initialize Vite project with React and TypeScript template
- [x] 1.2 Configure tsconfig.json with strict mode
- [x] 1.3 Install and configure Tailwind CSS
- [x] 1.4 Install react-router-dom for routing
- [x] 1.5 Install @tanstack/react-query for data fetching
- [x] 1.6 Install zod for schema validation
- [x] 1.7 Install react-hook-form and @hookform/resolvers

## 2. QueryClient Configuration
- [x] 2.1 Create src/queryClient.ts with configured QueryClient
- [x] 2.2 Set retry: 1 for queries
- [x] 2.3 Set staleTime: 5 minutes
- [x] 2.4 Set gcTime: 10 minutes
- [x] 2.5 Set retry: 0 for mutations

## 3. API Client Generation
- [x] 3.1 Add openapi-typescript-codegen as dev dependency
- [x] 3.2 Create npm script generate-api pointing to ../api/openapi.yaml
- [x] 3.3 Generate client to src/api/generated/
- [x] 3.4 Create src/api/client.ts with configured API instance

## 4. Form Validation Schemas
- [x] 4.1 Create src/schemas/fund.ts with Zod schema
- [x] 4.2 Add name, totalUnits, initialOwner validation
- [x] 4.3 Create src/schemas/transfer.ts with Zod schema
- [x] 4.4 Add fromOwner, toOwner, units validation
- [x] 4.5 Add refine rule for self-transfer prevention
- [x] 4.6 Export TypeScript types from schemas

## 5. Error Boundary
- [x] 5.1 Create src/components/ErrorBoundary.tsx
- [x] 5.2 Implement componentDidCatch for error logging
- [x] 5.3 Implement recovery UI with "Try again" button
- [x] 5.4 Add role="alert" and aria-live="assertive"

## 6. Query Hooks with Optimistic Updates
- [x] 6.1 Create src/hooks/useFunds.ts with fund query hooks
- [x] 6.2 Implement useFunds for listing
- [x] 6.3 Implement useFund for single fund
- [x] 6.4 Implement useCreateFund mutation with cache invalidation
- [x] 6.5 Create src/hooks/useCapTable.ts
- [x] 6.6 Create src/hooks/useTransfers.ts with transfer hooks
- [x] 6.7 Implement useCreateTransfer with optimistic update
- [x] 6.8 Implement onMutate to save previous state
- [x] 6.9 Implement onError to rollback
- [x] 6.10 Implement onSettled to invalidate queries

## 7. Layout and Navigation with Accessibility
- [x] 7.1 Create src/components/Layout.tsx with skip link
- [x] 7.2 Add "Skip to main content" link (sr-only, focus:visible)
- [x] 7.3 Add id="main-content" to main element
- [x] 7.4 Create src/App.tsx with router and ErrorBoundary
- [x] 7.5 Define routes: /, /funds/:id
- [x] 7.6 Wrap routes with ErrorBoundary

## 8. Dashboard Page
- [x] 8.1 Create src/pages/Dashboard.tsx
- [x] 8.2 Create src/components/FundList.tsx
- [x] 8.3 Create src/components/FundCard.tsx (keyboard accessible)
- [x] 8.4 Create src/components/CreateFundModal.tsx with focus trap
- [x] 8.5 Add aria-modal="true" and escape key handling
- [x] 8.6 Implement fund creation with Zod validation

## 9. Fund Detail Page
- [x] 9.1 Create src/pages/FundPage.tsx
- [x] 9.2 Create src/components/CapTable.tsx with proper table semantics
- [x] 9.3 Add th scope="col" and aria-label
- [x] 9.4 Create src/components/TransferForm.tsx with React Hook Form
- [x] 9.5 Integrate Zod schema with zodResolver
- [x] 9.6 Add accessible error messages (role="alert", aria-describedby)
- [x] 9.7 Add aria-invalid to inputs with errors
- [x] 9.8 Create src/components/TransferHistory.tsx

## 10. Styling and Accessibility
- [x] 10.1 Create consistent color scheme with WCAG AA contrast
- [x] 10.2 Add responsive design for mobile
- [x] 10.3 Add loading states and skeletons
- [x] 10.4 Add error states and messages with role="alert"
- [x] 10.5 Verify all interactive elements are keyboard accessible
- [x] 10.6 Add focus visible styles

## 11. Testing Infrastructure
- [x] 11.1 Install @testing-library/react and @testing-library/user-event
- [x] 11.2 Install msw (Mock Service Worker)
- [x] 11.3 Create src/__mocks__/handlers.ts with API mocks
- [x] 11.4 Create src/__mocks__/server.ts with MSW setup
- [x] 11.5 Create src/__tests__/FundList.test.tsx
- [x] 11.6 Create src/__tests__/TransferForm.test.tsx
- [x] 11.7 Install @playwright/test for E2E
- [x] 11.8 Create playwright.config.ts
- [x] 11.9 Create e2e/fund-flow.spec.ts

## 12. Build Configuration
- [x] 12.1 Create Dockerfile for production build with nginx
- [x] 12.2 Configure vite.config.ts for API proxy in dev
- [x] 12.3 Configure environment variable handling
