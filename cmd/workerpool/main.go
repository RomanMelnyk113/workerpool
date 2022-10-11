package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanMelnyk113/workerpool/internal/workerpool"
	"github.com/RomanMelnyk113/workerpool/pkg/querier"
	"github.com/RomanMelnyk113/workerpool/pkg/reader"
	"github.com/sirupsen/logrus"
)

func main() {
	workersCount := flag.Int("workers", 10, "workers pool size")
	tasksLimit := flag.Int("tasks", 30, "max allowed tasks number to be executed")
	flag.Parse()

	log := logrus.New()

	if err := run(log, *workersCount, *tasksLimit); err != nil {
		log.Fatalf("workerpool exited with error: %v", err)
	}

	log.Info("workerpool shutdown")
}

// run will start listening for new incoming tasks to be executed by workers
func run(log logrus.FieldLogger, size, limit int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	totalStats := &querier.TotalStats{}

	log.Infof("start pool with %d and tasks limit %d\n", size, limit)
	wp := workerpool.NewPool(ctx, log, size)
	urlChan := make(chan string)

	// simulate streaming by reading the TOP 1m sites list by chunks
	go reader.ProcessTestFile(ctx, log, urlChan)

	// process received urls and send tasks to the worker pool
	go func() {
		i := 0
		for url := range urlChan {
			if i >= limit {
				log.Info("reached limit, stopping tasks processing")
				break
			}
			url := url
			wp.Execute(func() error {
				stats, err := querier.GetAndPrintPage(ctx, log, url)
				totalStats.UpdateStats(stats)
				return err
			})
			i++
		}
		wp.Stop()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-wp.Stopped:
		log.Info("workerpool stopped, printing summary")
		printSummary(totalStats)
	case <-signals:
		log.Info("interrupt signal received, initiating workerpool shutdown")
		cancel()
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func printSummary(stats *querier.TotalStats) {
	fmt.Println("==================================")
	fmt.Printf(
		"Total tasks: %d\nSuccess: %d\nFailure: %d\nAverage body size: %v bytes\nAverage response time: %v\n",
		stats.Total, stats.Succeed, stats.Failed, stats.AvgResponseSize, stats.AvgResponseTime,
	)
	fmt.Println("==================================")
}
