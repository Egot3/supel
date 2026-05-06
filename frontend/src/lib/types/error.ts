export type FormError = {
	status: number;
	error: string;
	errorMessage: string;
	cause?: string;
};

export class ServerError extends Error {
	public readonly status: number;

	constructor(message: string, status: number) {
		super(message);
		this.name = this.constructor.name;
		this.status = status;
		Object.setPrototypeOf(this, ServerError.prototype);
		Error.captureStackTrace(this, this.constructor);
	}
}

export class ImageSizeError extends ServerError {
	constructor(message: string = 'Image is too big') {
		super(message, 422);
	}
}

export class InternalServerError extends ServerError {
	constructor(message: string = 'Internal server error') {
		super(message, 500);
	}
}
