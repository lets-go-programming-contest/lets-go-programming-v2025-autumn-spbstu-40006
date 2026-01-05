package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	internaldb "github.com/vitsh1/task-6/internal/db"
)

func TestDB_GetNames_OK(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("vita").
		AddRow("nika")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := internaldb.New(mockDB)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"vita", "nika"}, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	expected := errors.New("query error")

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(expected)

	service := internaldb.New(mockDB)

	names, err := service.GetNames()
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Contains(t, err.Error(), "db query")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	scanErr := errors.New("scan error")

	rows := sqlmock.NewRows([]string{"name"}).AddRow("ok")

	rows.RowError(0, scanErr)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := internaldb.New(mockDB)

	names, err := service.GetNames()
	require.Error(t, err)
	require.ErrorIs(t, err, scanErr)
	require.Contains(t, err.Error(), "rows scanning")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockDB.Close()

	rowsErr := errors.New("rows error")

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("vita").
		AddRow("nika")
	rows.RowError(1, rowsErr)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := internaldb.New(mockDB)

	names, err := service.GetNames()
	require.Error(t, err)
	require.ErrorIs(t, err, rowsErr)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetUniqueNames_OK(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("vita").
		AddRow("nika")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := internaldb.New(mockDB)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"vita", "nika"}, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockDB.Close()

	expected := errors.New("query error")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(expected)

	service := internaldb.New(mockDB)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Contains(t, err.Error(), "db query")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockDB.Close()

	scanErr := errors.New("scan error")

	rows := sqlmock.NewRows([]string{"name"}).AddRow("ok")
	rows.RowError(0, scanErr)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := internaldb.New(mockDB)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.ErrorIs(t, err, scanErr)
	require.Contains(t, err.Error(), "rows scanning")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_GetUniqueNames_RowsErr(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockDB.Close()

	rowsErr := errors.New("rows error")

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("vita").
		AddRow("nika")
	rows.RowError(1, rowsErr)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := internaldb.New(mockDB)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.ErrorIs(t, err, rowsErr)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}
