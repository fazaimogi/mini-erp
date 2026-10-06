import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import DashboardPage from '@/app/dashboard/page';
import { apiClient } from '@/lib/api';

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
  usePathname: () => '/dashboard',
}));

jest.mock('@/lib/api', () => ({
  apiClient: { get: jest.fn() },
}));

const get = apiClient.get as jest.Mock;

const metrics = {
  total_products: 142,
  total_customers: 38,
  total_sales: 256,
  total_revenue: 45200.5,
  low_stock_count: 1,
};

const products = [
  { id: 1, name: 'Kopi Arabika', sku: 'KP-01', stock: 2, minimum_stock: 10 },
  { id: 2, name: 'Teh Hijau', sku: 'TH-02', stock: 80, minimum_stock: 10 },
];

const sales = [
  {
    id: 1,
    invoice_number: 'INV-001',
    total_amount: 125000,
    status: 'completed',
    customer: { id: 1, name: 'Siti' },
  },
];

const ok = (data: unknown) => Promise.resolve({ data });

describe('DashboardPage', () => {
  beforeEach(() => {
    get.mockReset();
    localStorage.setItem(
      'user',
      JSON.stringify({ id: 1, name: 'Budi', email: 'budi@example.com', role_id: 1, role: 'admin' })
    );
  });

  it('renders metrics, filters low stock, and lists recent sales', async () => {
    get.mockImplementation((endpoint: string) => {
      if (endpoint === '/reports/dashboard') return ok(metrics);
      if (endpoint === '/products') return ok(products);
      if (endpoint === '/sales') return ok(sales);
      return ok([]);
    });

    render(<DashboardPage />);

    expect(await screen.findByText('142')).toBeInTheDocument();
    expect(screen.getByText('256')).toBeInTheDocument();
    // only the product at/below minimum_stock reaches the table
    expect(await screen.findByText('Kopi Arabika')).toBeInTheDocument();
    expect(screen.queryByText('Teh Hijau')).not.toBeInTheDocument();
    expect(screen.getByText('INV-001')).toBeInTheDocument();
    expect(screen.getByText('Siti')).toBeInTheDocument();
  });

  it('shows a dash instead of a fake zero when the API fails', async () => {
    get.mockRejectedValue(new Error('invalid or expired token'));

    render(<DashboardPage />);

    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('invalid or expired token'));
    expect(await screen.findAllByText('—')).toHaveLength(5);
    expect(screen.getByText('Unable to load stock data')).toBeInTheDocument();
    expect(screen.getByText('Unable to load sales data')).toBeInTheDocument();
    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });
});