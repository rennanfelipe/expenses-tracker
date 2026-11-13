// Package expense provides CRUD operations for the expenses table.
package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound indicates that an expense with the requested ID does not exist.
var ErrNotFound = errors.New("expense not found")

// Expense represents a row in the expenses table.
type Expense struct {
	ID         uint64
	CreatedAt  time.Time
	Item       string
	CategoryID *uint32
	Price      string
	Date       time.Time
}

// Input contains the fields that can be supplied when creating or updating an expense.
type Input struct {
	Item       string
	CategoryID *uint32
	Price      string
	Date       time.Time
}

// Model performs operations on expenses using the supplied database connection.
type Model struct {
	db *sql.DB
}

// New creates an expense Model backed by db.
func New(db *sql.DB) *Model {
	return &Model{db: db}
}

// Create inserts an expense and returns the complete persisted record.
func (m *Model) Create(ctx context.Context, input Input) (Expense, error) {
	result, err := m.db.ExecContext(ctx, `
		INSERT INTO expenses (item, category_id, price, date)
		VALUES (?, ?, ?, ?)
	`, input.Item, input.CategoryID, input.Price, input.Date)
	if err != nil {
		return Expense{}, fmt.Errorf("create expense: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Expense{}, fmt.Errorf("get created expense ID: %w", err)
	}

	return m.Get(ctx, uint64(id))
}

// Get retrieves an expense by its primary key.
func (m *Model) Get(ctx context.Context, id uint64) (Expense, error) {
	var expense Expense

	err := m.db.QueryRowContext(ctx, `
		SELECT id, created_at, item, category_id, price, date
		FROM expenses
		WHERE id = ?
	`, id).Scan(
		&expense.ID,
		&expense.CreatedAt,
		&expense.Item,
		&expense.CategoryID,
		&expense.Price,
		&expense.Date,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Expense{}, ErrNotFound
	}
	if err != nil {
		return Expense{}, fmt.Errorf("get expense %d: %w", id, err)
	}

	return expense, nil
}

// GetAll retrieves all expenses.
func (m *Model) GetAll(ctx context.Context) ([]Expense, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT id, created_at, item, category_id, price, date
		FROM expenses
	`)
	if err != nil {
		return nil, fmt.Errorf("get expenses: %w", err)
	}
	defer rows.Close()

	expenses := make([]Expense, 0)
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(
			&expense.ID,
			&expense.CreatedAt,
			&expense.Item,
			&expense.CategoryID,
			&expense.Price,
			&expense.Date,
		); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}
		expenses = append(expenses, expense)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expenses: %w", err)
	}

	return expenses, nil
}

// Update changes the mutable fields of an expense and returns the persisted record.
func (m *Model) Update(ctx context.Context, id uint64, input Input) (Expense, error) {
	_, err := m.db.ExecContext(ctx, `
		UPDATE expenses
		SET item = ?, category_id = ?, price = ?, date = ?
		WHERE id = ?
	`, input.Item, input.CategoryID, input.Price, input.Date, id)
	if err != nil {
		return Expense{}, fmt.Errorf("update expense %d: %w", id, err)
	}

	return m.Get(ctx, id)
}

// Delete removes an expense by its primary key.
func (m *Model) Delete(ctx context.Context, id uint64) error {
	result, err := m.db.ExecContext(ctx, `DELETE FROM expenses WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete expense %d: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted expense count: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
