import { Link } from "@tanstack/react-router";
import { NavbarTabs } from "./NavbarTabs";
import { NavbarUserMenu } from "./NavbarUserMenu";

export function Header() {
  return (
    <div className="bg-base-100 border-b border-base-300">
      <div className="navbar max-w-4xl mx-auto">
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
