package database_test

import (
	"JWT/pkg/database"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func TestConnectMySQL_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	cfg := MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "arti2002",
		Database: "forum_db",
	}

	mock.ExpectPing()

	sqlOpen := func(driverName, dataSourceName string) (*sql.DB, error) {

		expectedDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

		assert.Equal(t, "mysql", driverName)
		assert.Equal(t, expectedDSN, dataSourceName)

		return db, nil
	}

	originalSQLOpen := sqlOpenFunc
	sqlOpenFunc = sqlOpen

	defer func() {
		sqlOpenFunc = originalSQLOpen
	}()

	resultDB, err := database.ConnectMySQL(database.MySQLConfig(cfg))
	require.NoError(t, err)
	assert.NotNil(t, resultDB)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	stats := resultDB.Stats()
	assert.Equal(t, 10, stats.MaxOpenConnections)

}

var sqlOpenFunc = func(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}
