package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mydb "github.com/MoneyprogerISG/task-6/internal/db"
)

var (
	errQuery = errors.New("query error")
	errRows  = errors.New("rows error")
)

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := mydb.New(dbMock)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(errQuery)

	service := mydb.New(dbMock)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
}

func TestGetNames_RowsError(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRows)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := mydb.New(dbMock)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := mydb.New(dbMock)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	service := mydb.New(dbMock)
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(errQuery)

	service := mydb.New(dbMock)
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errRows)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	service := mydb.New(dbMock)
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	dbMock, mock, _ := sqlmock.New()
	defer dbMock.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	service := mydb.New(dbMock)
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
}
