import { fail, redirect, type Actions } from '@sveltejs/kit';
import axios, { isAxiosError } from 'axios';
import { type FormError } from '$lib/types/error';
import { checkPassword } from '$lib/passwordUtils/checkPassword';

/**
 @satisfies {import{'./$types'}.actions;}
 */
export const actions: Actions = {
	login: async ({ request, url, cookies }) => {
		//ох, зря я сюда полез

		console.log('Login fired');

		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password'); //перехват пароля уже проблема пользователя

		if (!email) {
			return fail(422, { email, missing: true });
		}
		if (!password) {
			return fail(422, { password, missing: true });
		}
		const passwordString = password?.toString();
		console.log(passwordString);

		const err: FormError = checkPassword(passwordString);
		if (err.status != 200) {
			return fail(err.status, {
				error: err.error,
				cause: err.cause ?? '',
				errorMessage: err.errorMessage
			});
		}

		try {
			const tokenResponse = await axios.post('http://localhost/v1/public/login', {
				email: email,
				password: password
			});

			const setCookies: Array<string> = tokenResponse.headers['set-cookie'] ?? [];
			setCookies.forEach((cookie) => {
				const eqIdx = cookie.indexOf('=');
				const seIdx = cookie.indexOf(';');

				const cookieName = eqIdx !== -1 ? cookie.substring(0, eqIdx) : cookie;
				const cookieValue = seIdx !== -1 ? cookie.substring(eqIdx + 1, seIdx) : '';
				cookies.set(cookieName, cookieValue, { path: '/', httpOnly: true, sameSite: 'lax' });
			});
		} catch (err) {
			if (isAxiosError(err)) {
				console.log('Err while logging in: ', err.message, err.status);
			}
			return fail(500, 'Try logging in later');
		}

		if (url.searchParams.has('redirectTo')) {
			throw redirect(303, url.searchParams.get('redirectTo')!);
		} else {
			console.log('threw a redirect to /news');
			throw redirect(303, '/news');
		}
	}
};
