# go-error-simulator

This simple app is an HTTP server which lets you control the response and behavior in several ways.
Options are:

* returned HTTP status code
* print text to STDOUT
* print text to STDERR
* add a (random) delay before returning the response

## Usage

**Run locally**

```
go run main.go
```

Then you can access the server at `http://localhost:8080`.

**Run via Docker**

```
docker run -p 8080:8080 -it --rm ghcr.io/aeimer/go-error-simulator:latest
```

Then you can access the server at `http://localhost:8080`.

## API

### `GET /`

Returns a simple HTML page with a link to the API documentation.

### `GET /simulate`

Simulates an error based on the query parameters. The following parameters are supported:
- `status`: HTTP status code to return (default: 200)
- `stdout`: text to print to STDOUT (default: empty)
- `stderr`: text to print to STDERR (default: empty)
- `latency`: delay in milliseconds before returning the response (default: 0)
  - exact number: `latency=100`
  - random range: `latency=100-200`
  - random range shortcut: `latency=-100` (random between 0 and 100)
