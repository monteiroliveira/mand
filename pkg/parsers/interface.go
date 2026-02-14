package parsers

// TODO: Adjust interface content
type Parser interface {
	Extract() ([]byte, error)
	Download() error
}
