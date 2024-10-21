package dao_test

import (
	"go-postgres-boilerplate/dao"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	testConfig, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Test the Connect function
	db, err := dao.Connect(testConfig)
	assert.NoError(t, err, "Connect should not return an error")
	assert.NotNil(t, db, "Connect should return a non-nil database connection")

	// Test the connection
	err = db.Ping()
	assert.NoError(t, err, "Should be able to ping the database")

	// Close the connection
	err = db.Close()
	assert.NoError(t, err, "Should be able to close the connection without error")
}

func TestNewConfig(t *testing.T) {
	// Set environment variables
	t.Setenv("POSTGRES_HOST", "testhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "testuser")
	t.Setenv("POSTGRES_PASSWORD", "testpass")
	t.Setenv("POSTGRES_DB", "testdb")

	// Call NewConfig
	config := dao.NewConfig()

	// Assert that the config values match the environment variables
	assert.Equal(t, "testhost", config.Host)
	assert.Equal(t, "5432", config.Port)
	assert.Equal(t, "testuser", config.User)
	assert.Equal(t, "testpass", config.Password)
	assert.Equal(t, "testdb", config.DBName)
	assert.Equal(t, "disable", config.SSLMode)
}
