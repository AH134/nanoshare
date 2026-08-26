import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { useMemo } from "react";
import { File } from "#/components/File";
import { PageTitle } from "#/components/PageTitle";
import { UploadZone } from "#/components/UploadZone";
import { useFiles } from "#/hooks/use-files";

export const Route = createFileRoute("/_authenticated/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { files, isLoading } = useFiles();
  const recentFiles = useMemo(() => {
    if (!files) return [];

    return files.slice(0, 5);
  }, [files]);

  return (
    <div>
      <PageTitle
        title="Upload your files"
        description="Set a download limit
        and expiry date, then drag your file(s) into the zone below or click to browse. Then share the link."
      />
      <UploadZone />

      <div className="mt-8">
        <div className="flex justify-between items-center mb-3">
          <h2 className="text-sm font-medium text-base-content">
            Recent uploads
            <span className="text-base-content/60 ml-2">
              ({recentFiles?.length ?? 0})
            </span>
          </h2>
          <Link
            to="/files"
            className="text-sm font-medium btn btn-ghost btn-sm inline-flex items-center gap-1"
          >
            View all
            <ArrowRight className="size-4" />
          </Link>
        </div>
        {isLoading || recentFiles.length <= 0 ? (
          <div className="rounded-box border border-base-300 text-center text-sm text-base-content/60 py-10 px-5">
            <p>No files have been uploaded yet.</p>
          </div>
        ) : (
          <ul className="list bg-base-100 rounded-box border border-base-300">
            {recentFiles.map((file) => (
              <File key={file.id} file={file} />
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
