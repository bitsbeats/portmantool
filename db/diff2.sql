SELECT
  address,
  port,
  protocol,
  a.state state_a,
  b.state state_b,
  a.scan_id scan_a,
  b.scan_id scan_b
FROM
  (SELECT * FROM actual_states WHERE scan_id = ?) a
FULL JOIN (SELECT * FROM actual_states WHERE scan_id = ?) b
USING (
  address,
  port,
  protocol
)
WHERE
  a.state IS DISTINCT FROM b.state
;
