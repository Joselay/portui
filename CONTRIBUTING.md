# Contributing to portui

Thanks for your interest in contributing!

## Getting Started

1. Fork and clone the repo
2. Make sure you have Go 1.26+ installed
3. Build and run:

```bash
go build -o portui
./portui
```

## Making Changes

1. Create a branch from `main`
2. Make your changes
3. Run checks before submitting:

```bash
go build ./...
go vet ./...
go test ./...
```

4. Open a pull request against `main`

## Reporting Bugs

Open an issue with:
- What you expected to happen
- What actually happened
- Your OS and Go version

## Adding Features

Please open an issue first to discuss the feature before starting work. This helps avoid duplicate effort and ensures the feature fits the project direction.
