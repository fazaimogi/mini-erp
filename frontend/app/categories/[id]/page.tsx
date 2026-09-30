"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Category, ApiResponse, User } from "@/lib/types";

export default function CategoryFormPage() {
  const router = useRouter();
  const params = useParams();
  const { success, error: showError } = useToast();
  const id = params?.id as string;
  const isEdit = !!id && id !== "new";

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState({
    name: "",
    description: "",
  });

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    const userData = JSON.parse(stored) as User;
    setUser(userData);

    const isAdmin = (typeof userData.role === "string" ? userData.role : userData.role?.name) === "admin";
    if (!isAdmin) {
      showError("Only admins can manage categories");
      router.push("/categories");
      return;
    }

    if (isEdit) {
      const fetchCategory = async () => {
        try {
          const res = await apiClient.get<ApiResponse<Category>>(`/categories/${id}`);
          if (res.data) {
            setForm(res.data);
          }
        } catch (err) {
          showError(err instanceof Error ? err.message : "Failed to load category");
          router.push("/categories");
        } finally {
          setLoading(false);
        }
      };
      fetchCategory();
    } else {
      setLoading(false);
    }
  }, [id, isEdit]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
    setError("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!form.name.trim()) {
      setError("Name is required");
      return;
    }

    try {
      setSubmitting(true);

      if (isEdit) {
        await apiClient.put(`/categories/${id}`, form);
        success("Category updated successfully");
      } else {
        await apiClient.post("/categories", form);
        success("Category created successfully");
      }

      router.push("/categories");
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to save category";
      setError(msg);
      showError(msg);
    } finally {
      setSubmitting(false);
    }
  };

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <div className="flex-1">
        <Header title={isEdit ? "Edit Category" : "New Category"} user={user} />
        <div className="p-8">
          <div className="max-w-2xl">
            <div className="mb-6">
              <h1 className="text-3xl font-bold text-gray-900">
                {isEdit ? "Edit Category" : "New Category"}
              </h1>
            </div>

            <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6">
              {error && (
                <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded text-red-700">
                  {error}
                </div>
              )}

              <div className="mb-6">
                <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
                  Category Name *
                </label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={form.name}
                  onChange={handleChange}
                  placeholder="e.g., Electronics"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={submitting}
                />
              </div>

              <div className="mb-6">
                <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
                  Description
                </label>
                <textarea
                  id="description"
                  name="description"
                  value={form.description}
                  onChange={handleChange}
                  placeholder="Category description..."
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={submitting}
                />
              </div>

              <div className="flex gap-4">
                <button
                  type="submit"
                  disabled={submitting}
                  className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submitting ? "Saving..." : isEdit ? "Update Category" : "Create Category"}
                </button>
                <button
                  type="button"
                  onClick={() => router.push("/categories")}
                  disabled={submitting}
                  className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <ToastContainer />
    </div>
  );
}
