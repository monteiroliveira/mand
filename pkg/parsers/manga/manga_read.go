package manga

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"github.com/monteiroliveira/mand/internal"
	"github.com/monteiroliveira/mand/internal/scraper"
)

var (
	MangaReadValidLink string = "www.mangaread.org"
)

type MangaReadParser struct {
	args         *MangaParserArgs
	log          internal.Logger
	client       *scraper.HttpClient
	regex        *scraper.RegexParser
	imageManager *internal.ImageManager
	htmlManager  *scraper.HtmlManager
	errorChan    chan error
}

func NewMangaReadParser(args *MangaParserArgs) *MangaReadParser {
	return &MangaReadParser{
		args:         args,
		log:          args.Log,
		client:       scraper.NewHttpClient(),
		regex:        scraper.NewRegexParser(),
		imageManager: internal.NewImageManager(),
		htmlManager:  scraper.NewHtmlManager(),
	}
}

func (p *MangaReadParser) buildContentPages(links []string) ([][]byte, error) {
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
		return p.args.Source.String(), nil
	}
	return chapterName, nil
}

func (p *MangaReadParser) extractChapterContent(source string) ([][]byte, error) {
	p.log.Trace("Fetching chapter page: %s", source)
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
		link := p.htmlManager.FindHtmlContent(doc, "img", "id", fmt.Sprintf("^image-%d$", i))
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

	p.log.Info("Downloading %d pages", len(links))
	pages, err := p.buildContentPages(links)
	if err != nil {
		return nil, err
	}
	p.log.Debug("Downloaded %d/%d pages", len(pages), len(links))
	return pages, nil
}

func (p *MangaReadParser) ExtractChapterContentByList(batchSize int) error {
	p.log.Trace("Fetching manga page: %s", p.args.Source.String())
	downloadInfo, err := p.client.Get(context.Background(), p.args.Source.String())
	if err != nil {
		return err
	}

	doc, err := html.Parse(bytes.NewReader(downloadInfo))
	if err != nil {
		return err
	}

	urlTemplate, err := p.regex.Normalize(p.args.Source.String())
	if err != nil {
		return err
	}

	links := p.htmlManager.ListHtmlContent(doc, "a", "href", fmt.Sprintf("%schapter.*", urlTemplate))
	if len(links) == 0 {
		return fmt.Errorf("Cant extract manga links, link list is empty")
	}
	p.log.Info("Found %d chapters (batch size: %d)", len(links), batchSize)

	wg := new(sync.WaitGroup)
	extractor := func(link string, wg *sync.WaitGroup, ch chan error) {
		defer wg.Done()

		p.log.Debug("Extracting chapter: %s", link)
		errFmt := fmt.Errorf("Error in link content extraction: %s", link)
		chp, err := p.extractChapterContent(link)
		if err != nil {
			ch <- errors.Join(errFmt, err)
			return
		}
		chn, err := p.getChapterName(link)
		if err != nil || chn == "" {
			chn = p.args.Source.String()
		}
		if err = p.DownloadPages(chp, chn); err != nil {
			ch <- errors.Join(errFmt, err)
			return
		}
		p.log.Debug("Chapter done: %s", chn)
	}

	if err = ExtractList(links, batchSize, extractor, wg, p.args.ErrorChan); err != nil {
		return err
	}
	return nil
}

func (p *MangaReadParser) ExtractChapterName() (string, error) {
	return p.getChapterName(p.args.Source.String())
}

func (p *MangaReadParser) ExtractChapterContent() ([][]byte, error) {
	return p.extractChapterContent(p.args.Source.String())
}

func (p *MangaReadParser) DownloadPages(pages [][]byte, chapterName string) error {
	p.log.Info("Saving PDF: %s.pdf", chapterName)
	if err := p.imageManager.SavePdfInSystem(pages, chapterName); err != nil {
		return err
	}
	p.log.Debug("PDF saved: %s.pdf (%d pages)", chapterName, len(pages))
	return nil
}

var _ MangaParser = &MangaReadParser{}
