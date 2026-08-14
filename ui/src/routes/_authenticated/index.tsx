import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { PageTitle } from "#/components/PageTitle";
import { UploadZone } from "#/components/UploadZone";
import { useAuth } from "#/hooks/use-auth";
import { useTheme } from "#/hooks/use-theme";

export const Route = createFileRoute("/_authenticated/")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const { logoutMutation } = useAuth();
  const { theme, setTheme } = useTheme();

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => navigate({ to: "/login", search: { redirect: "/" } }),
    });
  };
  return (
    <div>
      <PageTitle
        title="Upload your files"
        description="Drag files into the zone below or click to browse. Set a download limit
        and expiry date, then share the link."
      />
      <button
        type="button"
        onClick={handleLogout}
        disabled={logoutMutation.isPending}
      >
        Log out
      </button>
      <div>hello theme: {theme}</div>
      <div className="flex gap-1">
        <button type="button" onClick={() => setTheme("light")}>
          light
        </button>
        <button type="button" onClick={() => setTheme("dark")}>
          dark
        </button>
        <button type="button" onClick={() => setTheme("system")}>
          system
        </button>
      </div>
    </div>
  );
}
