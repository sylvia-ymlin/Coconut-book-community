# Quick Reference Card

## Daily Development Commands

```bash
# Before starting work
make ci                    # Run full CI locally

# During development
make test                  # Run tests
make lint                  # Check code quality
make lint-fix              # Auto-fix issues
make build                 # Build binary

# Before committing
make ci                    # Final check
git add .
git commit -m "type: description"
git push
```

## CI/CD Workflows

### Automatic Triggers
- **CI Pipeline**: Every push/PR to main/develop
- **Docker Build**: Every push to main
- **Release**: Push tag `v*`

### Manual Triggers
```bash
# Create release
git tag v1.0.0
git push origin v1.0.0

# View workflows
https://github.com/sylvia-ymlin/Coconut-book-community/actions
```

## Common Issues & Fixes

### Lint Failures
```bash
make lint-fix              # Auto-fix
make fmt                   # Format code
```

### Test Failures
```bash
go test -v ./...           # Verbose output
go test -run TestName      # Run specific test
```

### Build Failures
```bash
go mod tidy                # Clean dependencies
make build                 # Test build
```

## Project Status

✅ **Implemented**
- CI/CD Pipeline (GitHub Actions)
- Code Quality (golangci-lint)
- Security Scanning (Gosec, Trivy)
- Automated Releases (GoReleaser)
- Dependency Updates (Dependabot)

🚧 **Next Steps**
- Increase test coverage (target: 80%)
- API documentation (Swagger)
- Integration tests
