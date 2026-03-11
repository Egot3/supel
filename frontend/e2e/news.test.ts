import { expect, test } from '@playwright/test';
import { randomBytes } from 'crypto';

test('posting unauthorized', async ({ page, request }) => {
	const randomPosts: { caption: string; text: string }[] = [];
	for (let i = 0; i < 20; i++) {
		randomPosts.push({
			caption: randomBytes(8).toString('hex'),
			text: randomBytes(255).toString('hex')
		});
	}

	randomPosts.forEach(async (post) => {
		await request.post('http://localhost:5004/api/post', {
			data: post
		});
	});

	page.waitForTimeout(204);
	const found = (await (await request.get('http://localhost:5004/api/post')).json()).total;
	expect(found).toBeFalsy();
	await request.post('http://localhost:5004/clear');
});

test('loading content with authorization', async ({ request, browser }) => {
	const token = (await (await request.post('http://localhost:5003/api/user/register')).json())
		.token;

	const ctx = await browser.newContext({
		extraHTTPHeaders: {
			Authorization: `Bearer ${token}`
		}
	});

	const page = await ctx.newPage();

	const randomPosts: { caption: string; text: string }[] = [];
	for (let i = 0; i < 20; i++) {
		randomPosts.push({
			caption: randomBytes(8).toString('hex'),
			text: randomBytes(255).toString('hex')
		});
	}

	randomPosts.forEach(async (post) => {
		await request.post('http://localhost:5004/api/post', {
			data: post
		});
	});

	await page.waitForTimeout(200);
	await page.goto('/news');

	const actualPosts = await page.locator('.new').all();

	actualPosts.forEach(async (actualPost) => {
		expect(
			randomPosts.includes({
				caption: String(await actualPost.locator('.caption').textContent()),
				text: String(await actualPost.locator('.text').textContent())
			})
		).toBe(true);
	});

	await request.post('http://localhost:5004/clear');
	await request.post('http://localhost:5003/clear');
});
