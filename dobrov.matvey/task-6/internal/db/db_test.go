package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/netwite/task-6/internal/db"
	"github.com/stretchr/testify/suite"
)

type DBServiceTestSuite struct {
	suite.Suite
	mockDB *sql.DB
	mock   sqlmock.Sqlmock
}

func (s *DBServiceTestSuite) SetupTest() {
	var err error
	s.mockDB, s.mock, err = sqlmock.New()
	s.Require().NoError(err, "failed to create sqlmock")
}

func (s *DBServiceTestSuite) TearDownTest() {
	if s.mockDB != nil {
		s.mockDB.Close()
	}
}

func (s *DBServiceTestSuite) TestNew() {
	service := db.New(s.mockDB)
	s.Equal(s.mockDB, service.DB)
}

func (s *DBServiceTestSuite) TestGetNames_Success() {
	service := db.DBService{DB: s.mockDB}

	expectedRows := []string{"Alice", "Bob", "Charlie"}
	rows := sqlmock.NewRows([]string{"name"})

	for _, name := range expectedRows {
		rows.AddRow(name)
	}

	s.mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Equal(expectedRows, result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetNames_EmptyResult() {
	service := db.DBService{DB: s.mockDB}

	rows := sqlmock.NewRows([]string{"name"})
	s.mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	result, err := service.GetNames()

	s.Require().NoError(err)
	s.Empty(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetNames_QueryError() {
	service := db.DBService{DB: s.mockDB}

	testError := errors.New("connection failed") //nolint:err113
	s.mock.ExpectQuery("SELECT name FROM users").WillReturnError(testError)

	result, err := service.GetNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "db query")
	s.Contains(err.Error(), testError.Error())
	s.Nil(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetNames_ScanError() {
	service := db.DBService{DB: s.mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	s.mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	result, err := service.GetNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "rows scanning")
	s.Nil(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetNames_RowsError() {
	service := db.DBService{DB: s.mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("row error")) //nolint:err113
	s.mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	result, err := service.GetNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "rows error")
	s.Nil(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetUniqueNames_Success() {
	service := db.DBService{DB: s.mockDB}

	uniqueNames := []string{"Alice", "Bob", "Charlie"}

	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range uniqueNames {
		rows.AddRow(name)
	}

	s.mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	result, err := service.GetUniqueNames()

	s.Require().NoError(err)
	s.Equal(uniqueNames, result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetUniqueNames_EmptyResult() {
	service := db.DBService{DB: s.mockDB}

	rows := sqlmock.NewRows([]string{"name"})
	s.mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	result, err := service.GetUniqueNames()

	s.Require().NoError(err)
	s.Empty(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetUniqueNames_QueryError() {
	service := db.DBService{DB: s.mockDB}

	testError := errors.New("connection failed") //nolint:err113
	s.mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(testError)

	result, err := service.GetUniqueNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "db query")
	s.Contains(err.Error(), testError.Error())
	s.Nil(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetUniqueNames_ScanError() {
	service := db.DBService{DB: s.mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	s.mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	result, err := service.GetUniqueNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "rows scanning")
	s.Nil(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetUniqueNames_RowsError() {
	service := db.DBService{DB: s.mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice").RowError(0, errors.New("row error")) //nolint:err113
	s.mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	result, err := service.GetUniqueNames()

	s.Require().Error(err)
	s.Require().ErrorContains(err, "rows error")
	s.Nil(result)

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func (s *DBServiceTestSuite) TestGetUniqueNames_ReturnsOnlyUnique() {
	service := db.DBService{DB: s.mockDB}

	uniqueRows := []string{"John", "Jane"}
	rows := sqlmock.NewRows([]string{"name"})

	for _, name := range uniqueRows {
		rows.AddRow(name)
	}

	s.mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	result, err := service.GetUniqueNames()

	s.Require().NoError(err)
	s.Equal(uniqueRows, result)
	s.Len(result, 2, "Должно вернуть только уникальные значения")

	err = s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func TestDBServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DBServiceTestSuite))
}
