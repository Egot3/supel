import { withAuth } from '$lib/requestUtils/axiosConfigs';
import { fail, type Actions } from '@sveltejs/kit';
import axios, { isAxiosError } from 'axios';
import type { PageServerLoad } from './$types.js';
import type { newCooked, newRaw, ImageMeta } from '$lib/types/new.js';
import sharp from 'sharp';
import { ImageSizeError } from '$lib/types/error.js';

export const actions: Actions = {
	post: async ({ request, cookies }) => {
		const data = await request.formData();

		const caption = data.get('caption');
		const text = data.get('textArea') as string;
		const image = data.get('image') as File;

		const token = cookies.get('auth_token');

		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}

		try {
			const bodyPutRequestCfg = withAuth('http://localhost/v1/news/body', 'post', token, {
				bodyName: caption?.slice(0, 16)
			});

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

			let imageKey: string = '';
			if (image.size > 0) {
				console.log('getting an image: ', image);
				console.log('image name:', image.name);
				const imagePutRequestCfg = withAuth('http://localhost/v1/news/images', 'post', token, {
					images: [
						{
							fileName: image.name,
							mime: 'image/webp'
						} as ImageMeta
					]
				});

				const imagePutResponse = await axios(imagePutRequestCfg);
				console.log('img resp: ', imagePutResponse.data);
				const putUrl = imagePutResponse.data.targets[0].uploadUrl;
				imageKey = imagePutResponse.data.targets[0].fileKey;

				console.log('key and url: ', imageKey, putUrl);
				const buffer = Buffer.from(await image.arrayBuffer());
				const { width, height } = await sharp(buffer).metadata();
				if (height / width > 20 || width / height > 20) {
					throw new ImageSizeError("Image's aspect ratio is too big");
				}

				const webpBuffer = await sharp(buffer).webp({ quality: 75 }).toBuffer();

				await axios.put(putUrl, webpBuffer, {
					headers: {
						'Content-Type': `image/webp`,
						'Content-Length': webpBuffer.length
					}
				});
			}

			const cfg = withAuth('http://localhost/v1/news', 'POST', token, {
				caption: caption,
				bodyKey: bodyKey,
				bodySize: bodyBytes.byteLength,
				imageKeys: imageKey === '' ? [] : [imageKey]
			});

			axios(cfg);
		} catch (err) {
			console.log('failed to make a req: ', err);

			if (isAxiosError(err) && err.status) return fail(err.status, 'axios error');
			if (err instanceof ImageSizeError) {
				return fail(err.status, err.message);
			}
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

	const response = await fetch(`http://localhost/v1/news?page=0&size=50`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!response.ok) {
		console.log(response);
		console.log('response is not ok: ', response.status);
		return { news: [] };
	}

	const { news } = (await response.json()) as { news: newRaw[] };
	console.log('news: ', news);

	const cookedNews = news.map(async (item: newRaw) => {
		const body = item.bodyUrl ? await (await fetch(item.bodyUrl)).text() : null;
		return { ...item, body } as newCooked;
	});

	const promisedNews = await Promise.all(cookedNews);
	console.log('cooked: ', promisedNews);

	return { news: promisedNews };
};
