import { User } from "@/lib/types";

interface HeaderProps {
  title: string;
  user: User | null;
  action?: React.ReactNode;
}

export default function Header({ title, user, action }: HeaderProps) {
  return (
    <div className="bg-white rounded-lg shadow-sm p-6 mb-6 flex justify-between items-center">
      <div>
        <h1 className="text-3xl font-bold">{title}</h1>
        {user && <p className="text-gray-600 text-sm mt-1">Welcome, {user.name}</p>}
      </div>
      <div className="text-right">
        {user && (
          <>
            <p className="text-sm text-gray-600">{user.email}</p>
            <span className="inline-block bg-blue-600 text-white px-3 py-1 rounded text-xs font-semibold mt-2">
              {user.role.name}
            </span>
          </>
        )}
      </div>
      {action && <div>{action}</div>}
    </div>
  );
}
