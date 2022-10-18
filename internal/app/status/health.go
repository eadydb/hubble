package status

import (
	"io"
	"net/http"
	"strings"
)

type Health struct {
	Address     string // 应用地址, 例如 http://127.0.0.1:8080
	ContextPath string // 应用上下文
	Url         string // 监控检查url, 默认 /actutor/info
}

// NewDefaultHealth 默认监控检查初始化
func NewDefaultHealth(addr, ctxPath string) *Health {
	return &Health{
		Address: addr,
		Url:     "/actutor/info",
	}
}

// NewHealth 监控检查初始化
func NewHealth(addr, ctxPath, url string) *Health {
	return &Health{
		Address:     addr,
		ContextPath: ctxPath,
		Url:         url,
	}
}

// Healthz 健康检查
func (h *Health) Healthz() bool {
	return healthProbe(strings.Join([]string{h.Address, h.ContextPath, h.Url}, ""))
}

// healthProbe 健康检查探测
func healthProbe(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	if !successfulStatusCode(resp.StatusCode) {
		return false
	}
	return true
}

// successfulStatusCode 正确响应结果状态码
func successfulStatusCode(code int) bool {
	return code >= 200 && code < 300
}
