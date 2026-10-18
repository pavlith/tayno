# Sync server

```text
apis/       Gin HTTP handlers, JWT middleware, DTO conversion and OpenAPI
  generated/  Generated API types and Go client
core/       Use cases, access checks, repository interfaces and errors
models/     User, Keystore and Invitation entities with GORM schema tags
storage/    SQLite connection, migrations and repository implementations
```

`cmd/server` assembles these layers. The existing HTTP routes and JSON formats
are retained, including the invitation acceptance/finalization workflow and
opaque encrypted keystore payloads.

## Run

```sh
go run ./cmd/server --db data/myst.sqlite --port 8080 --jwt-signing-key '<base64-key>'
```

Flags can also be provided using the existing environment variable support:
`DB`, `PORT`, `JWT_SIGNING_KEY`. The default database path is `data/myst.sqlite`;
its parent directories are created automatically. The signing key is a
base64-encoded byte string and should remain stable across restarts.

The server loads `.env` from the current working directory using
[`godotenv`](https://github.com/joho/godotenv). The file is optional; invalid
contents cause startup to fail. Priority is: command-line flags, existing
environment variables, `.env`, defaults. `.env.development` is not loaded
automatically.

For example, create `.env` in the repository root:

```dotenv
PORT=8080
DB=data/myst.sqlite
JWT_SIGNING_KEY='<base64-key>'
```

Then run `go run ./cmd/server`. Variables named `MYST_SERVER_ADDRESS` and
`MYST_SERVER_NETWORK_TYPE` are not used by this entry point.

The official GORM SQLite driver uses CGO. Build with `CGO_ENABLED=1` and a C
compiler. Cross compilation requires a C compiler for the target platform;
the Go builder defaults to the host OS.

## Persistence

`storage.Open` enables WAL, a five-second busy timeout and foreign keys, then
runs GORM `AutoMigrate`. The pool uses one connection to serialize writes and
support `:memory:` databases. Production storage should use a file.

Users have unique usernames. Invitations have foreign keys to their keystore,
inviter and invitee, plus a partial unique index allowing only one active
invitation for each keystore/inviter/invitee combination. Deleting a keystore
removes its invitations in the same transaction. Invitation `DeletedAt` is a
workflow timestamp, not a GORM soft-delete field: deleted invitations remain
available to the inviter under the existing core access rules.

MongoDB and the old in-memory repository have been removed. Existing MongoDB
data is not imported automatically. `AutoMigrate` handles the SQLite schema;
it does not convert MongoDB documents.

## Tests

```sh
go test ./server/... ./test/...
```

Integration tests use a separate temporary SQLite file for each test, with no
MongoDB setup. Storage tests cover reopening/migration, binary data, invitation
timestamps, zero-value updates, missing records, foreign keys, deletion and
concurrent username uniqueness.

Reference: [GORM update semantics](https://gorm.io/docs/update.html). Repository
updates explicitly select columns, so zero values are persisted and updates of
missing records do not silently create new rows.
