import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('auth_token');
	if (!token) {
		console.log('No token on load');
		return { user: {} };
	}

	const response = await fetch(`http://localhost/v1/user/self`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!response.ok) {
		const decoder = new TextDecoder();
		for await (const chunk of response.body) {
			console.log(decoder.decode(chunk, { stream: true }));
		}
		console.log('response is not ok in layout: ', response.status);
		return { user: {} };
	}

	const user = await response.json();

	console.log('user in layout: ', user);

	return user;
};
