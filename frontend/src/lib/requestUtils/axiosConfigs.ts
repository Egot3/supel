import type { AxiosRequestConfig } from 'axios';

export const withoutAuth = (url: string, method: string): AxiosRequestConfig => {
	const config: AxiosRequestConfig = {
		url,
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

export const withAuth = (url: string, method: string, token: string): AxiosRequestConfig => {
	const config: AxiosRequestConfig = {
		url,
		method: method,
		withCredentials: true,
		timeout: 10000,
		headers: {
			'Content-Type': 'application/json',
			'Hello-From': 'frontend',
			Authorization: `Bearer ${token}`
		}
	};
	return config;
};
