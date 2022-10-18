package status

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Registory 注册中心
type Registory struct {
	Address     string // 应用地址, 例如 http://127.0.0.1:8080
	ContextPath string // 应用上下文
}

// RegistoryInfo 注册中心响应结果
type RegistoryInfo struct {
	Status string `json:"status,omitempty"` // 注册状态
}

// Info 注册中心应用状态
func (r *Registory) Info() string {
	return getInfo(strings.Join([]string{r.Address, r.ContextPath, "/actuator/health"}, ""))
}

// {"status":"UP"}
func getInfo(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "UNKOWN"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "UNKOWN"
	}

	if !successfulStatusCode(resp.StatusCode) {
		return "UNKOWN"
	}

	info := &RegistoryInfo{}
	if len(body) > 0 {
		_ = json.Unmarshal(body, &info)
	}
	return info.Status
}
