package server

import (
	"fmt"
	"os"
)

func Run() {
	err := buildDatabase()
	if err != nil {
		trace(_control, "main: %v", err)
	}
}

func buildDatabase() error {
	filename := "data/koidata.xml"
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	if err = DecodeDatabase(file); err != nil {
		return fmt.Errorf("failed to decode database in %s: %v", filename, err)
	}

	return nil

}
