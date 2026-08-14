import { createFileRoute } from "@tanstack/react-router";
import { PageTitle } from "#/components/PageTitle";

export const Route = createFileRoute("/_authenticated/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div>
      <PageTitle
        title="Settings"
        description="Manage your appearance and account preferences."
      />
    </div>
  );
}
