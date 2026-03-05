# Mirra

A lightweight file server with a modern web UI. Built with Go using server-side rendering.

<div align="center">

**[English](README.md)** | **[中文](README_zh.md)**

</div>

## Features

- Single binary deployment, no external dependencies
- Web UI with light/dark theme support
- Directory browsing with real-time search
- README.md rendering with syntax highlighting

## Configuration

### Automatic Configuration
If `config.json` doesn't exist, it will be automatically created with the current directory as the root path.

### Manual Configuration
Copy `config.example.json` to `config.json` and edit:

```bash
cp config.example.json config.json
```

Edit `config.json`:

```json
{
  "server": {
    "name": "Mirra",
    "host": "0.0.0.0",
    "port": "8080",
    "favicon": ""
  },
  "share": {
    "root_path": "."
  },
  "appearance": {
    "theme": "auto",
    "show_hidden": false
  }
}
```

- `server.name`: Display name in web UI
- `server.host`: Listen host (use `0.0.0.0` for all interfaces)
- `server.port`: Listen port
- `share.root_path`: Root directory to serve
- `appearance.theme`: Default theme (`light`, `dark`, or `auto` for system preference)
- `appearance.show_hidden`: Whether to show hidden files (starting with `.`)

## Usage

```bash
# Run with default config file (config.json)
./mirra

# Run with custom config file path
./mirra -c /path/to/config.json

# Show version information
./mirra -v
```

## Building from Source

### Prerequisites

- Go 1.25 or later

```bash
# Build Go binary
go build -o mirra ./cmd/server

# Run
./mirra
```

### Using Makefile

```bash
# Build binary for current platform
make build

# Cross-compile for all platforms (linux, darwin, windows)
make cross-build

# Create distribution packages
make dist

# Run server
make run

# Format code
make fmt

# Lint check
make lint

# Clean build artifacts
make clean
```

## Changelog

### v0.0.2 (2026-03-05)

**Added**
- Added `-c` flag to specify custom config file path
- Added Prism.js syntax highlighting support for 30+ programming languages

**Changed**
- Changed all code comments and user-facing messages to English
- Optimized page title display logic (removed redundant server name element)
- Improved theme switching to preserve Prism.js highlighting effect

For earlier versions, see [CHANGELOG.md](CHANGELOG.md).

## TODO

### Completed
- [x] File server with directory browsing
- [x] Light/dark theme support
- [x] README.md rendering
- [x] Real-time search
- [x] Syntax highlighting for code blocks

### Planned
- [ ] Mobile WebUI improvements
- [ ] Code preview feature
- [ ] File URL copy and folder download
- [ ] Multi-threaded download support

## License

MIT License - see [LICENSE](LICENSE) for details.
