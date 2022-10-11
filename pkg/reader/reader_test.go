package reader

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadCSVByLines(t *testing.T) {
	r := require.New(t)

	t.Run("succesfully read csv by lines", func(t *testing.T) {
		urlChan := make(chan string, 4)
		csv := bytes.NewBufferString("1,one\n2,two\n3,three\n4,four\n")
		counter := 0
		var mu sync.Mutex
		go func() {
			for range urlChan {
				mu.Lock()
				defer mu.Unlock()
				counter++
			}
		}()
		ReadCSVByLines(csv, urlChan, 10)
		time.Sleep(2 * time.Second)
		// make sure we processed same number of lines as we have settled up above in testing string
		r.Equal(4, counter)
	})
}
