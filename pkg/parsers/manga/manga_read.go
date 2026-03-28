package manga

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/internal/scraper"
)

var (
	MangaReadValidLink string = "www.mangaread.org"
)

type MangaReadParser struct {
	source       *url.URL
	client       *scraper.HttpClient
	imageManager *internal.ImageManager
	htmlManager  *scraper.HtmlManager
}

func NewMangaReadParser(source *url.URL) *MangaReadParser {
	return &MangaReadParser{
		source:       source,
		client:       scraper.NewHttpClient(),
		imageManager: internal.NewImageManager(),
		htmlManager:  scraper.NewHtmlManager(),
	}
}

func (p *MangaReadParser) buildContentPages(links []string) ([][]byte, error) {
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

func (p *MangaReadParser) getChapterName(source string) (string, error) {
	ch, err := p.client.Get(context.Background(), source)
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(bytes.NewReader(ch))
	if err != nil {
		return "", err
	}
	chapterName := p.htmlManager.FindHtmlContentData(doc, "h1", "id", "chapter-heading")
	if chapterName == "" {
		return p.source.String(), nil
	}
	return chapterName, nil
}

func (p *MangaReadParser) extractChapter(source string) ([][]byte, error) {
	downloadInfo, err := p.client.Get(context.Background(), source)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(downloadInfo))
	if err != nil {
		return nil, err
	}

	links := []string{}

	for i := 0; ; i++ {
		imageName := "image"
		link := p.htmlManager.FindHtmlContent(doc, "img", "id", fmt.Sprintf("%s-%d", imageName, i))
		if link == "" && i != 0 {
			break
		}
		if link != "" {
			links = append(links, strings.TrimSpace(link))
		}
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("cant extract manga links")
	}

	pages, err := p.buildContentPages(links)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (p *MangaReadParser) ExtractChapterList() error {
	downloadInfo, err := p.client.Get(context.Background(), p.source.String())
	if err != nil {
		return err
	}

	chaptersLinks := []string{}
	doc, err := html.Parse(bytes.NewReader(downloadInfo))
	if err != nil {
		return err
	}

	for i := 0; ; i++ {
		link := p.htmlManager.FindHtmlContent(doc, "a", "href", fmt.Sprintf("%schapter-%d/", p.source.String(), i))
		if link == "" && i != 0 {
			break
		}
		if link != "" {
			chaptersLinks = append(chaptersLinks, strings.TrimSpace(link))
		}
	}
	if len(chaptersLinks) == 0 {
		return fmt.Errorf("cant extract manga links")
	}

	for _, link := range chaptersLinks {
		chp, err := p.extractChapter(link)
		if err != nil {
			// TODO: Log error
			continue
		}
		chn, err := p.getChapterName(link)
		if err != nil || chn == "" {
			chn = p.source.String()
		}
		if err = p.DownloadPages(chp, chn); err != nil {
			continue
		}
	}
	return nil
}

func (p *MangaReadParser) ExtractChapterName() (string, error) {
	return p.getChapterName(p.source.String())
}

func (p *MangaReadParser) ExtractSingleChapter() ([][]byte, error) {
	return p.extractChapter(p.source.String())
}

func (p *MangaReadParser) DownloadPages(pages [][]byte, chapterName string) error {
	if err := p.imageManager.SavePdfInSystem(pages, chapterName); err != nil {
		return err
	}

	return nil
}

var _ MangaParser = &MangaReadParser{}
