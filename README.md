## Mand - Manga And Novel Downloader

CLI application to download Mangas & Novels from mapped sources.

Mand automatically detects the source from the provided URL and uses the appropriate parser to scrape or fetch chapter data. Each parser handles page extraction, naming, and output format specific to its source. Downloaded chapters are saved locally as images or PDFs depending on the source.

### How it works

1. You provide a manga URL from a supported source
2. Mand identifies the source by the URL hostname and selects the right parser
3. For single chapter downloads (`d`), it extracts all pages and saves them as a single file
4. For list downloads (`dl`), it extracts all chapter links from the manga page and downloads them concurrently using async workers

### Build

```bash
make build
```

### Install / Uninstall

```bash
make install            # installs to /usr/local/bin
make install PREFIX=~/.local  # custom prefix
make uninstall
```

### Usage

```bash
# Download a single chapter
mand manga d <SOURCE_URL>

# Download chapters list from source
mand manga dl <SOURCE_URL> [-b <batch_size>]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-b, --batch` | `5` | Batch size for async workers in list download |
| `-v, --verbose` | `0` | Verbosity level (1 info, 2 debug, 3 trace) |

> [!NOTE]
> Supported Parsers:
> - [x] Manga Dex
> - [x] Manga Read

### Development

```bash
make test    # run tests
make fmt     # format code
make lint    # run linter
make clean   # remove build artifacts
```

> [!NOTE]
> TODOs:
> - [x] Create a log system;
> - [ ] Add a way control to what chapter will be downloaded from download list (dl) operation;
> - [x] Add async workers control from download list (dl) operation;
> - [ ] Add more parsers;
> - [ ] Init novel parser package.
