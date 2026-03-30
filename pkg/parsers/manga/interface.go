package manga

import (
	"fmt"
	"sync"
)

type Operation int
type ExtractorFunc func(link string, wg *sync.WaitGroup, ch chan error)

const (
	DownloadOperation Operation = iota
	DownloadListOperation
)

type MangaParser interface {
	ExtractChapterName() (string, error)
	ExtractChapterContent() ([][]byte, error)
	DownloadPages(pages [][]byte, chapterName string) error
	ExtractChapterContentByList(batchSize int) error
}

func ExtractList(
	links []string, batchSize int, extractor ExtractorFunc, wg *sync.WaitGroup, ch chan error,
) error {
	batch := 0
	if batchSize <= 0 {
		return fmt.Errorf("Invalid batch size")
	}
	for _, link := range links {
		wg.Add(1)
		go extractor(link, wg, ch)
		batch++
		if batch == batchSize {
			wg.Wait()
			batch = 0
		}
	}
	if batch != 0 {
		wg.Wait()
	}
	return nil
}
