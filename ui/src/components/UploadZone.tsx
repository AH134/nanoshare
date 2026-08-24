import { TriangleAlert, Upload } from "lucide-react";
import { useState } from "react";
import { useDropzone } from "react-dropzone";
import { useFiles } from "#/hooks/use-files";
import type { LinkPayload } from "#/services/link";
import { LinkOptions } from "./LinkOptions";

export function UploadZone() {
  const [linkOptions, setLinkOptions] = useState<LinkPayload>({
    maxDownloads: null,
    expiresAt: null,
  });
  const { uploadMutation } = useFiles();
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: async (acceptedFiles) => {
      for (const file of acceptedFiles) {
        try {
          await uploadMutation.mutateAsync({
            file,
            linkOptions,
          });
        } catch (err) {
          console.error(err);
        } finally {
          setLinkOptions({
            ...linkOptions,
            maxDownloads: null,
            expiresAt: null,
          });
        }
      }
    },
  });

  return (
    <div className="flex flex-col gap-4">
      {uploadMutation.isError && (
        <div role="alert" className="alert alert-error alert-soft mb-4">
          <TriangleAlert className="size-5" />
          <span>
            {uploadMutation.error?.message ?? "Invalid JSON payload."}
          </span>
        </div>
      )}
      <LinkOptions options={linkOptions} onChange={setLinkOptions} />
      <div
        {...getRootProps()}
        className={`border-2 border-dashed rounded-box py-20 px-4 text-center cursor-pointer transition-colors hover:border-primary/60 hover:bg-primary/10 ${isDragActive ? "border-primary/60 bg-primary/10" : "border-primary/20 bg-base-100"}`}
      >
        <input {...getInputProps()} />
        <div className="flex flex-col items-center gap-3">
          <div className="rounded-full bg-primary/10 p-3">
            <Upload className="size-5 text-primary " />
          </div>
          <div>
            <p className="font-semibold">Drag & drop files here</p>
            <p className="text-sm text-base-content/60">
              or click to browse from your device
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
