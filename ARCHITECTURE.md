# Testing Studio - Architecture

## Overview

Testing Studio (formerly CloudEvents Explorer) is a refactored, maintainable Go web application for exploring and testing Kafka and Google PubSub messages. The codebase has been restructured from a monolithic `main.go` into a clean, modular architecture following Go best practices.

## Project Structure

```
testing-studio/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go            # Load/save configs, manage PubSub/Kafka settings
│   ├── handlers/                # HTTP request handlers
│   │   ├── index.go             # Landing page handler
│   │   ├── pubsub.go            # PubSub page handler
│   │   ├── kafka.go             # Kafka page handler
│   │   ├── flowdiagram.go       # Flow diagram page handler
│   │   └── api.go               # API endpoints (pull, publish, config)
│   ├── kafka/                   # Kafka operations
│   │   └── kafka.go             # Pull messages, publish messages, Avro encoding
│   ├── pubsub/                  # PubSub operations
│   │   └── pubsub.go            # Pull messages from Google PubSub
│   ├── templates/               # HTML templates and components
│   │   ├── index.go             # Landing page template
│   │   ├── base.go              # Base HTML template with shared styles/JS
│   │   ├── components.go        # Reusable components (Base64 modal)
│   │   ├── pubsub.go            # PubSub page content and JavaScript
│   │   ├── kafka.go             # Kafka page content and JavaScript
│   │   └── flowdiagram.go       # COMMS EPIC flow diagram
│   └── types/                   # Shared data structures
│       └── cloudevent.go        # CloudEvent type definition
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── start.sh                     # Quick launcher script
├── README.md                    # User guide
├── ARCHITECTURE.md              # This file
└── configs.json                 # Saved configurations (auto-generated)
```

## Design Principles

### 1. Separation of Concerns

Each package has a single, well-defined responsibility:

- **config**: Manages configuration persistence and retrieval
- **handlers**: HTTP request routing and response handling
- **kafka**: Kafka-specific business logic (consume, produce, Avro encoding)
- **pubsub**: Google PubSub operations
- **templates**: UI rendering and HTML generation
- **types**: Shared data structures across packages

### 2. Dependency Flow

```
cmd/server/main.go
    ├─> handlers (HTTP routing)
    │       ├─> config (load/save settings)
    │       ├─> kafka (message operations)
    │       ├─> pubsub (message operations)
    │       └─> templates (UI rendering)
    │
    └─> config (initialize on startup)
```

### 3. Package Naming

- `cmd/server`: Executable entry point
- `internal/*`: Private packages, not importable by external projects
- Clear, descriptive names: `pubsub`, `kafka`, `handlers`, `templates`

## Key Components

### Configuration (`internal/config`)

**Purpose**: Centralized configuration management with thread-safe access

**Key Functions**:
- `Load()`: Load configs from disk, create defaults if missing
- `Save()`: Persist configs to disk
- `Get()`: Thread-safe config retrieval
- `AddOrUpdatePubSubConfig()`: Add or update PubSub configuration
- `AddOrUpdateKafkaConfig()`: Add or update Kafka configuration

**Thread Safety**: Uses `sync.RWMutex` for concurrent access

### Handlers (`internal/handlers`)

**Purpose**: HTTP request handling and routing

**Files**:
- `index.go`: Landing page with platform selection
- `pubsub.go`: PubSub message viewer page
- `kafka.go`: Kafka message viewer and publisher page
- `flowdiagram.go`: COMMS EPIC eventing flow diagram
- `api.go`: REST API endpoints for pulling/publishing messages

**Pattern**: Each handler is a simple function that delegates to business logic packages

### Kafka Operations (`internal/kafka`)

**Purpose**: Kafka-specific operations (consume, produce, Avro serialization)

**Key Functions**:
- `Pull()`: Consume messages from Kafka topics with Avro decoding
- `Publish()`: Produce messages to Kafka topics with Avro encoding
- `decodeAvroMessage()`: Decode Avro binary using schema registry

**Features**:
- Automatic schema registry lookup
- Confluent wire format support (magic byte + schema ID)
- Fallback to raw/JSON if Avro decoding fails

### PubSub Operations (`internal/pubsub`)

**Purpose**: Google Cloud PubSub operations

**Key Functions**:
- `Pull()`: Pull messages from PubSub subscriptions
- CloudEvent attribute parsing (`ce-type`, `ce-subject`, etc.)

**Features**:
- Emulator support via `PUBSUB_EMULATOR_HOST`
- Automatic message acknowledgment
- Timeout-based pull with configurable max messages

### Templates (`internal/templates`)

**Purpose**: HTML/JavaScript UI rendering

**Architecture**:
- **Base Template** (`base.go`): Shared layout, styles, message rendering logic
- **Page Templates**: Specific content and JavaScript for each page
- **Components**: Reusable UI elements (Base64 modal, etc.)

**Key Functions**:
- `GetBaseHTML(title, content, extraJS)`: Generate full HTML page
- Constants for each page's content and JavaScript

### Types (`internal/types`)

**Purpose**: Shared data structures

**CloudEvent**: Unified message format for both PubSub and Kafka messages

## How to Modify

### Adding a New Message Platform

1. Create `internal/newplatform/newplatform.go`:
   ```go
   package newplatform

   import "cloudevents-explorer/internal/types"

   type PullParams struct { /* ... */ }
   type PullResult struct { /* ... */ }

   func Pull(params PullParams) (*PullResult, error) {
       // Implementation
   }
   ```

2. Create `internal/templates/newplatform.go`:
   ```go
   package templates

   const NewPlatformContent = `<div>...</div>`
   const NewPlatformJS = `function pullMessages() { /*...*/ }`
   ```

3. Create `internal/handlers/newplatform.go`:
   ```go
   package handlers

   func HandleNewPlatform(w http.ResponseWriter, r *http.Request) {
       html := templates.GetBaseHTML("New Platform", templates.NewPlatformContent, templates.NewPlatformJS)
       w.Header().Set("Content-Type", "text/html; charset=utf-8")
       fmt.Fprint(w, html)
   }
   ```

4. Add route in `cmd/server/main.go`:
   ```go
   http.HandleFunc("/newplatform", handlers.HandleNewPlatform)
   http.HandleFunc("/api/newplatform/pull", handlers.HandlePullNewPlatform)
   ```

### Adding a New Feature to Existing Platform

1. Add business logic to `internal/kafka/` or `internal/pubsub/`
2. Update templates in `internal/templates/`
3. Add handler in `internal/handlers/api.go`
4. Register route in `cmd/server/main.go`

### Modifying UI

1. Find template in `internal/templates/`
2. Update HTML in content constants
3. Update JavaScript in JS constants
4. Rebuild: `go build cmd/server/main.go`

## Testing

```bash
# Run from project root
go test ./...

# Test specific package
go test ./internal/kafka/
go test ./internal/pubsub/
go test ./internal/config/
```

## Building

```bash
# Development mode (with hot reload if using air/fresh)
go run cmd/server/main.go

# Production build
go build -o testing-studio cmd/server/main.go
./testing-studio

# Quick start script
./start.sh
```

## Future Improvements

1. **Template System**: Consider using `html/template` for better separation
2. **Dependency Injection**: Use interfaces for better testability
3. **Middleware**: Add logging, metrics, error handling middleware
4. **Configuration**: Support environment variables and config files
5. **Observability**: Add structured logging, metrics, tracing
6. **Testing**: Add unit tests for each package
7. **API Documentation**: OpenAPI/Swagger specs for REST endpoints

## Migration from Old Structure

The refactoring maintained 100% feature parity:

- ✅ All original features preserved
- ✅ No breaking changes to APIs or UI
- ✅ Same port (8888), same endpoints, same behavior
- ✅ Configurations remain compatible

**Old**: Single `main.go` (2183 lines)
**New**: Modular structure (~300 lines per package)

**Benefits**:
- Easier to understand and modify
- Better code organization
- Improved testability
- Clearer separation of concerns
- Maintainable for team collaboration
