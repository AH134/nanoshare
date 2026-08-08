import { request } from "./client";

export interface User {
	id: number;
	username: string;
	createdAt: string;
}

export const authService = {
	async getMe(): Promise<User | null> {
		try {
			return await request<User>("/api/me");
		} catch {
			return null;
		}
	},

	async login(username: string, password: string): Promise<User> {
		return await request<User>("/api/auth/login", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ username, password }),
		});
	},

	async logout(): Promise<void> {
		await request<void>("/api/auth/logout", { method: "POST" });
	},
};
