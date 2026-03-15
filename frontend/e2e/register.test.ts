import { expect, test } from '@playwright/test';
import { randomBytes } from 'crypto';

test('register page has form', async ({ page }) => {
	await page.goto('/register');
	await expect(page.locator('form')).toBeVisible();
});

test('register with unequal passwords', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'PASSword12345!!!');
	await page.fill('input[name="passwordDup"]', 'NotAPASSword12345!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('unequal');
});

test('register with no digit password', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'PASSword!!!');
	await page.fill('input[name="passwordDup"]', 'PASSword!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('digits');
});

test('register with no lowercase password', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'PASSWORD123456!!!');
	await page.fill('input[name="passwordDup"]', 'PASSWORD123456!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('lowercase');
});

test('register with no uppercase password', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'asdasd123123!!!');
	await page.fill('input[name="passwordDup"]', 'asdasd123123!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('UPPERCASE');
});

test('register with small password length', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', '123!!!');
	await page.fill('input[name="passwordDup"]', '123!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('longer');
});

test('register with big password length', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill(
		'input[name="password"]',
		'passwASASord123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffupassword123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90password123456sfhasdfhdsajfhuieyr283523h43hdfspasswASASord123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffupassword123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90password123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90da8f90ahfasdfaufyad8fsd8fdafdauffuda8f90da8f90!!!'
	);
	await page.fill(
		'input[name="passwordDup"]',
		'passwASASord123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffupassword123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90password123456sfhasdfhdsajfhuieyr283523h43hdfspasswASASord123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffupassword123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90password123456sfhasdfhdsajfhuieyr283523h43hdfsahfasdfaufyad8fsd8fdafdauffuda8f90da8f90ahfasdfaufyad8fsd8fdafdauffuda8f90da8f90!!!'
	);

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('shorter');
});

test('register with whitespace in password', async ({ page }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'password 123456!!!');
	await page.fill('input[name="passwordDup"]', 'password 123456!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page.locator('form p#error')).toContainText('whitespace');
});

test('register with right credits', async ({ page, request }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'PASSword12345!!!');
	await page.fill('input[name="passwordDup"]', 'PASSword12345!!!');

	await page.waitForTimeout(100);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(100);
	await expect(page).toHaveURL('/');
	await request.post('http://localhost:5003/clear');
});
