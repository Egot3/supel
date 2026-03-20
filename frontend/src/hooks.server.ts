import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';
import axios from 'axios';

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
		console.log('found token ', `/api/user/${token}`);
		const resp = await axios.get(`http://localhost:5003/api/user/${token}`);
		if (resp.status == 401) {
			console.log('bad token');
			event.cookies.delete('auth_token', { path: '/' });
			if (!isPublicRoute) throw redirect(302, '/login');
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
