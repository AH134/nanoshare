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

export interface ChangePasswordPayload {
  currentPassword: string;
  newPassword: string;
}

export async function getMe(): Promise<User | null> {
  try {
    return await request<User>("/api/me");
  } catch {
    return null;
  }
}

export async function login(username: string, password: string): Promise<User> {
  return request<User>("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
}

export async function logout(): Promise<void> {
  return request<void>("/api/auth/logout", { method: "POST" });
}

export async function ChangePassword({
  currentPassword,
  newPassword,
}: ChangePasswordPayload): Promise<void> {
  return request<void>("/api/auth/password", {
    method: "POST",
    body: JSON.stringify({ currentPassword, newPassword }),
  });
}
