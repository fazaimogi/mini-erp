"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import MetricCard from "@/components/MetricCard";
import { apiClient } from "@/lib/api";
import { DashboardMetrics, ApiResponse, User, Product, Sale } from "@/lib/types";

const ICON = {
  fill: "none",
  stroke: "currentColor",
  strokeWidth: 1.8,
  strokeLinecap: "round" as const,
  strokeLinejoin: "round" as const,
};

const Svg = ({ d, className = "h-5 w-5" }: { d: string; className?: string }) => (
  <svg viewBox="0 0 24 24" className={className} {...ICON} aria-hidden>
    <path d={d} />
  </svg>
);

const money = (n: unknown) =>
  "$" + (Number(n) || 0).toLocaleString("en-US", { minimumFractionDigits: 2 });

function statusStyle(status: string) {
  const s = status?.toLowerCase();
  if (s === "completed" || s === "paid")
    return "bg-emerald-50 text-emerald-700 ring-emerald-600/20";
  if (s === "pending")
    return "bg-amber-50 text-amber-700 ring-amber-600/20";
  if (s === "cancelled" || s === "void")
    return "bg-rose-50 text-rose-700 ring-rose-600/20";
  return "bg-slate-100 text-slate-600 ring-slate-500/20";
}

const Th = ({ children, className = "" }: { children: React.ReactNode; className?: string }) => (
  <th
    scope="col"
    className={`px-4 py-3 text-left text-[11px] font-semibold uppercase tracking-wider text-slate-500 ${className}`}
  >
    {children}
  </th>
);

const Skeleton = () => (
  <div className="animate-pulse space-y-6" aria-busy="true" aria-label="Loading dashboard">
    <div className="grid grid-cols-1 gap-6 md:grid-cols-2 xl:grid-cols-5">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="h-28 rounded-xl border border-slate-200 bg-white p-5">
          <div className="h-3 w-24 rounded bg-slate-200" />
          <div className="mt-4 h-8 w-16 rounded bg-slate-200" />
        </div>
      ))}
    </div>
    {[0, 1].map(i => (
      <div key={i} className="rounded-xl border border-slate-200 bg-white p-6">
        <div className="h-4 w-40 rounded bg-slate-200" />
        <div className="mt-5 space-y-3">
          {Array.from({ length: 4 }).map((_, j) => (
            <div key={j} className="h-9 rounded bg-slate-100" />
          ))}
        </div>
      </div>
    ))}
  </div>
);

const EmptyState = ({ text, error }: { text: string; error?: boolean }) => (
  <div className="flex flex-col items-center gap-2 px-6 py-8 text-center">
    <span
      className={`grid h-10 w-10 place-items-center rounded-full ${
        error ? "bg-amber-50 text-amber-600" : "bg-emerald-50 text-emerald-600"
      }`}
    >
      <Svg
        d={
          error
            ? "M12 9v4m0 4h.01M10.3 3.9 2.4 18a2 2 0 0 0 1.7 3h15.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z"
            : "M20 6 9 17l-5-5"
        }
      />
    </span>
    <p className="text-sm text-slate-500">{text}</p>
  </div>
);

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [lowStockProducts, setLowStockProducts] = useState<Product[]>([]);
  const [recentSales, setRecentSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetchDashboard = async () => {
      setLoading(true);
      setError("");
      try {
        const [metricsRes, productsRes, salesRes] = await Promise.all([
          apiClient.get<ApiResponse<DashboardMetrics>>("/reports/dashboard"),
          apiClient.get<ApiResponse<Product[]>>("/products"),
          apiClient.get<ApiResponse<Sale[]>>("/sales"),
        ]);

        if (metricsRes.data) setMetrics(metricsRes.data);

        if (Array.isArray(productsRes.data)) {
          const lowStock = productsRes.data.filter(
            p => p.stock <= p.minimum_stock
          );
          setLowStockProducts(lowStock.slice(0, 5));
        }

        if (Array.isArray(salesRes.data)) {
          setRecentSales(salesRes.data.slice(0, 5));
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load dashboard");
      } finally {
        setLoading(false);
      }
    };

    fetchDashboard();
    // refetch when reloadKey bumps (retry); mount reads stored session once
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [reloadKey]);

  if (!user) return null;

  // don't render 0 as if it were real data on failure (ponytail: minimal — swap for retry UI when needed)
  const val = (n: unknown): string | number => (error ? "—" : (n as string | number));

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar user={user} />

      <main className="flex-1 p-6 lg:p-8">
        <Header title="Dashboard" user={user} />

        <div className="animate-fade-up">
          {loading ? (
            <Skeleton />
          ) : (
            <>
              <section className="mb-8 grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-5">
                <MetricCard
                  label="Total Products"
                  value={val(metrics?.total_products ?? 0)}
                  icon={<Svg d="M3 7l9-4 9 4-9 4-9-4Zm0 0v10l9 4 9-4V7" />}
                  hint="Active catalog items"
                />
                <MetricCard
                  label="Total Customers"
                  value={val(metrics?.total_customers ?? 0)}
                  variant="teal"
                  icon={<Svg d="M16 19v-1a4 4 0 0 0-4-4H7a4 4 0 0 0-4 4v1M9.5 10a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm11 9v-1a4 4 0 0 0-3-3.9M16.5 4.2a3 3 0 0 1 0 5.6" />}
                  hint="Registered buyers"
                />
                <MetricCard
                  label="Total Sales"
                  value={val(metrics?.total_sales ?? 0)}
                  variant="emerald"
                  icon={<Svg d="M6 2h9l4 4v16H6V2Zm3 7h6M9 13h6" />}
                  hint="Orders recorded"
                />
                <MetricCard
                  label="Total Revenue"
                  value={error ? "—" : money(metrics?.total_revenue)}
                  variant="violet"
                  icon={<Svg d="M12 2v20M17 6H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />}
                  hint="All recorded sales"
                />
                <MetricCard
                  label="Low Stock"
                  value={val(metrics?.low_stock_count ?? 0)}
                  variant="rose"
                  icon={<Svg d="M12 9v4m0 4h.01M10.3 3.9 2.4 18a2 2 0 0 0 1.7 3h15.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" />}
                  hint="Needs restocking"
                />
              </section>

              {error && (
                <div
                  role="alert"
                  className="mb-6 flex flex-wrap items-center gap-3 rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700"
                >
                  <Svg d="M12 9v4m0 4h.01M10.3 3.9 2.4 18a2 2 0 0 0 1.7 3h15.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" className="h-4 w-4 shrink-0" />
                  <span className="flex-1 font-medium">
                    Couldn&rsquo;t load dashboard data. {error}
                  </span>
                  <button
                    type="button"
                    onClick={() => setReloadKey(k => k + 1)}
                    className="rounded-lg border border-rose-300 bg-white px-3 py-1.5 text-xs font-semibold text-rose-700 transition hover:bg-rose-100"
                  >
                    Retry
                  </button>
                </div>
              )}

              {/* Low Stock Products */}
              <section className="mb-6 overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
                <div className="flex items-center justify-between gap-3 border-b border-slate-100 px-6 py-4">
                  <h2 className="flex items-center gap-2 text-base font-semibold text-slate-900">
                    <span className="grid h-7 w-7 place-items-center rounded-lg bg-rose-50 text-rose-600">
                      <Svg d="M12 9v4m0 4h.01M10.3 3.9 2.4 18a2 2 0 0 0 1.7 3h15.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" className="h-4 w-4" />
                    </span>
                    Low Stock Products
                  </h2>
                  {lowStockProducts.length > 0 && (
                    <span className="rounded-full bg-rose-50 px-2.5 py-0.5 text-xs font-semibold text-rose-700">
                      {lowStockProducts.length}
                    </span>
                  )}
                </div>

                {lowStockProducts.length === 0 ? (
                  <EmptyState text={error ? "Unable to load stock data" : "All products in stock"} error={!!error} />
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full min-w-[520px] text-sm">
                      <thead className="bg-slate-50">
                        <tr>
                          <Th>Product</Th>
                          <Th>SKU</Th>
                          <Th className="text-right">Stock</Th>
                          <Th className="text-right">Min Stock</Th>
                          <Th>Status</Th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-100">
                        {lowStockProducts.map(product => {
                          const critical = product.stock <= product.minimum_stock / 2;
                          return (
                            <tr
                              key={product.id}
                              className="transition-colors even:bg-slate-50/60 hover:bg-indigo-50/50"
                            >
                              <td className="px-4 py-3 font-medium text-slate-800">
                                {product.name}
                              </td>
                              <td className="px-4 py-3 font-mono text-xs text-slate-500">
                                {product.sku}
                              </td>
                              <td className="px-4 py-3 text-right font-semibold tabular-nums text-rose-600">
                                {product.stock}
                              </td>
                              <td className="px-4 py-3 text-right tabular-nums text-slate-500">
                                {product.minimum_stock}
                              </td>
                              <td className="px-4 py-3">
                                <span
                                  className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-semibold ring-1 ring-inset ${
                                    critical
                                      ? "bg-rose-50 text-rose-700 ring-rose-600/20"
                                      : "bg-amber-50 text-amber-700 ring-amber-600/20"
                                  }`}
                                >
                                  {critical ? "Critical" : "Low"}
                                </span>
                              </td>
                            </tr>
                          );
                        })}
                      </tbody>
                    </table>
                  </div>
                )}
              </section>

              {/* Recent Sales */}
              <section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
                <div className="flex items-center justify-between gap-3 border-b border-slate-100 px-6 py-4">
                  <h2 className="flex items-center gap-2 text-base font-semibold text-slate-900">
                    <span className="grid h-7 w-7 place-items-center rounded-lg bg-indigo-50 text-indigo-600">
                      <Svg d="M6 2h9l4 4v16H6V2Zm3 7h6M9 13h6" className="h-4 w-4" />
                    </span>
                    Recent Sales
                  </h2>
                </div>

                {recentSales.length === 0 ? (
                  <EmptyState text={error ? "Unable to load sales data" : "No sales yet"} error={!!error} />
                ) : (
                  <div className="overflow-x-auto">
                    <table className="w-full min-w-[520px] text-sm">
                      <thead className="bg-slate-50">
                        <tr>
                          <Th>Invoice</Th>
                          <Th>Customer</Th>
                          <Th className="text-right">Amount</Th>
                          <Th>Status</Th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-slate-100">
                        {recentSales.map(sale => (
                          <tr
                            key={sale.id}
                            className="transition-colors even:bg-slate-50/60 hover:bg-indigo-50/50"
                          >
                            <td className="px-4 py-3 font-mono text-xs font-medium text-slate-700">
                              {sale.invoice_number}
                            </td>
                            <td className="px-4 py-3 text-slate-600">
                              {sale.customer?.name ?? "—"}
                            </td>
                            <td className="px-4 py-3 text-right font-semibold tabular-nums text-slate-900">
                              {money(sale.total_amount)}
                            </td>
                            <td className="px-4 py-3">
                              <span
                                className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-semibold capitalize ring-1 ring-inset ${statusStyle(
                                  sale.status
                                )}`}
                              >
                                {sale.status}
                              </span>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                )}
              </section>
            </>
          )}
        </div>
      </main>
    </div>
  );
}