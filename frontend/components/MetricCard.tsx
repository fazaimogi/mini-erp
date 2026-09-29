interface MetricCardProps {
  label: string;
  value: string | number;
  variant?: "default" | "success" | "warning" | "danger";
}

export default function MetricCard({ label, value, variant = "default" }: MetricCardProps) {
  const variantStyles = {
    default: "text-blue-600",
    success: "text-green-600",
    warning: "text-orange-600",
    danger: "text-red-600",
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md hover:-translate-y-0.5 transition">
      <p className="text-gray-600 text-sm uppercase font-semibold tracking-wide">{label}</p>
      <p className={`text-4xl font-bold mt-3 ${variantStyles[variant]}`}>{value}</p>
    </div>
  );
}
