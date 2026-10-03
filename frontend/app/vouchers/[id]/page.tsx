"use client";

import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Voucher, DiscountType, ApiResponse, User } from "@/lib/types";

interface VoucherForm {
  code: string;
  discount_type: DiscountType;
  discount_value: number;
  min_purchase: number;
  usage_limit: number;
  expires_at: string;
  is_active: boolean;
}

const emptyForm: VoucherForm = {
  code: "",
  discount_type: "percentage",
  discount_value: 0,
  min_purchase: 0,
  usage_limit: 0,
  expires_at: "",
  is_active: true,
};

export default function VoucherFormPage() {
  const router = useRouter();
  const params = useParams();
  const { success, error: showError } = useToast();
  const id = params?.id as string;
  const isEdit = !!id && id !== "new";

  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [form, setForm] = useState<VoucherForm>(emptyForm);

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
      showError("Only admins can manage vouchers");
      router.push("/vouchers");
      return;
    }

    if (isEdit) {
      const fetchVoucher = async () => {
        try {
          const res = await apiClient.get<ApiResponse<Voucher>>(`/vouchers/${id}`);
          if (res.data) {
            const v = res.data;
            setForm({
              code: v.code,
              discount_type: v.discount_type,
              discount_value: Number(v.discount_value),
              min_purchase: Number(v.min_purchase),
              usage_limit: v.usage_limit,
              expires_at: v.expires_at ? v.expires_at.slice(0, 10) : "",
              is_active: v.is_active,
            });
          }
        } catch (err) {
          showError(err instanceof Error ? err.message : "Failed to load voucher");
          router.push("/vouchers");
        } finally {
          setLoading(false);
        }
      };
      fetchVoucher();
    } else {
      setLoading(false);
    }
  }, [id, isEdit]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!form.code.trim()) {
      setError("Code is required");
      return;
    }
    if (form.discount_value <= 0) {
      setError("Discount value must be greater than 0");
      return;
    }
    if (form.discount_type === "percentage" && form.discount_value > 100) {
      setError("Percentage discount cannot exceed 100");
      return;
    }

    const payload = {
      code: form.code.trim().toUpperCase(),
      discount_type: form.discount_type,
      discount_value: form.discount_value,
      min_purchase: form.min_purchase,
      usage_limit: form.usage_limit,
      expires_at: form.expires_at ? new Date(form.expires_at).toISOString() : null,
      is_active: form.is_active,
    };

    try {
      setSubmitting(true);
      if (isEdit) {
        await apiClient.put(`/vouchers/${id}`, payload);
        success("Voucher updated successfully");
      } else {
        await apiClient.post("/vouchers", payload);
        success("Voucher created successfully");
      }
      router.push("/vouchers");
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Failed to save voucher";
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
        <Header title={isEdit ? "Edit Voucher" : "New Voucher"} user={user} />
        <div className="p-8">
          <div className="max-w-2xl">
            <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6">
              {error && (
                <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded text-red-700">
                  {error}
                </div>
              )}

              <div className="mb-6">
                <label htmlFor="code" className="block text-sm font-medium text-gray-700 mb-2">
                  Code *
                </label>
                <input
                  type="text"
                  id="code"
                  value={form.code}
                  onChange={(e) => setForm(prev => ({ ...prev, code: e.target.value.toUpperCase() }))}
                  placeholder="e.g., PROMO2026"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono"
                  disabled={submitting}
                />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-6">
                <div>
                  <label htmlFor="discount_type" className="block text-sm font-medium text-gray-700 mb-2">
                    Discount Type *
                  </label>
                  <select
                    id="discount_type"
                    value={form.discount_type}
                    onChange={(e) => setForm(prev => ({ ...prev, discount_type: e.target.value as DiscountType }))}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    disabled={submitting}
                  >
                    <option value="percentage">Percentage (%)</option>
                    <option value="fixed">Fixed ($)</option>
                  </select>
                </div>
                <div>
                  <label htmlFor="discount_value" className="block text-sm font-medium text-gray-700 mb-2">
                    {form.discount_type === "percentage" ? "Discount (%) *" : "Discount ($) *"}
                  </label>
                  <input
                    type="number"
                    id="discount_value"
                    min="0"
                    step="0.01"
                    value={form.discount_value}
                    onChange={(e) => setForm(prev => ({ ...prev, discount_value: parseFloat(e.target.value) || 0 }))}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    disabled={submitting}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-6">
                <div>
                  <label htmlFor="min_purchase" className="block text-sm font-medium text-gray-700 mb-2">
                    Min Purchase ($)
                  </label>
                  <input
                    type="number"
                    id="min_purchase"
                    min="0"
                    step="0.01"
                    value={form.min_purchase}
                    onChange={(e) => setForm(prev => ({ ...prev, min_purchase: parseFloat(e.target.value) || 0 }))}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    disabled={submitting}
                  />
                </div>
                <div>
                  <label htmlFor="usage_limit" className="block text-sm font-medium text-gray-700 mb-2">
                    Usage Limit (0 = unlimited)
                  </label>
                  <input
                    type="number"
                    id="usage_limit"
                    min="0"
                    value={form.usage_limit}
                    onChange={(e) => setForm(prev => ({ ...prev, usage_limit: parseInt(e.target.value) || 0 }))}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    disabled={submitting}
                  />
                </div>
              </div>

              <div className="mb-6">
                <label htmlFor="expires_at" className="block text-sm font-medium text-gray-700 mb-2">
                  Expires At (optional)
                </label>
                <input
                  type="date"
                  id="expires_at"
                  value={form.expires_at}
                  onChange={(e) => setForm(prev => ({ ...prev, expires_at: e.target.value }))}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={submitting}
                />
              </div>

              <div className="mb-6 flex items-center gap-3">
                <input
                  type="checkbox"
                  id="is_active"
                  checked={form.is_active}
                  onChange={(e) => setForm(prev => ({ ...prev, is_active: e.target.checked }))}
                  className="w-4 h-4 rounded"
                  disabled={submitting}
                />
                <label htmlFor="is_active" className="text-sm font-medium text-gray-700">
                  Active
                </label>
              </div>

              <div className="flex gap-4">
                <button
                  type="submit"
                  disabled={submitting}
                  className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submitting ? "Saving..." : isEdit ? "Update Voucher" : "Create Voucher"}
                </button>
                <button
                  type="button"
                  onClick={() => router.push("/vouchers")}
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