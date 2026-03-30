## Mand - Manga And Novel Downloader

CLI application to Download Mangas & Novel from mapped sources.

```bash
go build cmd/main.go -o mand
```

```bash
mand manga d {{ MANGA_URL }}
```

> [!NOTE]
> Supported Parsers:
> - [x] Manga Dex
> - [x] Manga Read

> [!NOTE]
> TODOs:
> - [ ] Create a log system;
> - [ ] Add a way control to what chapter will be downloaded from download list (dl) operation;
> - [ ] Add async workers control from download list (dl) operation;
> - [ ] Add more parsers;
> - [ ] Init novel parser package.
