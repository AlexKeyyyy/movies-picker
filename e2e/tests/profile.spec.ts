import { test, expect } from '@playwright/test';

test.describe('Профиль без авторизации', () => {
  test('TC-010: Доступ к странице профиля без авторизации', async ({ page }) => {
    await page.goto('/profile/edit');
    await page.waitForURL(/.*\/login/);
    await expect(page.locator('.ant-card-head-title', { hasText: 'Вход' })).toBeVisible();
  });
});