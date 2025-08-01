# GitHub Actions Workflows

This directory contains the GitHub Actions workflows for the terraform-provider-extip project. These workflows have been modernized and optimized for better performance, security, and maintainability.

## Workflows Overview

### 🔧 Core Workflows

#### `ci.yml` - Continuous Integration
- **Triggers**: Push to main/master, Pull requests
- **Purpose**: Primary CI pipeline with testing, building, and validation
- **Features**:
  - Multi-version Go testing (1.21, 1.22, 1.23)
  - Linting with golangci-lint (only on latest Go version)
  - Unit tests with coverage reporting
  - Acceptance tests (on main branch only)
  - Cross-platform builds (Linux, Darwin, Windows, FreeBSD)
  - Security scanning with gosec and govulncheck
  - Code validation and formatting checks

#### `release.yml` - Release Management
- **Triggers**: Git tags matching `v*`
- **Purpose**: Automated releases using GoReleaser
- **Features**:
  - Multi-platform binary builds
  - GPG signing of release artifacts
  - Automated changelog generation
  - Release validation and notification
  - Integration with GitHub Releases

### 🔄 Maintenance Workflows

*Note: Dependency management is handled by Dependabot (configured separately)*

### 📚 Quality Assurance Workflows

#### `docs.yml` - Documentation Validation
- **Triggers**: Changes to docs/, examples/, or Markdown files
- **Purpose**: Validate and maintain documentation quality
- **Features**:
  - Terraform provider documentation validation
  - Example configuration validation
  - Markdown linting
  - Auto-generated documentation sync check

#### `registry.yml` - Provider Registry Publishing
- **Triggers**: Git tags, Manual workflow dispatch
- **Purpose**: Prepare and validate releases for Terraform Registry
- **Features**:
  - Registry metadata validation
  - Release artifact verification
  - GPG signature validation
  - Publishing readiness checklist

### 🔧 Reusable Components

#### `reusable-setup.yml` - Shared Setup Logic
- **Purpose**: Reusable workflow for Go environment setup
- **Features**:
  - Configurable Go version
  - Optional dependency installation
  - Consistent caching strategy

## Usage

### Manual Triggers
```bash
# Trigger full test suite
gh workflow run ci.yml
```

### Branch Protection
Recommended branch protection rules:
- Require status checks to pass before merging
- Require branches to be up to date before merging
- Require pull request reviews before merging
- Require conversation resolution before merging

### Secrets Required
- `GPG_PRIVATE_KEY` - GPG private key for release signing
- `PASSPHRASE` - GPG key passphrase
- `GITHUB_TOKEN` - Automatically provided by GitHub

## Performance

### CI Pipeline
- **Fast feedback**: ~5-10 minutes for basic tests
- **Full validation**: ~15-20 minutes for complete pipeline
- **Parallel execution**: Jobs run in parallel where possible
- **Caching**: Go modules and build cache enabled

### Release Pipeline
- **Release time**: ~10-15 minutes
- **Artifact validation**: Automatic testing of release binaries
- **Cross-platform**: Supports Linux, macOS, Windows

## Monitoring

### Code Coverage
- Uploaded to Codecov on every run
- Coverage reports available in PRs
- Coverage badges in README

### Security
- SARIF reports uploaded to GitHub Security
- Vulnerability alerts for high/critical issues
- Dependency vulnerability scanning

### Notifications
- Release success/failure notifications
- Security issue alerts
- Build status notifications

## Troubleshooting

### Common Issues

1. **Build failures on Windows**
   - Check for Windows-specific code paths
   - Verify cross-platform compatibility

2. **Security scan failures**
   - Review gosec and govulncheck reports
   - Address high/critical vulnerabilities

3. **Dependabot PR failures**
   - Review automated dependency update PRs
   - Check for breaking changes in dependencies

### Debugging
- Enable debug logging in workflows
- Check workflow run logs
- Use `actions/upload-artifact` for debugging files 