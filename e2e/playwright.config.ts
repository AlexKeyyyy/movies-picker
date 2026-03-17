/// <reference types="node" />

import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  // Директория с тестами
  testDir: './tests',

  // Полностью параллельный запуск тестов
  fullyParallel: true,

  // Запрет на использование test.only в CI (чтобы случайно не запустить только один тест)
  forbidOnly: !!process.env.CI,

  // Количество повторных попыток при падении теста (2 в CI, 0 локально)
  retries: process.env.CI ? 2 : 0,

  // Количество параллельных воркеров (1 в CI для стабильности, максимально локально)
  workers: process.env.CI ? 1 : undefined,

  // Репортёры результатов
  reporter: [
    ['html', { outputFolder: 'playwright-report' }], // генерируем HTML-отчёт
    ['list']                                         // выводим прогресс в консоль
  ],

  // Настройки, общие для всех тестов
  use: {
    // Базовый URL фронтенда (все относительные ссылки будут дополняться этим адресом)
    baseURL: 'http://localhost:3000',

    // Сохранять трейс (логи сети, DOM и т.д.) только при первом повторном запуске
    trace: 'on-first-retry',

    // Делать скриншот только при падении теста
    screenshot: 'only-on-failure',
  },

  // Проекты для разных браузеров
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    // Можно добавить Firefox и WebKit для кросс-браузерного тестирования, но для начала достаточно Chromium
    // {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
  ],
});