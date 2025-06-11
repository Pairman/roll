#!/usr/bin/env bash

TARGETS=(
    "linux x86_64"
    "linux aarch64"
    "linux riscv64"
    "linux loongarch64"

    "freebsd x86_64"

    "darwin x86_64"
    "darwin aarch64"

    "windows x86_64"
    "windows aarch64"
)

compare_versions() {
    local version="$1"
    local version_prev="$2"
    local v1 v2
    IFS='.' read -r -a v1 <<< "$version"
    IFS='.' read -r -a v2 <<< "$version_prev"
    for i in 0 1 2; do
        if (( $((10#${v1[i]:-0})) > $((10#${v2[i]:-0})) )); then
            echo "$version is newer than $version_prev"
            return 1
        fi
    done
    echo "$version is equal to $version_prev"
    return 0
}


mkdir -p build/ cache/

# Setup environment

UPX_VERSION="5.0.1"
wget -q https://github.com/upx/upx/releases/download/v$UPX_VERSION/upx-$UPX_VERSION-amd64_linux.tar.xz -O cache/upx.tar.xz
tar -xf cache/upx.tar.xz -C cache/
export PATH="cache/upx-$UPX_VERSION-amd64_linux/:$PATH"

# Get version string

go build -o cache/roll
VERSION=`cache/roll version`

build() {
    local GOOS="$1"
    local ARCH="$2"

    local GOARCH
    case "$ARCH" in
        x86_64)      GOARCH="amd64"  ;;
        i386|i686)   GOARCH="386"    ;;
        aarch64)     GOARCH="arm64"  ;;
        armv7)       GOARCH="arm"    ;;
        loongarch64) GOARCH="loong64";;
        *)           GOARCH="$ARCH"  ;;
    esac

    local EXT=""
    if [[ "$GOOS" == "windows" ]]; then
        EXT=".exe"
    fi

    local FNAME="roll_${VERSION}_${GOOS}_${ARCH}${EXT}"
    local FPATH="build/$FNAME"
    echo "Building '$FNAME'"
    GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build \
        -gcflags=all="-B" -ldflags="-w -s"  -o "$FPATH"

    upx -q --best --ultra-brute "$FPATH" 1>/dev/null
    echo "Building '$FNAME' done"
}

# Build binaries

for target in "${TARGETS[@]}"; do
    build $target &
done
wait

rm -rf cache/
