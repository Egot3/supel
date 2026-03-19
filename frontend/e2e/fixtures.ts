import { test as base } from '@playwright/test';

export const test = base.extend({
	page: async ({ page }, use, testInfo) => {
		await use(page);
		if (testInfo.status !== testInfo.expectedStatus) {
			const videoPath = `${testInfo.outputDir}/video.webm`;

			await page.video()?.saveAs(videoPath);

			await testInfo.attach('video', {
				path: videoPath,
				contentType: 'video/webm'
			});
		}
	}
});
