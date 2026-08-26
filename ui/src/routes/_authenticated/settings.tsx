import { createFileRoute } from "@tanstack/react-router";
import { ChangePasswordForm } from "#/components/ChangePasswordForm";
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
          <ChangePasswordForm />
        </SettingsCard>
      </div>
    </div>
  );
}
