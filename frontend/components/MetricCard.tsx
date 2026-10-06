import type { ReactNode } from "react";

type Variant = "indigo" | "teal" | "emerald" | "violet" | "rose";

interface MetricCardProps {
  label: string;
  value: string | number;
  variant?: Variant;
  icon?: ReactNode;
  hint?: string;
}

const styles: Record<Variant, { chip: string; bar: string; ring: string }> = {
  indigo: {
    chip: "bg-indigo-50 text-indigo-600",
    bar: "bg-indigo-500",
    ring: "hover:border-indigo-200 hover:shadow-indigo-100",
  },
  teal: {
    chip: "bg-teal-50 text-teal-600",
    bar: "bg-teal-500",
    ring: "hover:border-teal-200 hover:shadow-teal-100",
  },
  emerald: {
    chip: "bg-emerald-50 text-emerald-600",
    bar: "bg-emerald-500",
    ring: "hover:border-emerald-200 hover:shadow-emerald-100",
  },
  violet: {
    chip: "bg-violet-50 text-violet-600",
    bar: "bg-violet-500",
    ring: "hover:border-violet-200 hover:shadow-violet-100",
  },
  rose: {
    chip: "bg-rose-50 text-rose-600",
    bar: "bg-rose-500",
    ring: "hover:border-rose-200 hover:shadow-rose-100",
  },
};

export default function MetricCard({
  label,
  value,
  variant = "indigo",
  icon,
  hint,
}: MetricCardProps) {
  const s = styles[variant];

  return (
    <div
      className={`group relative flex flex-col overflow-hidden rounded-xl border border-slate-200 bg-white p-5 pt-6 shadow-sm transition duration-200 hover:-translate-y-1 hover:shadow-lg ${s.ring}`}
    >
      {/* top accent bar: the at-a-glance differentiator between cards */}
      <span className={`absolute inset-x-0 top-0 h-1 ${s.bar}`} aria-hidden />

      <div className="relative flex items-start justify-between gap-3">
        {/* fixed two-line slot: wrapping labels no longer shift the row */}
        <p className="min-h-[2rem] text-xs font-semibold uppercase leading-4 tracking-wider text-slate-500">
          {label}
        </p>
        {icon ? (
          <span
            className={`grid h-9 w-9 shrink-0 place-items-center rounded-lg transition-transform duration-200 group-hover:scale-110 ${s.chip}`}
          >
            {icon}
          </span>
        ) : null}
      </div>

      <p className="relative mt-2 text-3xl font-bold tabular-nums tracking-tight text-slate-900">
        {value}
      </p>

      {hint ? (
        <p className="relative mt-auto pt-1 text-xs text-slate-600">{hint}</p>
      ) : null}
    </div>
  );
}