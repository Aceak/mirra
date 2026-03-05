# Mirra

一个轻量级的文件服务器，拥有现代化的 Web UI。使用 Go 构建，采用服务端渲染。

<div align="center">

**[English](README.md)** | **[中文](README_zh.md)**

</div>

## 功能特性

- 单二进制文件部署，无外部依赖
- 支持明/暗主题的 Web UI
- 目录浏览，支持实时搜索
- README.md 渲染，支持代码高亮

## 配置

### 自动配置
如果 `config.json` 不存在，将自动创建并使用当前目录作为根路径。

### 手动配置
复制 `config.example.json` 到 `config.json` 并编辑：

```bash
cp config.example.json config.json
```

编辑 `config.json`：

```json
{
  "server": {
    "name": "Mirra",
    "host": "0.0.0.0",
    "port": "8080",
    "favicon": ""
  },
  "share": {
    "root_path": "."
  },
  "appearance": {
    "theme": "auto",
    "show_hidden": false
  }
}
```

- `server.name`: Web UI 中显示的名称
- `server.host`: 监听地址（使用 `0.0.0.0` 监听所有接口）
- `server.port`: 监听端口
- `share.root_path`: 要提供服务的根目录
- `appearance.theme`: 默认主题（`light`、`dark` 或 `auto` 跟随系统）
- `appearance.show_hidden`: 是否显示隐藏文件（以 `.` 开头）

## 使用方法

```bash
# 使用默认配置文件运行（config.json）
./mirra

# 使用自定义配置文件路径运行
./mirra -c /path/to/config.json

# 显示版本信息
./mirra -v
```

## 源码编译

### 前置条件

- Go 1.25 或更高版本

```bash
# 编译 Go 二进制
go build -o mirra ./cmd/server

# 运行
./mirra
```

### 使用 Makefile

```bash
# 编译当前平台的二进制文件
make build

# 跨平台编译（linux, darwin, windows）
make cross-build

# 创建分发包
make dist

# 运行服务器
make run

# 格式化代码
make fmt

# 代码检查
make lint

# 清理构建产物
make clean
```

## 更新日志

### v0.0.3 (2026-03-05)

**修复**
- 修复代码块 Copy 按钮字体使用等宽字体
- 修复 Copy 按钮布局问题导致代码下移
- 修复 Clipboard API 回退方案以提升浏览器兼容性

### v0.0.2 (2026-03-05)

**新增**
- 添加 `-c` 标志以指定自定义配置文件路径
- 添加 Prism.js 语法高亮支持 30+ 种编程语言

**变更**
- 将所有代码注释和用户可见消息改为英文
- 优化页面标题显示逻辑（移除冗余服务器名称元素）
- 改进主题切换以保持 Prism.js 高亮效果

更多历史版本请查看 [CHANGELOG.md](CHANGELOG.md)。

## TODO

### 已完成
- [x] 文件服务器和目录浏览
- [x] 明/暗主题支持
- [x] README.md 渲染
- [x] 实时搜索
- [x] 代码块语法高亮

### 计划中
- [ ] 优化移动端 WebUI 体验
- [ ] 代码预览功能
- [ ] 文件 URL 复制和文件夹打包下载
- [ ] 多线程下载支持

## 许可证

MIT 许可证 - 详情查看 [LICENSE](LICENSE)。
