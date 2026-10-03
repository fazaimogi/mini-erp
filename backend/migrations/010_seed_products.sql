-- Seed starter products for the Electronics and Furniture categories.
-- category_id is resolved by category name rather than hard-coded IDs so this
-- works regardless of the sequence values the 009 seed landed on.
-- products.sku is UNIQUE, so ON CONFLICT keeps a hand-run `psql -f` idempotent;
-- a category that is missing from the DB simply yields no row for that product.
INSERT INTO products (name, sku, category_id, description, price, stock, minimum_stock, created_at, updated_at)
SELECT seed.name, seed.sku, c.id, seed.description, seed.price, seed.stock, seed.minimum_stock, NOW(), NOW()
FROM (VALUES
    ('Laptop Pro 15', 'ELEC-001', 'Electronics', 'Laptop high-performance untuk developer dan desainer.', 1200, 25, 5),
    ('Wireless Mechanical Keyboard', 'ELEC-002', 'Electronics', 'Keyboard RGB switch mekanikal dengan koneksi bluetooth.', 150, 50, 10),
    ('UltraWide Monitor 29 inch', 'ELEC-003', 'Electronics', 'Monitor IPS resolusi tinggi untuk produktivitas multitasking.', 350, 15, 3),
    ('Ergonomic Office Chair', 'FURN-001', 'Furniture', 'Kursi kantor ergonomis dengan penopang lumbar yang nyaman.', 250, 20, 4),
    ('Minimalist Standing Desk', 'FURN-002', 'Furniture', 'Meja kerja modern dengan fitur pengatur ketinggian.', 450, 10, 2),
    ('Wooden Bookshelf 4 Tier', 'FURN-003', 'Furniture', 'Rak buku kayu minimalis serbaguna untuk ruang kerja.', 120, 30, 5)
) AS seed(name, sku, category, description, price, stock, minimum_stock)
JOIN categories c ON c.name = seed.category
ON CONFLICT (sku) DO NOTHING;