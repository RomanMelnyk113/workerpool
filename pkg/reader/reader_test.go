package reader

import (
	"bytes"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestProcessTestFile(t *testing.T) {
	log := logrus.New()

	r := require.New(t)

	t.Run("succesfully read csv by lines", func(t *testing.T) {
		chunkSize := 1
		var mu sync.Mutex
		counter := 0
		fn := func(data []string) {
			mu.Lock()
			defer mu.Unlock()
			log.Infof("processsing: %v", data)
			counter++
		}
		csv := bytes.NewBufferString("1,one\n2,two\n3,three\n4,four")
		ReadCSVByLines(chunkSize, csv, fn)

		// make sure we processed same number of lines as we have settled up above in testing string
		r.Equal(counter, 4)
	})
}
