package migrator

import (
	"errors"
	"io"
)

type State string

// All the different migration states.
const (
	NEW = iota // Important this is 0.
	TESTING
	TESTING_FAILED
	MIGRATING
	MIGRATION_FAILED
	VERIFYING
	VERIFICATION_FAILED
	MIGRATION_SUCCEEDED
)

// An error that can be returned if
var NOT_IMPLEMENTED = errors.New("The feature is not implemented by this driver.")

type Driver interface {
	// Initialize is the first function to be called.
	// Check the url string and open and verify any connection
	// that has to be made.
	Initialize(url string) error

	// Close is the last function to be called.
	// Close any open connection here.
	Close() error

	// FilenameExtension returns the extension of the migration files.
	// The returned string must not begin with a dot.
	FilenameExtension() string
}

type Upgrader interface {
	Driver

	// Apply a migration. It will receive a file which the driver should
	// apply to its backend or whatever. The migration function should
	// return an error if something failed, nil otherwise.
	Upgrade(io.Reader) error
}

type Downgrader interface {
	Driver

	// Undo a migration.
	Downgrade(io.Reader) error
}

type PreVerifyer interface {
	Driver

	// Test whether it looks like migration will apply cleanly. May
	// return `NOT_IMPLEMENTED` if not implemented.
	PreVerify(io.Reader) error
}

type PostVerifyer interface {
	Driver

	// Verify that a migration was applied cleanly after performed. May
	// return `NOT_IMPLEMENTED` if not implemented.
	PostVerify(io.Reader) error
}
