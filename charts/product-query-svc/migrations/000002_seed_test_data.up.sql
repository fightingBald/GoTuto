-- Seed test data for development/demo
-- Inserts are conditional to avoid duplicates on re-run

INSERT INTO products (name, price, tags)
SELECT 'Basic Plan', 9900, ARRAY['starter','subscription']
WHERE NOT EXISTS (
  SELECT 1 FROM products WHERE name = 'Basic Plan'
);

INSERT INTO products (name, price, tags)
SELECT 'Pro Plan', 19900, ARRAY['professional','subscription']
WHERE NOT EXISTS (
  SELECT 1 FROM products WHERE name = 'Pro Plan'
);

INSERT INTO products (name, price, tags)
SELECT 'Enterprise Plan', 49900, ARRAY['enterprise','subscription']
WHERE NOT EXISTS (
  SELECT 1 FROM products WHERE name = 'Enterprise Plan'
);
