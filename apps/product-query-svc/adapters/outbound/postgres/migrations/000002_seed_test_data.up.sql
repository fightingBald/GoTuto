-- Seed initial product test data for development and integration tests
INSERT INTO products (name, price, tags) VALUES
  ('Basic Plan', 9900, ARRAY['starter', 'subscription']),
  ('Pro Plan', 19900, ARRAY['professional', 'subscription']),
  ('Enterprise Plan', 49900, ARRAY['enterprise', 'subscription']);

