package utils

import (
	"fmt"
	"html"
	"io"
	"path/filepath"
	"strings"

	markdownhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// RenderMarkdown 将 Markdown 内容渲染为 HTML
func RenderMarkdown(content []byte) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	htmlFlags := markdownhtml.CommonFlags | markdownhtml.HrefTargetBlank
	opts := markdownhtml.RendererOptions{
		Flags: htmlFlags,
		RenderNodeHook: renderCodeBlock,
	}
	renderer := markdownhtml.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

// renderCodeBlock 渲染代码块，添加语言类名供 prism.js 高亮
func renderCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	codeBlock, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.WalkStatus(0), false
	}

	if entering {
		var lang string
		if codeBlock.Info != nil {
			lang = strings.TrimSpace(string(codeBlock.Info))
		}
		if lang == "" {
			lang = "plaintext"
		}

		// 写入 pre 标签，添加 language-前缀供 prism.js 使用
		fmt.Fprintf(w, "<pre class=\"language-%s\"><code class=\"language-%s\">", lang, lang)

		// 转义并写入代码内容
		escaped := html.EscapeString(string(codeBlock.Literal))
		w.Write([]byte(escaped))
		fmt.Fprint(w, "</code></pre>")
	}

	return ast.SkipChildren, true
}

// GetIconClass 根据文件名获取图标类名
func GetIconClass(filename string, isDir bool) string {
	if isDir {
		return "fas fa-folder"
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "fas fa-file"
	}

	// 移除点号
	ext = ext[1:]

	switch ext {
	case "txt", "md", "markdown", "rst":
		return "fas fa-file-alt"
	case "go", "js", "ts", "py", "java", "cpp", "c", "rs", "rb", "php", "swift", "kt":
		return "fas fa-file-code"
	case "html", "htm", "xml", "json", "yaml", "yml", "toml":
		return "fas fa-file-code"
	case "css", "scss", "less", "sass":
		return "fas fa-file-code"
	case "pdf":
		return "fas fa-file-pdf"
	case "doc", "docx":
		return "fas fa-file-word"
	case "xls", "xlsx":
		return "fas fa-file-excel"
	case "ppt", "pptx":
		return "fas fa-file-powerpoint"
	case "zip", "rar", "7z", "tar", "gz", "bz2":
		return "fas fa-file-archive"
	case "jpg", "jpeg", "png", "gif", "bmp", "svg", "webp":
		return "fas fa-file-image"
	case "mp3", "wav", "ogg", "flac":
		return "fas fa-file-audio"
	case "mp4", "avi", "mov", "mkv", "webm":
		return "fas fa-file-video"
	case "exe", "bin", "app":
		return "fas fa-file-executable"
	case "db", "sqlite", "sql":
		return "fas fa-database"
	default:
		return "fas fa-file"
	}
}

// FormatSize 格式化文件大小
func FormatSize(size int64) string {
	if size == 0 {
		return "-"
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
