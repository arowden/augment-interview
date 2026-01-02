## 1. Project Setup
- [ ] 1.1 Initialize Vite project with React and TypeScript template
- [ ] 1.2 Configure tsconfig.json with strict mode
- [ ] 1.3 Install and configure Tailwind CSS
- [ ] 1.4 Install react-router-dom for routing
- [ ] 1.5 Install @tanstack/react-query for data fetching
- [ ] 1.6 Install zod for schema validation
- [ ] 1.7 Install react-hook-form and @hookform/resolvers

## 2. QueryClient Configuration
- [ ] 2.1 Create src/queryClient.ts with configured QueryClient
- [ ] 2.2 Set retry: 1 for queries
- [ ] 2.3 Set staleTime: 5 minutes
- [ ] 2.4 Set gcTime: 10 minutes
- [ ] 2.5 Set retry: 0 for mutations

## 3. API Client Generation
- [ ] 3.1 Add openapi-typescript-codegen as dev dependency
- [ ] 3.2 Create npm script generate-api pointing to ../api/openapi.yaml
- [ ] 3.3 Generate client to src/api/generated/
- [ ] 3.4 Create src/api/client.ts with configured API instance

## 4. Form Validation Schemas
- [ ] 4.1 Create src/schemas/fund.ts with Zod schema
- [ ] 4.2 Add name, totalUnits, initialOwner validation
- [ ] 4.3 Create src/schemas/transfer.ts with Zod schema
- [ ] 4.4 Add fromOwner, toOwner, units validation
- [ ] 4.5 Add refine rule for self-transfer prevention
- [ ] 4.6 Export TypeScript types from schemas

## 5. Error Boundary
- [ ] 5.1 Create src/components/ErrorBoundary.tsx
- [ ] 5.2 Implement componentDidCatch for error logging
- [ ] 5.3 Implement recovery UI with "Try again" button
- [ ] 5.4 Add role="alert" and aria-live="assertive"

## 6. Query Hooks with Optimistic Updates
- [ ] 6.1 Create src/hooks/useFunds.ts with fund query hooks
- [ ] 6.2 Implement useFunds for listing
- [ ] 6.3 Implement useFund for single fund
- [ ] 6.4 Implement useCreateFund mutation with cache invalidation
- [ ] 6.5 Create src/hooks/useCapTable.ts
- [ ] 6.6 Create src/hooks/useTransfers.ts with transfer hooks
- [ ] 6.7 Implement useCreateTransfer with optimistic update
- [ ] 6.8 Implement onMutate to save previous state
- [ ] 6.9 Implement onError to rollback
- [ ] 6.10 Implement onSettled to invalidate queries

## 7. Layout and Navigation with Accessibility
- [ ] 7.1 Create src/components/Layout.tsx with skip link
- [ ] 7.2 Add "Skip to main content" link (sr-only, focus:visible)
- [ ] 7.3 Add id="main-content" to main element
- [ ] 7.4 Create src/App.tsx with router and ErrorBoundary
- [ ] 7.5 Define routes: /, /funds/:id
- [ ] 7.6 Wrap routes with ErrorBoundary

## 8. Dashboard Page
- [ ] 8.1 Create src/pages/Dashboard.tsx
- [ ] 8.2 Create src/components/FundList.tsx
- [ ] 8.3 Create src/components/FundCard.tsx (keyboard accessible)
- [ ] 8.4 Create src/components/CreateFundModal.tsx with focus trap
- [ ] 8.5 Add aria-modal="true" and escape key handling
- [ ] 8.6 Implement fund creation with Zod validation

## 9. Fund Detail Page
- [ ] 9.1 Create src/pages/FundPage.tsx
- [ ] 9.2 Create src/components/CapTable.tsx with proper table semantics
- [ ] 9.3 Add th scope="col" and aria-label
- [ ] 9.4 Create src/components/TransferForm.tsx with React Hook Form
- [ ] 9.5 Integrate Zod schema with zodResolver
- [ ] 9.6 Add accessible error messages (role="alert", aria-describedby)
- [ ] 9.7 Add aria-invalid to inputs with errors
- [ ] 9.8 Create src/components/TransferHistory.tsx

## 10. Styling and Accessibility
- [ ] 10.1 Create consistent color scheme with WCAG AA contrast
- [ ] 10.2 Add responsive design for mobile
- [ ] 10.3 Add loading states and skeletons
- [ ] 10.4 Add error states and messages with role="alert"
- [ ] 10.5 Verify all interactive elements are keyboard accessible
- [ ] 10.6 Add focus visible styles

## 11. Testing Infrastructure
- [ ] 11.1 Install @testing-library/react and @testing-library/user-event
- [ ] 11.2 Install msw (Mock Service Worker)
- [ ] 11.3 Create src/__mocks__/handlers.ts with API mocks
- [ ] 11.4 Create src/__mocks__/server.ts with MSW setup
- [ ] 11.5 Create src/__tests__/FundList.test.tsx
- [ ] 11.6 Create src/__tests__/TransferForm.test.tsx
- [ ] 11.7 Install @playwright/test for E2E
- [ ] 11.8 Create playwright.config.ts
- [ ] 11.9 Create e2e/fund-flow.spec.ts

## 12. Build Configuration
- [ ] 12.1 Create Dockerfile for production build with nginx
- [ ] 12.2 Configure vite.config.ts for API proxy in dev
- [ ] 12.3 Configure environment variable handling
