import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('auth_token');
	if (!token) {
		console.log('No token on load');
		return { news: [] };
	}

	const response = await fetch(`http://localhost/v1/user/self`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!response.ok) {
		console.log(response);
		console.log('response is not ok: ', response.status);
		return { user: {} };
	}

	return await response.json();
};
