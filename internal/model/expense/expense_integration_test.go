//go:build integration

package expense_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rennanfelipe/expenses-tracker/internal/model/expense"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("EXPENSES_TEST_DSN")
	if dsn == "" {
		t.Skip("EXPENSES_TEST_DSN is not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping test database: %v", err)
	}

	return db
}

func testContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func createCategory(t *testing.T, db *sql.DB) uint32 {
	t.Helper()

	name := fmt.Sprintf("integration-category-%d", time.Now().UnixNano())
	result, err := db.ExecContext(testContext(t),
		"INSERT INTO categories (name, description) VALUES (?, ?)",
		name, "synthetic integration-test category")
	if err != nil {
		t.Fatalf("create test category: %v", err)
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get test category ID: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM categories WHERE id = ?", categoryID)
	})

	return uint32(categoryID)
}

func inputForTest(categoryID *uint32) expense.Input {
	return expense.Input{
		Item:       "Synthetic integration expense",
		CategoryID: categoryID,
		Price:      "1234567890.12",
		Date:       time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
	}
}

func assertExpenseFields(t *testing.T, got expense.Expense, want expense.Input) {
	t.Helper()

	if got.ID == 0 {
		t.Error("created expense ID = 0, want an auto-generated ID")
	}
	if got.CreatedAt.IsZero() {
		t.Error("created_at is zero, want the MariaDB default timestamp")
	}
	if got.Item != want.Item {
		t.Errorf("item = %q, want %q", got.Item, want.Item)
	}
	if got.Price != want.Price {
		t.Errorf("price = %q, want exact DECIMAL value %q", got.Price, want.Price)
	}
	if got.Date.Year() != want.Date.Year() || got.Date.Month() != want.Date.Month() || got.Date.Day() != want.Date.Day() {
		t.Errorf("date = %v, want calendar date %v", got.Date, want.Date)
	}
	if want.CategoryID == nil {
		if got.CategoryID != nil {
			t.Errorf("category_id = %v, want NULL", *got.CategoryID)
		}
		return
	}
	if got.CategoryID == nil || *got.CategoryID != *want.CategoryID {
		t.Errorf("category_id = %v, want %d", got.CategoryID, *want.CategoryID)
	}
}

func TestModelCreateAndGetPreserveSchemaValues(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)

	categoryID := createCategory(t, db)
	wantInput := inputForTest(&categoryID)
	created, err := model.Create(ctx, wantInput)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Cleanup(func() { _ = model.Delete(context.Background(), created.ID) })

	assertExpenseFields(t, created, wantInput)

	got, err := model.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != created {
		t.Errorf("Get() = %#v, want %#v", got, created)
	}
}

func TestModelCreateAllowsNullCategory(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)
	wantInput := inputForTest(nil)

	created, err := model.Create(ctx, wantInput)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Cleanup(func() { _ = model.Delete(context.Background(), created.ID) })

	assertExpenseFields(t, created, wantInput)
}

func TestModelGetAllIncludesCreatedExpenses(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)

	wantIDs := make(map[uint64]struct{}, 2)
	for _, item := range []string{"Synthetic expense A", "Synthetic expense B"} {
		created, err := model.Create(ctx, expense.Input{
			Item:  item,
			Price: "0.01",
			Date:  time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("Create(%q) error = %v", item, err)
		}
		wantIDs[created.ID] = struct{}{}
		id := created.ID
		t.Cleanup(func() { _ = model.Delete(context.Background(), id) })
	}

	got, err := model.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	seen := make(map[uint64]expense.Expense)
	for _, item := range got {
		if _, wanted := wantIDs[item.ID]; wanted {
			seen[item.ID] = item
		}
	}
	if len(seen) != len(wantIDs) {
		t.Fatalf("GetAll() returned %d of %d created expenses", len(seen), len(wantIDs))
	}
}

func TestModelUpdateChangesMutableFieldsOnly(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)

	created, err := model.Create(ctx, inputForTest(nil))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Cleanup(func() { _ = model.Delete(context.Background(), created.ID) })

	wantInput := expense.Input{
		Item:  "Updated synthetic expense",
		Price: "0.10",
		Date:  time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC),
	}
	updated, err := model.Update(ctx, created.ID, wantInput)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	assertExpenseFields(t, updated, wantInput)
	if updated.ID != created.ID {
		t.Errorf("updated ID = %d, want unchanged ID %d", updated.ID, created.ID)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("updated created_at = %v, want unchanged value %v", updated.CreatedAt, created.CreatedAt)
	}
}

func TestModelDeleteRemovesExpense(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)

	created, err := model.Create(ctx, inputForTest(nil))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := model.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := model.Get(ctx, created.ID); !errors.Is(err, expense.ErrNotFound) {
		t.Fatalf("Get() after Delete() error = %v, want ErrNotFound", err)
	}
}

func TestModelMissingRecordsReturnNotFound(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)
	missingID := uint64(18446744073709551615)

	if _, err := model.Get(ctx, missingID); !errors.Is(err, expense.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
	if _, err := model.Update(ctx, missingID, inputForTest(nil)); !errors.Is(err, expense.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
	if err := model.Delete(ctx, missingID); !errors.Is(err, expense.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestModelCreateRejectsInvalidData(t *testing.T) {
	db := openTestDB(t)
	model := expense.New(db)
	ctx := testContext(t)
	invalidCategory := uint32(4294967295)

	tests := []struct {
		name  string
		input expense.Input
	}{
		{
			name: "foreign key violation",
			input: expense.Input{
				Item:       "Synthetic invalid category",
				CategoryID: &invalidCategory,
				Price:      "1.00",
				Date:       time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "item exceeds varchar limit",
			input: expense.Input{
				Item:  strings.Repeat("x", 256),
				Price: "1.00",
				Date:  time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "price exceeds decimal precision",
			input: expense.Input{
				Item:  "Synthetic oversized price",
				Price: "10000000000.00",
				Date:  time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := model.Create(ctx, tt.input); err == nil {
				t.Fatal("Create() error = nil, want database constraint error")
			}
		})
	}
}

func TestModelPropagatesClosedDatabaseErrors(t *testing.T) {
	dsn := os.Getenv("EXPENSES_TEST_DSN")
	if dsn == "" {
		t.Skip("EXPENSES_TEST_DSN is not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close test database: %v", err)
	}

	model := expense.New(db)
	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "create",
			call: func() error {
				_, err := model.Create(ctx, inputForTest(nil))
				return err
			},
		},
		{
			name: "get",
			call: func() error {
				_, err := model.Get(ctx, 1)
				return err
			},
		},
		{
			name: "get all",
			call: func() error {
				_, err := model.GetAll(ctx)
				return err
			},
		},
		{
			name: "update",
			call: func() error {
				_, err := model.Update(ctx, 1, inputForTest(nil))
				return err
			},
		},
		{
			name: "delete",
			call: func() error {
				return model.Delete(ctx, 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil {
				t.Fatal("operation error = nil, want closed-database error")
			}
		})
	}
}

