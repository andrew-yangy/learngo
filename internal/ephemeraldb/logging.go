package ephemeraldb

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"
)

// Logging is enabled by an environment variable `ENABLE_DOCKERTEST_LOGGING`;
// (its value doesn't matter, as long as it is defined in the calling environment)
// by default, when this var is undefined, the logging structure won't initialize the log output
// directory, and won't set the postgres command line parameters;
// this design is to simplify the caller's work.
// The unit tests of ephemeraldb package should call `newLoggingManager(dir, basename, true)` to bypass
// this implicit behavior (to always enable logging).
type Logging struct {
	dir      string
	basename string
	enabled  bool
}

// WithDir joins dir and basename to get the full path of the log file, i.e. <dir>/<basename>
// it will truncate the existing log file;
// it will create the directory if it doesn't exist;
// it will set the directory writable;
func WithDir(dir string, basename string) (*Logging, error) {
	return newLoggingManager(dir, basename, os.Getenv("ENABLE_DOCKERTEST_LOGGING") != "")
}

func newLoggingManager(dir string, basename string, enabled bool) (*Logging, error) {
	l := &Logging{dir, basename, enabled}
	if !enabled {
		return l, nil
	}
	if err := l.ensureDirExists(); err != nil {
		return nil, err
	}
	if err := l.touchOrTruncate(); err != nil {
		return nil, err
	}
	return l, nil
}

// WithTempDir call the standard library's TempDir function to get a writable location;
// this depends on the operating system's libc implementation.
// on POSIX systems this is usually /var/tmp (or /tmp)
// parent: specify the parent directory, i.e. <OS temp dir>/<parent>/<basename>
func WithTempDir(parent string, basename string) (*Logging, error) {
	dir := path.Join(os.TempDir(), parent)
	return WithDir(dir, basename)
}

func newLoggingManagerTempDir(parent string, basename string, enabled bool) (*Logging, error) {
	dir := path.Join(os.TempDir(), parent)
	return newLoggingManager(dir, basename, enabled)
}

func (l *Logging) Teardown() error {
	return os.Remove(l.Filename())
}

func (l *Logging) ensureDirExists() error {
	oldUmask := syscall.Umask(000)
	defer func() {
		syscall.Umask(oldUmask)
	}()
	if fileInfo, statError := os.Stat(l.dir); statError != nil {
		if errors.Is(statError, os.ErrNotExist) {
			creatError := os.Mkdir(l.dir, 0777)
			if creatError != nil {
				return creatError
			}
			return nil // create successfully
		}
		return statError
	} else {
		if fileInfo.IsDir() {
			if chmodError := os.Chmod(l.dir, 0777); chmodError != nil {
				return chmodError
			}
			return nil // chmod successfully
		}
		return fmt.Errorf("invalid directory path (is a file): %s", l.dir)
	}
}

func (l *Logging) touchOrTruncate() error {
	oldUmask := syscall.Umask(000)
	defer func() {
		syscall.Umask(oldUmask)
	}()
	filename := l.Filename()
	if fileInfo, statError := os.Stat(filename); statError != nil {
		if errors.Is(statError, os.ErrNotExist) {
			if _, createError := os.Create(filename); createError != nil {
				return createError
			}
			if chmodError := os.Chmod(filename, 0666); chmodError != nil {
				return chmodError
			}
			return nil // create + chmod successfully
		} else {
			return statError
		}
	} else {
		if fileInfo.IsDir() {
			return fmt.Errorf("failed to create log file (file exists and is a directory): %s", filename)
		}
		if chmodError := os.Chmod(filename, 0666); chmodError != nil {
			return chmodError
		}
		return nil // chmod successfully
	}
}

func (l *Logging) Filename() string {
	return path.Join(l.dir, l.basename)
}

func (l *Logging) PGLoggingParams() []string {
	if !l.enabled {
		return []string{}
	}
	return []string{
		"-c", "logging_collector=on",
		"-c", "log_directory=/pg_logs",
		"-c", "log_filename=" + l.basename,
		"-c", "log_statement=all",
	}
}

func (l *Logging) MountParams() []string {
	if !l.enabled {
		return []string{}
	}
	return []string{
		fmt.Sprintf("%s:/pg_logs", l.dir),
	}
}
