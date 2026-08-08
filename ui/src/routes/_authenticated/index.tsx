import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuth } from "#/hooks/use-auth";

export const Route = createFileRoute("/_authenticated/")({
	component: RouteComponent,
});

function RouteComponent() {
	const navigate = useNavigate();
	const { user, logoutMutation } = useAuth();

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
		</div>
	);
}
