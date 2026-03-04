# File Server

A simple file server with web UI, similar to university mirror sites. Built with Go using server-side rendering.

## Features

- Single binary deployment
- Configurable listening address and root directory
- Web UI with GitHub-style light/dark themes
- Theme toggle button (follows system preference by default)
- Server-side rendered directory listing
- README.md rendering (like GitHub)
- Statistics (directories, files, total size)

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
    "name": "File Server",
    "host": "0.0.0.0",
    "port": "8080"
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

- Go 1.22 or later

### Using Makefile

```bash
# Build binary
make build

# Run server
make run

# Format code
make fmt

# Run tests
make test

# Lint check
make lint
```

The binary `fileserver` will be created in the project root.

### Manual Build

```bash
# Build Go binary
go build -o fileserver main.go

# Run
./fileserver
```

## Architecture

The server uses Go's `html/template` for server-side rendering. All HTML, CSS, and JavaScript are embedded in the binary via `go:embed`.

- `main.go` – Main server logic with directory handling and Markdown rendering
- `template.html` – HTML template with embedded CSS and JavaScript
- `config.json` – Configuration file

### Theme System

The UI supports light and dark themes based on CSS custom properties. The theme defaults to system preference, but can be overridden by:
1. User toggle (saved in localStorage)
2. Server configuration (`appearance.theme` in `config.json`)

Supported theme values: `light`, `dark`, or `auto` (follows system preference).

## Offline Usage (No Internet Required)

This server is designed to work completely offline in internal network environments:

1. **FontAwesome icons are served locally** – No external CDN dependencies
2. **All web resources are embedded** – CSS, fonts, and scripts are included in the `static/` directory

### Setting Up Local FontAwesome

```bash
# Download FontAwesome resources
./download-fontawesome.sh

# Or manually download from:
# https://fontawesome.com/download
# Place the CSS files in static/css/ and font files in static/fonts/
```

### Static File Structure
```
static/
├── css/
│   └── font-awesome.min.css    # FontAwesome CSS
└── fonts/
    ├── fa-solid-900.woff2      # Solid icons font
    ├── fa-regular-400.woff2    # Regular icons font
    └── fa-brands-400.woff2     # Brand icons font
```

## Building and Distribution

### Makefile Commands
```bash
make build          # Build for current platform (output: dist/fileserver)
make cross-build    # Cross-compile for all platforms
make dist           # Create distribution packages (.tar.gz/.zip)
make clean          # Clean build artifacts
make help           # Show all available commands
```

### Supported Platforms
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)
- FreeBSD (amd64)

### Distribution Packages
Distribution packages include:
- The compiled binary
- `config.example.json`
- `README.md`

### GitHub Actions
CI/CD pipeline is configured in `.github/workflows/build.yml`:
- Tests on push/PR
- Cross-compilation for all platforms
- Automatic release creation when tags are pushed

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Format code
make fmt

# Run with hot reload (using air or similar)
go run main.go
```

## License

MIT