import { expect } from '@playwright/test';
import { test } from './fixtures';

test('login page has form', async ({ page }) => {
	await page.goto('/login');
	await expect(page.locator('form')).toBeVisible();
});

test('login with no digit password', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[name="email"]', 'test@example.com');
	await page.fill('input[name="password"]', 'PASSword!!!');
	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page.locator('form div .passwordTooltip')).toContainText('digits');
});

test('login with no lowercase password', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[name="email"]', 'test@example.com');
	await page.fill('input[name="password"]', 'PASSWORD123456!!!');
	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page.locator('form div .passwordTooltip')).toContainText('lowercase');
});

test('login with no uppercase password', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[name="email"]', 'test@example.com');
	await page.fill('input[name="password"]', 'asdasd123123!!!');
	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page.locator('form div .passwordTooltip')).toContainText('UPPERCASE');
});

test('login with small password length', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[name="email"]', 'test@example.com');
	await page.fill('input[name="password"]', '123!!!');
	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page.locator('form div .passwordTooltip')).toContainText('longer');
});

test('login with big password length', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[name="email"]', 'test@example.com');
	await page.fill(
		'input[name="password"]',
		'passwASASord123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffupassword123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90password123456sfhasdfhdsajfhuieyr283523h43hdfspasswASASord123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffupassword123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90password123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90da8f90ahfasdfaufyad8fsd8fdafdauffuda8f90da8f90!!!'
	);
	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page.locator('form div .passwordTooltip')).toContainText('shorter');
});

test('login with whitespace in password', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[name="email"]', 'test@example.com');
	await page.fill('input[name="password"]', 'password 123456!!!');
	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page.locator('form div .passwordTooltip')).toContainText('whitespace');
});
