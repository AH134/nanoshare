import dayjs from "dayjs";
import { Files, FileText, Trash2 } from "lucide-react";
import prettyBytes from "pretty-bytes";
import { useFiles } from "#/hooks/use-files";
import type { UploadedFile } from "#/services/file";
import { CopyLinkButton } from "./CopyLinkButton";

export interface FileProps {
  file: UploadedFile;
}

export function File({ file }: FileProps) {
  const { deleteMutation } = useFiles();

  const handleDelete = async (fileID: number) => {
    try {
      await deleteMutation.mutateAsync(fileID);
    } catch (err) {
      console.error(err);
    }
  };

  const link = file.links?.[0];
  const firstLink = link
    ? `${window.location.origin}/d/${file.links[0].token}`
    : `${window.location.origin}/d/no-link`;

  const expiryBadge = (() => {
    if (!link.expiresAt) {
      return (
        <span className="badge badge-soft badge-sm badge-primary">
          {" "}
          Never expire
        </span>
      );
    }

    const isExpired = dayjs(link.expiresAt).isBefore(dayjs());
    return (
      <span
        className={`badge badge-soft badge-sm ${isExpired ? "badge-error" : "badge-warning"}`}
      >
        {isExpired ? "Expired" : `Expires ${dayjs(link.expiresAt).fromNow()}`}
      </span>
    );
  })();

  return (
    <li className="list-row">
      <div className="p-2 rounded-lg size-10 flex items-center justify-center bg-primary/10">
        <FileText className="size-4 text-primary" />
      </div>
      <div className="min-w-0">
        <div className="flex items-center justify-between gap-1">
          <p className="truncate">{file.originalFilename}</p>
          <span className="shrink-0">{expiryBadge}</span>
        </div>
        <div className="flex items-center gap-1 text-xs text-base-content/60">
          <span className="shrink-0">
            {prettyBytes(file.sizeBytes, { binary: true })}
          </span>
          <span>·</span>
          <span className="shrink-0">{dayjs(file.uploadedAt).fromNow()}</span>
          <span>·</span>
          <span className="truncate">{file.mimeType}</span>
        </div>
      </div>
      <div className="flex gap-1">
        <CopyLinkButton url={firstLink} />
        <div className="tooltip" data-tip="Delete">
          <button
            type="button"
            className="btn btn-square btn-ghost group hover:bg-error/20"
            onClick={() => handleDelete(file.id)}
          >
            <Trash2 className="size-4 group-hover:text-error" />
          </button>
        </div>
      </div>
    </li>
  );
}
