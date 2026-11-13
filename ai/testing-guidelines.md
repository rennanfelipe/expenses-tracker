# Go Testing Guidelines

## 1. General Principles

The project MUST follow Go testing conventions and prioritize correctness, reliability, maintainability, and reproducibility.

The coding agent MUST follow these principles:

1. All significant business logic MUST have automated tests.
2. Tests MUST verify observable behavior rather than implementation details.
3. Tests MUST be deterministic and reproducible.
4. Tests MUST be independent of execution order.
5. Tests MUST cover expected behavior, invalid inputs, and relevant edge cases.
6. Tests MUST NOT depend on external services unless explicitly designed as integration tests.
7. Tests MUST NOT be modified merely to make an incorrect implementation pass.
8. Prefer the Go standard library testing package.
9. Avoid unnecessary testing frameworks and dependencies.
10. Every bug fix SHOULD include a regression test.

The primary objective is to ensure that application behavior remains correct as the system evolves.

## 2. Testing Tools

Use the standard Go testing package.

Recommended commands:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests without cached results
go test -count=1 ./...

# Run tests with race detection
go test -race ./...

# Generate coverage information
go test -cover ./...

# Generate a coverage profile
go test -coverprofile=coverage.out ./...

# Display coverage by function
go tool cover -func=coverage.out

# Generate an HTML coverage report
go tool cover -html=coverage.out
```

Use `go test -race ./...` when concurrency is involved or when validating shared mutable state.

## 3. Test Organization

Follow standard Go project conventions.

### Test Files

- Test files MUST use the `_test.go` suffix.
- Test files SHOULD reside alongside the source files they validate.
- Test function names MUST start with `Test`.
- Benchmark functions MUST start with `Benchmark`.
- Fuzzing functions MUST start with `Fuzz`.
- Example functions MUST start with `Example`.

Example:

```text
internal/
└── transaction/
    ├── transaction.go
    ├── transaction_test.go
    ├── category.go
    └── category_test.go
```

### Test Packages

Use the same package when testing internal implementation details.

Use the `_test` package suffix when testing the public behavior of a package from an external consumer's perspective.

Example:

```go
package transaction

// Tests internal implementation details.
```

```go
package transaction_test

// Tests the exported package API.
```

Do not expose private implementation details solely to make testing easier.

## 4. Unit Testing

Unit tests MUST validate individual functions, methods, and business rules independently.

Guidelines:

1. Test one logical behavior per test case.
2. Avoid dependencies on databases, networks, or external services.
3. Use explicit input and expected output values.
4. Test error conditions.
5. Test boundary conditions.
6. Prefer straightforward assertions.
7. Avoid excessive test setup.
8. Keep tests readable and maintainable.

Example:

```go
func TestCalculateTotal(t *testing.T) {
    tests := []struct {
        name  string
        items []int64
        want  int64
    }{
        {
            name:  "multiple amounts",
            items: []int64{1000, 2500, 500},
            want:  4000,
        },
        {
            name:  "empty amounts",
            items: []int64{},
            want:  0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculateTotal(tt.items)

            if got != tt.want {
                t.Errorf("CalculateTotal() = %d, want %d", got, tt.want)
            }
        })
    }
}
```

## 5. Table-Driven Tests

Prefer table-driven tests when validating multiple scenarios for the same behavior.

Guidelines:

- Each test case MUST have a descriptive name.
- Inputs and expected results MUST be explicit.
- Include valid, invalid, and boundary cases.
- Use `t.Run()` to execute individual scenarios.
- Avoid overly complicated test structures.

Table-driven tests SHOULD be used for:

- CSV parsing.
- Date parsing.
- Monetary value conversion.
- Input validation.
- Transaction categorization.
- Business rule evaluation.

Do not force table-driven tests when a test has only one meaningful scenario or requires substantially different setup.

## 6. Financial Data Testing

Financial calculations and data transformations require particular attention.

The following rules are mandatory:

1. Never use floating-point arithmetic in financial test expectations.
2. Use integer-based monetary values or the project's selected decimal representation.
3. Test positive, negative, and zero amounts where applicable.
4. Test decimal precision.
5. Test rounding behavior explicitly.
6. Test currency consistency.
7. Ensure aggregation operations do not lose precision.
8. Verify that refunds and reversed transactions follow established business rules.

Example:

```go
func TestSumTransactions(t *testing.T) {
    transactions := []Transaction{
        {AmountCents: 1050},
        {AmountCents: 2099},
        {AmountCents: -500},
    }

    got := SumTransactions(transactions)
    want := int64(2649)

    if got != want {
        t.Errorf("SumTransactions() = %d, want %d", got, want)
    }
}
```

Financial tests MUST verify that calculations preserve the original precision.

## 7. CSV and ETL Testing

CSV processing is a critical part of Expenses Tracker and MUST receive comprehensive test coverage.

### CSV Parsing

Test scenarios SHOULD include:

- Valid CSV files.
- Empty files.
- Missing required columns.
- Unexpected columns.
- Different column ordering.
- Empty fields.
- Invalid delimiters.
- Malformed records.
- Quoted fields.
- Fields containing delimiters.
- Fields containing line breaks.
- Different character encodings when supported.
- Unexpected header names.

### Bank-Specific Formats

Each supported financial institution MUST have dedicated test cases.

Initially:

- Banco Inter.
- Bradesco.

Tests MUST validate that each bank's CSV format is converted into the expected normalized representation.

Use representative CSV fixtures whenever possible.

### Date Parsing

Test scenarios MUST include:

- Valid dates.
- Invalid dates.
- Different supported date formats.
- Ambiguous date representations.
- Leap years.
- Month and year boundaries.
- Missing dates.
- Invalid calendar dates.

Dates MUST be parsed explicitly according to the bank-specific format.

Never rely on the host machine's locale to interpret financial statement dates.

### Monetary Value Parsing

Test scenarios SHOULD include:

- Brazilian decimal notation (1.234,56).
- Decimal notation with no thousands separator.
- Integer amounts.
- Negative amounts.
- Zero amounts.
- Empty amounts.
- Invalid numeric representations.
- Currency symbols when supported.

Tests MUST ensure that monetary conversions do not introduce precision loss.

### ETL Pipeline

Test the extraction, transformation, validation, and loading stages independently when practical.

Also include integration tests covering the complete pipeline.

Verify that:

1. Valid records are transformed correctly.
2. Invalid records are handled according to established business rules.
3. Original transaction information is preserved when required.
4. Duplicate detection behaves consistently.
5. Database failures do not produce misleading success results.
6. Failed imports do not leave partially persisted data when atomicity is required.
7. Repeated imports behave according to the defined idempotency rules.

## 8. Database Integration Testing

Database integration tests MUST validate actual persistence behavior.

Guidelines:

1. Use an isolated test database.
2. Never run integration tests against production databases.
3. Ensure database state is predictable.
4. Clean up test data after execution.
5. Test database constraints.
6. Test transactions and rollback behavior.
7. Test unique constraints and duplicate detection.
8. Test foreign key relationships.
9. Test repository error handling.
10. Test migrations against supported database versions when applicable.

Use the same database engine as the application in integration tests.

For example, if production uses MariaDB, integration tests SHOULD use MariaDB rather than substituting SQLite.

Tests SHOULD use independent database schemas or containers to prevent interference between concurrent test executions.

## 9. Mocking and Test Doubles

Use mocks, stubs, and fakes only when they provide a meaningful testing advantage.

Guidelines:

- Prefer real implementations for simple, deterministic dependencies.
- Use fakes when they provide a lightweight alternative to external dependencies.
- Use mocks when interaction verification is necessary.
- Avoid mocking every dependency.
- Avoid excessive expectations and implementation-specific assertions.
- Do not create interfaces solely to accommodate a mocking framework.

Prefer dependency injection through explicit constructors or function parameters.

The testing approach MUST NOT force unnecessary architectural complexity into production code.

## 10. Error and Edge Case Testing

All significant error-handling paths SHOULD have dedicated tests.

Consider:

- Nil inputs where permitted.
- Empty inputs.
- Missing configuration.
- Invalid file paths.
- Permission errors.
- Database connection failures.
- Constraint violations.
- Context cancellation.
- Timeouts.
- Unexpected external input.

Use `errors.Is()` and `errors.As()` when testing wrapped errors.

Example:

```go
func TestImportStatementInvalidFile(t *testing.T) {
    err := ImportStatement(context.Background(), "")

    if err == nil {
        t.Fatal("expected an error, got nil")
    }

    if !errors.Is(err, ErrInvalidFile) {
        t.Fatalf("expected ErrInvalidFile, got %v", err)
    }
}
```

Tests MUST validate meaningful error conditions rather than depending on complete error message strings unless the message itself is part of the public contract.

## 11. Test Fixtures and Test Data

Use fixtures for complex or representative input data.

Recommended structure:

```text
tests/
├── fixtures/
│   ├── inter/
│   │   ├── valid-statement.csv
│   │   ├── malformed-statement.csv
│   │   └── empty-statement.csv
│   │
│   └── bradesco/
│       ├── valid-statement.csv
│       ├── malformed-statement.csv
│       └── empty-statement.csv
│
└── integration/
```

Guidelines:

1. Fixtures MUST contain synthetic financial information.
2. Never commit real credit card statements.
3. Keep fixtures small and focused.
4. Include representative edge cases.
5. Avoid generating test data from external services.
6. Document unusual fixture characteristics when necessary.

Use `testdata/` directories for package-specific fixtures when appropriate.

## 12. Test Independence and Determinism

Tests MUST produce consistent results regardless of execution order or environment.

Guidelines:

- Avoid shared mutable state.
- Avoid depending on the current working directory.
- Avoid relying on system locale.
- Avoid relying on external network connectivity.
- Avoid hardcoded assumptions about the system clock.
- Use temporary directories for filesystem tests.
- Use `t.TempDir()` for temporary files.
- Use `t.Cleanup()` for resource cleanup.
- Use fixed dates and timestamps when testing date-dependent behavior.

When testing time-dependent logic, inject a clock abstraction only when necessary to make the behavior deterministic.

## 13. Fuzz Testing

Use Go fuzzing for input-processing functions where unexpected inputs could cause panics, incorrect parsing, or inconsistent behavior.

Fuzz testing is particularly valuable for:

- CSV parsers.
- Date parsers.
- Monetary value parsers.
- Input validators.
- Serialization and deserialization.

Example:

```go
func FuzzParseAmount(f *testing.F) {
    f.Add("125,50")
    f.Add("1.250,00")
    f.Add("0,00")

    f.Fuzz(func(t *testing.T, input string) {
        _, _ = ParseAmount(input)
    })
}
```

Fuzz targets MUST NOT panic for arbitrary input unless a panic is explicitly part of the documented contract.

Fuzz tests SHOULD verify meaningful invariants whenever possible, rather than merely checking that a function does not panic.

Run fuzz tests explicitly:

```bash
go test -fuzz=FuzzParseAmount -fuzztime=30s ./internal/money/
```

## 14. Benchmarks

Benchmarks SHOULD be introduced when performance is relevant or when a performance regression needs investigation.

Use Go's standard benchmarking facilities.

Potential benchmark targets:

- CSV parsing.
- Transaction normalization.
- Duplicate detection.
- Report aggregation.
- Large transaction queries.

Example:

```go
func BenchmarkParseAmount(b *testing.B) {
    for b.Loop() {
        _, err := ParseAmount("1.250,50")
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Benchmarks MUST:

- Use representative input data.
- Avoid unnecessary setup inside the measured operation.
- Report allocations when relevant.
- Be compared against a baseline before claiming performance improvements.

Do not introduce benchmarks for every function.

## 15. Coverage Expectations

Code coverage is an indicator of tested code, not proof of correctness.

Guidelines:

1. Prioritize critical business logic.
2. Prioritize financial calculations and transformations.
3. Prioritize CSV parsing and validation.
4. Prioritize duplicate detection.
5. Prioritize persistence and transaction handling.
6. Avoid writing meaningless tests solely to increase coverage.

The project SHOULD aim for at least 80% statement coverage for critical business packages, while evaluating coverage quality rather than relying exclusively on a percentage.

Coverage requirements MUST NOT encourage testing trivial implementation details.

## 16. Continuous Integration

Every change submitted for integration SHOULD pass automated verification.

Recommended CI checks:

```bash
gofmt -l .
go vet ./...
go test -count=1 ./...
go test -race ./...
```

Additional checks MAY include:

- Coverage reporting.
- Static analysis.
- Dependency vulnerability scanning.
- Integration tests.
- Database migration validation.

CI MUST fail when formatting, compilation, or required tests fail.

Tests MUST NOT be skipped or weakened merely to obtain a successful CI result.

## 17. AI Coding Agent Testing Responsibilities

The coding agent MUST follow these requirements when implementing or modifying code:

1. Inspect existing tests before changing implementation.
2. Identify the relevant testing strategy for the requested task.
3. Write tests for new business behavior.
4. Add regression tests for bug fixes.
5. Update tests when approved requirements change.
6. Run relevant tests after implementation.
7. Run the complete test suite when practical.
8. Investigate test failures instead of ignoring them.
9. Never claim that tests passed without executing them.
10. Report tests that could not be executed and explain why.
11. Avoid modifying unrelated tests.
12. Never delete tests simply because they fail.

When a test exposes a contradiction between implementation and requirements, the agent MUST report the conflict rather than arbitrarily changing either side.

## 18. Definition of Done

A task involving code changes is considered adequately tested when:

- Relevant unit tests have been implemented or updated.
- Integration tests have been added when external component interactions require verification.
- Edge cases and error conditions have been considered.
- Financial data integrity has been verified where applicable.
- Tests are deterministic and independent.
- Relevant test suites pass.
- No known test failures are left unexplained.
- Regression tests exist for resolved defects.
- Coverage of critical business behavior has been considered.

The coding agent MUST include a summary of the tests executed, their results, and any remaining testing limitations in its final response.