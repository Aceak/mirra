package utils

import (
	"fmt"
	"html"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	markdownhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// HTML 块内的 Markdown 链接正则
var inlineMarkdownRegex = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

// 任务列表正则：匹配 - [ ] 或 - [x] 开头的行
var taskListRegex = regexp.MustCompile(`(?m)^([ ]*)- \[([ x])\] (.*)$`)

// RenderMarkdown 将 Markdown 内容渲染为 HTML
func RenderMarkdown(content []byte) string {
	// 预处理：将 HTML 块内的 Markdown 链接转换为 HTML
	content = convertInlineMarkdown(content)
	// 预处理：转换任务列表语法
	content = convertTaskLists(content)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoIntraEmphasis
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	htmlFlags := markdownhtml.CommonFlags | markdownhtml.HrefTargetBlank | markdownhtml.Smartypants
	opts := markdownhtml.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: renderCodeBlock,
	}
	renderer := markdownhtml.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

// convertInlineMarkdown 将 HTML 块内的 Markdown 链接转换为 HTML
func convertInlineMarkdown(content []byte) []byte {
	// 匹配 <div ...>...</div> 内的内容
	divRegex := regexp.MustCompile(`(<div[^>]*>)([\s\S]*?)(</div>)`)

	return divRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		// 提取开始标签、内容和结束标签
		parts := divRegex.FindSubmatch(match)
		if len(parts) != 4 {
			return match
		}

		openTag := parts[1]
		innerContent := parts[2]
		closeTag := parts[3]

		// 将内部的 Markdown 链接 [text](url) 转换为 HTML
		converted := inlineMarkdownRegex.ReplaceAll(innerContent, []byte(`<a href="$2">$1</a>`))

		// 处理粗体 **text**
		boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
		converted = boldRegex.ReplaceAll(converted, []byte(`<strong>$1</strong>`))

		result := append([]byte{}, openTag...)
		result = append(result, converted...)
		result = append(result, closeTag...)
		return result
	})
}

// convertTaskLists 将任务列表语法转换为 HTML
// - [ ] 转换为未选中复选框
// - [x] 转换为选中复选框
func convertTaskLists(content []byte) []byte {
	return taskListRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := taskListRegex.FindSubmatch(match)
		if len(parts) != 4 {
			return match
		}

		indent := string(parts[1])
		checked := string(parts[2])
		text := string(parts[3])

		var isChecked string
		var itemClass string
		if checked == "x" {
			isChecked = "checked"
			itemClass = "completed"
		} else {
			isChecked = ""
			itemClass = ""
		}

		// 计算嵌套层级
		indentLevel := len(indent) / 2
		var padding string
		if indentLevel > 0 {
			padding = strings.Repeat("  ", indentLevel)
		}

		// 生成 HTML
		if itemClass != "" {
			return []byte(fmt.Sprintf("%s<li class=\"task-list-item %s\"><input type=\"checkbox\" %s disabled>%s</li>", padding, itemClass, isChecked, text))
		}
		return []byte(fmt.Sprintf("%s<li class=\"task-list-item\"><input type=\"checkbox\" %s disabled>%s</li>", padding, isChecked, text))
	})
}

// renderCodeBlock 渲染代码块，添加语言类名供 prism.js 高亮
func renderCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	codeBlock, ok := node.(*ast.CodeBlock)
	if !ok {
		return ast.WalkStatus(0), false
	}

	if !entering {
		return ast.SkipChildren, true
	}

	var lang string
	if codeBlock.Info != nil {
		lang = strings.TrimSpace(string(codeBlock.Info))
	}
	if lang == "" {
		lang = "plaintext"
	}

	// 写入 pre 标签，code 标签添加 language-前缀供 prism.js 使用
	// Prism.js 标准格式：<pre><code class="language-xxx">
	fmt.Fprintf(w, "<pre><code class=\"language-%s\">", lang)

	// 转义并写入代码内容
	escaped := html.EscapeString(string(codeBlock.Literal))
	w.Write([]byte(escaped))
	fmt.Fprint(w, "</code></pre>")

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
