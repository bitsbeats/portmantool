CREATE TABLE `actual_states` (
	`address`	inet,
	`port`		integer		CHECK (`port` > 0 AND `port` < 65536),
	`protocol`	protocol,
	`state`		state		NOT NULL,
	`scan_id`	timestamptz,

	PRIMARY KEY (`address`, `port`, `protocol`, `scan_id`),
	FOREIGN KEY (`scan_id`)		REFERENCES `scans` ON UPDATE CASCADE ON DELETE CASCADE,
);

CREATE TABLE `expected_states` (
	`address`	inet,
	`port`		integer		CHECK (`port` > 0 AND `port` < 65536),
	`protocol`	protocol,
	`state`		state		NOT NULL,
	`comment`	text		NOT NULL,

	PRIMARY KEY (`address`, `port`, `protocol`),
);

CREATE TABLE `scans` (
	`id`		timestamptz	PRIMARY KEY,
);
