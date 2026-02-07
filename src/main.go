package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get configuration from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	domain := os.Getenv("DOMAIN")
	pathURL := os.Getenv("PATH_URL")

	// Log startup information
	log.Printf("Starting Go application on port %s", port)
	log.Printf("Domain: %s, Path: %s", domain, pathURL)

	// Set up routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/api/info", handleInfo)

	// Create server with timeouts
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	log.Printf("Server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// handleHome serves the home page
func handleHome(w http.ResponseWriter, r *http.Request) {
	domain := os.Getenv("DOMAIN")
	pathURL := os.Getenv("PATH_URL")

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Application - YunoHost</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: #333;
        }
        .container {
            background: white;
            border-radius: 10px;
            padding: 40px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
        }
        h1 {
            color: #667eea;
            margin-top: 0;
        }
        .info {
            background: #f5f7fa;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .info-item {
            margin: 10px 0;
        }
        .label {
            font-weight: bold;
            color: #667eea;
        }
        a {
            color: #667eea;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Welcome to Your Go Application</h1>
        <p>This is a skeleton YunoHost application built with Go!</p>
        
        <div class="info">
            <div class="info-item">
                <span class="label">Domain:</span> %s
            </div>
            <div class="info-item">
                <span class="label">Path:</span> %s
            </div>
            <div class="info-item">
                <span class="label">Server Time:</span> %s
            </div>
        </div>
        
        <h2>Available Endpoints</h2>
        <ul>
            <li><a href="/health">/health</a> - Health check endpoint</li>
            <li><a href="/api/info">/api/info</a> - JSON API info</li>
        </ul>
        
        <h2>Next Steps</h2>
        <p>Edit <code>src/main.go</code> to customize your application and add new functionality!</p>
    </div>
</body>
</html>`, domain, pathURL, time.Now().Format(time.RFC1123))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// handleHealth provides a health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

// handleInfo provides API information
func handleInfo(w http.ResponseWriter, r *http.Request) {
	domain := os.Getenv("DOMAIN")
	pathURL := os.Getenv("PATH_URL")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
	"service": "YunoHost Go Application",
	"version": "1.0.0",
	"domain": "%s",
	"path": "%s",
	"timestamp": "%s"
}`, domain, pathURL, time.Now().Format(time.RFC3339))
}
