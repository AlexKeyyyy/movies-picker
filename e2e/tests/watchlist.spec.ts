import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Смотреть позже', () => {
  const email = generateTestEmail();
  const password = 'Password123!';

  test.beforeAll(async ({ browser }) => {
    const page = await browser.newPage();
    await registerUser(page, email, password);
    await loginUser(page, email, password);
    await page.close();
  });

  test('TC-005: Добавление фильма в список «Смотреть позже»', async ({ page }) => {
    await page.goto('/');
    await page.waitForSelector('.ant-list-items .ant-card');
    await page.getByPlaceholder('Поиск фильмов...').fill('Inception');
    await page.waitForTimeout(1000);
    await page.locator('.ant-list-items .ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.locator('.ant-message-success')).toBeVisible();

    // Переход через меню
    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);

    const watchlistCards = page.locator('.ant-card');
    expect(await watchlistCards.count()).toBeGreaterThan(0);
    await expect(watchlistCards.first()).toContainText('Inception');
  });

  test('TC-006: Удаление фильма из «Смотреть позже»', async ({ page }) => {
    // Добавляем фильм
    await page.goto('/');
    await page.waitForSelector('.ant-list-items .ant-card');
    await page.getByPlaceholder('Поиск фильмов...').fill('Inception');
    await page.waitForTimeout(1000);
    await page.locator('.ant-list-items .ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.locator('.ant-message-success')).toBeVisible();

    // Переход через меню
    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);

    // Удаляем
    const firstCard = page.locator('.ant-card').first();
    await firstCard.getByRole('button', { name: 'Удалить' }).click();
    await expect(page.locator('.ant-message-success')).toBeVisible();

    // Проверяем, что список пуст (или фильм исчез)
    await expect(page.locator('.ant-card')).toHaveCount(0);
  });
});