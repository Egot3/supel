import { expect, test } from '@playwright/test';
import { randomBytes } from 'crypto';

test('loading content', async ({ page, request }) => {
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
				caption: String(await actualPost.locator('.new-caption').textContent()),
				text: String(await actualPost.locator('.new-text').textContent())
			})
		).toBe(true);
	});

	await request.post('http://localhost:5004/clear');
});
