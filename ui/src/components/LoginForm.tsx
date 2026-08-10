import { useNavigate, useSearch } from "@tanstack/react-router";
import { TriangleAlert } from "lucide-react";
import { useAuth } from "#/hooks/use-auth";
import PasswordInput from "./PasswordInput";

export default function LoginForm() {
  const { loginMutation } = useAuth();
  const navigate = useNavigate();
  const search = useSearch({ from: "/login" });

  const handleLogin = async (formData: FormData) => {
    const username = formData.get("username") as string;
    const password = formData.get("password") as string;
    console.log(username, password);
    try {
      await loginMutation.mutateAsync(
        { username, password },
        {
          onSuccess: () => {
            navigate({ href: search.redirect });
          },
        },
      );
    } catch {}
  };

  return (
    <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-sm">
      <form action={handleLogin}>
        <div className="card-body">
          {loginMutation.isError && (
            <div role="alert" className="alert alert-error alert-soft mb-4">
              <TriangleAlert className="size-5" />
              <span>
                {loginMutation.error?.message ??
                  "Invalid username or password."}
              </span>
            </div>
          )}
          <fieldset className="fieldset">
            <label className="label text-primary text-sm" htmlFor="username">
              Username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              className="input bg-base-200/40 focus:outline-none"
              placeholder="Enter your username"
              required
            />
            <label
              className="label text-primary text-sm mt-2"
              htmlFor="password"
            >
              Password
            </label>
            <PasswordInput id="password" name="password" required />
            <button
              className="btn btn-primary hover:bg-primary/80 mt-4"
              type="submit"
              disabled={loginMutation.isPending}
            >
              {loginMutation.isPending ? "Logging in..." : "Login"}
            </button>
          </fieldset>
        </div>
      </form>
    </div>
  );
}
