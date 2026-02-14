package internal

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/alexflint/go-arg"
)

type Source struct {
	URL *url.URL
}

func (s *Source) UnmarshalText(b []byte) error {
	value := string(b)
	source, err := url.Parse(value)
	if err != nil {
		return errors.Join(SetSyntaxError(), fmt.Errorf("Failed to parse source url"), err)
	}
	if source.Scheme == "" || source.Host == "" {
		parseSourceErr := fmt.Errorf("Failed to parse source url, URL is not valid")
		return errors.Join(SetSyntaxError(), parseSourceErr)
	}
	s.URL = source
	return nil
}

type ConsoleArgs struct {
	Source  Source `arg:"positional,required" help:"Link to the source of manga or novel"`
	Verbose bool   `arg:"-v,--verbose" help:"verbosity level"`
}

func NewArgs() *ConsoleArgs {
	args := &ConsoleArgs{}
	arg.MustParse(args)

	return args
}
