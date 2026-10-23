# PortPassport

PortPassport is a macOS command-line tool for finding processes that are
listening on TCP ports. It can inspect one port, list all active TCP listeners,
or interactively terminate a process using a selected port.

## Requirements

- macOS
- Go 1.27.1 or later to build from source
- `lsof` available on your `PATH`

## Build

From the repository root:

```bash
go build -o PortPassport .
```

You can also run it directly with Go:

```bash
go run . <port>
```

## Usage

```text
PortPassport <port>
PortPassport --list
PortPassport --kill <port>
```

| Command | Description |
| --- | --- |
| `PortPassport <port>` | Show TCP listeners on the specified port. |
| `PortPassport --list` | Show all active TCP listeners. |
| `PortPassport --kill <port>` | Show listeners on the specified port, then prompt you to select and confirm a process to terminate. |

Ports must be integers between 1 and 65535.

Examples:

```bash
./PortPassport 8080
./PortPassport --list
./PortPassport --kill 8080
```

Port inspection output includes the port, process ID, process name, user,
protocol, and local address. If no listener is found, the command reports that
the port has no active TCP listeners.

### Terminating a process

`--kill` asks for the PID to terminate and requires you to type `yes` to
confirm. It rechecks that the selected process is still listed on the port
before sending the termination signal. This action forcibly terminates the
process; only confirm it when you are sure it is safe to stop.

## How it works

The macOS inspector uses `lsof` to find TCP sockets in the `LISTEN` state.
`PortInspector` exposes `Find` for a single port and `List` for all active
listeners. The CLI uses the inspector to display results and, in `--kill` mode,
to recheck the selected process before termination.

## Project layout

```text
.
├── main.go
├── go.mod
└── internal/port/
    ├── darwin.go
    ├── internal_darwin.go
    ├── port.go
    ├── process_darwin.go
    └── validation.go
```
