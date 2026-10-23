# PortPassport

PortPassport is a command-line tool that looks up which process is using a given port. You pass a port number and it prints the protocol, bind address, process ID, process name, and user.

Lookups go through a `PortInspector` interface. The current implementation is a mock that returns sample data for any port you ask about, so the CLI shape is in place while a real system inspector is still to be added.

## Requirements

- Go 1.27.1 or later

## Build and run

From the repository root:

```bash
go build -o PortPassport .
./PortPassport 8080
```

Or run it without building a binary:

```bash
go run . 8080
```

## Usage

```text
PortPassport <port>
```

| Argument | Description |
| --- | --- |
| `port` | Port number to look up (required integer) |

Example:

```bash
./PortPassport 8080
```

```text
PortPassport!
Port: 8080
Protocol: TCP
Address: 127.0.0.1:8080
PID: 4218
Process: dotnet
User: Aryan
```

If you omit the port, the program prints usage and exits. If the argument is not an integer, it prints `Invalid port`.

## Project layout

```text
.
├── main.go                 # CLI entry point: parse the port and print results
├── go.mod
└── internal/port/
    ├── port.go             # PortInfo and the PortInspector interface
    └── mock.go             # Mock inspector that returns sample data
```

## How it works

`main` reads the first argument, parses it as an integer, and calls `Find` on a `PortInspector`.

```go
type PortInfo struct {
    Port     int
    Protocol string
    Address  string
    PId      int
    Process  string
    User     string
}

type PortInspector interface {
    Find(port int) ([]PortInfo, error)
}
```

`MockPortInspector` implements that interface and always returns one sample `PortInfo` for the requested port. A real inspector (for example one that reads the operating system's socket table) can replace the mock without changing the CLI, as long as it implements `Find`.
