# portmantool

Port scanning and monitoring tool

## Components

### runner

```
while true
do
	run.sh <nmap args...> &
	sleep <scan interval>
	wait $!
done
```

### scanalyzer

This program imports the scan reports from the shared directory into a database
and, if it has been successful, archives them.

The API provides endpoints for retrieving scan results as well as updating the
expected state. A Prometheus `/metrics` endpoint provides metrics useful for
alerting.

#### v1 Endpoints

| Path (excl. /v1)               | Method | Description                                                                                                        |
| ------------------------------ | ------ | ------------------------------------------------------------------------------------------------------------------ |
| /diff                          | GET    | Compute difference between current and expected state                                                              |
| /diff/{id1}[/{id2}]            | GET    | Compute difference between scans {id1} and {id2}, if it is given, or the expected state, otherwise                 |
| /expected                      | GET    | Get expected state                                                                                                 |
| /expected                      | PATCH  | Update expected state                                                                                              |
| /hello                         | *      | hello                                                                                                              |
| /run[/{id}]                    | POST   | Run scanner {id} or, in case {id} is not given, all scanners immediately                                           |
| /scans                         | GET    | Get list of scan timestamps                                                                                        |
| /scans[/{keep}]                | DELETE | Delete entries that do not contribute to the current state and, optionally, are older than {keep}                  |
| /scan/{id}                     | GET    | Get result of scan at timestamp {id}                                                                               |
|                                |        |                                                                                                                    |

#### Metrics

| Name                                     | Description                                                                                                       |
| ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| portmantool_ports                        | Number of unique host/protocol/port combinations in database (labels: host, protocol, state)                      |
| portmantool_ports_rogue                  | Number of ports with a state different from the expected (labels: host, protocol)                                 |
| portmantool_imports_failed_total         | Total number of failed imports                                                                                    |
| portmantool_                             |                                                                                                                   |

##### Labels

* host (e.g. "10.23.42.127", "host42.bitsbeats.io")
* protocol (e.g. "tcp", "udp")
* state (e.g. "open", "closed")

## Database

see db/types.sql, db/schema.sql
