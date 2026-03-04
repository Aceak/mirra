package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/Aceak/mirra/internal/config"
	"github.com/Aceak/mirra/internal/handlers"
	"github.com/Aceak/mirra/internal/version"
)

//go:embed static/template.html static/css/* static/js/* static/webfonts/*
var StaticFS embed.FS

func main() {
	// 命令行参数解析
	showVersion := flag.Bool("v", false, "显示版本信息")
	flag.Parse()

	// 加载配置文件
	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		fmt.Printf("Error loading config: %v\n", cfgErr)
		return
	}

	// 如果指定了 -v 参数，显示版本信息并退出
	if *showVersion {
		fmt.Printf("%s\n", version.FormatVersion())
		return
	}

	// 初始化模板
	tmpl, tmplErr := handlers.InitTemplate(StaticFS)
	if tmplErr != nil {
		fmt.Printf("Error initializing template: %v\n", tmplErr)
		return
	}

	// 设置静态文件服务（使用嵌入的文件系统）
	// 创建static子目录的文件系统
	staticSubFS, subErr := fs.Sub(StaticFS, "static")
	if subErr != nil {
		fmt.Printf("Error creating static sub filesystem: %v\n", subErr)
		return
	}

	// 创建静态文件服务器
	staticHandler := http.FileServer(http.FS(staticSubFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// 设置路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRequest(w, r, cfg, tmpl)
	})

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	fmt.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
