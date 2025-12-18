package imgutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	"github.com/gofrs/uuid"
	"golang.org/x/sync/errgroup"
)

// ImageReplacer 图片替换器
type ImageReplacer struct {
	MaxSize       int64                                  // 最大图片大小（字节）
	Timeout       time.Duration                          // 下载超时
	Concurrent    int                                    // 并发数
	MaxRetries    int                                    // 最大重试次数
	AllowedTypes  []string                               // 允许的图片类型
	uploadFileFun func(localPath string) (string, error) // 上传文件函数

	tempDir     string // 临时目录路径
	client      *http.Client
	clientOnce  sync.Once
	tempDirOnce sync.Once
}

// Option 配置选项函数类型
type Option func(*ImageReplacer)

// NewImageReplacer 创建图片替换器
func NewImageReplacer(uploadFileFun func(localPath string) (string, error), opts ...Option) *ImageReplacer {
	replacer := &ImageReplacer{
		MaxSize:       10 * 1024 * 1024, // 10MB
		Timeout:       30 * time.Second,
		Concurrent:    5,
		MaxRetries:    3,
		uploadFileFun: uploadFileFun,
		AllowedTypes:  []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml"},
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
		// 在系统临时目录下创建唯一的子目录
		tempBase := os.TempDir()
		uniqueDir, createErr := os.MkdirTemp(tempBase, "imgutil-")
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
			log.Printf("图片下载失败 [%s]: %v", result.task.src, result.err)
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
	// 获取临时目录
	tempDir, err := r.getTempDir()
	if err != nil {
		return "", err
	}

	// 创建HTTP客户端
	client := r.getHTTPClient()

	// 发送请求
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 检查内容类型
	contentType := resp.Header.Get("Content-Type")
	if !r.isAllowedType(contentType) {
		return "", fmt.Errorf("不支持的图片类型: %s", contentType)
	}

	// 检查内容大小
	if resp.ContentLength > r.MaxSize {
		return "", fmt.Errorf("图片太大: %d bytes (最大允许: %d bytes)",
			resp.ContentLength, r.MaxSize)
	}

	// 生成文件名
	ext := r.getExtension(contentType, imageURL)
	v7, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("生成UUID失败: %w", err)
	}
	filename := fmt.Sprintf("%s%s", v7.String(), ext)

	// 创建日期子目录
	datePath := time.Now().Format("2006/01/02")
	saveDir := filepath.Join(tempDir, datePath)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 保存文件
	filePath := filepath.Join(saveDir, filename)
	if err := r.saveImageToFile(resp.Body, filePath, r.MaxSize); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 验证图片
	if err := r.validateImageFile(filePath); err != nil {
		os.Remove(filePath) // 删除无效文件
		return "", fmt.Errorf("图片验证失败: %w", err)
	}

	// 返回相对路径并上传
	relativePath := filepath.Join(datePath, filename)
	return r.uploadFileFun(relativePath)
}

// saveImageToFile 将图片保存到文件
func (r *ImageReplacer) saveImageToFile(src io.Reader, filePath string, maxSize int64) error {
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用LimitReader限制大小
	limitedReader := io.LimitReader(src, maxSize)

	// 同时写入文件和内存缓冲区（用于验证）
	var buf bytes.Buffer
	teeReader := io.TeeReader(limitedReader, &buf)

	// 先读取一部分用于验证
	header := make([]byte, 512)
	n, err := teeReader.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("读取图片头部失败: %w", err)
	}

	// 验证图片头部
	if !r.isValidImageHeader(header[:n]) {
		return fmt.Errorf("无效的图片格式")
	}

	// 继续写入文件
	if _, err := file.Write(header[:n]); err != nil {
		return err
	}

	// 复制剩余内容
	_, err = io.Copy(file, teeReader)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// validateImageFile 验证图片文件
func (r *ImageReplacer) validateImageFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取文件头部
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return err
	}

	// 检查文件大小
	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		return fmt.Errorf("文件为空")
	}

	if info.Size() == r.MaxSize {
		return fmt.Errorf("文件可能被截断")
	}

	// 验证图片格式
	if !r.isValidImageHeader(header[:n]) {
		return fmt.Errorf("无效的图片格式")
	}

	return nil
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

// Cleanup 清理临时目录
func (r *ImageReplacer) Cleanup() error {
	if r.tempDir == "" {
		return nil
	}

	log.Printf("清理临时目录: %s", r.tempDir)
	if err := os.RemoveAll(r.tempDir); err != nil {
		return fmt.Errorf("清理临时目录失败: %w", err)
	}
	r.tempDir = ""
	r.tempDirOnce = sync.Once{}
	return nil
}

// GetTempDir 获取临时目录路径（用于调试）
func (r *ImageReplacer) GetTempDir() string {
	return r.tempDir
}

// Finalize 终结器，确保资源被清理
func (r *ImageReplacer) Finalize() {
	_ = r.Cleanup()
}
