# NMEA Simulator

A versatile marine data simulator supporting both NMEA 0183 and NMEA 2000 protocols simultaneously.

## Features

### NMEA 0183
- **TCP Server** (default port 10110)
- **WebSocket Server** with web interface (default port 8080)
- **Configurable Baud Rates**: 4800, 9600, 19200, 38400
- **Supported Sentences**
  - Position: GGA (GPS Fix), GLL (Geographic Position)
  - Navigation: RMC (Recommended Minimum), HDT (True Heading), VTG (Track & Speed), XTE (Cross-Track Error)
  - Environment: DBT (Depth Below Transducer), MTW (Water Temperature), MWV (Wind), VHW (Water Speed & Heading), DPT (Depth)

### NMEA 2000
- **TCP Server** (default port 10200)
- **WebSocket Server** with web interface (default port 8081)
- **Supported PGNs**
  - 127250 (Vessel Heading)
  - 128259 (Speed)
  - 128267 (Water Depth)
  - 129025 (Position Rapid Update)
  - 129026 (COG & SOG Rapid Update)
  - 130306 (Wind Data)

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

By default, the simulator runs both NMEA 0183 and NMEA 2000 protocols simultaneously:
```bash
nmeasim
```

### Run Specific Protocol

Run only NMEA 0183:
```bash
nmeasim --protocol nmea0183
```

Run only NMEA 2000:
```bash
nmeasim --protocol nmea2000
```

### Command Line Options

Protocol Selection:
- `--protocol`: Protocol to use ("both", "nmea0183", or "nmea2000", default: "both")

NMEA 0183 Options:
- `--nmea0183-ws-port`: WebSocket server port (default: 8080)
- `--nmea0183-tcp-port`: TCP server port (default: 10110)
- `--baud`: Baud rate for TCP output (default: 4800)

NMEA 2000 Options:
- `--nmea2000-ws-port`: WebSocket server port (default: 8081)
- `--nmea2000-tcp-port`: TCP port (default: 10200)

Common Options:
- `--host`: Host to bind servers to (default: "0.0.0.0")
- `--interval`: Data update interval (default: 1s)

## Web Interface

The web interface now supports viewing both NMEA 0183 and NMEA 2000 data simultaneously. Access it at:
- NMEA 0183: http://localhost:8080
- NMEA 2000: http://localhost:8081

## Viewing TCP Data

You can use common terminal commands to view the NMEA data streams directly:

### Using netcat (nc)

For NMEA 0183:
```bash
nc localhost 10110
```

For NMEA 2000:
```bash
nc localhost 10200
```

### Using telnet

For NMEA 0183:
```bash
telnet localhost 10110
```

For NMEA 2000:
```bash
telnet localhost 10200
```

### Using socat with hex dump

To view data with timestamps and hex dump:

For NMEA 0183:
```bash
socat TCP:localhost:10110 STDOUT | hexdump -C
```

For NMEA 2000:
```bash
socat TCP:localhost:10200 STDOUT | hexdump -C
```

Note: You may need to install these tools first:
- macOS: `brew install netcat socat`
- Ubuntu/Debian: `sudo apt install netcat-openbsd socat`
- Windows: Use PowerShell's `Test-NetConnection` or install WSL

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