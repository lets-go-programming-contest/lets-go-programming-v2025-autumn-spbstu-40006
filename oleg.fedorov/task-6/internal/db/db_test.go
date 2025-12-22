package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dizey5k/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DatabaseTestCase struct {
	name        string
	query       string
	mockRows    []string
	expectError bool
	errorType   error
	setupMock   func(sqlmock.Sqlmock, []string)
}

func TestDatabaseOperations(t *testing.T) {
	testError := errors.New("database operation failed")

	testCases := []DatabaseTestCase{
		{
			name:      "getting users list",
			query:     "SELECT name FROM users",
			mockRows:  []string{"Анна", "Борис", "Виктор"},
			setupMock: setupQueryMock,
		},
		{
			name:      "empty db",
			query:     "SELECT name FROM users",
			mockRows:  []string{},
			setupMock: setupQueryMock,
		},
		{
			name:        "err while proccess",
			query:       "SELECT name FROM users",
			expectError: true,
			errorType:   testError,
			setupMock:   setupErrorMock,
		},
		{
			name:        "err str scan",
			query:       "SELECT name FROM users",
			mockRows:    []string{"Некорректные данные"},
			expectError: true,
			setupMock:   setupScanErrorMock,
		},
		{
			name:      "get unique names",
			query:     "SELECT DISTINCT name FROM users",
			mockRows:  []string{"Анна", "Борис", "Виктор"},
			setupMock: setupQueryMock,
		},
		{
			name:        "err while proccessing str",
			query:       "SELECT DISTINCT name FROM users",
			mockRows:    []string{"Анна"},
			expectError: true,
			errorType:   testError,
			setupMock:   setupRowsErrorMock,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sqlDB, mock, err := sqlmock.New()
			require.NoError(t, err, "creating mock db")
			defer sqlDB.Close()

			tc.setupMock(mock, tc.mockRows)

			service := db.New(sqlDB)
			var result []string
			var execErr error

			if tc.query == "SELECT DISTINCT name FROM users" {
				result, execErr = service.GetUniqueNames()
			} else {
				result, execErr = service.GetNames()
			}

			if tc.expectError {
				require.Error(t, execErr)
				if tc.errorType != nil {
					assert.ErrorIs(t, execErr, tc.errorType)
				}
			} else {
				require.NoError(t, execErr)
				assert.Equal(t, tc.mockRows, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDatabaseEdgeScenarios(t *testing.T) {
	t.Run("process nil", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Валидное имя").
			AddRow(nil).
			AddRow("Еще одно имя")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		_, err = service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
	})

	t.Run("checking connection close", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Пользователь 1").
			AddRow("Пользователь 2")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		service := db.New(sqlDB)
		result, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"Пользователь 1", "Пользователь 2"}, result)
	})

	t.Run("service with zero db", func(t *testing.T) {
		t.Parallel()

		service := db.DBService{}
		result, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "db query")
	})
}

func setupQueryMock(mock sqlmock.Sqlmock, rows []string) {
	sqlRows := sqlmock.NewRows([]string{"name"})
	for _, row := range rows {
		sqlRows.AddRow(row)
	}
	mock.ExpectQuery("SELECT (name|DISTINCT name) FROM users").WillReturnRows(sqlRows)
}

func setupErrorMock(mock sqlmock.Sqlmock, _ []string) {
	mock.ExpectQuery("SELECT (name|DISTINCT name) FROM users").
		WillReturnError(errors.New("database operation failed"))
}

func setupScanErrorMock(mock sqlmock.Sqlmock, _ []string) {
	sqlRows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(sqlRows)
}

func setupRowsErrorMock(mock sqlmock.Sqlmock, rows []string) {
	sqlRows := sqlmock.NewRows([]string{"name"})
	for _, row := range rows {
		sqlRows.AddRow(row)
	}
	sqlRows.CloseError(errors.New("rows processing error"))
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(sqlRows)
}

func TestDatabaseIntegration(t *testing.T) {
	t.Run("method calls", func(t *testing.T) {
		t.Parallel()

		sqlDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		rows1 := sqlmock.NewRows([]string{"name"}).
			AddRow("Иван").
			AddRow("Мария").
			AddRow("Иван")

		rows2 := sqlmock.NewRows([]string{"name"}).
			AddRow("Иван").
			AddRow("Мария")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows1)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows2)

		service := db.New(sqlDB)

		allNames, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Иван", "Мария", "Иван"}, allNames)

		uniqueNames, err := service.GetUniqueNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"Иван", "Мария"}, uniqueNames)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

type CustomDBMock struct {
	queryFunc func(query string, args ...any) (*sql.Rows, error)
}

func (m *CustomDBMock) Query(query string, args ...any) (*sql.Rows, error) {
	return m.queryFunc(query, args...)
}

func TestWithCustomMock(t *testing.T) {
	t.Run("custom mock", func(t *testing.T) {
		t.Parallel()

		mockDB := &CustomDBMock{
			queryFunc: func(query string, args ...any) (*sql.Rows, error) {
				rows := &sql.Rows{}
				return rows, nil
			},
		}

		service := db.New(mockDB)
		assert.NotNil(t, service)
	})
}
