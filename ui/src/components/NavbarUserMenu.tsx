import { Link, useNavigate } from "@tanstack/react-router";
import { LogOut } from "lucide-react";
import { useRef } from "react";
import { useAuth } from "#/hooks/use-auth";
import { navLinks } from "./nav-links";

export function NavbarUserMenu() {
  const { user, logoutMutation } = useAuth();
  const navigate = useNavigate();
  const popoverRef = useRef<HTMLUListElement>(null);

  const handleHideMenu = () => {
    popoverRef.current?.hidePopover();
  };

  const handleLogout = async () => {
    try {
      await logoutMutation.mutateAsync();
    } catch (err) {
      console.error("Logout failed:", err);
    } finally {
      handleHideMenu();
      navigate({ to: "/login", search: { redirect: "/" } });
    }
  };

  return (
    <>
      <button
        type="button"
        className="rounded-full"
        popoverTarget="user-menu-popover"
        style={{ anchorName: "--user-menu-anchor" } as React.CSSProperties}
      >
        <div className="avatar avatar-placeholder">
          <div className="bg-accent/80 text-accent-content w-8 rounded-full">
            <span className="text-sm font-semibold capitalize">
              {user?.username[0]}
            </span>
          </div>
        </div>
      </button>

      <ul
        className="dropdown dropdown-end menu w-52 rounded-box bg-base-100 p-2 shadow-sm"
        ref={popoverRef}
        popover="auto"
        id="user-menu-popover"
        style={{ positionAnchor: "--user-menu-anchor" } as React.CSSProperties}
      >
        {navLinks.map(({ to, label, icon: Icon }) => (
          <li key={to}>
            <Link to={to} onClick={handleHideMenu}>
              <Icon className="size-4" />
              <span>{label}</span>
            </Link>
          </li>
        ))}
        <div className="divider m-0"></div>
        <li>
          <button type="button" onClick={handleLogout}>
            <LogOut className="size-4" /> Log out
          </button>
        </li>
      </ul>
    </>
  );
}
