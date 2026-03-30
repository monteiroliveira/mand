package manga

import (
	"sync"
	"testing"
)

func TestExtractList_BasicBatching(t *testing.T) {
	links := []string{"a", "b", "c", "d", "e"}
	extracted := make([]string, 0)
	mu := sync.Mutex{}

	wg := new(sync.WaitGroup)
	ch := make(chan error, 10)
	extractor := func(link string, wg *sync.WaitGroup, ch chan error) {
		defer wg.Done()
		mu.Lock()
		extracted = append(extracted, link)
		mu.Unlock()
	}

	err := ExtractList(links, 2, extractor, wg, ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(extracted) != 5 {
		t.Errorf("expected 5 extracted links, got %d", len(extracted))
	}
}

func TestExtractList_BatchSizeOne(t *testing.T) {
	links := []string{"a", "b", "c"}
	count := 0
	mu := sync.Mutex{}

	wg := new(sync.WaitGroup)
	ch := make(chan error, 10)
	extractor := func(link string, wg *sync.WaitGroup, ch chan error) {
		defer wg.Done()
		mu.Lock()
		count++
		mu.Unlock()
	}

	err := ExtractList(links, 1, extractor, wg, ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3, got %d", count)
	}
}

func TestExtractList_EmptyLinks(t *testing.T) {
	wg := new(sync.WaitGroup)
	ch := make(chan error, 10)
	extractor := func(link string, wg *sync.WaitGroup, ch chan error) {
		defer wg.Done()
	}

	err := ExtractList([]string{}, 5, extractor, wg, ch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExtractList_InvalidBatchSize(t *testing.T) {
	wg := new(sync.WaitGroup)
	ch := make(chan error, 10)
	extractor := func(link string, wg *sync.WaitGroup, ch chan error) {
		defer wg.Done()
	}

	err := ExtractList([]string{"a"}, 0, extractor, wg, ch)
	if err == nil {
		t.Fatal("expected error for zero batch size")
	}

	err = ExtractList([]string{"a"}, -1, extractor, wg, ch)
	if err == nil {
		t.Fatal("expected error for negative batch size")
	}
}

func TestExtractList_ErrorReporting(t *testing.T) {
	links := []string{"a", "b"}
	wg := new(sync.WaitGroup)
	ch := make(chan error, 10)
	extractor := func(link string, wg *sync.WaitGroup, ch chan error) {
		defer wg.Done()
		ch <- &testError{msg: "extraction failed for " + link}
	}

	err := ExtractList(links, 5, extractor, wg, ch)
	if err != nil {
		t.Fatalf("unexpected error from ExtractList: %v", err)
	}

	if len(ch) != 2 {
		t.Errorf("expected 2 errors in channel, got %d", len(ch))
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
