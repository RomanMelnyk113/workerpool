package workerpool

import (
	"context"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWorkerPool(t *testing.T) {
	ctx := context.Background()
	log := logrus.New()

	t.Run("succesfully handle all tasks", func(t *testing.T) {
		tasksCount := 10
		var wg sync.WaitGroup
		wg.Add(tasksCount)
		p := NewPool(ctx, log, 2)
		for i := 0; i < tasksCount; i++ {
			p.Execute(func() error {
				defer wg.Done()
				return nil
			})
		}
		// it will fail here in case if number of executed taks is not equal to expected number
		wg.Wait()
	})
}
