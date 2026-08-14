import { createFileRoute } from "@tanstack/react-router";
import { PageTitle } from "#/components/PageTitle";

export const Route = createFileRoute("/_authenticated/files")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div>
      <PageTitle title="Files" description="n files - n bytes" />
    </div>
  );
}
