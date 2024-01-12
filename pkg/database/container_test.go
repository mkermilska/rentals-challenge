package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName = "testingwithrentals"
	testDBUser = "root"
	testDBPass = "root"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	c, host, port := RunDBContainer(ctx)
	defer func() {
		_ = c.Terminate(ctx)
	}()

	var err error
	db, err = StartDbStore(StartUpOptions{
		DBHost:     host,
		DBPort:     port,
		DBName:     testDBName,
		DBUsername: testDBUser,
		DBPassword: testDBPass,
	})
	if err != nil {
		log.Fatalf("Failed to start the database: %s", err)
	}
	code := m.Run()
	os.Exit(code)
}

func RunDBContainer(ctx context.Context) (dbC testcontainers.Container, host string, port int) {
	basePort, err := nat.NewPort("tcp", "5432")
	if err != nil {
		log.Fatal(err)
	}

	sourceMountAbsPath, err := filepath.Abs("../../sql-init.sql")
	if err != nil {
		panic(err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432"},
		Env: map[string]string{
			"POSTGRES_USER":     testDBUser,
			"POSTGRES_PASSWORD": testDBPass,
			"POSTGRES_DB":       testDBName,
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: sourceMountAbsPath,
					Target: "/docker-entrypoint-initdb.d/sql-init.sql",
				},
			}
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForLog("listening on IPv4 address"),
		),
	}

	dbC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	host, err = dbC.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	natPort, err := dbC.MappedPort(ctx, basePort)
	if err != nil {
		log.Fatalf("Could not get test container port: %s", err)
	}

	port, err = strconv.Atoi(string(natPort.Port()))
	if err != nil {
		log.Fatalf("Could not parse test container port: %s", err)
	}

	return
}
