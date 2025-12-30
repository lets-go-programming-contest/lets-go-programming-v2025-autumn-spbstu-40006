package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetNames_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetNames_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(errors.New("query error"))

	service := New(db)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}

func TestGetNames_RowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}

func TestGetNames_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}

func TestGetUniqueNames_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(errors.New("query error"))

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
}
