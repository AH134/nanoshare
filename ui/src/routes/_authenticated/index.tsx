import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuth } from "#/hooks/use-auth";
import { useTheme } from "#/hooks/use-theme";

export const Route = createFileRoute("/_authenticated/")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const { user, logoutMutation } = useAuth();
  const { theme, setTheme } = useTheme();

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => navigate({ to: "/login", search: { redirect: "/" } }),
    });
  };
  return (
    <div>
      Hello {user?.username}!
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
