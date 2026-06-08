# totion

Small Go demo app demonstrating a simple user management system with registration, login, sessions, and password reset.

## Run

- Run demo (prints flows):

```bash
go run .
```

- Run interactive CLI:

```bash
go run . -- cli
# or
go run . cli
```

## Commands (interactive CLI)

- `register <username> <email> <password>`
- `login <username> <password>`
- `logout <sessionID>`
- `info <sessionID>`
- `list`
- `delete <username>`
- `changepassword <sessionID> <oldPassword> <newPassword>`
- `requestreset <email>`
- `resetpassword <token> <newPassword>`
- `exit`

## Tests

Run tests with:

```bash
go test ./...
```

## Notes

- This repo is a single-package example. Consider splitting into packages (`auth`, `store`, `cli`) for larger projects.
- Password policy and TTLs are defined in `config.go`.
