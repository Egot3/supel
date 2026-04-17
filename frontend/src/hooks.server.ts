import { redirect } from '@sveltejs/kit';
import type { Handle, HandleFetch } from '@sveltejs/kit';
import axios, { isAxiosError } from 'axios';

const PUBLIC_ROUTES = ['/login', '/register'];

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('auth_token');
	console.log('token:', token);

	const isPublicRoute = PUBLIC_ROUTES.some((route) => event.url.pathname.startsWith(route));
	console.log('ispub:', isPublicRoute);

	if (isPublicRoute) {
		return resolve(event);
	}

	if (token) {
		console.log('found token ', token);
		try {
			await axios.get(`http://localhost/v1/public/validate/${token}`);
		} catch (err) {
			if (isAxiosError(err) && err.status === 401) {
				console.log('bad token');
				event.cookies.delete('auth_token', { path: '/' });
				throw redirect(302, '/login');
			}
		}
	}

	console.log('!isPublicRoute && !token', !isPublicRoute && !token);
	if (!isPublicRoute && !token) {
		console.log('not pub and no token');
		throw redirect(302, `/login?redirectTo=${event.url.pathname}`);
	}

	console.log('resolved with no issue');
	return resolve(event);
};

export const handleFetch: HandleFetch = async ({ request, fetch, event }) => {
	const response = await fetch(request);

	if (response.status === 401) {
		event.cookies.delete('auth_token', { path: '/' });
		throw redirect(302, '/login');
	}

	return response;
};
