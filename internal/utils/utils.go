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

// convertTaskLists 将任务列表语法转换为 HTML 注释标记
// 这样 Markdown 解析器会解析列表项内的 Markdown 语法
func convertTaskLists(content []byte) []byte {
	// 匹配 - [ ] 或 - [x] 开头的行
	taskRegex := regexp.MustCompile(`(?m)^([ ]*)- \[([ x])\] `)
	return taskRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := taskRegex.FindSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		indent := parts[1]
		checked := parts[2]
		// 使用特殊标记 <!--task: x--> 或 <!--task:  -->
		if string(checked) == "x" {
			return []byte(string(indent) + "- <!--task: x--> ")
		}
		return []byte(string(indent) + "- <!--task:  --> ")
	})
}

// postProcessHTML 在 HTML 渲染后处理任务列表标记
// 将 <!--task: x--> 替换为复选框，并给父 <li> 添加 task-list-item 类
func postProcessHTML(htmlContent string) string {
	// 先将 <!--task: x--> 替换为带有标记的复选框
	// 标记 <!--task-item:xxx--> 用于后续给 <li> 添加类名
	htmlContent = strings.ReplaceAll(htmlContent, "<!--task: x-->", `<input type="checkbox" checked disabled class="task-list-checkbox"><!--task-item:completed-->`)
	htmlContent = strings.ReplaceAll(htmlContent, "<!--task:  -->", `<input type="checkbox" disabled class="task-list-checkbox"><!--task-item:-->`)

	// 使用正则匹配：<li><input...><!--task-item:xxx--> text</li>
	// 捕获复选框和后续文本，只移除 <!--task-item:xxx--> 标记
	liCompletedRegex := regexp.MustCompile(`<li>(<input[^>]*>)<!--task-item:completed-->(\s*.*)</li>`)
	htmlContent = liCompletedRegex.ReplaceAllString(htmlContent, `<li class="task-list-item completed">$1$2</li>`)

	liRegex := regexp.MustCompile(`<li>(<input[^>]*>)<!--task-item:-->(\s*.*)</li>`)
	htmlContent = liRegex.ReplaceAllString(htmlContent, `<li class="task-list-item">$1$2</li>`)

	return htmlContent
}

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

	htmlContent := string(markdown.Render(doc, renderer))
	// 后处理：替换任务列表标记为复选框
	return postProcessHTML(htmlContent)
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
