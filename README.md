# Simple User API

`docker compose up` - starting server, postgresql, migrations

## Configuration

### Environment variables

For database settings:

- `DB_PASSWORD` - Password for PostgreSQL database connection.
- `DB_USER` - Username for PostgreSQL database connection.
- `DB_NAME` - Name of the PostgreSQL database.
- `DB_PORT` - Port for PostgreSQL database connection, default is `5432`.
- `DB_HOST` - Host for PostgreSQL database connection. For containerized applications (like Docker), it may be `db`. For local setups, it might be `localhost`.

Server settings:

- `SRV_PORT` - Port on which the server will listen for incoming requests. Default is `8080`.

## Technologies used

- Go
- PostgreSQL
- Docker, docker-compose
- sql - std library sql
- go-chi - router
- migrate - migrations
- swagger - documentation
