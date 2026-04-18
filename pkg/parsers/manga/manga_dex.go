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
	log          internal.Logger
	client       *scraper.HttpClient
	imageManager *internal.ImageManager
	htmlManager  *scraper.HtmlManager
}

func NewMangaDexParser(args *MangaParserArgs) *MangaDexParser {
	return &MangaDexParser{
		source:       args.Source,
		chapterId:    "",
		log:          args.Log,
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
	for i, link := range links {
		p.log.Trace("Fetching page %d/%d: %s", i+1, len(links), link)
		page, err := p.client.Get(context.Background(), link)
		if err != nil {
			p.log.Debug("Failed to fetch page %d: %s", i+1, err)
			continue
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (p *MangaDexParser) ExtractChapterName() (string, error) {
	p.log.Trace("Scraping chapter name from: %s", p.source.String())
	ch, err := p.client.Get(context.Background(), p.source.String())
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(bytes.NewReader(ch))
	if err != nil {
		return "", err
	}
	chapterName := p.htmlManager.FindHtmlContent(doc, "meta", "property", "^og:title$")
	if chapterName == "" {
		p.log.Debug("Chapter name not found, falling back to source URL")
		return p.source.String(), nil
	}
	p.log.Debug("Chapter name: %s", chapterName)
	return chapterName, nil
}

// TODO: create a flow to parse novels and mangas
func (p *MangaDexParser) ExtractChapterContent() ([][]byte, error) {
	chapterId, err := getChapterId(p.source)
	if err != nil && chapterId == "" {
		return nil, err
	}
	p.chapterId = chapterId
	p.log.Debug("Extracted chapter ID: %s", chapterId)

	downloadEndpoint := mangaDexChDownEndpoint + chapterId
	p.log.Trace("Fetching download info from: %s", downloadEndpoint)
	downloadInfo, err := p.client.Get(context.Background(), downloadEndpoint)
	if err != nil {
		return nil, err
	}

	var chDownInfo models.MangaDexChapterDownloadInfo
	if err = json.Unmarshal(downloadInfo, &chDownInfo); err != nil {
		return nil, err
	}

	links := buildDownloadList(&chDownInfo)
	p.log.Info("Downloading %d pages", len(links))
	pages, err := p.buildContentPages(links)
	p.log.Debug("Downloaded %d/%d pages", len(pages), len(links))

	return pages, nil
}

func (p *MangaDexParser) DownloadPages(pages [][]byte, chapterName string) error {
	p.log.Info("Saving PDF: %s.pdf", chapterName)
	if err := p.imageManager.SavePdfInSystem(pages, chapterName); err != nil {
		return err
	}
	p.log.Debug("PDF saved: %s.pdf (%d pages)", chapterName, len(pages))
	return nil
}

func (p *MangaDexParser) ExtractChapterContentByList(batchSize int) error {
	return fmt.Errorf("Not supported")
}

var _ MangaParser = &MangaDexParser{}
