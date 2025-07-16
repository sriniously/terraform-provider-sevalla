# Contributing to Terraform Provider for Sevalla

We welcome contributions to the Terraform Provider for Sevalla! This document provides guidelines for contributing to this project.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes and commit them
5. Push your changes to your fork
6. Submit a pull request

## Development Setup

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [golangci-lint](https://golangci-lint.run/usage/install/) for code linting

### Building the Provider

```bash
# Clone the repository
git clone https://github.com/sriniously/terraform-provider-sevalla.git
cd terraform-provider-sevalla

# Build the provider
make build

# Install the provider locally
make install
```

### Running Tests

```bash
# Run unit tests
make test

# Run acceptance tests (requires Sevalla API credentials)
make testacc
```

### Code Style

This project follows Go coding conventions. Please ensure your code is formatted properly:

```bash
# Format code
make fmt

# Run linter
make lint
```

## Making Changes

### Adding New Resources

1. Create the resource file in `internal/provider/`
2. Implement the resource using the Terraform Plugin Framework
3. Add comprehensive tests
4. Update documentation

### Adding New Data Sources

1. Create the data source file in `internal/provider/`
2. Implement the data source using the Terraform Plugin Framework
3. Add comprehensive tests
4. Update documentation

### Testing

All changes should include appropriate tests:

- Unit tests for individual functions
- Acceptance tests for resources and data sources
- Integration tests where applicable

### Documentation

Documentation is auto-generated from the code. Use appropriate comments and examples in your resource/data source implementations.

## Code Review Process

1. All changes must be submitted via pull request
2. Pull requests must pass all automated tests
3. Code must be reviewed by at least one maintainer
4. All feedback must be addressed before merging

## Reporting Issues

When reporting issues, please include:

1. Terraform version
2. Provider version
3. Operating system
4. Configuration that reproduces the issue
5. Expected vs actual behavior

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.