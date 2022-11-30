// Package app provides primitives to initialise crucial files for Lominus.
package app

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	appDir "github.com/beebeeoii/lominus/internal/app/dir"
	appConstants "github.com/beebeeoii/lominus/internal/constants"
	"github.com/beebeeoii/lominus/internal/file"
	logs "github.com/beebeeoii/lominus/internal/log"
	"github.com/boltdb/bolt"
)

// Init initialises and ensures log and preference files that Lominus requires are available.
// Directory in Preferences defaults to empty string ("").
// Frequency in Preferences defaults to -1.
func Init() error {
	baseDir, retrieveBaseDirErr := appDir.GetBaseDir()
	if retrieveBaseDirErr != nil {
		return retrieveBaseDirErr
	}

	if !file.Exists(baseDir) {
		os.Mkdir(baseDir, os.ModePerm)
	}

	dbFName := filepath.Join(baseDir, appConstants.DATABASE_FILE_NAME)
	db, dbErr := bolt.Open(dbFName, 0600, &bolt.Options{Timeout: 3 * time.Second})

	if dbErr != nil {
		return dbErr
	}
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Auth"))
		prefBucket, prefBucketErr := tx.CreateBucketIfNotExists([]byte("Preferences"))
		if prefBucketErr != nil {
			return prefBucketErr
		}

		logLevel := prefBucket.Get([]byte("logLevel"))

		logInitErr := logs.Init(string(logLevel))
		if logInitErr != nil {
			return logInitErr
		}

		return nil
	})

	// TODO Consider moving this to its own module in the future.
	gradesPath := filepath.Join(baseDir, appConstants.GRADES_FILE_NAME)

	if !file.Exists(gradesPath) {
		gradeFileErr := file.EncodeStructToFile(gradesPath, time.Now())

		if gradeFileErr != nil {
			return gradeFileErr
		}
	}

	return err
}

// GetOs returns user's running program's operating system target:
// one of darwin, freebsd, linux, and so on.
// To view possible combinations of GOOS and GOARCH, run "go tool dist list".
func GetOs() string {
	return runtime.GOOS
}
