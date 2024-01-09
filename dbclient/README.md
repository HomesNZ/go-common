# dbclient/v4

This package provides a client for interacting with a PostgreSQL database using the pgxpool package.

# Usage

To use this package, you need to import it in your Go code:
```
import "github.com/HomesNZ/go-common/dbclient/v4"
```

Then, you can create a new database client useing `New` or `NewFromEnv`

# Configuration
The Connect function takes a *config.Config parameter which is used to configure the database connection. Here are the available configuration options:

DB_HOST: The host of the database.
DB_USER: The user to connect to the database.
DB_NAME: The name of the database.
DB_PASSWORD: The password to connect to the database.
DB_PORT: The port of the database. // default 5432
DB_MAX_CONNECT: The maximum number of connections in the pool. // default 3
DB_HEALTH_CHECK_PERIOD: seconds - how often to check health of the connection // default 30
DB_MAX_CONN_IDLE_TIME: mins - how long connection can be idle before it'll be closed // default 5
DB_PING_BEFORE_USE:  if true, t'll be used to check connection before use and if the connection is not alive, it'll be reconnected