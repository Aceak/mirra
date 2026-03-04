package version

// BaseVersion 基础版本号
const BaseVersion = "0.0.1"

// CommitHash 提交哈希，由构建时注入
var CommitHash = "unknown"

// FormatVersion 返回格式化的版本字符串
// 格式：Mirra version x.x.x, Build abc1234
func FormatVersion() string {
	return "Mirra version " + BaseVersion + ", Build " + CommitHash
}
