# procfile-run

Minimal Procfile runner for local development with colored output, process supervision, and port conflict detection.

## Installation

```bash
go install github.com/yourusername/procfile-run@latest
```

Or build from source:

```bash
go build -o procfile-run .
```

## Usage

Create a `Procfile` in your project root:

```
web:    go run main.go
worker: go run worker/main.go
redis:  redis-server --port 6379
```

Then run:

```bash
procfile-run
```

### Options

```
-f, --file      Path to Procfile (default: ./Procfile)
-p, --port      Base port for $PORT assignment (default: 5000)
    --no-color  Disable colored output
```

Each process gets a unique color in the terminal output, is automatically restarted on failure, and `procfile-run` will warn you if any required ports are already in use before starting.

## Example Output

```
12:01:00 web    | Listening on port 5000
12:01:00 worker | Worker started, waiting for jobs...
12:01:01 redis  | Ready to accept connections
```

## Requirements

- Go 1.21+
- A valid `Procfile` in your project directory

## License

MIT © [yourusername](https://github.com/yourusername)