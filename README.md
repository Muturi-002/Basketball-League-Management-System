# Basketball-League-Management-System
School project for managing a basketball league, with a Go backend, a static frontend, and Oracle ATP database connectivity.

## Directory Structure

```text
Basketball-League-Management-System/
├── Backend/  # Go backend for routing, authentication, and database access
│   ├── admin.go  # Admin page handlers and admin-facing logic
│   ├── admin_auth.go  # Admin authentication and access control helpers
│   ├── go.mod  # Go module definition and dependency list
│   ├── go.sum  # Checksums for Go module dependencies
│   ├── web.go  # HTTP server setup and route registration
│   └── Database/  # Database access layer and SQL helpers
│       ├── auth.go  # Database-backed authentication queries
│       ├── db_connect.go  # Oracle database connection setup
│       ├── fixtures.go  # Fixture management queries
│       ├── injury.go  # Injury tracking queries
│       ├── managers.go  # Manager management queries
│       ├── players.go  # Player management queries
│       ├── query_helpers.go  # Shared query and scan helpers
│       ├── stats.go  # Statistics and reporting queries
│       ├── teams.go  # Team management queries
│       ├── validation.go  # Input validation logic for database operations
│       └── sql/  # SQL scripts for schema and data setup
│           ├── blms_league_automation.sql  # Automation script for league database tasks
│           ├── blms_schema_create.sql  # Schema creation script
│           └── blms_schema_populate.sql  # Seed data population script
├── Frontend/  # Static HTML and CSS files for the user interface
│   ├── admin-login.html  # Admin login page
│   ├── admin.html  # Admin dashboard page
│   ├── auth.html  # Authentication-related page
│   ├── fixtures.html  # Fixtures listing page
│   ├── home.html  # Landing page for the application
│   ├── injury.html  # Injury management page
│   ├── manager.html  # Manager management page
│   ├── players.html  # Player management page
│   ├── stats.html  # Statistics and standings page
│   ├── styles.css  # Shared stylesheet for the frontend
│   ├── team.html  # Single team details page
│   └── teams.html  # Teams overview page
```

The frontend is served from the HTML and CSS files in the `Frontend/` directory, while the backend handles routing, authentication, and database access.


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
go run .
```

Persistent options

Notes and troubleshooting
- The project uses the `go-ora` driver which does not require the Oracle Instant Client native libraries for typical TCP/TCPS connections. Ensure `DB_USER`, `DB_PASSWORD`, and `ORACLE_CONNECTION_STRING` are set correctly. If you are using a wallet-based (TCPS) connection, set `ORACLE_CONFIG_DIR` to the wallet directory.
- See the database connector source at [Backend/Database/db_connect.go](Backend/Database/db_connect.go) for how `ORACLE_CONFIG_DIR` is used by the application.

