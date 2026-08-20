import { request } from "./client";

export interface User {
  id: number;
  username: string;
  createdAt: string;
}

export interface LoginPayload {
  username: string;
  password: string;
}

export async function getMe(): Promise<User | null> {
  try {
    return await request<User>("/api/me");
  } catch {
    return null;
  }
}

export async function login(username: string, password: string): Promise<User> {
  return await request<User>("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
}

export async function logout(): Promise<void> {
  await request<void>("/api/auth/logout", { method: "POST" });
}
