"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Product, ApiResponse, User } from "@/lib/types";

export default function ProductFormPage() {
  const router = useRouter();
  const params = useParams();
  const { success, error: showError } = useToast();
  const id = params?.id as string;
  const isEdit = !!id && id !== "new";

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [categories, setCategories] = useState<{ id: number; name: string }[]>([]);

  const [form, setForm] = useState({
    name: "",
    sku: "",
    category_id: 0,
    description: "",
    price: 0,
    stock: 0,
    minimum_stock: 0,
  });

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    const userData = JSON.parse(stored);
    setUser(userData);

    const isAdmin = (typeof userData.role === "string" ? userData.role : userData.role?.name) === "admin";
    if (!isAdmin) {
      router.push("/products");
      return;
    }

    const fetchCategories = async () => {
      try {
        const res = await apiClient.get<ApiResponse<{ id: number; name: string }[]>>(
          "/categories"
        );
        if (Array.isArray(res.data)) {
          setCategories(res.data);
        }
      } catch (err) {
        console.error(err);
      }
    };

    fetchCategories();

    if (isEdit) {
      const fetchProduct = async () => {
        try {
          const res = await apiClient.get<ApiResponse<Product>>(`/products/${id}`);
          if (res.data) {
            setForm(res.data);
          }
        } catch (err) {
          showError(err instanceof Error ? err.message : "Failed to load product");
        } finally {
          setLoading(false);
        }
      };
      fetchProduct();
    }
  }, [id, isEdit]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setForm(prev => ({
      ...prev,
      [name]: ["price", "stock", "minimum_stock", "category_id"].includes(name) 
        ? parseInt(value) 
        : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (form.category_id === 0) {
      setError("Please select a category");
      showError("Please select a category");
      return;
    }

    setSubmitting(true);
    setError("");

    try {
      if (isEdit) {
        await apiClient.put(`/products/${id}`, form);
        success("Product updated successfully");
      } else {
        await apiClient.post("/products", form);
        success("Product created successfully");
      }
      router.push("/products");
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : "Submission failed";
      setError(errMsg);
      showError(errMsg);
      setSubmitting(false);
    }
  };

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <ToastContainer />

      <main className="flex-1 p-8 max-w-2xl">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">{isEdit ? "Edit Product" : "New Product"}</h1>
          <p className="text-gray-600 mt-2">
            {isEdit ? "Update product details" : "Create a new product"}
          </p>
        </div>

        {error && <div className="text-red-600 mb-6 bg-red-50 p-4 rounded">{error}</div>}

        <div className="bg-white rounded-lg shadow-sm p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Product Name *
              </label>
              <input
                type="text"
                name="name"
                value={form.name}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Laptop Pro 15"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                SKU *
              </label>
              <input
                type="text"
                name="sku"
                value={form.sku}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="ELEC-001"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Category *
              </label>
              <select
                name="category_id"
                value={form.category_id}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value={0}>Select category</option>
                {categories.map(cat => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
              {form.category_id === 0 && <p className="text-red-600 text-xs mt-1">Category required</p>}
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Description
              </label>
              <textarea
                name="description"
                value={form.description}
                onChange={handleChange}
                rows={4}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Product description..."
              />
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Price $ *
                </label>
                <input
                  type="number"
                  name="price"
                  value={form.price}
                  onChange={handleChange}
                  required
                  step="0.01"
                  min="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="0.00"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Stock *
                </label>
                <input
                  type="number"
                  name="stock"
                  value={form.stock}
                  onChange={handleChange}
                  required
                  min="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="0"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Min Stock *
                </label>
                <input
                  type="number"
                  name="minimum_stock"
                  value={form.minimum_stock}
                  onChange={handleChange}
                  required
                  min="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="0"
                />
              </div>
            </div>

            <div className="flex gap-4 pt-6">
              <button
                type="submit"
                disabled={submitting || form.category_id === 0}
                className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-2 px-4 rounded-lg font-medium transition"
              >
                {submitting ? "Saving..." : isEdit ? "Update Product" : "Create Product"}
              </button>
              <button
                type="button"
                onClick={() => router.push("/products")}
                className="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-800 py-2 px-4 rounded-lg font-medium transition"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
