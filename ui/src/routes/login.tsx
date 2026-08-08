import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { authQueryOptions, useAuth } from "#/hooks/use-auth";

export const Route = createFileRoute("/login")({
	validateSearch: (search) => ({
		redirect: (search.redirect as string) || "/",
	}),
	beforeLoad: async ({ context, search }) => {
		const user = await context.queryClient.ensureQueryData(authQueryOptions);

		if (user) {
			throw redirect({ to: search.redirect });
		}
	},
	component: RouteComponent,
});

function RouteComponent() {
	const { redirect } = Route.useSearch();
	const navigate = useNavigate();
	const { loginMutation } = useAuth();

	const handleLogin = async () => {
		loginMutation.mutate(
			{ username: "admin", password: "admin" },
			{
				onSuccess: () => navigate({ to: redirect }),
			},
		);
	};

	return (
		<div>
			<button
				type="button"
				onClick={handleLogin}
				disabled={loginMutation.isPending}
			>
				{loginMutation.isPending ? "Logging in..." : "Log in"}
			</button>
			{loginMutation && <p>{loginMutation.error?.message}</p>}
		</div>
	);
}
