"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Sale, ApiResponse, User } from "@/lib/types";

export default function SalesPage() {
  const router = useRouter();
  const { error: showError } = useToast();
  const [user, setUser] = useState<User | null>(null);
  const [sales, setSales] = useState<Sale[]>([]);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<"" | "completed" | "pending">("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetch = async () => {
      try {
        const res = await apiClient.get<ApiResponse<Sale[]>>("/sales");
        if (Array.isArray(res.data)) setSales(res.data);
      } catch (err) {
        const errMsg = err instanceof Error ? err.message : "Failed to load";
        setError(errMsg);
        showError(errMsg);
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, []);

  const filtered = sales.filter(s => {
    const matchSearch =
      s.invoice_number.toLowerCase().includes(search.toLowerCase()) ||
      (s.customer?.name || "").toLowerCase().includes(search.toLowerCase());
    const matchStatus = statusFilter === "" || s.status === statusFilter;
    return matchSearch && matchStatus;
  });

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <ToastContainer />
      <main className="flex-1 p-8">
        <div className="flex justify-between items-center mb-8 bg-white rounded-lg shadow-sm p-6">
          <h1 className="text-3xl font-bold">Sales</h1>
          <button
            onClick={() => router.push("/sales/new")}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded font-medium transition"
          >
            + New Sale
          </button>
        </div>

        {error && <div className="text-red-600 mb-4 bg-red-50 p-4 rounded">{error}</div>}

        <div className="mb-6 bg-white rounded-lg shadow-sm p-4 space-y-4">
          <input
            type="text"
            placeholder="Search by invoice or customer..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <div className="flex gap-2">
            <button
              onClick={() => setStatusFilter("")}
              className={`px-4 py-2 rounded font-medium transition ${
                statusFilter === ""
                  ? "bg-blue-600 text-white"
                  : "bg-gray-200 text-gray-700 hover:bg-gray-300"
              }`}
            >
              All
            </button>
            <button
              onClick={() => setStatusFilter("completed")}
              className={`px-4 py-2 rounded font-medium transition ${
                statusFilter === "completed"
                  ? "bg-green-600 text-white"
                  : "bg-gray-200 text-gray-700 hover:bg-gray-300"
              }`}
            >
              Completed
            </button>
            <button
              onClick={() => setStatusFilter("pending")}
              className={`px-4 py-2 rounded font-medium transition ${
                statusFilter === "pending"
                  ? "bg-yellow-600 text-white"
                  : "bg-gray-200 text-gray-700 hover:bg-gray-300"
              }`}
            >
              Pending
            </button>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Invoice</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Customer</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Total</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Status</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Date</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {filtered.map((sale) => (
                <tr key={sale.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium">{sale.invoice_number}</td>
                  <td className="px-6 py-4">{sale.customer?.name || "—"}</td>
                  <td className="px-6 py-4">${(Number(sale.total_amount) || 0).toFixed(2)}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`text-xs font-semibold px-2 py-1 rounded ${
                        sale.status === "completed"
                          ? "bg-green-100 text-green-700"
                          : "bg-yellow-100 text-yellow-700"
                      }`}
                    >
                      {sale.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {new Date(sale.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 space-x-2">
                    <button
                      onClick={() => router.push(`/sales/${sale.id}`)}
                      className="text-blue-600 hover:underline text-sm"
                    >
                      View
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {filtered.length === 0 && (
            <div className="p-8 text-center text-gray-500">
              {search || statusFilter ? "No sales match filters" : "No sales found"}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
