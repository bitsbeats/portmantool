DELETE FROM
  scans
WHERE
  id NOT IN (
    SELECT
      MAX(scan_id)
    FROM
      actual_states
    GROUP BY
      address,
      port,
      protocol
  )
  AND
  id < ?
;
DELETE FROM
  actual_states
WHERE
  ROW(address, port, protocol, scan_id) NOT IN (
    SELECT
      address,
      port,
      protocol,
      MAX(scan_id)
    FROM
      actual_states
    GROUP BY
      address,
      port,
      protocol
  )
  AND
  scan_id < ?
;
