# portmantool

Port scanning and alerting tool

## Components

### (nmap "cron")

```
<program name tbd> <interval>['s'|'m'|'h'|'d'] <nmap executable> [nmap arg...]
```

Run one process that starts `nmap` and another one that sleeps for `interval`
_s_econds/_m_inutes/_h_ours/_d_ays. Then wait for `nmap` to terminate, move the
generated report to a directory shared with the next component and, finally,
wait for the sleeping process. Repeat.

### scanalyzer

This program imports the scan reports from the shared directory into a database
and, if it has been successful, archives them.

### (api server)

The API provides endpoints for retrieving scan results as well as updating the
expected state. A Prometheus `/metrics` endpoint lists any deviations.

## Database

see types.sql, schema.sql
