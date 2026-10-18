# Go client

```text
apis/             Local HTTP handlers, DTO conversion and OpenAPI
  generated/      Generated local API types and Go client
core/             Use cases, sessions, sync logic, adapter interfaces and errors
models/           Entry, Keystore, Credentials, User and Invitation entities
storage/          Encrypted file repository and persistence DTO conversion
  enclave/        Low-level encrypted enclave format and integrity checks
remote/           HTTP adapter for the sync server API
scheduler/        Periodic sync and remote reauthentication
```

`cmd/client` assembles `core`, `storage`, `remote`, `apis` and `scheduler`.
The `core` package depends on `models` and declares the `Enclave` and `Remote`
interfaces. Both adapters have compile-time checks against these interfaces.
`remote` uses the server's generated API client from `server/apis/generated`.

Client entities are separate from server entities. HTTP DTOs and persisted JSON
are converted explicitly; encrypted files and HTTP contracts retain their
existing formats. Local storage reports errors from the client's `core`.

Run from the repository root:

```sh
go test ./client/golang/... ./test/...
go vet ./client/golang/... ./cmd/client
```

The local API's OpenAPI schema is `apis/openapi.json`. Its generation directives
write TypeScript API code to `client/web/src/api/generated`.
