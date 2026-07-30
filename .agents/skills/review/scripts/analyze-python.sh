#!/usr/bin/env bash
#
# Python Code Analysis Script
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

    local python_install="brew install python3"
    if [[ "$PLATFORM" == "linux" ]]; then
        python_install="apt install python3 python3-pip"
    fi

    if ! check_tool "python3" "$python_install"; then
        missing=1
    else
        success "python3 $(python3 --version 2>&1 | awk '{print $2}')"
    fi

    # Determine installer strategy
    local INSTALLER=""
    if command -v uv &>/dev/null; then
        success "uv $(uv --version 2>&1 | awk '{print $2}')"
        INSTALLER="uv pip install"
    elif [[ -n "${VIRTUAL_ENV:-}" ]]; then
        log "Virtual environment detected: $VIRTUAL_ENV"
        INSTALLER="pip install"
    elif command -v pipx &>/dev/null; then
         # We can't use pipx to install libraries for analysis usually, unless they are CLI tools
         # But ruff, mypy, bandit are CLI tools.
         log "pipx detected"
         INSTALLER="pipx install"
    else
        warn "No virtual environment or 'uv' detected."
        warn "Installing packages globally with 'pip' is risky."
        warn "Recommendation: Use 'uv' or create a venv first."
        INSTALLER="pip install --user"
    fi

    # Helper to install if missing (but respect safety)
    install_if_missing() {
        local tool="$1"
        if ! command -v "$tool" &>/dev/null; then
            if [[ -n "$INSTALLER" ]]; then
                # Only auto-install if we are in a safe environment (venv or uv) OR user explicitly allows?
                # For now, we WARN and show command if unsafe.
                if [[ "$INSTALLER" == "pip install --user" ]]; then
                    warn "$tool missing. To install: $INSTALLER $tool"
                    return 1
                else
                    log "Installing $tool with $INSTALLER..."
                    $INSTALLER "$tool" 2>/dev/null || warn "Failed to install $tool"
                fi
            else
                warn "$tool missing."
                return 1
            fi
        fi
        
        # Check again
        if command -v "$tool" &>/dev/null; then
            success "$tool installed"
            return 0
        else
            return 1
        fi
    }

    # Check/Install tools
    install_if_missing "ruff" || log "ruff missing (optional but recommended)"
    install_if_missing "mypy" || log "mypy missing (optional)"
    install_if_missing "bandit" || log "bandit missing (optional)"
    install_if_missing "radon" || log "radon missing (optional)"

    return $missing
}

# Find Python source directory
find_python_src() {
    if [[ -d "src" ]]; then
        echo "src"
    elif [[ -d "lib" ]]; then
        echo "lib"
    else
        echo "."
    fi
}

# Run ruff (linter + formatter)
run_ruff() {
    header "Running Ruff (Linter)"

    if ! command -v ruff &>/dev/null; then
        warn "ruff not available, skipping"
        return 0
    fi

    local src_dir
    src_dir=$(find_python_src)

    local fix_flag=""
    if [[ "$FIX" == "true" ]]; then
        fix_flag="--fix"
    fi

    if ruff check $fix_flag "$src_dir" 2>&1; then
        success "ruff check passed"
    else
        warn "ruff found issues"
        return 1
    fi

    header "Running Ruff (Formatter Check)"
    if [[ "$FIX" == "true" ]]; then
        ruff format "$src_dir" 2>&1
        success "ruff format applied"
    else
        if ruff format --check "$src_dir" 2>&1; then
            success "ruff format check passed"
        else
            warn "ruff format found issues. Run 'ruff format .' to fix"
            return 1
        fi
    fi

    return 0
}

# Run mypy (type checker)
run_mypy() {
    header "Running MyPy (Type Checker)"

    if ! command -v mypy &>/dev/null; then
        warn "mypy not available, skipping"
        return 0
    fi

    local src_dir
    src_dir=$(find_python_src)

    # Use loose settings if no mypy.ini/pyproject.toml config
    local config_flag=""
    if [[ ! -f "mypy.ini" ]] && ! grep -q "\[tool.mypy\]" pyproject.toml 2>/dev/null; then
        config_flag="--ignore-missing-imports --no-error-summary"
    fi

    if mypy $config_flag "$src_dir" 2>&1; then
        success "mypy passed"
        return 0
    else
        warn "mypy found type issues"
        return 1
    fi
}

# Run bandit (security scanner)
run_bandit() {
    header "Running Bandit (Security Scanner)"

    if ! command -v bandit &>/dev/null; then
        warn "bandit not available, skipping"
        return 0
    fi

    local src_dir
    src_dir=$(find_python_src)

    if bandit -r -q "$src_dir" 2>&1; then
        success "bandit passed (no security issues)"
        return 0
    else
        warn "bandit found potential security issues"
        return 1
    fi
}

# Check complexity with radon
run_radon() {
    header "Checking Cyclomatic Complexity (radon)"

    if ! command -v radon &>/dev/null; then
        warn "radon not available, skipping"
        return 0
    fi

    local src_dir
    src_dir=$(find_python_src)
    local threshold="${RADON_THRESHOLD:-C}"  # A-F scale, C = 11-20

    local result
    result=$(radon cc -s --min "$threshold" "$src_dir" 2>/dev/null || true)

    if [[ -z "$result" ]]; then
        success "No functions with complexity >= $threshold"
        return 0
    else
        warn "Functions with high complexity (>= $threshold):"
        echo "$result"
        return 1
    fi
}

# Check maintainability index
run_radon_mi() {
    header "Checking Maintainability Index (radon)"

    if ! command -v radon &>/dev/null; then
        warn "radon not available, skipping"
        return 0
    fi

    local src_dir
    src_dir=$(find_python_src)

    local result
    result=$(radon mi -s "$src_dir" 2>/dev/null | grep -E "^.+ - [CD]" || true)

    if [[ -z "$result" ]]; then
        success "All files have good maintainability (A or B)"
        return 0
    else
        warn "Files with low maintainability (C or D):"
        echo "$result"
        return 1
    fi
}

# Run tests with pytest
run_tests() {
    header "Running Tests"

    if ! command -v pytest &>/dev/null; then
        if [[ -f "pytest.ini" ]] || [[ -f "pyproject.toml" ]]; then
            warn "pytest not installed but config found. Run: pip install pytest"
        else
            log "pytest not available, skipping tests"
        fi
        return 0
    fi

    if pytest --tb=short -q 2>&1; then
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

    local src_dir
    src_dir=$(find_python_src)

    local count
    count=$(grep -rn --include="*.py" -E "(TODO|FIXME|XXX|HACK)" "$src_dir" 2>/dev/null | wc -l | tr -d ' ')

    if [[ "$count" -gt 0 ]]; then
        warn "Found $count TODO/FIXME comments:"
        grep -rn --include="*.py" -E "(TODO|FIXME|XXX|HACK)" "$src_dir" 2>/dev/null | head -20
        if [[ "$count" -gt 20 ]]; then
            echo "... and $((count - 20)) more"
        fi
    else
        success "No TODO/FIXME comments found"
    fi
}

# Check dependencies for vulnerabilities
check_deps() {
    header "Checking Dependencies for Vulnerabilities"

    if command -v pip-audit &>/dev/null; then
        if pip-audit 2>&1; then
            success "pip-audit: no vulnerabilities found"
        else
            warn "pip-audit found vulnerabilities"
        fi
    elif command -v safety &>/dev/null; then
        if safety check 2>&1; then
            success "safety: no vulnerabilities found"
        else
            warn "safety found vulnerabilities"
        fi
    else
        warn "Neither pip-audit nor safety installed. Run: pip install pip-audit"
    fi
}

# Main execution
main() {
    local exit_code=0

    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║             Python Code Analysis Report                       ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════════╝${NC}"

    # Check for Python project indicators
    if [[ ! -f "pyproject.toml" ]] && [[ ! -f "setup.py" ]] && [[ ! -f "requirements.txt" ]]; then
        warn "No Python project files found (pyproject.toml, setup.py, requirements.txt)"
    fi

    ensure_tools || exit_code=1

    run_ruff || exit_code=1
    run_mypy || exit_code=1
    run_bandit || exit_code=1
    run_radon || exit_code=1
    run_radon_mi || exit_code=1

    if [[ "${SKIP_TESTS:-false}" != "true" ]]; then
        run_tests || exit_code=1
    fi

    check_todos
    check_deps

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
    RADON_THRESHOLD=C      Set complexity threshold (A-F, default: C)

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
