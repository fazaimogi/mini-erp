import React from 'react';
import { render, screen } from '@testing-library/react';
import Sidebar from '@/components/Sidebar';
import { User } from '@/lib/types';

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
  usePathname: () => '/dashboard',
}));

describe('Sidebar', () => {
  const adminUser: User = {
    id: 1,
    name: 'Admin User',
    email: 'admin@example.com',
    role_id: 1,
    role: 'admin',
  };

  const staffUser: User = {
    id: 2,
    name: 'Staff User',
    email: 'staff@example.com',
    role_id: 2,
    role: 'staff',
  };

  it('renders navigation links', () => {
    render(<Sidebar user={adminUser} />);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Products')).toBeInTheDocument();
    expect(screen.getByText('Inventory')).toBeInTheDocument();
    expect(screen.getByText('Customers')).toBeInTheDocument();
    expect(screen.getByText('Sales')).toBeInTheDocument();
    expect(screen.getByText('Reports')).toBeInTheDocument();
  });

  it('shows admin-only links for admin user', () => {
    render(<Sidebar user={adminUser} />);
    expect(screen.getByText('Categories')).toBeInTheDocument();
    expect(screen.getByText('Users')).toBeInTheDocument();
  });

  it('hides admin-only links for staff user', () => {
    render(<Sidebar user={staffUser} />);
    expect(screen.queryByText('Categories')).not.toBeInTheDocument();
    expect(screen.queryByText('Users')).not.toBeInTheDocument();
  });

  it('renders logout button', () => {
    render(<Sidebar user={adminUser} />);
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });

  it('handles null user', () => {
    render(<Sidebar user={null} />);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('handles role as object', () => {
    const userWithRoleObj: User = {
      ...adminUser,
      role: { id: 1, name: 'admin' },
    };
    render(<Sidebar user={userWithRoleObj} />);
    expect(screen.getByText('Categories')).toBeInTheDocument();
  });
});
