---
title: "Configure witches.env"
weight: 3
bookCollapseSection: true
---

Below is a detailed table describing each environment variable in the `witches.env` file:

## Environment Variables Reference

| Variable | Description | Type | Default Value |
|----------|-------------|------|---------------|
| **SERVER CONFIGURATION** |
| `APP_PORT` | HTTP server port | `int` | `8080` |
| `APP_HOST` | Server binding hostname | `string` | `localhost` |
| `METRIC_PORT` | Prometheus metrics port | `int` | `8088` |
| `UID_BITS` | Number of bits for ID generation (Snowflake/ULID) | `int` | `26` |
| `REQUEST_KEY` | Key for storing request context | `string` | `request_context` |
| **LOGGING** |
| `LOG_PATH` | Directory for log files | `string` | `./logs` |
| **CORS** |
| `CORS_ALLOW_ORIGINS` | Allowed origins (* = all) | `string` | `*` |
| `CORS_ALLOW_METHODS` | Allowed HTTP methods | `string` | `GET,POST,PUT,DELETE,OPTIONS` |
| `CORS_ALLOW_HEADERS` | Allowed request headers | `string` | `Content-Type,Authorization` |
| **TIMEZONE** |
| `TIMEZONE` | Application timezone | `string` | `UTC` |
| `SUPPORTED_LANGUAGES` | Supported languages (for i18n) | `string` | `en-US,vi-VN` |
| **SECURITY** |
| `JWT_SECRET` | Secret key for JWT signing | `string` | `your_secret_key` |
| **DATABASE CONFIGURATION** |
| `DB_DRIVER` | Database driver (postgres/mysql/mssql) | `string` | `your_driver` |
| `DB_USER` | Database username | `string` | `your_user` |
| `DB_HOST` | Database host | `string` | `localhost` |
| `DB_PORT` | Database port | `int` | `3306` |
| `DB_NAME` | Database name | `string` | `your_database` |
| `DB_PASSWORD` | Database password | `string` | `your_password` |
| **Database Connection Pool** |
| `DB_MAX_OPEN_CONNS` | Maximum number of open connections | `int` | `100` |
| `DB_MAX_IDLE_CONNS` | Maximum number of idle connections | `int` | `10` |
| `DB_CONN_MAX_LIFETIME` | Maximum lifetime of a connection (seconds) | `int` | `60` |
| `DB_CONN_MAX_IDLE_TIME` | Maximum idle time of a connection (seconds) | `int` | `600` |
| `SLOW_THRESHOLD` | Threshold for slow query logging (seconds) | `float` | `5` |
| **REDIS CONFIGURATION** |
| `REDIS_HOST` | Redis server host | `string` | `localhost` |
| `REDIS_PORT` | Redis server port | `int` | `6379` |
| `REDIS_PASSWORD` | Redis password (leave empty if none) | `string` | _empty_ |
| **TOKEN EXPIRATION** |
| `ACCESS_TOKEN_TTL` | Access token lifetime (seconds) | `int` | `900` (15 minutes) |
| `REFRESH_TOKEN_TTL` | Refresh token lifetime (seconds) | `int` | `604800` (7 days) |
| **SESSION SETTINGS** |
| `SESSION_TTL` | Session lifetime (seconds) | `int` | `604800` (7 days) |
| `IDLE_TIMEOUT` | Session idle timeout (seconds) | `int` | `1800` (30 minutes) |
| **BLACKLIST SETTINGS** |
| `REVOKED_TTL` | Revoked token storage time (seconds) | `int` | `300` (5 minutes) |
| **RATE LIMIT** |
| `RATE_LIMIT_PERIOD` | Rate limit period (seconds) | `int` | `60` |
| `RATE_LIMIT_MAX` | Maximum requests per period | `int` | `100` |

---

## Notes

- **Time values**: All TTL/time values are in **seconds**.
- **DB_DRIVER**: Valid values: `postgres`, `postgresql`, `mysql`, `mssql`, `sqlserver`.
- **CORS_ALLOW_ORIGINS**: Use `*` for all origins, or comma-separate multiple origins.
- **JWT_SECRET**: Should be changed to a strong secret key in production environments.