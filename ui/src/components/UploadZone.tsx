import { Upload } from "lucide-react";
import { useState } from "react";
import { useDropzone } from "react-dropzone";
import { useFiles } from "#/hooks/use-files";

export function UploadZone() {
  const [maxDownloads, setMaxDownloads] = useState("");
  const [expiresAt, setExpiresAt] = useState("");

  const { uploadMutation } = useFiles();
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop: async (acceptedFiles) => {
      for (const file of acceptedFiles) {
        try {
          await uploadMutation.mutateAsync({
            file,
            linkOptions: {
              maxDownloads: maxDownloads ? Number(maxDownloads) : null,
              expiresAt:
                expiresAt !== "" ? new Date(expiresAt).toISOString() : null,
            },
          });
          console.log(maxDownloads ? Number(maxDownloads) : null);
          console.log(
            expiresAt !== "" ? new Date(expiresAt).toISOString() : null,
          );
        } catch (err) {
          console.error(`Failed to upload/link ${file.name}:`, err);
        }
      }
    },
  });

  return (
    <div className="flex flex-col gap-4">
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
      <div className="card bg-base-100 w-full border border-base-300">
        <div className="card-body">
          <div className="grid gap-2 sm:grid-cols-2 mb-1">
            <fieldset className="fieldset">
              <label
                className="label text-base-content"
                htmlFor="max-downloads"
              >
                Max downloads test={maxDownloads}
              </label>
              <input
                type="number"
                id="max-downloads"
                className="input"
                min={1}
                placeholder="Unlimited"
                value={maxDownloads}
                onChange={(e) => setMaxDownloads(e.target.value)}
              />
            </fieldset>
            <fieldset className="fieldset">
              <label className="label text-base-content" htmlFor="expiry-date">
                Expiry date test={new Date().toISOString()}---
                {/* {new Date().toUTCString()} */}
              </label>
              <input
                type="datetime-local"
                id="expiry-date"
                className="input"
                min={(() => {
                  const offset = new Date().getTimezoneOffset() * 60000;
                  return new Date(Date.now() - offset)
                    .toISOString()
                    .slice(0, -8);
                })()}
                value={expiresAt}
                onChange={(e) => setExpiresAt(e.target.value)}
              />
            </fieldset>
          </div>
          <p className="text-base-content/60 text-xs">
            These limits apply to files you add next. Leave max downloads or
            expiry empty for no limit.
          </p>
        </div>
      </div>
    </div>
  );
}
