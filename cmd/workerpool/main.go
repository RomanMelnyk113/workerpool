package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanMelnyk113/workerpool/internal/workerpool"
	"github.com/RomanMelnyk113/workerpool/pkg/querier"
	"github.com/RomanMelnyk113/workerpool/pkg/reader"
	"github.com/sirupsen/logrus"
)

func main() {
	size := flag.Int("size", 100, "workers pool size")

	log := logrus.New()

	if err := run(log, *size); err != nil {
		log.Fatalf("workerpool exited with error: %v", err)
	}

	log.Info("workerpool shutdown")
}

func run(log logrus.FieldLogger, size int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := workerpool.NewPool(ctx, log, size)
	urlChan := make(chan string)
	go reader.GetTestFileByChunks(urlChan)

	go func() {
		for url := range urlChan {
			url := url
			wp.Execute(func() error {
				return querier.GetAndPrintPage(url)
			})
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	select {
	//case url := <-urlChan:
	//log.Info("hanle URL", url)
	//wp.Execute(func() error {
	//return querier.GetAndPrintPage(url)
	//})
	case <-signals:
		log.Info("interrupt signal received, initiating workerpool shutdown")
		cancel()
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
