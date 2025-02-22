/*
controls embedded postgres database via pg_ctl
sets locale for messages (LC_MESSAGES) for parsing call outputs
*/
package embedded

import (
	"path/filepath"
	"r3/config"
	"runtime"
)

var (
	// paths
	dbBin    string // database binary
	dbBinCtl string // database control binary
	dbData   string // database data directory
	locale   string // locale setting for database messages

	// database state messages
	msgState0  string = "is stopped"
	msgState1  string = "is running"
	msgStarted string = "started"
	msgStopped string = "stopped"
)

func GetDbBinPath() string {
	return dbBin
}

// SetPaths sets the paths for the embedded database
func SetPaths() {
	dbBin = config.File.Paths.EmbeddedDbBin
	dbData = config.File.Paths.EmbeddedDbData

	if runtime.GOOS == "windows" {
		dbBinCtl = filepath.Join(dbBin, "pg_ctl.exe")
	} else {
		dbBinCtl = filepath.Join(dbBin, "pg_ctl")
	}
}
