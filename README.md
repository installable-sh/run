# RUN

Script runner for [installable.sh](https://github.com/installable-sh).

## Installation

### Docker (recommended)

```bash
docker pull ghcr.io/installable-sh/run:v1
```

### From source

```bash
go install github.com/installable-sh/run@latest
```

## Usage

```bash
RUN [+env] [+raw] [+nocache] <url> [args...]
```

### Options

- `+env` - Send environment variables as X-Env-* headers
- `+raw` - Print the script without executing
- `+nocache` - Bypass CDN caches

### Examples

```bash
# Run a script
RUN https://example.com/install.sh

# Run with arguments
RUN https://example.com/install.sh --version 1.0.0

# Preview script without executing
RUN +raw https://example.com/install.sh

# Run with environment forwarding
RUN +env https://example.com/install.sh
```

## Development

```bash
# Build
make build

# Run tests
make test

# Run all CI checks
make ci
```

## License

Apache 2.0 - see [LICENSE](LICENSE)
