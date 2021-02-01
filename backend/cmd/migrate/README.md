<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Using `migrate`](#using-migrate)
  - [Migrating Up](#migrating-up)
  - [Migrating Down](#migrating-down)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Using `migrate`

## Migrating Up

Ensure that the working directory is `migrate` so the binary can discover the migration files, e.g.

```bash
cd backend/cmd/migrate
go run migrate.go -template -c path/to/my/clutch-config.yaml
```

Note: `migrate.go` accepts the same arguments as the main Clutch binary.

Migrate also accepts an `-f` option to skip user confirmation of migration (useful for CI).

## Migrating Down

To migrate down supply the `-down` flag.

A single invocation of the migrate down will only ever migrate one version down.
If you need to migrate down multiple versions, invoke the tool until you reach the desired version.

```
cd backend/cmd/migrate
go run migrate.go -template -down -c path/to/my/clutch-config.yaml
```

## Federated Migrations

:warning: **IT IS NOT SAFE TO MODIFY THE SAME TABLES FROM THE FEDERATED MIGRATIONS AND BUILT-IN MODULE MIGRATIONS. IT WILL CORRUPT THE DATABASE.**

If you wish to run migrations for independent tables from federated code, provide the `-namespace` argument and `-migrationDir` argument.

For example, if I had created a scaffolded application in `clutch-custom-gateway/`, I could create migrations for a new module's table in `clutch-custom-gateway/migrations`,
e.g. `000001_create_my_new_table.up.sql` and `000001_create_my_new_table.down.sql`, and run the migrations using
```bash
$ go run migrate.go -c path/to/my/clutch-config.yaml -namespace custom_gateway -migrationDir path/to/clutch-custom-gateway/migrations
```
