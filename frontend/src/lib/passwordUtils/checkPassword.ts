import { type FormError } from '$lib/types/error';

export const checkPassword = (passwordString: string) => {
	const passwordLength = passwordString.length;
	if (passwordLength <= 8) {
		return {
			status: 422,
			cause: passwordString,
			error: 'SHORT_PASSWORD',
			errorMessage: 'password must be longer than 7 characters'
		} as FormError;
	}
	if (passwordLength >= 255) {
		return {
			status: 422,
			cause: passwordString,
			error: 'LONG_PASSWORD',
			errorMessage: 'password must be shorter than 256 characters'
		} as FormError;
	}

	const digitCount = (passwordString.match(/\d/g) || []).length;
	if (digitCount < 4) {
		return {
			status: 422,
			cause: passwordString,
			error: 'DIGIT_IN_PASSWORD_IS_REQUIRED',
			errorMessage: '5 digits in password are required' // in bank account too
		} as FormError;
	}

	const whiteSpace = /\s+/.exec(passwordString);
	if (whiteSpace !== null) {
		return {
			status: 422,
			cause: passwordString,
			error: 'WHITESPACE_IN_PASSWORD',
			errorMessage: `Password mustn't have any whitespace characters`
		} as FormError;
	}

	const uppercaseCount = (passwordString.match(/[A-Z]/g) || []).length;
	if (uppercaseCount < 1) {
		return {
			status: 422,
			cause: passwordString,
			error: 'UPPERCASE_IN_PASSWORD_IS_REQUIRED',
			errorMessage: 'Password must have at least 2 UPPERCASE letters'
		} as FormError;
	}

	const lowercaseCount = (passwordString.match(/([a-z])/g) || []).length;
	if (lowercaseCount < 1) {
		return {
			status: 422,
			cause: passwordString,
			error: 'LOWERCASE_IN_PASSWORD_IS_REQUIRED',
			errorMessage: 'Password must have at least 2 lowercase letters'
		} as FormError;
	}

	return {
		status: 200,
		error: '',
		errorMessage: ''
	} as FormError;
};
