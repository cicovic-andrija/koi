package server

import (
	"fmt"
	"os"
	"path/filepath"
)

// Run configures the system, builds the database, and boots the server.
func Run() {
	trace(_control, "main: start: %s v0.1", filepath.Base(os.Args[0]))
	readEnvironment()
	buildDatabase()
	_serverControl.boot()
}

// TODO: Implement this function.
func readEnvironment() {
	_serverControl.endpoint = "localhost:8072"
	_database.filePath = "data/koidata.xml"
}

func buildDatabase() error {
	file, err := os.Open(_database.filePath)
	if err != nil {
		panic(fmt.Errorf("failed to open file %s: %v", _database.filePath, err))
	}
	defer file.Close()

	if err = DecodeDatabase(file); err != nil {
		panic(fmt.Errorf("failed to decode database in %s: %v", _database.filePath, err))
	}

	return nil
}
