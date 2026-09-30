"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import DeleteModal from "@/components/DeleteModal";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Category, ApiResponse, User } from "@/lib/types";

export default function CategoriesPage() {
  const router = useRouter();
  const { success, error } = useToast();
  const [user, setUser] = useState<User | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
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

    const fetchCategories = async () => {
      try {
        const res = await apiClient.get<ApiResponse<Category[]>>("/categories");
        if (res.data && Array.isArray(res.data)) {
          setCategories(res.data);
        }
      } catch (err) {
        error(err instanceof Error ? err.message : "Failed to load categories");
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, []);

  const handleDeleteClick = (id: number, name: string) => {
    setDeleteModal({ open: true, id, name });
  };

  const handleDeleteConfirm = async () => {
    try {
      setDeleting(true);
      await apiClient.delete(`/categories/${deleteModal.id}`);
      setCategories(categories.filter(c => c.id !== deleteModal.id));
      success("Category deleted successfully");
      setDeleteModal({ open: false, id: 0, name: "" });
    } catch (err) {
      error("Failed to delete category");
    } finally {
      setDeleting(false);
    }
  };

  const filtered = categories.filter(c =>
    c.name.toLowerCase().includes(search.toLowerCase()) ||
    c.description.toLowerCase().includes(search.toLowerCase())
  );

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  const isAdmin = (typeof user.role === "string" ? user.role : user.role?.name) === "admin";

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <div className="flex-1">
        <Header title="Categories" user={user} />
        <div className="p-8">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Categories</h1>
            {isAdmin && (
              <button
                onClick={() => router.push("/categories/new")}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                + New Category
              </button>
            )}
          </div>

          <div className="mb-6">
            <input
              type="text"
              placeholder="Search categories..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {filtered.length === 0 ? (
            <div className="text-center text-gray-500 py-12">
              <p>No categories found</p>
            </div>
          ) : (
            <div className="bg-white rounded-lg shadow overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-100 border-b">
                  <tr>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Name</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Description</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Created</th>
                    {isAdmin && (
                      <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Actions</th>
                    )}
                  </tr>
                </thead>
                <tbody>
                  {filtered.map((category) => (
                    <tr key={category.id} className="border-b hover:bg-gray-50 transition">
                      <td className="px-6 py-4 text-sm text-gray-900">
                        <button
                          onClick={() => router.push(`/categories/${category.id}`)}
                          className="text-blue-600 hover:underline"
                        >
                          {category.name}
                        </button>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600">{category.description}</td>
                      <td className="px-6 py-4 text-sm text-gray-600">
                        {new Date(category.created_at).toLocaleDateString()}
                      </td>
                      {isAdmin && (
                        <td className="px-6 py-4 text-sm">
                          <button
                            onClick={() => router.push(`/categories/${category.id}`)}
                            className="text-blue-600 hover:text-blue-900 mr-4"
                          >
                            Edit
                          </button>
                          <button
                            onClick={() => handleDeleteClick(category.id, category.name)}
                            className="text-red-600 hover:text-red-900"
                          >
                            Delete
                          </button>
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>

      <DeleteModal
        isOpen={deleteModal.open}
        title="Delete Category"
        message={`Are you sure you want to delete "${deleteModal.name}"?`}
        onCancel={() => setDeleteModal({ open: false, id: 0, name: "" })}
        onConfirm={handleDeleteConfirm}
      />

      <ToastContainer />
    </div>
  );
}
