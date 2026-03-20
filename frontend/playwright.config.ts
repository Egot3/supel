import { defineConfig } from '@playwright/test';

export default defineConfig({
	reporter: [['html', { outputFolder: 'playwright-report' }]],
	testDir: 'e2e',
	workers: 1,
	use: {
		video: 'retain-on-failure',
		baseURL: 'http://localhost:5173/',
		trace: 'on-first-retry'
	},
	webServer: [
		{
			command: 'ts-node e2e/mocks/mock-user.ts',
			url: 'http://localhost:5003/',
			timeout: 10000,
			stdout: 'pipe'
		},
		{
			command: 'ts-node e2e/mocks/mock-post.ts',
			url: 'http://localhost:5004/',
			timeout: 10000,
			stdout: 'pipe' //не по дефолту кстати
		},
		{
			command: 'yarn run dev',
			url: 'http://localhost:5173/'
		}
	]
});
