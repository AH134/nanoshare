import { createFileRoute } from "@tanstack/react-router";
import { PageTitle } from "#/components/PageTitle";
import { UploadZone } from "#/components/UploadZone";

export const Route = createFileRoute("/_authenticated/")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div>
      <PageTitle
        title="Upload your files"
        description="Drag files into the zone below or click to browse. Set a download limit
        and expiry date, then share the link."
      />
      <UploadZone />
    </div>
  );
}
