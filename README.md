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

This program receives scan reports from the runner and imports them into a
database.

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
| /info                          | GET    | Get server info (currently only timestamp of last successful import)                                               |
| /run[/{id}]                    | POST   | **(NYI)** Run scanner {id} or, in case {id} is not given, all scanners immediately                                 |
| /scans                         | GET    | Get list of scan timestamps                                                                                        |
| /scans[/{keep}]                | DELETE | Delete entries that do not contribute to the current state and, optionally, are older than {keep} (UNIX timestamp) |
| /scan                          | POST   | Upload a new scan report                                                                                           |
| /scan/{id}                     | GET    | Get result of scan at timestamp {id}                                                                               |
|                                |        |                                                                                                                    |

#### Metrics

| Name                                     | Description                                                                                                       |
| ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| portmantool_ports                        | Number of unique host/protocol/port combinations in database (labels: host, protocol, state)                      |
| portmantool_ports_rogue                  | Number of ports with a state different from the expected (labels: host, protocol)                                 |
| portmantool_imports_failed_total         | Total number of failed imports since server was running                                                           |
| portmantool_imports_last                 | Timestamp of last successful import                                                                               |
| portmantool_                             |                                                                                                                   |

##### Labels

* host (e.g. "10.23.42.127", "host42.bitsbeats.io")
* protocol (e.g. "tcp", "udp")
* state (e.g. "open", "closed")

### web

* Show diff between expected and current state, updated every 5 seconds
* Show expected state, add and update independent of diff
* Show (list of) scan results, prune obsolete

#### TODO

* Compare scan(s) to current state (backend currently compares to expected state)/each other

#### Backlog

* Edit expected state
  * Delete (needs implementation in backend)

## Database

see db/types.sql, db/schema.sql

## [License](https://github.com/bitsbeats/portmantool/blob/main/LICENSE)

```
Copyright 2020-2022 Thomann Bits & Beats GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
