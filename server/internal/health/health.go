package health

// Status 表示应用当前的基础运行状态。
type Status string

const (
	// StatusOK 表示应用进程正常运行。
	StatusOK Status = "ok"

	// StatusNotReady 表示应用还没有准备好提供完整服务。
	StatusNotReady Status = "not_ready"
)

// Response 是健康检查接口统一返回的数据结构。
type Response struct {
	Status Status `json:"status"`
}
