# 🏥 Go CLI Health Checker

A production-ready CLI tool for health checking multiple HTTP endpoints concurrently, built with Go. This project demonstrates Go best practices, concurrent programming with goroutines, professional CLI design with Cobra, and Docker multi-stage build optimization.

## 📋 Project Overview

This health checker performs concurrent HTTP GET requests to multiple endpoints and reports their health status, response times, and any errors. It's packaged as both a standalone binary and an optimized Docker container.

## ✨ Features

- ⚡ **Concurrent health checks** using goroutines and WaitGroups
- 🎯 **Professional CLI** with Cobra framework
- ⏱️ **Response time tracking** for performance monitoring
- 🔧 **Configurable timeouts** via command-line flags
- 🌐 **Custom URL support** for checking any HTTP endpoint
- 📊 **Verbose output mode** for detailed diagnostics
- 🐳 **Optimized Docker images** (98% size reduction: 550MB → 15MB)
- 🔒 **Security-focused** (non-root user, static binary)
- 📝 **Automatic help documentation** via Cobra

## 🎓 What Was Built & Learned

This project is a small learning journey through Go development:

### Core Go Concepts
- **Structs & Types**: Creating custom data structures (`Endpoint`, `HealthResult`)
- **Error Handling**: Go's explicit error handling pattern (no try/catch)
- **HTTP Client**: Making requests with timeouts and proper cleanup
- **Concurrency**: Goroutines for parallel execution
- **Synchronization**: WaitGroups for coordinating concurrent operations
- **defer**: Proper resource cleanup
- **Slices**: Dynamic arrays for endpoint collections

### CLI Development
- **Cobra Framework**: Professional command-line interface
- **Command Structure**: Root commands and subcommands
- **Flags & Arguments**: IntVarP, StringSliceVarP, BoolVarP
- **Auto-generated Help**: Documentation from code

### Docker & Optimization
- **Multi-stage Builds**: Separating build and runtime environments
- **Layer Caching**: Optimizing rebuild times
- **Static Binaries**: CGO_ENABLED=0 for portability
- **Binary Stripping**: -ldflags for size reduction
- **Security**: Non-root users in containers
- **Base Image Selection**: Alpine vs scratch

## 🚀 Quick Start

### Prerequisites
- Go 1.24+ installed
- Docker (optional, for containerized usage)

### Installation
```bash
# Clone the repository
git clone <your-repo-url>
cd healthcheck

# Download dependencies
go mod download

# Build the binary
go build -o healthcheck

# Run it
./healthcheck check
```

## 📖 Usage

### Basic Health Check
```bash
./healthcheck check
```

Checks default endpoints:
- GitHub API Status
- JSONPlaceholder API
- Dog CEO API

### Custom Timeout
```bash
./healthcheck check --timeout 5
# or short form
./healthcheck check -t 5
```

### Custom URLs
```bash
./healthcheck check --urls https://api.github.com,https://google.com,https://example.com
# or short form
./healthcheck check -u https://api.github.com,https://google.com
```

### Verbose Output
```bash
./healthcheck check --verbose
# or short form
./healthcheck check -v
```

### Combine Flags
```bash
./healthcheck check -t 3 -v --urls https://api.github.com,https://dog.ceo/api/breeds/list/all
```

### Help
```bash
./healthcheck --help
./healthcheck check --help
```

## 🐳 Docker Usage

### Build Options

We've created three Docker configurations to demonstrate optimization:

#### 1. Bloated Image (❌ Don't use in production)
```bash
docker build -f Dockerfile.bloated -t healthcheck:bloated .
# Size: ~550MB
```

#### 2. Optimized Multi-Stage (✅ Recommended)
```bash
docker build -t healthcheck:latest .
# Size: ~15MB (98% reduction!)
```

#### 3. Scratch-based (🚀 Smallest)
```bash
docker build -f Dockerfile.scratch -t healthcheck:scratch .
# Size: ~8MB (99% reduction!)
```

### Run with Docker
```bash
# Default health check
docker run healthcheck:latest

# Custom timeout
docker run healthcheck:latest check --timeout 5

# Custom URLs
docker run healthcheck:latest check --urls https://api.github.com,https://dog.ceo/api/breeds/list/all

# Verbose mode
docker run healthcheck:latest check -v

# Help
docker run healthcheck:latest --help
```

## 📊 Docker Optimization Breakdown

### Size Comparison

| Build Method | Image Size | Reduction | Notes |
|-------------|-----------|-----------|-------|
| Single-stage (bloated) | ~550MB | 0% | Includes Go compiler, source code, build tools |
| Multi-stage (Alpine) | ~15MB | 98% | **Recommended** - Small with debugging tools |
| Multi-stage (scratch) | ~8MB | 99% | Absolute minimum - No shell or utilities |

### What Makes It Small?

1. **Multi-stage builds**: Build in one stage, copy only binary to runtime stage
2. **Alpine Linux base**: 5MB vs Ubuntu's 100MB
3. **Static binary**: `CGO_ENABLED=0` removes C library dependencies
4. **Stripped binary**: `-ldflags="-w -s"` removes debug symbols
5. **Layer caching**: Optimized COPY order for faster rebuilds

### Build Process Visualization
```
Stage 1: Builder (~550MB)          Stage 2: Runtime (~15MB)
┌─────────────────────────┐       ┌──────────────────────┐
│ Go 1.24 compiler        │       │ Alpine Linux (5MB)   │
│ Source code             │  -->  │ Binary only (10MB)   │
│ Dependencies            │       │ CA certificates      │
│ Build tools             │       │ Non-root user        │
└─────────────────────────┘       └──────────────────────┘
      (discarded)                      (final image)
```

## 📁 Project Structure
```
healthcheck/
├── cmd/
│   ├── root.go              # Root command definition
│   └── check.go             # Health check subcommand & logic
├── main.go                  # Application entry point (3 lines!)
├── go.mod                   # Module definition & dependencies
├── go.sum                   # Dependency checksums
├── Dockerfile               # Optimized multi-stage build (Alpine)
├── Dockerfile.bloated       # Single-stage build (for comparison)
├── Dockerfile.scratch       # Ultra-minimal build (scratch base)
├── .dockerignore            # Files to exclude from Docker context
└── README.md                # This file
```

## 🔧 Technical Details

### Concurrency Model

The health checker uses goroutines to check all endpoints simultaneously:
```go
var wg sync.WaitGroup

for _, endpoint := range endpoints {
    wg.Add(1)
    
    go func(ep Endpoint) {
        defer wg.Done()
        result := checkEndpoint(ep)
        printResult(result)
    }(endpoint)
}

wg.Wait() // Block until all checks complete
```

**Benefits:**
- 3 endpoints taking 100ms each: Sequential = 300ms, Concurrent = 100ms
- Scales efficiently to hundreds of endpoints
- Proper synchronization with WaitGroups

### Error Handling Philosophy

Go's explicit error handling (no exceptions):
```go
resp, err := client.Get(endpoint.URL)
if err != nil {
    // Handle error immediately
    return HealthResult{Error: err}
}
defer resp.Body.Close()
```

**Why this is better:**
- Errors are visible in function signatures
- Forces explicit handling at each step
- No hidden control flow from exceptions

### HTTP Client Configuration
```go
client := &http.Client{
    Timeout: time.Duration(timeout) * time.Second,
}
```

**Key decisions:**
- Timeout prevents hanging forever on unresponsive endpoints
- Configurable via flags for different use cases
- Standard library HTTP client (no external dependencies)

## 🏗️ Build Commands

### Local Development
```bash
# Run without building
go run main.go check

# Build for current OS
go build -o healthcheck

# Build with optimizations (what Docker uses)
CGO_ENABLED=0 go build -ldflags="-w -s" -o healthcheck

# Cross-compile for Linux (from Mac/Windows)
GOOS=linux GOARCH=amd64 go build -o healthcheck-linux

# Cross-compile for Windows (from Mac/Linux)
GOOS=windows GOARCH=amd64 go build -o healthcheck.exe
```

### Testing Build Performance
```bash
# Time the build
time go build -o healthcheck

# Check binary size
ls -lh healthcheck

# Check what's in the binary
go tool nm healthcheck | head -20
```

## 🔐 Security Features

1. **Non-root user in Docker**: Runs as UID 1000 (not root)
2. **Static binary**: No shared library dependencies = smaller attack surface
3. **Minimal base image**: Less software = fewer vulnerabilities
4. **CA certificates included**: Proper SSL/TLS verification
5. **No secrets in image**: .dockerignore prevents accidental inclusion

## 📈 Performance Characteristics

- **Startup time**: <10ms (compiled binary)
- **Memory usage**: ~5MB RSS (minimal footprint)
- **Concurrent checks**: Limited only by network and timeout
- **Binary size**: ~10MB (stripped)
- **Docker image size**: 15MB (Alpine) or 8MB (scratch)

## 🎯 Use Cases

This health checker is suitable for:

- **Monitoring**: Periodic health checks of APIs and services
- **CI/CD**: Pre-deployment smoke tests
- **Kubernetes**: Init containers or health check sidecars
- **Development**: Quick endpoint validation during development
- **Learning**: Understanding Go concurrency and Docker optimization

## 🚦 Exit Codes

- `0`: All health checks passed
- `1`: Error occurred (check failed, invalid flags, etc.)

## 📝 Example Output
```
🏥 Health Checker v1.0
━━━━━━━━━━━━━━━━━━━━━━━

✓ HEALTHY [GitHub API]
  URL: https://api.github.com/status
  Status: 200
  Response Time: 145ms

✓ HEALTHY [JSONPlaceholder]
  URL: https://jsonplaceholder.typicode.com/posts/1
  Status: 200
  Response Time: 89ms

✓ HEALTHY [Dog API]
  URL: https://dog.ceo/api/breeds/list/all
  Status: 200
  Response Time: 112ms

━━━━━━━━━━━━━━━━━━━━━━━
✓ Checked 3 endpoints in 152ms
```

## 🛠️ Dependencies
```go
require github.com/spf13/cobra v1.10.1

require (
    github.com/inconshreveable/mousetrap v1.1.0 // indirect
    github.com/spf13/pflag v1.0.9 // indirect
)
```

Only one direct dependency (Cobra), everything else is standard library!

## 📚 Learning Resources

If you're using this project to learn Go, here are key concepts to explore further:

1. **Effective Go**: https://go.dev/doc/effective_go
2. **Go by Example**: https://gobyexample.com/
3. **Cobra Documentation**: https://cobra.dev/
4. **Docker Multi-stage Builds**: https://docs.docker.com/build/building/multi-stage/
5. **Go Concurrency Patterns**: https://go.dev/blog/pipelines

## 🤝 Contributing

This is a learning project, but improvements are welcome! Key areas:

- Additional endpoint protocols (gRPC, WebSocket)
- Better error messages
- Unit tests
- Integration tests
- CI/CD pipeline examples

## 📄 License

MIT License - Feel free to use this for learning or production!

## 🙏 Acknowledgments

Built as a learning project exploring:
- Go fundamentals and concurrency
- Professional CLI design with Cobra
- Docker optimization techniques
- Production-ready Go applications

---