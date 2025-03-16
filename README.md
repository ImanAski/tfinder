# tfinder

A command-line tool to find and track TODO, FIXME, BUG, and other actionable comments in your codebase.

## Overview

`tfinder` scans your codebase for pattern-based comments like TODOs and FIXMEs, making it easy to track and manage technical debt across your projects. It automatically detects comments in various programming languages and presents them in a colorized, readable format.

## Features

- 🔍 Find TODO, FIXME, BUG, and other custom pattern comments
- 🎨 Colorized output for easy visual scanning
- ⚙️ Configurable via JSON config file
- 🚫 Ignore specific directories or files
- 💾 Works in any codebase with minimal setup

## Installation

### From Releases

Download the latest binary for your platform from the [releases page](https://github.com/yourusername/tfinder/releases).

### From Source

```bash
git clone https://github.com/yourusername/tfinder.git
cd tfinder
go build
```

## Usage

Basic usage:

```bash
tfinder
```

This will scan the current directory for comments matching default patterns.

### Command-line Options

```
-dir string
    The path to look for comments (default ".")
-ignore string
    Comma-separated list of patterns to ignore (default ".git,node_modules,vendor")
-pattern string
    Comma-separated list of patterns to search for (default "TODO,FIXME,DO,BUG")
-only-comments
    Only search in comments (default true)
-create-config
    Create default config file
```

Example:

```bash
tfinder -dir ./src -pattern "TODO,CRITICAL" -ignore ".git,dist,build"
```

## Configuration

You can create a default configuration file by running:

```bash
tfinder -create-config
```

This will create a `tfinder.json` file with default settings:

```json
{
  "dir": ".",
  "ignore": [".git", "node_modules", "vendor"],
  "pattern": ["TODO", "FIXME", "DO", "BUG"],
  "colors": {
    "TODO": "yellow:bg-black",
    "FIXME": "red:bg-white",
    "DO": "green",
    "BUG": "red:bg-yellow"
  }
}
```

### Color Configuration

Colors are specified using the format `foreground:bg-background`. Available colors:

**Foreground Colors:**
- red
- green
- yellow
- blue
- purple
- cyan

**Background Colors:**
- bg-red
- bg-green
- bg-yellow
- bg-blue
- bg-purple
- bg-cyan
- bg-white
- bg-black

## Supported Comment Formats

tfinder detects comments in various formats:

- `// TODO: Fix this` (C, C++, Go, Java, JavaScript, etc.)
- `# TODO: Fix this` (Python, Ruby, Bash, etc.)
- `-- TODO: Fix this` (SQL, Lua, etc.)
- `% TODO: Fix this` (Matlab, LaTeX, etc.)
- `; TODO: Fix this` (Assembly, Lisp, etc.)
- `/* TODO: Fix this */` (C, C++, Java, etc.)
- `<!-- TODO: Fix this -->` (HTML, XML, etc.)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
