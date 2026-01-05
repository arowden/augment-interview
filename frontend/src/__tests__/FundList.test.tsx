import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';

import { FundList } from '../components/FundList';
import type { Fund } from '../api/client';

const mockFunds: Fund[] = [
  {
    id: '550e8400-e29b-41d4-a716-446655440000',
    name: 'Growth Fund I',
    totalUnits: 1000000,
    createdAt: '2024-01-15T10:30:00Z',
  },
  {
    id: '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
    name: 'Venture Fund II',
    totalUnits: 500000,
    createdAt: '2024-02-20T14:45:00Z',
  },
];

function renderWithRouter(ui: React.ReactElement) {
  return render(<BrowserRouter>{ui}</BrowserRouter>);
}

describe('FundList', () => {
  it('renders loading skeletons when loading', () => {
    renderWithRouter(
      <FundList
        funds={undefined}
        isLoading={true}
        error={null}
        onRetry={() => {}}
      />
    );

    const skeletons = document.querySelectorAll('.skeleton');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it('renders error state with retry button', async () => {
    const user = userEvent.setup();
    const onRetry = vi.fn();

    renderWithRouter(
      <FundList
        funds={undefined}
        isLoading={false}
        error={new Error('Failed to load')}
        onRetry={onRetry}
      />
    );

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Failed to load')).toBeInTheDocument();

    const retryButton = screen.getByRole('button', { name: /try again/i });
    await user.click(retryButton);

    expect(onRetry).toHaveBeenCalled();
  });

  it('renders empty state when no funds', () => {
    renderWithRouter(
      <FundList
        funds={[]}
        isLoading={false}
        error={null}
        onRetry={() => {}}
      />
    );

    expect(
      screen.getByText(/no funds yet/i)
    ).toBeInTheDocument();
  });

  it('renders fund cards when funds are provided', () => {
    renderWithRouter(
      <FundList
        funds={mockFunds}
        isLoading={false}
        error={null}
        onRetry={() => {}}
      />
    );

    expect(screen.getByText('Growth Fund I')).toBeInTheDocument();
    expect(screen.getByText('Venture Fund II')).toBeInTheDocument();
    expect(screen.getByText('1,000,000')).toBeInTheDocument();
    expect(screen.getByText('500,000')).toBeInTheDocument();
  });

  it('fund cards are keyboard accessible', async () => {
    const user = userEvent.setup();

    renderWithRouter(
      <FundList
        funds={mockFunds}
        isLoading={false}
        error={null}
        onRetry={() => {}}
      />
    );

    const fundCards = screen.getAllByRole('button');
    expect(fundCards.length).toBe(2);

    await user.tab();
    expect(fundCards[0]).toHaveFocus();

    await user.tab();
    expect(fundCards[1]).toHaveFocus();
  });
});
