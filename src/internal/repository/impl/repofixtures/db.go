package repofixtures

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/CargoMan0/avito-tech-task/pkg/database"
	"github.com/CargoMan0/avito-tech-task/pkg/postgrescontainer"
	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "/Users/bernsteinmond/ProgrammingProjects/GolandProjects/avito-tech-task/migrations"

	testUser     = "test"
	testPassword = "test"
	testName     = "test"
)

// TestDB is a wrapper over *sql.DB.
// This wrapper is required for storing methods to setup,
// shutdown and interact with test data during repository tests.
type TestDB struct {
	db *sql.DB
}

func NewTestDB(ctx context.Context) (*TestDB, error) {
	containerCfg := &postgrescontainer.Config{
		User:     testUser,
		Password: testPassword,
		DBName:   testName,
	}
	postgresContainer, err := postgrescontainer.New(containerCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres container: new: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("postgres container: get externally mapped port: %w", err)
	}

	databaseCfg := database.Config{
		Host:     "127.0.0.1",
		Port:     port.Int(),
		User:     testUser,
		Password: testPassword,
		DBName:   testName,
		SSLMode:  "disable",
	}

	sqlDB, err := database.NewSQL(databaseCfg)
	if err != nil {
		return nil, fmt.Errorf("new sql: %w", err)
	}

	testDB := &TestDB{
		db: sqlDB,
	}
	err = testDB.setup()
	if err != nil {
		return nil, fmt.Errorf("setup test database: %w", err)
	}

	return testDB, nil
}

func (t *TestDB) GetSQLDB() *sql.DB {
	return t.db
}

// setup sets up test database, by running up migrations and cleaning up tables.
func (t *TestDB) setup() error {
	err := t.runUpMigrations()
	if err != nil {
		return fmt.Errorf("run up migrations: %w", err)
	}

	err = t.cleanupTables()
	if err != nil {
		return fmt.Errorf("cleanup tables: %w", err)
	}

	return nil
}

// runUpMigrations runs up migrations with goose package.
func (t *TestDB) runUpMigrations() error {
	err := goose.Up(t.db, migrationsDir)
	if err != nil && !errors.Is(err, goose.ErrNoNextVersion) {

		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

// ShutDown terminates test database, by running down migrations and closing database connection.
func (t *TestDB) ShutDown() error {
	err := t.runDownMigrations()
	if err != nil {
		return fmt.Errorf("run down migrations: %w", err)
	}

	err = t.db.Close()
	if err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	return nil
}

// runDownMigrations runs down migrations with goose package.
func (t *TestDB) runDownMigrations() error {
	err := goose.Down(t.db, migrationsDir)
	if err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return fmt.Errorf("goose down: %w", err)
	}

	return nil
}

// cleanupTables deletes all rows from all tables in the database.
func (t *TestDB) cleanupTables() error {
	tables, err := t.getAllTables()
	if err != nil {
		return fmt.Errorf("failed to get tables for cleanup: %w", err)
	}

	for _, table := range tables {
		_, err = t.db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to delete table %s: %w", table, err)
		}
	}

	return nil
}

func (t *TestDB) getAllTables() ([]string, error) {
	var tables []string
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`
	rows, err := t.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}

	var tableName string
	for rows.Next() {
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, fmt.Errorf("scan row into table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}
