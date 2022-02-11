package dbmigration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// to reliably locate the migration resource regardless of the caller's stack

type LocationServiceError struct {
	msg   string
	inner error
}

func (err LocationServiceError) Error() string {
	return fmt.Sprintf("DB-Migration location service error: %s ==> %s", err.msg, err.inner)
}

func isDirectoryOrError(filename string) error {
	if stat, err := os.Stat(filename); err != nil {
		return err
	} else {
		if stat.IsDir() {
			return nil
		} else {
			return fmt.Errorf("%s is not a directory", filename)
		}
	}
}

func isFileOrError(filename string) error {
	if stat, err := os.Stat(filename); err != nil {
		return err
	} else {
		if stat.Size() > 0 {
			return nil
		} else {
			return fmt.Errorf("%s is an invalid file", filename)
		}
	}
}

func migrationDir(pathSegs []string) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	// Take into account vendor - why is this not always required?
	if strings.Contains(filename, "/vendor/") {
		vendorSegs := []string{"..", "..", "..", "..", ".."}
		pathSegs = append(vendorSegs, pathSegs...)
	}
	pathSegs = append([]string{filename}, pathSegs...)
	relativePath := path.Join(pathSegs...)
	if absFilename, err := filepath.Abs(relativePath); err != nil {
		return "", err
	} else {
		if err := isDirectoryOrError(absFilename); err != nil {
			return "", err
		}
		return absFilename, nil
	}
}

func EngineInitScript() (string, error) {
	dirname, err := migrationDir([]string{"..", "..", "..", "engine-migration", "schemas"})
	if err != nil {
		return "", LocationServiceError{"Failed to resolve Engine init script directory", err}
	}
	filename := path.Join(dirname, "_init.sql")
	if err := isFileOrError(filename); err != nil {
		return "", LocationServiceError{"", err}
	}
	return filename, nil
}

func EngineMigrationDir() (string, error) {
	filename, err := migrationDir([]string{"..", "..", "..", "engine-migration", "schemas", "schema-migrations"})
	if err != nil {
		return "", LocationServiceError{"Failed to resolve Engine migration directory", err}
	}
	return filename, nil
}

func APIMigrationDir() (string, error) {
	filename, err := migrationDir([]string{"..", "..", "..", "public-api", "db-migrations"})
	if err != nil {
		return "", LocationServiceError{"Failed to resolve API migration directory", err}
	}
	return filename, nil
}

func SnapshotInitScript() (string, error) {
	dirname, err := migrationDir([]string{"..", "..", "..", "snapshot-migration", "schemas"})
	if err != nil {
		return "", LocationServiceError{"Failed to resolve Snapshot init script directory", err}
	}
	filename := path.Join(dirname, "_init.sql")
	if err := isFileOrError(filename); err != nil {
		return "", LocationServiceError{"", err}
	}
	return filename, nil
}

func SnapshotMigrationDir() (string, error) {
	filename, err := migrationDir([]string{"..", "..", "..", "snapshot-migration", "schemas", "schema-migrations"})
	if err != nil {
		return "", LocationServiceError{"Failed to resolve Snapshot migration directory", err}
	}
	return filename, nil
}
