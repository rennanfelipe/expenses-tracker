# Architectural Rules

## 1. Architectural Overview

The Expenses Tracker must follow the **Model–View–Controller (MVC)** architectural pattern.

The application is organized into three main components:

- **Model:** Owns the application's data, business rules, operations, and technical implementation.
- **View:** Responsible for presenting information and collecting user input.
- **Controller:** Receives user actions, coordinates operations, and connects the View to the Model.

The architecture must remain simple and pragmatic. Do not introduce additional architectural layers or patterns unless the existing structure demonstrates a concrete need.

### 1.1 Architectural Structure

```text
Expenses Tracker
│
├── Model
│   ├── Category
│   ├── Statement
│   └── Reports
│
├── Controller
│   ├── CLI
│   ├── TUI
│   └── Web
│
└── View
    ├── CLI
    ├── TUI
    └── Web
```

The application must support multiple presentation interfaces while sharing the same Model.

## 2. Model

The Model owns the application's functionality and implementation.

It is responsible for:

- Business rules and validations.
- Data structures and domain operations.
- Database connections and SQL operations.
- CSV parsing and data extraction.
- Data transformation and normalization.
- Filesystem operations.
- Report generation and calculations.
- Logging when required by model operations.

There must be no separate Infrastructure, Domain, or Application layer.

Technical operations must be implemented within their respective Model modules.

### 2.1 Model Modules

The initial Model must contain the following modules:

#### Category

Responsible for expense categories.

Responsibilities:

- Create, retrieve, update, and delete categories.
- Validate category information.
- Provide category information for expense classification.
- Manage category persistence.

#### Statement

Responsible for importing and managing credit card statements.

Responsibilities:

- Import CSV files.
- Parse Banco Inter and Bradesco statement formats.
- Validate imported data.
- Normalize transactions.
- Detect duplicate transactions.
- Persist transactions.
- Retrieve transaction information.
- Handle import errors and results.

Bank-specific parsing must remain isolated within the Statement module.

Adding support for another bank must not require modifications to unrelated Model modules.

#### Reports

Responsible for generating expense reports.

Responsibilities:

- Retrieve the necessary transaction information.
- Apply date and category filters.
- Calculate expense totals and aggregations.
- Prepare structured report data.
- Support the application's reporting requirements.

Report calculations must be centralized in the Model, not duplicated across presentation interfaces.

### 2.2 Model Rules

- Model modules must not depend on CLI, TUI, or Web presentation packages.
- Business rules must not be duplicated between modules or presentation interfaces.
- Database operations must remain inside the Model.
- Use parameterized SQL queries.
- Financial calculations must avoid floating-point inaccuracies. Use integer monetary units or an appropriate decimal representation.
- Database-specific implementation details must not leak into Controllers or Views.
- Keep modules cohesive and avoid unnecessary dependencies between them.
- Introduce additional Model modules only when a distinct responsibility justifies them.

## 3. Controller

Controllers receive user actions and coordinate the corresponding Model operations.

Each presentation interface must have its own Controllers.

### 3.1 Responsibilities

- Receive and interpret user requests.
- Validate presentation-level input.
- Invoke the appropriate Model operations.
- Coordinate multiple Model operations when necessary.
- Handle Model errors and determine the appropriate response.
- Provide data to the corresponding View.

### 3.2 Controller Rules

- Controllers must remain lightweight.
- Do not implement business rules inside Controllers.
- Do not execute SQL queries directly.
- Do not parse bank-specific CSV files.
- Do not perform report calculations.
- Do not duplicate functionality already implemented by the Model.
- Controllers must not depend on other presentation interfaces.

## 4. View

Views are responsible for presenting information and collecting user input.

The application initially supports three presentation interfaces:

- CLI: Command-line interface.
- TUI: Terminal user interface.
- Web: Web interface.

### 4.1 Responsibilities

- Display information returned by Controllers.
- Collect user input.
- Present validation errors and operation results.
- Format information for the corresponding interface.
- Provide an appropriate user experience for each interface.

### 4.2 View Rules

- Views must not contain business rules.
- Views must not execute SQL queries.
- Views must not access the database directly.
- Views must not parse bank statements.
- Views must not implement report calculations.
- Presentation-specific formatting is allowed.
- Views must not depend on other presentation interfaces.

## 5. Communication and Dependency Rules

The expected interaction flow is:

```text
User
  |
  v
View
  |
  v
Controller
  |
  v
Model
  |
  v
Controller
  |
  v
View
  |
  v
User
```

The following rules apply:

1. Views communicate with their respective Controllers.
2. Controllers invoke Model operations.
3. Model modules own their data operations and business behavior.
4. Model modules must not import Controller or View packages.
5. Controllers must not import other presentation interfaces.
6. Views must not bypass Controllers to access Model functionality directly.
7. Presentation interfaces must share the same Model functionality.

Avoid circular dependencies between Go packages.

## 6. Multiple Presentation Interfaces

CLI, TUI, and Web are different ways of interacting with the same application.

They must not become independent implementations of the business logic.

For example, importing a Banco Inter statement must use the same Model functionality regardless of whether the request originates from:

- A CLI command.
- A TUI action.
- A Web request.

The same principle applies to category management, transaction retrieval, and report generation.

Adding a new presentation interface must not require rewriting existing Model functionality.

## 7. Error Handling

Errors must be handled at the appropriate level.

- **Model:** Detects and returns business and technical errors.
- **Controller:** Interprets errors and determines the appropriate response.
- **View:** Presents errors in a human-readable format.

Rules:

- Never silently ignore errors.
- Preserve error context when propagating errors.
- Use Go's standard error handling mechanisms.
- Avoid exposing SQL errors, filesystem paths, or internal implementation details directly to end users.
- Do not use `panic` for expected application failures.

## 8. Configuration and Dependencies

- Application configuration must be explicit.
- Do not introduce global mutable state.
- Database credentials and other secrets must not be hardcoded.
- Use environment variables or an appropriate configuration mechanism.
- Initialize required dependencies explicitly.
- Avoid unnecessary dependency injection frameworks.

Use Go interfaces only when they provide a concrete benefit, such as decoupling a Model module from an implementation or enabling meaningful testing.

Do not create interfaces solely to satisfy an architectural pattern.

## 9. Testing

Testing must follow the project's `docs/ai/testing-guidelines.md`.

At minimum:

- Model business operations must have unit tests.
- CSV parsers must have tests using synthetic bank statement fixtures.
- Database operations must have integration tests against MariaDB.
- Controllers must have tests for request handling and error responses.
- Views should be tested when they contain meaningful presentation behavior.

Tests must not require real financial statements or production credentials.

## 10. Architectural Change Management

Codex must preserve the established MVC architecture.

Before introducing architectural changes:

1. Inspect the existing project structure.
2. Identify the responsibility affected by the proposed change.
3. Evaluate whether the existing MVC components can accommodate it.
4. Explain why the change is necessary.
5. Request approval before introducing additional architectural layers, frameworks, or major abstractions.

Do not introduce:

- A separate Infrastructure layer.
- Separate Domain and Application layers.
- Unnecessary repository, service, factory, or adapter abstractions.
- Additional presentation interfaces without an explicit requirement.
- Duplicate business logic across interfaces.

The architecture may evolve as the application grows, but complexity must be justified by actual requirements.

## 11. Definition of Done

An implementation is architecturally complete when:

- Changes respect MVC responsibilities.
- Model functionality is shared across presentation interfaces.
- Business rules remain inside the Model.
- Technical operations remain inside their respective Model modules.
- Controllers coordinate operations without implementing business logic.
- Views handle presentation without accessing persistence directly.
- No unnecessary architectural layers or abstractions are introduced.
- Relevant tests pass.
- Existing functionality remains intact.