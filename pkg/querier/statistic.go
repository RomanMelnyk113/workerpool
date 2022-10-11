package querier

import (
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
	s.AvgResponseSize = (s.AvgResponseSize + stats.responseSize) / s.Total
	avgTime := (s.AvgResponseTime.Milliseconds() + stats.responseTime.Milliseconds()) / int64(s.Total)
	s.AvgResponseTime = time.Duration(avgTime)
}
