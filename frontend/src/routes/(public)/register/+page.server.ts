import { fail, redirect } from '@sveltejs/kit';
import axios, { isAxiosError } from 'axios';

/**
 @satisfies {import{'./$types'}.actions;}
 */
export const actions = {
	register: async ({ request, url, cookies }) => {
		const data = await request.formData();

		const email = data.get('email');
		const password = data.get('password');
		const nickname = data.get('nickname');
		const passwordDup = data.get('passwordDup');

		if (!email) {
			return fail(422, { email, error: 'EMAIL_MISSING', errorMessage: 'Email is required' });
		}
		if (!password) {
			return fail(422, {
				password,
				error: 'PASSWORD_MISSING',
				errorMessage: 'Password is required'
			});
		}
		if (password !== passwordDup) {
			console.log('unequal passwords');
			return fail(422, {
				password,
				passwordDup,
				error: 'PASSWORDS_UNEQUALITY',
				errorMessage: 'Passwords are unequal'
			});
		}

		const passwordString = password?.toString();

		const passwordLength = passwordString.length;
		if (passwordLength <= 8) {
			console.log('short password');
			return fail(422, {
				password,
				error: 'SHORT_PASSWORD',
				errorMessage: 'password must be longer than 7 characters'
			});
		}
		if (passwordLength >= 255) {
			console.log('long password');
			return fail(422, {
				password,
				error: 'LONG_PASSWORD',
				errorMessage: 'password must be shorter than 256 characters'
			});
		}

		const digitCount = (passwordString.match(/\d/g) || []).length;
		if (digitCount < 4) {
			console.log('digit count');
			return fail(422, {
				password,
				error: 'DIGIT_IN_PASSWORD_IS_REQUIRED',
				errorMessage: '5 digits in password are required' // in bank account too
			});
		}

		const whiteSpace = /\s+/.exec(passwordString);
		if (whiteSpace !== null) {
			console.log('whitespace is in');
			return fail(422, {
				password,
				error: 'WHITESPACE_IN_PASSWORD',
				errorMessage: `Password mustn't have any whitespace characters`
			});
		}

		const uppercaseCount = (passwordString.match(/[A-Z]/g) || []).length;
		if (uppercaseCount < 1) {
			console.log('too little uppercase');
			return fail(422, {
				password,
				error: 'UPPERCASE_IN_PASSWORD_IS_REQUIRED',
				errorMessage: 'Password must have at least 2 UPPERCASE letters'
			});
		}

		const lowercaseCount = (passwordString.match(/([a-z])/g) || []).length;
		if (lowercaseCount < 1) {
			console.log('too little lowercase');
			return fail(422, {
				password,
				error: 'LOWERCASE_IN_PASSWORD_IS_REQUIRED',
				errorMessage: 'Password must have at least 2 lowercase letters'
			});
		}

		try {
			const tokenResponse = await axios.post('http://localhost/v1/public/register', {
				email: email,
				password: password,
				nickname: nickname
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
			console.log(err);
			if (isAxiosError(err)) {
				if (err.status === 500) {
					return fail(500, {
						err: 'SERVER_ERROR',
						errorMessage: 'Unexpected server error'
					});
				}

				if (err.status === 409) {
					return fail(409, {
						err: 'EMAIL_IN_USE',
						errorMessage: 'Email already in use'
					});
				}
				return fail(500);
			}
		}

		if (url.searchParams.has('redirectTo')) {
			throw redirect(303, url.searchParams.get('redirectTo')!);
		} else {
			throw redirect(303, '/news');
		}
	}
};
