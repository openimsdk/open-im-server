# OpenIM End-to-End (E2E) Testing Module

## Overview

This repository contains the End-to-End (E2E) testing suite for OpenIM, a comprehensive instant messaging platform. The E2E tests are designed to simulate real-world usage scenarios to ensure that all components of the OpenIM system are functioning correctly in an integrated environment.

The tests cover various aspects of the system, including API endpoints, chat services, web interfaces, and RPC components, as well as performance and scalability under different load conditions.

## Directory Structure

```bash
❯ tree e2e
test/e2e/
├── conformance/             # Contains tests for verifying OpenIM API conformance
├── framework/               # Provides auxiliary code and libraries for building and running E2E tests
│   ├── config/              # Test configuration files and management
│   ├── ginkgowrapper/       # Functions wrapping the testing library for handling test failures and skips
│   └── helpers/             # Helper functions such as user creation, message sending, etc.
├── api/                     # End-to-end tests for OpenIM API
├── chat/                    # Tests for the business server (including login, registration, and other logic)
├── web/                     # Tests for the web frontend (login, registration, message sending and receiving)
├── rpc/                     # End-to-end tests for various RPC components
│   ├── auth/                # Tests for the authentication service
│   ├── conversation/        # Tests for conversation management
│   ├── friend/              # Tests for friend relationship management
│   ├── group/               # Tests for group management
│   └── message/             # Tests for message handling
├── scalability/             # Tests for the scalability of the OpenIM system
├── performance/             # Performance tests such as load testing and stress testing
└── upgrade/                 # Tests for compatibility and stability during OpenIM upgrades
```

The E2E tests are organized into the following directory structure:

- `conformance/`: Contains tests to verify the conformance of OpenIM API implementations.
- `framework/`: Provides helper code for constructing and running E2E tests using the Ginkgo framework.
  - `config/`: Manages test configurations and options.
  - `ginkgowrapper/`: Wrappers for Ginkgo's `Fail` and `Skip` functions to handle structured data panics.
  - `helpers/`: Utility functions for common test actions like user creation, message dispatching, etc.
- `api/`: E2E tests for the OpenIM API endpoints.
- `chat/`: Tests for the chat service, including authentication, session management, and messaging logic.
- `web/`: Tests for the web interface, including user interactions and information exchange.
- `rpc/`: E2E tests for each of the RPC components.
  - `auth/`: Tests for the authentication service.
  - `conversation/`: Tests for conversation management.
  - `friend/`: Tests for friend relationship management.
  - `group/`: Tests for group management.
  - `message/`: Tests for message handling.
- `scalability/`: Tests for the scalability of the OpenIM system.
- `performance/`: Performance tests, including load and stress tests.
- `upgrade/`: Tests for the upgrade process of OpenIM, ensuring compatibility and stability.

## Prerequisites

Since the deployment of OpenIM requires some components such as Mongo and Kafka, you should think a bit before using E2E tests

```bash
docker compose up -d
```

OR User [kubernetes deployment](https://github.com/openimsdk/helm-charts)

Before running the E2E tests, ensure that you have the following prerequisites installed:

- Docker
- Kubernetes
- Ginkgo test framework
- Go (version 1.19 or higher)

## Configuration

Test configurations can be customized via the `config/` directory. The configuration files are in YAML format and allow you to set parameters such as API endpoints, user credentials, and test data.

## Running the Tests

To run a single test or set of tests, you'll need the [Ginkgo](https://github.com/onsi/ginkgo) tool installed on your machine:

```
ginkgo --help
  --focus value
    	If set, ginkgo will only run specs that match this regular expression. Can be specified multiple times, values are ORed.
```

To run the entire suite of E2E tests, use the following command:

```sh
ginkgo -v --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
```

You can also run a specific test or group of tests by specifying the path to the test directory:

```bash
ginkgo -v ./test/e2e/chat
```

Or you can use Makefile to run the tests:

```bash
make test-e2e
```

## Test Development

To contribute to the E2E tests:

1. Clone the repository and navigate to the `test/e2e/` directory.
2. Create a new test file or modify an existing test to cover a new scenario.
3. Write test cases using the Ginkgo BDD style, ensuring that they are clear and descriptive.
4. Run the tests locally to ensure they pass.
5. Submit a pull request with your changes.

Please refer to the `CONTRIBUTING.md` file for more detailed instructions on contributing to the test suite.


## Reporting Issues

If you encounter any issues while running the E2E tests, please open an issue on the GitHub repository with the following information:

Open issue: https://github.com/openimsdk/open-im-server/issues/new/choose, choose "Failing Test" template.

+ A clear and concise description of the issue.
+ Steps to reproduce the behavior.
+ Relevant logs and test output.
+ Any other context that could be helpful in troubleshooting.


## Continuous Integration (CI)

The E2E test suite is integrated with CI, which runs the tests automatically on each code commit. The results are reported back to the pull request or commit to provide immediate feedback on the impact of the changes.


## Contact

For any queries or assistance, please reach out to the OpenIM development team at [support@openim.com](mailto:support@openim.com).