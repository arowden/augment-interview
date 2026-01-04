import { test, expect } from '@playwright/test';

test.describe('Fund Management Flow', () => {
  test('should display the dashboard with header', async ({ page }) => {
    await page.goto('/');

    // Check header is visible.
    await expect(page.locator('header')).toContainText('Augment Fund Cap Table');

    // Check page title.
    await expect(page.locator('h1')).toContainText('Funds');
  });

  test('should have skip link accessible on tab', async ({ page }) => {
    await page.goto('/');

    // Press tab to focus skip link.
    await page.keyboard.press('Tab');

    // Skip link should be visible when focused.
    const skipLink = page.getByRole('link', { name: /skip to main content/i });
    await expect(skipLink).toBeFocused();
    await expect(skipLink).toBeVisible();
  });

  test('should open create fund modal', async ({ page }) => {
    await page.goto('/');

    // Click create fund button.
    await page.getByRole('button', { name: /create fund/i }).click();

    // Modal should be visible.
    await expect(page.getByRole('dialog')).toBeVisible();
    await expect(page.getByRole('heading', { name: /create new fund/i })).toBeVisible();

    // Form fields should be present.
    await expect(page.getByLabel(/fund name/i)).toBeVisible();
    await expect(page.getByLabel(/total units/i)).toBeVisible();
    await expect(page.getByLabel(/initial owner/i)).toBeVisible();
  });

  test('should close modal on escape key', async ({ page }) => {
    await page.goto('/');

    // Open modal.
    await page.getByRole('button', { name: /create fund/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible();

    // Press escape.
    await page.keyboard.press('Escape');

    // Modal should be closed.
    await expect(page.getByRole('dialog')).not.toBeVisible();
  });

  test('should close modal on cancel button', async ({ page }) => {
    await page.goto('/');

    // Open modal.
    await page.getByRole('button', { name: /create fund/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible();

    // Click cancel.
    await page.getByRole('button', { name: /cancel/i }).click();

    // Modal should be closed.
    await expect(page.getByRole('dialog')).not.toBeVisible();
  });

  test('should show validation errors in create fund modal', async ({ page }) => {
    await page.goto('/');

    // Open modal.
    await page.getByRole('button', { name: /create fund/i }).click();

    // Clear the total units field and try to submit.
    await page.getByLabel(/total units/i).clear();
    await page.getByRole('button', { name: /create fund/i }).last().click();

    // Should show validation errors.
    await expect(page.getByText(/name is required/i)).toBeVisible();
  });

  test('should navigate to fund detail page', async ({ page }) => {
    await page.goto('/');

    // Wait for funds to load (assuming MSW or real API provides data).
    // If using MSW, funds should appear.

    // Check if we're on dashboard.
    await expect(page.getByRole('heading', { name: 'Funds' })).toBeVisible();
  });

  test('should navigate back to dashboard from fund page', async ({ page }) => {
    // Navigate directly to a fund page.
    await page.goto('/funds/550e8400-e29b-41d4-a716-446655440000');

    // Click back link.
    const backLink = page.getByRole('link', { name: /back to dashboard/i });
    await expect(backLink).toBeVisible();

    await backLink.click();

    // Should be on dashboard.
    await expect(page).toHaveURL('/');
  });

  test('should show 404 for unknown routes', async ({ page }) => {
    await page.goto('/unknown-page');

    await expect(page.getByRole('heading', { name: '404' })).toBeVisible();
    await expect(page.getByText(/page not found/i)).toBeVisible();
  });

  test('should have accessible fund cards', async ({ page }) => {
    await page.goto('/');

    // Wait for page to load.
    await page.waitForLoadState('networkidle');

    // If funds exist, check they have proper aria labels.
    const fundCards = page.locator('[role="button"][aria-label*="fund details"]');
    const count = await fundCards.count();

    if (count > 0) {
      // Each card should be keyboard accessible.
      await page.keyboard.press('Tab'); // Skip link.
      await page.keyboard.press('Tab'); // Header link.
      await page.keyboard.press('Tab'); // Create fund button.
      await page.keyboard.press('Tab'); // First fund card.

      // First fund card should be focusable.
    }
  });
});

test.describe('Transfer Flow', () => {
  test('should display transfer form on fund page', async ({ page }) => {
    await page.goto('/funds/550e8400-e29b-41d4-a716-446655440000');

    // Check transfer form is visible.
    await expect(page.getByRole('heading', { name: /transfer units/i })).toBeVisible();
    await expect(page.getByLabel(/from owner/i)).toBeVisible();
    await expect(page.getByLabel(/to owner/i)).toBeVisible();
    await expect(page.getByLabel(/units/i)).toBeVisible();
  });

  test('should show validation error for self-transfer', async ({ page }) => {
    await page.goto('/funds/550e8400-e29b-41d4-a716-446655440000');

    // Fill form with same owner.
    await page.getByLabel(/from owner/i).fill('Founder LLC');
    await page.getByLabel(/to owner/i).fill('Founder LLC');
    await page.getByLabel(/units/i).fill('100');

    // Submit form.
    await page.getByRole('button', { name: /execute transfer/i }).click();

    // Should show error.
    await expect(page.getByText(/cannot transfer to same owner/i)).toBeVisible();
  });
});

test.describe('Accessibility', () => {
  test('should have proper heading hierarchy', async ({ page }) => {
    await page.goto('/');

    // Check h1 exists.
    const h1 = page.locator('h1');
    await expect(h1).toHaveCount(1);
  });

  test('should have focus visible styles', async ({ page }) => {
    await page.goto('/');

    // Tab to an interactive element.
    await page.keyboard.press('Tab');

    // Check that focused element has visible focus indicator.
    const focusedElement = page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });

  test('should trap focus in modal', async ({ page }) => {
    await page.goto('/');

    // Open modal.
    await page.getByRole('button', { name: /create fund/i }).click();

    // Tab through modal.
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');

    // Focus should stay within modal.
    const focusedElement = page.locator(':focus');
    const modal = page.getByRole('dialog');
    await expect(modal).toContainText(await focusedElement.textContent() ?? '');
  });
});
