// Package utils 提供 Markdown 渲染功能
package utils

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
	"regexp"
	"strings"

	markdownhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownRenderer Markdown 渲染器
type MarkdownRenderer struct {
	extensions parser.Extensions
	htmlFlags  markdownhtml.Flags
	enableCodeCopy bool
}

// MarkdownOptions Markdown 渲染选项
type MarkdownOptions struct {
	// 启用任务列表语法 (- [ ] 和 - [x])
	EnableTaskLists bool
	// 启用代码块复制按钮
	EnableCodeCopy bool
	// 链接在新窗口打开
	TargetBlank bool
	// 启用脚注
	EnableFootnotes bool
}

// DefaultOptions 默认选项
var DefaultOptions = MarkdownOptions{
	EnableTaskLists:  true,
	EnableCodeCopy:   true,
	TargetBlank:      true,
	EnableFootnotes:  false,
}

// NewMarkdownRenderer 创建新的 Markdown 渲染器
func NewMarkdownRenderer(opts MarkdownOptions) *MarkdownRenderer {
	// 设置解析器扩展
	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.NoIntraEmphasis |
		parser.DefinitionLists

	if opts.EnableFootnotes {
		extensions |= parser.Footnotes
	}

	// 设置 HTML 渲染标志
	htmlFlags := markdownhtml.CommonFlags | markdownhtml.Smartypants
	if opts.TargetBlank {
		htmlFlags |= markdownhtml.HrefTargetBlank
	}

	return &MarkdownRenderer{
		extensions:     extensions,
		htmlFlags:      htmlFlags,
		enableCodeCopy: opts.EnableCodeCopy,
	}
}

// Render 将 Markdown 内容渲染为 HTML
func (mr *MarkdownRenderer) Render(content []byte) string {
	// 预处理：转换特殊语法
	content = mr.preprocess(content)

	// 每次渲染时创建新的 Parser（Parser 不可重复使用）
	p := parser.NewWithExtensions(mr.extensions)
	doc := p.Parse(content)

	// 创建新的渲染器
	rendererOpts := markdownhtml.RendererOptions{
		Flags:          mr.htmlFlags,
		RenderNodeHook: createCodeBlockHook(mr.enableCodeCopy),
	}
	renderer := markdownhtml.NewRenderer(rendererOpts)

	// 渲染为 HTML
	htmlContent := string(markdown.Render(doc, renderer))

	// 后处理：转换任务列表等
	htmlContent = mr.postprocess(htmlContent)

	return htmlContent
}

// preprocess 预处理 Markdown 内容
func (mr *MarkdownRenderer) preprocess(content []byte) []byte {
	// 处理 HTML 块内的 Markdown 语法
	content = convertInlineMarkdown(content)
	// 处理任务列表语法
	content = convertTaskLists(content)
	return content
}

// postprocess 后处理 HTML 输出
func (mr *MarkdownRenderer) postprocess(htmlContent string) string {
	// 处理任务列表
	htmlContent = postProcessTaskLists(htmlContent)
	return htmlContent
}

// createCodeBlockHook 创建代码块渲染钩子
func createCodeBlockHook(enableCodeCopy bool) markdownhtml.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		codeBlock, ok := node.(*ast.CodeBlock)
		if !ok {
			return ast.WalkStatus(0), false
		}

		if !entering {
			return ast.SkipChildren, true
		}

		// 获取语言
		var lang string
		if codeBlock.Info != nil {
			lang = strings.TrimSpace(string(codeBlock.Info))
		}
		if lang == "" {
			lang = "plaintext"
		}

		// 渲染代码块
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "<pre><code class=\"language-%s\">", lang)
		escaped := html.EscapeString(string(codeBlock.Literal))
		buf.WriteString(escaped)
		buf.WriteString("</code></pre>")

		// 如果启用复制按钮，添加按钮
		if enableCodeCopy {
			fmt.Fprintf(w, `<div class="code-block-wrapper"><button class="code-copy-btn" onclick="copyCode(this)">Copy</button>%s</div>`, buf.String())
		} else {
			w.Write(buf.Bytes())
		}

		return ast.SkipChildren, true
	}
}

// ============================================================================
// 任务列表处理
// ============================================================================

var taskListRegex = regexp.MustCompile(`(?m)^([ ]*)- \[([ x])\] `)

// convertTaskLists 将任务列表语法转换为 HTML 注释标记
func convertTaskLists(content []byte) []byte {
	return taskListRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := taskListRegex.FindSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		indent := parts[1]
		checked := parts[2]
		if string(checked) == "x" {
			return []byte(string(indent) + "- <!--task: x--> ")
		}
		return []byte(string(indent) + "- <!--task:  --> ")
	})
}

// postProcessTaskLists 后处理任务列表 HTML
func postProcessTaskLists(htmlContent string) string {
	// 替换任务标记为复选框
	htmlContent = strings.ReplaceAll(htmlContent, "<!--task: x-->",
		`<input type="checkbox" checked disabled class="task-list-checkbox"><!--task-item:completed-->`)
	htmlContent = strings.ReplaceAll(htmlContent, "<!--task:  -->",
		`<input type="checkbox" disabled class="task-list-checkbox"><!--task-item:-->`)

	// 给 <li> 添加类名
	liCompletedRegex := regexp.MustCompile(`<li>(<input[^>]*>)<!--task-item:completed-->(\s*.*)</li>`)
	htmlContent = liCompletedRegex.ReplaceAllString(htmlContent, `<li class="task-list-item completed">$1$2</li>`)

	liRegex := regexp.MustCompile(`<li>(<input[^>]*>)<!--task-item:-->(\s*.*)</li>`)
	htmlContent = liRegex.ReplaceAllString(htmlContent, `<li class="task-list-item">$1$2</li>`)

	return htmlContent
}

// ============================================================================
// HTML 块内 Markdown 处理
// ============================================================================

var inlineMarkdownRegex = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

// convertInlineMarkdown 将 HTML 块内的 Markdown 链接转换为 HTML
func convertInlineMarkdown(content []byte) []byte {
	divRegex := regexp.MustCompile(`(<div[^>]*>)([\s\S]*?)(</div>)`)

	return divRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		parts := divRegex.FindSubmatch(match)
		if len(parts) != 4 {
			return match
		}

		openTag := parts[1]
		innerContent := parts[2]
		closeTag := parts[3]

		// 转换 Markdown 链接
		converted := inlineMarkdownRegex.ReplaceAll(innerContent, []byte(`<a href="$2">$1</a>`))

		// 转换粗体
		boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
		converted = boldRegex.ReplaceAll(converted, []byte(`<strong>$1</strong>`))

		result := append([]byte{}, openTag...)
		result = append(result, converted...)
		result = append(result, closeTag...)
		return result
	})
}

// ============================================================================
// 便捷函数
// ============================================================================

// globalRenderer 全局渲染器实例
var globalRenderer = NewMarkdownRenderer(DefaultOptions)

// RenderMarkdown 将 Markdown 内容渲染为 HTML（便捷函数）
func RenderMarkdown(content []byte) string {
	return globalRenderer.Render(content)
}

// RenderMarkdownString 将 Markdown 字符串渲染为 HTML 字符串
func RenderMarkdownString(content string) string {
	return RenderMarkdown([]byte(content))
}

// RenderMarkdownHTML 将 Markdown 内容渲染为 template.HTML（用于 Go 模板）
func RenderMarkdownHTML(content []byte) template.HTML {
	return template.HTML(RenderMarkdown(content))
}

// ============================================================================
// 代码块复制按钮 JavaScript
// ============================================================================

// CopyCodeScript 代码块复制按钮的 JavaScript 代码
const CopyCodeScript = `
function copyCode(button) {
    const wrapper = button.parentElement;
    const code = wrapper.querySelector('code');
    const text = code.textContent;

    navigator.clipboard.writeText(text).then(function() {
        button.textContent = 'Copied!';
        button.classList.add('copied');
        setTimeout(function() {
            button.textContent = 'Copy';
            button.classList.remove('copied');
        }, 2000);
    }).catch(function() {
        // Fallback for older browsers
        const textarea = document.createElement('textarea');
        textarea.value = text;
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
        button.textContent = 'Copied!';
        button.classList.add('copied');
        setTimeout(function() {
            button.textContent = 'Copy';
            button.classList.remove('copied');
        }, 2000);
    });
}
`
