import { User } from "@/lib/types";

interface HeaderProps {
  title: string;
  user: User | null;
  action?: React.ReactNode;
}

const roleLabel = (role: User["role"]) =>
  typeof role === "string" ? role : role?.name;

export default function Header({ title, user, action }: HeaderProps) {
  return (
    <header className="mb-6 flex flex-wrap items-center justify-between gap-4 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="min-w-0">
        <h1 className="text-2xl font-bold tracking-tight text-slate-900">
          {title}
        </h1>
        {user && (
          <p className="mt-0.5 truncate text-sm text-slate-500">
            Welcome, {user.name}
          </p>
        )}
      </div>

      <div className="flex items-center gap-4">
        {action && <div>{action}</div>}
        {user && (
          <div className="flex items-center gap-3">
            <div className="flex flex-col items-end gap-1">
              <p className="text-sm text-slate-600">{user.email}</p>
              <span className="inline-block rounded-full bg-indigo-50 px-2.5 py-0.5 text-xs font-semibold capitalize leading-none text-indigo-700">
                {roleLabel(user.role)}
              </span>
            </div>
            <span
              className="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-indigo-500 text-sm font-semibold uppercase text-white"
              aria-hidden
            >
              {user.name?.slice(0, 1) ?? "?"}
            </span>
          </div>
        )}
      </div>
    </header>
  );
}