# hbd

Birthday bot backend and frontend

## Migrations

### Up

```bash
migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations up
```

### Down

```bash
migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations down
```

### Down to last migration

```bash
migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations down 1
```

## SQLBoiler

Generate models

```bash
sqlboiler psql
```

Point to the `.sqlboiler.toml` file

```bash
sqlboiler psql --config .sqlboiler.toml
```

## All together

```bash
migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations down && migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations up && sqlboiler psql --config .sqlboiler.toml && air
```
