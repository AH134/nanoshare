import { createFileRoute } from "@tanstack/react-router";
import { PageTitle } from "#/components/PageTitle";
import { SettingsCard } from "#/components/SettingsCard";
import { ThemeToggle } from "#/components/ThemeToggle";

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

      <div className="flex flex-col gap-6">
        <SettingsCard
          title="Appearance"
          description="Select your preferred color theme"
        >
          <ThemeToggle />
        </SettingsCard>
        {/* TODO: backend stuff */}
        <SettingsCard title="Account" description="Change your password.">
          <form action="">
            <fieldset className="fieldset">
              <legend className="sr-only">Change password</legend>
              <label
                className="label text-base-content text-sm mt-2"
                htmlFor="current-password"
              >
                Current Password
              </label>
              <input type="text" className="input w-full" placeholder="" />

              <label
                className="label text-base-content text-sm mt-2"
                htmlFor="new-password"
              >
                New Password
              </label>
              <input type="text" className="input w-full" />

              <label
                className="label text-base-content text-sm mt-2"
                htmlFor="confirm-password"
              >
                Confirm Password
              </label>
              <input type="text" className="input w-full" />
              <button
                type="submit"
                className="btn btn-primary w-fit rounded mt-4"
              >
                Update password
              </button>
            </fieldset>
          </form>
        </SettingsCard>
      </div>
    </div>
  );
}
