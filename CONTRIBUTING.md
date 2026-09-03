# Contributing

Thank you for considering contributing to the OPNsense Go API client!

This document provides guidelines and instructions for developing, testing, and contributing changes to the project.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
  - [1. Fork and Clone](#1-fork-and-clone)
  - [2. Install Dependencies](#2-install-dependencies)
- [Making Changes](#making-changes)
  - [Branch Naming](#branch-naming)
  - [Code Organization](#code-organization)
  - [Adding an OPNsense Module](#adding-an-opnsense-module)
  - [Generating Code](#generating-code)
  - [Code Formatting](#code-formatting)
- [Testing](#testing)
  - [Acceptance Tests](#acceptance-tests)
  - [Writing Acceptance Tests](#writing-acceptance-tests)
- [Code Standards](#code-standards)
  - [Go Best Practices](#go-best-practices)
  - [API and Schema Conventions](#api-and-schema-conventions)
  - [Generated Code](#generated-code)
- [Pull Request Process](#pull-request-process)
  - [Before Submitting](#before-submitting)
  - [PR Guidelines](#pr-guidelines)
  - [Review Process](#review-process)
- [Project Structure](#project-structure)
- [Getting Help](#getting-help)

## Prerequisites

Before you begin, ensure that you have:

- **Go 1.24 or later**
- An **OPNsense instance** for acceptance tests
- **Make** (optional, but recommended)

You should also have a basic understanding of:

- Go programming
- OPNsense configuration
- The OPNsense API
- YAML

## Development Setup

### 1. Fork and Clone

Fork this repository on GitHub and clone your fork:

```bash
git clone https://github.com/YOUR-USERNAME/opnsense-go.git
cd opnsense-go
````

### 2. Install Dependencies

Download the Go dependencies:

```bash
go mod download
```

## Making Changes

### Branch Naming

Create a feature branch from `main`:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:

* `feature/` - New features or OPNsense resources
* `fix/` - Bug fixes
* `docs/` - Documentation changes
* `refactor/` - Code refactoring
* `test/` - Test changes

## Code Organization

Each OPNsense module is represented by a Go package under `pkg/`.

The generated resource files contain the Go types and conversion functions required to communicate with the OPNsense API.

### Adding an OPNsense Module

Adding support for a new OPNsense module generally consists of the following steps:

1. Add the module schema to `schema/`.
2. Create the corresponding package under `pkg/`.
3. Add a `generate.go` file.
4. Generate the API client code.
5. Add acceptance tests where applicable.
6. Run formatting and tests.

For example, a new `service` module contains:

```text
schema/
└── [service].yml

pkg/
└── [service]/
    ├── controller.go
    ├── generate.go
    ├── [resource1].go
    ├── [resource1]_test.go
    ├── [resource2].go
    └── [resource2]_test.go
```

Acceptance tests should normally be provided for every resource that can be created, modified, or deleted through the API.

### Schema Definition

OPNsense modules are primarily defined through YAML schemas in the `schema/` directory.

The schema describes:

* The OPNsense module
* The API endpoints
* The resources provided by the module
* The attributes of each resource
* The corresponding Go types

The generator uses these schemas to create the Go API implementation.

The general structure is:

```text
OPNsense
   │
   ├── XML model
   ├── API endpoints
   │
   ▼
schema/<module>.yml
   │
   ▼
internal/generate/api/main.go
   │
   ▼
pkg/<module>/
   ├── controller.go
   └── <resource>.go
```

The OPNsense MVC model is the primary source for resource and field definitions, including field types, defaults, validation constraints, and available options.

The OPNsense API should additionally be used to verify the actual request and response format.

The generated client should be tested against a real OPNsense instance whenever possible, as the API behavior can differ from what might be inferred from the MVC model alone.

When adding or changing a resource, the schema should be updated first whenever possible instead of manually modifying generated source files.

### Schema Types

The schema type determines the Go type generated for an attribute.

Common schema types include:

| Schema Type | Generated Go Type | Usage |
|-------------|-------------------|-------|
| `string` | `string` | Text, numbers, booleans and other scalar API values |
| `SelectedMap` | `api.SelectedMap` | A single selected OPNsense option |
| `SelectedMapList` | `api.SelectedMapList` | Multiple selected OPNsense options |

The schema type must be chosen based on how the value is represented by the OPNsense API, not only on the field type used in the OPNsense MVC model.

OPNsense often represents values such as integers and booleans as strings in its API. These values should therefore normally use the `string` schema type unless the generator provides a more specific type.

For OPNsense `OptionField` values, use `SelectedMap` when the field allows a single selection.

For fields with `Multiple="Y"`, use `SelectedMapList`.


### `generate.go`

Each generated package contains a `generate.go` file with the `go:generate` directive.

For example:

```go
//go:generate go run ../../internal/generate/api/main.go -controller [service]

package [service]
```

The `generate.go` file should contain **only the package declaration and generate directives**.

Do not add manually maintained implementation code to this file.

### Generating Code

After changing a schema, regenerate the affected package:

```bash
make all [service]
```

Generated files must not be edited manually.

Any changes to generated code should normally be made by changing the corresponding schema or generator instead.

### Code Formatting

Format the code before committing:

```bash
make fmt
```

or:

```bash
gofmt -w .
```

The CI pipeline checks formatting.

## Testing

### Acceptance Tests

Acceptance tests communicate with a real OPNsense instance through its API.

They are used to verify the complete API lifecycle of resources, including:

* Creating resources
* Reading resources
* Updating resources
* Deleting resources
* Reconfiguring OPNsense where required
* Verifying values returned by OPNsense

Acceptance tests should clean up resources they create, including resources created as prerequisites for other resources.

### Acceptance Test Environment

Set the following environment variables before running acceptance tests:

```bash
export OPNSENSE_URI="https://your-opnsense-host.example.com"
export OPNSENSE_API_KEY="your-api-key"
export OPNSENSE_API_SECRET="your-api-secret"
export OPNSENSE_ALLOW_INSECURE="true"
```

`OPNSENSE_ALLOW_INSECURE=true` should only be used when required, for example with a self-signed certificate.

### Running Acceptance Tests

Run all tests:

```bash
make testacc
```

Run acceptance tests for a specific package:

```bash
make testacc PKG=[service]
```

Run a specific acceptance test:

```bash
make testacc PKG=[service] TEST=Test[resource]
```

Alternatively, use `go test` directly:

```bash
OPNSENSE_URI="..." \
OPNSENSE_API_KEY="..." \
OPNSENSE_API_SECRET="..." \
go test -v -p 1 ./pkg/[service]/...
```

The `-p 1` option is important when running acceptance tests against a shared OPNsense instance.

Tests modify real OPNsense configuration and should therefore be executed serially unless the test environment explicitly supports concurrent execution.

### Writing Acceptance Tests

Acceptance tests should normally verify the complete lifecycle:

```text
Create
  ↓
Read
  ↓
Verify
  ↓
Update
  ↓
Read
  ↓
Verify
  ↓
Delete
  ↓
Verify cleanup
```

Tests should use realistic values accepted by the OPNsense API.

When a resource depends on another resource, create the dependency first and use the identifier returned by OPNsense.

For example, a resource using a `ModelRelationField` may require the UUID of another resource rather than its display name or numeric identifier.

## Code Standards

### Go Best Practices

Follow standard Go conventions:

* Use `gofmt`.
* Keep functions focused and small.
* Return errors instead of silently ignoring failures.
* Use descriptive names.
* Avoid unnecessary abstractions.
* Add tests for new functionality.

### API and Schema Conventions

The schema should reflect the OPNsense API as closely as possible.

Pay particular attention to the distinction between:

* API keys
* Display values
* Selected values
* Resource identifiers

For example, an OPNsense option may expose:

```text
Key:   Kbit
Value: kbit/s
```

When using `api.SelectedMap`, the API key must be used:

```go
api.SelectedMap("Kbit")
```

rather than the display value:

```go
api.SelectedMap("kbit/s")
```

The same principle applies to other `OptionField`, `InterfaceField`, and relation values.

### SelectedMap and SelectedMapList

Use `api.SelectedMap` for a single selected value.

Use `api.SelectedMapList` for fields that support multiple selected values.

For example:

```go
api.SelectedMap("tcp")
```

and:

```go
api.SelectedMapList{"tcp", "udp"}
```

The exact keys must match the values accepted by the corresponding OPNsense API field.

### Generated Code

Files generated by the API generator contain:

```text
// Code generated by internal/generate/api/main.go; DO NOT EDIT.
```

Do not modify generated files manually.

Instead:

1. Modify the YAML schema.
2. Regenerate the affected package.
3. Review the generated changes.
4. Add or update tests where required.

If the generator itself needs to change, modify the generator and regenerate the affected packages.

## Pull Request Process

### Before Submitting

Before opening a pull request, ensure that:

* [ ] Code is formatted with `gofmt`
* [ ] Acceptance tests pass where applicable
* [ ] Generated code is up to date
* [ ] No generated files were modified manually
* [ ] No unnecessary files are included
* [ ] Commit messages are clear and descriptive
* [ ] Documentation has been updated where necessary

### PR Guidelines

1. **Keep PRs focused**
   A pull request should ideally contain one feature, resource, or fix.

2. **Write a clear title**
   The title should briefly describe the change.

3. **Describe the change**
   Explain what was changed and why.

4. **Include tests**
   New functionality should include appropriate tests.

5. **Link related issues**
   Reference related issues where applicable.

### Review Process

Pull requests generally go through the following process:

1. Automated CI checks
2. Maintainer review
3. Requested changes, if necessary
4. Approval
5. Merge

## Project Structure

The main directories are:

```text
.
├── internal/
│   └── generate/
│       └── api/
├── pkg/
│   ├── <module>/
│   └── ...
├── schema/
│   ├── <module>.yml
│   └── ...
├── Makefile
├── go.mod
└── README.md
```

### Key Concepts

**`schema/`**

Contains the YAML definitions used by the API generator.

**`internal/generate/`**

Contains the code generator used to create the typed API implementations.

**`pkg/`**

Contains the public Go packages representing OPNsense modules and resources.

**Generated resource files**

Contain the generated resource types, endpoint configuration, and schema conversion logic.

**Acceptance tests**

Exercise the generated API against a real OPNsense instance.

## Getting Help

Before opening an issue, check the existing documentation and issues.

When reporting a problem, provide:

* A clear description of the problem
* Steps to reproduce it
* Relevant code or configuration
* OPNsense version
* `opnsense-go` version or commit
* Relevant API responses where possible

---

Thank you for contributing to `opnsense-go`!

Your contributions help improve programmatic access to OPNsense for everyone.
