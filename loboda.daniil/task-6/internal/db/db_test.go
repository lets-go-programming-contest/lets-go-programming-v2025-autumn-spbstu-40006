package db_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"loboda.daniil/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
)

var (
	errBoom     = errors.New("boom")
	errIterBoom = errors.New("iter boom")
)

func TestDBService_GetNames_Success(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	got, err := svc.GetNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"alice", "bob"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(errBoom)

	svc := db.New(dbConn)

	_, err = svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errBoom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}

	if !strings.Contains(err.Error(), "db query") {
		t.Fatalf("expected db query context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	// чтобы rows.Next() было true, но Scan упал -> кладём int вместо string
	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	_, err = svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "rows scanning") {
		t.Fatalf("expected rows scanning context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetNames_RowsErr(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		RowError(1, errIterBoom)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	_, err = svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errIterBoom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}

	if !strings.Contains(err.Error(), "rows error") {
		t.Fatalf("expected rows error context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	got, err := svc.GetUniqueNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"alice", "bob"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errBoom)

	svc := db.New(dbConn)

	_, err = svc.GetUniqueNames()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errBoom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}

	if !strings.Contains(err.Error(), "db query") {
		t.Fatalf("expected db query context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(123)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	_, err = svc.GetUniqueNames()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "rows scanning") {
		t.Fatalf("expected rows scanning context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_RowsErr(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		RowError(1, errIterBoom)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	_, err = svc.GetUniqueNames()
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, errIterBoom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}

	if !strings.Contains(err.Error(), "rows error") {
		t.Fatalf("expected rows error context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_Dedup(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	// Дубликат "alice" должен быть отброшен внутри логики UniqueNames.
	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob").
		AddRow("alice")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := db.New(dbConn)

	got, err := svc.GetUniqueNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем как множество (не завязываемся на порядок).
	if len(got) != 2 {
		t.Fatalf("expected 2 unique names, got %v", got)
	}

	set := map[string]bool{}
	for _, v := range got {
		set[v] = true
	}

	if !set["alice"] || !set["bob"] {
		t.Fatalf("expected alice and bob, got %v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

