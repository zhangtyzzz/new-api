# Database idle mode

Set `DATABASE_IDLE_MODE=true` explicitly so a low-traffic, single-node
deployment can let a serverless PostgreSQL database
(for example Neon) scale to zero between real requests.

The image and Compose defaults remain `false` because the mode has operational
trade-offs. When enabled, it removes periodic database activity from:

- option, channel-cache, authorization-policy, and task-plugin synchronization;
- system-instance heartbeats;
- independent scheduled system-task polling (channel tests, upstream model
  updates, and asynchronous task polling become request-driven);
- independent subscription quota, Codex credential refresh, and expired
  dashboard-auth timers (these become request-driven);
- model performance-metric collection and retention cleanup;
- the periodic quota-dashboard flush timer. Request-generated quota data is
  flushed immediately while the request already has the database awake;
- idle SQL connections: after startup, `SQL_MAX_IDLE_CONNS` defaults to `0` in
  this mode so the application does not hold the serverless compute open after
  a request. An explicit `0` is honored after initialization while startup
  migrations temporarily reuse a small connection pool;
- batch quota updates: `BATCH_UPDATE_ENABLED=true` is ignored in this mode and
  quota writes stay synchronous instead of relying on a background flush loop.

The system-task runner and maintenance jobs remain available in a wake-driven
form. After real API, dashboard, or relay traffic has already woken the
database, due subscription resets, Codex credential refresh, auth cleanup,
channel tests, upstream-model updates, asynchronous task polling, and stale
task-lock cleanup run with their normal throttles. Retrying an active manual
task also wakes the runner, allowing an expired lock left by a crash to recover.
The normal `/api/status` Docker health check reads in-memory state and does not
trigger maintenance or ping the database.

## Trade-offs

Use this mode only for a single application instance. Changes made by another
instance are not periodically reloaded. Maintenance is delayed until the next
real request, so it is not suitable when subscription boundaries or background
media tasks must progress while the service has no traffic. The System Info
page also does not receive live instance heartbeats. Restart the container
after out-of-band database changes.

Actual API traffic, dashboard operations, and `/api/status/test` still access
the database and will wake it as expected.

To restore the upstream scheduling behavior:

```yaml
environment:
  DATABASE_IDLE_MODE: "false"
```
