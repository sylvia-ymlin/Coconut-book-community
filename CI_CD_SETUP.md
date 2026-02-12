# CI/CD Pipeline Documentation

## Overview

BookCommunity uses GitHub Actions for continuous integration and deployment. The pipeline automatically tests, builds, and deploys the application on every commit.

## Workflows

### 1. CI Pipeline (`.github/workflows/ci.yaml`)

Runs on every push and pull request to `main` and `develop` branches.

**Jobs:**

#### Lint (Code Quality)
- Runs `golangci-lint` with 15+ linters
- Checks code style, security, and best practices
- Configuration: `.golangci.yml`

#### Test (Unit Tests)
- Runs all unit tests with race detection
- Generates coverage report
- Uploads coverage to Codecov
- Uses PostgreSQL and Redis service containers
- Coverage artifacts uploaded to GitHub

#### Build
- Compiles binary for Linux/amd64
- Includes version, commit, and build time
- Uploads binary as artifact

#### Security
- Runs Gosec security scanner
- Checks for common security vulnerabilities
- Uploads results to GitHub Security tab

### 2. Docker Workflow (`.github/workflows/docker.yaml`)

Builds and publishes Docker images to GitHub Container Registry (ghcr.io).

**Features:**
- Multi-platform builds (amd64, arm64)
- Image tagging strategy:
  - `main` → `latest`
  - `v1.2.3` → `1.2.3`, `1.2`, `v1.2.3`
  - PR → `pr-123`
- Vulnerability scanning with Trivy
- Caching for faster builds

### 3. Release Workflow (`.github/workflows/release.yaml`)

Automated release creation when you push a version tag.

**Process:**
1. Push tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GoReleaser creates:
   - GitHub Release with changelog
   - Binaries for Linux, macOS, Windows (amd64, arm64)
   - Checksums file
   - Archives (.tar.gz, .zip)

## Configuration Files

### `.golangci.yml`
Linter configuration with enabled checks:
- `errcheck` - Unchecked errors
- `gosimple` - Code simplification
- `govet` - Go vet analysis
- `gosec` - Security issues
- `revive` - Style guide
- `dupl` - Code duplication
- And 10+ more...

### `.goreleaser.yaml`
Release automation configuration:
- Multi-OS builds
- Archive creation
- Changelog generation
- GitHub release creation

## Local Development

### Prerequisites
```bash
# Install tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/goreleaser/goreleaser@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### Run CI Locally

```bash
# Full CI pipeline
make ci

# Individual steps
make lint        # Run linter
make test        # Run tests with coverage
make build       # Build binary
make security    # Security scan
```

### Test Workflow

```bash
# Run tests
make test

# View coverage
open coverage.html

# Run with services
docker-compose up -d postgres redis
go test ./... -v
```

## GitHub Secrets

Required secrets for workflows:

| Secret | Description | Required For |
|--------|-------------|--------------|
| `GITHUB_TOKEN` | Auto-provided | Docker, Release |
| `CODECOV_TOKEN` | Codecov upload | Coverage (optional) |

## Pipeline Status Badges

Add to README:

```markdown
[![CI](https://github.com/sylvia-ymlin/Coconut-book-community/workflows/CI%20Pipeline/badge.svg)](https://github.com/sylvia-ymlin/Coconut-book-community/actions)
[![Docker](https://github.com/sylvia-ymlin/Coconut-book-community/workflows/Docker%20Build%20%26%20Push/badge.svg)](https://github.com/sylvia-ymlin/Coconut-book-community/actions)
[![codecov](https://codecov.io/gh/sylvia-ymlin/Coconut-book-community/branch/main/graph/badge.svg)](https://codecov.io/gh/sylvia-ymlin/Coconut-book-community)
```

## Troubleshooting

### Lint Failures
```bash
# Auto-fix issues
make lint-fix

# Check specific issues
golangci-lint run --disable-all --enable=errcheck
```

### Test Failures
```bash
# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestUserHandler ./internal/app/handlers/user
```

### Docker Build Failures
```bash
# Test build locally
docker build -t bookcommunity:test .

# Check Dockerfile syntax
docker build --check -t bookcommunity:test .
```

## Performance Optimization

### Caching
- Go modules cached by GitHub Actions
- Docker layers cached using `cache-from/cache-to`
- golangci-lint cache enabled

### Parallel Execution
- Lint and Test jobs run in parallel
- Build waits for both to complete

### Service Containers
- PostgreSQL and Redis start in parallel
- Health checks ensure readiness

## Best Practices

1. **Always run `make ci` before pushing**
2. **Write tests for new features**
3. **Keep test coverage above 60%**
4. **Fix lint issues immediately**
5. **Use conventional commits** for changelog generation
6. **Tag releases** with semantic versioning

## Commit Message Format

For automatic changelog generation:

```
feat: add user authentication
fix: resolve cache invalidation bug
docs: update API documentation
test: add user service tests
chore: update dependencies
```

## Next Steps

1. **Increase Test Coverage**
   - Add integration tests
   - Mock external dependencies
   - Target 80%+ coverage

2. **Add Integration Tests**
   - API endpoint tests
   - Database integration tests
   - Cache integration tests

3. **Set up Codecov**
   - Sign up at codecov.io
   - Add `CODECOV_TOKEN` to secrets
   - Monitor coverage trends

4. **Deploy to Kubernetes**
   - Add deployment workflow
   - Use Helm for releases
   - Implement blue-green deployment
