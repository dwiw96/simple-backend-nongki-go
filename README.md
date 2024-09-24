# How To Run
- run postgres inside docker:
  ```
  sudo docker compose -f ./docker-compose.yml up -d
  ```
- run redis inside docker:
  ```
  sudo docker run --name redis-nongki -p 6379:6379 -d redis
  ```
- create table for postgres, if you use Makefile:
  ```
  make pgMigrateUsers
  make pgMigratePrivKey
  ```
  run without Makefile:
  ```
  sudo docker exec -i pg_nongki_container psql -U dwiw nongki_db < internal/utils/postgres/users.sql
  sudo docker exec -i pg_nongki_container psql -U dwiw nongki_db < internal/utils/postgres/priv_key.sql
  ```
- run the program
  linux/mac:
  ```
  ./simple-backend-nongki-go
  ```
  windows:
  ```
  simple-backend-nongki-go.exe
  ```

# Endpoint
you can see docs/user-api.json for open api specification
