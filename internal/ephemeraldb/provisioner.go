package ephemeraldb

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"io/ioutil"
)

func (edb EphemeralDB) Connect() (*EphemeralConnection, error) {
	sqldb, err := sql.Open("nrpostgres", edb.GetConnectionString())
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(sqldb, "postgres")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(edb.MaxConnectionLifetime)
	return &EphemeralConnection{edb, db}, nil
}

func (econ EphemeralConnection) RunScript(script string) error {
	data, err := ioutil.ReadFile(script)
	if err != nil {
		return err
	}
	_, err = econ.conn.Exec(string(data))
	return err
}

type MigrationError struct {
	path  string
	inner error
}

func (err MigrationError) Error() string {
	return fmt.Sprintf("Failed to load migration scripts from: %s ==> %s", err.path, err.inner)
}

func (econ EphemeralConnection) RunMigration(migrationDir string) error {
	driver, err := postgres.WithInstance(econ.conn.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	asUrl := fmt.Sprintf("file://%s", migrationDir)
	m, err := migrate.NewWithDatabaseInstance(asUrl, econ.edb.PgDBName, driver)
	if err != nil {
		return MigrationError{migrationDir, err}
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return MigrationError{migrationDir, err}
	}
	return nil
}

// Disconnect will disconnect all the existing connections.
// Do this before calling Clone()
func (econ EphemeralConnection) Disconnect() error {
	dbName := econ.edb.PgDBName
	cmd := fmt.Sprintf(
		"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s' AND pid <> pg_backend_pid();", dbName)
	if _, err := econ.conn.Exec(cmd); err != nil {
		return err
	}
	return nil
}

// disconnectAllFrom will disconnect all the existing connections from the given database while the caller
// is connected to a different db;
// It's mainly used to clean up the state of the blueprint database before the caller clones it.
// It has a stronger guarantee of eliminating connections than Disconnect
// Think of it as `sudo kill -9 <pid>`
// WARNING: it can be harmful if used too often, hence this function is NOT exported!
func (econ EphemeralConnection) disconnectAllFrom(dbName string) error {
	cmd := fmt.Sprintf(
		"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%s';", dbName)
	if _, err := econ.conn.Exec(cmd); err != nil {
		return err
	}
	return nil
}

// Reposition will detach the caller's from its current database, then reconnect it to the template1 database.
// This completely frees the original database from active connections.
// This is NOT a generic re-connection mechanism as it internally "kills" all the live connection to the
// original database, hence should ONLY be used once.
// At the moment the destination database is hardcoded to `template1`.
func (econ EphemeralConnection) Reposition() (*EphemeralConnection, error) {
	econ.edb.PgDBName = "template1"
	newEcon, err := econ.edb.Connect()
	if err != nil {
		return nil, err
	}
	if err := newEcon.disconnectAllFrom("postgres"); err != nil {
		return nil, err
	}
	return newEcon, nil

}

// Clone creates a deepcopy of the EphemeralConnection
// The new database has the exact same contents but a different name. It's named with a random numeric suffix.
//
// It expects the source database is disconnected meaning the source connection called Disconnect
// which didn't return an error, otherwise the it returns an error to the caller:
// "Can not create from template XYZ because other users are still connecting to it"
//
// Clone is typically use for performance optimization.
// Once a database has loaded many migration scripts, it is much faster to duplicate it instead of
// reloading the same set of scripts to create a brand new database with the same data.
// According to PostgreSQL document, duplication happens on the lowest file level.
func (econ EphemeralConnection) Clone() (*EphemeralConnection, error) {
	dbUser := econ.edb.PgUserName
	dbNameSrc := econ.edb.PgDBName
	dbNameDest := dbNameSrc + RandomNumericString()
	cmd := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;", dbNameDest, dbNameSrc, dbUser)
	if _, err := econ.conn.Exec(cmd); err != nil {
		return nil, err
	}
	econ.edb.PgDBName = dbNameDest
	return econ.edb.Connect()
}

// CloneFrom is the 3rd person version of Clone
// It clones a source database not in any active connections, and gives the new database a caller-provided name
func (econ EphemeralConnection) CloneFrom(dbNameSrc string, dbNameDest string) (*EphemeralConnection, error) {
	dbUser := econ.edb.PgUserName
	cmd := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s OWNER %s;", dbNameDest, dbNameSrc, dbUser)
	if _, err := econ.conn.Exec(cmd); err != nil {
		return nil, err
	}
	econ.edb.PgDBName = dbNameDest
	return econ.edb.Connect()
}
