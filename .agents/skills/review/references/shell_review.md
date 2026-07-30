---
title: Shell Script Code Review Best Practices
description: Shell script review for safety, quoting, error handling, and portability.
allowed-tools:
  # Shell linting
  - Bash(shellcheck:*)
  - Bash(shfmt:*)
  - Bash(checkbashisms:*)
---

# Shell Script Code Review Best Practices

## Critical Issues (Must Fix)

### Quote Variables

```bash
# ❌ BAD: Unquoted variables (word splitting, glob expansion)
for file in $files; do
    rm $file  # Fails on filenames with spaces!
done

# ❌ BAD: Unquoted command substitution
path=$(pwd)
cd $path  # Word splitting!

# ✅ GOOD: Always quote variables
for file in "$files"; do
    rm "$file"
done

path="$(pwd)"
cd "$path"
```

### Use Set Options

```bash
# ❌ BAD: No error handling
#!/bin/bash
cd /some/dir
rm -rf *  # Runs even if cd fails!

# ✅ GOOD: Strict mode
#!/bin/bash
set -euo pipefail

# -e: Exit on error
# -u: Error on undefined variables
# -o pipefail: Fail on pipe errors
```

### Avoid Eval

```bash
# ❌ BAD: Code injection risk
cmd="ls $user_input"
eval "$cmd"  # User can inject: "; rm -rf /"

# ✅ GOOD: Use arrays for commands
cmd=(ls "$user_input")
"${cmd[@]}"
```

### Check Command Success

```bash
# ❌ BAD: Assuming commands succeed
cd /some/directory
process_files

# ✅ GOOD: Check return codes
if ! cd /some/directory; then
    echo "Failed to cd" >&2
    exit 1
fi

# ✅ GOOD: Use || for inline handling
cd /some/directory || { echo "Failed" >&2; exit 1; }
```

## High Priority Issues

### Use [[ ]] Instead of [ ]

```bash
# ❌ BAD: [ ] has quirks
if [ $var = "value" ]; then  # Fails if var is empty!
if [ $a -gt $b ]; then

# ✅ GOOD: [[ ]] is safer
if [[ "$var" == "value" ]]; then
if (( a > b )); then  # Use (( )) for arithmetic
```

### ShellCheck Warnings

```bash
# ❌ BAD: Common shellcheck issues
cat file | grep pattern  # Useless use of cat
echo $var                # Missing quotes
test -f $file && source $file  # Missing quotes

# ✅ GOOD: Clean code
grep pattern file
echo "$var"
[[ -f "$file" ]] && source "$file"
```

### Temporary Files

```bash
# ❌ BAD: Predictable temp file (security risk)
tmpfile=/tmp/myapp.tmp
echo "$data" > "$tmpfile"

# ✅ GOOD: Use mktemp
tmpfile=$(mktemp)
trap 'rm -f "$tmpfile"' EXIT
echo "$data" > "$tmpfile"

# ✅ GOOD: Temp directory
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT
```

### Reading Files

```bash
# ❌ BAD: For loop on lines
for line in $(cat file); do  # Word splitting!
    echo "$line"
done

# ✅ GOOD: While read loop
while IFS= read -r line; do
    echo "$line"
done < file

# ✅ GOOD: Process substitution for commands
while IFS= read -r line; do
    echo "$line"
done < <(some_command)
```

## Code Quality

### Function Definitions

```bash
# ❌ BAD: No local variables
my_function() {
    result="computed value"  # Pollutes global namespace
}

# ✅ GOOD: Use local
my_function() {
    local result
    result="computed value"
    echo "$result"
}

# ✅ GOOD: Return values via stdout
get_name() {
    echo "computed_name"
}
name=$(get_name)
```

### Array Usage

```bash
# ✅ GOOD: Declare arrays properly
declare -a files=("file1" "file2" "file with spaces")

# ✅ GOOD: Iterate arrays
for file in "${files[@]}"; do
    process "$file"
done

# ✅ GOOD: Array length
echo "Count: ${#files[@]}"
```

### String Operations

```bash
# ✅ GOOD: Parameter expansion
filename="path/to/file.tar.gz"
echo "${filename##*/}"      # file.tar.gz (basename)
echo "${filename%.*}"       # path/to/file.tar (remove extension)
echo "${filename%%.*}"      # path/to/file (remove all extensions)
echo "${filename#*/}"       # to/file.tar.gz (remove first dir)
echo "${filename//\//\\}"   # Replace / with \
echo "${filename:-default}" # Default if empty
```

### Default Values

```bash
# ❌ BAD: Verbose check
if [ -z "$var" ]; then
    var="default"
fi

# ✅ GOOD: Parameter expansion
var="${var:-default}"      # Use default if unset/empty
var="${var:=default}"      # Set AND use default
var="${var:?Error msg}"    # Error if unset/empty
```

## Portability

### POSIX Compatibility

```bash
# ❌ BAD: Bash-only features in /bin/sh
#!/bin/sh
[[ "$a" == "$b" ]]  # [[ ]] is bash-only
((count++))         # (( )) is bash-only

# ✅ GOOD: POSIX for /bin/sh
#!/bin/sh
[ "$a" = "$b" ]
count=$((count + 1))

# ✅ GOOD: Use bash when needed
#!/usr/bin/env bash
[[ "$a" == "$b" ]]
```

### Command Availability

```bash
# ✅ GOOD: Check for required commands
check_requirements() {
    local missing=()
    for cmd in jq curl git; do
        if ! command -v "$cmd" &>/dev/null; then
            missing+=("$cmd")
        fi
    done
    if (( ${#missing[@]} > 0 )); then
        echo "Missing required commands: ${missing[*]}" >&2
        exit 1
    fi
}
```

## Error Handling

### Trap for Cleanup

```bash
#!/bin/bash
set -euo pipefail

cleanup() {
    rm -rf "$tmpdir"
    echo "Cleaned up" >&2
}

trap cleanup EXIT

tmpdir=$(mktemp -d)
# ... rest of script
```

### Error Messages

```bash
# ✅ GOOD: Send errors to stderr
error() {
    echo "ERROR: $*" >&2
}

die() {
    error "$@"
    exit 1
}

# Usage
[[ -f "$config" ]] || die "Config file not found: $config"
```

### Verbose Mode

```bash
#!/bin/bash
set -euo pipefail

VERBOSE="${VERBOSE:-false}"

log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" >&2
    fi
}

log "Starting process..."
```

## Testing

### Unit Testing with Bats

```bash
#!/usr/bin/env bats

setup() {
    tmpdir=$(mktemp -d)
    source ./my_script.sh
}

teardown() {
    rm -rf "$tmpdir"
}

@test "add_numbers returns correct sum" {
    result=$(add_numbers 2 3)
    [ "$result" -eq 5 ]
}

@test "validate_email rejects invalid email" {
    run validate_email "invalid"
    [ "$status" -eq 1 ]
    [[ "$output" == *"invalid email"* ]]
}
```

## Static Analysis Tools

| Tool | Purpose | Command |
|------|---------|---------|
| `shellcheck` | Linting | `shellcheck script.sh` |
| `shfmt` | Formatting | `shfmt -d script.sh` |
| `checkbashisms` | POSIX compliance | `checkbashisms script.sh` |
| `bats` | Unit testing | `bats tests/` |

## Script Template

```bash
#!/usr/bin/env bash
#
# Description: Brief description of script
# Usage: script.sh [options] <args>
#

set -euo pipefail

# Constants
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly SCRIPT_NAME="$(basename "$0")"

# Defaults
VERBOSE=false
DRY_RUN=false

usage() {
    cat <<EOF
Usage: $SCRIPT_NAME [options] <args>

Options:
    -h, --help      Show this help
    -v, --verbose   Verbose output
    -n, --dry-run   Dry run mode
EOF
}

log() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo "[INFO] $*" >&2
    fi
}

error() {
    echo "[ERROR] $*" >&2
}

die() {
    error "$@"
    exit 1
}

main() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help) usage; exit 0 ;;
            -v|--verbose) VERBOSE=true ;;
            -n|--dry-run) DRY_RUN=true ;;
            --) shift; break ;;
            -*) die "Unknown option: $1" ;;
            *) break ;;
        esac
        shift
    done

    log "Starting $SCRIPT_NAME"
    # Main logic here
}

main "$@"
```

## Security Checklist

- [ ] Use `set -euo pipefail`
- [ ] All variables quoted
- [ ] No eval with user input
- [ ] Temporary files created securely (mktemp)
- [ ] Cleanup with trap
- [ ] Input validation
- [ ] No hardcoded secrets
- [ ] Proper file permissions (umask)
- [ ] Use full paths for commands in cron/sudo
- [ ] Avoid running as root when not needed
