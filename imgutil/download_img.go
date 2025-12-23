package imgutil

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// -------------------------- 私有配置 --------------------------
type config struct {
	maxImageSize        int64
	downloadTimeout     time.Duration
	tempDirPrefix       string
	allowedContentTypes map[string]bool
	enableCompress      bool
	compressQuality     int
	redirectMaxCount    int
	requestHeaders      http.Header
	refererMap          map[string]string // 新增：域名到 Referer 的映射
}

func defaultConfig() *config {
	return &config{
		maxImageSize:        5 * 1024 * 1024,
		downloadTimeout:     10 * time.Second,
		tempDirPrefix:       "img_",
		allowedContentTypes: map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true, "image/gif": true},
		enableCompress:      false,
		compressQuality:     80,
		redirectMaxCount:    3,
		requestHeaders: http.Header{
			"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
			"Accept":     {"image/webp,image/png,image/jpeg,image/gif,*/*;q=0.8"},
		},
		refererMap: make(map[string]string), // 初始化 refererMap
	}
}

// -------------------------- Builder 结构体 --------------------------
type Builder struct {
	cfg *config
	err error
}

func NewBuilder() *Builder {
	return &Builder{cfg: defaultConfig()}
}

// 检查错误辅助方法
func (b *Builder) checkErr() bool {
	return b.err != nil
}

func (b *Builder) setErr(err error) {
	if b.err == nil {
		b.err = err
	}
}

// 链式配置方法
func (b *Builder) WithMaxImageSize(size int64) *Builder {
	if b.checkErr() {
		return b
	}
	if size <= 0 {
		b.setErr(fmt.Errorf("最大图片大小必须大于0"))
		return b
	}
	b.cfg.maxImageSize = size
	return b
}

func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	if b.checkErr() {
		return b
	}
	if timeout <= 0 {
		b.setErr(fmt.Errorf("超时时间必须大于0"))
		return b
	}
	b.cfg.downloadTimeout = timeout
	return b
}

func (b *Builder) WithAllowedImageFormats(formats ...string) *Builder {
	if b.checkErr() {
		return b
	}

	if len(formats) == 0 {
		b.setErr(fmt.Errorf("至少需要指定一种图片格式"))
		return b
	}

	allowed := make(map[string]bool)
	for _, format := range formats {
		switch strings.ToLower(format) {
		case "jpeg", "jpg":
			allowed["image/jpeg"] = true
		case "png":
			allowed["image/png"] = true
		case "webp":
			allowed["image/webp"] = true
		case "gif":
			allowed["image/gif"] = true
		default:
			b.setErr(fmt.Errorf("不支持的图片格式: %s", format))
			return b
		}
	}
	b.cfg.allowedContentTypes = allowed
	return b
}

func (b *Builder) WithCompress(enable bool) *Builder {
	if b.checkErr() {
		return b
	}
	b.cfg.enableCompress = enable
	return b
}

func (b *Builder) WithCompressQuality(quality int) *Builder {
	if b.checkErr() {
		return b
	}
	if quality < 1 || quality > 100 {
		b.setErr(fmt.Errorf("压缩质量必须在1-100之间"))
		return b
	}
	b.cfg.compressQuality = quality
	return b
}

func (b *Builder) WithReferer(referer string) *Builder {
	if b.checkErr() {
		return b
	}
	b.cfg.requestHeaders.Set("Referer", referer)
	return b
}

// 新增：设置 RefererMap
func (b *Builder) WithRefererMap(refererMap map[string]string) *Builder {
	if b.checkErr() {
		return b
	}
	if refererMap == nil {
		b.setErr(fmt.Errorf("RefererMap 不能为 nil"))
		return b
	}
	b.cfg.refererMap = refererMap
	return b
}

// 新增：添加或更新单个域名的 Referer
func (b *Builder) AddReferer(domain, referer string) *Builder {
	if b.checkErr() {
		return b
	}
	if domain == "" {
		b.setErr(fmt.Errorf("域名不能为空"))
		return b
	}
	if b.cfg.refererMap == nil {
		b.cfg.refererMap = make(map[string]string)
	}
	b.cfg.refererMap[domain] = referer
	return b
}

func (b *Builder) WithUserAgent(userAgent string) *Builder {
	if b.checkErr() {
		return b
	}
	b.cfg.requestHeaders.Set("User-Agent", userAgent)
	return b
}

// Build 构建下载器实例
func (b *Builder) Build() (*Downloader, error) {
	if b.err != nil {
		return nil, b.err
	}
	if len(b.cfg.allowedContentTypes) == 0 {
		return nil, fmt.Errorf("未配置允许的图片格式")
	}

	// 创建可复用的 HTTP 客户端
	client := &http.Client{
		Timeout: b.cfg.downloadTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > b.cfg.redirectMaxCount {
				return fmt.Errorf("超过最大重定向次数: %d", b.cfg.redirectMaxCount)
			}
			return nil
		},
	}

	// 返回下载器实例
	return &Downloader{
		cfg:    b.cfg,
		client: client,
	}, nil
}

// -------------------------- 下载器结构体 --------------------------
type Downloader struct {
	cfg       *config
	client    *http.Client
	tempPaths []string
}

// DownloadImg 下载单张图片
func (d *Downloader) DownloadImg(imgURL string) (string, error) {
	// 1. 校验 URL
	if imgURL == "" {
		return "", fmt.Errorf("图片URL不能为空")
	}
	parsedURL, err := url.ParseRequestURI(imgURL)
	if err != nil {
		return "", fmt.Errorf("URL格式非法: %w", err)
	}

	// 2. 下载图片数据
	imgData, ct, err := d.download(imgURL, parsedURL)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}

	// 3. 校验格式
	if !d.isContentTypeAllowed(ct) {
		return "", fmt.Errorf("不支持的图片格式: %s", ct)
	}

	// 4. 压缩（可选）
	if d.cfg.enableCompress {
		imgData, err = d.compress(imgData, ct)
		if err != nil {
			return "", fmt.Errorf("压缩失败: %w", err)
		}
	}

	// 5. 保存到临时文件
	tempPath, err := d.saveToTemp(imgData, ct)
	if err != nil {
		return "", fmt.Errorf("保存失败: %w", err)
	}

	// 记录临时路径
	d.tempPaths = append(d.tempPaths, tempPath)
	return tempPath, nil
}

// BatchDownload 批量下载
func (d *Downloader) BatchDownload(imgURLs []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, url := range imgURLs {
		path, err := d.DownloadImg(url)
		if err != nil {
			return result, fmt.Errorf("下载%s失败: %w", url, err)
		}
		result[url] = path
	}
	return result, nil
}

// CleanAllTemp 清理所有下载的临时文件
func (d *Downloader) CleanAllTemp() error {
	var errs []error
	for _, path := range d.tempPaths {
		if err := os.Remove(path); err != nil {
			errs = append(errs, fmt.Errorf("清理%s失败: %w", path, err))
		}
	}
	d.tempPaths = nil
	if len(errs) > 0 {
		return fmt.Errorf("清理临时文件失败: %v", errs)
	}
	return nil
}

// -------------------------- 下载器私有方法 --------------------------
// isContentTypeAllowed 检查内容类型是否允许
func (d *Downloader) isContentTypeAllowed(ct string) bool {
	return d.cfg.allowedContentTypes[ct]
}

// getRefererForURL 根据 URL 获取对应的 Referer
func (d *Downloader) getRefererForURL(parsedURL *url.URL) string {
	// 先检查是否有全局 Referer
	globalReferer := d.cfg.requestHeaders.Get("Referer")

	// 然后检查是否有针对该域名的 Referer
	if d.cfg.refererMap != nil {
		host := parsedURL.Host

		// 尝试匹配完整域名
		if referer, ok := d.cfg.refererMap[host]; ok {
			return referer
		}

		// 尝试匹配主域名（去掉 www. 前缀）
		hostWithoutWWW := strings.TrimPrefix(host, "www.")
		if referer, ok := d.cfg.refererMap[hostWithoutWWW]; ok {
			return referer
		}

		// 尝试匹配包含通配符的域名（如 *.example.com）
		domainParts := strings.Split(host, ".")
		if len(domainParts) >= 2 {
			wildcardDomain := "*." + strings.Join(domainParts[len(domainParts)-2:], ".")
			if referer, ok := d.cfg.refererMap[wildcardDomain]; ok {
				return referer
			}
		}
	}

	return globalReferer
}

// download 执行下载逻辑
func (d *Downloader) download(imgURL string, parsedURL *url.URL) ([]byte, string, error) {
	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		return nil, "", err
	}

	// 设置请求头
	for k, v := range d.cfg.requestHeaders {
		req.Header[k] = v
	}

	// 根据 URL 设置对应的 Referer
	if referer := d.getRefererForURL(parsedURL); referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("响应状态码: %d", resp.StatusCode)
	}

	// 解析 Content-Type
	ct := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(ct)
	if err == nil {
		ct = mediaType
	} else if idx := strings.Index(ct, ";"); idx != -1 {
		// 如果无法解析，则简单处理
		ct = strings.TrimSpace(ct[:idx])
	}

	// 限制大小
	buf := bytes.NewBuffer(nil)
	n, err := io.CopyN(buf, resp.Body, d.cfg.maxImageSize+1)
	if err != nil && err != io.EOF {
		return nil, "", err
	}
	if n > d.cfg.maxImageSize {
		return nil, "", fmt.Errorf("图片超过大小限制: %d字节", d.cfg.maxImageSize)
	}

	return buf.Bytes(), ct, nil
}

// compress 压缩图片
func (d *Downloader) compress(imgData []byte, ct string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	switch ct {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: d.cfg.compressQuality})
	case "image/png":
		err = png.Encode(&buf, img)
	//case "image/webp":
	//	err = webp.Encode(&buf, img, &webp.Options{Quality: float32(d.cfg.compressQuality)})
	case "image/gif":
		return imgData, nil // GIF 暂不压缩
	default:
		return nil, fmt.Errorf("不支持压缩格式: %s", ct)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// saveToTemp 保存临时文件
func (d *Downloader) saveToTemp(imgData []byte, ct string) (string, error) {
	// 获取文件扩展名
	ext := getFileExt(ct)

	// 使用 os.CreateTemp 自动生成临时文件
	tempFile, err := os.CreateTemp("", d.cfg.tempDirPrefix+"*"+ext)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 写入数据
	if _, err := tempFile.Write(imgData); err != nil {
		// 如果写入失败，删除已创建的临时文件
		os.Remove(tempFile.Name())
		return "", err
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return absPath, nil
}

// -------------------------- 辅助函数 --------------------------
func getFileExt(contentType string) string {
	// 使用 mime 包获取标准扩展名
	exts, _ := mime.ExtensionsByType(contentType)
	if len(exts) > 0 {
		return exts[0]
	}

	// 后备方案
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".jpg"
	}
}

// CleanTemp 单独清理某个临时文件
func CleanTemp(tempPath string) error {
	if tempPath == "" {
		return fmt.Errorf("临时路径不能为空")
	}
	return os.Remove(tempPath)
}

// GetTempPaths 获取当前所有临时文件路径（用于调试或手动管理）
func (d *Downloader) GetTempPaths() []string {
	return d.tempPaths
}

// GetRefererMap 获取当前的 RefererMap（用于调试或查看配置）
func (d *Downloader) GetRefererMap() map[string]string {
	return d.cfg.refererMap
}
