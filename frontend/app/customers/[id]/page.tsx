"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Customer, ApiResponse, User } from "@/lib/types";

export default function CustomerFormPage() {
  const router = useRouter();
  const params = useParams();
  const id = params?.id as string;
  const isEdit = !!id && id !== "new";
  const { success, error: showError } = useToast();

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const [form, setForm] = useState({
    name: "",
    email: "",
    phone: "",
    address: "",
  });

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    if (isEdit) {
      const fetchCustomer = async () => {
        try {
          const res = await apiClient.get<ApiResponse<Customer>>(`/customers/${id}`);
          if (res.data) {
            setForm(res.data);
          }
        } catch (err) {
          const errMsg = err instanceof Error ? err.message : "Failed to load";
          setError(errMsg);
          showError(errMsg);
        } finally {
          setLoading(false);
        }
      };
      fetchCustomer();
    }
  }, [id, isEdit]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      if (isEdit) {
        await apiClient.put(`/customers/${id}`, form);
        success("Customer updated");
      } else {
        await apiClient.post("/customers", form);
        success("Customer created");
      }
      router.push("/customers");
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
          <h1 className="text-3xl font-bold">{isEdit ? "Edit Customer" : "New Customer"}</h1>
          <p className="text-gray-600 mt-2">
            {isEdit ? "Update customer details" : "Create a new customer"}
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
                Phone
              </label>
              <input
                type="tel"
                name="phone"
                value={form.phone}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="+62 812 3456 7890"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Address
              </label>
              <textarea
                name="address"
                value={form.address}
                onChange={handleChange}
                rows={4}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Street address..."
              />
            </div>

            <div className="flex gap-4 pt-6">
              <button
                type="submit"
                disabled={submitting}
                className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-2 px-4 rounded-lg font-medium transition"
              >
                {submitting ? "Saving..." : isEdit ? "Update" : "Create"}
              </button>
              <button
                type="button"
                onClick={() => router.push("/customers")}
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
