import { withAuth } from '$lib/requestUtils/axiosConfigs';
import { fail } from '@sveltejs/kit';
import axios from 'axios';

export const actions = {
	post: async ({ request, cookies }) => {
		const data = await request.formData();

		const caption = data.get('caption');
		const text = data.get('textArea');

		const token = cookies.get('auth_token');

		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}
		const cfg = withAuth('http://localhost:5004/api/post', 'POST', token, {
			text: text,
			caption: caption
		});

		const response = await axios(cfg);
		if (response.status == 401) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!' //token not found
			});
		}

		return { success: true };
	}
};
