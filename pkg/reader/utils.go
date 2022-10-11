package reader

import (
	"archive/zip"
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// GetTestFileByChunks downloading testing file and read all files from zip
func ProcessTestFile(ctx context.Context, log logrus.FieldLogger, processingChan chan<- string) error {
	fileUrl := "https://s3.amazonaws.com/alexa-static/top-1m.csv.zip"
	log.Printf("downloading file %v\n", fileUrl)
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// assume we are dealing with not a big files
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}
	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		log.Println("reading file:", zipFile.Name)
		err := readZipFile(ctx, log, processingChan, zipFile)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func readZipFile(ctx context.Context, log logrus.FieldLogger, processingChan chan<- string, zf *zip.File) error {
	f, err := zf.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	return ReadCSVByLines(ctx, log, f, processingChan)

}
