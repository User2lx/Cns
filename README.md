# gobang [![builds.sr.ht status](https://builds.sr.ht/~delthas/gobang.svg)](https://builds.sr.ht/~delthas/gobang?)

A small custom DuckDuckGo proxy for defining your own bangs.

If no bangs are used, will redirect to Google by default.

Setup:
- `gobang -port 8080 -url "http://example.com:8080"`

Usage:
- visit `/` to make your browser autodetect the search engine
- simply use it as your browser search engine
- add bangs by filling the form at `/add`
- manage bangs by stopping gobang and editing `bangs.txt`

Building:
- `go install git.saucisseroyale.cc/delthas/picopacker`
- make sure to `go generate` after editing any HTML files

| OS | URL |
|---|---|
| Linux x64 | https://delthas.fr/gobang/linux/gobang |
| Mac OS X x64 | https://delthas.fr/gobang/mac/gobang |
| Windows x64 | https://delthas.fr/gobang/windows/gobang.exe |
