// Package sqlmock is a minimal subset of the API of github.com/DATA-DOG/go-sqlmock.
//
// It is implemented here only to make this training exercise self-contained in
// an offline environment. It supports just enough functionality for the unit
// tests in this repository.
package sqlmock

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Sqlmock is the mock controller returned from New().
type Sqlmock interface {
	ExpectQuery(query string) *ExpectedQuery
	ExpectationsWereMet() error
}

// ExpectedQuery represents a single expected query.
type ExpectedQuery struct {
	query string
	rows  *Rows
	err   error
}

// WillReturnRows defines the rows to be returned for this expected query.
func (e *ExpectedQuery) WillReturnRows(rows *Rows) *ExpectedQuery {
	e.rows = rows
	return e
}

// WillReturnError defines the error to be returned for this expected query.
func (e *ExpectedQuery) WillReturnError(err error) *ExpectedQuery {
	e.err = err
	return e
}

// Rows mimics a small subset of sqlmock.Rows.
type Rows struct {
	columns   []string
	data      [][]driver.Value
	nextErrAt map[int]error // error returned from Next when moving to this row index
	mu        sync.Mutex
}

// NewRows creates a new Rows with given column names.
func NewRows(columns []string) *Rows {
	colsCopy := append([]string(nil), columns...)
	return &Rows{columns: colsCopy, nextErrAt: make(map[int]error)}
}

// AddRow appends a row.
func (r *Rows) AddRow(values ...any) *Rows {
	r.mu.Lock()
	defer r.mu.Unlock()

	row := make([]driver.Value, 0, len(values))
	for _, v := range values {
		row = append(row, driver.Value(v))
	}
	r.data = append(r.data, row)
	return r
}

// RowError sets an error returned from Rows.Next when it is asked to move to the
// specified row index. This is useful to simulate rows.Err().
func (r *Rows) RowError(row int, err error) *Rows {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextErrAt[row] = err
	return r
}

type mock struct {
	mu      sync.Mutex
	expects []*ExpectedQuery
	used    int
}

func (m *mock) ExpectQuery(query string) *ExpectedQuery {
	m.mu.Lock()
	defer m.mu.Unlock()
	e := &ExpectedQuery{query: query}
	m.expects = append(m.expects, e)
	return e
}

func (m *mock) popExpected(query string) (*ExpectedQuery, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.used >= len(m.expects) {
		return nil, fmt.Errorf("unexpected query: %q", query)
	}
	e := m.expects[m.used]
	m.used++
	if e.query != query {
		return nil, fmt.Errorf("query mismatch: got %q want %q", query, e.query)
	}
	return e, nil
}

func (m *mock) ExpectationsWereMet() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.used != len(m.expects) {
		return fmt.Errorf("not all expectations were met: %d remaining", len(m.expects)-m.used)
	}
	return nil
}

var (
	drvCounter uint64
	driversMu  sync.Mutex
	drivers    = map[string]*mockDriver{}
)

type mockDriver struct {
	m *mock
}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{m: d.m}, nil
}

type mockConn struct {
	m *mock
}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}
func (c *mockConn) Close() error              { return nil }
func (c *mockConn) Begin() (driver.Tx, error) { return nil, errors.New("not implemented") }

// QueryContext is used by database/sql for *sql.DB.QueryContext / Query.
func (c *mockConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	exp, err := c.m.popExpected(query)
	if err != nil {
		return nil, err
	}
	if exp.err != nil {
		return nil, exp.err
	}
	if exp.rows == nil {
		return &driverRows{cols: []string{}, data: nil, nextErrAt: map[int]error{}}, nil
	}
	return &driverRows{cols: exp.rows.columns, data: exp.rows.data, nextErrAt: exp.rows.nextErrAt}, nil
}

var _ driver.QueryerContext = (*mockConn)(nil)

type driverRows struct {
	cols      []string
	data      [][]driver.Value
	idx       int
	nextErrAt map[int]error
	closed    bool
}

func (r *driverRows) Columns() []string { return append([]string(nil), r.cols...) }

func (r *driverRows) Close() error {
	r.closed = true
	return nil
}

func (r *driverRows) Next(dest []driver.Value) error {
	if r.closed {
		return errors.New("rows closed")
	}
	// If an error is set for this row index, return it as Next error.
	// database/sql will surface it via rows.Err().
	if err, ok := r.nextErrAt[r.idx]; ok {
		return err
	}
	if r.idx >= len(r.data) {
		return io.EOF
	}
	row := r.data[r.idx]
	r.idx++
	for i := range dest {
		if i < len(row) {
			dest[i] = row[i]
		} else {
			dest[i] = nil
		}
	}
	return nil
}

// New creates a new mock database and controller.
// The returned *sql.DB can be passed into application code.
func New() (*sql.DB, Sqlmock, error) {
	m := &mock{}

	id := atomic.AddUint64(&drvCounter, 1)
	drvName := fmt.Sprintf("sqlmock-%d", id)

	driversMu.Lock()
	drivers[drvName] = &mockDriver{m: m}
	driversMu.Unlock()

	sql.Register(drvName, drivers[drvName])
	db, err := sql.Open(drvName, "")
	if err != nil {
		return nil, nil, err
	}
	return db, m, nil
}
