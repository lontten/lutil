package imgutil

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

// ImageReplacer 图片替换器
type ImageReplacer struct {
	MaxSize         int64                                   // 最大图片大小（字节）
	Timeout         time.Duration                           // 下载超时
	Concurrent      int                                     // 并发数
	MaxRetries      int                                     // 最大重试次数
	AllowedTypes    []string                                // 允许的图片类型
	uploadFileFun   func(localPath string) (string, error)  // 上传文件函数
	downloadFileFun func(remotePath string) (string, error) // 下载文件函数

	tempDir     string // 临时目录路径
	client      *http.Client
	clientOnce  sync.Once
	tempDirOnce sync.Once
}

// Option 配置选项函数类型
type Option func(*ImageReplacer)

// NewImageReplacer 创建图片替换器
func NewImageReplacer(uploadFileFun, downloadFileFun func(path string) (string, error), opts ...Option) *ImageReplacer {
	replacer := &ImageReplacer{
		MaxSize:         10 * 1024 * 1024, // 10MB
		Timeout:         30 * time.Second,
		Concurrent:      5,
		MaxRetries:      3,
		uploadFileFun:   uploadFileFun,
		downloadFileFun: downloadFileFun,
		AllowedTypes:    []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml"},
	}

	// 应用选项
	for _, opt := range opts {
		opt(replacer)
	}

	return replacer
}

// WithMaxSize 设置最大文件大小
func WithMaxSize(size int64) Option {
	return func(r *ImageReplacer) {
		r.MaxSize = size
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(r *ImageReplacer) {
		r.Timeout = timeout
	}
}

// WithConcurrent 设置并发数
func WithConcurrent(concurrent int) Option {
	return func(r *ImageReplacer) {
		r.Concurrent = concurrent
	}
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(retries int) Option {
	return func(r *ImageReplacer) {
		r.MaxRetries = retries
	}
}

// WithAllowedTypes 设置允许的图片类型
func WithAllowedTypes(types []string) Option {
	return func(r *ImageReplacer) {
		r.AllowedTypes = types
	}
}

// replaceTask 替换任务
type replaceTask struct {
	index int
	src   string
	node  *goquery.Selection
}

// downloadResult 下载结果
type downloadResult struct {
	task    replaceTask
	newPath string
	err     error
}

// getTempDir 获取临时目录（懒加载，线程安全）
func (r *ImageReplacer) getTempDir() (string, error) {
	var err error
	r.tempDirOnce.Do(func() {
		// 使用os.MkdirTemp创建临时目录，由操作系统管理
		uniqueDir, createErr := os.MkdirTemp("", "imgutil-*")
		if createErr != nil {
			err = fmt.Errorf("创建临时目录失败: %w", createErr)
			return
		}
		r.tempDir = uniqueDir
		log.Printf("临时目录已创建: %s", r.tempDir)
	})
	return r.tempDir, err
}

// getHTTPClient 获取HTTP客户端（懒加载，线程安全）
func (r *ImageReplacer) getHTTPClient() *http.Client {
	r.clientOnce.Do(func() {
		r.client = &http.Client{
			Timeout: r.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  true,
			},
		}
	})
	return r.client
}

// ReplaceRichText 替换富文本中的远程图片
func (r *ImageReplacer) ReplaceRichText(html string) (string, error) {
	return r.ReplaceRichTextWithContext(context.Background(), html)
}

// ReplaceRichTextWithContext 替换富文本中的远程图片（支持上下文）
func (r *ImageReplacer) ReplaceRichTextWithContext(ctx context.Context, html string) (string, error) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(ctx, r.Timeout*time.Duration(r.MaxRetries+1))
	defer cancel()

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %w", err)
	}

	// 提取所有远程图片
	imgNodes := doc.Find("img")
	if imgNodes.Length() == 0 {
		return html, nil
	}

	// 准备替换任务
	var tasks []replaceTask
	imgNodes.Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && r.isRemoteURL(src) {
			tasks = append(tasks, replaceTask{
				index: i,
				src:   src,
				node:  s,
			})
		}
	})

	if len(tasks) == 0 {
		return html, nil
	}

	// 并发下载图片
	log.Printf("开始处理 %d 张远程图片", len(tasks))
	start := time.Now()

	results := r.downloadImagesConcurrentlyWithContext(ctx, tasks)

	// 替换图片URL
	successCount := 0
	for _, result := range results {
		if result.err == nil && result.newPath != "" {
			result.task.node.SetAttr("src", result.newPath)
			successCount++
		} else if result.err != nil {
			log.Printf("图片处理失败 [%s]: %v", result.task.src, result.err)
		}
	}

	log.Printf("图片处理完成: 成功 %d/%d, 耗时 %v",
		successCount, len(tasks), time.Since(start))

	// 返回处理后的HTML
	htmlStr, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("生成HTML失败: %w", err)
	}

	// 提取body内容
	bodyContent := r.extractBodyContent(htmlStr)
	return bodyContent, nil
}

// extractBodyContent 提取body内容
func (r *ImageReplacer) extractBodyContent(htmlStr string) string {
	// 如果没有完整的HTML结构，直接返回
	if !strings.Contains(htmlStr, "<body>") {
		return htmlStr
	}

	bodyStart := strings.Index(htmlStr, "<body>")
	if bodyStart == -1 {
		return htmlStr
	}
	bodyStart += len("<body>")

	bodyEnd := strings.LastIndex(htmlStr, "</body>")
	if bodyEnd == -1 || bodyEnd <= bodyStart {
		return htmlStr
	}

	return strings.TrimSpace(htmlStr[bodyStart:bodyEnd])
}

// isRemoteURL 判断是否为远程URL
func (r *ImageReplacer) isRemoteURL(src string) bool {
	u, err := url.Parse(src)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// downloadImagesConcurrentlyWithContext 并发下载图片（支持上下文）
func (r *ImageReplacer) downloadImagesConcurrentlyWithContext(ctx context.Context, tasks []replaceTask) []downloadResult {
	results := make([]downloadResult, len(tasks))
	resultCh := make(chan downloadResult, len(tasks))

	// 使用errgroup管理并发
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(r.Concurrent)

	// 启动下载任务
	for _, task := range tasks {
		t := task // 创建副本
		g.Go(func() error {
			select {
			case <-ctx.Done():
				resultCh <- downloadResult{
					task: t,
					err:  ctx.Err(),
				}
				return nil
			default:
			}

			newPath, err := r.downloadAndSaveImageWithRetry(t.src)
			resultCh <- downloadResult{
				task:    t,
				newPath: newPath,
				err:     err,
			}
			return nil
		})
	}

	// 等待所有任务完成
	go func() {
		_ = g.Wait()
		close(resultCh)
	}()

	// 收集结果
	for result := range resultCh {
		results[result.task.index] = result
	}

	return results
}

// downloadAndSaveImageWithRetry 下载并保存图片（带重试）
func (r *ImageReplacer) downloadAndSaveImageWithRetry(imageURL string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= r.MaxRetries; attempt++ {
		if attempt > 0 {
			// 指数退避
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			time.Sleep(backoff)
			log.Printf("重试下载图片 [%s] 第 %d 次", imageURL, attempt)
		}

		newPath, err := r.downloadAndSaveImage(imageURL)
		if err == nil {
			return newPath, nil
		}

		lastErr = err

		// 如果是客户端错误（4xx），不再重试
		if strings.Contains(err.Error(), "HTTP错误: 4") {
			break
		}
	}

	return "", fmt.Errorf("下载失败（重试%d次）: %w", r.MaxRetries, lastErr)
}

// downloadAndSaveImage 下载并保存单张图片
func (r *ImageReplacer) downloadAndSaveImage(imageURL string) (string, error) {
	localPath, err := r.downloadFileFun(imageURL)
	if err != nil {
		return "", fmt.Errorf("下载文件失败: %w", err)
	}
	// 调用上传函数
	newURL, err := r.uploadFileFun(localPath)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}
	// 无论上传成功还是失败，都删除本地临时文件
	go func() {
		if err := os.Remove(localPath); err != nil {
			log.Printf("删除临时文件失败: %v (文件: %s)", err, localPath)
		} else {
			log.Printf("临时文件已删除: %s", localPath)
		}
	}()
	return newURL, nil
}

// imageMagicNumbers 常见图片格式的魔数
var imageMagicNumbers = []struct {
	magic []byte
	name  string
}{
	{[]byte{0xFF, 0xD8, 0xFF}, "jpeg"},
	{[]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, "png"},
	{[]byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, "gif87a"},
	{[]byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}, "gif89a"},
	{[]byte{0x52, 0x49, 0x46, 0x46}, "webp"},
	{[]byte{0x42, 0x4D}, "bmp"},
	{[]byte{0x49, 0x49, 0x2A, 0x00}, "tiff (little endian)"},
	{[]byte{0x4D, 0x4D, 0x00, 0x2A}, "tiff (big endian)"},
	{[]byte{0x3C, 0x3F, 0x78, 0x6D, 0x6C}, "svg"},
}

// isValidImageHeader 检查是否为有效的图片头部
func (r *ImageReplacer) isValidImageHeader(header []byte) bool {
	if len(header) < 4 {
		return false
	}

	// 检查常见图片格式的魔数
	for _, imgType := range imageMagicNumbers {
		if len(header) >= len(imgType.magic) && bytes.Equal(header[:len(imgType.magic)], imgType.magic) {
			return true
		}
	}

	return false
}

// isAllowedType 检查是否为允许的图片类型
func (r *ImageReplacer) isAllowedType(contentType string) bool {
	// 解析MIME类型
	mimeType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	for _, allowedType := range r.AllowedTypes {
		if mimeType == allowedType {
			return true
		}
	}
	return false
}

// getExtension 获取文件扩展名
func (r *ImageReplacer) getExtension(contentType, urlStr string) string {
	// 从Content-Type获取
	mimeType, _, err := mime.ParseMediaType(contentType)
	if err == nil {
		exts, err := mime.ExtensionsByType(mimeType)
		if err == nil && len(exts) > 0 {
			// 优先使用常见的扩展名
			for _, ext := range exts {
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" ||
					ext == ".gif" || ext == ".webp" || ext == ".svg" {
					return ext
				}
			}
			return exts[0]
		}
	}

	// 从URL路径获取
	u, err := url.Parse(urlStr)
	if err == nil {
		ext := filepath.Ext(u.Path)
		if ext != "" {
			// 规范化扩展名
			ext = strings.ToLower(ext)
			if len(ext) > 1 && ext[0] == '.' {
				return ext
			}
		}
	}

	// 默认使用.jpg
	return ".jpg"
}

// GetTempDir 获取临时目录路径（用于调试）
func (r *ImageReplacer) GetTempDir() string {
	return r.tempDir
}
