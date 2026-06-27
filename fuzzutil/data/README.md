# 国家级功能区映射数据

- `functional_zones_national.json`：397 个国家级主体（173 高新区 + 224 经开区）及地址别名
- 由 `scripts/bootstrap_functional_zones.py` 从公开名单生成，可手工修正 `district` / `aliases`
- 生成 Go 代码：`go generate ./fuzzutil/` 或 `go run gen_functional_zones.go ./fuzzutil`

更新流程：

1. 修正 `scripts/bootstrap_functional_zones.py` 中的 `DISTRICT` / `EXTRA` 覆盖表
2. 运行 bootstrap（需 Wikipedia / 商务部名单源文本，或直改 JSON）
3. `go run gen_functional_zones.go .`（在 fuzzutil 目录下）
4. `go test ./fuzzutil/`
