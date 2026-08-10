import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import Loading from "#/components/Loading";
import { authQueryOptions } from "#/hooks/use-auth";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async ({ context, location }) => {
    const user = await context.queryClient.ensureQueryData(authQueryOptions);

    if (!user) {
      throw redirect({
        to: "/login",
        search: { redirect: location.href },
      });
    }

    return { user };
  },
  pendingComponent: Loading,
  component: () => <Outlet />,
});
