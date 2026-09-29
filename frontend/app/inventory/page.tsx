"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import { apiClient } from "@/lib/api";
import { Product, ApiResponse, User } from "@/lib/types";

export default function InventoryPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [inventory, setInventory] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState<number | null>(null);

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetch = async () => {
      try {
        const res = await apiClient.get<ApiResponse<Product[]>>("/inventory");
        if (Array.isArray(res.data)) setInventory(res.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load");
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, []);

  const handleStockIn = async (productId: number) => {
    const qty = prompt("Enter quantity to add:");
    if (!qty || isNaN(parseInt(qty))) return;

    setSubmitting(productId);
    try {
      await apiClient.post(`/inventory/stock-in`, {
        product_id: productId,
        quantity: parseInt(qty),
      });

      setInventory(inventory.map(item =>
        item.id === productId
          ? { ...item, stock: item.stock + parseInt(qty) }
          : item
      ));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Stock in failed");
    } finally {
      setSubmitting(null);
    }
  };

  const handleStockOut = async (productId: number) => {
    const qty = prompt("Enter quantity to remove:");
    if (!qty || isNaN(parseInt(qty))) return;

    setSubmitting(productId);
    try {
      await apiClient.post(`/inventory/stock-out`, {
        product_id: productId,
        quantity: parseInt(qty),
      });

      setInventory(inventory.map(item =>
        item.id === productId
          ? { ...item, stock: Math.max(0, item.stock - parseInt(qty)) }
          : item
      ));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Stock out failed");
    } finally {
      setSubmitting(null);
    }
  };

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <main className="flex-1 p-8">
        <Header title="Inventory Management" user={user} />

        {error && <div className="text-red-600 mb-4 bg-red-50 p-4 rounded">{error}</div>}

        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Product</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Stock</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Min</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Status</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {inventory.map((item) => (
                <tr key={item.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{item.name}</td>
                  <td className="px-6 py-4 font-semibold">{item.stock}</td>
                  <td className="px-6 py-4">{item.minimum_stock}</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs font-semibold px-2 py-1 rounded ${
                      item.stock <= item.minimum_stock
                        ? "bg-red-100 text-red-700"
                        : "bg-green-100 text-green-700"
                    }`}>
                      {item.stock <= item.minimum_stock ? "Low" : "OK"}
                    </span>
                  </td>
                  <td className="px-6 py-4 space-x-2">
                    <button
                      onClick={() => handleStockIn(item.id)}
                      disabled={submitting === item.id}
                      className="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white px-3 py-1 rounded text-sm transition"
                    >
                      {submitting === item.id ? "..." : "Stock In"}
                    </button>
                    <button
                      onClick={() => handleStockOut(item.id)}
                      disabled={submitting === item.id}
                      className="bg-orange-600 hover:bg-orange-700 disabled:bg-gray-400 text-white px-3 py-1 rounded text-sm transition"
                    >
                      {submitting === item.id ? "..." : "Stock Out"}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {inventory.length === 0 && <div className="p-8 text-center text-gray-500">No items</div>}
        </div>
      </main>
    </div>
  );
}
