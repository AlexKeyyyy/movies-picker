import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Фильмы и обзоры', () => {
  const email = generateTestEmail();
  const password = 'Password123!';

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await registerUser(page, email, password);
    await loginUser(page, email, password);
    await page.close();
  });

  test('TC-004: Поиск фильма по названию', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.ant-list-items .ant-card');
    const searchInput = page.getByPlaceholder('Поиск фильмов...');
    await searchInput.fill('Inception');
    await page.waitForTimeout(1000); // debounce
    const movieCards = page.locator('.ant-list-items .ant-card');
    await expect(movieCards.first()).toBeVisible({ timeout: 10000 });
    expect(await movieCards.count()).toBeGreaterThan(0);
  });

  test('TC-007: Просмотр YouTube-обзоров к фильму', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.ant-list-items .ant-card');
    await page.getByPlaceholder('Поиск фильмов...').fill('Inception');
    await page.waitForTimeout(1000);
    await page.locator('.ant-list-items .ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    const reviewsHeader = page.getByRole('heading', { name: 'Обзоры с YouTube' });
    await expect(reviewsHeader).toBeVisible();
    const reviewCards = page.locator('h4:has-text("Обзоры с YouTube") + div .ant-card');
    expect(await reviewCards.count()).toBeGreaterThan(0);
  });
});