# Database idle mode

This fork enables `DATABASE_IDLE_MODE=true` in its Docker image so a
low-traffic, single-node deployment can let a serverless PostgreSQL database
(for example Neon) scale to zero between real requests.

The mode removes periodic database activity from:

- option, channel-cache, authorization-policy, and task-plugin synchronization;
- system-instance heartbeats;
- scheduled system-task polling (channel tests, upstream model updates, and
  asynchronous task polling);
- subscription quota maintenance and Codex credential pre-refresh;
- expired dashboard-auth artifact cleanup;
- model performance-metric collection and retention cleanup;
- the periodic quota-dashboard flush timer. Request-generated quota data is
  flushed immediately while the request already has the database awake.

The system-task runner remains available in a wake-driven form, so tasks
created by this same process through the management API can still run without
polling the database every 15 seconds. The normal `/api/status` Docker health
check reads in-memory state and does not ping the database.

## Trade-offs

Use this mode only for a single application instance. Changes made by another
instance are not periodically reloaded. Scheduled channel tests and upstream
model updates do not run, subscription resets are not automatic, asynchronous
media-task polling is paused, and the System Info page does not receive live
instance heartbeats. Restart the container after out-of-band database changes.

Actual API traffic, dashboard operations, and `/api/status/test` still access
the database and will wake it as expected.

To restore the upstream scheduling behavior:

```yaml
environment:
  DATABASE_IDLE_MODE: "false"
```
