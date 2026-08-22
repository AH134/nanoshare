import { request } from "./client";

export interface Link {
  id: number;
  fileID: number;
  token: string;
  maxDownloads: number | null;
  downloadCount: number;
  createdAt: string;
  expiresAt: string | null;
  revokedAt: string | null;
}

export interface LinkPayload {
  maxDownloads?: number | null;
  expiresAt?: string | null;
}

export async function createLink(
  fileID: number,
  payload: LinkPayload,
): Promise<Link> {
  return await request<Link>(`/api/files/${fileID}/links`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function getAllFileLinks(fileID: number): Promise<Link[]> {
  return await request<Link[]>(`/api/files/${fileID}/links`);
}
