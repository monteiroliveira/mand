package parsers

// TODO: Adjust interface content
type Parser interface {
	ExtractChapterName() (string, error)
	ExtractSingleChapter() ([][]byte, error)
	DownloadPages(pages [][]byte) error
}
