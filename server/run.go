package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Run configures the system, builds the database, and boots the server.
func Run() {
	trace(_control, "main: start: %s v1.0", filepath.Base(os.Args[0]))
	readEnvironment()
	buildDatabase()
	_serverControl.boot()
}

func readEnvironment() {
	const (
		modeEnvVar = "KOIPOND_MODE"
		portEnvVar = "KOIPOND_PORT"
	)

	mode := os.Getenv(modeEnvVar)
	trace(_env, "%s = %q", modeEnvVar, mode)
	if mode == "" {
		mode = "prod"
	}

	if mode == "dev" {
		_serverControl.endpoint = "localhost:8072"
	} else if mode == "prod" {
		port := os.Getenv(portEnvVar)
		trace(_env, "%s = %q", portEnvVar, port)
		if num, err := strconv.Atoi(port); err != nil || num < 1 || num > 65535 {
			trace(_error, "value of %s is invalid or is not a valid TCP port number", portEnvVar)
			os.Exit(1)
		}
		_serverControl.endpoint = "127.0.0.1:" + port
	} else {
		trace(_error, "value of %s is invalid", modeEnvVar)
		os.Exit(1)
	}
	trace(_control, "main: in mode %q (HTTP): endpoint will be http://%s", mode, _serverControl.endpoint)

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
