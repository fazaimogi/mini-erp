-- Seed default admin user
-- Email: admin@example.com
-- Password: password (bcrypt hashed)
INSERT INTO users (name, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    'Admin',
    'admin@example.com',
    '$2a$10$sPssmhG7jXCpwae435eNz.apYqOGj8weSiyMEEO64RoA0w/7U8FGu', -- hashed 'password'
    1, -- admin role
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO UPDATE SET role_id = 1;
