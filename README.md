# Grextor

Grextor is a proof-of-concept engine that combines vector similarity search
with graph traversal to retrieve semantically relevant and structurally valid
information.

## Why Grextor?

Vector search answers:
- "What is similar?"

Graph search answers:
- "What is connected, allowed, or dependent?"

Grextor answers both.

## Core Idea

1. Embed documents into vectors
2. Store relationships in a graph
3. Query by meaning
4. Constrain by structure

## Use Cases

- Knowledge base search with access control
- Code intelligence with dependency awareness
- Recommendation systems with business rules
- Trust-based discovery systems

## Status

Experimental PoC. Not production-ready.

## Development

### Prerequisites
- Go 1.22+
- Docker & Docker Compose (for dependencies)

### Quick Start

```bash
# Clone the repo
git clone https://github.com/bondzai/grextor.git
cd grextor

# Run tests
make test

# Build the project (creates grextor-ingest and grextor-query)
make build

# Run the application (example)
./grextor-query
```

### Make Commands
- `make test`: Run unit tests
- `make test-cover`: Run tests with coverage report
- `make build`: Compile the binaries (`grextor-ingest` & `grextor-query`)
- `make run`: Show run instructions
- `make clean`: Remove build artifacts and coverage files
- `make fmt`: Format code

