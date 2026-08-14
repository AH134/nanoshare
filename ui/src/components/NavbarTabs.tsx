import { Link } from "@tanstack/react-router";
import { navLinks } from "./nav-links";

export function NavbarTabs() {
  const activeTab = { className: "tab tab-active text-primary" };

  return (
    <div role="tablist" className="tabs tabs-border">
      {navLinks.map(({ to, label, exact }) => (
        <Link
          key={to}
          to={to}
          role="tab"
          className="tab"
          activeProps={activeTab}
          activeOptions={exact ? { exact: true } : undefined}
        >
          {label}
        </Link>
      ))}
    </div>
  );
}
