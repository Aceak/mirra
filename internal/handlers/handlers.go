package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Aceak/mirra/internal/config"
	"github.com/Aceak/mirra/internal/types"
	"github.com/Aceak/mirra/internal/utils"
)

// InitTemplate 初始化HTML模板
func InitTemplate(staticFS embed.FS) (*template.Template, error) {
	templateBytes, err := fs.ReadFile(staticFS, "static/template.html")
	if err != nil {
		return nil, fmt.Errorf("error reading template: %v", err)
	}
	return template.New("page").Parse(string(templateBytes))
}

// HandleRequest 处理HTTP请求
func HandleRequest(w http.ResponseWriter, r *http.Request, cfg *config.Config, tmpl *template.Template) {
	// 清理路径，防止目录遍历
	requestPath := filepath.Clean(r.URL.Path)
	if requestPath == "/" {
		requestPath = "."
	} else {
		requestPath = requestPath[1:] // 移除开头的斜杠
	}

	// 构建完整路径
	fullPath := filepath.Join(cfg.Share.RootPath, requestPath)

	// 检查路径是否存在
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		// 目录：显示目录列表
		serveDirectory(w, r, requestPath, fullPath, cfg, tmpl)
	} else {
		// 文件：直接提供下载
		http.ServeFile(w, r, fullPath)
	}
}

func serveDirectory(w http.ResponseWriter, _ *http.Request, requestPath, fullPath string, cfg *config.Config, tmpl *template.Template) {
	// 读取目录内容
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	// 构建文件项列表
	var items []types.FileItem
	var totalDirs, totalFiles int
	var totalSize int64
	var latestModTime string
	var latestModTimeUnix int64 = 0

	for _, entry := range entries {
		// 跳过隐藏文件（如果配置为不显示）
		if !cfg.Appearance.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		itemPath := filepath.Join(requestPath, entry.Name())
		isDir := entry.IsDir()
		size := info.Size()
		sizeHuman := utils.FormatSize(size)
		modTime := info.ModTime()

		// 追踪最新的修改时间
		if modTime.Unix() > latestModTimeUnix {
			latestModTimeUnix = modTime.Unix()
			latestModTime = modTime.Format("2006-01-02 15:04:05")
		}

		if isDir {
			totalDirs++
			sizeHuman = "-"
		} else {
			totalFiles++
			totalSize += size
		}

		items = append(items, types.FileItem{
			Name:        entry.Name(),
			Path:        itemPath,
			IsDirectory: isDir,
			Size:        size,
			SizeHuman:   sizeHuman,
			Modified:    modTime.Format("2006-01-02 15:04:05"),
			IconClass:   utils.GetIconClass(entry.Name(), isDir),
		})
	}

	// 排序：目录在前，文件在后，按名称升序
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsDirectory && !items[j].IsDirectory {
			return true
		}
		if !items[i].IsDirectory && items[j].IsDirectory {
			return false
		}
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	// 检查README.md
	var readmeHTML template.HTML
	hasReadme := false
	readmePath := filepath.Join(fullPath, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		content, err := os.ReadFile(readmePath)
		if err == nil {
			// 渲染Markdown为HTML
			html := utils.RenderMarkdown(content)
			readmeHTML = template.HTML(html)
			hasReadme = true
		}
	}

	// 构建面包屑导航
	breadcrumbs := buildBreadcrumbs(requestPath)

	// 确定主题
	isDarkTheme := cfg.Appearance.Theme == "dark"
	if cfg.Appearance.Theme == "" || cfg.Appearance.Theme == "auto" {
		// 这里可以添加自动检测系统主题的逻辑
		// 暂时默认使用深色主题
		isDarkTheme = true
	}

	// 准备页面数据
	// 计算页面标题：根目录使用配置的 Name，子目录使用当前文件夹名
	pageTitle := cfg.Server.Name
	if requestPath != "." {
		// 提取当前目录名
		parts := strings.Split(requestPath, string(filepath.Separator))
		if len(parts) > 0 {
			pageTitle = parts[len(parts)-1]
		}
	}

	data := types.PageData{
		ServerName:    cfg.Server.Name,
		ServerFavicon: cfg.Server.Favicon,
		PageTitle:     pageTitle,
		Path:          requestPath,
		Items:         items,
		HasReadme:     hasReadme,
		ReadmeHTML:    readmeHTML,
		Stats: types.Stats{
			Directories: totalDirs,
			Files:       totalFiles,
			TotalSize:   utils.FormatSize(totalSize),
		},
		IsDarkTheme:   isDarkTheme,
		ShowHidden:    cfg.Appearance.ShowHidden,
		Breadcrumbs:   breadcrumbs,
		SortOrder:     "asc",
		LatestModTime: latestModTime,
	}

	// 渲染模板
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

func buildBreadcrumbs(path string) []types.Breadcrumb {
	var breadcrumbs []types.Breadcrumb

	// 总是包含根目录
	breadcrumbs = append(breadcrumbs, types.Breadcrumb{Name: "/", Path: ".", Last: path == "."})

	if path == "." {
		// 根目录，只有 "/" 一项
		return breadcrumbs
	}

	parts := strings.Split(path, string(filepath.Separator))
	for i, part := range parts {
		crumbPath := strings.Join(parts[:i+1], string(filepath.Separator))
		last := i == len(parts)-1
		breadcrumbs = append(breadcrumbs, types.Breadcrumb{
			Name: part,
			Path: crumbPath,
			Last: last,
		})
	}

	return breadcrumbs
}
