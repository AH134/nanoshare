export interface Envelope<T> {
	success: boolean;
	data?: T;
	error?: {
		code: string;
		message: string;
		fields?: Record<string, string>;
	};
}

export class APIError extends Error {
	code: string;
	fields?: Record<string, string>;

	constructor(code: string, message: string, fields?: Record<string, string>) {
		super(message);

		this.name = "APIError";
		this.code = code;
		this.fields = fields;
	}
}

export async function request<T>(
	url: string,
	options?: RequestInit,
): Promise<T> {
	const res = await fetch(url, { credentials: "include", ...options });
	const body: Envelope<T> = await res.json();

	if (!body.success) {
		throw new APIError(
			body.error?.code || "UNKNOWN_ERROR",
			body.error?.message || "API request failed.",
			body.error?.fields,
		);
	}

	return body.data as T;
}
