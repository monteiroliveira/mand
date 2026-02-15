package manga

// TODO: Adjust interface content
type MangaParser interface {
	ExtractChapterName() (string, error)
	ExtractSingleChapter() ([][]byte, error)
	DownloadPages(pages [][]byte) error
}
