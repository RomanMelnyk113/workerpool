package querier

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetAndPrintPage(url string) error {
	url = "https://" + url
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = string(body)
	log.Printf("%s processed in %v, body size %v\n", url, elapsed, len(body))
	//log.Printf(sb)
	return nil
}
