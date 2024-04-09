#!/bin/bash

# 获取当前目录的绝对路径
ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

# 定义输出目录的基础部分
OUTPUT_BASE="${ROOT_DIR}/_output/bin/platforms"

compile_for_platform() {
    local platform=$1
    local os=${platform%_*}
    local arch=${platform#*_}

    # 最终输出目录
    OUTPUT_DIR="${OUTPUT_BASE}/${os}/${arch}"

    # 创建输出目录（如果不存在）
    mkdir -p "$OUTPUT_DIR"

    compile_dir() {
        for dir in $1/*; do
            if [ -d "$dir" ] && [ -f "$dir/main.go" ]; then
                # 获取目录名
                DIR_NAME=$(basename "$dir")
                echo "Compiling $DIR_NAME for $platform..."

                # 编译生成二进制文件
                GOOS=$os GOARCH=$arch go build -o "${OUTPUT_DIR}/${DIR_NAME}" "$dir/main.go"
                if [ $? -ne 0 ]; then
                    echo "Failed to compile $DIR_NAME for $platform"
                    exit 1
                fi
            fi
        done
    }

    echo "Compiling cmd for $platform..."
    compile_dir "${ROOT_DIR}/cmd"

    echo "Compiling tools for $platform..."
    compile_dir "${ROOT_DIR}/tools"
}

# 检测并编译指定的平台
if [ -z "$PLATFORMS" ]; then
    # 默认使用当前系统的平台
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64)
            ARCH="arm64"
            ;;
        *)
            echo "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    compile_for_platform "${OS}_${ARCH}"
else
    # 遍历并编译每个指定的平台
    IFS=' ' read -ra ADDR <<< "$PLATFORMS"
    for platform in "${ADDR[@]}"; do
        compile_for_platform $platform
    done
fi

echo "Compilation complete. Binaries are located in ${OUTPUT_BASE}"
