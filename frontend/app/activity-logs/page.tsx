"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import { apiClient } from "@/lib/api";
import { ActivityLog, ApiResponse, User } from "@/lib/types";

export default function ActivityLogsPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [logs, setLogs] = useState<ActivityLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [limit] = useState(20);
  
  const [filterModule, setFilterModule] = useState("");
  const [filterAction, setFilterAction] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    const userData = JSON.parse(stored);
    const role = typeof userData.role === "string" ? userData.role : userData.role?.name;
    if (role !== "admin") {
      router.push("/dashboard");
      return;
    }
    setUser(userData);
  }, [router]);

  useEffect(() => {
    if (!user) return;
    fetchLogs();
  }, [user, page, filterModule, filterAction]);

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const params = new URLSearchParams({
        page: page.toString(),
        limit: limit.toString(),
      });
      if (filterModule) params.append("module", filterModule);
      if (filterAction) params.append("action", filterAction);

      const res = await apiClient.get<ApiResponse<ActivityLog[]>>(
        `/activity-logs?${params.toString()}`
      );
      
      if (Array.isArray(res.data)) {
        setLogs(res.data);
        setTotal(res.meta?.total || 0);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load logs");
    } finally {
      setLoading(false);
    }
  };

  const totalPages = Math.ceil(total / limit);

  if (!user) return null;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <div className="flex-1 flex flex-col">
        <Header user={user} title="Activity Logs" />
        <main className="flex-1 p-6">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="mb-4 flex gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Module</label>
                <select
                  value={filterModule}
                  onChange={(e) => { setFilterModule(e.target.value); setPage(1); }}
                  className="px-3 py-2 border rounded"
                >
                  <option value="">All Modules</option>
                  <option value="Products">Products</option>
                  <option value="Sales">Sales</option>
                  <option value="Inventory">Inventory</option>
                  <option value="Customers">Customers</option>
                  <option value="Users">Users</option>
                  <option value="Vouchers">Vouchers</option>
                  <option value="Categories">Categories</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Action</label>
                <select
                  value={filterAction}
                  onChange={(e) => { setFilterAction(e.target.value); setPage(1); }}
                  className="px-3 py-2 border rounded"
                >
                  <option value="">All Actions</option>
                  <option value="CREATE">CREATE</option>
                  <option value="UPDATE">UPDATE</option>
                  <option value="DELETE">DELETE</option>
                  <option value="LOGIN">LOGIN</option>
                </select>
              </div>
            </div>

            {error && <div className="bg-red-100 text-red-700 p-3 rounded mb-4">{error}</div>}

            {loading ? (
              <div className="text-center py-8">Loading...</div>
            ) : logs.length === 0 ? (
              <div className="text-center py-8 text-gray-500">No activity logs found</div>
            ) : (
              <>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Time</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">User</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Module</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                        <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">IP</th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {logs.map((log) => (
                        <tr key={log.id} className="hover:bg-gray-50">
                          <td className="px-4 py-3 text-sm text-gray-900 whitespace-nowrap">
                            {new Date(log.created_at).toLocaleString('id-ID')}
                          </td>
                          <td className="px-4 py-3 text-sm text-gray-900">{log.user_name}</td>
                          <td className="px-4 py-3 text-sm">
                            <span className={`px-2 py-1 rounded text-xs font-medium ${
                              log.action === 'CREATE' ? 'bg-green-100 text-green-800' :
                              log.action === 'UPDATE' ? 'bg-blue-100 text-blue-800' :
                              log.action === 'DELETE' ? 'bg-red-100 text-red-800' :
                              'bg-gray-100 text-gray-800'
                            }`}>
                              {log.action}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-sm text-gray-900">{log.module}</td>
                          <td className="px-4 py-3 text-sm text-gray-700">{log.description}</td>
                          <td className="px-4 py-3 text-sm text-gray-500">{log.ip_address}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>

                {totalPages > 1 && (
                  <div className="mt-4 flex items-center justify-between">
                    <div className="text-sm text-gray-700">
                      Showing {(page - 1) * limit + 1} to {Math.min(page * limit, total)} of {total} logs
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                        className="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
                      >
                        Previous
                      </button>
                      <span className="px-3 py-1">
                        Page {page} of {totalPages}
                      </span>
                      <button
                        onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                        disabled={page === totalPages}
                        className="px-3 py-1 border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
                      >
                        Next
                      </button>
                    </div>
                  </div>
                )}
              </>
            )}
          </div>
        </main>
      </div>
    </div>
  );
}
