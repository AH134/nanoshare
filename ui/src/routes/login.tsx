import { createFileRoute, redirect } from "@tanstack/react-router";
import Loading from "#/components/Loading";
import LoginForm from "#/components/LoginForm";
import { authQueryOptions } from "#/hooks/use-auth";

export const Route = createFileRoute("/login")({
  validateSearch: (search): { redirect: string } => ({
    redirect: (search.redirect as string) || "/",
  }),
  beforeLoad: async ({ context, search }) => {
    const user = await context.queryClient.ensureQueryData(authQueryOptions);

    if (user) {
      throw redirect({ to: search.redirect });
    }
  },
  pendingComponent: Loading,
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="hero bg-base-200 min-h-screen">
      <div className="hero-content text-center">
        <div className="max-w-md">
          <h1 className="text-5xl font-extrabold">NANOSHARE</h1>
          <div className="p-5 mt-2 mb-1">
            <p className="text-xl ">Log in to your account</p>
            <p className="text-md text-base-content/60 mt-1">
              Welcome back! Please enter your details
            </p>
          </div>
          <LoginForm />
        </div>
      </div>
    </div>
  );
}
