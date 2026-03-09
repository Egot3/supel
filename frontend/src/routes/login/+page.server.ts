import { fail, redirect, type Actions } from '@sveltejs/kit';
import axios from 'axios';

/**
 @satisfies {import{'./$types'}.actions;}
 */
export const actions: Actions = {
	login: async ({ request, url }) => {
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

		const passwordLength = passwordString.length;
		if (passwordLength <= 8) {
			return fail(422, {
				password,
				error: 'SHORT_PASSWORD',
				errorMessage: 'password must be longer than 7 characters'
			});
		}
		if (passwordLength >= 255) {
			return fail(422, {
				password,
				error: 'LONG_PASSWORD',
				errorMessage: 'password must be shorter than 256 characters'
			});
		}

		const digitCount = (passwordString.match(/\d/g) || []).length;
		if (digitCount < 4) {
			return fail(422, {
				password,
				error: 'DIGIT_IN_PASSWORD_IS_REQUIRED',
				errorMessage: '5 digits in password are required' // in bank account too
			});
		}

		const whiteSpace = /\s+/.exec(passwordString);
		if (whiteSpace !== null) {
			return fail(422, {
				password,
				error: 'WHITESPACE_IN_PASSWORD',
				errorMessage: `Password mustn't have any whitespace characters`
			});
		}

		const uppercaseCount = (passwordString.match(/[A-Z]/g) || []).length;
		if (uppercaseCount < 1) {
			return fail(422, {
				password,
				error: 'UPPERCASE_IN_PASSWORD_IS_REQUIRED',
				errorMessage: 'Password must have at least 2 UPPERCASE letters'
			});
		}

		const lowercaseCount = (passwordString.match(/([a-z])/g) || []).length;
		if (lowercaseCount < 1) {
			return fail(422, {
				password,
				error: 'LOWERCASE_IN_PASSWORD_IS_REQUIRED',
				errorMessage: 'Password must have at least 2 lowercase letters'
			});
		}

		const tokenResponse = await axios.post('http://localhost:5003/api/user/login', {
			email: email,
			password: password
		});
		if (tokenResponse.data == null) {
			return fail(500);
		}
		if (tokenResponse.status == 401) {
			return fail(401, {
				error: tokenResponse.data.error,
				errorMessage: tokenResponse.data.errorMessage
			});
		}
		const token = tokenResponse.data.token;

		if (url.searchParams.has('redirectTo')) {
			redirect(303, url.searchParams.get('redirectTo')!);
		}

		//localStorage.setItem("token", token)

		return {
			success: true,
			token: token
		};
	}
};
