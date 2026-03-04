package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed static/css/* static/webfonts/* static/template.html
var StaticFS embed.FS

var tmpl *template.Template

func init() {
	// 加载模板
	templateBytes, err := fs.ReadFile(StaticFS, "static/template.html")
	if err != nil {
		panic(fmt.Sprintf("Error reading template: %v", err))
	}
	tmpl, err = template.New("page").Parse(string(templateBytes))
	if err != nil {
		panic(fmt.Sprintf("Error parsing template: %v", err))
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 清理路径，防止目录遍历
	requestPath := filepath.Clean(r.URL.Path)
	if requestPath == "/" {
		requestPath = "."
	} else {
		requestPath = requestPath[1:] // 移除开头的斜杠
	}

	// 构建完整路径
	fullPath := filepath.Join(config.Share.RootPath, requestPath)

	// 检查路径是否存在
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		// 目录：显示目录列表
		serveDirectory(w, r, requestPath, fullPath)
	} else {
		// 文件：直接提供下载
		http.ServeFile(w, r, fullPath)
	}
}

func serveDirectory(w http.ResponseWriter, _ *http.Request, requestPath, fullPath string) {
	// 读取目录内容
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	// 构建文件项列表
	var items []FileItem
	var totalDirs, totalFiles int
	var totalSize int64

	for _, entry := range entries {
		// 跳过隐藏文件（如果配置为不显示）
		if !config.Appearance.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		itemPath := filepath.Join(requestPath, entry.Name())
		isDir := entry.IsDir()
		size := info.Size()
		sizeHuman := formatSize(size)

		if isDir {
			totalDirs++
			sizeHuman = "-"
		} else {
			totalFiles++
			totalSize += size
		}

		items = append(items, FileItem{
			Name:        entry.Name(),
			Path:        itemPath,
			IsDirectory: isDir,
			Size:        size,
			SizeHuman:   sizeHuman,
			Modified:    info.ModTime().Format("2006-01-02 15:04:05"),
			IconClass:   getIconClass(entry.Name(), isDir),
		})
	}

	// 检查README.md
	var readmeHTML template.HTML
	hasReadme := false
	readmePath := filepath.Join(fullPath, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		content, err := os.ReadFile(readmePath)
		if err == nil {
			// 渲染Markdown为HTML
			html := renderMarkdown(content)
			readmeHTML = template.HTML(html)
			hasReadme = true
		}
	}

	// 构建面包屑导航
	breadcrumbs := buildBreadcrumbs(requestPath)

	// 确定主题
	isDarkTheme := config.Appearance.Theme == "dark"
	if config.Appearance.Theme == "" || config.Appearance.Theme == "auto" {
		// 这里可以添加自动检测系统主题的逻辑
		// 暂时默认使用深色主题
		isDarkTheme = true
	}

	// 准备页面数据
	data := PageData{
		ServerName:  config.Server.Name,
		Path:        requestPath,
		Items:       items,
		HasReadme:   hasReadme,
		ReadmeHTML:  readmeHTML,
		Stats: Stats{
			Directories: totalDirs,
			Files:       totalFiles,
			TotalSize:   formatSize(totalSize),
		},
		IsDarkTheme: isDarkTheme,
		ShowHidden:  config.Appearance.ShowHidden,
		Breadcrumbs: breadcrumbs,
	}

	// 渲染模板
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

func buildBreadcrumbs(path string) []Breadcrumb {
	var breadcrumbs []Breadcrumb

	// 总是包含根目录
	breadcrumbs = append(breadcrumbs, Breadcrumb{Name: "/", Path: ".", Last: path == "."})

	if path == "." {
		// 根目录，只有 "/" 一项
		return breadcrumbs
	}

	parts := strings.Split(path, string(filepath.Separator))
	for i, part := range parts {
		crumbPath := strings.Join(parts[:i+1], string(filepath.Separator))
		last := i == len(parts)-1
		// 目录名不加斜杠，让分隔符处理
		name := part
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Name: name,
			Path: crumbPath,
			Last: last,
		})
	}

	return breadcrumbs
}