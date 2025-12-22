package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	Mdb "github.com/filon6/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

func TestNewDBServiceStoresDB(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectationsWereMet()

	service := Mdb.New(db)
	require.Equal(t, db, service.DB, "expected DB to be set")
}

func TestGetNames_ReturnsAllNames(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Angelina").
		AddRow("Vladislav").
		AddRow("Vadim").
		AddRow("Petr").
		AddRow("Alexey")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.NoError(t, err)

	require.Len(t, names, 5, "expected 5 names")
	require.Equal(t, "Angelina", names[0], "first name should be Angelina")
	require.Equal(t, "Vladislav", names[1], "second name should be Vladislav")
	require.Equal(t, "Vadim", names[2], "third name should be Vadim")
	require.Equal(t, "Petr", names[3], "fourth name should be Petr")
	require.Equal(t, "Alexey", names[4], "fifth name should be Alexey")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetNames_KeepsDuplicates(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Angelina").
		AddRow("Vladislav").
		AddRow("Angelina").
		AddRow("Petr").
		AddRow("Vladislav")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.NoError(t, err)

	require.Len(t, names, 5, "expected 5 names with duplicates")
	require.Equal(t, []string{"Angelina", "Vladislav", "Angelina", "Petr", "Vladislav"}, names)

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetNames_EmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.NoError(t, err)
	require.Empty(t, names, "expected empty slice")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetNames_QueryErrorWrapped(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(sql.ErrConnDone)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err, "expected error")
	require.Nil(t, names, "expected nil result on error")
	require.Contains(t, err.Error(), "db query", "error should contain 'db query'")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetNames_ScanErrorWrapped(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Angelina").
		AddRow(nil)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err, "expected error")
	require.Nil(t, names, "expected nil result on error")
	require.Contains(t, err.Error(), "rows scanning", "error should contain 'rows scanning'")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetNames_RowsErrorWrapped(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Angelina")
	rows.RowError(0, sql.ErrTxDone)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err, "expected error")
	require.Nil(t, names, "expected nil result on error")
	require.Contains(t, err.Error(), "rows error", "error should contain 'rows error'")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetUniqueNames_FullSetReturned(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Angelina").
		AddRow("Vladislav").
		AddRow("Vadim").
		AddRow("Petr").
		AddRow("Alexey").
		AddRow("Angelina").
		AddRow("Petr")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Len(t, names, 7, "expected 7 names (with duplicates from DISTINCT)")

	expected := []string{"Angelina", "Vladislav", "Vadim", "Petr", "Alexey", "Angelina", "Petr"}
	for i, name := range names {
		require.Equal(t, expected[i], name, "name mismatch at index %d", i)
	}

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetUniqueNames_AllAngelina(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Angelina").
		AddRow("Angelina").
		AddRow("Angelina")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Len(t, names, 3, "expected 3 Angelina names")

	for i, name := range names {
		require.Equal(t, "Angelina", name, "all names should be Angelina at index %d", i)
	}

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetUniqueNames_EmptyResult(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Empty(t, names, "expected empty slice")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetUniqueNames_QueryErrorWrapped(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnError(sql.ErrConnDone)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.Error(t, err, "expected error")
	require.Nil(t, names, "expected nil result on error")
	require.Contains(t, err.Error(), "db query", "error should contain 'db query'")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetUniqueNames_ScanErrorWrapped(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Angelina").
		AddRow(nil)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.Error(t, err, "expected error")
	require.Nil(t, names, "expected nil result on error")
	require.Contains(t, err.Error(), "rows scanning", "error should contain 'rows scanning'")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetUniqueNames_RowsErrorWrapped(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Angelina")
	rows.RowError(0, sql.ErrTxDone)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := Mdb.New(db)
	names, err := service.GetUniqueNames()
	require.Error(t, err, "expected error")
	require.Nil(t, names, "expected nil result on error")
	require.Contains(t, err.Error(), "rows error", "error should contain 'rows error'")

	require.NoError(t, mock.ExpectationsWereMet(), "unfulfilled expectations")
}

func TestGetNames_QueryError_NoPanic(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnError(sql.ErrConnDone)

	service := Mdb.New(db)
	names, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, names)
	require.Contains(t, err.Error(), "db query")

	require.NoError(t, mock.ExpectationsWereMet())
}
