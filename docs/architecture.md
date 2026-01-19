# Architecture

## Overview

RoundTiming is a server-side rendered Go web application that tracks game round timing, likely for esports or competitive gaming scenarios.

## Technology Stack

| Layer | Technology |
| ------- | ------------ |
| Language | Go 1.22 |
| Router | Gorilla Mux |
| Templates | Templ (type-safe HTML) |
| Styling | TailwindCSS |
| Database | MySQL 8.0 |
| Auth | Goth (OAuth2), gorilla/sessions |
| i18n | ctxi18n |

## Directory Structure

```text
round-timing/
├── config/              # App configuration
│   └── config.go        # Config loading from env
├── database/
│   └── migration/       # SQL migrations (encrypted)
├── handlers/            # HTTP route handlers
│   ├── admin.go         # Admin panel routes
│   ├── auth.go          # Authentication routes
│   ├── match.go         # Match management
│   ├── player.go        # Player management
│   └── ...
├── helper/              # Utility functions
├── i18n/
│   └── locales/         # Translation files
├── model/               # Database models
│   ├── db.go            # Database connection
│   ├── user.go          # User model
│   ├── match.go         # Match model
│   └── ...
├── service/
│   └── auth/            # Auth business logic
├── shared/
│   └── components/      # Reusable templ components
├── views/
│   └── page/            # Page templates
├── public/              # Static assets
└── main.go              # Application entry point
```

## Request Flow

```text
HTTP Request
    │
    ▼
main.go (Gorilla Mux Router)
    │
    ▼
handlers/*.go (Route Handler)
    │
    ├──▶ service/*.go (Business Logic)
    │
    ├──▶ model/*.go (Database Query)
    │
    ▼
views/page/*.templ (Render Template)
    │
    ▼
HTTP Response (HTML)
```

## Key Design Decisions

### Server-Side Rendering

All HTML is rendered on the server using Templ templates. This provides:

- Fast initial page loads
- SEO-friendly pages
- Simple state management

### Encrypted Migrations

Database migrations are encrypted with GPG to protect sensitive schema information. Decrypt before development, encrypt before committing.

### OAuth Authentication

Users authenticate via Discord or Google OAuth. Sessions are managed with gorilla/sessions and stored in cookies.
