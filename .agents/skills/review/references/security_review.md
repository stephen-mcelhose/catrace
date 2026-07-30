---
title: Security Code Review Guide
description: Security review checklist covering OWASP Top 10, injection, authentication, secrets, and common vulnerabilities.
---

# Security Code Review Guide

## OWASP Top 10 (2021) Checklist

### 1. Broken Access Control (A01)

```
✓ Authorization checks on every endpoint
✓ Deny by default
✓ Rate limiting on sensitive operations
✓ CORS properly configured
✓ JWT validation includes expiration and signature
✓ No IDOR vulnerabilities (direct object references)
```

### 2. Cryptographic Failures (A02)

```
✓ No plaintext sensitive data storage
✓ Strong hashing for passwords (bcrypt, argon2)
✓ TLS 1.2+ for data in transit
✓ No hardcoded secrets
✓ Secure random for tokens
✓ No weak cryptographic algorithms (MD5, SHA1, DES)
```

### 3. Injection (A03)

```
✓ Parameterized queries (SQL)
✓ Sanitized command arguments
✓ No eval() with user input
✓ Escaped HTML output (XSS)
✓ LDAP injection prevention
✓ XML external entity prevention
```

### 4. Insecure Design (A04)

```
✓ Threat modeling performed
✓ Defense in depth
✓ Principle of least privilege
✓ Fail securely
✓ No security through obscurity
```

### 5. Security Misconfiguration (A05)

```
✓ Default credentials changed
✓ Unnecessary features disabled
✓ Error messages don't leak info
✓ Security headers set
✓ Cloud permissions minimal
```

### 6. Vulnerable Components (A06)

```
✓ Dependencies up to date
✓ Vulnerability scanning (npm audit, cargo audit, etc.)
✓ SBOM maintained
✓ No abandoned libraries
```

### 7. Authentication Failures (A07)

```
✓ Multi-factor authentication option
✓ Brute force protection
✓ Session invalidation on logout
✓ Secure password requirements
✓ No credential stuffing vulnerability
```

### 8. Data Integrity Failures (A08)

```
✓ Signed updates and packages
✓ CI/CD pipeline secured
✓ Dependency integrity verified
✓ No insecure deserialization
```

### 9. Logging Failures (A09)

```
✓ Security events logged
✓ Logs don't contain sensitive data
✓ Log injection prevented
✓ Tamper-evident logging
✓ Alerting configured
```

### 10. SSRF (A10)

```
✓ URL validation
✓ Allowlist for outbound requests
✓ No internal network access via user input
✓ Metadata service access blocked
```

## Input Validation

### Never Trust User Input

```
All data from:
- Request parameters (query, body, headers)
- URL paths
- Cookies
- File uploads
- WebSocket messages
- Environment variables (in some contexts)
- Database data (may have been user-supplied)
```

### Validation Strategies

```
1. Allowlist validation (preferred)
   - Define exactly what's allowed
   - Reject everything else

2. Type validation
   - Ensure correct data type
   - Validate against schema (JSON Schema, Zod, etc.)

3. Range/length checks
   - Minimum/maximum values
   - String length limits

4. Format validation
   - Regex for structured data
   - Use well-tested validators for email, URL, etc.

5. Encoding/decoding
   - Normalize before validation
   - Encode output for context (HTML, URL, SQL, etc.)
```

## Secrets Management

### ❌ BAD Patterns

```
- Hardcoded secrets in source code
- Secrets in git history
- Plaintext secrets in configs
- Secrets in Docker images
- Secrets in error messages/logs
- Same secrets across environments
```

### ✅ GOOD Patterns

```
- Environment variables (12-factor app)
- Secret management systems (Vault, AWS Secrets Manager, Google Secret Manager)
- Encrypted configuration files
- CI/CD secret injection
- Regular secret rotation
- Different secrets per environment
```

### Scanning for Secrets

```bash
# Tools to use
gitleaks detect --source .
trufflehog git file://.
detect-secrets scan
```

## Authentication Security

### Password Handling

```
✓ Hash with bcrypt, argon2, or scrypt
✓ Minimum complexity requirements
✓ Check against breach databases (HaveIBeenPwned)
✓ No maximum length restrictions (beyond very large)
✓ Rate limit failed attempts
✓ Secure password reset flow
```

### Session Management

```
✓ Secure, HttpOnly, SameSite cookies
✓ Regenerate session ID after login
✓ Absolute and idle timeouts
✓ Invalidate sessions on logout
✓ Invalidate all sessions on password change
```

### Token Security

```
✓ JWTs signed with strong algorithm (RS256 or EdDSA)
✓ Short expiration times
✓ Secure storage (not localStorage for auth tokens)
✓ Refresh token rotation
✓ Token revocation capability
```

## Authorization Patterns

### Check at Every Layer

```
1. API Gateway / Reverse Proxy
   - Rate limiting
   - Basic authentication

2. Controller / Handler
   - Role/permission checks
   - Ownership verification

3. Service Layer
   - Business rule authorization

4. Data Layer
   - Row-level security
   - Column-level permissions
```

### Common Vulnerabilities

```
❌ Checking permission once, using data multiple times
❌ Relying on client-side authorization only
❌ Exposing internal IDs without access checks
❌ Horizontal privilege escalation (accessing other users' data)
❌ Vertical privilege escalation (admin functions as regular user)
```

## API Security

### Headers to Set

```http
Content-Security-Policy: default-src 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-XSS-Protection: 0  (deprecated, rely on CSP)
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), camera=(), microphone=()
```

### CORS Configuration

```
✓ Specify allowed origins (never *)
✓ Limit allowed methods
✓ Limit allowed headers
✓ Set credentials requirement correctly
✓ Don't trust Origin header for authorization
```

### Rate Limiting

```
✓ Per-user/IP limits
✓ Endpoint-specific limits
✓ Exponential backoff for failures
✓ 429 Too Many Requests response
✓ Retry-After header
```

## File Upload Security

### Validation

```
✓ Check MIME type (not just extension)
✓ Validate magic bytes
✓ Limit file size
✓ Scan for malware
✓ Sanitize filenames
✓ Don't execute uploaded files
```

### Storage

```
✓ Store outside web root
✓ Generate random filenames
✓ Separate domain for user content
✓ Set Content-Disposition: attachment
✓ No directory listing
```

## Logging Best Practices

### What to Log

```
✓ Authentication events (login, logout, failures)
✓ Authorization failures
✓ Input validation failures
✓ Application errors
✓ Sensitive operations (password change, permission change)
✓ High-value transactions
```

### What NOT to Log

```
✗ Passwords (even hashed)
✗ Session tokens
✗ Credit card numbers
✗ Personal identification (SSN, etc.)
✗ Full stack traces in production
✗ Any data subject to compliance (GDPR PII)
```

### Log Injection Prevention

```
✓ Sanitize user input in logs
✓ Use structured logging (JSON)
✓ Encode special characters
✓ Don't log user-controlled newlines
```

## Dependency Security

### Regular Audits

```bash
# Node.js
npm audit
yarn audit
pnpm audit

# Python
pip-audit
safety check

# Go
govulncheck ./...

# Rust
cargo audit

# Java
mvn dependency-check:check
```

### Update Strategy

```
✓ Automated dependency updates (Dependabot, Renovate)
✓ CI checks for vulnerabilities
✓ Regular security reviews
✓ Lock file versioning
✓ Test after updates
```

## Secure Development Lifecycle

### Code Review Focus

```
1. Authentication/authorization logic
2. Data validation and sanitization
3. Cryptography usage
4. Error handling
5. Logging content
6. Third-party integrations
7. Configuration management
```

### Pre-Commit Checks

```yaml
# Example pre-commit config
repos:
  - repo: https://github.com/gitleaks/gitleaks
    hooks:
      - id: gitleaks
  - repo: https://github.com/Yelp/detect-secrets
    hooks:
      - id: detect-secrets
```

### CI Security Gates

```
✓ SAST (static analysis)
✓ Dependency scanning
✓ Secret scanning
✓ Container scanning
✓ DAST (dynamic analysis) in staging
✓ License compliance
```

## Language-Specific Security

### JavaScript/TypeScript

```
✗ eval(), new Function()
✗ innerHTML with user content
✗ document.write()
✗ Prototype pollution
✓ Content Security Policy
✓ Subresource Integrity
```

### Python

```
✗ pickle.loads() with untrusted data
✗ eval(), exec()
✗ os.system() with user input
✗ yaml.load() (use safe_load)
✓ Parameterized SQL queries
```

### Go

```
✗ Unsanitized template execution
✗ Command injection via exec
✗ Weak random (math/rand for crypto)
✓ crypto/rand for tokens
✓ html/template for escaping
```

### Shell

```
✗ Unquoted variables
✗ eval with user input
✗ Backticks for command substitution
✓ set -euo pipefail
✓ Proper quoting
```

## Quick Security Audit

```markdown
## Security Review Checklist

### Authentication
- [ ] Strong password hashing
- [ ] Brute force protection
- [ ] Secure session management
- [ ] MFA available for sensitive actions

### Authorization
- [ ] Every endpoint has access control
- [ ] Principle of least privilege
- [ ] No IDOR vulnerabilities

### Data Protection
- [ ] TLS everywhere
- [ ] No sensitive data in logs
- [ ] Secrets properly managed
- [ ] Encryption at rest for sensitive data

### Input/Output
- [ ] All input validated
- [ ] Output encoded for context
- [ ] Parameterized queries only

### Dependencies
- [ ] No known vulnerabilities
- [ ] Regular updates planned

### Infrastructure
- [ ] Security headers set
- [ ] CORS configured correctly
- [ ] Rate limiting enabled
```
