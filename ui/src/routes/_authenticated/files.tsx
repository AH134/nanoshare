import { createFileRoute } from "@tanstack/react-router";
import { Search } from "lucide-react";
import prettyBytes from "pretty-bytes";
import { useMemo, useState } from "react";
import { File } from "#/components/File";
import { PageTitle } from "#/components/PageTitle";
import { useFiles } from "#/hooks/use-files";

export const Route = createFileRoute("/_authenticated/files")({
  component: RouteComponent,
});

type FileTypeFilter = "all" | "image" | "video" | "audio" | "other";
type SortFilter = "newest" | "oldest" | "largest" | "smallest" | "name";

function RouteComponent() {
  const { files, isLoading } = useFiles();
  const [searchQuery, setSearchQuery] = useState("");
  const [fileTypeFilter, setFileTypeFilter] = useState<FileTypeFilter>("all");
  const [sortFilter, setSortFilter] = useState<SortFilter>("newest");

  const filteredFiles = useMemo(() => {
    if (!files) return [];

    let searchResult = files.filter((file) =>
      file.originalFilename.includes(searchQuery.toLowerCase()),
    );

    switch (fileTypeFilter) {
      case "all":
        break;
      case "other":
        searchResult = searchResult.filter((file) =>
          file.mimeType.includes("application"),
        );
        break;
      default:
        searchResult = searchResult.filter((file) =>
          file.mimeType.includes(fileTypeFilter),
        );
        break;
    }

    switch (sortFilter) {
      case "newest":
        return searchResult;
      case "oldest":
        return searchResult.sort(
          (a, b) =>
            new Date(a.uploadedAt).getTime() - new Date(b.uploadedAt).getTime(),
        );
      case "largest":
        return searchResult.sort((a, b) => b.sizeBytes - a.sizeBytes);
      case "smallest":
        return searchResult.sort((a, b) => a.sizeBytes - b.sizeBytes);
      case "name":
        return searchResult.sort((a, b) =>
          a.originalFilename.localeCompare(b.originalFilename),
        );
    }
  }, [files, searchQuery, fileTypeFilter, sortFilter]);

  const totalBytes = useMemo(() => {
    if (!files) return 0;

    return files
      .map((file) => file.sizeBytes)
      .reduce((acc, bytes) => acc + bytes, 0);
  }, [files]);

  return (
    <div>
      <PageTitle
        title="Files"
        description={`${files ? files.length : 0} total file(s) · ${prettyBytes(totalBytes, { binary: true })}`}
      />
      <div className="grid grid-cols-1 sm:grid-cols-[2fr_1fr_1fr_auto] gap-2 mb-4">
        <label className="input w-full">
          <Search className="size-4.5" />
          <input
            type="search"
            className="grow"
            placeholder="Search by file name"
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </label>
        <select
          className="select w-full"
          value={fileTypeFilter}
          onChange={(e) => setFileTypeFilter(e.target.value as FileTypeFilter)}
        >
          <option value={"all"}>All types</option>
          <option value={"image"}>Images</option>
          <option value={"video"}>Videos</option>
          <option value={"audio"}>Audio</option>
          <option value={"other"}>Others</option>
        </select>
        <select
          className="select w-full"
          value={sortFilter}
          onChange={(e) => setSortFilter(e.target.value as SortFilter)}
        >
          <option value={"newest"}>Newest first</option>
          <option value={"oldest"}>Oldest first</option>
          <option value={"largest"}>Largest first</option>
          <option value={"smallest"}>Smallest first</option>
          <option value={"name"}>Name (A-Z)</option>
        </select>
        <button
          type="button"
          disabled={
            !(
              searchQuery !== "" ||
              sortFilter !== "newest" ||
              fileTypeFilter !== "all"
            )
          }
          className="btn btn-soft btn-error"
          onClick={() => {
            setFileTypeFilter("all");
            setSortFilter("newest");
            setSearchQuery("");
          }}
        >
          Clear filters
        </button>
      </div>
      {isLoading || filteredFiles.length === 0 ? (
        <div className="rounded-box border border-base-300 text-center text-sm text-base-content/60 py-10 px-5">
          {searchQuery !== "" ? (
            <p>No files found with such filters.</p>
          ) : (
            <p>No files have been uploaded yet.</p>
          )}
        </div>
      ) : (
        <ul className="list bg-base-100 rounded-box border border-base-300">
          {filteredFiles.map((file) => (
            <File key={file.id} file={file} />
          ))}
        </ul>
      )}
    </div>
  );
}
