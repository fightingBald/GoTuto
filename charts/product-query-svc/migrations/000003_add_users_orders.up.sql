CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email CITEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create orders table referencing users
CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  product_name TEXT NOT NULL,
  total BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS orders_user_id_idx ON orders(user_id);

-- Seed sample data
INSERT INTO users (name, email)
VALUES ('Alice', 'alice@example.com'),
       ('Bob', 'bob@example.com')
ON CONFLICT (email) DO NOTHING;

INSERT INTO orders (user_id, product_name, total)
SELECT u.id, p.product_name, p.total
FROM users u
JOIN (VALUES
  ('alice@example.com', 'Starter Pack', 4999),
  ('alice@example.com', 'Premium Pack', 19999),
  ('bob@example.com', 'Gift Card', 2500)
) AS p(email, product_name, total)
  ON u.email = p.email
LEFT JOIN orders o ON o.user_id = u.id AND o.product_name = p.product_name
WHERE o.id IS NULL;
