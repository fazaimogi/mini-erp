CREATE TABLE vouchers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed')),
    discount_value NUMERIC(15,2) NOT NULL CHECK (discount_value > 0),
    min_purchase NUMERIC(15,2) NOT NULL DEFAULT 0 CHECK (min_purchase >= 0),
    usage_limit INTEGER NOT NULL DEFAULT 0 CHECK (usage_limit >= 0),
    used_count INTEGER NOT NULL DEFAULT 0 CHECK (used_count >= 0),
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- usage_limit = 0 means unlimited redemptions.
-- sale records which voucher was applied and how much it took off.
ALTER TABLE sales ADD COLUMN IF NOT EXISTS voucher_id INTEGER REFERENCES vouchers(id);
ALTER TABLE sales ADD COLUMN IF NOT EXISTS discount_amount NUMERIC(15,2) NOT NULL DEFAULT 0;