import type { AxiosRequestConfig } from 'axios';

export const withoutAuth = (url: string, method: string, body: unknown): AxiosRequestConfig => {
	const config: AxiosRequestConfig = {
		url,
		data: body,
		method: method,
		withCredentials: true,
		timeout: 10000,
		headers: {
			'Content-Type': 'application/json',
			'Hello-From': 'frontend'
		}
	};
	return config;
};

export const withAuth = (
	url: string,
	method: string,
	token: string,
	body: unknown
): AxiosRequestConfig => {
	const config: AxiosRequestConfig = {
		url,
		method: method,
		withCredentials: true,
		timeout: 10000,
		data: body,
		headers: {
			'Content-Type': 'application/json',
			'Hello-From': 'frontend',
			Authorization: `Bearer ${token}`
		}
	};
	return config;
};
