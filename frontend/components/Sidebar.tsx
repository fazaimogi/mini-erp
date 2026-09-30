import Link from "next/link";
import { useRouter } from "next/navigation";
import { User } from "@/lib/types";

interface SidebarProps {
  user: User | null;
}

export default function Sidebar({ user }: SidebarProps) {
  const router = useRouter();

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    router.push("/login");
  };

  const isAdmin = (typeof user?.role === "string" ? user.role : user?.role?.name) === "admin";
  return (
    <aside className="w-64 bg-slate-900 text-white p-5 flex flex-col min-h-screen">
      <h2 className="text-2xl font-bold mb-8">ERP</h2>

      <nav className="flex-1 space-y-2">
        <Link href="/dashboard" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
          Dashboard
        </Link>
        <Link href="/products" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
          Products
        </Link>
        {isAdmin && (
          <Link href="/categories" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
            Categories
          </Link>
        )}
        <Link href="/inventory" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
          Inventory
        </Link>
        <Link href="/customers" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
          Customers
        </Link>
        <Link href="/sales" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
          Sales
        </Link>
        {isAdmin && (
          <Link href="/users" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
            Users
          </Link>
        )}
        <Link href="/reports" className="block px-3 py-2 rounded hover:bg-slate-800 transition">
          Reports
        </Link>
      </nav>

      <button
        onClick={handleLogout}
        className="w-full bg-red-600 hover:bg-red-700 py-2 px-4 rounded text-sm font-medium transition"
      >
        Logout
      </button>
    </aside>
  );
}
