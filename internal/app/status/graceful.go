package status

// Graceful 优雅状态控制
type Graceful struct{}

// Shutdown 优雅下线
func (g *Graceful) Shutdown() bool {
	return true
}

// OutOfService 服务下线
func (g *Graceful) OutOfService() bool {
	return true
}

// Up 服务上线
func (g *Graceful) Up() bool {
	return true
}
