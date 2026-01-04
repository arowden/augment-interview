## ADDED Requirements

### Requirement: Project Structure
The frontend SHALL be a Vite-based React TypeScript project in the /frontend directory.

#### Scenario: Project initialization
- **WHEN** the frontend directory is examined
- **THEN** it contains package.json, tsconfig.json, vite.config.ts, and src/ directory

#### Scenario: TypeScript strict mode
- **WHEN** tsconfig.json is examined
- **THEN** strict mode is enabled

#### Scenario: Build command
- **WHEN** npm run build is executed
- **THEN** a production bundle is created in dist/

### Requirement: Generated API Client
The frontend SHALL use a TypeScript API client generated from the OpenAPI spec.

#### Scenario: API generation script
- **WHEN** npm run generate-api is executed
- **THEN** src/api/generated/ is populated with typed services

#### Scenario: Type safety
- **WHEN** API methods are called
- **THEN** request and response types are enforced by TypeScript

#### Scenario: Regeneration
- **WHEN** the OpenAPI spec changes and generate-api runs
- **THEN** the generated client reflects the new spec

### Requirement: API Client Configuration
The frontend SHALL provide a configured API client instance with base URL support.

#### Scenario: Base URL configuration
- **WHEN** VITE_API_URL is set to "http://api.example.com"
- **THEN** all API requests use that base URL

#### Scenario: Default base URL
- **WHEN** VITE_API_URL is not set
- **THEN** API requests use relative URLs (same origin)

### Requirement: TanStack Query Integration with Configured Defaults
The frontend SHALL use TanStack Query with configured retry, stale time, and garbage collection.

#### Scenario: Query provider
- **WHEN** the app mounts
- **THEN** QueryClientProvider wraps the application with configured queryClient

#### Scenario: Query retry configuration
- **WHEN** a query fails
- **THEN** it retries once before showing error (retry: 1)

#### Scenario: Stale time configuration
- **WHEN** data is fetched
- **THEN** it remains fresh for 5 minutes (staleTime: 300000)

#### Scenario: Garbage collection time
- **WHEN** a query is no longer used
- **THEN** cached data is garbage collected after 10 minutes (gcTime: 600000)

#### Scenario: Cache invalidation
- **WHEN** a mutation succeeds
- **THEN** related queries are invalidated and refetched

#### Scenario: Mutation retry
- **WHEN** a mutation fails
- **THEN** it does not automatically retry (retry: 0)

### Requirement: Fund Query Hooks
The frontend SHALL provide hooks for fund operations.

#### Scenario: useFunds hook
- **WHEN** useFunds is called
- **THEN** it returns data, isLoading, and error for all funds

#### Scenario: useFund hook
- **WHEN** useFund(id) is called
- **THEN** it returns data for a single fund

#### Scenario: useCreateFund mutation
- **WHEN** useCreateFund().mutate is called with fund data
- **THEN** a new fund is created and fund list is refetched

### Requirement: Cap Table Query Hook
The frontend SHALL provide a hook for cap table retrieval.

#### Scenario: useCapTable hook
- **WHEN** useCapTable(fundId) is called
- **THEN** it returns the cap table entries for that fund

### Requirement: Transfer Hooks
The frontend SHALL provide hooks for transfer operations.

#### Scenario: useTransfers hook
- **WHEN** useTransfers(fundId) is called
- **THEN** it returns transfer history for that fund

#### Scenario: useCreateTransfer mutation
- **WHEN** useCreateTransfer().mutate is called with transfer data
- **THEN** the transfer is executed and cap table/history are refetched

#### Scenario: Optimistic update for transfers
- **WHEN** useCreateTransfer().mutate is called
- **THEN** cap table UI updates optimistically before server response

#### Scenario: Optimistic rollback on error
- **WHEN** transfer mutation fails
- **THEN** cap table is rolled back to previous state

### Requirement: Application Routing
The frontend SHALL use react-router-dom for client-side routing.

#### Scenario: Dashboard route
- **WHEN** user navigates to /
- **THEN** the Dashboard page is rendered

#### Scenario: Fund detail route
- **WHEN** user navigates to /funds/:id
- **THEN** the FundPage is rendered with that fund's data

#### Scenario: 404 handling
- **WHEN** user navigates to unknown route
- **THEN** a not found message is displayed

### Requirement: Dashboard Page
The frontend SHALL display a dashboard with all funds and fund creation.

#### Scenario: Fund list display
- **WHEN** dashboard loads
- **THEN** all funds are displayed as cards

#### Scenario: Empty state
- **WHEN** no funds exist
- **THEN** a message prompts user to create first fund

#### Scenario: Create fund button
- **WHEN** user clicks "Create Fund"
- **THEN** a modal opens with fund creation form

#### Scenario: Fund navigation
- **WHEN** user clicks a fund card
- **THEN** they navigate to /funds/:id

### Requirement: Fund Creation Modal
The frontend SHALL provide a modal for creating new funds.

#### Scenario: Form fields
- **WHEN** create fund modal opens
- **THEN** it displays inputs for name, totalUnits, and initialOwner

#### Scenario: Validation
- **WHEN** user submits with empty name
- **THEN** an error message is shown

#### Scenario: Successful creation
- **WHEN** user submits valid fund data
- **THEN** modal closes, fund list updates, success message shown

### Requirement: Fund Detail Page
The frontend SHALL display fund details with cap table, transfer form, and history.

#### Scenario: Fund info display
- **WHEN** fund page loads
- **THEN** fund name and total units are displayed

#### Scenario: Cap table display
- **WHEN** fund page loads
- **THEN** ownership entries are shown in a table

#### Scenario: Transfer form display
- **WHEN** fund page loads
- **THEN** a form for executing transfers is visible

#### Scenario: Transfer history display
- **WHEN** fund page loads
- **THEN** past transfers are listed chronologically

### Requirement: Cap Table Component
The frontend SHALL render cap table entries as a sortable table.

#### Scenario: Table columns
- **WHEN** cap table renders
- **THEN** columns show Owner, Units, and percentage of total

#### Scenario: Percentage calculation
- **WHEN** an owner has 300 of 1000 total units
- **THEN** 30% is displayed

#### Scenario: Empty cap table
- **WHEN** cap table has no entries (impossible in practice)
- **THEN** an empty state message is shown

### Requirement: Transfer Form Component
The frontend SHALL provide a form for executing transfers.

#### Scenario: Form fields
- **WHEN** transfer form renders
- **THEN** it has inputs for fromOwner, toOwner, and units

#### Scenario: Submit validation
- **WHEN** user submits with units = 0
- **THEN** client-side error is shown

#### Scenario: Successful transfer
- **WHEN** user submits valid transfer
- **THEN** form clears, cap table updates, history updates

#### Scenario: Transfer error
- **WHEN** transfer fails (insufficient units)
- **THEN** error message from API is displayed

### Requirement: Transfer History Component
The frontend SHALL display transfer history for a fund.

#### Scenario: History list
- **WHEN** transfers exist
- **THEN** they are listed with from, to, units, and date

#### Scenario: Chronological order
- **WHEN** transfer history renders
- **THEN** most recent transfers appear first

#### Scenario: Empty history
- **WHEN** no transfers exist
- **THEN** "No transfers yet" message is shown

### Requirement: Layout Component
The frontend SHALL have a consistent layout with header and navigation.

#### Scenario: Header content
- **WHEN** any page renders
- **THEN** header shows "Augment Fund Cap Table"

#### Scenario: Navigation
- **WHEN** user is on fund detail page
- **THEN** a link to return to dashboard is available

### Requirement: Loading States
The frontend SHALL display loading indicators during data fetching.

#### Scenario: Initial load
- **WHEN** data is being fetched
- **THEN** a loading spinner or skeleton is shown

#### Scenario: Mutation pending
- **WHEN** a mutation is in progress
- **THEN** submit button is disabled with loading indicator

### Requirement: Error Handling
The frontend SHALL display user-friendly error messages.

#### Scenario: Network error
- **WHEN** API request fails due to network
- **THEN** "Unable to connect to server" message is shown

#### Scenario: API error
- **WHEN** API returns 400 with error message
- **THEN** that message is displayed to user

#### Scenario: Retry option
- **WHEN** a query fails
- **THEN** a retry button is available

### Requirement: Tailwind CSS Styling
The frontend SHALL use Tailwind CSS for styling.

#### Scenario: Tailwind configuration
- **WHEN** tailwind.config.js is examined
- **THEN** it includes content paths for all components

#### Scenario: Responsive design
- **WHEN** viewed on mobile viewport
- **THEN** layout adapts appropriately

### Requirement: Docker Build
The frontend SHALL be buildable as a Docker image serving static files.

#### Scenario: Dockerfile exists
- **WHEN** frontend/Dockerfile is examined
- **THEN** it builds the app and serves via nginx

#### Scenario: Production build
- **WHEN** docker build is run
- **THEN** an optimized production bundle is created

### Requirement: Error Boundaries
The frontend SHALL use React Error Boundaries to handle rendering errors gracefully.

#### Scenario: ErrorBoundary component exists
- **WHEN** src/components/ErrorBoundary.tsx is examined
- **THEN** it contains a class component implementing componentDidCatch

#### Scenario: Error boundary wraps routes
- **WHEN** App.tsx is examined
- **THEN** routes are wrapped with ErrorBoundary

#### Scenario: Error recovery UI
- **WHEN** a component throws during render
- **THEN** ErrorBoundary shows error message with "Try again" button

#### Scenario: Error logging
- **WHEN** ErrorBoundary catches an error
- **THEN** it logs error and component stack to console

### Requirement: Form Validation with Zod
The frontend SHALL validate forms using Zod schemas with React Hook Form.

#### Scenario: Transfer form schema
- **WHEN** src/schemas/transfer.ts is examined
- **THEN** it exports a Zod schema for fromOwner, toOwner, and units

#### Scenario: Self-transfer validation
- **WHEN** user enters same owner for from and to
- **THEN** validation error "Cannot transfer to same owner" is shown

#### Scenario: Required field validation
- **WHEN** user submits with empty fromOwner
- **THEN** validation error "From owner is required" is shown

#### Scenario: Minimum units validation
- **WHEN** user enters units < 1
- **THEN** validation error "Units must be at least 1" is shown

#### Scenario: Type inference
- **WHEN** TransferInput type is used
- **THEN** it is inferred from Zod schema (z.infer)

### Requirement: Accessibility (WCAG 2.1 AA)
The frontend SHALL comply with WCAG 2.1 AA accessibility guidelines.

#### Scenario: Skip link
- **WHEN** user presses Tab on page load
- **THEN** "Skip to main content" link becomes visible and focusable

#### Scenario: Form labels
- **WHEN** form inputs are examined
- **THEN** each input has an associated label with htmlFor/id

#### Scenario: Error announcements
- **WHEN** form validation error occurs
- **THEN** error has role="alert" for screen reader announcement

#### Scenario: Invalid state indication
- **WHEN** form field has validation error
- **THEN** input has aria-invalid="true" and aria-describedby pointing to error

#### Scenario: Table semantics
- **WHEN** cap table renders
- **THEN** it uses proper th scope="col" and aria-label

#### Scenario: Modal focus trap
- **WHEN** modal opens
- **THEN** focus is trapped within modal until closed

#### Scenario: Keyboard navigation
- **WHEN** user navigates with Tab key
- **THEN** all interactive elements are reachable

#### Scenario: Color contrast
- **WHEN** text elements are examined
- **THEN** contrast ratio meets 4.5:1 for normal text, 3:1 for large text

### Requirement: Testing Infrastructure
The frontend SHALL have testing infrastructure with RTL, MSW, and Playwright.

#### Scenario: RTL tests exist
- **WHEN** src/__tests__/ is examined
- **THEN** component tests using @testing-library/react exist

#### Scenario: MSW handlers exist
- **WHEN** src/__mocks__/handlers.ts is examined
- **THEN** it contains API mock handlers for tests

#### Scenario: E2E test configuration
- **WHEN** playwright.config.ts is examined
- **THEN** it configures E2E tests in e2e/ directory

#### Scenario: npm test script
- **WHEN** package.json scripts are examined
- **THEN** test script runs vitest or jest
