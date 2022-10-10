package reader

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"sync"
)

func ReadCSVByLines(file io.Reader, urlChan chan<- string, limit int) error {
	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	i := 0
	for {
		if limit > 0 && i > limit {
			break
		}
		i++
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// 1st value is the row number so we can read directly 2nd value
		urlChan <- rec[1]
	}
	return nil
}

func rawReadFileByChunk(file io.Reader, fn func(string)) error {
	var wg sync.WaitGroup
	r := bufio.NewReader(file)
	for {
		buf := make([]byte, 4*1024) //the chunk size
		n, err := r.Read(buf)       //loading chunk into buffer
		buf = buf[:n]
		if n == 0 {
			if err != nil {
				log.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return err
		}
		nextUntillNewline, err := r.ReadBytes('\n')
		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}
		wg.Add(1)
		go func() {
			fn(string(buf))
			wg.Done()
		}()
		break
	}
	wg.Wait()
	return nil
}
