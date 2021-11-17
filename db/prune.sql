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
;
DELETE FROM
  actual_states
WHERE
  --
;
