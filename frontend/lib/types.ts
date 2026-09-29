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
  status: string;
  created_at: string;
  updated_at: string;
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
