package password

import "testing"

// TestHash 测试密码哈希功能。
func TestHash(t *testing.T) {
	// 准备测试数据。
	plainPassword := "123456"

	// 调用我们刚才编写的 Hash 函数。
	hash, err := Hash(plainPassword)

	// 如果发生错误，测试失败。
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	// 哈希结果不应该为空。
	if hash == "" {
		t.Fatal("Hash() returned empty hash")
	}

	// 哈希结果不应该等于原始密码。
	if hash == plainPassword {
		t.Fatal("hash should not equal plain password")
	}
}

// TestHashProducesDifferentResults 测试同一个密码
// 多次 Hash 后结果是否不同。
func TestHashProducesDifferentResults(t *testing.T) {
	plainPassword := "123456"

	hash1, err := Hash(plainPassword)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	hash2, err := Hash(plainPassword)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	// bcrypt 会随机生成 salt。
	//
	// 因此即使密码相同，
	// 每次生成的 hash 通常也不同。
	if hash1 == hash2 {
		t.Fatal("same password should produce different hashes")
	}
}

// TestCompare 测试密码比较功能。
func TestCompare(t *testing.T) {
	plainPassword := "123456"

	hash, err := Hash(plainPassword)
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	// 正确密码应该验证成功。
	if !Compare(hash, plainPassword) {
		t.Fatal("Compare() should return true for correct password")
	}

	// 错误密码应该验证失败。
	if Compare(hash, "wrong-password") {
		t.Fatal("Compare() should return false for wrong password")
	}
}
