"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import DeleteModal from "@/components/DeleteModal";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { User, ApiResponse } from "@/lib/types";

interface UserData extends User {
  status?: string;
}

export default function UsersPage() {
  const router = useRouter();
  const { success, error } = useToast();
  const [user, setUser] = useState<User | null>(null);
  const [users, setUsers] = useState<UserData[]>([]);
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
    const userData = JSON.parse(stored);
    setUser(userData);

    const isAdmin = (typeof userData.role === "string" ? userData.role : userData.role?.name) === "admin";
    if (!isAdmin) {
      router.push("/dashboard");
      return;
    }

    const fetch = async () => {
      try {
        const res = await apiClient.get<ApiResponse<UserData[]>>("/users");
        if (Array.isArray(res.data)) setUsers(res.data);
      } catch (err) {
        error(err instanceof Error ? err.message : "Failed to load");
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, []);

  const handleDeleteClick = (id: number, name: string) => {
    setDeleteModal({ open: true, id, name });
  };

  const handleDeleteConfirm = async () => {
    setDeleting(true);
    try {
      await apiClient.delete(`/users/${deleteModal.id}`);
      setUsers(users.filter(u => u.id !== deleteModal.id));
      success(`User "${deleteModal.name}" deleted`);
      setDeleteModal({ open: false, id: 0, name: "" });
    } catch (err) {
      error(err instanceof Error ? err.message : "Delete failed");
    } finally {
      setDeleting(false);
    }
  };

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <ToastContainer />
      <main className="flex-1 p-8">
        <div className="flex justify-between items-center mb-8 bg-white rounded-lg shadow-sm p-6">
          <h1 className="text-3xl font-bold">Users</h1>
          <button
            onClick={() => router.push("/users/new")}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded font-medium transition"
          >
            + New User
          </button>
        </div>
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Name</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Email</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Role</th>
                <th className="text-left px-6 py-4 font-semibold text-gray-700 text-sm uppercase">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {users.map((u) => (
                <tr key={u.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{u.name}</td>
                  <td className="px-6 py-4">{u.email}</td>
                  <td className="px-6 py-4">
                    <span className={`text-xs font-semibold px-2 py-1 rounded ${
                      (typeof u.role === "string" ? u.role : u.role?.name) === "admin"
                        ? "bg-orange-100 text-orange-700"
                        : "bg-blue-100 text-blue-700"
                    }`}>
                      {typeof u.role === "string" ? u.role : u.role?.name}
                    </span>
                  </td>
                  <td className="px-6 py-4 space-x-2">
                    <button
                      onClick={() => router.push(`/users/${u.id}`)}
                      className="text-blue-600 hover:underline text-sm"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteClick(u.id, u.name)}
                      className="text-red-600 hover:underline text-sm"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </main>

      <DeleteModal
        isOpen={deleteModal.open}
        title="Delete User"
        message={`Are you sure you want to delete "${deleteModal.name}"? This action cannot be undone.`}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setDeleteModal({ open: false, id: 0, name: "" })}
        isLoading={deleting}
      />
    </div>
  );
}
