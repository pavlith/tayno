# Myst

Encrypted keystores with a Go client, a Svelte web interface and a sync server.

```text
cmd/            Go entry points: client, server, builder
client/
  golang/       Client: apis, core, models, storage, remote and scheduler
  web/          Svelte 5 web client
server/         Sync server: apis, core, models and GORM/SQLite storage
pkgs/           Shared Go packages and build utilities
test/           Integration tests
script/         Development and build scripts
```

The repository uses one Go module (`myst`) in the root.
See [client/golang/README.md](client/golang/README.md) for the Go client layers.

Run commands from the repository root:

```sh
npm --prefix client/web ci
npm --prefix client/web run dev
npm --prefix client/web run check
npm --prefix client/web run build
go run ./cmd/server
go test ./pkgs/... ./client/golang/... ./server/...
```

The web development server uses port 8082 and calls the Go client API on port
8081. The sync server uses port 8080. Run these in separate terminals:

```sh
go run ./cmd/server
go run ./cmd/client --remote http://localhost:8080/api --port 8081
npm --prefix client/web run dev
```

Set `PORT=8080` in the server's `.env`. The `/api/enclave` endpoint belongs to
the Go client, so the web UI must connect to the client on 8081, not directly
to the sync server. Restart the sync server after changing `.env`.

The web client reads `VITE_API_URL` from env files in the repository root.
Use `VITE_API_URL=/api` in `.env` for a production build served by the Go client,
and `VITE_API_URL=http://localhost:8081/api` in `.env.development` for local
development. Without a value, the default is `/api` on the current origin.
The URL must point to the Go client API, including `/api`.
Restart Vite after changing env files; production changes require a rebuild.
Only variables prefixed with `VITE_` are exposed to browser code; these values
are public and must not contain secrets.

Start by signing in with your server username and account password, then
unlock the local vault with its master password. Use **Create an account**
for registration; this creates a local vault or registers an existing
unregistered vault using its current master password. A registered vault
can only be unlocked by its owner's account. Each client data directory
contains one vault; this does not add vault migration between devices.

Automatic vault locking keeps the account session, so only the master
password is requested again. **Logout** ends both the account session and
its vault sessions, then returns to the account login screen. Account
sessions expire after 12 hours and are kept in the current browser tab;
restarting the Go client requires signing in again.

The web production build is written to `static/` in the repository root.
The build scripts copy these files to `cmd/client/static/` for embedding in the
Go client binary. `static/` is generated and ignored by Git.

`script/build.sh` builds the web client and both Go binaries. The Go builder
(`go run ./cmd/builder build`) provides the same build workflow.

The server stores data in SQLite (`data/myst.sqlite` by default). Configure it
with `go run ./cmd/server --db data/myst.sqlite --port 8080`.
See [server/README.md](server/README.md) for configuration and persistence details.

Integration tests (`go test ./test/...` or `script/test`) use isolated temporary
SQLite databases and do not require external services.

The web workspace uses the shadcn-svelte dashboard-01 shell and sidebar-07
navigation. Collapse the sidebar with its header button (or Ctrl/Cmd+B); on
mobile it opens as a drawer. PaneForge lets you resize the entry list and
preview, saving pane sizes locally. The header theme button uses mode-watcher
to persist light/dark preference; the initial theme follows the OS.
