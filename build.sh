#!/bin/bash

# Excel to JSON 跨平台构建脚本
# 用于编译 macOS、Windows、Linux 版本

set -e

echo "=========================================="
echo "Excel2JSON 构建脚本"
echo "=========================================="

# 获取用户输入的版本号
read -p "请输入版本号 (例如: 1.0.0): " VERSION

# 验证版本号不为空
if [ -z "$VERSION" ]; then
    echo "❌ 错误: 版本号不能为空"
    exit 1
fi

# 清理 dist 目录中的旧版本
if [ -d "./dist" ] && [ "$(ls -A ./dist 2>/dev/null)" ]; then
    echo ""
    echo "📦 dist 目录中已有以下文件："
    ls -1 ./dist/ | while read file; do
        if [ -f "./dist/$file" ]; then
            size=$(ls -lh "./dist/$file" | awk '{print $5}')
            echo "  - $file ($size)"
        fi
    done
    echo ""
    echo "🧹 清理 dist 目录中的旧版本文件..."
    rm -f ./dist/excel2json-*
    echo "✅ dist 目录已清理"
    echo ""
fi

# 询问是否打 git tag
read -p "是否创建 Git Tag (版本号: $VERSION)? [y/N]: " CREATE_TAG
CREATE_TAG=${CREATE_TAG:-N}

BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo ""
echo "版本: $VERSION"
echo "构建时间: $BUILD_TIME"
echo "Git提交: $GIT_COMMIT"
if [ "$CREATE_TAG" = "y" ] || [ "$CREATE_TAG" = "Y" ]; then
    echo "将创建 Git Tag: $VERSION"
fi
echo "=========================================="

# 创建输出目录
mkdir -p ./dist

# 定义链接器标志
LD_FLAGS="-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME' -X main.gitCommit=$GIT_COMMIT"

# 定义文件名（包含版本号）
BINARY_MACOS_INTEL="excel2json-v${VERSION}-macos-intel"
BINARY_MACOS_ARM64="excel2json-v${VERSION}-macos-arm64"
BINARY_WINDOWS="excel2json-v${VERSION}-windows-amd64.exe"
BINARY_LINUX="excel2json-v${VERSION}-linux-amd64"

# macOS Intel 版本
echo ""
echo "📦 编译 macOS (Intel x86_64)..."
GOOS=darwin GOARCH=amd64 go build \
  -ldflags "$LD_FLAGS" \
  -o ./dist/$BINARY_MACOS_INTEL main.go
echo "✅ 完成: ./dist/$BINARY_MACOS_INTEL"

# macOS ARM64 版本（M1/M2/M3）
echo ""
echo "📦 编译 macOS (Apple Silicon ARM64)..."
GOOS=darwin GOARCH=arm64 go build \
  -ldflags "$LD_FLAGS" \
  -o ./dist/$BINARY_MACOS_ARM64 main.go
echo "✅ 完成: ./dist/$BINARY_MACOS_ARM64"

# Windows 版本
echo ""
echo "📦 编译 Windows (x86_64)..."
GOOS=windows GOARCH=amd64 go build \
  -ldflags "$LD_FLAGS" \
  -o ./dist/$BINARY_WINDOWS main.go
echo "✅ 完成: ./dist/$BINARY_WINDOWS"

# Linux 版本
echo ""
echo "📦 编译 Linux (x86_64)..."
GOOS=linux GOARCH=amd64 go build \
  -ldflags "$LD_FLAGS" \
  -o ./dist/$BINARY_LINUX main.go
echo "✅ 完成: ./dist/$BINARY_LINUX"

echo ""
echo "=========================================="
echo "✅ 构建完成！所有文件已保存到 ./dist/"
echo "=========================================="
echo ""
ls -lh ./dist/

# 创建 Git Tag
if [ "$CREATE_TAG" = "y" ] || [ "$CREATE_TAG" = "Y" ]; then
    echo ""
    echo "=========================================="
    echo "🏷️  创建 Git Tag: $VERSION"
    echo "=========================================="
    
    # 检查 tag 是否已存在
    if git rev-parse "$VERSION" >/dev/null 2>&1; then
        echo "⚠️  警告: Tag '$VERSION' 已存在"
        read -p "是否覆盖现有 Tag? [y/N]: " OVERWRITE_TAG
        OVERWRITE_TAG=${OVERWRITE_TAG:-N}
        
        if [ "$OVERWRITE_TAG" = "y" ] || [ "$OVERWRITE_TAG" = "Y" ]; then
            git tag -d "$VERSION" 2>/dev/null || true
            git push origin ":refs/tags/$VERSION" 2>/dev/null || true
            git tag "$VERSION"
            echo "✅ 已覆盖并创建 Tag: $VERSION"
        else
            echo "⏭️  跳过创建 Tag"
        fi
    else
        git tag "$VERSION"
        echo "✅ 已创建 Tag: $VERSION"
    fi
    
    # 询问是否推送 tag
    read -p "是否推送 Tag 到远程仓库? [y/N]: " PUSH_TAG
    PUSH_TAG=${PUSH_TAG:-N}
    
    if [ "$PUSH_TAG" = "y" ] || [ "$PUSH_TAG" = "Y" ]; then
        git push origin "$VERSION"
        echo "✅ 已推送 Tag 到远程仓库"
    fi
fi

echo ""
echo "📝 文件说明："
echo "  - $BINARY_MACOS_INTEL   : macOS Intel 版本"
echo "  - $BINARY_MACOS_ARM64   : macOS Apple Silicon 版本"
echo "  - $BINARY_WINDOWS       : Windows 版本"
echo "  - $BINARY_LINUX         : Linux 版本"
echo ""
echo "🚀 使用示例："
echo "  ./dist/$BINARY_MACOS_ARM64 -file=./assets/abc.xlsx -output_path=./output/"
echo "  ./dist/$BINARY_WINDOWS -file=./assets/abc.xlsx -output_path=./output/"
echo ""
