# portmantool

Port scanning and alerting tool

## Components

### (nmap "cron")

```
<program name tbd> <interval>['s'|'m'|'h'|'d'] [nmap arg...]
```

Run one process that starts `nmap` and another one that sleeps for `interval` _s_econds/_m_inutes/_h_ours/_d_ays. Then wait for `nmap` to terminate, move the generated report to a directory shared with the next component and, finally, wait for the sleeping process. Repeat.

### (report analysis)

Another program imports the scan reports from the shared directory into a database and, if it has been successful, deletes (or archives) them. If the currently processed state is found to be identical to the previous state, no new data will be inserted into the database and only metadata will be updated. This data is then compared to the expected state and any differences are flagged.

### (api server)

The API provides endpoints for retrieving scan results as well as updating the expected state. A Prometheus /metrics endpoint lists any deviations.
