package db_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mishaa105/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { dbMock.Close() })
	return dbMock, mock
}

const queryNames = "SELECT name FROM users"

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success with data", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Misha").
			AddRow("Masha")

		mock.ExpectQuery(regexp.QuoteMeta(queryNames)).WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.NotEmpty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(regexp.QuoteMeta(queryNames)).WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).WillReturnError(assert.AnError)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(regexp.QuoteMeta(queryNames)).WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Misha").
			RowError(0, assert.AnError)

		mock.ExpectQuery(regexp.QuoteMeta(queryNames)).WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	t.Run("success with data", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Misha").
			AddRow("Masha")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.NotEmpty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Empty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnError(assert.AnError)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		dbMock, mock := newMockDB(t)
		service := db.New(dbMock)

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Misha").
			RowError(0, assert.AnError)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		assert.Nil(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
