# portmantool

Port scanning and alerting tool

## Components

### runner

```
<program name tbd> <interval>['s'|'m'|'h'|'d'] <nmap executable> [nmap arg...]
```

Run one process that starts `nmap` and another one that sleeps for `interval`
_s_econds/_m_inutes/_h_ours/_d_ays. Then wait for `nmap` to terminate, move the
generated report to a directory shared with the next component and, finally,
wait for the sleeping process. Repeat.

### scanalyzer & api server

Previously planned as separate components

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
| /run                           | POST   | Run all scanners immediately                                                                                       |
| /scans                         | GET    | Get list of scan timestamps                                                                                        |
| /scans[/{keep}]                | DELETE | Delete entries that do not contribute to the current state and, optionally, are older than {keep}                  |
| /scan/{id}                     | GET    | Get result of scan at timestamp {id}                                                                               |
|                                |        |                                                                                                                    |

#### Metrics

| Name                                     | Description                                                                                                       |
| ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| portmantool_ports                        | Number of unique host/protocol/port combinations in database (labels: host, protocol, state)                      |
| portmantool_ports_rogue                  | Number of ports with a state different from the expected (labels: host, protocol)                                 |
| portmantool_                             |                                                                                                                   |

##### Labels

* host (e.g. "10.23.42.127", "host42.bitsbeats.io")
* protocol (e.g. "tcp", "udp")
* state (e.g. "open", "closed")

## Database

see types.sql, schema.sql
