# Using `migrate`

Ensure that the working directory is `migrate` so the binary can discover the migration files, e.g.

```bash
cd backend/cmd/migrate
go run migrate.go -template -c path/to/my/clutch-config.yaml
```

Note: `migrate.go` accepts the same arguments as the main Clutch binary.
