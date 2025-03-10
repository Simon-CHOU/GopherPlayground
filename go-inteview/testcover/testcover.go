// 可以使用 go test -cover 命令生成代码覆盖率报告，用户评估测试覆盖率。
package testcover

// Add 返回两个整数的和
func Add(a, b int) int {
	return a + b
}

// Subtract 返回两个整数的差
func Subtract(a, b int) int {
	return a - b
}

// Multiply 返回两个整数的乘积
func Multiply(a, b int) int {
	return a * b
}

// Divide 返回两个整数的商
// 如果除数为0，返回0
func Divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}
