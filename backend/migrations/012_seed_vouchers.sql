-- Demo vouchers so the checkout flow has something to apply.
INSERT INTO vouchers (code, discount_type, discount_value, min_purchase, usage_limit, expires_at, is_active)
SELECT v.code, v.discount_type, v.discount_value, v.min_purchase, v.usage_limit, v.expires_at, v.is_active
FROM (VALUES
    ('PROMO2026', 'percentage', 10.00, 100.00, 100, NULL::timestamptz, TRUE),
    ('FLAT50',    'fixed',      50.00, 200.00,   0, NULL::timestamptz, TRUE)
) AS v(code, discount_type, discount_value, min_purchase, usage_limit, expires_at, is_active)
WHERE NOT EXISTS (SELECT 1 FROM vouchers WHERE vouchers.code = v.code);