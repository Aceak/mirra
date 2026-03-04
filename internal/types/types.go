package types

import "html/template"

// FileItem 文件项
type FileItem struct {
	Name        string
	Path        string
	IsDirectory bool
	Size        int64
	SizeHuman   string
	Modified    string
	IconClass   string
}

// PageData 页面数据
type PageData struct {
	ServerName    string
	ServerFavicon string
	PageTitle     string
	Path          string
	Items         []FileItem
	HasReadme     bool
	ReadmeHTML    template.HTML
	Stats         Stats
	IsDarkTheme   bool
	ShowHidden    bool
	Breadcrumbs   []Breadcrumb
	SortOrder     string
	LatestModTime string
}

// Stats 统计信息
type Stats struct {
	Directories int
	Files       int
	TotalSize   string
}

// Breadcrumb 面包屑导航
type Breadcrumb struct {
	Name string
	Path string
	Last bool
}
