"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Sidebar from "@/components/Sidebar";
import { apiClient } from "@/lib/api";
import { Customer, ApiResponse, User } from "@/lib/types";

export default function CustomersPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [customers, setCustomers] = useState<Customer[]>([]);
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
        const res = await apiClient.get<ApiResponse<Customer[]>>("/customers");
        if (Array.isArray(res.data)) setCustomers(res.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load");
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, [router]);

  const handleDelete = async (id: number) => {
    if (!confirm("Delete customer?")) return;
    try {
      await apiClient.delete(`/customers/${id}`);
      setCustomers(customers.filter(c => c.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Delete failed");
    }
  };

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <main className="flex-1 p-8">
        <div className="flex justify-between items-center mb-8 bg-white rounded-lg shadow-sm p-6">
          <h1 className="text-3xl font-bold">Customers</h1>
          <button
            onClick={() => router.push("/customers/new")}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded font-medium transition"
          >
            + New Customer
          </button>
        </div>
        {error && <div className="text-red-600 mb-4 bg-red-50 p-4 rounded">{error}</div>}
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Name</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Email</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Phone</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {customers.map((cust) => (
                <tr key={cust.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{cust.name}</td>
                  <td className="px-6 py-4">{cust.email}</td>
                  <td className="px-6 py-4">{cust.phone}</td>
                  <td className="px-6 py-4 space-x-2">
                    <button
                      onClick={() => router.push(`/customers/${cust.id}`)}
                      className="text-blue-600 hover:underline text-sm"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(cust.id)}
                      className="text-red-600 hover:underline text-sm"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {customers.length === 0 && <div className="p-8 text-center text-gray-500">No customers</div>}
        </div>
      </main>
    </div>
  );
}
