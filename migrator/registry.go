package migrator

import (
	"errors"
	"fmt"
)

var (
	DriverByFilenameExtension map[string]Driver

	UpgraderByFilenameExtension     map[string]Upgrader
	DowngraderByFilenameExtension   map[string]Downgrader
	PreVerifyerByFilenameExtension  map[string]PreVerifyer
	PostVerifyerByFilenameExtension map[string]PostVerifyer
)

var alreadyRegistered = errors.New("Already registered.")

func RegisterUpgrader(u Upgrader) {
	registerDriver(u)
	if e, exist := UpgraderByFilenameExtension[u.FilenameExtension()]; exist {
		panic(fmt.Sprintf("Upgrader with the same extension already registered. Previous: %s New: %s", e, u))
	}
	UpgraderByFilenameExtension[u.FilenameExtension()] = u
}

func RegisterDowngrader(u Downgrader) {
	registerDriver(u)
	if e, exist := DowngraderByFilenameExtension[u.FilenameExtension()]; exist {
		panic(fmt.Sprintf("Downgrader with the same extension already registered. Previous: %s New: %s", e, u))
	}
	DowngraderByFilenameExtension[u.FilenameExtension()] = u
}

func RegisterPreVerifyer(u PreVerifyer) {
	registerDriver(u)
	if e, exist := PreVerifyerByFilenameExtension[u.FilenameExtension()]; exist {
		panic(fmt.Sprintf("PreVerifyer with the same extension already registered. Previous: %s New: %s", e, u))
	}
	PreVerifyerByFilenameExtension[u.FilenameExtension()] = u
}

func RegisterPostVerifyer(u PostVerifyer) {
	registerDriver(u)
	if e, exist := PostVerifyerByFilenameExtension[u.FilenameExtension()]; exist {
		panic(fmt.Sprintf("PostVerifyer with the same extension already registered. Previous: %s New: %s", e, u))
	}
	PostVerifyerByFilenameExtension[u.FilenameExtension()] = u
}

func registerDriver(u Driver) {
	filenameExtension := u.FilenameExtension()

	if filenameExtension == "" {
		panic(fmt.Sprintf("Filename extension must not be empty string for %s.", u))
	}

	if e, exist := DriverByFilenameExtension[filenameExtension]; exist && e != u {
		panic(fmt.Sprintf("A different driver has been registered with the same file extension. Previous: %s New: %s", e, u))
	}

	DriverByFilenameExtension[filenameExtension] = u
}
