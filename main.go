package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
)




// FileItem 文件项


func main() {
	// 命令行参数解析
	showVersion := flag.Bool("v", false, "显示版本信息")
	flag.Parse()

	// 加载配置文件
	if err := loadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// 如果指定了 -v 参数，显示版本信息并退出
	if *showVersion {
		fmt.Printf("%s version %s (commit: %s, built: %s)\n", config.Server.Name, version, commitHash, buildTime)
		return
	}


	// 设置静态文件服务（使用嵌入的文件系统）
	// 创建static子目录的文件系统
	staticSubFS, err := fs.Sub(StaticFS, "static")
	if err != nil {
		fmt.Printf("Error creating static sub filesystem: %v\n", err)
		return
	}

	// 创建静态文件服务器
	staticHandler := http.FileServer(http.FS(staticSubFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// 设置路由
	http.HandleFunc("/", handleRequest)

	addr := config.Server.Host + ":" + config.Server.Port
	fmt.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}







