package workerpool

import (
	"context"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestWorkerPool(t *testing.T) {
	log := logrus.New()
	r := require.New(t)

	t.Run("succesfully handle all tasks", func(t *testing.T) {
		ctx := context.Background()
		tasksCount := 10
		var wg sync.WaitGroup
		wg.Add(tasksCount)
		p := NewPool(ctx, log, 2)
		p.Run()
		for i := 0; i < tasksCount; i++ {
			p.Execute(func() error {
				defer wg.Done()
				return nil
			})
		}
		// it will fail here in case if number of executed taks is not equal to expected number
		wg.Wait()
	})

	t.Run("context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		p := NewPool(ctx, log, 1)
		p.Run()

		cancel()

		stopped := <-p.Stopped
		r.Equal(stopped, struct{}{})
	})

	t.Run("stop pool manually", func(t *testing.T) {
		ctx := context.Background()
		p := NewPool(ctx, log, 1)
		p.Run()

		p.Stop()

		stopped := <-p.Stopped
		r.Equal(stopped, struct{}{})
	})

	t.Run("make sure at least one worker exist", func(t *testing.T) {
		ctx := context.Background()
		var wg sync.WaitGroup
		wg.Add(1)
		// when workers count less than 1 it will always create at least 1 worker
		p := NewPool(ctx, log, 0)
		p.Run()
		p.Execute(func() error {
			defer wg.Done()
			return nil
		})
		// it will fail here in case if number of executed taks is not equal to expected number
		wg.Wait()
	})

	t.Run("succesfully handle 3 of 5 tasks", func(t *testing.T) {
		ctx := context.Background()
		tasksCount := 5
		var wg sync.WaitGroup
		wg.Add(3)
		p := NewPool(ctx, log, 1)
		p.Run()
		for i := 0; i < tasksCount; i++ {
			if i == 3 {
				p.Stop()
			}
			p.Execute(func() error {
				defer wg.Done()
				return nil
			})
		}
		// it will fail here in case if number of executed taks is not equal to expected number
		wg.Wait()
	})

	// TODO: add more test cases

}
