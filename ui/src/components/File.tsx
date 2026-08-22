import dayjs from "dayjs";
import { FileText, Trash2 } from "lucide-react";
import prettyBytes from "pretty-bytes";
import { useLinks } from "#/hooks/use-links";
import type { UploadedFile } from "#/services/file";
import { CopyLinkButton } from "./CopyLinkButton";

export interface FileProps {
  file: UploadedFile;
}

export function File({ file }: FileProps) {
  const { links } = useLinks(file.id);
  console.log(links?.[0]);
  const firstLink = links?.[0]
    ? `${window.location.origin}/d/${links[0].token}`
    : `${window.location.origin}/d/no-link`;

  return (
    <li className="list-row">
      <div className="p-2 rounded-lg size-10 flex items-center justify-center bg-primary/10">
        <FileText className="size-4 text-primary" />
      </div>
      <div>
        <div>
          {file.originalFilename}
          {file.mimeType}
        </div>
        <div className="flex items-center gap-1 text-xs text-base-content/60">
          <span>{prettyBytes(file.sizeBytes, { binary: true })}</span>
          <span>·</span>
          <span>{dayjs(file.uploadedAt).fromNow()}</span>
        </div>
      </div>
      <div className="flex gap-1">
        <CopyLinkButton url={firstLink} />
        <div className="tooltip" data-tip="Delete">
          <button
            type="button"
            className="btn btn-square btn-ghost group hover:bg-error/20"
          >
            <Trash2 className="size-4 group-hover:text-error" />
          </button>
        </div>
      </div>
    </li>
  );
}
