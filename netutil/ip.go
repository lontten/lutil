package netutil

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type contextKey string

const ipContextKey contextKey = "client_ip"

// IPConfig 简化配置
type IPConfig struct {
	TrustedCIDRs []string     // 信任的代理CIDR
	IPHeaders    []string     // 检查的IP头
	trustedNets  []*net.IPNet // 缓存的网络
}

// 默认配置
var DefaultConfig = &IPConfig{
	TrustedCIDRs: []string{
		"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
		"127.0.0.0/8", "fc00::/7", "::1/128",
	},
	IPHeaders: []string{
		"X-Real-IP",
		"X-Forwarded-For",
		"CF-Connecting-IP",
		"True-Client-IP",
	},
}

// 初始化信任网络
func (c *IPConfig) init() error {
	if c.trustedNets != nil {
		return nil
	}

	c.trustedNets = make([]*net.IPNet, 0, len(c.TrustedCIDRs))
	for _, cidr := range c.TrustedCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		c.trustedNets = append(c.trustedNets, network)
	}
	return nil
}

// 检查IP是否在信任网络中
func (c *IPConfig) isTrusted(ipStr string) bool {
	if err := c.init(); err != nil {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, network := range c.trustedNets {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// RealIP 获取客户端真实IP
func RealIP(r *http.Request, config *IPConfig) string {
	if config == nil {
		config = DefaultConfig
	}

	// 获取远程IP
	remoteIP := extractRemoteIP(r.RemoteAddr)

	// 收集头部中的候选IP
	candidates := collectIPsFromHeaders(r, config.IPHeaders)

	// 选择最终IP
	return selectIP(candidates, remoteIP, config)
}

// 从RemoteAddr提取IP
func extractRemoteIP(remoteAddr string) string {
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	return remoteAddr
}

// 从HTTP头收集IP
func collectIPsFromHeaders(r *http.Request, headers []string) []string {
	var ips []string
	for _, header := range headers {
		if value := r.Header.Get(header); value != "" {
			for _, ip := range strings.Split(value, ",") {
				if ip = strings.TrimSpace(ip); ip != "" && net.ParseIP(ip) != nil {
					ips = append(ips, ip)
				}
			}
		}
	}
	return ips
}

// 选择最终IP
func selectIP(candidates []string, remoteIP string, config *IPConfig) string {
	if len(candidates) == 0 {
		return remoteIP
	}

	// 如果远程IP是信任的代理，从右向左找第一个非信任IP
	if config.isTrusted(remoteIP) {
		for i := len(candidates) - 1; i >= 0; i-- {
			if !config.isTrusted(candidates[i]) {
				return candidates[i]
			}
		}
		// 全是信任的代理，返回第一个（最原始客户端）
		return candidates[0]
	}

	// 远程IP不信任，直接返回
	return remoteIP
}

// 简化版本
func RealIPSimple(r *http.Request) string {
	return RealIP(r, DefaultConfig)
}

// 中间件
func IPMiddleware(config *IPConfig) func(http.Handler) http.Handler {
	if config == nil {
		config = DefaultConfig
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := RealIP(r, config)
			ctx := context.WithValue(r.Context(), ipContextKey, ip)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// 从上下文获取IP
func IPFromContext(ctx context.Context) string {
	if ip, ok := ctx.Value(ipContextKey).(string); ok {
		return ip
	}
	return ""
}

// 基础工具函数
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IPVersion(ip string) int {
	if parsed := net.ParseIP(ip); parsed != nil {
		if parsed.To4() != nil {
			return 4
		}
		return 6
	}
	return 0
}
