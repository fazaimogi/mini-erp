"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import DeleteModal from "@/components/DeleteModal";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Product, ApiResponse, User } from "@/lib/types";

export default function ProductsPage() {
  const router = useRouter();
  const { success, error } = useToast();
  const [user, setUser] = useState<User | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [deleteModal, setDeleteModal] = useState<{ open: boolean; id: number; name: string }>({
    open: false,
    id: 0,
    name: "",
  });
  const [deleting, setDeleting] = useState(false);

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
        error(err instanceof Error ? err.message : "Failed to load products");
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, []);

  const handleDeleteClick = (id: number, name: string) => {
    setDeleteModal({ open: true, id, name });
  };

  const handleDeleteConfirm = async () => {
    setDeleting(true);
    try {
      await apiClient.delete(`/products/${deleteModal.id}`);
      setProducts(products.filter(p => p.id !== deleteModal.id));
      success(`Product "${deleteModal.name}" deleted`);
      setDeleteModal({ open: false, id: 0, name: "" });
    } catch (err) {
      error(err instanceof Error ? err.message : "Delete failed");
    } finally {
      setDeleting(false);
    }
  };

  const filtered = products.filter(p =>
    p.name.toLowerCase().includes(search.toLowerCase()) ||
    p.sku.toLowerCase().includes(search.toLowerCase())
  );

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  const isAdmin = (typeof user.role === "string" ? user.role : user.role?.name) === "admin";

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <ToastContainer />

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

        <div className="mb-6 bg-white rounded-lg shadow-sm p-4">
          <input
            type="text"
            placeholder="Search by name or SKU..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

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
              {filtered.map((product) => (
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
                          onClick={() => handleDeleteClick(product.id, product.name)}
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
          {filtered.length === 0 && (
            <div className="p-8 text-center text-gray-500">
              {search ? "No products match your search" : "No products found"}
            </div>
          )}
        </div>
      </main>

      <DeleteModal
        isOpen={deleteModal.open}
        title="Delete Product"
        message={`Are you sure you want to delete "${deleteModal.name}"? This action cannot be undone.`}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setDeleteModal({ open: false, id: 0, name: "" })}
        isLoading={deleting}
      />
    </div>
  );
}
