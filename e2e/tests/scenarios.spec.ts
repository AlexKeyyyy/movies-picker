import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Сквозные сценарии пользователя', () => {
  test('Сценарий 1: Полный цикл нового пользователя', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);

    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.getByRole('button', { name: 'Убрать из списка' })).toBeVisible();

    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);
    const watchlistCards = page.locator('.ant-card');
    expect(await watchlistCards.count()).toBeGreaterThan(0);
    await expect(watchlistCards.first()).toContainText(/Inception|Начало/i);

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 2: Оценка фильма после просмотра', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
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

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 3: Управление списком – добавить и удалить', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.getByRole('button', { name: 'Убрать из списка' })).toBeVisible();

    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);
    const firstCard = page.locator('.ant-card').first();
    await firstCard.getByRole('button', { name: 'Удалить' }).click();
    await expect(page.locator('.ant-card')).toHaveCount(0);

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 4: Доступ к профилю без авторизации, затем вход', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await page.goto('/profile/edit');
    await page.waitForURL(/.*\/login/);
    await expect(page.getByRole('button', { name: 'Войти' })).toBeVisible();

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.goto('/profile/edit');
    await page.waitForURL(/\/profile\/edit/);
    await expect(page.locator('.ant-card-head-title', { hasText: 'Редактировать профиль' })).toBeVisible();

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 5: Регистрация -> вход -> просмотр деталей фильма', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);

    const title = page.locator('h3');
    await expect(title).toContainText(/Inception|Начало/i);

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 6: Неверный пароль, затем успешный вход', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);

    await page.goto('/login');
    await page.getByLabel('Email').fill(email);
    await page.getByLabel('Пароль').fill('wrongpassword');
    await page.getByRole('button', { name: 'Войти' }).click();
    await expect(page.locator('.ant-message-error')).toBeVisible();
    await expect(page).toHaveURL(/.*\/login/);

    await page.getByLabel('Пароль').fill(password);
    await page.getByRole('button', { name: 'Войти' }).click();
    await page.waitForURL('http://localhost:3000/');
    await expect(page.getByRole('menuitem', { name: 'Выйти' })).toBeVisible();

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await expect(page.locator('.ant-card').first()).toBeVisible();

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 7: Добавление двух фильмов, удаление одного', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    // Добавление Inception
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.getByRole('button', { name: 'Убрать из списка' })).toBeVisible();

    // Добавление Interstellar
    await page.goto('/');
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Интерстеллар');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.getByRole('button', { name: 'Убрать из списка' })).toBeVisible();

    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);
    const cards = page.locator('.ant-card');
    await expect(cards).toHaveCount(2);

    // Удаляем первый (Inception)
    const firstCard = cards.first();
    await firstCard.getByRole('button', { name: 'Удалить' }).click();
    await expect(cards).toHaveCount(1);
    await expect(cards.first()).toContainText(/Interstellar|Интерстеллар/i);

    // Удаляем второй
    await cards.first().getByRole('button', { name: 'Удалить' }).click();
    await expect(page.locator('.ant-card')).toHaveCount(0);

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 8: Оценка двух фильмов и проверка в рейтингах', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    // Оценка Inception
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.locator('.ant-rate-star').nth(3).click();

    // Оценка Interstellar
    await page.goto('/');
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Интерстеллар');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.locator('.ant-rate-star').nth(4).click();

    await page.getByRole('menuitem', { name: 'Мои оценки' }).click();
    await page.waitForURL(/\/ratings/);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    const ratingCards = page.locator('.ant-card');
    await expect(ratingCards).toHaveCount(2);

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 9: Попытка регистрации с существующим email', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await page.goto('/login');
    await loginUser(page, email, password);
    await page.getByRole('menuitem', { name: 'Выйти' }).click();

    await page.goto('/register');
    await page.getByLabel('Email').fill(email);
    await page.getByLabel('Пароль').fill(password);
    await page.getByRole('button', { name: 'Зарегистрироваться' }).click();
    await expect(page.locator('.ant-message-error')).toBeVisible();

    const newEmail = generateTestEmail();
    await page.getByLabel('Email').fill(newEmail);
    await page.getByRole('button', { name: 'Зарегистрироваться' }).click();
    await page.waitForURL(/.*\/login/);
    await expect(page.locator('.ant-message-success')).toBeVisible();

    await loginUser(page, newEmail, password);
    await expect(page.getByRole('menuitem', { name: 'Выйти' })).toBeVisible();

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });

  test('Сценарий 10: Максимальный полный цикл (регистрация -> добавление -> оценка -> удаление -> рейтинг)', async ({ page }) => {
    const email = generateTestEmail();
    const password = 'Password123!';

    await registerUser(page, email, password);
    await loginUser(page, email, password);

    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().waitFor();
    await page.getByPlaceholder('Поиск фильмов...').fill('Начало');
    await page.waitForTimeout(1000);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    await page.locator('.ant-card').first().click();
    await page.waitForURL(/\/movies\/\d+/);
    await page.getByRole('button', { name: 'К просмотру' }).click();
    await expect(page.getByRole('button', { name: 'Убрать из списка' })).toBeVisible();
    await page.locator('.ant-rate-star').nth(3).click();

    await page.getByRole('menuitem', { name: 'К просмотру' }).click();
    await page.waitForURL(/\/watchlist/);
    const firstCard = page.locator('.ant-card').first();
    await firstCard.getByRole('button', { name: 'Удалить' }).click();
    await expect(page.locator('.ant-card')).toHaveCount(0);

    await page.getByRole('menuitem', { name: 'Мои оценки' }).click();
    await page.waitForURL(/\/ratings/);
    await page.waitForSelector('.ant-spin', { state: 'hidden' });
    const ratingCards = page.locator('.ant-card');
    await expect(ratingCards.first()).toBeVisible();
    expect(await ratingCards.count()).toBeGreaterThan(0);

    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
  });
});