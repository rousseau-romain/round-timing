# /templ - Templ Template Commands

Work with Templ HTML templates.

## Usage

```text
/templ generate   # Generate all templ files
/templ watch      # Watch and regenerate on changes
/templ fmt        # Format all templ files
```

## Instructions

### /templ generate

Run `templ generate` to compile all `.templ` files to `*_templ.go` files.

### /templ watch

Run `templ generate -watch` to watch for changes and auto-regenerate.
This is automatically included in `make live`.

### /templ fmt

Run `templ fmt .` to format all templ files in the project.

## Templ File Locations

- `views/page/` - Full page templates
- `shared/components/` - Reusable UI components

## Tips

- Never edit `*_templ.go` files directly - they are generated
- Templ syntax is similar to Go with HTML
- Use `@component()` to include other components
- Use `{ expression }` for Go expressions in templates
