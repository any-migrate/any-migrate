package repository

import (
	"fmt"
)

var RepositoryByScheme map[string]Repository

func Register(r Repository) {
	if e, exist := RepositoryByScheme[r.Scheme()]; exist {
		panic(fmt.Sprintf("Repository with scheme '%s' already registered. Previous: %s New: %s", r.Scheme(), e, r))
	}
	RepositoryByScheme[r.Scheme()] = r
}
