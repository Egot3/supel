import { json, type RequestHandler } from '@sveltejs/kit';
import { error } from 'console';
import { type newRaw } from '$lib/types/new';

export const GET: RequestHandler = async ({ url, cookies, fetch }) => {
	const token = cookies.get('auth_token');

	if (!token) {
		throw error(401, 'Unauthorized');
	}

	const page = url.searchParams.get('page');
	const size = url.searchParams.get('size');

	const response = await fetch(`http://localhost/v1/news?page=${page}&size=${size}`, {
		headers: {
			Authorization: `Bearer ${token}`,
			'Content-Type': 'application/json'
		}
	});

	if (!response.ok) {
		throw error(response.status, 'Failed to fetch news');
	}

	const data = (await response.json()) as { news: newRaw[] }; // readablitity++

	//console.log(data);

	return json(data);
};
