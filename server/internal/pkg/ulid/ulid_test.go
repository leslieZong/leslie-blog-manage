package ulid

import (
	"testing"
)

// TestNew 测试 ULID 是否能够正常生成。
func TestNew(t *testing.T) {

	id := New()

	// ULID 应该是 26 个字符。
	if len(id) != 26 {
		t.Fatalf(
			"expected ULID length 26, got %d",
			len(id),
		)
	}
}

// TestNewProducesDifferentIDs 测试连续生成的 ID
// 是否不同。
func TestNewProducesDifferentIDs(t *testing.T) {

	id1 := New()
	id2 := New()

	if id1 == id2 {
		t.Fatal("generated ULIDs should be different")
	}
}
