# Go UI

A simple Go boilerplate for building a server rendered web app.

## Includes

- Sqlite DB
- Simple user signup and login flows
- Dockerfile for deployment
- Tests
- Database migration cli

## Docker

```sh
export PORT=4000
export DB="./db.sqlite"

# Build
docker build . \
--file ./docker/Dockerfile \
--tag go-ui

# Run migration commands
docker run \
--publish $PORT:$PORT \
--volume $(pwd)/$DB:/app/$DB \
go-ui migrate -database-dsn $DB up

# Run the web server
docker run \
--publish $PORT:$PORT \
--volume $(pwd)/$DB:/app/$DB \
go-ui web -addr ":$PORT" -database-dsn $DB
```