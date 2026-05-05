import { fail, type Actions } from '@sveltejs/kit';
import { withAuth } from '$lib/requestUtils/axiosConfigs';
import type { TokenPayload } from '$lib/types/token';
import axios, { isAxiosError } from 'axios';
import sharp from 'sharp';

export const actions: Actions = {
	updateProfile: async ({ request, cookies, fetch }) => {
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

		let ownUUID = '';

		if (pfp) {
			try {
				const ownInfoResponse = (await (
					await fetch(`http://localhost/v1/public/validate/${token}`)
				).json()) as TokenPayload;
				ownUUID = ownInfoResponse.uuid;
			} catch (err) {
				fail(500, err); //так и живем
			}

			const pfpGetPUTUrlCfg = withAuth('http://localhost/v1/user/avatar', 'post', token, {
				uuid: ownUUID
			});

			let pfpPUTUrl = '';
			try {
				const pfpGetPutUrlResponse = await axios(pfpGetPUTUrlCfg);
				pfpPUTUrl = pfpGetPutUrlResponse.data.avatarUrl;
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
			uuid: ownUUID,
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
	}
};
