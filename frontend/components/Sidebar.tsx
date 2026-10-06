import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { User } from "@/lib/types";

interface SidebarProps {
  user: User | null;
}

const icon = {
  fill: "none",
  stroke: "currentColor",
  strokeWidth: 1.8,
  strokeLinecap: "round" as const,
  strokeLinejoin: "round" as const,
};

const NAV = [
  { href: "/dashboard", label: "Dashboard", d: "M3 12 12 4l9 8M5 10v10h14V10" },
  { href: "/products", label: "Products", d: "M3 7l9-4 9 4-9 4-9-4Zm0 0v10l9 4 9-4V7" },
  { href: "/inventory", label: "Inventory", d: "M4 6h16v4H4V6Zm0 8h16v4H4v-4ZM8 6v12" },
  { href: "/customers", label: "Customers", d: "M16 19v-1a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v1M9.5 10a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm11 9v-1a4 4 0 0 0-3-3.9M16.5 4.2a3 3 0 0 1 0 5.6" },
  { href: "/sales", label: "Sales", d: "M6 2h9l4 4v16H6V2Zm3 7h6M9 13h6M9 17h3" },
  { href: "/reports", label: "Reports", d: "M4 20V10m5 10V4m5 16v-7m5 7V8" },
] as const;

const ADMIN_NAV = [
  { href: "/categories", label: "Categories", d: "M4 5h7v7H4V5Zm9 0h7v7h-7V5ZM4 14h7v7H4v-7Zm9 0h7v7h-7v-7Z" },
  { href: "/vouchers", label: "Vouchers", d: "M3 7h18v3a2 2 0 0 0 0 4v3H3v-3a2 2 0 0 0 0-4V7Zm7 0v10" },
  { href: "/users", label: "Users", d: "M15 19v-1a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v1M9 10a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm12 9v-1a4 4 0 0 0-3-3.9M18 4.2a3 3 0 0 1 0 5.6" },
  { href: "/activity-logs", label: "Activity Logs", d: "M12 8v4l3 2M3 12a9 9 0 1 0 9-9 9 9 0 0 0-9 9Zm0 0H1m2 0 2-3" },
] as const;

export default function Sidebar({ user }: SidebarProps) {
  const router = useRouter();
  const pathname = usePathname();

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    router.push("/login");
  };

  const role = typeof user?.role === "string" ? user.role : user?.role?.name;
  const isAdmin = role === "admin";

  const NavItem = ({
    href,
    label,
    d,
  }: {
    href: string;
    label: string;
    d: string;
  }) => {
    const active = pathname === href;
    return (
      <Link
        href={href}
        aria-current={active ? "page" : undefined}
        className={`group flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition ${
          active
            ? "bg-indigo-500/15 text-white"
            : "text-slate-400 hover:bg-white/5 hover:text-white"
        }`}
      >
        <svg
          viewBox="0 0 24 24"
          className={`h-[18px] w-[18px] shrink-0 transition-colors ${
            active ? "text-indigo-300" : "text-slate-500 group-hover:text-indigo-300"
          }`}
          {...icon}
          aria-hidden
        >
          <path d={d} />
        </svg>
        {label}
        {active && (
          <span className="ml-auto h-1.5 w-1.5 rounded-full bg-indigo-300" aria-hidden />
        )}
      </Link>
    );
  };

  return (
    <aside className="flex min-h-screen w-64 shrink-0 flex-col bg-slate-900 p-5 text-white">
      <div className="mb-8 flex items-center gap-2.5">
        <span className="grid h-8 w-8 place-items-center rounded-lg bg-indigo-500 text-sm font-bold">
          E
        </span>
        <span className="text-lg font-bold tracking-tight">ERP</span>
      </div>

      <nav className="flex-1 space-y-1" aria-label="Main">
        {NAV.map(item => (
          <NavItem key={item.href} {...item} />
        ))}

        {isAdmin && (
          <>
            <p className="px-3 pb-1 pt-5 text-[11px] font-semibold uppercase tracking-wider text-slate-500">
              Admin
            </p>
            {ADMIN_NAV.map(item => (
              <NavItem key={item.href} {...item} />
            ))}
          </>
        )}
      </nav>

      {user && (
        <div className="mt-4 flex items-center gap-3 rounded-lg border border-white/10 bg-white/5 px-3 py-2">
          <span
            className="grid h-8 w-8 shrink-0 place-items-center rounded-full bg-indigo-500 text-xs font-semibold uppercase text-white"
            aria-hidden
          >
            {user.name?.slice(0, 1) ?? "?"}
          </span>
          <div className="min-w-0">
            <p className="truncate text-sm font-medium text-white">{user.name}</p>
            <p className="truncate text-xs capitalize text-slate-400">{role}</p>
          </div>
        </div>
      )}

      <button
        onClick={handleLogout}
        className="mt-3 w-full rounded-lg border border-rose-500/30 px-4 py-2 text-sm font-medium text-rose-300 transition hover:border-rose-500 hover:bg-rose-500 hover:text-white"
      >
        Logout
      </button>
    </aside>
  );
}