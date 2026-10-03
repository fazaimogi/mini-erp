export interface User {
  id: number;
  name: string;
  email: string;
  role_id: number;
  role: string | { id: number; name: string };
  status?: string;
  created_at?: string;
  updated_at?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  token: string;
}

export interface AuthResponse {
  data: User | LoginResponse;
  message: string;
}

export interface Product {
  id: number;
  name: string;
  sku: string;
  category_id: number;
  description: string;
  price: number;
  stock: number;
  minimum_stock: number;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Customer {
  id: number;
  name: string;
  email: string;
  phone: string;
  address: string;
  created_at: string;
  updated_at: string;
}

export interface Sale {
  id: number;
  invoice_number: string;
  customer_id: number;
  user_id: number;
  total_amount: number;
  discount_amount: number;
  voucher_id: number | null;
  status: string;
  customer?: { id: number; name: string };
  created_at: string;
  updated_at: string;
}

export type DiscountType = "percentage" | "fixed";

export interface Voucher {
  id: number;
  code: string;
  discount_type: DiscountType;
  discount_value: number;
  min_purchase: number;
  usage_limit: number;
  used_count: number;
  expires_at: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface VoucherValidation {
  code: string;
  discount_type: DiscountType;
  discount_value: number;
  discount_amount: number;
  final_amount: number;
}

export interface DashboardMetrics {
  total_products: number;
  total_customers: number;
  total_sales: number;
  total_revenue: number;
  low_stock_count: number;
}

export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: {
    code: string;
    message: string;
  };
  meta?: {
    page: number;
    limit: number;
    total: number;
  };
}
