# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.2] - 2026-03-05

### Added
- Added `-c` flag to specify custom config file path
- Added Prism.js syntax highlighting support for 30+ programming languages

### Changed
- Changed all code comments and user-facing messages to English
- Optimized page title display logic (removed redundant server name element)
- Improved theme switching to preserve Prism.js highlighting effect

## [0.0.1] - 2026-03-04
### Added
- Initial stable release
- Single binary deployment, no additional dependencies required
- GitHub-style light/dark theme support with system preference detection
- Directory listing with file type icons
- README.md rendering (GitHub-like)
- Breadcrumb navigation
- Statistics display (directory count, file count, total size)
- Sortable file list (name, size, modified time)
- Real-time search filtering
- SPA-style navigation (no-refresh directory switching)
- Code block syntax highlighting (Prism.js localized)
- Responsive design with mobile support

### Changed
- Refactored project structure to Go standard layout (cmd/, internal/)
- Separated CSS and JS from template.html into standalone files
- All static resources localized, no CDN dependencies
- Optimized SPA navigation logic to support all directory links
- Optimized README area dynamic rendering logic
