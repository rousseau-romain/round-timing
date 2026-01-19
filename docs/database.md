# Database

## Overview

RoundTiming uses MySQL 8.0 running in Docker. The schema is managed through versioned migrations.

## Connection

Database connection is configured via environment variables:

```.env
DB_HOST=localhost
DB_PORT=3306
DB_DRIVER=mysql
DB_NAME=round_timing
DB_USER=user
DB_PASSWORD=password
```

## Models

### Core Entities

| Model | File | Description |
| ------- | ------ | ------------- |
| User | `model/user.go` | User accounts and profiles |
| Match | `model/match.go` | Game matches/rounds |
| Player | `model/player.go` | Players in matches |
| Team | `model/team.go` | Team groupings |
| Spell | `model/spells.go` | Game abilities/spells |
| Class | `model/class.go` | Character classes |

### Supporting Entities

| Model | File | Description |
| ------- | ------ - ------------|
| FavoriteSpell | `model/favorite-spell.go` | User's favorite spells |
| MatchPlayerSpell | `model/match-player-spell.go` | Spells used in matches |
| UserConfiguration | `model/user-configuration.go` | User settings |
| UserSpectate | `model/user-spectate.go` | Spectator relationships |
| EmailWhiteListed | `model/email-white-listed.go` | Allowed emails |
| FeatureFlag | `model/feature-flag.go` | Feature toggles |
| Language | `model/language.go` | Supported languages |

## Query Building

The project uses `huandu/go-sqlbuilder` for constructing SQL queries:

```go
sb := sqlbuilder.NewSelectBuilder()
sb.Select("id", "name", "email")
sb.From("users")
sb.Where(sb.Equal("id", userID))
sql, args := sb.Build()
```

## Migrations

Migrations are stored in `database/migration/` and managed with `golang-migrate`.

### Commands

```bash
make migration_create <name>  # Create new migration
make migration_up             # Apply migrations
make migration_down           # Rollback one migration
make migration_fix VERSION=N  # Force version (for fixing)
```

### Security

Migrations are encrypted with GPG:

```bash
make db/decrypt  # Decrypt for editing
make db/encrypt  # Encrypt before commit
```

Never commit decrypted `.sql` files to the repository.
