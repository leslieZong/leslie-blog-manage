package password

import "golang.org/x/crypto/bcrypt"

// DefaultCost 是 bcrypt 默认的计算成本。
//
// bcrypt 的 cost 越高，计算密码哈希所需要的时间越长，
// 同时暴力破解的成本也越高。
//
// 对于当前项目，我们直接使用 bcrypt 官方提供的默认值。
// 不建议初学阶段自己随意修改。
const DefaultCost = bcrypt.DefaultCost

// Hash 对用户输入的明文密码进行哈希。
//
// 注意：
//
// 这个函数不会返回原密码。
// 它返回的是 bcrypt 生成的密码哈希字符串。
//
// 例如：
//
//	输入：
//	123456
//
//	输出：
//	$2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//
// 参数：
//
//	plainPassword：用户输入的明文密码
//
// 返回值：
//
//	string：密码哈希
//	error：哈希过程中出现的错误
func Hash(plainPassword string) (string, error) {
	// bcrypt.GenerateFromPassword 接收的是 []byte。
	//
	// Go 中 string 和 []byte 可以进行转换：
	//
	// string
	//   ↓
	// []byte
	//
	// 因为密码本质上需要按照字节参与密码算法计算。
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		DefaultCost,
	)

	// Go 不使用 try/catch。
	//
	// 函数执行过程中出现错误，
	// 通常通过返回值返回 error。
	if err != nil {
		return "", err
	}

	// bcrypt 返回的是 []byte。
	//
	// 我们转换成 string，
	// 因为数据库中的 password_hash 字段是 VARCHAR。
	return string(hash), nil
}

// Compare 比较明文密码和数据库中的密码哈希。
//
// 注意参数顺序：
//
//	hashedPassword：数据库中保存的密码哈希
//	plainPassword：用户登录时输入的明文密码
//
// 如果密码正确：
//
//	return true
//
// 如果密码错误：
//
//	return false
//
// 为什么这里不直接返回 error？
//
// 因为对于登录场景来说，
// “密码错误”本身就是一个正常的业务结果，
// 并不是系统异常。
func Compare(
	hashedPassword string,
	plainPassword string,
) bool {
	// bcrypt.CompareHashAndPassword 会完成：
	//
	// 1. 解析 hashedPassword
	// 2. 读取其中保存的 salt 和 cost
	// 3. 使用相同参数计算 plainPassword
	// 4. 比较最终结果
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)

	// err == nil 表示密码匹配。
	return err == nil
}
