-- Seed test data for development/demo
-- Safe to run multiple times; uses NOT EXISTS guards to prevent duplicates

INSERT INTO products (name, price, tags)
SELECT 'Blue Widget', 1999, ARRAY['demo','blue']
WHERE NOT EXISTS (
  SELECT 1 FROM products WHERE name = 'Blue Widget'
);

INSERT INTO products (name, price, tags)
SELECT 'Red Gizmo', 2999, ARRAY['demo','red']
WHERE NOT EXISTS (
  SELECT 1 FROM products WHERE name = 'Red Gizmo'
);

