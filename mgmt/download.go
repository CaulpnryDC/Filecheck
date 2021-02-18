package mgmt

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	} ()

	out, err := os.Create(filepath)
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err := out.Close()
		if err != nil {
			log.Fatal()
		}
	} ()

	_, err = io.Copy(out, resp.Body)
	return err
}