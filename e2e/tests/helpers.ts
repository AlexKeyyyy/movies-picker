import { expect, Page } from '@playwright/test';

export function generateTestEmail(): string {
  return `test+${Date.now()}@example.com`;
}

export async function registerUser(page: Page, email: string, password: string) {
  await page.goto('/register');
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Пароль').fill(password);
  await page.getByRole('button', { name: 'Зарегистрироваться' }).click();
  await page.waitForURL(/.*\/login/);
}

export async function loginUser(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Пароль').fill(password);
  await page.getByRole('button', { name: 'Войти' }).click();
  await page.waitForURL('http://localhost:3000/');
  // Дожидаемся появления меню для авторизованного пользователя (пункт "Выйти")
  await expect(page.getByRole('menuitem', { name: 'Выйти' })).toBeVisible();
}