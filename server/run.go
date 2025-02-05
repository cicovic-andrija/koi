package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Run configures the system, builds the database, and boots the server.
func Run() {
	trace(_control, "main: start: %s v1.3", filepath.Base(os.Args[0]))
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
	} else if mode == "prod" || mode == "prod-local-listener" {
		port := os.Getenv(portEnvVar)
		trace(_env, "%s = %q", portEnvVar, port)
		if num, err := strconv.Atoi(port); err != nil || num < 1 || num > 65535 {
			panic(fmt.Errorf("value of %s is invalid or is not a valid TCP port number", portEnvVar))
		}
		if mode == "prod-local-listener" {
			_serverControl.endpoint = "127.0.0.1:" + port
		} else {
			_serverControl.endpoint = "0.0.0.0:" + port
		}
	} else {
		panic(fmt.Errorf("value of %s is invalid", modeEnvVar))
	}
	trace(_control, "main: in mode %q (HTTP): endpoint will be http://%s", mode, _serverControl.endpoint)

	// always read koidata.xml from store/ relative to the working directory
	// <#hardcoded#>
	if abs, err := filepath.Abs("store/koidata.xml"); err != nil {
		panic(fmt.Errorf("failed to compose full store path: %v", err))
	} else {
		_database.filePath = abs
	}
}

func buildDatabase() error {
	trace(_decoder, "decoding %s", _database.filePath)
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
