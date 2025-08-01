#!/bin/bash
set -e

echo "🚀 Setting up Terraform Provider ExtIP development environment..."

# Install development tools
echo "📦 Installing development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
go install github.com/securecode/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
go install github.com/goreleaser/goreleaser@latest

# Install pre-commit if available
if command -v pip >/dev/null 2>&1; then
    echo "🪝 Installing pre-commit..."
    pip install pre-commit
    pre-commit install || echo "Pre-commit hooks will be set up when you first commit"
fi

# Create necessary directories
echo "📁 Creating cache directories..."
mkdir -p /tmp/.terraform-plugin-cache
mkdir -p ~/.terraform.d/plugins

# Verify installations
echo "✅ Verifying installations..."
go version
terraform version
golangci-lint version
gosec --version || echo "gosec installed"
govulncheck -version || echo "govulncheck installed"
tfplugindocs version || echo "tfplugindocs installed"

echo "🎉 Development environment setup complete!"
echo ""
echo "Available commands:"
echo "  make help          - Show all available make targets"
echo "  make pre-commit    - Run pre-commit checks"
echo "  make test          - Run tests"
echo "  make ci-test       - Run full CI test suite"
echo "  make security-scan - Run security scans"
echo ""