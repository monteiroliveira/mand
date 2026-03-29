package manga

import "net/url"

type Operation int

const (
	DownloadOperation Operation = iota
	DownloadListOperation
)

// TODO: Adjust interface content
type MangaParser interface {
	ExtractChapterName() (string, error)
	ExtractSingleChapter() ([][]byte, error)
	DownloadPages(pages [][]byte, chapterName string) error
	ExtractChapterList() error
}

type MangaParserArgs struct {
	Source    *url.URL
	Verbose   bool
	Operation Operation
}
