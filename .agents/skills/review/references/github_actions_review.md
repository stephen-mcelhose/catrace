---
title: GitHub Actions Code Review Best Practices
description: GitHub Actions workflow security, permissions, secrets handling, and CI/CD best practices.
---

# GitHub Actions Code Review Best Practices

## Check csgda-platform-actions First

Before creating new workflows or actions, **always search csgda-platform-actions** for existing reusable actions:

```
ReferenceSearch(query="build docker image", repo="csgda-platform-actions")
ReferenceSearch(query="deploy to GKE", repo="csgda-platform-actions")
ReferenceSearch(query="semantic release", repo="csgda-platform-actions")
ReferenceSearch(query="terraform apply", repo="csgda-platform-actions")
```

### Common Reusable Actions

| Need | Search Query |
|------|-------------|
| Docker builds | `"docker build" OR "container"` |
| Kubernetes deploy | `"deploy" OR "GKE" OR "kubernetes"` |
| Terraform | `"terraform" OR "infrastructure"` |
| Go builds | `"go build" OR "golang"` |
| Python builds | `"python" OR "pip" OR "poetry"` |
| Release automation | `"release" OR "semantic" OR "changelog"` |
| Security scanning | `"security" OR "trivy" OR "snyk"` |
| Code quality | `"lint" OR "test" OR "coverage"` |

### Pattern Alignment Questions

- ❓ Is there an existing action in csgda-platform-actions for this?
- ❓ Does this workflow follow our standard patterns?
- ❓ Should this be contributed back as a reusable action?
- ❓ Are we using the organization's standard action versions?

---

## Critical Issues (Must Fix)

### Pin Action Versions

```yaml
# ❌ BAD: Unpinned versions (security risk, breaking changes)
- uses: actions/checkout@main
- uses: actions/setup-node@v4

# ❌ BAD: Using latest tag
- uses: docker/build-push-action@latest

# ✅ GOOD: Pin to specific SHA
- uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1
- uses: actions/setup-node@60edb5dd545a775178f52524783378180af0d1f8 # v4.0.2

# ✅ GOOD: For internal actions, use version tags
- uses: bayer-int/csgda-platform-actions/go-build@v1.2.3
```

### Secrets Security

```yaml
# ❌ BAD: Hardcoded secrets
env:
  API_KEY: "sk-1234567890abcdef"

# ❌ BAD: Secrets in logs
- run: echo "Using key ${{ secrets.API_KEY }}"

# ❌ BAD: Secrets passed to untrusted action
- uses: random-org/unknown-action@main
  with:
    token: ${{ secrets.GITHUB_TOKEN }}

# ✅ GOOD: Use GitHub secrets
env:
  API_KEY: ${{ secrets.API_KEY }}

# ✅ GOOD: Mask secrets in logs
- run: |
    echo "::add-mask::${{ secrets.API_KEY }}"
    ./deploy.sh
```

### Injection Prevention

```yaml
# ❌ BAD: Direct injection of untrusted input (command injection!)
- run: echo "Processing ${{ github.event.issue.title }}"

# ❌ BAD: Pull request title/body can contain malicious code
- run: |
    TITLE="${{ github.event.pull_request.title }}"
    echo "$TITLE"

# ✅ GOOD: Use environment variables (shell escapes properly)
- run: echo "Processing $TITLE"
  env:
    TITLE: ${{ github.event.issue.title }}

# ✅ GOOD: Use intermediate file
- run: |
    cat > issue.json << 'EOF'
    ${{ toJSON(github.event.issue) }}
    EOF
    jq '.title' issue.json
```

### GITHUB_TOKEN Permissions

```yaml
# ❌ BAD: Default (overly permissive in some cases)
# No permissions block

# ✅ GOOD: Explicit minimum permissions
permissions:
  contents: read
  pull-requests: write

# ✅ GOOD: Read-only for most jobs
jobs:
  build:
    permissions:
      contents: read

  deploy:
    permissions:
      contents: read
      id-token: write  # For OIDC
```

## High Priority Issues

### Use Workload Identity (OIDC)

```yaml
# ❌ BAD: Long-lived service account keys
- uses: google-github-actions/auth@v2
  with:
    credentials_json: ${{ secrets.GCP_SA_KEY }}

# ✅ GOOD: Workload Identity Federation (no stored secrets!)
permissions:
  id-token: write
  contents: read

- uses: google-github-actions/auth@v2
  with:
    workload_identity_provider: 'projects/123/locations/global/workloadIdentityPools/github/providers/github'
    service_account: 'github-actions@project.iam.gserviceaccount.com'
```

### Conditional Execution

```yaml
# ❌ BAD: Running expensive jobs unnecessarily
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - run: ./deploy.sh  # Runs on every push!

# ✅ GOOD: Use conditions
jobs:
  deploy:
    if: github.ref == 'refs/heads/main' && github.event_name == 'push'
    runs-on: ubuntu-latest

# ✅ GOOD: Path filters for monorepos
on:
  push:
    paths:
      - 'services/api/**'
      - '.github/workflows/api.yml'

# ✅ GOOD: Skip CI for docs-only changes
on:
  push:
    paths-ignore:
      - '**.md'
      - 'docs/**'
```

### Concurrency Control

```yaml
# ❌ BAD: Multiple deploys can run simultaneously
on:
  push:
    branches: [main]

# ✅ GOOD: Cancel in-progress runs for same ref
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

# ✅ GOOD: For deployments, wait (don't cancel)
concurrency:
  group: deploy-${{ github.ref }}
  cancel-in-progress: false
```

### Error Handling

```yaml
# ❌ BAD: No error handling
- run: |
    ./might-fail.sh
    ./depends-on-above.sh

# ✅ GOOD: Use set -e and proper error handling
- run: |
    set -euo pipefail
    ./might-fail.sh
    ./depends-on-above.sh

# ✅ GOOD: Use continue-on-error for non-critical steps
- run: ./optional-check.sh
  continue-on-error: true

# ✅ GOOD: Cleanup on failure
- run: ./deploy.sh
  id: deploy
- run: ./rollback.sh
  if: failure() && steps.deploy.outcome == 'failure'
```

## Code Quality

### Reusable Workflows

```yaml
# ❌ BAD: Duplicated workflow logic across repos

# ✅ GOOD: Use reusable workflows
# In shared repo:
# .github/workflows/go-build.yml
name: Go Build (Reusable)
on:
  workflow_call:
    inputs:
      go-version:
        type: string
        default: '1.22'
    secrets:
      token:
        required: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
      - run: go build ./...

# In consuming repo:
jobs:
  build:
    uses: bayer-int/csgda-platform-actions/.github/workflows/go-build.yml@v1
    with:
      go-version: '1.22'
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

### Composite Actions

```yaml
# ✅ GOOD: Create composite actions for repeated steps
# action.yml
name: 'Setup Project'
description: 'Setup Go and dependencies'
inputs:
  go-version:
    default: '1.22'
runs:
  using: 'composite'
  steps:
    - uses: actions/setup-go@v5
      with:
        go-version: ${{ inputs.go-version }}
    - run: go mod download
      shell: bash
    - run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
      shell: bash
```

### Job Dependencies

```yaml
# ✅ GOOD: Explicit job dependencies
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: go test ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - run: golangci-lint run

  build:
    needs: [test, lint]  # Wait for both
    runs-on: ubuntu-latest
    steps:
      - run: go build ./...

  deploy:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
```

### Matrix Builds

```yaml
# ✅ GOOD: Test across multiple versions
jobs:
  test:
    strategy:
      matrix:
        go-version: ['1.21', '1.22']
        os: [ubuntu-latest, macos-latest]
      fail-fast: false  # Don't cancel other jobs on failure
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      - run: go test ./...
```

## Performance

### Caching

```yaml
# ✅ GOOD: Cache dependencies
- uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    cache: true  # Built-in caching

# ✅ GOOD: Custom cache for other tools
- uses: actions/cache@v4
  with:
    path: |
      ~/.cache/golangci-lint
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-

# ✅ GOOD: Cache Docker layers
- uses: docker/build-push-action@v5
  with:
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

### Artifact Management

```yaml
# ✅ GOOD: Share artifacts between jobs
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: go build -o app ./...
      - uses: actions/upload-artifact@v4
        with:
          name: app-binary
          path: app
          retention-days: 1  # Short retention for CI artifacts

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/download-artifact@v4
        with:
          name: app-binary
```

### Parallel Jobs

```yaml
# ❌ BAD: Sequential steps that could be parallel
jobs:
  ci:
    steps:
      - run: go test ./...
      - run: golangci-lint run
      - run: go build ./...

# ✅ GOOD: Parallel jobs
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - run: go test ./...

  lint:
    runs-on: ubuntu-latest
    steps:
      - run: golangci-lint run

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
```

## Security Best Practices

### Pull Request Workflows

```yaml
# ❌ BAD: pull_request_target with checkout of PR code
on:
  pull_request_target:
jobs:
  build:
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{ github.event.pull_request.head.sha }}  # DANGEROUS!
      - run: ./build.sh  # Runs untrusted code with secrets access!

# ✅ GOOD: Use pull_request (no secrets, fork-safe)
on:
  pull_request:
jobs:
  build:
    steps:
      - uses: actions/checkout@v4
      - run: ./build.sh

# ✅ GOOD: If you need pull_request_target, don't checkout PR code
on:
  pull_request_target:
jobs:
  label:
    steps:
      - uses: actions/github-script@v7
        with:
          script: |
            // Safe: only uses API, doesn't run PR code
            github.rest.issues.addLabels(...)
```

### Environment Protection

```yaml
# ✅ GOOD: Use environments for sensitive deployments
jobs:
  deploy-prod:
    environment:
      name: production
      url: https://app.example.com
    runs-on: ubuntu-latest
    steps:
      - run: ./deploy.sh

# Environment settings (in GitHub UI):
# - Required reviewers
# - Wait timer
# - Deployment branches (only main)
# - Environment secrets
```

### Supply Chain Security

```yaml
# ✅ GOOD: Generate and verify SBOMs
- uses: anchore/sbom-action@v0
  with:
    artifact-name: sbom.spdx.json

# ✅ GOOD: Sign artifacts with Sigstore
- uses: sigstore/cosign-installer@v3
- run: cosign sign-blob --yes artifact.tar.gz

# ✅ GOOD: Use artifact attestations (GitHub native)
- uses: actions/attest-build-provenance@v1
  with:
    subject-path: 'dist/*'
```

## Standard Workflow Template

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

permissions:
  contents: read
  pull-requests: write

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1

      - uses: actions/setup-go@0c52d547c9bc32b1aa3301fd7a9cb496313a4491 # v5.0.0
        with:
          go-version-file: 'go.mod'
          cache: true

      - name: Run tests
        run: |
          set -euo pipefail
          go test -race -coverprofile=coverage.txt ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: coverage.txt

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1

      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1

      - uses: actions/setup-go@0c52d547c9bc32b1aa3301fd7a9cb496313a4491 # v5.0.0
        with:
          go-version-file: 'go.mod'
          cache: true

      - name: Build
        run: go build -o app ./...

      - uses: actions/upload-artifact@v4
        with:
          name: app
          path: app
```

## Security Checklist

- [ ] All third-party actions pinned to SHA
- [ ] GITHUB_TOKEN has minimum required permissions
- [ ] No secrets in workflow logs
- [ ] Using OIDC for cloud authentication (no stored keys)
- [ ] Pull request workflows don't access secrets
- [ ] Injection-safe handling of untrusted input
- [ ] Environment protection for production deploys
- [ ] Concurrency controls to prevent race conditions
- [ ] Using reusable actions from csgda-platform-actions
