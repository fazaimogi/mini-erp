"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import MetricCard from "@/components/MetricCard";
import { apiClient } from "@/lib/api";
import { DashboardMetrics, ApiResponse, User, Product, Sale } from "@/lib/types";

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [lowStockProducts, setLowStockProducts] = useState<Product[]>([]);
  const [recentSales, setRecentSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetchDashboard = async () => {
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
  }, []);

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />

      <main className="flex-1 p-8">
        <Header title="Dashboard" user={user} />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
          <MetricCard label="Total Products" value={metrics?.total_products || 0} />
          <MetricCard label="Total Customers" value={metrics?.total_customers || 0} />
          <MetricCard label="Total Sales" value={metrics?.total_sales || 0} variant="success" />
          <MetricCard label="Total Revenue" value={`$${(Number(metrics?.total_revenue) || 0).toFixed(2)}`} />
          {metrics?.low_stock_count ? (
            <MetricCard
              label="Low Stock"
              value={metrics.low_stock_count}
              variant="danger"
            />
          ) : null}
        </div>

        {error && <div className="text-red-600 mb-4 bg-red-50 p-4 rounded">{error}</div>}

        {/* Low Stock Products Table */}
        <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
          <h3 className="text-xl font-bold mb-6 pb-4 border-b">Low Stock Products</h3>
          {lowStockProducts.length === 0 ? (
            <div className="text-gray-500 py-4">All products in stock</div>
          ) : (
            <table className="w-full text-sm">
              <thead className="bg-gray-50">
                <tr>
                  <th className="text-left px-4 py-3 font-semibold text-gray-700">Product</th>
                  <th className="text-left px-4 py-3 font-semibold text-gray-700">Stock</th>
                  <th className="text-left px-4 py-3 font-semibold text-gray-700">Min Stock</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {lowStockProducts.map(product => (
                  <tr key={product.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3">{product.name}</td>
                    <td className="px-4 py-3 text-red-600 font-semibold">{product.stock}</td>
                    <td className="px-4 py-3">{product.minimum_stock}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>

        {/* Recent Sales Table */}
        <div className="bg-white rounded-lg shadow-sm p-6">
          <h3 className="text-xl font-bold mb-6 pb-4 border-b">Recent Sales</h3>
          {recentSales.length === 0 ? (
            <div className="text-gray-500 py-4">No sales yet</div>
          ) : (
            <table className="w-full text-sm">
              <thead className="bg-gray-50">
                <tr>
                  <th className="text-left px-4 py-3 font-semibold text-gray-700">Invoice</th>
                  <th className="text-left px-4 py-3 font-semibold text-gray-700">Amount</th>
                  <th className="text-left px-4 py-3 font-semibold text-gray-700">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {recentSales.map(sale => (
                  <tr key={sale.id} className="hover:bg-gray-50">
                    <td className="px-4 py-3">{sale.invoice_number}</td>
                    <td className="px-4 py-3 font-semibold">${(Number(sale.total_amount) || 0).toFixed(2)}</td>
                    <td className="px-4 py-3">
                      <span className={`text-xs font-semibold px-2 py-1 rounded ${
                        sale.status === "completed"
                          ? "bg-green-100 text-green-700"
                          : "bg-yellow-100 text-yellow-700"
                      }`}>
                        {sale.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </main>
    </div>
  );
}
