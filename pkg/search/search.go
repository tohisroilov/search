package search

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

// Result struct
type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

// All ...
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}

	for _, filePath := range files {
		wg.Add(1)
		go func(ctx context.Context, filePath string, ch chan<- []Result) {
			defer wg.Done()

			file, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Print(err)
				return
			}

			result := []Result{}
			lines := strings.Split(string(file), "\n")
			for ln, line := range lines {
				if strings.Contains(line, phrase) {
					res := Result{
						Phrase:  phrase,
						Line:    line,
						LineNum: int64(ln + 1),
						ColNum:  int64(strings.Index(line, phrase)) + 1,
					}
					result = append(result, res)
				}
			}
			if len(result) > 0 {
				ch <- result
			}
		}(ctx, filePath, ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}
