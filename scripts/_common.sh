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
    
    pushd "$install_dir/src" > /dev/null
    # Build as current user (root)
    env GOPATH="$install_dir/go" go build -o "$install_dir/bin/$app" .
    popd > /dev/null
    
    # Ensure correct ownership
    chown "$app:$app" "$install_dir/bin/$app"
    chmod 755 "$install_dir/bin/$app"
}

# Download application source
download_app_source() {
    local install_dir=$1
    
    # For now, we'll use the local source
    # In production, this would download from upstream
    mkdir -p "$install_dir/src"
    cp -r ../src/* "$install_dir/src/"
}

#=================================================
# FUTURE OFFICIAL HELPERS
#=================================================
