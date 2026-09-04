package ulid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// New 生成一个新的 ULID。
//
// ULID 的最终格式类似：
//
//	01K4XXXXXXXXXXXXXXX
//
// 长度固定为 26 个字符。
//
// 这里使用 crypto/rand 作为随机源，
// 而不是普通的 math/rand。
//
// 因为 ID 属于系统级唯一标识，
// 我们希望随机部分具有足够好的随机性。
func New() string {

	// 使用当前时间作为 ULID 的时间部分。
	timestamp := ulid.Timestamp(time.Now())

	// 使用 crypto/rand 生成随机部分。
	entropy := ulid.Monotonic(
		rand.Reader,
		0,
	)

	// 生成 ULID。
	id := ulid.MustNew(
		timestamp,
		entropy,
	)

	// ULID 本质上是一个类型，
	// 最终转换成 string 返回。
	return id.String()
}
