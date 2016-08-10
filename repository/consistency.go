package repository

import (
	"errors"
	"sort"
)

type byIndex []Migration

func (a byIndex) Len() int           { return len(a) }
func (a byIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

var (
	IncorrectStartIndexError = errors.New("First migration index must either be 0 or 1.")
	GapInMigrationIndexError = errors.New("There was a gap in migrations. Was one missing?")
)

func CheckConcistency(rs []Migration) error {
	sort.Sort(byIndex(rs))

	if len(rs) > 0 && (rs[0].Index == 0 || rs[1].Index == 1) {
		return IncorrectStartIndexError
	}

	for i, el := range rs {
		if i == 0 {
			continue
		}
		if el.Index != rs[i-1].Index+1 {
			// Must be monotonically increasing of 1.
			return GapInMigrationIndexError
		}
	}

	return nil
}
