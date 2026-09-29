"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { User, ApiResponse } from "@/lib/types";

interface UserData extends User {
  password?: string;
}

export default function UserFormPage() {
  const router = useRouter();
  const params = useParams();
  const id = params?.id as string;
  const isEdit = !!id && id !== "new";
  const { success, error: showError } = useToast();

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [roles, setRoles] = useState<{ id: number; name: string }[]>([]);

  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    role_id: 2,
  });

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    const userData = JSON.parse(stored);
    setUser(userData);

    const roleCheck = typeof userData.role === "string" ? userData.role : userData.role?.name;
    if (roleCheck !== "admin") {
      router.push("/users");
      return;
    }

    const fetchRoles = async () => {
      try {
        const res = await apiClient.get<ApiResponse<{ id: number; name: string }[]>>("/roles");
        if (Array.isArray(res.data)) setRoles(res.data);
      } catch (err) {
        showError(err instanceof Error ? err.message : "Failed to load roles");
      }
    };
    fetchRoles();

    if (isEdit) {
      const fetchUser = async () => {
        try {
          const res = await apiClient.get<ApiResponse<UserData>>(`/users/${id}`);
          if (res.data) {
            setForm({
              name: res.data.name,
              email: res.data.email,
              password: "",
              role_id: typeof res.data.role === "string" ? 2 : res.data.role.id,
            });
          }
        } catch (err) {
          showError(err instanceof Error ? err.message : "Failed to load");
        } finally {
          setLoading(false);
        }
      };
      fetchUser();
    }
  }, [id, isEdit]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setForm(prev => ({
      ...prev,
      [name]: name === "role_id" ? parseInt(value) : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const payload = isEdit
        ? { name: form.name, email: form.email, role_id: form.role_id }
        : form;

      if (isEdit) {
        await apiClient.put(`/users/${id}`, payload);
        success("User updated");
      } else {
        await apiClient.post("/users", form);
        success("User created");
      }
      router.push("/users");
    } catch (err) {
      const errMsg = err instanceof Error ? err.message : "Failed";
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
          <h1 className="text-3xl font-bold">{isEdit ? "Edit User" : "New User"}</h1>
          <p className="text-gray-600 mt-2">
            {isEdit ? "Update user details" : "Create a new user"}
          </p>
        </div>

        {error && <div className="text-red-600 mb-6 bg-red-50 p-4 rounded">{error}</div>}

        <div className="bg-white rounded-lg shadow-sm p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Name *
              </label>
              <input
                type="text"
                name="name"
                value={form.name}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="John Doe"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Email *
              </label>
              <input
                type="email"
                name="email"
                value={form.email}
                onChange={handleChange}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="john@example.com"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Password {!isEdit && "*"}
              </label>
              <input
                type="password"
                name="password"
                value={form.password}
                onChange={handleChange}
                required={!isEdit}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder={isEdit ? "Leave blank to keep current" : "••••••••"}
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Role *
              </label>
              <select
                name="role_id"
                value={form.role_id}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value={0}>Select role</option>
                {roles.map(role => (
                  <option key={role.id} value={role.id}>
                    {role.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex gap-4 pt-6">
              <button
                type="submit"
                disabled={submitting}
                className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-2 px-4 rounded-lg font-medium transition"
              >
                {submitting ? "Saving..." : isEdit ? "Update User" : "Create User"}
              </button>
              <button
                type="button"
                onClick={() => router.push("/users")}
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
