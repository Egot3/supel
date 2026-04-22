import { expect } from '@playwright/test';
import { randomBytes } from 'crypto';
import { test } from './fixtures';

test('manual postion', async ({ page, context }) => {
	await page.goto('/register');
	await page.fill('input[name="email"]', `test-${randomBytes(4).toString('hex')}@example.com`);
	await page.fill('input[name="password"]', 'PASSword12345!!!');
	await page.fill('input[name="passwordDup"]', 'PASSword12345!!!');

	await page.waitForTimeout(300);

	await page.click('button[name="submit"]');
	await page.waitForTimeout(300);
	await expect(page).toHaveURL('/news');

	let cookies = await context.cookies();
	console.log(cookies);

	//getting to new
	await expect(page.locator('button[name="postingButton"]')).toBeVisible();
	await page.click('button[name="postingButton"]');
	await page.waitForTimeout(400);

	//filling the new
	await page.fill('input[name="caption"]', 'Extremly original and interesting caption');
	await page.fill(
		'textarea[name="textArea"]',
		'An extremly long text with no grammatical(and logical) mistakes, anyway: once upon a time...'
	);
	cookies = await context.cookies();
	console.log(cookies);
	await page.click('button[name="goPostIt"]');

	await page.waitForTimeout(400);

	await page.click('button[name="reload"]');
	//checking my genius new
	await expect(page.locator('.new .caption')).toHaveText(
		'Extremly original and interesting caption'
	);
	await expect(page.locator('.new .text')).toHaveText(
		'An extremly long text with no grammatical(and logical) mistakes, anyway: once upon a time...'
	);
});
