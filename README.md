*Sorry there is no unit test since i was focusing on the essential functions first and also due to the limited time.
also i didn't put much time on detailed validations.

The app and database (postgres) are already dockerized so you can just run it using docker compose.

## How to run

```shell
docker compose build
```

### reload
```shell
docker compose down && docker compose up -d
```

```shell
docker compose up -d && docker compose logs -f
``` 

### access database
```shell
docker exec -it postgres psql -U postgres -d moneydb
```