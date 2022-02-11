package ephemeraldb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/url"
	"time"
)

// EphemeralDB represents the database running inside the docker container.
// It is created by the provisioner and passed to the test runners.
// It's supposed to treated as an immutable value.
type EphemeralDB struct {
	PgVersion             string
	PgHost                string
	PgPort                string
	PgDBName              string
	PgUserName            string
	PgUserPassword        string
	MaxConnectionLifetime time.Duration
}

func Default() EphemeralDB {
	return EphemeralDB{
		"10.15",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"hunter1",
		time.Minute * 10,
	}
}
func (edb EphemeralDB) GetConnectionURL() *url.URL {
	pgURL := &url.URL{
		Host:   fmt.Sprintf("%s:%s", edb.PgHost, edb.PgPort),
		Scheme: "postgres",
		User:   url.UserPassword(edb.PgUserName, edb.PgUserPassword),
		Path:   edb.PgDBName,
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()
	return pgURL
}

func (edb EphemeralDB) GetConnectionString() string {
	return edb.GetConnectionURL().String()
}

// EphemeralConnection is populated by RunWithContainer and is automatically passed to
// to the test runner functions. This ensure each test group (usually a Go package) is allocated
// its own connection to an ephemeral database.
// Here is the recommended approach to construct this value:
// - Each test group should call NewInvalid to create this value in the package scope (making it a global var).
//   As the function name suggests, this connection at this moment is not yet initialized hence considered invalid.
// - The test orchestrator, RunWithContainer, initializes this value.
// - Once orchestration is done, all the tests in this group share the same underlying connection to
//   the database resource.
type EphemeralConnection struct {
	edb  EphemeralDB
	conn *sqlx.DB
}

// NewInvalid creates an invalid connection value waiting to be initialized by RunWithContainer
func NewInvalid() EphemeralConnection {
	return EphemeralConnection{Default(), nil}
}

// Conn returns the underlying connection to the caller.
// It will panic if used outside the scope of RunWithContainer
func (econ EphemeralConnection) Conn() *sqlx.DB {
	if econ.conn == nil {
		panic("Invalid connection! This is likely caused by using the connection outside the scope of RunWithContainer()!")
	}
	return econ.conn
}

// CopyTo shallow-copy (field-by-field copy) `self` to `other`
// See the unit test for a demonstration of this behaviour.
func (econ EphemeralConnection) CopyTo(other *EphemeralConnection) {
	other.edb = econ.edb
	other.conn = econ.conn
}
