package db_test

import (
	"errors"
	"task-6/internal/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errTestQuery = errors.New("query failed")
	errTestRows  = errors.New("rows error")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errTestQuery)

		service := db.New(sqlDB)
		_, err = service.GetNames()

		require.ErrorIs(t, err, errTestQuery)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		_, err = service.GetNames()

		require.Error(t, err)
	})

	t.Run("rows_err", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errTestRows)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		_, err = service.GetNames()

		require.ErrorIs(t, err, errTestRows)
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").AddRow("Bob")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
	})

	t.Run("query_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errTestQuery)

		service := db.New(sqlDB)
		_, err = service.GetUniqueNames()

		require.ErrorIs(t, err, errTestQuery)
	})

	t.Run("scan_error", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		_, err = service.GetUniqueNames()

		require.Error(t, err)
	})

	t.Run("rows_err", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").CloseError(errTestRows)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		_, err = service.GetUniqueNames()

		require.ErrorIs(t, err, errTestRows)
	})
}