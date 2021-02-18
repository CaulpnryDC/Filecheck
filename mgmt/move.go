package mgmt

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func MoveFiles(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		log.Error("Couldn't open source file")
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		log.Error("Couldn't open destination file")
	}
	defer func() {
		err := outputFile.Close()
		if err != nil {
			log.Error(err)
		}
	} ()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		log.Error("Writing to output file failed")
	}

	err = os.Remove(sourcePath)
	if err != nil {
		log.Error("Failed removing original file")
	}
	return nil
}
