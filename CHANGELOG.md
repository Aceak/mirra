# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2026-03-04

### Added
- 第一个正式版本
- 单二进制部署，无需额外依赖
- GitHub 风格亮/暗主题支持，可跟随系统偏好
- 目录列表显示，带文件类型图标
- README.md 渲染（类似 GitHub）
- 面包屑导航
- 统计信息显示（目录数、文件数、总大小）
- 文件列表可排序（名称、大小、修改时间）
- 实时搜索过滤
- SPA 风格导航（无刷新切换目录）
- 代码块语法高亮（Prism.js 本地化）
- 响应式设计，支持移动端

### Changed
- 重构项目结构，使用 Go 标准布局（cmd/, internal/）
- 分离 template.html 中的 CSS 和 JS 到独立文件
- 所有静态资源本地化，不再依赖 CDN
- 优化 SPA 导航逻辑，支持所有目录链接
- 优化 README 区域动态渲染逻辑
