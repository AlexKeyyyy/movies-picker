import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Рейтинги', () => {
  test('TC-008: Выставление рейтинга фильму', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.goto('/');
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Inception');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);

    const stars = page.locator('.ant-rate-star');
    await stars.nth(3).click();

    await page.getByRole('menuitem', { name: 'Мои оценки' }).click();
    await page.waitForURL(/\/ratings/);

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    const ratingCards = page.locator('.ant-card');
    await expect(ratingCards.first()).toBeVisible({ timeout: 10000 });
    expect(await ratingCards.count()).toBeGreaterThan(0);
  });
});