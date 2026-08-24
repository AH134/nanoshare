import { request } from "./client";
import type { Link, LinkPayload } from "./link";

export interface UploadedFile {
  id: number;
  originalFilename: string;
  sizeBytes: number;
  mimeType: string;
  uploadedAt: string;
  links: Link[];
}

export async function uploadFile(
  file: File,
  linkOptions: LinkPayload,
): Promise<UploadedFile> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("linkOptions", JSON.stringify(linkOptions));

  return request<UploadedFile>("/api/files", {
    method: "POST",
    body: formData,
  });
}

export async function getFiles(): Promise<UploadedFile[]> {
  return request<UploadedFile[]>("/api/files");
}
