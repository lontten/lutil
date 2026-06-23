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
| `codeutil` | Encoding, hashing, random strings, password helpers |
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

## Development

```bash
go mod verify
go test -race -count=1 ./...
```

## License

Apache-2.0. See [LICENSE](LICENSE).
