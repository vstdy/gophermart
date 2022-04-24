# Gophermart loyalty program

A loyalty points program that uses points as a currency that customers can build up 
and spend. Customers can earn a points for every order, and they can then use these points 
to buy their next order at a discount.  
The project is the diploma work of the [course](https://practicum.yandex.ru/promo/go-profession/)

## REST API
### Gophermart service

By default, server starts at `8080` HTTP port with the following endpoints:

- `POST /api/user/register` — register user;
- `POST /api/user/login` — login user;
- `POST /api/user/orders` — add order to program;
- `GET /api/user/orders` — get user's orders status;
- `GET /api/user/balance` — get user's balance;
- `POST /api/user/balance/withdraw` — add withdrawal;
- `GET /api/user/balance/withdrawals` — get current user's withdrawals.

### Accrual service

- `POST /api/goods` — register cashback rule;
- `POST /api/orders` — create order to count cashback;
- `GET /api/orders/{order}` — get orders cashback.

For details check out [***http-client.http***](./http-client.http) file


## Code
### Packages used

App architecture and configuration:

- [viper](https://github.com/spf13/viper) - app configuration;
- [cobra](https://github.com/spf13/cobra) - CLI;
- [zerolog](https://github.com/rs/zerolog) - logger;

Networking:

- [go-chi](https://github.com/go-chi/chi) - HTTP router;

SQL database interface provider:

- [bun](https://github.com/uptrace/bun) - database client;

## CLI

All CLI commands have the following flags:
- `--log_level`: (optional) logging level (default: `info`);
- `--config`: (optional) path to configuration file (default: `./config.toml`);
- `--timeout`: (optional) request timeout (default: `5s`);
- `-d --database_uri`: (optional) database source name (default: `postgres://user:password@localhost:5432/gophermart?sslmode=disable`);

Root only command flags:
- `-a --run_address`: (optional) server address (default: `0.0.0.0:8080`);
- `-r --accrual_system_address`: (optional) accrual system address (default: `http://127.0.0.1:8080`);
- `-s --storage_type`: (optional) storage type (default: `psql`);

If config file not specified, defaults are used. Defaults can be overwritten using ENV variables.

### Migrations

    gophermart migrate --config ./my-confs/config-1.toml

Command migrates DB to the latest version

## How to run
### Docker

    docker-compose -f build/docker-compose.yml up
