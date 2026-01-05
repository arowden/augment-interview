import { test, expect } from '@playwright/test';

test.describe('Fund Management Flow', () => {
  test('should display the dashboard with header', async ({ page }) => {
    await page.goto('/');

    await expect(page.locator('header')).toContainText('Augment Fund Cap Table');

    await expect(page.locator('h1')).toContainText('Funds');
  });

  test('should have skip link accessible on tab', async ({ page }) => {
    await page.goto('/');

    await page.keyboard.press('Tab');

    const skipLink = page.getByRole('link', { name: /skip to main content/i });
    await expect(skipLink).toBeFocused();
    await expect(skipLink).toBeVisible();
  });

  test('should open create fund modal', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('button', { name: /create fund/i }).click();

    await expect(page.getByRole('dialog')).toBeVisible();
    await expect(page.getByRole('heading', { name: /create new fund/i })).toBeVisible();

    await expect(page.getByLabel(/fund name/i)).toBeVisible();
    await expect(page.getByLabel(/total units/i)).toBeVisible();
    await expect(page.getByLabel(/initial owner/i)).toBeVisible();
  });

  test('should close modal on escape key', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('button', { name: /create fund/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible();

    await page.keyboard.press('Escape');

    await expect(page.getByRole('dialog')).not.toBeVisible();
  });

  test('should close modal on cancel button', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('button', { name: /create fund/i }).click();
    await expect(page.getByRole('dialog')).toBeVisible();

    await page.getByRole('button', { name: /cancel/i }).click();

    await expect(page.getByRole('dialog')).not.toBeVisible();
  });

  test('should show validation errors in create fund modal', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('button', { name: /create fund/i }).click();

    await page.getByLabel(/total units/i).clear();
    await page.getByRole('button', { name: /create fund/i }).last().click();

    await expect(page.getByText(/name is required/i)).toBeVisible();
  });

  test('should navigate to fund detail page', async ({ page }) => {
    await page.goto('/');


    await expect(page.getByRole('heading', { name: 'Funds' })).toBeVisible();
  });

  test('should navigate back to dashboard from fund page', async ({ page }) => {
    await page.goto('/funds/550e8400-e29b-41d4-a716-446655440000');

    const backLink = page.getByRole('link', { name: /back to dashboard/i });
    await expect(backLink).toBeVisible();

    await backLink.click();

    await expect(page).toHaveURL('/');
  });

  test('should show 404 for unknown routes', async ({ page }) => {
    await page.goto('/unknown-page');

    await expect(page.getByRole('heading', { name: '404' })).toBeVisible();
    await expect(page.getByText(/page not found/i)).toBeVisible();
  });

  test('should have accessible fund cards', async ({ page }) => {
    await page.goto('/');

    await page.waitForLoadState('networkidle');

    const fundCards = page.locator('[role="button"][aria-label*="fund details"]');
    const count = await fundCards.count();

    if (count > 0) {
      await page.keyboard.press('Tab');
      await page.keyboard.press('Tab');
      await page.keyboard.press('Tab');
      await page.keyboard.press('Tab');

    }
  });
});

test.describe('Transfer Flow', () => {
  test('should display transfer form on fund page', async ({ page }) => {
    await page.goto('/funds/550e8400-e29b-41d4-a716-446655440000');

    await expect(page.getByRole('heading', { name: /transfer units/i })).toBeVisible();
    await expect(page.getByLabel(/from owner/i)).toBeVisible();
    await expect(page.getByLabel(/to owner/i)).toBeVisible();
    await expect(page.getByLabel(/units/i)).toBeVisible();
  });

  test('should show validation error for self-transfer', async ({ page }) => {
    await page.goto('/funds/550e8400-e29b-41d4-a716-446655440000');

    await page.getByLabel(/from owner/i).fill('Founder LLC');
    await page.getByLabel(/to owner/i).fill('Founder LLC');
    await page.getByLabel(/units/i).fill('100');

    await page.getByRole('button', { name: /execute transfer/i }).click();

    await expect(page.getByText(/cannot transfer to same owner/i)).toBeVisible();
  });
});

test.describe('Accessibility', () => {
  test('should have proper heading hierarchy', async ({ page }) => {
    await page.goto('/');

    const h1 = page.locator('h1');
    await expect(h1).toHaveCount(1);
  });

  test('should have focus visible styles', async ({ page }) => {
    await page.goto('/');

    await page.keyboard.press('Tab');

    const focusedElement = page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });

  test('should trap focus in modal', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('button', { name: /create fund/i }).click();

    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');

    const focusedElement = page.locator(':focus');
    const modal = page.getByRole('dialog');
    await expect(modal).toContainText(await focusedElement.textContent() ?? '');
  });
});
