```shell
docker compose build
```

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