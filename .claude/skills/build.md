# /build - Build Project

Build the project for deployment.

## Usage

```text
/build          # Full production build
/build templ    # Generate templ files only
/build css      # Build TailwindCSS only
/build check    # Check for build errors without building
```

## Instructions

### /build (full)

Run a complete production build:

1. `make build/templ` - Generate Go code from .templ files
2. `make build/tailwind` - Build and minify CSS
3. `go build -o tmp/main` - Compile Go binary

### /build templ

Run `templ generate` or `make build/templ` to regenerate all `*_templ.go` files from `.templ` sources.

### /build css

Run `make build/tailwind` to build minified TailwindCSS output.

### /build check

Run `go build ./...` to check for compilation errors without producing a binary.
Also run `go vet ./...` to check for common issues.
