#!/usr/bin/env bash
#
# Shell Script Analysis Script
# Runs shellcheck and other linters on shell scripts
# Compatible with macOS and Linux
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERBOSE="${VERBOSE:-false}"
FIX="${FIX:-false}"

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Darwin) PLATFORM="macos" ;;
    Linux)  PLATFORM="linux" ;;
    *)      PLATFORM="unknown" ;;
esac

# Colors for output (disable if not a tty)
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${BLUE}[INFO]${NC} $*" >&2
    fi
}

error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

success() {
    echo -e "${GREEN}[OK]${NC} $*"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

header() {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  $*${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
}

# Check if a command exists
check_tool() {
    local tool="$1"
    local install_cmd="$2"

    if ! command -v "$tool" &>/dev/null; then
        warn "$tool not found."
        warn "To install: $install_cmd"
        return 1
    fi
    return 0
}

# Ensure required tools are installed
ensure_tools() {
    header "Checking Required Tools"

    local missing=0

    local shellcheck_install="brew install shellcheck"
    local shfmt_install="brew install shfmt"
    if [[ "$PLATFORM" == "linux" ]]; then
        shellcheck_install="apt install shellcheck"
        shfmt_install="go install mvdan.cc/sh/v3/cmd/shfmt@latest"
    fi

    if ! check_tool "shellcheck" "$shellcheck_install"; then
        missing=1
    else
        success "shellcheck $(shellcheck --version | grep version: | awk '{print $2}')"
    fi

    if ! check_tool "shfmt" "$shfmt_install"; then
        log "shfmt not available"
    else
        success "shfmt $(shfmt --version 2>&1 || echo 'installed')"
    fi

    return $missing
}

# Find shell scripts
find_scripts() {
    local search_path="${1:-.}"

    # Optimized find command:
    # 1. Prune (ignore) node_modules, vendor, .git to avoid scanning massive directories
    # 2. Look for .sh files
    # 3. Look for files with shell shebangs
    
    # We use a temporary file to store results to avoid pipe subshell issues
    local tmp_scripts
    tmp_scripts=$(mktemp)

    find "$search_path" \
        \( -name "node_modules" -o -name "vendor" -o -name ".git" -o -name "dist" -o -name "build" \) -prune \
        -o -type f \( -name "*.sh" -o -name "*.bash" -o -name "*.zsh" \) -print > "$tmp_scripts"

    # Also check files without extension for shebang
    # This is slower, so we only check files that are likely scripts (executable or common names?)
    # For safety/speed, we'll stick to extension OR standard locations if needed.
    # But to match previous behavior safely:
    
    find "$search_path" \
        \( -name "node_modules" -o -name "vendor" -o -name ".git" -o -name "dist" -o -name "build" \) -prune \
        -o -type f ! -name "*.*" -exec grep -l -m 1 '^#!.*\(bash\|sh\|zsh\)' {} + >> "$tmp_scripts" 2>/dev/null || true

    sort -u "$tmp_scripts"
    rm "$tmp_scripts"
}

# Run shellcheck
run_shellcheck() {
    header "Running ShellCheck"

    if ! command -v shellcheck &>/dev/null; then
        error "shellcheck not available"
        return 1
    fi

    local scripts
    scripts=$(find_scripts .)

    if [[ -z "$scripts" ]]; then
        log "No shell scripts found"
        return 0
    fi

    local script_count
    script_count=$(echo "$scripts" | wc -l | tr -d ' ')
    log "Found $script_count shell scripts"

    local failed=0
    local checked=0

    while IFS= read -r script; do
        if [[ -f "$script" ]]; then
            # Run shellcheck with colors forced if tty, but we capture output so maybe not?
            if shellcheck -x "$script" 2>&1; then
                ((checked++)) || true
            else
                ((failed++)) || true
                warn "Issues in: $script"
            fi
        fi
    done <<< "$scripts"

    if [[ $failed -eq 0 ]]; then
        success "ShellCheck passed ($checked scripts checked)"
        return 0
    else
        warn "ShellCheck found issues in $failed/$((checked + failed)) scripts"
        return 1
    fi
}

# Run shfmt
run_shfmt() {
    header "Running shfmt (Format Check)"

    if ! command -v shfmt &>/dev/null; then
        warn "shfmt not available, skipping format check"
        return 0
    fi

    local scripts
    scripts=$(find_scripts .)

    if [[ -z "$scripts" ]]; then
        log "No shell scripts found"
        return 0
    fi

    local needs_format=()

    while IFS= read -r script; do
        if [[ -f "$script" ]]; then
            if ! shfmt -d "$script" &>/dev/null; then
                needs_format+=("$script")
            fi
        fi
    done <<< "$scripts"

    if [[ ${#needs_format[@]} -eq 0 ]]; then
        success "All scripts properly formatted"
        return 0
    else
        if [[ "$FIX" == "true" ]]; then
            for script in "${needs_format[@]}"; do
                shfmt -w "$script"
                log "Formatted: $script"
            done
            success "Formatted ${#needs_format[@]} scripts"
            return 0
        else
            warn "Scripts need formatting:"
            printf '%s\n' "${needs_format[@]}"
            echo ""
            echo "Run with --fix to auto-format, or: shfmt -w <script>"
            return 1
        fi
    fi
}

# Check POSIX compatibility
check_posix() {
    header "Checking POSIX Compatibility"

    if ! command -v checkbashisms &>/dev/null; then
        log "checkbashisms not available. Install with: brew install devscripts (optional)"
        return 0
    fi

    local scripts
    scripts=$(find_scripts .)

    if [[ -z "$scripts" ]]; then
        return 0
    fi

    local posix_scripts=0

    while IFS= read -r script; do
        if [[ -f "$script" ]] && head -1 "$script" | grep -q '^#!/bin/sh'; then
            ((posix_scripts++)) || true
            if ! checkbashisms "$script" 2>&1; then
                warn "Bashisms in POSIX script: $script"
            fi
        fi
    done <<< "$scripts"

    if [[ $posix_scripts -eq 0 ]]; then
        log "No /bin/sh scripts found (POSIX check skipped)"
    else
        success "Checked $posix_scripts POSIX scripts"
    fi

    return 0
}

# Check for TODO/FIXME comments
check_todos() {
    header "Checking TODO/FIXME Comments"

    local count
    count=$(grep -rn --include="*.sh" -E "(TODO|FIXME|XXX|HACK)" . 2>/dev/null | wc -l | tr -d ' ')

    if [[ "$count" -gt 0 ]]; then
        warn "Found $count TODO/FIXME comments:"
        grep -rn --include="*.sh" -E "(TODO|FIXME|XXX|HACK)" . 2>/dev/null | head -20
        if [[ "$count" -gt 20 ]]; then
            echo "... and $((count - 20)) more"
        fi
    else
        success "No TODO/FIXME comments found"
    fi
}

# Main execution
main() {
    local exit_code=0

    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║            Shell Script Analysis Report                       ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════════╝${NC}"

    ensure_tools || exit_code=1

    run_shellcheck || exit_code=1
    run_shfmt || exit_code=1
    check_posix
    check_todos

    header "Summary"
    if [[ $exit_code -eq 0 ]]; then
        success "All checks passed!"
    else
        warn "Some checks failed. Review the output above."
    fi

    return $exit_code
}

# Show usage
usage() {
    cat <<EOF
Usage: $(basename "$0") [options]

Options:
    -h, --help      Show this help
    -v, --verbose   Verbose output
    -f, --fix       Attempt to fix issues automatically (formatting only)

Environment Variables:
    VERBOSE=true    Enable verbose output
    FIX=true       Enable auto-fix mode (formatting only)

Examples:
    $(basename "$0")        # Run all checks
    $(basename "$0") --fix  # Run and auto-format
EOF
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help) usage; exit 0 ;;
        -v|--verbose) VERBOSE=true ;;
        -f|--fix) FIX=true ;;
        *) error "Unknown option: $1"; usage; exit 1 ;;
    esac
    shift
done

main
