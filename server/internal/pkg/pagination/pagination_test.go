package pagination

import "testing"

// TestNewParams 测试分页参数规范化。
func TestNewParams(t *testing.T) {

	// 模拟用户传入：
	//
	// page = 0
	// pageSize = 1000
	//
	// 这两个参数都不符合我们的规范。
	params := NewParams(0, 1000)

	// page 应该被修正成默认值 1。
	if params.Page != DefaultPage {
		t.Fatalf(
			"expected page %d, got %d",
			DefaultPage,
			params.Page,
		)
	}

	// pageSize 超过最大值，
	// 应该被限制为 MaxPageSize。
	if params.PageSize != MaxPageSize {
		t.Fatalf(
			"expected pageSize %d, got %d",
			MaxPageSize,
			params.PageSize,
		)
	}
}
