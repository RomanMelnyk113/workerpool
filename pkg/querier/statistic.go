package querier

import (
	"fmt"
	"sync"
	"time"
)

type TotalStats struct {
	mu sync.Mutex

	Total   int
	Succeed int
	Failed  int

	AvgResponseTime time.Duration
	AvgResponseSize int
}

func (s *TotalStats) UpdateStats(stats *Stats) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Total++

	if stats == nil {
		s.Failed++
		return
	}

	s.Succeed++
	if s.AvgResponseTime.Milliseconds() != 0 {
		avgTime := (s.AvgResponseTime.Milliseconds() + stats.responseTime.Milliseconds()) / int64(2)
		s.AvgResponseTime = time.Duration(avgTime) * time.Millisecond
	} else {
		s.AvgResponseTime = stats.responseTime
	}
	if s.AvgResponseTime != 0 {
		s.AvgResponseSize = (s.AvgResponseSize + stats.responseSize) / s.Total
	} else {
		s.AvgResponseSize = stats.responseSize
	}
}

func (s *TotalStats) PrintSummary() {
	fmt.Println("==================================")
	fmt.Printf(
		`Total tasks: %d
		Success: %d
		Failure: %d
		Average body size: %v bytes
		Average response time: %v
		`,
		s.Total, s.Succeed, s.Failed, s.AvgResponseSize, s.AvgResponseTime,
	)
	fmt.Println("==================================")
}
