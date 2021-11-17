SELECT
  address,
  port,
  protocol,
  a.state expected_state,
  b.state actual_state,
  scan_id,
  comment
FROM
  expected_states a
FULL JOIN (
  SELECT
    actual_states.*
  FROM
    actual_states
  JOIN (
    SELECT
      address,
      port,
      protocol,
      MAX(scan_id) max_scan_id
    FROM
      actual_states
    GROUP BY
      address,
      port,
      protocol
  ) latest_scans
  ON
    actual_states.address = latest_scans.address
    AND
    actual_states.port = latest_scans.port
    AND
    actual_states.protocol = latest_scans.protocol
    AND
    scan_id = max_scan_id
) b
USING (
  address,
  port,
  protocol
)
WHERE
  a.state IS DISTINCT FROM b.state
;
