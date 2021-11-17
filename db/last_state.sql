SELECT	a.address,
	a.port,
	a.protocol,
	state
FROM actual_states a
JOIN (
	SELECT	address,
		port,
		protocol,
		MAX(scan_id) max_scan_id
	FROM actual_states
	GROUP BY address, port, protocol
) b
ON	a.address = b.address
	AND a.port = b.port
	AND a.protocol = b.protocol
	AND scan_id = max_scan_id
;
