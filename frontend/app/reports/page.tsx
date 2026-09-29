"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import MetricCard from "@/components/MetricCard";
import { apiClient } from "@/lib/api";
import { DashboardMetrics, ApiResponse, User } from "@/lib/types";

export default function ReportsPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));

    const fetch = async () => {
      try {
        const res = await apiClient.get<ApiResponse<DashboardMetrics>>(
          "/reports/dashboard"
        );
        if (res.data) setMetrics(res.data);
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    fetch();
  }, []);

  if (!user) return null;
  if (loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar user={user} />
      <main className="flex-1 p-8">
        <Header title="Reports" user={user} />
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
          <MetricCard label="Total Products" value={metrics?.total_products || 0} />
          <MetricCard label="Total Customers" value={metrics?.total_customers || 0} />
          <MetricCard label="Total Sales" value={metrics?.total_sales || 0} variant="success" />
          <MetricCard label="Total Revenue" value={`$${metrics?.total_revenue || 0}`} />
          <MetricCard label="Low Stock" value={metrics?.low_stock_count || 0} variant="danger" />
        </div>
      </main>
    </div>
  );
}
