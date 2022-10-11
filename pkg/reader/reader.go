package reader

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/sirupsen/logrus"
)

func ReadCSVByLines(ctx context.Context, log logrus.FieldLogger, file io.Reader, urlChan chan<- string) error {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// 1st value is the row number so we can read directly 2nd value
		urlChan <- rec[1]

		// Check if the context is expired.
		select {
		default:
		case <-ctx.Done():
			log.Info("stop file parsing")
			return ctx.Err()
		}
	}
	return nil
}
