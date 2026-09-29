"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import { apiClient } from "@/lib/api";
import { Product, ApiResponse, User } from "@/lib/types";

export default function ProductsPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetchProducts = async () => {
      try {
        const res = await apiClient.get<ApiResponse<Product[]>>("/products");
        if (Array.isArray(res.data)) {
          setProducts(res.data);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load products");
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [router]);

  const handleDelete = async (id: number) => {
    if (!confirm("Delete product?")) return;
    try {
      await apiClient.delete(`/products/${id}`);
      setProducts(products.filter(p => p.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Delete failed");
    }
  };

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  const isAdmin = (typeof user.role === "string" ? user.role : user.role?.name) === "admin";

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />

      <main className="flex-1 p-8">
        <div className="flex justify-between items-center mb-8 bg-white rounded-lg shadow-sm p-6">
          <h1 className="text-3xl font-bold">Products</h1>
          {isAdmin && (
            <button
              onClick={() => router.push("/products/new")}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded font-medium transition"
            >
              + New Product
            </button>
          )}
        </div>

        {error && <div className="text-red-600 mb-4 bg-red-50 p-4 rounded">{error}</div>}

        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase tracking-wide">
                  Name
                </th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase tracking-wide">
                  SKU
                </th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase tracking-wide">
                  Price
                </th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase tracking-wide">
                  Stock
                </th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase tracking-wide">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {products.map((product) => (
                <tr key={product.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{product.name}</td>
                  <td className="px-6 py-4 text-gray-600">{product.sku}</td>
                  <td className="px-6 py-4">${product.price}</td>
                  <td className="px-6 py-4">
                    <span
                      className={
                        product.stock <= product.minimum_stock
                          ? "text-red-600 font-semibold"
                          : "text-green-600 font-semibold"
                      }
                    >
                      {product.stock}
                    </span>
                  </td>
                  <td className="px-6 py-4 space-x-2">
                    {isAdmin && (
                      <>
                        <button
                          onClick={() => router.push(`/products/${product.id}`)}
                          className="text-blue-600 hover:underline text-sm"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDelete(product.id)}
                          className="text-red-600 hover:underline text-sm"
                        >
                          Delete
                        </button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {products.length === 0 && (
            <div className="p-8 text-center text-gray-500">No products found</div>
          )}
        </div>
      </main>
    </div>
  );
}
