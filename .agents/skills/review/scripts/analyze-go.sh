#!/usr/bin/env bash
#
# Go Code Review Script
# Runs review-specific checks: module health, tests, and TODO scanning.
#
# Linting (go vet, golangci-lint, CSGDAA org linters) is handled by /lint.
# Run /lint first, then this script for the full review.
#
# Compatible with macOS and Linux.
#
set -euo pipefail

VERBOSE="${VERBOSE:-false}"
FIX="${FIX:-false}"

# Colors (disable if not a tty)
if [[ -t 1 ]]; then
    RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
else
    RED=''; GREEN=''; YELLOW=''; BLUE=''; NC=''
fi

log()     { [[ "$VERBOSE" == "true" ]] && echo -e "${BLUE}[INFO]${NC} $*" >&2 || true; }
error()   { echo -e "${RED}[ERROR]${NC} $*" >&2; }
success() { echo -e "${GREEN}[OK]${NC} $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $*"; }
header()  {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $*${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
}

# ── Module health ────────────────────────────────────────────────────────────

run_mod_checks() {
    header "Checking Module Health"

    # Check if go.mod is tidy (non-destructive: backup and restore)
    log "Running go mod tidy (check only)..."
    cp go.mod go.mod.bak 2>/dev/null || true
    cp go.sum go.sum.bak 2>/dev/null || true

    go mod tidy

    if diff -q go.mod go.mod.bak &>/dev/null && diff -q go.sum go.sum.bak &>/dev/null; then
        success "go.mod and go.sum are tidy"
        rm -f go.mod.bak go.sum.bak
    elif [[ "$FIX" == "true" ]]; then
        success "go.mod and go.sum have been tidied"
        rm -f go.mod.bak go.sum.bak
    else
        warn "go.mod or go.sum need tidying. Run 'go mod tidy' or use --fix"
        mv go.mod.bak go.mod 2>/dev/null || true
        mv go.sum.bak go.sum 2>/dev/null || true
    fi

    # Verify module integrity
    log "Verifying module integrity..."
    if go mod verify 2>&1; then
        success "Module dependencies verified"
    else
        error "Module verification failed"
        return 1
    fi

    return 0
}

# ── Tests ────────────────────────────────────────────────────────────────────

run_tests() {
    header "Running Tests with Race Detector"

    if go test -race -short ./... 2>&1; then
        success "Tests passed"
        return 0
    else
        error "Tests failed"
        return 1
    fi
}

# ── TODOs ────────────────────────────────────────────────────────────────────

check_todos() {
    header "Checking TODO/FIXME Comments"

    local count
    count=$(grep -rn --include="*.go" -E "(TODO|FIXME|XXX|HACK)" . 2>/dev/null | wc -l | tr -d ' ')

    if [[ "$count" -gt 0 ]]; then
        warn "Found $count TODO/FIXME comments:"
        grep -rn --include="*.go" -E "(TODO|FIXME|XXX|HACK)" . 2>/dev/null | head -20
        if [[ "$count" -gt 20 ]]; then
            echo "... and $((count - 20)) more"
        fi
    else
        success "No TODO/FIXME comments found"
    fi
}

# ── Main ─────────────────────────────────────────────────────────────────────

main() {
    local exit_code=0

    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║              Go Code Review                                   ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════════╝${NC}"

    if [[ ! -f "go.mod" ]]; then
        error "No go.mod found in current directory"
        exit 1
    fi

    if ! command -v go &>/dev/null; then
        error "go is not installed"
        exit 1
    fi

    run_mod_checks || exit_code=1

    if [[ "${SKIP_TESTS:-false}" != "true" ]]; then
        run_tests || exit_code=1
    fi

    check_todos

    header "Summary"
    if [[ $exit_code -eq 0 ]]; then
        success "All review checks passed!"
    else
        warn "Some checks failed. Review the output above."
    fi

    return $exit_code
}

usage() {
    cat <<EOF
Usage: $(basename "$0") [options]

Run this after /lint for the full review. This script handles review-specific
checks that /lint does not cover: module health, tests, and TODO scanning.

Options:
    -h, --help      Show this help
    -v, --verbose   Verbose output
    -f, --fix       Attempt to fix issues (e.g., go mod tidy)
    --skip-tests    Skip running tests

Environment Variables:
    VERBOSE=true            Enable verbose output
    FIX=true               Enable auto-fix mode
    SKIP_TESTS=true        Skip running tests

Examples:
    $(basename "$0")                    # Run all review checks
    $(basename "$0") --fix              # Run with auto-fix
    SKIP_TESTS=true $(basename "$0")    # Skip tests
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help) usage; exit 0 ;;
        -v|--verbose) VERBOSE=true ;;
        -f|--fix) FIX=true ;;
        --skip-tests) SKIP_TESTS=true ;;
        *) error "Unknown option: $1"; usage; exit 1 ;;
    esac
    shift
done

main
