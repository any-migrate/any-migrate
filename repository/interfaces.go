package repository

import (
	"github.com/any-migrate/any-migrate/migration"
)

type Migration struct {
	// The index of the migration. Each index in a repository must be unique.
	Index int

	// The current state of the migration.
	State migration.State
}

// A Repository stores a lineage of migrations and their states.
type Repository interface {
	Update(Migration) error
	GetAll() ([]Migration, error)
	Scheme() string
}
