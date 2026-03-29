import { Page, expect } from '@playwright/test';

export function generateTestEmail(): string {
  return `test+${Date.now()}@example.com`;
}

export async function registerUser(page: Page, email: string, password: string) {
  await page.goto('/register');
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Пароль').fill(password);
  await page.getByRole('button', { name: 'Зарегистрироваться' }).click();

  await page.waitForURL(/.*\/login/, { timeout: 10000 }).catch(async () => {
    const error = page.locator('.ant-message-error');
    if (await error.isVisible()) {
      throw new Error('Registration failed: ' + (await error.textContent()));
    }
  });

  await expect(page.getByRole('button', { name: 'Войти' })).toBeVisible({ timeout: 5000 });
}

export async function loginUser(page: Page, email: string, password: string) {
  await page.goto('/login');
  await page.getByLabel('Email').fill(email);
  await page.getByLabel('Пароль').fill(password);
  await page.getByRole('button', { name: 'Войти' }).click();
  await page.waitForURL('http://localhost:3000/');
  await page.waitForSelector('.ant-spin', { state: 'hidden' });
  await page.locator('.ant-card').first().waitFor({ timeout: 10000 });
  await expect(page.getByRole('menuitem', { name: 'Выйти' })).toBeVisible();
}