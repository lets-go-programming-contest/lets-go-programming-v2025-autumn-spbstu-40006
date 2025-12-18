package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zaouiomar/task-6/internal/db"
)

var (
	errTest    = errors.New("test error")
	errRowTest = errors.New("row test error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	service := db.New(dbMock)

	assert.NotNil(t, service)
	assert.Equal(t, dbMock, service.DB)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("a").AddRow("b"))

		svc := db.New(dbMock)
		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		svc := db.New(dbMock)
		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnError(errTest)

		svc := db.New(dbMock)
		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		svc := db.New(dbMock)
		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("x")
		rows.RowError(0, errRowTest)
		mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)

		svc := db.New(dbMock)
		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("x").AddRow("y"))

		svc := db.New(dbMock)
		names, err := svc.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"x", "y"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		svc := db.New(dbMock)
		names, err := svc.GetUniqueNames()

		require.NoError(t, err)
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnError(errTest)

		svc := db.New(dbMock)
		names, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		svc := db.New(dbMock)
		names, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer dbMock.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("z")
		rows.RowError(0, errRowTest)
		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).WillReturnRows(rows)

		svc := db.New(dbMock)
		names, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
