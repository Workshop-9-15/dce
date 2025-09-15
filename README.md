# Disposable Cloud Environment<sup>TM</sup>

DCE is an AWS account management platform that provides temporary, budget-controlled access to AWS accounts for safe cloud exploration and development.

> **DCE<sup>TM</sup> is your playground in the cloud**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Optum/dce)](https://goreportcard.com/report/github.com/Optum/dce)

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Requirements](#requirements)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Overview

DCE helps organizations provide secure, temporary AWS environments for development, testing, training, and experimentation. It manages a pool of AWS child accounts, automatically provisioning them to users with defined budgets and time limits, then cleaning and returning them for reuse.

**Common use cases:**
- **Development & Testing**: Isolated environments for application development and CI/CD pipelines
- **Training & Education**: Hands-on AWS learning without production risk
- **Experimentation**: Safe exploration of new AWS services and architectures
- **Compliance**: Controlled access with automatic cleanup and audit trails

## Key Features

- **🔒 Secure Isolation**: Each user gets a dedicated AWS account with IAM-enforced boundaries
- **💰 Budget Controls**: Automatic enforcement of spending limits with real-time monitoring
- **⏰ Time-based Leases**: Configurable lease durations with automatic expiration
- **🧹 Automatic Cleanup**: Complete resource removal using aws-nuke after lease ends
- **📊 Usage Tracking**: Comprehensive cost and resource utilization reporting
- **🚀 Self-Service**: CLI and API access for seamless account provisioning
- **⚡ Serverless Architecture**: Built on AWS Lambda, DynamoDB, and SQS for scalability

## Architecture

DCE follows a serverless, event-driven architecture:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   DCE CLI/API   │───▶│   API Gateway    │───▶│  Lambda Functions│
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
┌─────────────────┐    ┌──────────────────┐             │
│  Account Pool   │◀───│    DynamoDB      │◀────────────┘
│ (Child Accounts)│    │   (State Store)  │
└─────────────────┘    └──────────────────┘
         │                       │
         │              ┌──────────────────┐
         └──────────────▶│   Reset Queue    │
                        │ (SQS + CodeBuild)│
                        └──────────────────┘
```

**Core Components:**
- **API Layer**: REST endpoints for account, lease, and usage management
- **Account Service**: Manages AWS account lifecycle and pool operations
- **Lease Service**: Handles user account assignments with budget/time controls
- **Reset Pipeline**: Automated cleanup using aws-nuke via CodeBuild
- **Usage Tracking**: Real-time cost monitoring and reporting

## Quick Start

### Prerequisites

- AWS CLI configured with appropriate permissions
- Go 1.13+ (for development)
- Terraform 0.12+ (for infrastructure deployment)

### Installation

1. **Install the DCE CLI:**
   ```bash
   # Download from releases or build from source
   go install github.com/Optum/dce-cli@latest
   ```

2. **Deploy DCE Infrastructure:**
   ```bash
   # Initialize and deploy with Terraform
   cd modules
   terraform init
   terraform apply
   
   # Deploy application code
   cd ..
   make deploy
   ```

3. **Add AWS Accounts to Pool:**
   ```bash
   dce accounts add \
       --account-id 123456789012 \
       --admin-role-arn arn:aws:iam::123456789012:role/OrganizationAccountAccessRole
   ```

### Basic Usage

```bash
# Create a lease
dce leases create \
    --principal-id user@example.com \
    --budget-amount 100 \
    --budget-currency USD

# Login to your leased account
dce leases login <lease-id>

# Check lease status
dce leases list

# End lease early (optional)
dce leases end <lease-id>
```

## Requirements

**Runtime Requirements:**
- AWS Organization with multiple child accounts
- IAM permissions for cross-account access
- AWS CLI v1.x or v2.x (see [AWS CLI 2 configuration](./docs/develop.md#configuring-aws-cli-2))

**Development Requirements:**
- Go 1.13+
- Terraform 0.12+
- GNU Make 3.x+
- AWS account for testing deployments

**Supported Platforms:**
- macOS
- Linux
- Windows 10 with WSL

## Project Structure

```
dce/
├── cmd/                    # Application entry points
│   ├── lambda/            # AWS Lambda function handlers
│   │   ├── accounts/      # Account management API
│   │   ├── leases/        # Lease management API
│   │   ├── usage/         # Usage tracking and reporting
│   │   └── ...           # Additional Lambda functions
│   └── codebuild/         # CodeBuild projects for account reset
├── pkg/                   # Shared business logic
│   ├── account/          # Account service and models
│   ├── lease/            # Lease service and business rules
│   ├── db/               # DynamoDB data access layer
│   └── ...              # Additional packages
├── modules/              # Terraform infrastructure definitions
│   ├── *.tf             # AWS resource configurations
│   └── fixtures/        # IAM policies and templates
├── tests/               # Test suites
│   └── acceptance/      # Integration tests
├── docs/                # Documentation
└── scripts/             # Build and deployment automation
```

## Documentation

- **[Quick Start Guide](./docs/quickstart.md)** - Get up and running quickly
- **[Development Guide](./docs/develop.md)** - Local development setup and workflows
- **[Deployment Guide](./docs/terraform.md)** - Infrastructure deployment with Terraform
- **[API Documentation](./docs/api-auth.md)** - REST API reference and authentication
- **[Concepts](./docs/concepts.md)** - Core concepts and terminology
- **[IAM Policies](./docs/iam-policies.md)** - Security model and permissions
- **[AWS Nuke Support](./docs/awsnuke-support.md)** - Account cleanup configuration

## Development

### Local Setup

```bash
# Clone and setup
git clone https://github.com/Optum/dce.git
cd dce
make setup

# Build application
make build

# Run tests
make test

# Run linting
make lint
```

### Testing

```bash
# Unit tests
make test

# Functional tests (requires deployed DCE)
make test_functional
```

See the [Development Guide](./docs/develop.md) for detailed setup instructions and development workflows.

## Contributing

DCE was born at Optum but belongs to the community. We welcome contributions!

**Getting Started:**
1. Read our [Contributor Guidelines](./CONTRIBUTING.md)
2. Review our [Code of Conduct](./CODE_OF_CONDUCT.md)
3. Check out [open issues](https://github.com/Optum/dce/issues)
4. Submit a [pull request](https://github.com/Optum/dce/pulls)

**Development Process:**
- Fork the repository and create a feature branch
- Follow Go coding standards (gofmt, golint)
- Add tests for new functionality
- Ensure all tests pass before submitting PR

## License

[Apache License v2.0](./LICENSE)

---

*Originally written and maintained by contributors and [Devin](https://app.devin.ai), with updates from the core team.*
