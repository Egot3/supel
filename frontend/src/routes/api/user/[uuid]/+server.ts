import { type RequestHandler, json } from '@sveltejs/kit';
import { error } from 'console';

export const GET: RequestHandler = async ({ params, fetch, cookies }) => {
	const token = cookies.get('auth_token');
	if (!token) {
		throw error(401, 'Unauthorized');
	}

	const { uuid } = params;
	const response = await fetch(`http://localhost/v1/user/${uuid}`, {
		method: 'GET',
		headers: {
			Authorization: `Bearer ${token}`,
			'Content-Type': 'application/json'
		}
	});
	if (!response.ok) {
		throw error(response.status, 'Failed to fetch news');
	}

	const responseJson = await response.json();

	console.log('responeJson:', responseJson);

	return json(responseJson);
};
