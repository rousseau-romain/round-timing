#!/usr/bin/env bash
set -euo pipefail

# Automated version bumping script
# Usage: ./scripts/release.sh [auto|major|minor|patch] [--dry-run]

DRY_RUN=false
BUMP_TYPE="auto"
CONFIG_FILE="config/config.go"

for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=true ;;
    *) BUMP_TYPE="$arg" ;;
  esac
done

# --- Preflight checks ---
if [[ "$DRY_RUN" == false ]]; then
  CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
  if [[ "$CURRENT_BRANCH" != "master" ]]; then
    echo "Error: must be on master branch (currently on '$CURRENT_BRANCH')."
    exit 1
  fi

  if [[ -n "$(git status --porcelain)" ]]; then
    echo "Error: working directory is not clean. Commit or stash changes first."
    exit 1
  fi
fi

# --- Get latest tag ---
LATEST_TAG=$(git tag -l 'v[0-9]*' --sort=-v:refname | head -1)
LATEST_TAG="${LATEST_TAG:-v0.0.0}"
echo "Latest tag: $LATEST_TAG"

# Parse version components
VERSION_RE='^v([0-9]+)\.([0-9]+)\.([0-9]+)$'
if [[ ! "$LATEST_TAG" =~ $VERSION_RE ]]; then
  echo "Error: latest tag '$LATEST_TAG' does not match vMAJOR.MINOR.PATCH"
  exit 1
fi
MAJOR="${BASH_REMATCH[1]}"
MINOR="${BASH_REMATCH[2]}"
PATCH="${BASH_REMATCH[3]}"

# --- Determine bump type ---
if [[ "$BUMP_TYPE" == "auto" ]]; then
  COMMITS=$(git log "${LATEST_TAG}..HEAD" --pretty=format:"%s" --no-merges)

  if [[ -z "$COMMITS" ]]; then
    echo "Error: no new commits since $LATEST_TAG"
    exit 1
  fi

  # Check for breaking changes
  if echo "$COMMITS" | grep -qiE 'BREAKING CHANGE|^[a-z]+(\(.+\))?!:'; then
    BUMP_TYPE="major"
  elif echo "$COMMITS" | grep -qE '^feat(\(.+\))?:'; then
    BUMP_TYPE="minor"
  elif echo "$COMMITS" | grep -qE '^fix(\(.+\))?:'; then
    BUMP_TYPE="patch"
  else
    echo "No conventional commits (feat/fix/BREAKING CHANGE) found since $LATEST_TAG."
    echo "Defaulting to patch bump."
    BUMP_TYPE="patch"
  fi
fi

# --- Compute next version ---
case "$BUMP_TYPE" in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
  *)
    echo "Error: invalid bump type '$BUMP_TYPE'. Use: auto, major, minor, patch"
    exit 1
    ;;
esac

NEXT_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
echo "Bump type: $BUMP_TYPE"
echo "Next version: $NEXT_VERSION"

if [[ "$DRY_RUN" == true ]]; then
  echo ""
  echo "[dry-run] Would update $CONFIG_FILE with VERSION = \"${NEXT_VERSION}\""
  echo "[dry-run] Would generate CHANGELOG.md"
  echo "[dry-run] Would commit and tag ${NEXT_VERSION}"
  exit 0
fi

# --- Update config/config.go ---
sed -i "s/VERSION = \"v[0-9]*\.[0-9]*\.[0-9]*\"/VERSION = \"${NEXT_VERSION}\"/" "$CONFIG_FILE"
echo "Updated $CONFIG_FILE"

# --- Generate CHANGELOG.md ---
git cliff --tag "$NEXT_VERSION" -o CHANGELOG.md
echo "Generated CHANGELOG.md"

# --- Commit and tag ---
git add "$CONFIG_FILE" CHANGELOG.md
git commit -m "chore: release ${NEXT_VERSION}"
git tag -a "$NEXT_VERSION" -m "Release ${NEXT_VERSION}"

echo ""
echo "Release ${NEXT_VERSION} created successfully!"
echo ""
echo "Next steps:"
echo "  git push origin master --tags    # or: make release/push"
echo "  goreleaser release --clean       # or: make release/github (optional, requires GITHUB_TOKEN)"
