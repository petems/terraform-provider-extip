#!/bin/bash
set -euo pipefail

echo "🚀 Setting up Terraform Provider ExtIP development environment..."

# Check if Go is available
if ! command -v go >/dev/null 2>&1; then
    echo "❌ Go is not installed or not in PATH"
    exit 1
fi

echo "✅ Go version: $(go version)"

# Install development tools
echo "📦 Installing development tools..."
echo "  Installing golangci-lint..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
echo "  Installing gosec..."
go install github.com/securego/gosec/v2/cmd/gosec@latest
echo "  Installing govulncheck..."
go install golang.org/x/vuln/cmd/govulncheck@latest
echo "  Installing tfplugindocs..."
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
echo "  Installing goreleaser..."
go install github.com/goreleaser/goreleaser@latest

# Install pre-commit if available
if command -v pip >/dev/null 2>&1; then
    echo "🪝 Installing pre-commit..."
    pip install pre-commit
    if [ -f .pre-commit-config.yaml ]; then
        pre-commit install || echo "⚠️  Pre-commit hooks will be set up when you first commit"
    else
        echo "⚠️  No .pre-commit-config.yaml found, skipping hook installation"
    fi
else
    echo "⚠️  pip not found, skipping pre-commit installation"
fi

# Create necessary directories
echo "📁 Creating cache directories..."
mkdir -p /tmp/.terraform-plugin-cache
mkdir -p ~/.terraform.d/plugins

# Verify installations
echo "✅ Verifying installations..."
echo "  Go: $(go version)"
if command -v terraform >/dev/null 2>&1; then
    echo "  Terraform: $(terraform version | head -n1)"
else
    echo "  ⚠️  Terraform not found in PATH"
fi
if command -v golangci-lint >/dev/null 2>&1; then
    echo "  golangci-lint: $(golangci-lint version | head -n1)"
else
    echo "  ⚠️  golangci-lint not found in PATH (may need to restart shell)"
fi
echo "  gosec: $(gosec --version 2>/dev/null || echo "installed")"
echo "  govulncheck: $(govulncheck -version 2>/dev/null || echo "installed")"
echo "  tfplugindocs: $(tfplugindocs version 2>/dev/null || echo "installed")"

echo "🎉 Development environment setup complete!"
echo ""
echo "Available commands:"
echo "  make help          - Show all available make targets"
echo "  make pre-commit    - Run pre-commit checks"
echo "  make test          - Run tests"
echo "  make ci-test       - Run full CI test suite"
echo "  make security-scan - Run security scans"
echo ""