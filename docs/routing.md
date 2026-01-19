# Routing & API

## Overview

RoundTiming uses [Gorilla Mux](https://github.com/gorilla/mux) for routing. Routes are defined in `main.go` and handlers are in the `handlers/` directory.

## Route Structure

### Public Routes

No authentication required:

| Method | Path | Handler | Description |
| ------ | ---- | ------- | ----------- |
| GET | `/` | `HandlersHome` | Home page |
| GET | `/privacy` | `HandlerPrivacy` | Privacy policy |
| GET | `/cgu` | `HandlerCGU` | Terms of use |
| GET | `/commit-id` | `HandlerCommitId` | Current commit |
| GET | `/version` | `HandlerVersion` | App version |
| GET | `/404` | `HandlersNotFound` | Not found page |
| GET | `/403` | `HandlersForbidden` | Forbidden page |

### Authentication Routes

| Method | Path | Handler | Description |
| ------ | ---- | ------- | ----------- |
| GET | `/signup` | `HandleSignupEmail` | Signup page |
| GET | `/signin` | `HandleLogin` | Login page |
| POST | `/signup` | `HandleCreateUser` | Create account |
| POST | `/signin` | `HandleLoginEmail` | Email login |
| GET | `/auth/{provider}` | `HandleProviderLogin` | OAuth login |
| GET | `/auth/{provider}/callback` | `HandleAuthCallbackFunction` | OAuth callback |
| GET | `/auth/logout/{provider}` | `HandleLogout` | Logout |

### Match Routes (Authenticated)

| Method | Path | Handler | Auth |
| ------ | ---- | ------- | ---- |
| GET | `/match` | `HandlersListMatch` | User |
| POST | `/match` | `HandlersCreateMatch` | User |
| GET | `/match/{idMatch}` | `HandlersMatch` | Owner |
| DELETE | `/match/{idMatch}` | `HandlersDeleteMatch` | Owner |
| GET | `/match/{idMatch}/spectate` | `HandlerSpectateMatch` | Spectator |
| GET | `/match/{idMatch}/start` | `HandlerStartMatchPage` | Owner |
| PATCH | `/match/{idMatch}/reset` | `HandlerResetMatchPage` | Owner |
| GET | `/match/{idMatch}/increase-round` | `HandlerMatchNextRound` | Owner |
| GET | `/match/{idMatch}/table-live` | `HandlerMatchTableLive` | Public |
| GET | `/match/{idMatch}/toggle-mastery/{toggleBool}` | `HandlerToggleMatchMastery` | Owner |

### Player Routes (Authenticated)

| Method | Path | Handler | Auth |
| ------ | ---- | ------- | ---- |
| POST | `/match/{idMatch}/player` | `HandlersCreatePlayer` | Owner |
| PATCH | `/match/{idMatch}/player/{idPlayer}` | `HandlersUpdatePlayer` | Owner |
| DELETE | `/match/{idMatch}/player/{idPlayer}` | `HandlersDeletePlayer` | Owner |

### Spell Routes (Authenticated)

| Method | Path | Handler | Auth |
| ------ | ---- | ------- | ---- |
| GET | `/match/{idMatch}/player-spell/{idPlayerSpell}/use` | `HandlerUsePlayerSpell` | Owner |
| GET | `/match/{idMatch}/player-spell/{idPlayerSpell}/remove-round-recovery` | `HandlerRemoveRoundRecoveryPlayerSpell` | Owner |

### Profile Routes (Authenticated)

| Method | Path | Handler | Description |
| ------ | ---- | ------- | ----------- |
| GET | `/profile` | `HandlersProfile` | Profile page |
| PATCH | `/profile/configuration/{idConfiguration}/toggle-configuration` | `HandlersProfileToggleUserConfiguration` | Toggle config |
| PATCH | `/profile/spell-favorite/{idSpell}/toggle-favorite` | `HandlersToggleSpellFavorite` | Toggle favorite |
| POST | `/profile/user-spectate` | `HandlersProfileAddSpectate` | Add spectator |
| DELETE | `/profile/user-spectate` | `HandlersProfileDeleteSpectate` | Remove spectator |
| PATCH | `/user/{idUser}/locale/{code}` | `HandlersPlayerLanguage` | Change language |

### Admin Routes

| Method | Path | Handler | Description |
| ------ | ---- | ------- | ----------- |
| GET | `/admin/user` | `HandlersListUser` | User list |
| PATCH | `/admin/user/{idUser}/toggle-enabled/{toggleEnabled}` | `HandlersUserEnabled` | Enable/disable user |

### Static Files

```go
r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
```

Serves files from the `public/` directory.

## Middleware

### Language Middleware

Applied to all routes, sets the locale in context:

```go
func languageMiddleware(handler http.Handler, auth *auth.AuthService, slog *slog.Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get language from user or browser preference
        ctx, err := ctxi18n.WithLocale(r.Context(), lang)
        handler.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Auth Middleware

See [authentication.md](./authentication.md) for details on auth middleware functions.

## Handler Structure

Handlers are organized by feature in `handlers/`:

```text
handlers/
├── handlers.go      # Handler struct and constructor
├── admin.go         # Admin routes
├── auth.go          # Authentication routes
├── match.go         # Match management
├── match-spectate.go # Spectator features
├── page.go          # Static pages
├── player.go        # Player management
├── profile.go       # Profile routes
└── favorite-spell.go # Spell favorites
```

### Handler Definition

```go
// handlers/handlers.go
type Handlers struct {
    auth    *auth.AuthService
    slog    *slog.Logger
}

func New(auth *auth.AuthService, slog *slog.Logger) *Handlers {
    return &Handlers{auth: auth, slog: slog}
}
```

### Handler Example

```go
func (h *Handlers) HandlersHome(w http.ResponseWriter, r *http.Request) {
    user, _ := h.auth.GetAuthenticateUserFromRequest(r, h.slog)

    languages, _ := model.GetLanguages()
    navItems := getNavItems(user)
    popinMessage := getPopinMessageFromRequest(r)

    page.HomePage(user, popinMessage, navItems, languages, r.URL.Path).Render(r.Context(), w)
}
```

## HTMX Integration

Many routes are designed for HTMX partial updates:

```go
// Return partial HTML for HTMX swap
r.Handle("/match/{idMatch}/table-live", auth.AllowToBeAuth(handler.HandlerMatchTableLive, authService, versionLogger))
```

Used with:

```html
<div hx-get="/match/1/table-live" hx-trigger="every 1s" hx-swap="innerHTML">
```

## URL Parameters

Gorilla Mux extracts URL parameters:

```go
func (h *Handlers) HandlersMatch(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    matchId, _ := strconv.Atoi(vars["idMatch"])
    // Use matchId...
}
```

## Error Handling

```go
// 404 Not Found
page.NotFoundPage(errorMessage, navItems, languages, r.URL.Path, user).Render(r.Context(), w)

// 403 Forbidden
page.ForbidenPage(errorMessage, navItems, languages, r.URL.Path, user).Render(r.Context(), w)

// 500 Internal Server Error
http.Error(w, err.Error(), http.StatusInternalServerError)

// Redirect
http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
```
