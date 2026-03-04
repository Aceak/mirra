package main

// 基础版本号
const baseVersion = "v0.0.1"

// 版本信息变量，由构建时注入
var (
	version    = baseVersion + ".dev"
	commitHash = "unknown"
	buildTime  = "unknown"
)