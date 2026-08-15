import type React from "react";

export interface SettingsCardProps {
  title: string;
  description: string;
  children: React.ReactNode;
}
export function SettingsCard({
  title,
  description,
  children,
}: SettingsCardProps) {
  return (
    <div className="card bg-base-100 w-full border border-base-300">
      <div className="card-body">
        <div className="mb-4">
          <h2 className="card-title">{title}</h2>
          <p className="text-base-content/60">{description}</p>
        </div>
        {children}
      </div>
    </div>
  );
}
