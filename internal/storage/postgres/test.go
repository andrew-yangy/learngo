package postgres

import (
	"fmt"
	"github.com/ddvkid/learngo/internal/dbmigration"
	"github.com/ddvkid/learngo/internal/ephemeraldb"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var lg = log.New()

func WithEphemeralDB(m *testing.M, blueprint *ephemeraldb.EphemeralConnection) {
	mainFunc := func(econ *ephemeraldb.EphemeralConnection) error {
		if filename, err := dbmigration.EngineInitScript(); err != nil {
			return err
		} else {
			if err := econ.RunScript(filename); err != nil {
				return err
			}
		}
		if dirname, err := dbmigration.EngineMigrationDir(); err != nil {
			return err
		} else {
			if err := econ.RunMigration(dirname); err != nil {
				return err
			}
		}
		if con, err := econ.Reposition(); err != nil {
			return err
		} else {
			con.CopyTo(blueprint)
		}
		if 0 != m.Run() {
			return fmt.Errorf("FAILED")
		}
		return nil
	}
	err := ephemeraldb.RunWithContainer(randomizedContainerName, mainFunc)
	if err != nil {
		lg.Fatalln(err)
	} else {
		os.Exit(0)
	}
}

// LocalConnection gives each test a connection to an ephemeral database
// The ephemeral database is created by cloning the blueprint db (efficient).
func LocalConnection(t *testing.T, blueprint ephemeraldb.EphemeralConnection) *TestStore {
	econ, err := blueprint.CloneFrom("postgres", "postgres"+ephemeraldb.RandomNumericString())
	assert.NoError(t, err)
	store := &PgStore{DB: econ.Conn(), Source: Source{Queryable: econ.Conn()}}
	return &TestStore{PgStore: store}
}

func RandomID() string {
	return ephemeraldb.RandomNumericString()
}

type TestStore struct {
	*PgStore
}

func randomizedContainerName() string {
	return "engine-inte-test-" + RandomID()
}
