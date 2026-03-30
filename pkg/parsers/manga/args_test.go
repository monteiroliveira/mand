package manga

import (
	"testing"
)

func TestValidate_DownloadOperation(t *testing.T) {
	args := &MangaParserArgs{
		Operation: DownloadOperation,
	}
	if err := args.Validate(); err != nil {
		t.Fatalf("expected no error for download operation, got: %v", err)
	}
}

func TestValidate_DownloadListOperation_ValidBatch(t *testing.T) {
	args := &MangaParserArgs{
		Operation:     DownloadListOperation,
		ListBatchSize: 5,
	}
	if err := args.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_DownloadListOperation_ZeroBatch(t *testing.T) {
	args := &MangaParserArgs{
		Operation:     DownloadListOperation,
		ListBatchSize: 0,
	}
	if err := args.Validate(); err == nil {
		t.Fatal("expected error for zero batch size")
	}
}

func TestValidate_DownloadListOperation_NegativeBatch(t *testing.T) {
	args := &MangaParserArgs{
		Operation:     DownloadListOperation,
		ListBatchSize: -1,
	}
	if err := args.Validate(); err == nil {
		t.Fatal("expected error for negative batch size")
	}
}
