package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/identicalaffiliation/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errDB      = errors.New("db error")
	errRow     = errors.New("row iteration error")
	errSQLMock = errors.New("fail to create sqlmock")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err, errSQLMock.Error())
	defer mockDB.Close()

	service := db.New(mockDB)
	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("abc").AddRow("123"))

		service := db.New(mockDB)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"abc", "123"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()
		service := db.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errDB)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()
		service := db.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("1")
		rows.RowError(0, errRow)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()
	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("abc").AddRow("cba"))

		service := db.New(mockDB)
		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"abc", "cba"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()
		service := db.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errDB)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()
		service := db.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)

		defer mockDB.Close()

		service := db.New(mockDB)
		rows := sqlmock.NewRows([]string{"name"}).AddRow("1")
		rows.RowError(0, errRow)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
