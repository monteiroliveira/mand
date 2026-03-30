package manga

import (
	"bytes"
	"context"
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
	client       *scraper.HttpClient
	regex        *scraper.RegexParser
	imageManager *internal.ImageManager
	htmlManager  *scraper.HtmlManager
}

func NewMangaReadParser(args *MangaParserArgs) *MangaReadParser {
	return &MangaReadParser{
		args:         args,
		client:       scraper.NewHttpClient(),
		regex:        scraper.NewRegexParser(),
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
		return p.args.Source.String(), nil
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

	pages, err := p.buildContentPages(links)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (p *MangaReadParser) ExtractChapterList() error {
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
		return fmt.Errorf("cant extract manga links")
	}

	ch := make(chan error)
	var wg sync.WaitGroup

	batch := 0
	for _, link := range links {
		wg.Add(1)
		go func() {
			fmt.Printf("Init job in link: %s\n", link)
			chp, err := p.extractChapter(link)
			if err != nil {
				ch <- err
			}
			chn, err := p.getChapterName(link)
			if err != nil || chn == "" {
				chn = p.args.Source.String()
			}
			if err = p.DownloadPages(chp, chn); err != nil {
				ch <- err
			}
			fmt.Println("Done")
			wg.Done()
		}()
		batch++
		if batch == 5 {
			wg.Wait()
			batch = 0
		}
	}

	wg.Wait()
	close(ch)

	return nil
}

func (p *MangaReadParser) ExtractChapterName() (string, error) {
	return p.getChapterName(p.args.Source.String())
}

func (p *MangaReadParser) ExtractSingleChapter() ([][]byte, error) {
	return p.extractChapter(p.args.Source.String())
}

func (p *MangaReadParser) DownloadPages(pages [][]byte, chapterName string) error {
	if err := p.imageManager.SavePdfInSystem(pages, chapterName); err != nil {
		return err
	}

	return nil
}

var _ MangaParser = &MangaReadParser{}
