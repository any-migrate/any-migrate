package migrate

import (
	"os"

	"github.com/any-migrate/any-migrate/migration"
	"github.com/any-migrate/any-migrate/repository"
)

type MigrationState interface {
	Upgrade() (MigrationState, done bool)
	Downgrade() (MigrationState, done bool)
	StateCode() migration.State
}

type StopPolicy interface {
	Done(Repository) bool
}

func Upgrade(path os.File, p StopPolicy) {
	// TODO: Read up repositoryconfig.
	// TODO: Initialize repository.
	// TODO: Verify repository. (already implemented)

	// TODO: Read up all migrations from filesystem.
	// TODO: Verify all migrations from filesystem.

	// TODO: Find out which migration is the first one we should start with.
	// TODO: Make sure all migrations above that have not been migrated.

	// TODO: Migrate each unfinnished migration until `p` tells us we are done, or we have migrated all of them.
}
