package ephemeraldb

import (
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"os"
	"time"
)

var (
	// TO BE DEPRECATED!!
	pgPort = "5432"
)

const (
	pgDefaultPort string = "5432/tcp"
)

// GetEphemeralDBPort returns the local port mapped to the postgreSQL 5432 port in the docker container;
// the callers of RunWithContainer function must only use the returned port to compose the connection string;
// do not make any assumption on the port number as it is managed by the dockertest library and is subjective
// to changes.
func GetEphemeralDBPort() string {
	// pgPort global var is populated once the (new) docker containre is up and running.
	// it will likely be a different port number in each run.
	return pgPort
}

// RunWithContainer let the given function `proc` (procedure) to run in a well-prepared environment where:
// - A dockertest container (postgres container) is up and running;
// - The container has passed the connection test (Ping() returns no error)
// - A deferred teardown routine will clean up the container, as long as `proc` does not panic() nor calling os.Exit(..)
//
// The typical use case is to run this function with the TestMain(m *testing.M) function.
// Read this article to find out what TestMain can do
// https://medium.com/goingogo/why-use-testmain-for-testing-in-go-dafb52b406bc
// The most important thing to remember is:
// TestMain() runs once **per-package**.
// It means all the tests in one package will have access to the SAME container.
//
// There are a few disastrous situations that cause early-termination, thus the container(s) are left in the background
// - proc panic() or calling os.Exit(..)
// - terminated by a debugger
// - the SetUp routine fails half way
// - the teardown routine fails before completion
// For the last two cases, a well crafted error will be printed out for troubleshooting.
// **You must manually clean up these leftover containers.** (docker ps -q | xargs docker kill)
//
// nameGenerator:
// a function to generate the name for the dockertest container;
// if RunWithContainer is called concurrently, and there are N dockertest containers running in parallel, this function
// must make sure each container has a unique name.
// for example, it is not hard to implement a RandomString() function and return: "some-prefix" + RandomString()
//
// proc:
// a function that performs certain tasks and communicates with the dockertest container;
// this function must call GetEphemeralDBPort to compose the connection string;
// this function takes a pointer to the EphemeralConnection struct that provides a connection to the database.
func RunWithContainer(nameGenerator func() string, proc func(*EphemeralConnection) error) (reterror error) {
	// create a default EphemeralDB struct. We use it to spin up the container
	ephemeraldb := Default()

	// deferred teardown routine can populate this error
	reterror = nil

	// initialize a pool using the smart detection mechanism provided by dockertest
	pool, err := dockertest.NewPool("")
	if err != nil {
		reterror = err
		return
	}

	dockerContainerName := nameGenerator()
	loggingManager, err := WithTempDir("ephemeraldb", dockerContainerName)
	if err != nil {
		reterror = err
		return
	}

	var mountParams []string
	mountParams = append(mountParams, loggingManager.MountParams()...)

	command := []string{
		// NOTE: the argument -F (no-fsync) and --synchronous_commit=false seem to introduce test flakes
		//       hence commented out - we need to investigate another optimization approach not compromising
		//       the (strict) serialization promise
		"postgres",
		//"-F", // disable fsync to improve perf (https://www.postgresql.org/docs/9.5/app-postgres.html)
		// WARNING: it brings flake
		"-N", "500",
		//"--synchronous_commit=false", WARNING: it brings flake
		"-B", "2048MB",
	}
	command = append(command, loggingManager.PGLoggingParams()...)

	// spin up the postgres container
	pgURL := ephemeraldb.GetConnectionURL()
	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Name:       dockerContainerName,
		Tag:        ephemeraldb.PgVersion,
		Hostname:   ephemeraldb.PgHost,
		Labels:     map[string]string{"dockertest": "1"},
		Mounts:     mountParams,
		Cmd:        command,
		Env: []string{
			"POSTGRES_USER=" + ephemeraldb.PgUserName,
			"POSTGRES_PASSWORD=" + ephemeraldb.PgUserPassword,
			"POSTGRES_DB=" + pgURL.Path,
		},
	}
	resource, err := pool.RunWithOptions(&runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true // stopped container goes away by itself
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	// the container-teardown routine; will pass any error to the named return error `reterror`
	defer func() {
		if err := pool.Purge(resource); err != nil {
			reterror = ContainerError{"could not purge resource", err}
		}
	}()

	if err != nil {
		reterror = ContainerError{"could not connect to docker", err}
		return
	}

	// Update the port to the host-bound one
	ephemeraldb.PgPort = resource.GetPort(pgDefaultPort)
	// ONLY FOR COMPATIBILITY! WILL DEPRECATE THIS VARIABLE SOON!
	pgPort = ephemeraldb.PgPort

	// set up container log connection
	logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stderr:       true,
		Stdout:       true,
		Stream:       false,
	})
	if err != nil {
		reterror = ContainerError{"could not connect to postgres container log output", err}
		return
	}

	// container log connection teardown routine
	defer func() {
		if err := logWaiter.Close(); err != nil {
			reterror = ContainerError{"could not close container log", err}
			return
		}
		if err := logWaiter.Wait(); err != nil {
			reterror = ContainerError{"could not wait for container log to close", err}
			return
		}
	}()

	pool.MaxWait = 10 * time.Second

	// connection-tester will populate this pointer if successful
	var econ *EphemeralConnection = nil
	connectionTester := func() error {
		if ephemeralConn, err := ephemeraldb.Connect(); err != nil {
			return err
		} else {
			if err := ephemeralConn.conn.Ping(); err != nil {
				return err
			} else {
				// the ephemeral database is up and running; it has passed the connection test
				econ = ephemeralConn
				return nil
			}
		}
	}
	if err := pool.Retry(connectionTester); err != nil {
		reterror = ContainerError{"could not connect to postgres server", err}
		return
	}

	reterror = proc(econ)
	return
}
