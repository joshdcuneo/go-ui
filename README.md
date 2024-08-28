# Go UI

A simple Go boilerplate for building a server rendered web app.

## Includes

- Sqlite DB
- Simple user signup and login flows
- Dockerfile for deployment
- Tests
- Database migration cli

## Docker

```bash
# Build
docker build . --tag go-ui

# Run migration commands
docker run \
  --publish 4000:4000 \
  --volume $(pwd)/db.sqlite:/app/db.sqlite \
  --entrypoint /app/migrate \
  go-ui up

# Run the web server
docker run \
  --publish 4000:4000 \
  --volume $(pwd)/db.sqlite:/app/db.sqlite \
  go-ui
```