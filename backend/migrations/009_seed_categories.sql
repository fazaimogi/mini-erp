-- Seed starter categories so the product form has options on a fresh database.
-- categories.name has no unique constraint, so guard with NOT EXISTS: this
-- migration is already protected by schema_migrations, but keeping the insert
-- idempotent means a hand-run `psql -f` cannot create duplicates either.
INSERT INTO categories (name, created_at, updated_at)
SELECT name, NOW(), NOW()
FROM (VALUES ('Electronics'), ('Furniture')) AS seed(name)
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE categories.name = seed.name);