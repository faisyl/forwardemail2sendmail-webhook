# YunoHost Go Application

A skeleton project for creating YunoHost-installable applications written in Go.

## Overview

This is a template/starter project for building web applications in Go that can be easily installed and managed on YunoHost servers. It includes all the necessary packaging files and scripts required by YunoHost v2 packaging format.

## Features

- ✅ YunoHost v2 packaging format (`manifest.toml`)
- ✅ Complete installation, removal, backup, restore, and upgrade scripts
- ✅ Systemd service integration
- ✅ Nginx reverse proxy configuration
- ✅ Basic HTTP server with health check endpoint
- ✅ Environment-based configuration
- ✅ Multi-instance support

## Project Structure

```
.
├── manifest.toml              # YunoHost app manifest (v2 format)
├── conf/
│   ├── nginx.conf            # Nginx reverse proxy configuration
│   └── systemd.service       # Systemd service template
├── scripts/
│   ├── _common.sh            # Shared helper functions
│   ├── install               # Installation script
│   ├── remove                # Removal script
│   ├── backup                # Backup script
│   ├── restore               # Restore script
│   ├── upgrade               # Upgrade script
│   └── change_url            # Domain/path change script
├── src/
│   ├── main.go               # Main Go application
│   └── go.mod                # Go module file
├── doc/
│   └── DESCRIPTION.md        # Extended description for YunoHost
├── LICENSE                   # License file
└── README.md                 # This file
```

## Development

### Prerequisites

- Go 1.21 or higher
- Basic understanding of YunoHost packaging

### Building Locally

```bash
cd src
go build -o ../bin/goapp .
```

### Running Locally

```bash
PORT=8080 DOMAIN=localhost PATH_URL=/ ./bin/goapp
```

Then visit `http://localhost:8080` in your browser.

### Testing

```bash
# Health check endpoint
curl http://localhost:8080/health

# API info endpoint
curl http://localhost:8080/api/info
```

## Installation on YunoHost

### From Command Line

```bash
yunohost app install /path/to/this/directory
```

### Configuration

During installation, you'll be prompted for:

- **Domain**: The domain where the app will be accessible
- **Path**: The URL path (default: `/goapp`)
- **Admin**: The YunoHost user who will be the admin
- **Public access**: Whether visitors can access the app

## Customization

### Modifying the Application

1. Edit `src/main.go` to add your application logic
2. Add dependencies to `src/go.mod` as needed
3. Update `manifest.toml` with your app details:
   - Change `id`, `name`, and `description`
   - Update `upstream` section with your repository URL
   - Modify default port if needed

### Adding Features

- **Database**: Add database resources in `manifest.toml` resources section
- **LDAP/SSO**: Set `ldap = true` or `sso = true` in the integration section
- **Additional packages**: Add to `resources.apt.packages`

### Build Configuration

The build process is handled in `scripts/_common.sh` via the `build_go_app()` function. Customize this if you need:
- Additional build flags
- Cross-compilation
- Asset embedding
- Multi-binary builds

## YunoHost Helper Scripts

All scripts are located in the `scripts/` directory:

- **install**: Sets up the application for the first time
- **remove**: Completely removes the application
- **backup**: Creates a backup of app data and configuration
- **restore**: Restores from a backup
- **upgrade**: Updates to a new version
- **change_url**: Changes the domain or path

## Environment Variables

The application receives these environment variables from systemd:

- `PORT`: The port the application should listen on (managed by YunoHost)
- `DOMAIN`: The domain where the app is installed
- `PATH_URL`: The URL path where the app is accessible

## License

MIT License - see LICENSE file

## Resources

- [YunoHost v2 Packaging Documentation](https://yunohost.org/en/packaging_apps)
- [YunoHost App Helpers](https://yunohost.org/en/packaging_apps_helpers)
- [Go Documentation](https://go.dev/doc/)

## Contributing

Feel free to customize this template for your own applications. This is meant to be a starting point!
