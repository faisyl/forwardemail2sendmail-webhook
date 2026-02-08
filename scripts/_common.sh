#!/bin/bash

#=================================================
# COMMON VARIABLES
#=================================================

go_version="1.21"

#=================================================
# PERSONAL HELPERS
#=================================================

# Build the Go application
build_go_app() {
    local install_dir=$1
    local app=$2
 
    mkdir -p "$install_dir/go-cache"
    
    pushd "$install_dir/sources" > /dev/null
    # Build as current user (root) with explicit cache location
    env GOPATH="$install_dir/go" GOCACHE="$install_dir/go-cache" go build -o "$install_dir/bin/$app" .
    popd > /dev/null
    
    # Clean up build cache and GOPATH files
    rm -rf "$install_dir/go-cache"
    rm -rf "$install_dir/go"
    
    # Ensure correct ownership
    chown "$app:$app" "$install_dir/bin/$app"
    chmod 755 "$install_dir/bin/$app"
}

# Download application source
download_app_source() {
    local install_dir=$1
    
    # For now, we'll use the local source
    # In production, this would download from upstream
    mkdir -p "$install_dir/sources"
    cp -r ../sources/* "$install_dir/sources/"
}

#=================================================
# FUTURE OFFICIAL HELPERS
#=================================================
