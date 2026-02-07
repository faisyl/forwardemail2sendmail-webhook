# Go Application for YunoHost

This is a skeleton template for building web applications in Go that can be easily installed and managed on YunoHost servers.

## What This Does

This application provides a basic HTTP web server written in Go with:

- **Health monitoring**: Built-in health check endpoint for service monitoring
- **Environment-based configuration**: Automatically configures itself based on YunoHost settings
- **RESTful API**: JSON API endpoints ready for expansion
- **Modern web interface**: Clean, responsive web interface

## Features

- Simple HTTP server that can be extended with your own logic
- Multi-instance support - run multiple copies on different domains/paths
- Automatic backup and restore functionality
- Easy upgrades through YunoHost interface
- Systemd service management for reliability
- Nginx reverse proxy integration with SSL

## Use Cases

Perfect starting point for:

- Custom web applications
- API services
- Microservices
- Webhook receivers
- Internal tools and dashboards
- Automation services

## Technical Details

- **Language**: Go 1.21+
- **Web Server**: Built-in Go HTTP server
- **Process Management**: Systemd
- **Reverse Proxy**: Nginx
- **Security**: Runs as unprivileged system user with sandboxing

## Getting Started

After installation, the application will be available at your chosen domain and path. The default homepage provides links to available endpoints and basic information.

Customize the application by editing the source code in `/var/www/goapp/src/main.go` (or your chosen install directory).
