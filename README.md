# NMEA Simulator

A versatile NMEA 0183 sentence simulator that generates realistic marine navigation and environmental data over TCP and WebSocket protocols.

## Features

- **Multiple Protocol Support**
  - TCP server (default NMEA port 10110)
  - WebSocket server with web interface
- **Configurable Baud Rates**: 4800, 9600, 19200, 38400
- **NMEA 0183 Sentences**
  - Position: GGA (GPS Fix), GLL (Geographic Position)
  - Navigation: RMC (Recommended Minimum), HDT (True Heading), VTG (Track & Speed), XTE (Cross-Track Error)
  - Environment: DBT (Depth Below Transducer), MTW (Water Temperature), MWV (Wind), VHW (Water Speed & Heading), DPT (Depth)

## Installation

```bash
go install github.com/captv89/nmea-simulator/cmd/nmeasim@latest
```

Or clone and build:

```bash
git clone https://github.com/captv89/nmea-simulator.git
cd nmea-simulator
make build
```

## Usage

Run with default settings:
```bash
nmeasim
```

### Command Line Options

- `--ws-port`: WebSocket server port (default: 8080)
- `--tcp-port`: TCP server port (default: 10110)
- `--host`: Host to bind servers to (default: "0.0.0.0")
- `--interval`: NMEA sentence update interval (default: 1s)
- `--baud`: Baud rate for TCP output (default: 4800)

Example with custom settings:
```bash
nmeasim --ws-port 8081 --tcp-port 10111 --interval 500ms --baud 9600
```

## Development

### Prerequisites

- Go 1.24 or higher
- Make (optional, for using Makefile commands)

### Available Make Commands

- `make build`: Build the binary
- `make test`: Run tests
- `make coverage`: Generate test coverage report
- `make clean`: Clean build files
- `make release`: Create a new release
- `make release-dry-run`: Test release process without publishing

### Running Tests

```bash
make test
```

Generate coverage report:
```bash
make coverage
```

## License

This project is licensed under the [LICENSE](LICENSE) file in the repository.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request