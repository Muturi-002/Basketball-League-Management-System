# Basketball-League-Management-System
School Project. For the fans of the game, by the fans of the game,


### Steps to connect to the Oracle 19c ATP Database successfully

Prerequisites
- Install Oracle Instant Client on the machine or provide a repo-local wrapper directory named `oracle-client` containing the Instant Client shared libraries (e.g. `libclntsh.so*`) and any compatibility symlinks such as `libaio64`. Use this [link](https://www.oracle.com/database/technologies/instant-client/linux-x86-64-downloads.html) to download the package.
- Exttract contents from the downloaded zip file and move to a new folder in the root directory of the project (e.g., `oracle-client`).

Environment variables
 - A project `.env` is used by the code loader. Set the following values in your project `.env` (or in your shell):

```
DB_USER=your_db_user
DB_PASSWORD=your_db_password
ORACLE_CONNECTION_STRING="mydb_high"
ORACLE_DRIVER=godror
ORACLE_CONFIG_DIR=/full/path/to/wallet
# Path to the Instant Client directory or repo wrapper (recommended)
ORACLE_LIB_DIR=/full/path/to/oracle-client
```

Recommended startup (developer-friendly)
- Preferred: export `LD_LIBRARY_PATH` before starting the server so the OS dynamic loader can find the Oracle native libraries. Example (from project root):

```bash
export LD_LIBRARY_PATH=/full/path/to/oracle-client:${LD_LIBRARY_PATH:-}
cd Backend
go run web.go
```

- The application also includes a best-effort re-exec that attempts to set `LD_LIBRARY_PATH` from `ORACLE_LIB_DIR` (or `../oracle-client`) when starting `Backend/web.go`. This can help in many setups, but the most reliable method is to export `LD_LIBRARY_PATH` externally before launching the process.

Persistent options
- To make the library path persistent for your user session, add the `export LD_LIBRARY_PATH=...` line to your shell startup (for example `~/.bashrc` or `~/.profile`).
- For system-wide persistence (requires root), add the Instant Client path to `/etc/ld.so.conf.d/oracle-client.conf` and run `sudo ldconfig`.

Notes and troubleshooting
- If you see errors mentioning `DPI-1047` or `libclntsh.so` not found, it means the native Oracle client libraries were not visible to the dynamic loader when the process started. Ensure `LD_LIBRARY_PATH` contains the Instant Client path before launching.
- See the database connector source at [Backend/Database/db_connect.go](Backend/Database/db_connect.go) for how `ORACLE_LIB_DIR` and `ORACLE_CONFIG_DIR` are used by the application.

