package cache

import "context"

// InvalidateHome 删除首页缓存。
//
// Home 是一个聚合数据：
// Post、Project、Category、TechStack
// 任意一个模块发生变化，都可能影响首页。
//
// 因此统一通过这个方法删除首页缓存，
// 避免业务层到处直接操作具体的 Redis Key。
func InvalidateHome(
	ctx context.Context,
	c Cache,
) error {

	return c.Delete(
		ctx,
		KeyHome,
	)
}
