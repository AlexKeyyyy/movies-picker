import { test, expect } from '@playwright/test';
import { generateTestEmail, registerUser, loginUser } from './helpers';

test.describe('Авторизация', () => {
  const password = 'Password123!';
  let email: string;

  test.beforeEach(() => {
    email = generateTestEmail();
  });

  test('TC-001: Регистрация нового пользователя', async ({ page }) => {
    await registerUser(page, email, password);
    await expect(page).toHaveURL(/.*\/login/);
    const successMessage = page.locator('.ant-message-success');
    await expect(successMessage).toBeVisible();
  });

  test('TC-002: Вход с корректными данными', async ({ page }) => {
    await registerUser(page, email, password);
    await loginUser(page, email, password);
    // Проверяем, что на главной есть карточки фильмов
    await expect(page.locator('.ant-card').first()).toBeVisible();
  });

  test('TC-003: Вход с неверным паролем', async ({ page }) => {
    await registerUser(page, email, password);
    await page.goto('/login');
    await page.getByLabel('Email').fill(email);
    await page.getByLabel('Пароль').fill('wrongpassword');
    await page.getByRole('button', { name: 'Войти' }).click();
    const errorMessage = page.locator('.ant-message-error');
    await expect(errorMessage).toBeVisible();
    await expect(page).toHaveURL(/.*\/login/);
  });

  test('TC-009: Выход из системы', async ({ page }) => {
    await registerUser(page, email, password);
    await loginUser(page, email, password);
    await page.getByRole('menuitem', { name: 'Выйти' }).click();
    await expect(page.getByRole('menuitem', { name: 'Войти' })).toBeVisible();
    await expect(page).toHaveURL('http://localhost:3000/');
  });
});