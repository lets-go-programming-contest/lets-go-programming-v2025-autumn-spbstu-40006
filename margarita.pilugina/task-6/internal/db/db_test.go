package db

import (
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDBService_GetNames_OK(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(
			sqlmock.NewRows([]string{"name"}).
				AddRow("Alice").
				AddRow("Bob"),
		)

	service := New(sqlDB)
	got, err := service.GetNames()
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(got) != 2 || got[0] != "Alice" || got[1] != "Bob" {
		t.Fatalf("unexpected result: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnError(errors.New("db down"))

	service := New(sqlDB)
	_, err = service.GetNames()
	if err == nil || !strings.Contains(err.Error(), "db query") {
		t.Fatalf("expected wrapped db query error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "extra"}).
				AddRow("Alice", "boom"),
		)

	service := New(sqlDB)
	_, err = service.GetNames()
	if err == nil || !strings.Contains(err.Error(), "rows scanning") {
		t.Fatalf("expected wrapped rows scanning error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetNames_RowsErr(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		RowError(1, errors.New("row iteration error")) // error occurs during iteration

	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(rows)

	service := New(sqlDB)
	_, err = service.GetNames()
	if err == nil || !strings.Contains(err.Error(), "rows error") {
		t.Fatalf("expected wrapped rows error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(
			sqlmock.NewRows([]string{"name"}).
				AddRow("Alice").
				AddRow("Bob"),
		)

	service := New(sqlDB)
	got, err := service.GetUniqueNames()
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(got) != 2 || got[0] != "Alice" || got[1] != "Bob" {
		t.Fatalf("unexpected result: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnError(errors.New("db down"))

	service := New(sqlDB)
	_, err = service.GetUniqueNames()
	if err == nil || !strings.Contains(err.Error(), "db query") {
		t.Fatalf("expected wrapped db query error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(
			sqlmock.NewRows([]string{"name", "extra"}).
				AddRow("Alice", "boom"),
		)

	service := New(sqlDB)
	_, err = service.GetUniqueNames()
	if err == nil || !strings.Contains(err.Error(), "rows scanning") {
		t.Fatalf("expected wrapped rows scanning error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDBService_GetUniqueNames_RowsErr(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		RowError(1, errors.New("row iteration error"))

	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(rows)

	service := New(sqlDB)
	_, err = service.GetUniqueNames()
	if err == nil || !strings.Contains(err.Error(), "rows error") {
		t.Fatalf("expected wrapped rows error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
