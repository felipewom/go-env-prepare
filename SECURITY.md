# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest  | ✅        |
| < 0.0.1 | ❌        |

## Reporting a Vulnerability

**Please do not open a public issue for security vulnerabilities.**

Report security issues by emailing the maintainer directly or by using [GitHub private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing/privately-reporting-a-security-vulnerability) for this repository.

Include:
- A description of the vulnerability and potential impact.
- Steps to reproduce or proof-of-concept (if safe to share).
- Any suggested mitigations.

You will receive a response within **7 days**. If the issue is confirmed, a fix will be released as soon as possible and credit will be given in the changelog.

## Dependency Updates

Dependencies are monitored via [Dependabot](.github/dependabot.yml) and updated weekly. To check for vulnerabilities locally:

```bash
go list -json -m all | nancy sleuth
# or
govulncheck ./...
```
