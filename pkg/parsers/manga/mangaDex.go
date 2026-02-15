// Manga Dex manga parser, extracting the manga content from the API and
// the manga title from the chapter page via Web Scrap.
//
// Starting with Manga Dex as the first parser because they have a good
// API to work with, and I'm too lazy to fully scrap a site for now.
// REFERENCE: https://api.mangadex.org/docs/?ref=public_apis
package manga

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/internal/scraper"
	"github.com/monteiroliveira/mand/internal/scraper/models"
	"golang.org/x/net/html"
)

var (
	MangaDexValidLink      string = "mangadex.org"
	mangaDexChDownEndpoint string = "https://api.mangadex.org/at-home/server/"
)

type MangaDexParser struct {
	source       *url.URL
	chapterId    string
	client       *scraper.HttpClient
	imageManager *internal.ImageManager
	htmlManager  *scraper.HtmlManager
}

func NewMangaDexParser(source *url.URL) *MangaDexParser {
	return &MangaDexParser{
		source:       source,
		chapterId:    "",
		client:       scraper.NewHttpClient(),
		imageManager: internal.NewImageManager(),
		htmlManager:  scraper.NewHtmlManager(),
	}
}

func getChapterId(url *url.URL) (string, error) {
	pathContent := strings.Split(url.Path, "/")
	if len(pathContent) <= 1 {
		return "", fmt.Errorf("Failed to extract chapter id")
	}

	chapterId := pathContent[len(pathContent)-1]

	return chapterId, nil
}

func buildDownloadList(downloadInfo *models.MangaDexChapterDownloadInfo) []string {
	downloadList := []string{}
	for _, imageLink := range downloadInfo.Chapter.Data {
		downloadLink := downloadInfo.BaseUrl + "/data/" + downloadInfo.Chapter.Hash + "/" + imageLink
		downloadList = append(downloadList, downloadLink)
	}

	return downloadList
}

func (p *MangaDexParser) buildContentPages(links []string) ([][]byte, error) {
	var pages [][]byte
	for _, link := range links {
		page, err := p.client.Get(context.Background(), link)
		if err != nil {
			// build up error list with join
			continue
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (p *MangaDexParser) ExtractChapterName() (string, error) {
	ch, err := p.client.Get(context.Background(), p.source.String())
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(bytes.NewReader(ch))
	if err != nil {
		return "", err
	}
	return p.htmlManager.FindHtmlContent(doc, "meta", "property", "og:title"), nil
}

// TODO: create a flow to parse novels and mangas
func (p *MangaDexParser) ExtractSingleChapter() ([][]byte, error) {
	chapterId, err := getChapterId(p.source)
	if err != nil && chapterId == "" {
		return nil, err
	}
	p.chapterId = chapterId

	downloadEndpoint := mangaDexChDownEndpoint + chapterId
	downloadInfo, err := p.client.Get(context.Background(), downloadEndpoint)
	if err != nil {
		return nil, err
	}

	var chDownInfo models.MangaDexChapterDownloadInfo
	if err = json.Unmarshal(downloadInfo, &chDownInfo); err != nil {
		return nil, err
	}

	links := buildDownloadList(&chDownInfo)
	pages, err := p.buildContentPages(links)

	// TODO: Show failed pages download

	return pages, nil
}

func (p *MangaDexParser) DownloadPages(pages [][]byte) error {
	content, err := p.imageManager.ConcatPages(pages)
	if err != nil {
		return err
	}

	chn, err := p.ExtractChapterName()
	if err != nil && chn != "" {
		chn = p.chapterId
	}

	if err = p.imageManager.SaveImageInSystem(content, chn); err != nil {
		return err
	}

	return nil
}

var _ MangaParser = &MangaDexParser{}
