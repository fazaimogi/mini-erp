import { test } from '@e2e-dev/web';
import { expect } from 'e2e';

test('login lalu buka dashboard', async ({ app, screen, browser }) => {
  await app.open('/login');
  await expect(screen.getByRole('heading', 'ERP Mini System')).toBeVisible();

  await screen.getByLabel('Email').fill('admin@example.com');
  await screen.getByLabel('Password').fill('password');
  await screen.getByRole('button', 'Login').tap();

  await expect(browser).toHaveURL('/dashboard');
  await expect(screen.getByRole('heading', 'Dashboard')).toBeVisible();
});