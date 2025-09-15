# Disposable Cloud Environment<sup>TM</sup>

DCE manages temporary AWS accounts with automated budget and time controls, providing secure, isolated cloud environments for development, testing, and learning.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/)

> **DCE<sup>TM</sup> is your playground in the cloud**

## Table of Contents

- [Overview](#overview)
- [Use Cases](#use-cases)
- [Quick Start](#quick-start)
- [Requirements](#requirements)
- [Installation & Setup](#installation--setup)
- [Project Structure](#project-structure)
- [Development](#development)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Overview

DCE (Disposable Cloud Environment) is an AWS account management platform that provides temporary, budget-controlled access to AWS accounts. Organizations can provision isolated AWS environments on-demand with automatic cleanup and cost controls.

**Key Features:**
- Automated provisioning of AWS accounts from a managed pool
- Budget and time-based lease limits with automatic enforcement
- Secure isolation between users via IAM policies
- Automatic resource cleanup using aws-nuke after lease expiration
- Self-service account access through CLI and web interfaces
- Comprehensive usage tracking and cost reporting

The system operates by maintaining a pool of AWS child accounts, leasing them to users temporarily, then automatically cleaning and returning them to the pool for reuse.

## Use Cases

DCE is ideal for:
- **Developers** needing isolated AWS environments for testing and development
- **Students** learning AWS without affecting production resources
- **Training organizations** providing hands-on AWS experience
- **CI/CD pipelines** requiring clean AWS environments for integration testing
- **Experimentation** with cloud-native services and architectures
- **Cost-controlled** cloud resource exploration

## Quick Start

The easiest way to get started with DCE is with the DCE CLI:

### 1. Install DCE CLI

Download from [github.com/Optum/dce-cli](https://github.com/Optum/dce-cli/releases/latest):

```bash
# Download and install (macOS example)
curl -L -o dce_darwin_amd64.zip https://github.com/Optum/dce-cli/releases/latest/download/dce_darwin_amd64.zip
unzip dce_darwin_amd64.zip -d /usr/local/bin
```

### 2. Deploy DCE

```bash
# Initialize DCE configuration
dce init

# Deploy DCE infrastructure
export AWS_ACCESS_KEY_ID=XXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=XXXXXXXXXXXXXXXXXXXX
dce system deploy
```

### 3. Manage Accounts

```bash
# Add an account to the pool
dce accounts add \
    --account-id 123456789012 \
    --admin-role-arn arn:aws:iam::123456789012:role/OrganizationAccountAccessRole

# Lease an account
dce leases create \
    --principal-id jdoe@example.com \
    --budget-amount 100 --budget-currency USD

# Login to your leased account
dce leases login <lease-id>
```

## Requirements

### System Requirements
- **Go**: 1.22 or later
- **Terraform**: 0.12.x or later
- **AWS CLI**: v1 or v2 (with proper configuration)
- **GNU Make**: 3.x or later
- **GNU Bash**: For shell scripts

### AWS Requirements
- AWS account with administrative access
- IAM user with command line access
- Child AWS accounts for the account pool
- Cross-account IAM roles configured

### Supported Platforms
- macOS
- Linux
- Windows 10+ (with WSL)

## Installation & Setup

### For End Users

1. **Install DCE CLI** from the [latest release](https://github.com/Optum/dce-cli/releases/latest)
2. **Configure AWS credentials** with administrative access
3. **Deploy DCE** using `dce system deploy`
4. **Add child accounts** to the account pool

See the [Quickstart Guide](./docs/quickstart.md) for detailed instructions.

### For Contributors

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Optum/dce.git
   cd dce
   ```

2. **Install dependencies:**
   ```bash
   make setup
   ```

3. **Build the project:**
   ```bash
   make build
   ```

4. **Run tests:**
   ```bash
   make test
   ```

See [Development Guide](./docs/develop.md) for detailed setup instructions.

## Project Structure

```
dce/
├── cmd/                    # Application entry points
│   ├── lambda/            # AWS Lambda function handlers
│   │   ├── accounts/      # Account management API
│   │   ├── leases/        # Lease management API
│   │   ├── usage/         # Usage tracking and reporting
│   │   └── ...           # Other Lambda functions
│   └── codebuild/         # CodeBuild projects
│       └── reset/         # Account reset automation
├── pkg/                   # Shared business logic
│   ├── account/          # Account service and models
│   ├── lease/            # Lease service and business logic
│   ├── usage/            # Usage tracking and cost calculation
│   ├── db/               # DynamoDB data access layer
│   └── ...              # Other shared packages
├── modules/              # Terraform infrastructure
│   ├── dynamodb.tf      # Database table definitions
│   ├── gateway.tf       # API Gateway configuration
│   ├── *_lambda.tf      # Lambda function definitions
│   └── ...             # Other infrastructure modules
├── docs/                # Documentation
├── scripts/             # Build and deployment scripts
└── tests/              # Integration and acceptance tests
```

### Key Components

- **API Layer**: REST API handlers for account, lease, and usage management
- **Business Logic**: Core services for account lifecycle and lease management
- **Data Layer**: DynamoDB abstractions with consistent interfaces
- **Reset Pipeline**: Automated account cleanup using aws-nuke
- **Infrastructure**: Terraform modules for AWS resource provisioning

## Development

### Building

```bash
# Build all components
make build

# Build specific Lambda function
./scripts/build.sh <lambda-name>
```

### Testing

```bash
# Run unit tests and linting
make test

# Run functional tests (requires deployed DCE)
make test_functional
```

### Deployment

```bash
# Deploy to AWS (requires Terraform backend)
make deploy
```

### Code Standards

- **Language**: Go 1.22+ with standard formatting (`gofmt`)
- **Linting**: Uses `golangci-lint` with configuration in `.golangci.yml`
- **Testing**: Unit tests co-located with source code
- **Infrastructure**: Terraform with `tflint` validation

## Documentation

- **[Quickstart Guide](./docs/quickstart.md)** - Get up and running quickly
- **[Development Guide](./docs/develop.md)** - Local development setup
- **[Concepts](./docs/concepts.md)** - Core DCE concepts and terminology
- **[API Documentation](./docs/swagger.json)** - REST API reference
- **[Terraform Guide](./docs/terraform.md)** - Infrastructure deployment

### Building Documentation

```bash
# Generate HTML documentation
make documentation

# Serve documentation locally
make serve_docs
```

## Contributing

DCE was born at Optum, but belongs to the community. We welcome contributions!

### Getting Started

1. Read our [Contributor Guidelines](./CONTRIBUTING.md)
2. Review the [Code of Conduct](./CODE_OF_CONDUCT.md)
3. Sign the [Contributor License Agreement](./INDIVIDUAL_CONTRIBUTOR_LICENSE.md)

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make test` to verify
5. Submit a pull request

### Reporting Issues

- Use [GitHub Issues](https://github.com/Optum/dce/issues) for bug reports and feature requests
- Provide detailed information about your environment and use case
- Include steps to reproduce any issues

## License

This project is licensed under the [Apache License 2.0](./LICENSE).

---

_Originally written and maintained by contributors and [Devin](https://app.devin.ai), with updates from the core team._
