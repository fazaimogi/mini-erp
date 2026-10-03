"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import { apiClient } from "@/lib/api";
import { Customer, Product, ApiResponse, User, VoucherValidation } from "@/lib/types";

interface SaleItem {
  product_id: number;
  quantity: number;
}

export default function SalesFormPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [products, setProducts] = useState<Product[]>([]);

  const [form, setForm] = useState({
    customer_id: 0,
    items: [] as SaleItem[],
  });

  const [voucherCode, setVoucherCode] = useState("");
  const [applied, setApplied] = useState<VoucherValidation | null>(null);
  const [applying, setApplying] = useState(false);
  const [voucherError, setVoucherError] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetch = async () => {
      try {
        const [custRes, prodRes] = await Promise.all([
          apiClient.get<ApiResponse<Customer[]>>("/customers"),
          apiClient.get<ApiResponse<Product[]>>("/products"),
        ]);
        if (Array.isArray(custRes.data)) setCustomers(custRes.data);
        if (Array.isArray(prodRes.data)) setProducts(prodRes.data);
      } catch (err) {
        console.error(err);
      }
    };
    fetch();
  }, []);

  const handleAddItem = () => {
    setForm(prev => ({
      ...prev,
      items: [...prev.items, { product_id: 0, quantity: 1 }],
    }));
  };

  const handleRemoveItem = (idx: number) => {
    setForm(prev => ({
      ...prev,
      items: prev.items.filter((_, i) => i !== idx),
    }));
  };

  const handleItemChange = (idx: number, field: string, value: any) => {
    setForm(prev => {
      const items = [...prev.items];
      items[idx] = { ...items[idx], [field]: value };
      return { ...prev, items };
    });
  };

  const handleApplyVoucher = async () => {
    if (!voucherCode.trim()) return;
    setApplying(true);
    setVoucherError("");
    setApplied(null);
    try {
      const res = await apiClient.post<ApiResponse<VoucherValidation>>("/vouchers/validate", {
        code: voucherCode.trim(),
        total_amount: total,
      });
      if (res.data) {
        setApplied(res.data);
      }
    } catch (err) {
      setVoucherError(err instanceof Error ? err.message : "Invalid voucher");
    } finally {
      setApplying(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (form.items.length === 0) {
      setError("Add at least one item");
      return;
    }
    if (form.customer_id === 0) {
      setError("Select a customer");
      return;
    }

    setSubmitting(true);
    setError("");
    try {
      const payload: Record<string, unknown> = {
        customer_id: form.customer_id,
        items: form.items.map(item => ({
          product_id: item.product_id,
          quantity: item.quantity,
        })),
      };
      if (applied) payload.voucher_code = applied.code;
      await apiClient.post("/sales", payload);
      router.push("/sales");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed");
      setSubmitting(false);
    }
  };

  if (!user) return null;

  const total = form.items.reduce((sum, item) => {
    const prod = products.find(p => p.id === item.product_id);
    return sum + (item.quantity * (prod?.price || 0));
  }, 0);

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <main className="flex-1 p-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold">New Sale</h1>
        </div>
        {error && <div className="text-red-600 mb-6 bg-red-50 p-4 rounded">{error}</div>}
        <div className="bg-white rounded-lg shadow-sm p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">Customer *</label>
              <select
                value={form.customer_id}
                onChange={(e) => setForm(prev => ({ ...prev, customer_id: parseInt(e.target.value) }))}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value={0}>Select customer</option>
                {customers.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </div>

            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-semibold">Items</h3>
                <button
                  type="button"
                  onClick={handleAddItem}
                  className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm"
                >
                  + Add Item
                </button>
              </div>

              {form.items.length === 0 ? (
                <div className="text-gray-500 text-center py-4">No items added</div>
              ) : (
                <div className="space-y-3">
                  {form.items.map((item, idx) => (
                    <div key={idx} className="flex gap-3 p-4 bg-gray-50 rounded-lg">
                      <select
                        value={item.product_id}
                        onChange={(e) => handleItemChange(idx, "product_id", parseInt(e.target.value))}
                        className="flex-1 px-3 py-2 border border-gray-300 rounded"
                      >
                        <option value={0}>Product</option>
                        {products.map(p => (
                          <option key={p.id} value={p.id}>{p.name}</option>
                        ))}
                      </select>
                      <input
                        type="number"
                        min="1"
                        value={item.quantity}
                        onChange={(e) => handleItemChange(idx, "quantity", parseInt(e.target.value))}
                        className="w-20 px-3 py-2 border border-gray-300 rounded"
                        placeholder="Qty"
                      />
                      <span className="w-24 px-3 py-2 text-right font-semibold">
                        ${(item.quantity * (products.find(p => p.id === item.product_id)?.price || 0)).toFixed(2)}
                      </span>
                      <button
                        type="button"
                        onClick={() => handleRemoveItem(idx)}
                        className="text-red-600 hover:text-red-700 font-semibold"
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="pt-4 border-t space-y-4">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-2">
                  Promo Code / Voucher
                </label>
                <div className="flex gap-3">
                  <input
                    type="text"
                    value={voucherCode}
                    onChange={(e) => {
                      setVoucherCode(e.target.value.toUpperCase());
                      setVoucherError("");
                    }}
                    placeholder="e.g., PROMO2026"
                    className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono"
                  />
                  <button
                    type="button"
                    onClick={handleApplyVoucher}
                    disabled={applying || !voucherCode.trim()}
                    className="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white rounded-lg font-medium transition"
                  >
                    {applying ? "Checking..." : "Apply"}
                  </button>
                  {applied && (
                    <button
                      type="button"
                      onClick={() => {
                        setApplied(null);
                        setVoucherCode("");
                      }}
                      className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition"
                    >
                      Remove
                    </button>
                  )}
                </div>
                {voucherError && (
                  <p className="mt-2 text-sm text-red-600">{voucherError}</p>
                )}
                {applied && (
                  <p className="mt-2 text-sm text-green-700">
                    Voucher <span className="font-mono font-semibold">{applied.code}</span> applied
                  </p>
                )}
              </div>

              <div className="text-right">
                <p className="text-gray-600">Subtotal: ${total.toFixed(2)}</p>
                {applied && (
                  <p className="text-green-600">
                    Discount: -${Number(applied.discount_amount).toFixed(2)}
                  </p>
                )}
                <p className="text-gray-600 mt-1">Total:</p>
                <p className="text-3xl font-bold text-blue-600">
                  ${(applied ? Number(applied.final_amount) : total).toFixed(2)}
                </p>
              </div>
            </div>

            <div className="flex gap-4 pt-6">
              <button
                type="submit"
                disabled={submitting}
                className="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white py-2 px-4 rounded-lg font-medium transition"
              >
                {submitting ? "Creating..." : "Create Sale"}
              </button>
              <button
                type="button"
                onClick={() => router.push("/sales")}
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
