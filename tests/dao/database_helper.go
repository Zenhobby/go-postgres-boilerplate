package dao_test

import (
	"context"
	"go-postgres-boilerplate/dao"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T) (*dao.Config, func()) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:11",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err, "failed to start container")

	host, err := postgresC.Host(ctx)
	assert.NoError(t, err, "failed to get container host")
	port, err := postgresC.MappedPort(ctx, "5432")
	assert.NoError(t, err, "failed to get container port")

	testConfig := &dao.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	cleanup := func() {
		if err := postgresC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	}

	return testConfig, cleanup
}
