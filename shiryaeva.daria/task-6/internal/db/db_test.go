package db_test

import (
    "database/sql"
    "io"
    "regexp"
    "strings"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    db "github.com/Dora-shi/task-6/internal/db"
)

func TestDBService_GetNames_OK(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    rows := sqlmock.NewRows([]string{"name"}).
        AddRow("Alice").
        AddRow("Bob")
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
        WillReturnRows(rows)
    
    service := db.New(sqlDB)
    
    got, err := service.GetNames()
    require.NoError(t, err, "GetNames should not fail")
    
    assert.Equal(t, []string{"Alice", "Bob"}, got, "names should match")
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetNames_QueryError(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
        WillReturnError(sql.ErrConnDone)
    
    service := db.New(sqlDB)
    
    _, gotErr := service.GetNames()
    require.Error(t, gotErr, "should return error")
    assert.Contains(t, gotErr.Error(), "db query:", "error should be wrapped")
    
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetNames_ScanError(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
        WillReturnRows(rows)
    
    service := db.New(sqlDB)
    
    _, gotErr := service.GetNames()
    require.Error(t, gotErr, "should return error")
    assert.Contains(t, gotErr.Error(), "rows scanning:", "error should be wrapped")
    
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetNames_RowsErr(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    rows := sqlmock.NewRows([]string{"name"}).
        AddRow("Alice").
        AddRow("Bob").
        RowError(1, io.ErrUnexpectedEOF)
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM users")).
        WillReturnRows(rows)
    
    service := db.New(sqlDB)
    
    _, gotErr := service.GetNames()
    require.Error(t, gotErr, "should return error")
    assert.Contains(t, gotErr.Error(), "rows error:", "error should be wrapped")
    
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    rows := sqlmock.NewRows([]string{"name"}).
        AddRow("Alice").
        AddRow("Alice"). // Дубликат
        AddRow("Bob")
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
        WillReturnRows(rows)
    
    service := db.New(sqlDB)
    
    got, err := service.GetUniqueNames()
    require.NoError(t, err, "GetUniqueNames should not fail")
    
    assert.Equal(t, []string{"Alice", "Alice", "Bob"}, got, "values should match")
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
        WillReturnError(sql.ErrConnDone)
    
    service := db.New(sqlDB)
    
    _, gotErr := service.GetUniqueNames()
    require.Error(t, gotErr, "should return error")
    assert.Contains(t, gotErr.Error(), "db query:", "error should be wrapped")
    
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
        WillReturnRows(rows)
    
    service := db.New(sqlDB)
    
    _, gotErr := service.GetUniqueNames()
    require.Error(t, gotErr, "should return error")
    assert.Contains(t, gotErr.Error(), "rows scanning:", "error should be wrapped")
    
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_GetUniqueNames_RowsErr(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    rows := sqlmock.NewRows([]string{"name"}).
        AddRow("Alice").
        AddRow("Bob").
        RowError(1, io.ErrUnexpectedEOF)
    
    mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT name FROM users")).
        WillReturnRows(rows)
    
    service := db.New(sqlDB)
    
    _, gotErr := service.GetUniqueNames()
    require.Error(t, gotErr, "should return error")
    assert.Contains(t, gotErr.Error(), "rows error:", "error should be wrapped")
    
    assert.NoError(t, mock.ExpectationsWereMet(), "all expectations should be met")
}

func TestDBService_New(t *testing.T) {
    t.Parallel()
    
    sqlDB, mock, err := sqlmock.New()
    require.NoError(t, err, "sqlmock.New should not fail")
    t.Cleanup(func() { sqlDB.Close() })
    
    mock.ExpectationsWereMet() // Нет ожиданий
    
    service := db.New(sqlDB)
    assert.NotNil(t, service, "service should not be nil")
    assert.Equal(t, sqlDB, service.DB, "DB should be set correctly")
}
