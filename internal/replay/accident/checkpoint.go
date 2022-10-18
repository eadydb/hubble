package accident

import "time"

// CheckPoint
type CheckPoint struct {
	AppName      string    `json:"app_name"`      // 应用名称
	InstanceName string    `json:"instance_name"` // 实例名称
	StartTime    time.Time `json:"start_time"`    // 开始时间
	EndTime      time.Time `json:"end_time"`      // 结束时间
	Kube         bool      `json:"kube"`          // kubernetes 部署
	Namespace    string    `json:"namespace"`     // kuberntes 命名空间
}
