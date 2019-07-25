# Database migrations

gAPI database is based on migrations.
Oracle migrations are stored inside folder _migrations/oracle_.

## How it works

There is a migration table which stores all migrations already run.
When the server starts, all new migrations are run and the table updated.

## Enable/Disable

Database migrations can be enabled/disabled by adding an environment variable: **RUN_MIGRATIONS**

## Migration naming convention

```
<current_timestamp>_<description_separated_by_underscore>.sql
```

Example:

```
1557849123509_create_users_table.sql
```
