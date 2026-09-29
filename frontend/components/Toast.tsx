import { useEffect, useState } from "react";

export type ToastType = "success" | "error" | "info";

interface Toast {
  id: string;
  message: string;
  type: ToastType;
}

let toastId = 0;
const listeners: Set<(toasts: Toast[]) => void> = new Set();
let toasts: Toast[] = [];

export const useToast = () => {
  const [toastList, setToastList] = useState<Toast[]>([]);

  useEffect(() => {
    listeners.add(setToastList);
    return () => {
      listeners.delete(setToastList);
    };
  }, []);

  const show = (message: string, type: ToastType = "info", duration = 3000) => {
    const id = `toast-${++toastId}`;
    const newToast = { id, message, type };
    toasts = [...toasts, newToast];
    listeners.forEach(listener => listener(toasts));

    setTimeout(() => {
      toasts = toasts.filter(t => t.id !== id);
      listeners.forEach(listener => listener(toasts));
    }, duration);
  };

  return {
    toasts: toastList,
    success: (msg: string) => show(msg, "success"),
    error: (msg: string) => show(msg, "error"),
    info: (msg: string) => show(msg, "info"),
  };
};

export default function ToastContainer() {
  const { toasts } = useToast();

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2">
      {toasts.map(toast => (
        <div
          key={toast.id}
          className={`px-4 py-3 rounded-lg shadow-lg text-white animate-slide-in ${
            toast.type === "success"
              ? "bg-green-600"
              : toast.type === "error"
              ? "bg-red-600"
              : "bg-blue-600"
          }`}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
}
