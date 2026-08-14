import { Link } from "@tanstack/react-router";
import { NavbarTabs } from "./NavbarTabs";
import { NavbarUserMenu } from "./NavbarUserMenu";

export function Header() {
  return (
    <div className="bg-base-100 shadow-sm">
      <div className="navbar max-w-6xl mx-auto">
        <div className="navbar-start">
          <Link to="/" className="text-xl font-semibold">
            Nanoshare
          </Link>
        </div>
        <div className="navbar-center hidden sm:flex font-medium">
          <NavbarTabs />
        </div>
        <div className="navbar-end">
          <NavbarUserMenu />
        </div>
      </div>
    </div>
  );
}
