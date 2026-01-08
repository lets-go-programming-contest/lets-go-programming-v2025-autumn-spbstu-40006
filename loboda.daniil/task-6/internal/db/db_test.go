package db

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDBService_GetNames_Success(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").AddRow("bob")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := New(dbConn)
	got, err := svc.GetNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"alice", "bob"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("names mismatch: got %v want %v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	boom := errors.New("boom")
	mock.ExpectQuery("SELECT name FROM users").WillReturnError(boom)

	svc := New(dbConn)
	_, err = svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
	if !strings.Contains(err.Error(), "db query") {
		t.Fatalf("expected context in error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	// Return 2 columns so rows.Scan(&name) fails with a destination count mismatch.
	rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("alice", "x")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := New(dbConn)
	_, err = svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "rows scanning") {
		t.Fatalf("expected scan context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetNames_RowsErr(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	boom := errors.New("iter boom")
	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").RowError(1, boom)
	// RowError(1, ...) means: when trying to advance to row index 1 (the 2nd row), Next returns boom.
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	svc := New(dbConn)
	_, err = svc.GetNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, boom) {
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
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").AddRow("alice").AddRow("bob")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := New(dbConn)
	got, err := svc.GetUniqueNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"alice", "alice", "bob"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("values mismatch: got %v want %v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	boom := errors.New("boom")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(boom)

	svc := New(dbConn)
	_, err = svc.GetUniqueNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
	if !strings.Contains(err.Error(), "db query") {
		t.Fatalf("expected context in error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name", "extra"}).AddRow("alice", "x")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := New(dbConn)
	_, err = svc.GetUniqueNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "rows scanning") {
		t.Fatalf("expected scan context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_RowsErr(t *testing.T) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer dbConn.Close()

	boom := errors.New("iter boom")
	rows := sqlmock.NewRows([]string{"name"}).AddRow("alice").RowError(1, boom)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	svc := New(dbConn)
	_, err = svc.GetUniqueNames()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
	if !strings.Contains(err.Error(), "rows error") {
		t.Fatalf("expected rows error context, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
