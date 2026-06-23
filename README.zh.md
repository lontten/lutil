# lutil

[English](README.md)

[![test](https://github.com/lontten/lutil/actions/workflows/test.yml/badge.svg)](https://github.com/lontten/lutil/actions/workflows/test.yml)

[lontten](https://github.com/lontten) 的 Go 工具库集合，基于 [lcore/v2](https://github.com/lontten/lcore) 的 `LocalDate`、`LocalDateTime`、`Decimal` 等类型。

## 环境要求

- Go 1.25+

## 安装

```bash
go get -u github.com/lontten/lutil
```

## 子包说明

| 包 | 说明 |
|----|------|
| `lutil` | 协程池与按键互斥锁等基础工具 |
| `codeutil` | 编码、哈希、随机字符串与密码工具 |
| `dateutil` | `LocalDate` 的比较与聚合工具 |
| `datetimeutil` | `LocalDateTime` 的比较与聚合工具 |
| `decimalutil` | `decimal.Decimal` 的运算工具 |
| `fileutil` | 临时文件、文件复制与路径解析工具 |
| `fuzzutil` | 字符串模糊匹配（Like）与关系链词表提取 |
| `imgutil` | 图片下载、Base64 与 HTML 富文本图片处理 |
| `jsonutil` | JSON 序列化与反序列化便捷函数 |
| `listutil` | 切片集合运算与 `ListTool` 条件检查工具 |
| `logutil` | 简单的日志输出工具 |
| `moneyutil` | 金额运算与折扣计算工具 |
| `netutil` | HTTP 请求、IP 解析与文件下载工具 |
| `numutil` | 数值相关的工具函数 |
| `perfutil` | 简单的性能计时工具 |
| `structutil` | 结构体与 map 之间的转换工具 |
| `strutil` | 字符串判断与截取工具 |

## 开发与测试

```bash
go mod verify
go test -race -count=1 ./...
```

## 许可证

Apache-2.0，详见 [LICENSE](LICENSE)。
