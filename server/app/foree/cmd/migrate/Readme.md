Prerequisite:
1. Migrate CLI: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Migrate CLI Cheet Sheet:
- Create migration file:
    migrate create -ext sql -dir migrate/migrations {{TABLE_NAME}}

Run CLI:
- go build -o main && ./main {{args}}