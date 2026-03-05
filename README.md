# Mirra

A simple file server with web UI, similar to university mirror sites. Built with Go using server-side rendering.

## Features

- Single binary deployment
- Configurable listening address and root directory
- Web UI with GitHub-style light/dark themes
- Theme toggle button (follows system preference by default)
- Server-side rendered directory listing
- README.md rendering (like GitHub)
- Statistics (directories, files, total size)
- Breadcrumb navigation
- File list with sorting (name, size, modified time)
- Real-time search filtering
- SPA-style navigation (no-refresh directory switching)
- Code block syntax highlighting (Prism.js)
- Responsive design with mobile support

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

## Building and Running

### Prerequisites

- Go 1.25 or later

### Usage

```bash
# Run with default config file (config.json)
./mirra

# Run with custom config file path
./mirra -c /path/to/config.json

# Show version information
./mirra -v
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

### Manual Build

```bash
# Build Go binary
go build -o mirra ./cmd/server

# Run
./mirra
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       ├── main.go           # Main entry point
│       └── static/
│           ├── template.html # HTML template with embedded CSS/JS
│           ├── css/          # Stylesheets
│           ├── js/           # JavaScript files
│           └── webfonts/     # FontAwesome fonts
├── internal/
│   ├── config/               # Configuration handling
│   ├── handlers/             # HTTP handlers
│   ├── types/                # Type definitions
│   ├── utils/                # Utility functions
│   └── version/              # Version information
├── config.json               # Configuration file
├── config.example.json       # Example configuration
├── Makefile                  # Build scripts
└── go.mod                    # Go module definition
```

## Architecture

The server uses Go's `html/template` for server-side rendering. All HTML, CSS, and JavaScript are embedded in the binary via `go:embed`.

### Key Components

- `cmd/server/main.go` - Main server logic with directory handling
- `cmd/server/static/template.html` - HTML template with embedded CSS and JavaScript
- `internal/config/config.go` - Configuration loading and management
- `internal/handlers/handlers.go` - HTTP request handlers
- `internal/utils/utils.go` - Utility functions (size formatting, etc.)

### Theme System

The UI supports light and dark themes based on CSS custom properties. The theme defaults to system preference, but can be overridden by:
1. User toggle (saved in localStorage)
2. Server configuration (`appearance.theme` in `config.json`)

Supported theme values: `light`, `dark`, or `auto` (follows system preference).

## Cross-Compilation

The Makefile supports cross-compilation for multiple platforms:

```bash
# Build for all supported platforms
make cross-build

# Output files will be in dist/ directory:
# - mirra_linux_amd64
# - mirra_linux_arm64
# - mirra_mac_amd64
# - mirra_mac_arm64
# - mirra_windows_amd64.exe
```

Supported platforms:
- Linux (amd64, arm64, 386, arm, riscv64)
- macOS (amd64, arm64)
- Windows (amd64, arm64, 386, arm)

## Release Script

Use `release.sh` to manage version tags:

```bash
# Show current version and create tag
./release.sh

# Bump major version (1.0.0 -> 2.0.0)
./release.sh --major

# Bump minor version (1.0.0 -> 1.1.0)
./release.sh --minor

# Bump patch version (1.0.0 -> 1.0.1)
./release.sh --patch

# Only create tag without bumping version
./release.sh --tag-only

# Delete current version tag (local and remote)
./release.sh --revert
```

After creating a tag, push it to trigger the GitHub Actions release workflow.

## GitHub Actions

The Release workflow is configured in `.github/workflows/build.yml`:

- **Lint**: Code formatting check
- **Cross-compile**: Build for all supported platforms
- **Release**: Create GitHub release with binaries

The workflow triggers automatically when a version tag (e.g., `v0.0.1`) is pushed.

## Development

```bash
# Install dependencies
make deps

# Format code
make fmt

# Run with hot reload (using air or similar)
go run ./cmd/server
```

## License

MIT
