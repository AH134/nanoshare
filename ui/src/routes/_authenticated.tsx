import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { Footer } from "#/components/Footer";
import { Header } from "#/components/Header";
import { Loading } from "#/components/Loading";
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
  component: AuthenticatedComponent,
});

function AuthenticatedComponent() {
  return (
    <div className="min-h-dvh bg-base-200 grid grid-rows-[auto_1fr_auto]">
      <Header />
      <div className="max-w-3xl mx-auto px-2 py-12 w-full">
        <Outlet />
      </div>
      <Footer />
    </div>
  );
}
