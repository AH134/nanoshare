import { Check, type LucideIcon, Monitor, Moon, Sun } from "lucide-react";
import { useTheme } from "#/hooks/use-theme";
import type { Theme } from "#/providers/theme-context";

const THEMES: {
  value: Theme;
  label: string;
  description: string;
  icon: LucideIcon;
}[] = [
  { value: "light", label: "Light", description: "Bright colours", icon: Sun },
  { value: "dark", label: "Dark", description: "Dark colours", icon: Moon },
  {
    value: "system",
    label: "System",
    description: "Match your device",
    icon: Monitor,
  },
];

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  return (
    <fieldset className="grid gap-3 sm:grid-cols-3">
      <legend className="sr-only">Theme</legend>
      {THEMES.map(({ value, label, description, icon: Icon }) => {
        const isSelected = theme === value;
        return (
          <label
            key={value}
            className={`card border outline-none has-focus-visible:ring-2 has-focus-visible:ring-offset-2 has-focus-visible:ring-offset-base-100 has-focus-visible:ring-primary ${isSelected ? "border-primary bg-primary/5" : "border-base-300 hover:bg-base-200 hover:border-primary"}`}
          >
            <input
              type="radio"
              name="theme"
              value={value}
              checked={isSelected}
              onChange={() => setTheme(value)}
              className="sr-only"
            />
            <div className="card-body text-start">
              <div className="flex justify-between">
                <span className="flex items-center justify-center text-primary size-10 rounded-lg bg-primary/10">
                  <Icon className="size-5" />
                </span>
                {isSelected && (
                  <span className="flex size-5 items-center justify-center text-primary-content bg-primary rounded-full">
                    <Check className="size-3" />
                  </span>
                )}
              </div>
              <div>
                <p className="text-sm font-medium">{label}</p>
                <p className="text-xs text-base-content/60">{description}</p>
              </div>
            </div>
          </label>
        );
      })}
    </fieldset>
  );
}
