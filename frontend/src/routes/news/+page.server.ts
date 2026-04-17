import { withAuth } from '$lib/requestUtils/axiosConfigs';
import { fail, type Actions } from '@sveltejs/kit';
import axios, { isAxiosError } from 'axios';
import type { PageServerLoad } from './$types.js';
import type { newCooked, newRaw } from '$lib/types/new.js';

export const actions: Actions = {
	post: async ({ request, cookies }) => {
		const data = await request.formData();

		const caption = data.get('caption');
		const text = data.get('textArea') as string;

		const token = cookies.get('auth_token');

		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}

		const bodyPutRequestCfg = withAuth('http://localhost/v1/news/body', 'post', token, {
			bodyName: caption?.slice(0, 16)
		});
		try {
			const bodyPutResponse = await axios(bodyPutRequestCfg);
			const putUrl = bodyPutResponse.data.target.uploadUrl;
			const bodyKey = bodyPutResponse.data.target.fileKey;
			const bodyBytes = new TextEncoder().encode(text);
			console.log('put url: ', putUrl);

			await axios.put(putUrl, bodyBytes, {
				headers: {
					'Content-Type': 'text/markdown',
					'Content-Length': String(bodyBytes.byteLength)
				}
			});
			// console.log('resp: ', resp);

			const cfg = withAuth('http://localhost/v1/news', 'POST', token, {
				caption: caption,
				bodyKey: bodyKey,
				bodySize: bodyBytes.byteLength
			});

			axios(cfg);
		} catch (err) {
			console.log('failed to make a req: ', err);

			if (isAxiosError(err) && err.status) return fail(err.status, 'axios error');
			return fail(500, {
				err: 'INTERNAL_SERVER',
				errorMessage: 'Please, try posting later'
			});
		}

		return { success: true };
	},
	loadMore: async ({ fetch, cookies, request }) => {
		const token = cookies.get('auth_token');

		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}

		const data = await request.formData();
		const page = data.get('page');
		const size = data.get('size');
		console.log('size: %d, page: %d', size, page);

		const res = await fetch(`http://localhost/v1/news?page=${page}&size=${size}`, {
			headers: { Authorization: `Bearer ${token}` }
		});
		console.log(res.text);

		if (!res.ok) {
			return fail(res.status, { message: 'failed to serve more news' });
		}

		const { news } = (await res.json()) as { news: newRaw[] };
		const cookedNews = news.map(async (item: newRaw) => {
			const body = item.bodyUrl ? await (await fetch(item.bodyUrl)).text() : null;
			return { ...item, body } as newCooked;
		});
		console.log('cooked news type: ', typeof cookedNews);

		return { news: await Promise.all(cookedNews) };
	}
};

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('auth_token');

	if (!token) {
		console.log('No token on load');
		return { news: [] };
	}

	const response = await fetch(`http://localhost/v1/news?page=0&size=1`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!response.ok) {
		console.log('response is not ok: ', response.status);
		return { news: [] };
	}

	const { news } = (await response.json()) as { news: newRaw[] };
	const cookedNews = news.map(async (item: newRaw) => {
		const body = item.bodyUrl ? await (await fetch(item.bodyUrl)).text() : null;
		return { ...item, body } as newCooked;
	});
	console.log('cooked news type: ', typeof cookedNews);

	return { news: await Promise.all(cookedNews) };
};
