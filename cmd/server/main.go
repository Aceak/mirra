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
	// Parse command line arguments
	showVersion := flag.Bool("v", false, "Show version information")
	configPath := flag.String("c", "config.json", "Specify config file path")
	flag.Parse()

	// Load configuration
	cfg, cfgErr := config.LoadConfig(*configPath)
	if cfgErr != nil {
		fmt.Printf("Error loading config: %v\n", cfgErr)
		return
	}

	// Show version information if -v flag is set
	if *showVersion {
		fmt.Printf("%s\n", version.FormatVersion())
		return
	}

	// Initialize template
	tmpl, tmplErr := handlers.InitTemplate(StaticFS)
	if tmplErr != nil {
		fmt.Printf("Error initializing template: %v\n", tmplErr)
		return
	}

	// Setup static file server (using embedded filesystem)
	// Create static subdirectory filesystem
	staticSubFS, subErr := fs.Sub(StaticFS, "static")
	if subErr != nil {
		fmt.Printf("Error creating static sub filesystem: %v\n", subErr)
		return
	}

	// Create static file server
	staticHandler := http.FileServer(http.FS(staticSubFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRequest(w, r, cfg, tmpl)
	})

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	fmt.Printf("Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
