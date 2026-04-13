import { withAuth } from '$lib/requestUtils/axiosConfigs';
import { error, fail } from '@sveltejs/kit';
import axios, { isAxiosError } from 'axios';
import type { PageServerLoad } from './$types.js';
import type { newRaw } from '$lib/types/new.js';

export const actions = {
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

			axios.put(putUrl, bodyBytes, {
				headers: {
					'Content-Type': 'text/markdown',
					'Content-Length': String(bodyBytes.byteLength)
				}
			});

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
	}
};

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('auth_token');

	if (!token) {
		return fail(401, {
			err: 'NO_AUTH',
			errorMessage: 'You are not registered!'
		});
	}

	const response = await fetch(`http://localhost/v1/news?page=0&size=50`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!response.ok) throw error(response.status, 'failed to fetch news');

	const { news } = (await response.json()) as { news: newRaw[] };
	const cookedNews = news.map((item: newRaw) => ({
		...item,
		body: item.bodyUrl
			? fetch(item.bodyUrl).then((r) => {
					if (!r.ok) throw new Error('Failed to fetch body');
					return r.text();
				})
			: null
	}));

	return { news: cookedNews };
};
