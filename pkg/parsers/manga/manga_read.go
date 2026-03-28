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

func (p *MangaReadParser) ExtractChapterName() (string, error) {
	ch, err := p.client.Get(context.Background(), p.source.String())
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(bytes.NewReader(ch))
	if err != nil {
		return "", err
	}
	return p.htmlManager.FindHtmlContentData(doc, "h1", "id", "chapter-heading"), nil
}

func (p *MangaReadParser) ExtractSingleChapter() ([][]byte, error) {
	downloadInfo, err := p.client.Get(context.Background(), p.source.String())
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(downloadInfo))
	if err != nil {
		return nil, err
	}

	links := []string{}

	i := 1
	for {
		imageName := "image"
		link := p.htmlManager.FindHtmlContent(doc, "img", "id", fmt.Sprintf("%s-%d", imageName, i))
		if link == "" {
			break
		}
		links = append(links, strings.TrimSpace(link))
		i++
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

func (p *MangaReadParser) DownloadPages(pages [][]byte) error {
	content, err := p.imageManager.ConcatPages(pages)
	if err != nil {
		return err
	}

	chn, err := p.ExtractChapterName()
	if err != nil && chn == "" {
		chn = p.source.String()
	}

	if err = p.imageManager.SaveImageInSystem(content, chn); err != nil {
		return err
	}

	return nil
}

var _ MangaParser = &MangaReadParser{}
