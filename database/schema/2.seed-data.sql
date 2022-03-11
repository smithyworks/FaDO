INSERT INTO policies (name, default_value)
VALUES
  ('lb_policy',             '"round_robin"'),
  ('lb_match_header',       '"X-FaDO-Bucket"'),
  ('lb_upstreams',          '[]'),
  ('lb_routes',             '{}'),
  ('lb_route_overrides',    '{}'),
  ('caddy_config',          '""'),
  ('replica_locations', '[]'),
  ('target_replica_count',  '0'),
  ('zones',                 '[]');
