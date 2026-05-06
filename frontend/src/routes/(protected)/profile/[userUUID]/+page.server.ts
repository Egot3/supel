import { fail, type Actions } from '@sveltejs/kit';
import { withAuth } from '$lib/requestUtils/axiosConfigs';
import axios, { isAxiosError } from 'axios';
import sharp from 'sharp';
import type { PageServerLoad } from './$types.js';
import { type newCooked, type newRaw } from '$lib/types/new';
import type { User } from '$lib/types/user.js';

export const actions: Actions = {
	deleteUser: async ({ cookies, fetch, params }) => {
		const token = cookies.get('auth_token');
		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}
		const { userUUID } = params;

		try {
			const resp = await fetch(`http://localhost/v1/user/identity/${userUUID}`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (!resp.ok) {
				console.log(resp);
				throw resp.status;
			}
		} catch (err) {
			if (err === 404) {
				return fail(404, {
					err: 'USER_NOT_FOUND',
					errorMessage: 'user with this uuid was not found'
				});
			}
			return fail(500, {
				err: 'INTERNAL_SERVER',
				errorMessage: "Couldn't delete user, try again later"
			});
		}

		cookies.delete('auth_token', {
			path: '/'
		});
		return { success: true };
	},
	updateProfile: async ({ request, cookies, fetch, params }) => {
		const token = cookies.get('auth_token');

		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}

		const data = await request.formData();

		const nickname = data.get('nickname');
		const description = data.get('description');
		const pfp = data.get('pfp') as File;
		const size = Number(data.get('clipSize') ?? 0);
		const posX = Number(data.get('clipX') ?? 0);
		const posY = Number(data.get('clipY') ?? 0);
		const { userUUID } = params;
		if (pfp) {
			const pfpGetPUTUrlCfg = withAuth('http://localhost/v1/user/avatar', 'post', token, {
				uuid: userUUID
			});

			let pfpPUTUrl = '';
			try {
				const pfpGetPutUrlresponseNews = await axios(pfpGetPUTUrlCfg);
				pfpPUTUrl = pfpGetPutUrlresponseNews.data.avatarUrl;
				console.log('put Url: ', pfpPUTUrl);
			} catch (err) {
				console.log('error while fetching avatar put url:', err);
				fail(500, 'axios error');
			}

			const image = sharp(await pfp.arrayBuffer());
			const { width, height } = await image.metadata();
			if (!(posX + size <= width && posX >= 0 && posY + size <= height && height >= 0)) {
				console.log('err: bad clip sizing');
			}
			const imageRefined = await image
				.extract({
					left: posX - size,
					top: posY - size,
					width: size * 2 || height || width,
					height: size * 2 || height || width
				})
				.resize(256, 256)
				.webp({
					quality: 100
				})
				.toBuffer();

			try {
				const resp = await fetch(pfpPUTUrl, {
					method: 'put',
					body: new Uint8Array(imageRefined),
					headers: {
						'content-type': 'image/webp',
						'content-length': String(imageRefined.byteLength)
					}
				});
				console.log('put url resp:', resp);
			} catch (err) {
				console.log('error while putting pfp: ', err);
				fail(500, 'Failed to PUT pfp, try again later');
			}
		}

		const cfg = withAuth('http://localhost/v1/user', 'patch', token, {
			uuid: userUUID,
			nickname: nickname,
			description: description
		});

		try {
			axios(cfg);
		} catch (err) {
			console.log('Axios error: ', err);
			if (isAxiosError(err)) {
				fail(err.status ?? 500, 'Error while updating profile');
			}
			fail(500, 'unknown error');
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

export const load: PageServerLoad = async ({ fetch, cookies, params }) => {
	const token = cookies.get('auth_token');
	if (!token) {
		console.log('No token on load');
		return { news: [], user: {} as User };
	}

	const { userUUID } = params;
	console.log('user uuid to fetch news && load user:', userUUID);

	const responseUser = await fetch(`http://localhost/v1/user/${userUUID}`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!responseUser.ok) {
		console.log(responseUser);
		console.log('response of users is not ok: ', responseUser.status);
		return { news: [], user: {} as User };
	}
	console.log('Response for user(nickname): ');

	const responseNews = await fetch(`http://localhost/v1/news/${userUUID}?page=0&size=5`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!responseNews.ok) {
		console.log(responseNews);
		console.log('responseNews is not ok: ', responseNews.status);
		return { news: [], user: {} as User };
	}

	const { news } = (await responseNews.json()) as { news: newRaw[] };
	console.log('news: ', news);

	const cookedNews = news.map(async (item: newRaw) => {
		const body = item.bodyUrl ? await (await fetch(item.bodyUrl)).text() : null;
		return { ...item, body } as newCooked;
	});

	const promisedNews = await Promise.all(cookedNews);
	console.log('cooked: ', promisedNews);

	const { user } = await responseUser.json();

	if (!user || Object.keys(user).length === 0) {
		return { news: [], user: {} as User };
	}

	return { news: promisedNews, user: user };
};
