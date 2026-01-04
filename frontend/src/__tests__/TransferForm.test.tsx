import { describe, it, expect, vi } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { TransferForm } from '../components/TransferForm';

describe('TransferForm', () => {
  it('renders all form fields', () => {
    render(
      <TransferForm
        onSubmit={async () => {}}
        isLoading={false}
        error={null}
      />
    );

    expect(screen.getByLabelText(/from owner/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/to owner/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/units/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /execute transfer/i })).toBeInTheDocument();
  });

  it('shows validation error for empty from owner', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(
      <TransferForm
        onSubmit={onSubmit}
        isLoading={false}
        error={null}
      />
    );

    // Fill only partial fields.
    await user.type(screen.getByLabelText(/to owner/i), 'Investor A');
    await user.clear(screen.getByLabelText(/units/i));
    await user.type(screen.getByLabelText(/units/i), '100');

    await user.click(screen.getByRole('button', { name: /execute transfer/i }));

    await waitFor(() => {
      expect(screen.getByText(/from owner is required/i)).toBeInTheDocument();
    });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('shows validation error for self-transfer', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(
      <TransferForm
        onSubmit={onSubmit}
        isLoading={false}
        error={null}
      />
    );

    await user.type(screen.getByLabelText(/from owner/i), 'Founder LLC');
    await user.type(screen.getByLabelText(/to owner/i), 'Founder LLC');
    await user.clear(screen.getByLabelText(/units/i));
    await user.type(screen.getByLabelText(/units/i), '100');

    await user.click(screen.getByRole('button', { name: /execute transfer/i }));

    await waitFor(() => {
      expect(screen.getByText(/cannot transfer to same owner/i)).toBeInTheDocument();
    });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('shows validation error for empty units field', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    render(
      <TransferForm
        onSubmit={onSubmit}
        isLoading={false}
        error={null}
      />
    );

    await user.type(screen.getByLabelText(/from owner/i), 'Founder LLC');
    await user.type(screen.getByLabelText(/to owner/i), 'Investor A');

    // Clear units field entirely - results in empty/NaN.
    const unitsInput = screen.getByLabelText(/units/i);
    await user.clear(unitsInput);

    await user.click(screen.getByRole('button', { name: /execute transfer/i }));

    // Wait for validation - should show error for invalid number.
    await waitFor(() => {
      const unitsField = screen.getByLabelText(/units/i);
      expect(unitsField).toHaveAttribute('aria-invalid', 'true');
    });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('submits valid form data', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(
      <TransferForm
        onSubmit={onSubmit}
        isLoading={false}
        error={null}
      />
    );

    await user.type(screen.getByLabelText(/from owner/i), 'Founder LLC');
    await user.type(screen.getByLabelText(/to owner/i), 'Investor A');
    await user.clear(screen.getByLabelText(/units/i));
    await user.type(screen.getByLabelText(/units/i), '100000');

    await user.click(screen.getByRole('button', { name: /execute transfer/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        fromOwner: 'Founder LLC',
        toOwner: 'Investor A',
        units: 100000,
      });
    });
  });

  it('disables submit button when loading', () => {
    render(
      <TransferForm
        onSubmit={async () => {}}
        isLoading={true}
        error={null}
      />
    );

    const submitButton = screen.getByRole('button', { name: /execute transfer/i });
    expect(submitButton).toBeDisabled();
  });

  it('displays API error message', () => {
    render(
      <TransferForm
        onSubmit={async () => {}}
        isLoading={false}
        error={new Error('Insufficient units')}
      />
    );

    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('Insufficient units')).toBeInTheDocument();
  });

  it('has proper accessibility attributes for invalid inputs', async () => {
    const user = userEvent.setup();

    render(
      <TransferForm
        onSubmit={async () => {}}
        isLoading={false}
        error={null}
      />
    );

    await user.click(screen.getByRole('button', { name: /execute transfer/i }));

    await waitFor(() => {
      const fromOwnerInput = screen.getByLabelText(/from owner/i);
      expect(fromOwnerInput).toHaveAttribute('aria-invalid', 'true');
      expect(fromOwnerInput).toHaveAttribute('aria-describedby');
    });
  });
});
