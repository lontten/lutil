# lutil

[中文](README.zh.md)

[![test](https://github.com/lontten/lutil/actions/workflows/test.yml/badge.svg)](https://github.com/lontten/lutil/actions/workflows/test.yml)

A collection of Go utility packages from [lontten](https://github.com/lontten), built on [lcore/v2](https://github.com/lontten/lcore) types such as `LocalDate`, `LocalDateTime`, and `Decimal`.

## Requirements

- Go 1.25+

## Installation

```bash
go get -u github.com/lontten/lutil
```

## Packages

| Package | Description |
|---------|-------------|
| `lutil` | Goroutine pool, key-based mutex (`KeyLock`) |
| `codeutil` | Encoding, hashing, random strings; `HashPassword`/`VerifyPassword` (bcrypt; prefer over deprecated `EnPwd`) |
| `dateutil` | `LocalDate` comparison and aggregation |
| `datetimeutil` | `LocalDateTime` comparison and aggregation |
| `decimalutil` | `decimal.Decimal` arithmetic helpers |
| `fileutil` | Temp files, copy, path helpers |
| `fuzzutil` | Fuzzy matching (Like) and vocabulary extraction |
| `imgutil` | Image download, Base64, HTML/richtext image handling |
| `jsonutil` | JSON marshal/unmarshal helpers |
| `listutil` | Slice set operations and `ListTool` |
| `logutil` | Simple logging helpers |
| `moneyutil` | Money/decimal operations and discount helpers |
| `netutil` | HTTP, IP resolution, file download |
| `numutil` | Numeric utilities |
| `perfutil` | Simple performance timing |
| `structutil` | Struct ↔ map conversion |
| `strutil` | String checks and substring helpers |

## Migrating to v0.5.0

This release includes breaking changes:

- `netutil.CleanString` was removed. Map spaces/controls yourself, or use `SafeFileName` for filenames.
- Prefer `netutil.SafeFileName` over deprecated `SafeURL` (output may differ: `filepath.Base`, extension-preserving truncate).
- `DownloadFileToLocal` accepts HTTP 200 only, caps at 500MB, and uses a 30s timeout; use `DownloadFileToLocalLimit` for other size limits.
- `netutil` HTTP helpers (`Get`, `PostJson*`, `PostForm*`) use a 30s default timeout.
- Pool `SubmitErr` returns sentinel errors: use `errors.Is(err, lutil.ErrQueueFull)` / `ErrPoolClosed` instead of matching error strings.

## Development

```bash
go mod verify
go test -race -count=1 ./...
```

## License

Apache-2.0. See [LICENSE](LICENSE).
