# Auth Service

## Docker Postgres Instance for Testing

Postgres docker container for local testing of the Auth microservice.

```bash
sudo docker run --rm --name postgres \
  -p 8080 \
  -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=auth \
  --network=host postgres:alpine3.18
```

## Create User

```bash
curl -X POST http://localhost:8090/api/v1/users -H 'content-type: application/json' -H 'Authorization: Bearer <token>' -d '{"email":"wakeward@control-plane.io","firstName":"wake","lastName":"ward","password":"simple", "secret": "secret"}'
```

## Update User

```bash
curl -X PUT http://localhost:8090/api/v1/users/1 -H 'application/json' -d '{"email":"wakeward@control-plane.io","firstName":"wake","lastName":"change","password":"simple"}'
```

## Delete User

```bash
curl -X DELETE http://localhost:8090/api/v1/users/1 -H 'application/json' -d '{"email":"wakeward@control-plane.io", "password":"simple"}'
```
