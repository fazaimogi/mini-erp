import React from 'react';
import { render, screen } from '@testing-library/react';
import Header from '@/components/Header';
import { User } from '@/lib/types';

describe('Header', () => {
  const mockUser: User = {
    id: 1,
    name: 'John Doe',
    email: 'john@example.com',
    role_id: 1,
    role: 'admin',
  };

  it('renders title', () => {
    render(<Header title="Dashboard" user={mockUser} />);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('renders user name', () => {
    render(<Header title="Dashboard" user={mockUser} />);
    expect(screen.getByText(/John Doe/)).toBeInTheDocument();
  });

  it('handles null user gracefully', () => {
    render(<Header title="Dashboard" user={null} />);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('renders action when provided', () => {
    const action = <button>Test Action</button>;
    render(<Header title="Dashboard" user={mockUser} action={action} />);
    expect(screen.getByText('Test Action')).toBeInTheDocument();
  });

  it('handles role object structure', () => {
    const userWithRoleObj: User = {
      ...mockUser,
      role: { id: 1, name: 'admin' },
    };
    render(<Header title="Dashboard" user={userWithRoleObj} />);
    expect(screen.getByText(/John Doe/)).toBeInTheDocument();
  });
});
