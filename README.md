# Basketball-League-Management-System
School Project. For the fans of the game, by the fans of the game,


### Steps to connect to the Oracle 19c ATP Database successfully

Environment variables
- A project `.env` is used by the code loader. Set the following values in your project `.env` (or in your shell):

```
DB_USER=your_db_user
DB_PASSWORD=your_db_password
ORACLE_CONNECTION_STRING="mydb_high"
ORACLE_DRIVER=oracle
ORACLE_CONFIG_DIR=/full/path/to/wallet
```

Recommended startup (developer-friendly)

From the project root run:

```bash
cd Backend
go run web.go
```

Persistent options

Notes and troubleshooting
- The project uses the `go-ora` driver which does not require the Oracle Instant Client native libraries for typical TCP/TCPS connections. Ensure `DB_USER`, `DB_PASSWORD`, and `ORACLE_CONNECTION_STRING` are set correctly. If you are using a wallet-based (TCPS) connection, set `ORACLE_CONFIG_DIR` to the wallet directory.
- See the database connector source at [Backend/Database/db_connect.go](Backend/Database/db_connect.go) for how `ORACLE_CONFIG_DIR` is used by the application.

