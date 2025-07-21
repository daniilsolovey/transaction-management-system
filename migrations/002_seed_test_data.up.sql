INSERT INTO transactions (user_id, transaction_type, amount, timestamp) VALUES
  ('11111111-1111-1111-1111-111111111111', 'bet',  50.00, NOW() - INTERVAL '5 minutes'),
  ('11111111-1111-1111-1111-111111111111', 'win', 100.00, NOW() - INTERVAL '3 minutes'),
  ('22222222-2222-2222-2222-222222222222', 'bet',  25.50, NOW() - INTERVAL '2 minutes'),
  ('33333333-3333-3333-3333-333333333333', 'win', 200.00, NOW() - INTERVAL '1 minute');
