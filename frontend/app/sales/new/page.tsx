"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import { apiClient } from "@/lib/api";
import { Category, Customer, Product, ApiResponse, User, VoucherValidation } from "@/lib/types";

interface CartLine {
  product_id: number;
  quantity: number;
}

const money = (n: number) => `$${n.toFixed(2)}`;

function Icon({ d, className = "h-5 w-5" }: { d: string; className?: string }) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={1.8}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden
    >
      <path d={d} />
    </svg>
  );
}

const ICON = {
  search: "M21 21l-4.35-4.35M17 11a6 6 0 1 1-12 0 6 6 0 0 1 12 0Z",
  cart: "M3 4h2l2.4 11.2a2 2 0 0 0 2 1.6h7.6a2 2 0 0 0 2-1.6L21 7H6M9 20h.01M18 20h.01",
  plus: "M12 5v14M5 12h14",
  minus: "M5 12h14",
  trash: "M4 7h16M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2m-8 0 1 13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1l1-13",
  check: "M20 6 9 17l-5-5",
  alert: "M12 9v4m0 4h.01M10.3 3.9 2.4 18a2 2 0 0 0 1.7 3h15.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z",
  tag: "M20.6 13.4 13.4 20.6a2 2 0 0 1-2.8 0l-8-8V3h9.6l8.4 8.4a1.9 1.9 0 0 1 0 2Z",
};

export default function SalesFormPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");

  const [customerId, setCustomerId] = useState(0);
  const [cart, setCart] = useState<CartLine[]>([]);
  const [search, setSearch] = useState("");
  const [activeCategory, setActiveCategory] = useState<number | "all">("all");
  const [notice, setNotice] = useState("");

  const [voucherCode, setVoucherCode] = useState("");
  const [applied, setApplied] = useState<VoucherValidation | null>(null);
  const [applying, setApplying] = useState(false);
  const [voucherError, setVoucherError] = useState("");

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    const stored = localStorage.getItem("user");
    if (!stored) {
      router.push("/login");
      return;
    }
    setUser(JSON.parse(stored));
    setLoading(true);
    setLoadError("");

    (async () => {
      try {
        const [custRes, prodRes, catRes] = await Promise.all([
          apiClient.get<ApiResponse<Customer[]>>("/customers"),
          apiClient.get<ApiResponse<Product[]>>("/products"),
          apiClient.get<ApiResponse<Category[]>>("/categories"),
        ]);
        if (Array.isArray(custRes.data)) setCustomers(custRes.data);
        if (Array.isArray(prodRes.data)) setProducts(prodRes.data);
        if (Array.isArray(catRes.data)) setCategories(catRes.data);
      } catch (err) {
        setLoadError(err instanceof Error ? err.message : "Failed to load catalog");
      } finally {
        setLoading(false);
      }
    })();
    // router identity is unstable under mocks; re-fetch is driven by reloadKey instead
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [reloadKey]);

  const productById = useMemo(() => new Map(products.map(p => [p.id, p])), [products]);

  const visibleProducts = useMemo(() => {
    const q = search.trim().toLowerCase();
    return products.filter(
      p =>
        (activeCategory === "all" || p.category_id === activeCategory) &&
        (!q || p.name.toLowerCase().includes(q) || p.sku.toLowerCase().includes(q))
    );
  }, [products, search, activeCategory]);

  const usedCategories = useMemo(
    () => categories.filter(c => products.some(p => p.category_id === c.id)),
    [categories, products]
  );

  const lines = cart.map(line => {
    const product = productById.get(line.product_id);
    return { ...line, product, lineTotal: line.quantity * Number(product?.price ?? 0) };
  });

  const subtotal = lines.reduce((sum, l) => sum + l.lineTotal, 0);
  // ponytail: discount previewed from discount_type/value so it tracks the live cart;
  // the backend recomputes the authoritative amount on submit.
  const discount = applied
    ? Math.min(
        subtotal,
        applied.discount_type === "percentage"
          ? (subtotal * Number(applied.discount_value)) / 100
          : Number(applied.discount_value)
      )
    : 0;
  const total = Math.max(0, subtotal - discount);
  const itemCount = cart.reduce((sum, l) => sum + l.quantity, 0);
  const canSubmit = !submitting && cart.length > 0 && customerId !== 0;

  const cartQty = (productId: number) => cart.find(l => l.product_id === productId)?.quantity ?? 0;

  const addToCart = (product: Product) => {
    if (product.stock <= 0) return;
    if (cartQty(product.id) >= product.stock) {
      setNotice(`${product.name}: only ${product.stock} in stock`);
      return;
    }
    setNotice("");
    setCart(prev =>
      prev.some(l => l.product_id === product.id)
        ? prev.map(l => (l.product_id === product.id ? { ...l, quantity: l.quantity + 1 } : l))
        : [...prev, { product_id: product.id, quantity: 1 }]
    );
  };

  const setQty = (productId: number, qty: number) => {
    const stock = productById.get(productId)?.stock ?? 1;
    const next = Math.max(1, Math.min(qty, Math.max(1, stock)));
    setNotice("");
    setCart(prev => prev.map(l => (l.product_id === productId ? { ...l, quantity: next } : l)));
  };

  const removeLine = (productId: number) => {
    setNotice("");
    setCart(prev => prev.filter(l => l.product_id !== productId));
  };

  const clearCart = () => {
    setCart([]);
    setNotice("");
  };

  const handleApplyVoucher = async () => {
    if (!voucherCode.trim()) return;
    setApplying(true);
    setVoucherError("");
    setApplied(null);
    try {
      const res = await apiClient.post<ApiResponse<VoucherValidation>>("/vouchers/validate", {
        code: voucherCode.trim(),
        total_amount: subtotal,
      });
      if (res.data) setApplied(res.data);
    } catch (err) {
      setVoucherError(err instanceof Error ? err.message : "Invalid voucher");
    } finally {
      setApplying(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (cart.length === 0) {
      setError("Add at least one item");
      return;
    }
    if (customerId === 0) {
      setError("Select a customer");
      return;
    }

    setSubmitting(true);
    setError("");
    try {
      const payload: Record<string, unknown> = {
        customer_id: customerId,
        items: cart.map(item => ({ product_id: item.product_id, quantity: item.quantity })),
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

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar user={user} />

      <main className="min-w-0 flex-1 p-6 lg:p-8">
        <Header
          title="New Sale"
          user={user}
          action={
            <span className="inline-flex items-center gap-1.5 rounded-full bg-indigo-50 px-3 py-1 text-xs font-semibold text-indigo-700">
              <Icon d={ICON.cart} className="h-3.5 w-3.5" />
              {itemCount} in cart
            </span>
          }
        />

        {loadError && (
          <div
            role="alert"
            className="mb-6 flex flex-wrap items-center gap-3 rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700"
          >
            <Icon d={ICON.alert} className="h-4 w-4 shrink-0" />
            <span className="flex-1 font-medium">Couldn&rsquo;t load the catalog. {loadError}</span>
            <button
              type="button"
              onClick={() => setReloadKey(k => k + 1)}
              className="rounded-lg border border-rose-300 bg-white px-3 py-1.5 text-xs font-semibold text-rose-700 transition hover:bg-rose-100"
            >
              Retry
            </button>
          </div>
        )}

        <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_380px]">
          {/* Catalog */}
          <section className="min-w-0" aria-label="Product catalog">
            <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
              <label htmlFor="catalog-search" className="sr-only">
                Search products
              </label>
              <div className="relative">
                <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-slate-400">
                  <Icon d={ICON.search} className="h-4 w-4" />
                </span>
                <input
                  id="catalog-search"
                  type="search"
                  value={search}
                  onChange={e => setSearch(e.target.value)}
                  placeholder="Search by name or SKU..."
                  className="w-full rounded-lg border border-slate-300 py-2 pl-9 pr-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100"
                />
              </div>
            </div>

            <div className="mt-4 flex flex-wrap gap-2" role="group" aria-label="Filter by category">
              {[{ id: "all" as const, name: "All" }, ...usedCategories].map(cat => {
                const active = activeCategory === cat.id;
                return (
                  <button
                    key={cat.id}
                    type="button"
                    aria-pressed={active}
                    onClick={() => setActiveCategory(cat.id)}
                    className={`rounded-full border px-4 py-1.5 text-sm font-medium transition ${
                      active
                        ? "border-indigo-500 bg-indigo-600 text-white shadow-sm"
                        : "border-slate-200 bg-white text-slate-600 hover:border-indigo-300 hover:text-indigo-700"
                    }`}
                  >
                    {cat.name}
                  </button>
                );
              })}
            </div>

            {loading ? (
              <div className="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                {Array.from({ length: 6 }).map((_, i) => (
                  <div key={i} className="h-40 animate-pulse rounded-xl border border-slate-200 bg-white" />
                ))}
              </div>
            ) : visibleProducts.length === 0 ? (
              <div className="mt-6 rounded-xl border border-dashed border-slate-300 bg-white p-10 text-center">
                <p className="text-sm font-medium text-slate-600">No products match</p>
                <p className="mt-1 text-xs text-slate-400">
                  {search ? "Try a different search term or category." : "This category is empty."}
                </p>
              </div>
            ) : (
              <div className="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                {visibleProducts.map(product => {
                  const inCart = cartQty(product.id);
                  const out = product.stock <= 0;
                  const low = !out && product.stock <= product.minimum_stock;
                  return (
                    <button
                      key={product.id}
                      type="button"
                      onClick={() => addToCart(product)}
                      disabled={out}
                      aria-label={`Add ${product.name} to cart`}
                      className="group relative flex flex-col rounded-xl border border-slate-200 bg-white p-4 text-left shadow-sm transition hover:-translate-y-0.5 hover:border-indigo-300 hover:shadow-md focus-visible:border-indigo-400 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0 disabled:hover:shadow-sm"
                    >
                      {inCart > 0 && (
                        <span className="absolute right-3 top-3 grid h-6 min-w-6 place-items-center rounded-full bg-indigo-600 px-1.5 text-xs font-bold text-white">
                          {inCart}
                        </span>
                      )}

                      <span className="grid h-9 w-9 place-items-center rounded-lg bg-indigo-50 text-xs font-bold uppercase text-indigo-600">
                        {product.name.slice(0, 2)}
                      </span>

                      <p className="mt-3 truncate pr-6 font-semibold text-slate-900">{product.name}</p>
                      <p className="mt-0.5 font-mono text-xs uppercase tracking-wide text-slate-400">
                        {product.sku}
                      </p>

                      <div className="mt-auto flex w-full items-end justify-between pt-4">
                        <span className="text-lg font-bold tabular-nums text-slate-900">
                          {money(Number(product.price))}
                        </span>
                        <span
                          className={`rounded-full px-2 py-0.5 text-xs font-semibold ${
                            out
                              ? "bg-slate-100 text-slate-500"
                              : low
                                ? "bg-amber-50 text-amber-700"
                                : "bg-emerald-50 text-emerald-700"
                          }`}
                        >
                          {out ? "Out of stock" : `${product.stock} in stock`}
                        </span>
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </section>

          {/* Cart */}
          <aside className="lg:sticky lg:top-6" aria-label="Shopping cart">
            <form
              onSubmit={handleSubmit}
              className="flex flex-col overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm"
            >
              <div className="flex items-center justify-between gap-3 border-b border-slate-200 p-4">
                <h2 className="flex items-center gap-2 font-semibold text-slate-900">
                  <span className="grid h-8 w-8 place-items-center rounded-lg bg-indigo-50 text-indigo-600">
                    <Icon d={ICON.cart} className="h-4 w-4" />
                  </span>
                  Cart
                  {itemCount > 0 && (
                    <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs font-semibold text-slate-600">
                      {itemCount}
                    </span>
                  )}
                </h2>
                {cart.length > 0 && (
                  <button
                    type="button"
                    onClick={clearCart}
                    className="text-xs font-semibold text-slate-500 transition hover:text-rose-600"
                  >
                    Clear
                  </button>
                )}
              </div>

              <div className="border-b border-slate-200 p-4">
                <label htmlFor="cart-customer" className="mb-1.5 block text-xs font-semibold uppercase tracking-wide text-slate-500">
                  Customer *
                </label>
                <select
                  id="cart-customer"
                  value={customerId}
                  onChange={e => setCustomerId(parseInt(e.target.value))}
                  className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100"
                >
                  <option value={0}>Select customer</option>
                  {customers.map(c => (
                    <option key={c.id} value={c.id}>
                      {c.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="max-h-[38vh] overflow-y-auto">
                {cart.length === 0 ? (
                  <div className="flex flex-col items-center gap-2 px-4 py-10 text-center">
                    <span className="grid h-11 w-11 place-items-center rounded-full bg-slate-100 text-slate-400">
                      <Icon d={ICON.cart} className="h-5 w-5" />
                    </span>
                    <p className="text-sm font-medium text-slate-600">Cart is empty</p>
                    <p className="text-xs text-slate-400">Click a product to add it</p>
                  </div>
                ) : (
                  <ul className="divide-y divide-slate-100">
                    {lines.map(line => {
                      const name = line.product?.name ?? "Unknown product";
                      const stock = line.product?.stock ?? 1;
                      return (
                        <li key={line.product_id} className="flex items-center gap-2 p-3">
                          <div className="min-w-0 flex-1">
                            <p className="truncate text-sm font-medium text-slate-900">{name}</p>
                            <p className="text-xs tabular-nums text-slate-500">
                              {money(Number(line.product?.price ?? 0))}
                            </p>
                          </div>

                          <div className="flex items-center gap-0.5 rounded-lg border border-slate-200 p-0.5">
                            <button
                              type="button"
                              aria-label={`Decrease quantity of ${name}`}
                              onClick={() => setQty(line.product_id, line.quantity - 1)}
                              disabled={line.quantity <= 1}
                              className="grid h-7 w-7 place-items-center rounded-md text-slate-600 transition hover:bg-slate-100 disabled:opacity-40 disabled:hover:bg-transparent"
                            >
                              <Icon d={ICON.minus} className="h-3.5 w-3.5" />
                            </button>
                            <span className="w-7 text-center text-sm font-semibold tabular-nums text-slate-900">
                              {line.quantity}
                            </span>
                            <button
                              type="button"
                              aria-label={`Increase quantity of ${name}`}
                              onClick={() => setQty(line.product_id, line.quantity + 1)}
                              disabled={line.quantity >= stock}
                              className="grid h-7 w-7 place-items-center rounded-md text-slate-600 transition hover:bg-slate-100 disabled:opacity-40 disabled:hover:bg-transparent"
                            >
                              <Icon d={ICON.plus} className="h-3.5 w-3.5" />
                            </button>
                          </div>

                          <span className="w-16 text-right text-sm font-semibold tabular-nums text-slate-900">
                            {money(line.lineTotal)}
                          </span>

                          <button
                            type="button"
                            aria-label={`Remove ${name} from cart`}
                            onClick={() => removeLine(line.product_id)}
                            className="grid h-7 w-7 place-items-center rounded-md text-slate-400 transition hover:bg-rose-50 hover:text-rose-600"
                          >
                            <Icon d={ICON.trash} className="h-4 w-4" />
                          </button>
                        </li>
                      );
                    })}
                  </ul>
                )}
              </div>

              <div className="space-y-3 border-t border-slate-200 p-4">
                <label htmlFor="cart-voucher" className="block text-xs font-semibold uppercase tracking-wide text-slate-500">
                  Promo Code / Voucher
                </label>
                <div className="flex gap-2">
                  <input
                    id="cart-voucher"
                    type="text"
                    value={voucherCode}
                    onChange={e => {
                      setVoucherCode(e.target.value.toUpperCase());
                      setVoucherError("");
                    }}
                    placeholder="e.g., PROMO2026"
                    className="min-w-0 flex-1 rounded-lg border border-slate-300 px-3 py-2 font-mono text-sm uppercase text-slate-900 outline-none transition placeholder:font-sans placeholder:text-slate-400 focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100"
                  />
                  {applied ? (
                    <button
                      type="button"
                      onClick={() => {
                        setApplied(null);
                        setVoucherCode("");
                      }}
                      className="rounded-lg border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-50"
                    >
                      Remove
                    </button>
                  ) : (
                    <button
                      type="button"
                      onClick={handleApplyVoucher}
                      disabled={applying || !voucherCode.trim()}
                      className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-300"
                    >
                      {applying ? "Checking..." : "Apply"}
                    </button>
                  )}
                </div>

                {voucherError && (
                  <p className="flex items-center gap-1.5 text-xs font-medium text-rose-600">
                    <Icon d={ICON.alert} className="h-3.5 w-3.5 shrink-0" />
                    {voucherError}
                  </p>
                )}
                {applied && (
                  <p className="flex items-center gap-1.5 text-xs font-medium text-emerald-700">
                    <Icon d={ICON.check} className="h-3.5 w-3.5 shrink-0" />
                    Voucher <span className="font-mono font-semibold">{applied.code}</span> applied
                  </p>
                )}

                <dl className="space-y-1.5 border-t border-dashed border-slate-200 pt-3 text-sm">
                  <div className="flex items-center justify-between">
                    <dt className="text-slate-500">Subtotal</dt>
                    <dd className="tabular-nums text-slate-700">{money(subtotal)}</dd>
                  </div>
                  {applied && (
                    <div className="flex items-center justify-between">
                      <dt className="text-slate-500">Discount</dt>
                      <dd className="tabular-nums font-medium text-emerald-600">-{money(discount)}</dd>
                    </div>
                  )}
                  <div className="flex items-end justify-between border-t border-slate-200 pt-2">
                    <dt className="font-semibold text-slate-900">Total</dt>
                    <dd className="text-2xl font-bold tabular-nums text-indigo-600">{money(total)}</dd>
                  </div>
                </dl>
              </div>

              <div className="space-y-2 border-t border-slate-200 p-4">
                {notice && (
                  <p className="flex items-center gap-1.5 text-xs font-medium text-amber-700">
                    <Icon d={ICON.alert} className="h-3.5 w-3.5 shrink-0" />
                    {notice}
                  </p>
                )}
                {error && (
                  <p className="flex items-center gap-1.5 text-xs font-medium text-rose-600" role="alert">
                    <Icon d={ICON.alert} className="h-3.5 w-3.5 shrink-0" />
                    {error}
                  </p>
                )}

                <button
                  type="submit"
                  disabled={!canSubmit}
                  className="flex w-full items-center justify-center gap-2 rounded-lg bg-indigo-600 py-2.5 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-700 disabled:cursor-not-allowed disabled:bg-slate-300"
                >
                  {submitting ? (
                    "Creating..."
                  ) : (
                    <>
                      <Icon d={ICON.check} className="h-4 w-4" />
                      Create Sale &middot; {money(total)}
                    </>
                  )}
                </button>
                <button
                  type="button"
                  onClick={() => router.push("/sales")}
                  className="w-full rounded-lg border border-slate-200 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50"
                >
                  Cancel
                </button>
                {cart.length === 0 && (
                  <p className="text-center text-xs text-slate-400">Add a product to start a sale</p>
                )}
                {cart.length > 0 && customerId === 0 && (
                  <p className="text-center text-xs text-slate-400">Select a customer to continue</p>
                )}
              </div>
            </form>
          </aside>
        </div>
      </main>
    </div>
  );
}