-- Seed default admin user
-- Email: admin@example.com
-- Password: password (bcrypt hashed)
INSERT INTO users (name, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    'Admin',
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoHyemKbz8t8.Ks0Nm8J.4e8C8Zy5YO4gVJG', -- hashed 'password'
    1, -- admin role
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE SET role_id = 1;
