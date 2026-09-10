package permission

import (
	"fmt"
	"strings"
)

// Parse 将：
//
// post:create
//
// 转换成：
//
// resource = post
// action   = create
func Parse(value string) (string, string, error) {

	parts := strings.SplitN(
		value,
		":",
		2,
	)

	// Permission 必须包含 resource:action。
	if len(parts) != 2 {
		return "", "", fmt.Errorf(
			"invalid permission: %s",
			value,
		)
	}

	resource := strings.TrimSpace(parts[0])
	action := strings.TrimSpace(parts[1])

	// resource 和 action 都不能为空。
	if resource == "" || action == "" {
		return "", "", fmt.Errorf(
			"invalid permission: %s",
			value,
		)
	}

	return resource, action, nil
}
