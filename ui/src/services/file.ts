import { request } from "./client";
import type { Link } from "./link";

export interface UploadedFile {
  id: number;
  originalFilename: string;
  sizeBytes: number;
  mimeType: string;
  uploadedAt: string;
  link: Link[];
}

export async function uploadFile(file: File): Promise<UploadedFile> {
  const formData = new FormData();
  formData.append("file", file);
  return request<UploadedFile>("/api/files", {
    method: "POST",
    body: formData,
  });
}

export async function getFiles(): Promise<UploadedFile[]> {
  return request<UploadedFile[]>("/api/files");
}
