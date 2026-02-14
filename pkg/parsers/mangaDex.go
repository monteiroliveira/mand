package parsers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/internal/scraper"
	"golang.org/x/net/html"
)

var (
	MangaDexValidLink               string = "mangadex.org"
	mangaDexChapterDownloadEndpoint string = "https://api.mangadex.org/at-home/server/"
)

type MangaDexParser struct {
	source       *url.URL
	client       *scraper.HttpClient
	imageManager *internal.ImageManager
	htmlManager  *scraper.HtmlManager
}

type mangaDexChapterDownloadInfo struct {
	Result  string `json:"result"`
	BaseUrl string `json:"baseUrl"`
	Chapter struct {
		Hash string   `json:"hash"`
		Data []string `json:"data"`
	} `json:"chapter"`
}

func NewMangaDexParser(source *url.URL) *MangaDexParser {
	return &MangaDexParser{
		source:       source,
		client:       scraper.NewHttpClient(),
		imageManager: internal.NewImageManager(),
		htmlManager:  scraper.NewHtmlManager(),
	}
}

func getChapterId(url *url.URL) (string, error) {
	pathContent := strings.Split(url.Path, "/")
	if len(pathContent) <= 1 {
		return "", nil
	}

	chapterId := pathContent[len(pathContent)-1]

	return chapterId, nil
}

func buildDownloadList(downloadInfo *mangaDexChapterDownloadInfo) []string {
	downloadList := []string{}
	for _, imageLink := range downloadInfo.Chapter.Data {
		downloadLink := downloadInfo.BaseUrl + "/data/" + downloadInfo.Chapter.Hash + "/" + imageLink
		downloadList = append(downloadList, downloadLink)
	}

	return downloadList
}

func (p *MangaDexParser) ExtractChapterName() (string, error) {
	chapter, err := p.client.Get(context.Background(), p.source.String())
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(bytes.NewReader(chapter))
	if err != nil {
		return "", err
	}
	return p.htmlManager.FindHtmlContent(doc, "meta", "property", "og:title"), nil
}

// TODO: create a flow to parse novels and mangas
func (p *MangaDexParser) ExtractSingleChapter() ([][]byte, error) {
	chapterId, err := getChapterId(p.source)
	if err != nil {
		return nil, err
	}

	downloadEndpoint := mangaDexChapterDownloadEndpoint + chapterId
	downloadInfo, err := p.client.Get(context.Background(), downloadEndpoint)
	if err != nil {
		return nil, err
	}

	var chapterDownloadInfo mangaDexChapterDownloadInfo
	if err = json.Unmarshal(downloadInfo, &chapterDownloadInfo); err != nil {
		return nil, err
	}

	downloadList := buildDownloadList(&chapterDownloadInfo)
	var pages [][]byte
	for _, downloadLink := range downloadList {
		page, err := p.client.Get(context.Background(), downloadLink)
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (p *MangaDexParser) DownloadPages(pages [][]byte) error {
	content, err := p.imageManager.ConcatPages(pages)
	if err != nil {
		return err
	}

	chn, err := p.ExtractChapterName()
	if err != nil && chn != "" {
		chn, err = getChapterId(p.source)
		if err != nil {
			chn = "mand_manga"
		}
	}

	if err = p.imageManager.SaveImageInSystem(content, chn); err != nil {
		return err
	}

	return nil
}

var _ Parser = &MangaDexParser{}
