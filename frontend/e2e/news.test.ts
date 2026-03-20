import { expect } from '@playwright/test';
import { randomBytes } from 'crypto';
import { test } from './fixtures';

test('posting unauthorized', async ({ page, request }) => {
	const randomPosts: { caption: string; text: string }[] = [];
	for (let i = 0; i < 20; i++) {
		randomPosts.push({
			caption: randomBytes(8).toString('hex'),
			text: randomBytes(255).toString('hex')
		});
	}

	const promises = randomPosts.map(async (post) => {
		await request.post('http://localhost:5004/api/post', {
			data: post
		});
	});

	await Promise.all(promises);
	const found = (await (await request.get('http://localhost:5004/api/post')).json()).total;

	await page.waitForTimeout(400);
	expect(found).toBeFalsy();
	await request.post('http://localhost:5004/clear');
});

test('loading content with authorization', async ({ request, page }) => {
	const regRes = await request.post('http://localhost:5003/api/user/register', {
		data: {
			email: 'example@email.com',
			password: 'PASSword123456789!!!',
			nickname: 'none'
		}
	});
	const { token } = await regRes.json();

	page.waitForTimeout(400);

	const randomPosts: { caption: string; text: string }[] = [];
	for (let i = 0; i < 20; i++) {
		randomPosts.push({
			caption: randomBytes(8).toString('hex'),
			text: randomBytes(255).toString('hex')
		});
	}

	await Promise.all(
		randomPosts.map(async (post) => {
			await request.post('http://localhost:5004/api/post', {
				data: post,
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
		})
	);

	await page.goto('/news');
	await expect(page.locator('.new')).toHaveCount(20);

	const actualPosts = await page.locator('.new').evaluateAll((elements) =>
		elements.map((e) => ({
			caption: e.querySelector('.caption')?.textContent ?? '',
			text: e.querySelector('.text')?.textContent ?? ''
		}))
	);

	expect(actualPosts).toEqual(expect.arrayContaining(randomPosts));

	await request.post('http://localhost:5004/clear');
	await request.post('http://localhost:5003/clear');
});

test('manual postion', async ({ page, request, context }) => {
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
	await expect(page.locator('.new>.caption')).toHaveText(
		'Extremly original and interesting caption'
	);
	await expect(page.locator('.new>.textArea')).toHaveText(
		'An extremly long text with no grammatical(and logical) mistakes, anyway: once upon a time...'
	);

	await request.post('http://localhost:5004/clear');
	await request.post('http://localhost:5003/clear');
});
