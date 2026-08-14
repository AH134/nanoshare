import { Files, Home, type LucideIcon, Settings } from "lucide-react";

export interface NavLink {
  to: string;
  label: string;
  icon: LucideIcon;
  exact?: boolean;
}

export const navLinks: NavLink[] = [
  { to: "/", label: "Upload", icon: Home, exact: true },
  { to: "/files", label: "Files", icon: Files },
  { to: "/settings", label: "Settings", icon: Settings },
];
