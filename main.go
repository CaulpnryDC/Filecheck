package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/akamensky/argparse"
	fc "github.com/caulpnryDC/filecheck/mgmt"
	"github.com/codingsince1985/checksum"
	log "github.com/sirupsen/logrus"
)

func main() {

	start := time.Now()
	fmt.Println(start)

	logFileName := "/var/log/file-check.json"
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.JSONFormatter)
	Formatter.TimestampFormat = "01-02-2006 15:04:05"
	log.SetFormatter(Formatter)

	if err != nil {
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}

	parser := argparse.NewParser("print", "prints string")

	l := parser.String("l", "string", &argparse.Options{Required: true, Help: "File check to run or download"})

	err = parser.Parse(os.Args)
	if err != nil {
		log.Error(parser.Usage(err))
		os.Exit(1)
	}

	if *l == "pd" {
		*l = "https://s3.amazonaws.com/pdpartner/PagerDuty+Outgoing+Numbers.vcf"
		if err := fc.DownloadFile("/app/files/newCard.vcf", *l); err != nil {
			log.Panic(err)
		} else {
			log.Info("Newest PagerDuty file downloaded")
		}
		staticCard := "/app/files/staticCard.vcf"
		staticCardMD5, _ := checksum.MD5sum(staticCard)
		log.Info("Static file: ", staticCardMD5)

		newCard := "/app/files/newCard.vcf"
		newCardMD5, _ := checksum.MD5sum(newCard)
		if err != nil {
			log.Fatal("File is missing: ", err)
		} else {
			log.Info("New file: ", newCardMD5)

			if staticCardMD5 != newCardMD5 {
				log.Error("Files MD5 sum does not match, new file has changes")
				err := fc.MoveFiles(staticCard, "/app/files/old/replaced.vcf")
				err = fc.MoveFiles(newCard, "/app/files/staticCard.vcf")

				alertURL := os.Getenv("alertURL")
				var jsonStr = []byte(`alert info here`)

				req, err := http.NewRequest("POST", alertURL, bytes.NewBuffer(jsonStr))
				req.Header.Set("Content-type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer func() {
					err := resp.Body.Close()
					if err != nil {
						log.Error(err)
					}
				}()

				log.Info("response status: ", resp.Status)
				log.Info("response header: ", resp.Header)
				body, _ := ioutil.ReadAll(resp.Body)
				log.Info("response body: ", string(body))
			} else {
				log.Info("Files are the same, no changes have been made")
			}
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
