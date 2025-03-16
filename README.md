# Redis Clone

A Redis server implementation in Go, supporting the Redis RESP (Redis Serialization Protocol).

## Features

Currently implemented:
- TCP server listening on port 6379 (default Redis port)
- Basic RESP protocol implementation
- Key-value store with TTL support
- Support for the following commands:
  - PING
  - ECHO
  - SET (with expiry via EX and PX)
  - GET

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Git

### Installation

1. Clone the repository
   ```
   git clone https://github.com/username/redis-clone.git
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

Currently, the server supports the following Redis commands:

- `PING`: Returns PONG
- `ECHO`: Returns the message
- `SET`: Sets a key to a value with optional expiry
- `GET`: Gets the value of a key

Example:
```
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> ECHO "Hello, World!"
"Hello, World!"
127.0.0.1:6379> SET mykey "myvalue"
OK
127.0.0.1:6379> GET mykey
"myvalue"
```

## Project Structure

- `app/main.go`: Entry point of the application
- `app/server/server.go`: TCP server implementation with Redis protocol handling
- `app/commands/commands.go`: Implementation of supported Redis commands
- `app/store/store.go`: In-memory key-value store with TTL support

## Implementation Details

The server implements the Redis RESP (Redis Serialization Protocol) for communication:
- Simple string responses are prefixed with '+'
- Error responses are prefixed with '-'

## Future Enhancements

- Support for more Redis commands
- Data persistence
- Replication
- Cluster mode

## License

This project is licensed under the MIT License - see the LICENSE file for details.
