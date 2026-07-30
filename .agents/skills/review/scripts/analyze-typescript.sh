#!/usr/bin/env bash
#
# TypeScript/JavaScript Code Analysis Script
# Runs multiple static analysis tools and reports findings
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

# Detect package manager
detect_pm() {
    if [[ -f "bun.lockb" ]]; then
        echo "bun"
    elif [[ -f "pnpm-lock.yaml" ]]; then
        echo "pnpm"
    elif [[ -f "yarn.lock" ]]; then
        echo "yarn"
    else
        echo "npm"
    fi
}

PM=$(detect_pm)

# Run command with package manager
pm_run() {
    case "$PM" in
        bun) bun run "$@" ;;
        pnpm) pnpm run "$@" ;;
        yarn) yarn run "$@" ;;
        npm) npm run "$@" ;;
    esac
}

pm_exec() {
    case "$PM" in
        bun) bunx "$@" ;;
        pnpm) pnpm exec "$@" ;;
        yarn) yarn exec "$@" ;;
        npm) npx "$@" ;;
    esac
}

# Run a tool using DLX (download and execute ephemeral)
# Falls back to pm_exec if locally installed, or downloads if needed
pm_dlx() {
    local pkg="$1"
    shift
    
    # Check if we should try to use local version first via pm_exec?
    # Actually pm_exec usually runs local binary if in node_modules/.bin
    # But npx/bunx also download if missing.
    # However, 'yarn exec' does NOT download. 'yarn dlx' does.
    # 'pnpm exec' does NOT download. 'pnpm dlx' does.
    
    case "$PM" in
        bun) bunx "$pkg" "$@" ;;  # bunx auto-downloads
        pnpm) pnpm dlx "$pkg" "$@" ;;
        yarn) yarn dlx "$pkg" "$@" ;;
        npm) npx -y "$pkg" "$@" ;; # -y to suppress confirmation
    esac
}

# Check if a command exists
check_tool() {
    local tool="$1"
    local install_cmd="$2"

    if ! command -v "$tool" &>/dev/null; then
        warn "$tool not found. Install with: $install_cmd"
        return 1
    fi
    return 0
}

# Ensure required tools are installed
ensure_tools() {
    header "Checking Required Tools"

    local missing=0

    case "$PM" in
        bun)
            if ! command -v bun &>/dev/null; then
                 warn "Bun not found. Please install it manually from https://bun.sh"
                 missing=1
            else
                 success "bun $(bun --version)"
            fi
            ;;
        pnpm)
            if ! command -v pnpm &>/dev/null; then
                warn "pnpm not found. Install with: npm install -g pnpm"
                missing=1
            else
                success "pnpm $(pnpm --version)"
            fi
            ;;
        yarn)
            if ! command -v yarn &>/dev/null; then
                warn "yarn not found. Install with: npm install -g yarn"
                missing=1
            else
                success "yarn $(yarn --version)"
            fi
            ;;
        npm)
             # npm usually comes with node
            if ! command -v npm &>/dev/null; then
                warn "npm not found."
                missing=1
            else
                success "npm $(npm --version)"
            fi
            ;;
    esac

    if check_tool "node" "Install Node.js from nodejs.org"; then
        success "node $(node --version)"
    else
        missing=1
    fi

    return $missing
}

# Run TypeScript compiler
run_tsc() {
    header "Running TypeScript Compiler"

    if [[ ! -f "tsconfig.json" ]]; then
        log "No tsconfig.json found, skipping TypeScript check"
        return 0
    fi

    if pm_exec tsc --noEmit 2>&1; then
        success "TypeScript compilation passed"
        return 0
    else
        error "TypeScript compilation failed"
        return 1
    fi
}

# Run ESLint
run_eslint() {
    header "Running ESLint"

    # Check for ESLint config
    if [[ ! -f "eslint.config.js" ]] && [[ ! -f "eslint.config.mjs" ]] && [[ ! -f ".eslintrc.js" ]] && [[ ! -f ".eslintrc.json" ]] && [[ ! -f ".eslintrc.yml" ]]; then
        if ! grep -q '"eslint"' package.json 2>/dev/null && ! grep -q '"@eslint' package.json 2>/dev/null; then
            warn "ESLint not configured, skipping"
            return 0
        fi
    fi

    local fix_flag=""
    if [[ "$FIX" == "true" ]]; then
        fix_flag="--fix"
    fi

    if pm_exec eslint $fix_flag . 2>&1; then
        success "ESLint passed"
        return 0
    else
        warn "ESLint found issues"
        return 1
    fi
}

# Run Prettier
run_prettier() {
    header "Running Prettier (Format Check)"

    # Check for Prettier config
    if [[ ! -f ".prettierrc" ]] && [[ ! -f ".prettierrc.js" ]] && [[ ! -f "prettier.config.js" ]]; then
        if ! grep -q '"prettier"' package.json 2>/dev/null; then
            log "Prettier not configured, skipping"
            return 0
        fi
    fi

    local src_pattern="${SRC_PATTERN:-src/**/*.{ts,tsx,js,jsx}}"

    if [[ "$FIX" == "true" ]]; then
        pm_exec prettier --write "$src_pattern" 2>&1
        success "Prettier formatting applied"
    else
        if pm_exec prettier --check "$src_pattern" 2>&1; then
            success "Prettier check passed"
        else
            warn "Prettier found formatting issues. Run with --fix to auto-format"
            return 1
        fi
    fi

    return 0
}

# Check for circular dependencies
run_madge() {
    header "Checking for Circular Dependencies"

    # Use dlx to run madge without installing if needed
    # But check if it is locally installed first to avoid download delay?
    # pm_exec usually fails if not installed for yarn/pnpm (unless using dlx).
    
    local src_dir="src"
    if [[ ! -d "$src_dir" ]]; then
        src_dir="."
    fi

    local result
    # Try pm_exec first (local), then fall back to dlx
    if ! result=$(pm_exec madge --circular --extensions ts,tsx,js,jsx "$src_dir" 2>&1); then
         log "madge not found locally, trying to run via $PM dlx..."
         result=$(pm_dlx madge --circular --extensions ts,tsx,js,jsx "$src_dir" 2>&1 || true)
    fi

    if echo "$result" | grep -q "No circular dependency found"; then
        success "No circular dependencies found"
        return 0
    elif [[ -n "$result" ]]; then
        warn "Circular dependencies detected:"
        echo "$result"
        return 1
    else
        success "No circular dependencies found"
        return 0
    fi
}

# Check for unused dependencies
run_depcheck() {
    header "Checking for Unused Dependencies"

    local result
    if ! result=$(pm_exec depcheck 2>&1); then
         log "depcheck not found locally, trying to run via $PM dlx..."
         result=$(pm_dlx depcheck 2>&1 || true)
    fi

    if echo "$result" | grep -q "No depcheck issue"; then
        success "No unused dependencies"
        return 0
    elif [[ -n "$result" ]]; then
        warn "Dependency issues found:"
        echo "$result" | head -30
        return 1
    else
        success "No unused dependencies"
        return 0
    fi
}

# Run npm audit
run_audit() {
    header "Checking for Security Vulnerabilities"

    case "$PM" in
        bun)
            # Bun doesn't have audit yet
            log "Bun doesn't support audit yet, using npm audit instead"
            npm audit --omit=dev 2>&1 || warn "npm audit found vulnerabilities"
            ;;
        pnpm)
            if pnpm audit --prod 2>&1; then
                success "pnpm audit: no vulnerabilities"
            else
                warn "pnpm audit found vulnerabilities"
            fi
            ;;
        yarn)
            if yarn audit --groups dependencies 2>&1; then
                success "yarn audit: no vulnerabilities"
            else
                warn "yarn audit found vulnerabilities"
            fi
            ;;
        npm)
            if npm audit --omit=dev 2>&1; then
                success "npm audit: no vulnerabilities"
            else
                warn "npm audit found vulnerabilities"
            fi
            ;;
    esac
}

# Run tests
run_tests() {
    header "Running Tests"

    # Check for test script in package.json
    if ! grep -q '"test"' package.json 2>/dev/null; then
        log "No test script found in package.json, skipping"
        return 0
    fi

    if pm_run test 2>&1; then
        success "Tests passed"
        return 0
    else
        error "Tests failed"
        return 1
    fi
}

# Check for TODO/FIXME comments
check_todos() {
    header "Checking TODO/FIXME Comments"

    local src_dir="src"
    if [[ ! -d "$src_dir" ]]; then
        src_dir="."
    fi

    local count
    count=$(grep -rn --include="*.ts" --include="*.tsx" --include="*.js" --include="*.jsx" -E "(TODO|FIXME|XXX|HACK)" "$src_dir" 2>/dev/null | wc -l | tr -d ' ')

    if [[ "$count" -gt 0 ]]; then
        warn "Found $count TODO/FIXME comments:"
        grep -rn --include="*.ts" --include="*.tsx" --include="*.js" --include="*.jsx" -E "(TODO|FIXME|XXX|HACK)" "$src_dir" 2>/dev/null | head -20
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
    echo -e "${BLUE}║           TypeScript Code Analysis Report                     ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Package manager: $PM"

    # Check for package.json
    if [[ ! -f "package.json" ]]; then
        error "No package.json found in current directory"
        exit 1
    fi

    ensure_tools || exit_code=1

    run_tsc || exit_code=1
    run_eslint || exit_code=1
    run_prettier || exit_code=1
    run_madge || exit_code=1
    run_depcheck || exit_code=1
    run_audit

    if [[ "${SKIP_TESTS:-false}" != "true" ]]; then
        run_tests || exit_code=1
    fi

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
    -f, --fix       Attempt to fix issues automatically
    --skip-tests    Skip running tests

Environment Variables:
    VERBOSE=true            Enable verbose output
    FIX=true               Enable auto-fix mode
    SKIP_TESTS=true        Skip running tests
    SRC_PATTERN="src/**/*" Pattern for source files

Examples:
    $(basename "$0")                    # Run all checks
    $(basename "$0") --fix              # Run and auto-fix issues
    SKIP_TESTS=true $(basename "$0")    # Skip tests
EOF
}

# Parse arguments
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
