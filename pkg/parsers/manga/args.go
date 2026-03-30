package manga

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/monteiroliveira/mand/internal"
)

type MangaParserArgs struct {
	Source        *url.URL
	Verbose       bool
	Operation     Operation
	ErrorChan     chan error
	ListBatchSize int
}

func (m *MangaParserArgs) Validate() error {
	if m.Operation == DownloadListOperation {
		if m.ListBatchSize <= 0 {
			err := fmt.Errorf("Invalid batch size for download list content")
			return errors.Join(internal.SetSyntaxError(), err)
		}
	}
	return nil
}
