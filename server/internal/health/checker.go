package health

import "context"

// Checker 用于检查某个基础设施是否已经准备好。
//
// 例如：
//
//	MySQL
//	Redis
//
// 只要实现这个接口，
// 就可以加入 Readiness 检查。
type Checker interface {

	// Check 执行健康检查。
	//
	// 返回 nil：
	//   表示依赖正常。
	//
	// 返回 error：
	//   表示依赖异常。
	Check(ctx context.Context) error
}
