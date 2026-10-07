import React from 'react';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SalesFormPage from '@/app/sales/new/page';

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
  usePathname: () => '/sales/new',
}));

jest.mock('@/lib/api', () => ({
  apiClient: { get: jest.fn(), post: jest.fn() },
}));

// eslint-disable-next-line @typescript-eslint/no-require-imports
const { apiClient } = require('@/lib/api');
const get = apiClient.get as jest.Mock;

const ok = (data: unknown) => Promise.resolve({ data });

const products = [
  { id: 1, name: 'Kopi Arabika', sku: 'KP-01', category_id: 1, description: '', price: 10, stock: 3, minimum_stock: 1 },
  { id: 2, name: 'Teh Hijau', sku: 'TH-02', category_id: 1, description: '', price: 5, stock: 80, minimum_stock: 1 },
  { id: 3, name: 'Cangkir', sku: 'CG-03', category_id: 2, description: '', price: 2, stock: 0, minimum_stock: 1 },
];

const customers = [{ id: 1, name: 'Siti', email: '', phone: '', address: '' }];
const categories = [
  { id: 1, name: 'Electronics', description: '' },
  { id: 2, name: 'Furniture', description: '' },
];

/** Grand total lives in the dl row whose dt is "Total"; read the whole row. */
const cartTotal = (cart: HTMLElement) => within(cart).getByText('Total').parentElement;

describe('SalesFormPage (POS)', () => {
  beforeEach(() => {
    get.mockReset();
    get.mockImplementation((endpoint: string) => {
      if (endpoint === '/products') return ok(products);
      if (endpoint === '/customers') return ok(customers);
      if (endpoint === '/categories') return ok(categories);
      return ok([]);
    });
    localStorage.setItem(
      'user',
      JSON.stringify({ id: 1, name: 'Budi', email: 'budi@example.com', role_id: 1, role: 'admin' })
    );
  });

  it('shows category filters and product cards with price + stock', async () => {
    render(<SalesFormPage />);

    expect(await screen.findByRole('button', { name: /Add Kopi Arabika to cart/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Electronics' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Furniture' })).toBeInTheDocument();
    expect(screen.getByText('$10.00')).toBeInTheDocument();
    expect(screen.getByText('3 in stock')).toBeInTheDocument();
    expect(screen.getByText('Out of stock')).toBeInTheDocument();
  });

  it('adds a product on card click and increments quantity on re-click', async () => {
    const user = userEvent.setup();
    render(<SalesFormPage />);

    const card = await screen.findByRole('button', { name: /Add Kopi Arabika to cart/i });
    await user.click(card);

    const cart = screen.getByRole('complementary', { name: /shopping cart/i });
    expect(cartTotal(cart)).toHaveTextContent('$10.00');

    await user.click(card);
    expect(cartTotal(cart)).toHaveTextContent('$20.00');
  });

  it('caps quantity at available stock', async () => {
    const user = userEvent.setup();
    render(<SalesFormPage />);

    const card = await screen.findByRole('button', { name: /Add Kopi Arabika to cart/i });
    await user.click(card);
    await user.click(card);
    await user.click(card);
    await user.click(card); // 4th click → stock is 3

    const cart = screen.getByRole('complementary', { name: /shopping cart/i });
    expect(cartTotal(cart)).toHaveTextContent('$30.00');
    expect(within(cart).getByText(/only 3 in stock/)).toBeInTheDocument();
  });

  it('adjusts quantity with +/- and removes a line', async () => {
    const user = userEvent.setup();
    render(<SalesFormPage />);

    await user.click(await screen.findByRole('button', { name: /Add Teh Hijau to cart/i }));
    const cart = screen.getByRole('complementary', { name: /shopping cart/i });

    await user.click(within(cart).getByRole('button', { name: /Increase quantity of Teh Hijau/i }));
    expect(cartTotal(cart)).toHaveTextContent('$10.00'); // 2 × $5

    await user.click(within(cart).getByRole('button', { name: /Decrease quantity of Teh Hijau/i }));
    expect(cartTotal(cart)).toHaveTextContent('$5.00');

    await user.click(within(cart).getByRole('button', { name: /Remove Teh Hijau from cart/i }));
    expect(within(cart).getByText('Cart is empty')).toBeInTheDocument();
  });

  it('requires a customer and a non-empty cart before submit is enabled', async () => {
    const user = userEvent.setup();
    render(<SalesFormPage />);

    await user.click(await screen.findByRole('button', { name: /Add Teh Hijau to cart/i }));
    const cart = screen.getByRole('complementary', { name: /shopping cart/i });
    const submit = within(cart).getByRole('button', { name: /Create Sale/i });

    expect(submit).toBeDisabled();
    await user.selectOptions(within(cart).getByLabelText(/Customer/i), '1');
    await waitFor(() => expect(submit).toBeEnabled());
  });
});