import { useNavigate } from "@tanstack/react-router";
import { TriangleAlert } from "lucide-react";
import { useState } from "react";
import { useAuth } from "#/hooks/use-auth";
import { PasswordInput } from "./PasswordInput";

export function ChangePasswordForm() {
  const { passwordMutation } = useAuth();
  const navigate = useNavigate();

  const [newPasswordError, setNewPasswordError] = useState("");

  const verifyNewPassword = (password: string, confirmPassword: string) => {
    if (password !== confirmPassword) {
      setNewPasswordError("New password do not match.");
      return false;
    }
    setNewPasswordError("");
    return true;
  };

  const handleChangePassword = (formData: FormData) => {
    const currentPassword = formData.get("current-password") as string;
    const newPassword = formData.get("new-password") as string;
    const confirmPassword = formData.get("confirm-password") as string;

    if (!verifyNewPassword(newPassword, confirmPassword)) {
      return;
    }

    passwordMutation.mutate(
      { currentPassword, newPassword },
      {
        onSuccess: () => {
          navigate({ to: "/login", search: { redirect: "/" } });
        },
        onError: (err) => {
          console.error(err);
        },
      },
    );
  };

  return (
    <form action={handleChangePassword}>
      {(newPasswordError || passwordMutation.isError) && (
        <div role="alert" className="alert alert-error alert-soft mb-4">
          <TriangleAlert className="size-5" />
          <span>
            {newPasswordError ||
              passwordMutation.error?.message ||
              "Something went wrong."}
          </span>
        </div>
      )}
      <fieldset className="fieldset">
        <legend className="sr-only">Change password</legend>
        <label
          className="label text-base-content text-sm mt-2"
          htmlFor="current-password"
        >
          Current Password
        </label>
        <PasswordInput
          id="current-password"
          name="current-password"
          placeholder="Enter your current password"
          required
        />

        <label
          className="label text-base-content text-sm mt-2"
          htmlFor="new-password"
        >
          New Password
        </label>
        <PasswordInput
          id="new-password"
          name="new-password"
          placeholder="Enter your new password"
          required
        />

        <label
          className="label text-base-content text-sm mt-2"
          htmlFor="confirm-password"
        >
          Confirm Password
        </label>
        <PasswordInput
          id="confirm-password"
          name="confirm-password"
          placeholder="Confirm your new password"
          required
        />
        <button
          type="submit"
          className="btn btn-primary w-fit rounded mt-4"
          disabled={passwordMutation.isPending}
        >
          {passwordMutation.isPending
            ? "Updating password..."
            : "Update password"}
        </button>
      </fieldset>
    </form>
  );
}
