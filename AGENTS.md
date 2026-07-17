# TerminalCharts — AI Agent Guide

## Project Overview
TerminalCharts is a Go CLI tool that generates ASCII charts from CSV and JSON data. It supports bar charts, line plots, pie charts, and histograms with configurable formatting.

## Build & Test
```bash
# Build
go build ./...
go build -o terminalcharts ./cmd/terminalcharts/

# Run tests
go test ./...

# Run with verbose output
go test -v ./...
```

## Architecture
- `cmd/terminalcharts/main.go` — CLI entry point with command routing and flag parsing
- `internal/data/data.go` — CSV and JSON data loaders with auto-detection
- `internal/chart/chart.go` — Chart generation (bar, line, pie, histogram) with options

## Key Design Decisions
- All chart rendering is done in-memory with strings.Builder
- No external dependencies — pure Go standard library
- Color output via ANSI escape codes (works in any terminal)
- Data auto-detection: CSV vs JSON, header detection, column name inference
- JSON output for programmatic use (CI/CD integration)

## Data Loading
- **CSV**: Auto-detects headers, supports label/value column specification
- **JSON**: Supports array of objects, simple objects, and simple arrays
- Column auto-detection uses common names (label/name/x/key for labels, value/y/count for values)

## Adding a New Chart Type
1. Add a new function in `internal/chart/chart.go`
2. Add a case in the switch statement in `main.go`
3. Add tests in `internal/chart/chart_test.go`

## Dependencies
- Standard library only (no external dependencies)