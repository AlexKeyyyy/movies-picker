import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Рейтинги', () => {
  const email = generateTestEmail();
  const password = 'Password123!';

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await registerUser(page, email, password);
    await loginUser(page, email, password);
    await page.close();
  });

  test('TC-008: Выставление рейтинга фильму', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.ant-list-items .ant-card');
    await page.getByPlaceholder('Поиск фильмов...').fill('Inception');
    await page.waitForTimeout(1000);
    await page.locator('.ant-list-items .ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);

    // Ставим 4 звезды
    const stars = page.locator('.ant-rate-star');
    await stars.nth(3).click();
    await expect(page.locator('.ant-message-success')).toBeVisible();

    // Переход через меню
    await page.getByRole('menuitem', { name: 'Мои оценки' }).click();
    await page.waitForURL(/\/ratings/);

    const ratingCards = page.locator('.ant-card');
    expect(await ratingCards.count()).toBeGreaterThan(0);
  });
});