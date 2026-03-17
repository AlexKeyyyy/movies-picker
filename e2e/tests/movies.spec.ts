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
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    const searchInput = page.getByPlaceholder('Поиск фильмов...');
    await searchInput.fill('Inception');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    const movieCards = page.locator('.ant-card');
    expect(await movieCards.count()).toBeGreaterThan(0);
  });

  test('TC-007: Просмотр деталей фильма', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Inception');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);

    // Проверяем, что заголовок фильма содержит название
    const title = page.locator('h3');
    await expect(title).toContainText(/Inception|Начало/i);
  });
});