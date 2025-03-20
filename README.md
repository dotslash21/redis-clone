# Redis Clone

A Redis server implementation in Go, supporting the Redis RESP (Redis Serialization Protocol).

## Features

Currently implemented:
- TCP server listening on port 6379 (default Redis port)
- RESP protocol implementation (both array format and inline commands)
- Thread-safe key-value store with TTL support
- Sharded in-memory data structure for improved concurrent performance
- Support for the following commands:
  - PING - Returns PONG
  - ECHO - Returns the message
  - SET - Sets a key to a value with optional expiry (via EX and PX)
  - GET - Gets the value of a key

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Git

### Installation

1. Clone the repository
   ```
   git clone https://github.com/dotslash21/redis-clone.git
   cd redis-clone
   ```

2. Build the application
   ```
   go build -o redis-clone ./app
   ```

3. Run the server
   ```
   ./redis-clone
   ```

## Usage

Once the server is running, you can connect to it using the Redis CLI or any Redis client:

```
redis-cli -p 6379
```

### Supported Commands

#### PING
Returns PONG
```
127.0.0.1:6379> PING
PONG
```

#### ECHO
Returns the message sent
```
127.0.0.1:6379> ECHO "Hello, World!"
"Hello, World!"
```

#### SET
Sets a key to a value with optional expiry
```
# Basic SET
127.0.0.1:6379> SET mykey "myvalue"
OK

# SET with expiry in seconds
127.0.0.1:6379> SET mykey "myvalue" EX 10
OK

# SET with expiry in milliseconds
127.0.0.1:6379> SET mykey "myvalue" PX 10000
OK
```

#### GET
Gets the value of a key
```
127.0.0.1:6379> GET mykey
"myvalue"
```

## Project Structure

- `app/` - Application code
  - `main.go` - Entry point of the application
  - `command/` - Implementation of Redis commands
  - `errors/` - Custom error types and handling
  - `resp/` - Redis Serialization Protocol formatting
  - `server/` - TCP server implementation
  - `store/` - In-memory key-value store with TTL support
  - `types/` - Shared data structures (ThreadSafeMap)
- `tests/` - Integration tests
  - `commands_test.go` - End-to-end command tests
  - `helpers/` - Test utilities including a Redis client

## Implementation Details

The server implements the Redis RESP (Redis Serialization Protocol) for communication:
- Simple string responses are prefixed with `+` (e.g., `+OK\r\n` for SET command)
- Error responses are prefixed with `-` (e.g., `-ERR message\r\n`)
- Integer responses are prefixed with `:` (e.g., `:1000\r\n`)
- Bulk string responses are prefixed with `$` followed by the string length (e.g., `$11\r\nHello,Redis!\r\n` for ECHO command)
- Null bulk strings are represented as `$-1\r\n` (returned for GET on a non-existent key)

### RESP Protocol Examples

- **PING**: Client sends `*1\r\n$4\r\nPING\r\n` and receives `+PONG\r\n`
- **ECHO**: Client sends `*2\r\n$4\r\nECHO\r\n$11\r\nHello,Redis!\r\n` and receives `$11\r\nHello,Redis!\r\n`
- **SET**: Client sends `*3\r\n$3\r\nSET\r\n$9\r\nmykey-resp\r\n$12\r\nmyvalue-resp\r\n` and receives `+OK\r\n`
- **GET**: Client sends `*2\r\n$3\r\nGET\r\n$9\r\nmykey-resp\r\n` and receives `$12\r\nmyvalue-resp\r\n`
- **GET (non-existent key)**: Client sends `*2\r\n$3\r\nGET\r\n$14\r\nnonexistentkey\r\n` and receives `$-1\r\n`

The server supports both RESP array format and inline commands for client communication.

### Key Features of the Implementation

- **Thread-safe store**: Uses a sharded map implementation for better concurrency
- **TTL support**: Keys can expire after a specified time (seconds or milliseconds)
- **Graceful shutdown**: Handles termination signals properly
- **Custom error handling**: Structured error types with context information
- **Comprehensive tests**: Unit and integration tests for all components

## Future Enhancements

- Support for more Redis commands
- Data persistence
- Replication
- Cluster mode

## License

This project is licensed under the MIT License - see the LICENSE file for details.
