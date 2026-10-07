import type { ProductSaleInfo } from "@/lib/types";

interface TopSellingChartProps {
  items: ProductSaleInfo[];
}

// Horizontal bar chart, pure CSS — one row per product, width % of the best seller.
export default function TopSellingChart({ items }: TopSellingChartProps) {
  const max = Math.max(1, ...items.map((i) => i.quantity));

  if (items.length === 0) {
    return (
      <div className="grid place-items-center rounded-xl border border-dashed border-slate-300 bg-white py-12 text-sm text-slate-500">
        No sales yet.
      </div>
    );
  }

  return (
    <div className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
      <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500">
        Top Selling Products
      </h2>
      <ul className="mt-5 space-y-4">
        {items.map((item, idx) => {
          const pct = Math.max(2, Math.round((item.quantity / max) * 100));
          return (
            <li key={item.id} className="group">
              <div className="flex items-baseline justify-between gap-3 text-sm">
                <span className="truncate font-medium text-slate-800">
                  <span className="mr-2 tabular-nums text-slate-400">
                    {idx + 1}.
                  </span>
                  {item.name}
                  <span className="ml-2 font-mono text-xs text-slate-400">
                    {item.sku}
                  </span>
                </span>
                <span className="shrink-0 tabular-nums text-slate-600">
                  {item.quantity} sold
                </span>
              </div>
              <div
                className="mt-1.5 h-2.5 w-full overflow-hidden rounded-full bg-slate-100"
                role="img"
                aria-label={`${item.name}: ${item.quantity} sold`}
              >
                <div
                  className="h-full rounded-full bg-indigo-500 transition-all duration-500 group-hover:bg-indigo-600"
                  style={{ width: `${pct}%` }}
                />
              </div>
              <p className="mt-1 text-xs text-slate-400">
                ${Number(item.revenue).toFixed(2)} revenue
              </p>
            </li>
          );
        })}
      </ul>
    </div>
  );
}