"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import DeleteModal from "@/components/DeleteModal";
import ToastContainer, { useToast } from "@/components/Toast";
import { apiClient } from "@/lib/api";
import { Voucher, ApiResponse, User } from "@/lib/types";

function formatDiscount(v: Voucher) {
  return v.discount_type === "percentage"
    ? `${v.discount_value}%`
    : `$${Number(v.discount_value).toFixed(2)}`;
}

function isExpired(v: Voucher) {
  return v.expires_at ? new Date(v.expires_at).getTime() < Date.now() : false;
}

export default function VouchersPage() {
  const router = useRouter();
  const { success, error } = useToast();
  const [user, setUser] = useState<User | null>(null);
  const [vouchers, setVouchers] = useState<Voucher[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [deleteModal, setDeleteModal] = useState({ open: false, id: 0, code: "" });
  const [togglingId, setTogglingId] = useState<number | null>(null);

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetchVouchers = async () => {
      try {
        const res = await apiClient.get<ApiResponse<Voucher[]>>("/vouchers?limit=100");
        if (res.data && Array.isArray(res.data)) setVouchers(res.data);
      } catch (err) {
        error(err instanceof Error ? err.message : "Failed to load vouchers");
      } finally {
        setLoading(false);
      }
    };

    fetchVouchers();
  }, []);

  const handleToggleActive = async (v: Voucher) => {
    setTogglingId(v.id);
    try {
      const res = await apiClient.put<ApiResponse<Voucher>>(`/vouchers/${v.id}`, {
        code: v.code,
        discount_type: v.discount_type,
        discount_value: Number(v.discount_value),
        min_purchase: Number(v.min_purchase),
        usage_limit: v.usage_limit,
        expires_at: v.expires_at,
        is_active: !v.is_active,
      });
      const updated = res.data;
      setVouchers(prev =>
        prev.map(item =>
          item.id === v.id ? (updated ?? { ...item, is_active: !item.is_active }) : item
        )
      );
      success(v.is_active ? "Voucher deactivated" : "Voucher activated");
    } catch (err) {
      error(err instanceof Error ? err.message : "Failed to update voucher");
    } finally {
      setTogglingId(null);
    }
  };

  const handleDeleteConfirm = async () => {
    try {
      await apiClient.delete(`/vouchers/${deleteModal.id}`);
      setVouchers(prev => prev.filter(v => v.id !== deleteModal.id));
      success("Voucher deleted successfully");
    } catch (err) {
      error(err instanceof Error ? err.message : "Failed to delete voucher");
    } finally {
      setDeleteModal({ open: false, id: 0, code: "" });
    }
  };

  const filtered = vouchers.filter(v =>
    v.code.toLowerCase().includes(search.toLowerCase())
  );

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  const isAdmin = (typeof user.role === "string" ? user.role : user.role?.name) === "admin";

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <div className="flex-1">
        <Header title="Vouchers" user={user} />
        <div className="p-8">
          <div className="flex justify-between items-center mb-6">
            <h1 className="text-3xl font-bold text-gray-900">Vouchers</h1>
            {isAdmin && (
              <button
                onClick={() => router.push("/vouchers/new")}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                + New Voucher
              </button>
            )}
          </div>

          <div className="mb-6">
            <input
              type="text"
              placeholder="Search by code..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          {filtered.length === 0 ? (
            <div className="text-center text-gray-500 py-12">No vouchers found</div>
          ) : (
            <div className="bg-white rounded-lg shadow overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-100 border-b">
                  <tr>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Code</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Discount</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Min Purchase</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Usage</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Expires</th>
                    <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Status</th>
                    {isAdmin && (
                      <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Actions</th>
                    )}
                  </tr>
                </thead>
                <tbody>
                  {filtered.map((v) => {
                    const expired = isExpired(v);
                    const exhausted = v.usage_limit > 0 && v.used_count >= v.usage_limit;
                    return (
                      <tr key={v.id} className="border-b hover:bg-gray-50 transition">
                        <td className="px-6 py-4 text-sm font-mono font-semibold text-gray-900">{v.code}</td>
                        <td className="px-6 py-4 text-sm text-gray-600">{formatDiscount(v)}</td>
                        <td className="px-6 py-4 text-sm text-gray-600">
                          ${Number(v.min_purchase).toFixed(2)}
                        </td>
                        <td className="px-6 py-4 text-sm text-gray-600">
                          {v.used_count} / {v.usage_limit === 0 ? "∞" : v.usage_limit}
                        </td>
                        <td className="px-6 py-4 text-sm text-gray-600">
                          {v.expires_at ? new Date(v.expires_at).toLocaleDateString() : "—"}
                        </td>
                        <td className="px-6 py-4">
                          <span
                            className={`text-xs font-semibold px-2 py-1 rounded ${
                              !v.is_active
                                ? "bg-gray-200 text-gray-600"
                                : expired
                                ? "bg-red-100 text-red-700"
                                : exhausted
                                ? "bg-yellow-100 text-yellow-700"
                                : "bg-green-100 text-green-700"
                            }`}
                          >
                            {!v.is_active ? "inactive" : expired ? "expired" : exhausted ? "exhausted" : "active"}
                          </span>
                        </td>
                        {isAdmin && (
                          <td className="px-6 py-4 text-sm space-x-4 whitespace-nowrap">
                            <button
                              onClick={() => router.push(`/vouchers/${v.id}`)}
                              className="text-blue-600 hover:text-blue-900"
                            >
                              Edit
                            </button>
                            <button
                              onClick={() => handleToggleActive(v)}
                              disabled={togglingId === v.id}
                              className="text-amber-600 hover:text-amber-900 disabled:opacity-50"
                            >
                              {v.is_active ? "Deactivate" : "Activate"}
                            </button>
                            <button
                              onClick={() => setDeleteModal({ open: true, id: v.id, code: v.code })}
                              className="text-red-600 hover:text-red-900"
                            >
                              Delete
                            </button>
                          </td>
                        )}
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>

      <DeleteModal
        isOpen={deleteModal.open}
        title="Delete Voucher"
        message={`Are you sure you want to delete voucher "${deleteModal.code}"?`}
        onCancel={() => setDeleteModal({ open: false, id: 0, code: "" })}
        onConfirm={handleDeleteConfirm}
      />

      <ToastContainer />
    </div>
  );
}