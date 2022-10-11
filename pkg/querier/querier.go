package querier

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Stats struct {
	responseTime time.Duration
	responseSize int
}

func GetAndPrintPage(ctx context.Context, log logrus.FieldLogger, url string) (*Stats, error) {
	url = "https://" + url
	log.Infof("proccessing %v\n", url)
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("%s finished in %v, body size %v\n", url, elapsed, len(body))
	stats := &Stats{responseTime: elapsed, responseSize: len(body)}
	return stats, nil
}
