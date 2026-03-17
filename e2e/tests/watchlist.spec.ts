import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Смотреть позже', () => {
  test('TC-005: Добавление фильма в список «Смотреть позже»', async ({ page }) => {
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

    await page.getByRole('button', { name: 'К просмотру' }).click();

    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);

    const watchlistCards = page.locator('.ant-card');
    expect(await watchlistCards.count()).toBeGreaterThan(0);
    await expect(watchlistCards.first()).toContainText('Начало');
  });

  test('TC-006: Удаление фильма из «Смотреть позже»', async ({ page }) => {
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

    await page.getByRole('button', { name: 'К просмотру' }).click();

    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);

    await expect(page.locator('.ant-card').first()).toBeVisible();

    const firstCard = page.locator('.ant-card').first();
    await firstCard.getByRole('button', { name: 'Удалить' }).click();

    await expect(page.locator('.ant-card')).toHaveCount(0, { timeout: 5000 });
  });
});