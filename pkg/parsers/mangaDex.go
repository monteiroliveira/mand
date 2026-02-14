package parsers

var MangaDexValidLink string = "mangadex.org"

type MangaDexParser struct{}

func NewMangaDexParser() *MangaDexParser {
	return &MangaDexParser{}
}

func (p *MangaDexParser) Extract() ([]byte, error) {
	return nil, nil
}

func (p *MangaDexParser) Download() error {
	return nil
}

var _ Parser = &MangaDexParser{}
