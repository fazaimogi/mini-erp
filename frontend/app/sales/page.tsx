"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import { apiClient } from "@/lib/api";
import { Sale, ApiResponse, User } from "@/lib/types";

export default function SalesPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [sales, setSales] = useState<Sale[]>([]);
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
        setError(err instanceof Error ? err.message : "Failed to load");
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, [router]);

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
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
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Invoice</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Amount</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Status</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Date</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {sales.map((sale) => (
                <tr key={sale.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{sale.invoice_number}</td>
                  <td className="px-6 py-4 font-semibold">${(Number(sale.total_amount) || 0).toFixed(2)}</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs font-semibold px-2 py-1 rounded ${
                      sale.status === "completed" ? "bg-green-100 text-green-700" : "bg-yellow-100 text-yellow-700"
                    }`}>
                      {sale.status}
                    </span>
                  </td>
                  <td className="px-6 py-4">{new Date(sale.created_at).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
          {sales.length === 0 && <div className="p-8 text-center text-gray-500">No sales</div>}
        </div>
      </main>
    </div>
  );
}
